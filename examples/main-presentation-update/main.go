package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brightsign/gopurple"
)

func main() {
	var (
		helpFlag        = flag.Bool("help", false, "Display usage information")
		verboseFlag     = flag.Bool("verbose", false, "Show detailed information")
		jsonFlag        = flag.Bool("json", false, "Output as JSON")
		timeoutFlag     = flag.Int("timeout", 30, "Request timeout in seconds")
		networkFlag     *string
		idFlag          = flag.Int("id", 0, "Presentation ID to update (required)")
		nameFlag        = flag.String("name", "", "New presentation name")
		descriptionFlag = flag.String("description", "", "New presentation description")
		tagsFlag        = flag.String("tags", "", "Comma-separated list of tags")
		modelFlag       = flag.String("model", "", "Player model (e.g., HD1024, XT1144, XD1034)")
		languageFlag    = flag.String("language", "", "Language (default: English)")
		statusFlag      = flag.String("status", "", "Status (Draft or Published)")
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A command-line tool for updating presentations on BSN.cloud.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Update presentation name:\n")
		fmt.Fprintf(os.Stderr, "    %s --id 12345 --name \"New Name\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Update with description and tags:\n")
		fmt.Fprintf(os.Stderr, "    %s --id 12345 --description \"Updated demo\" --tags \"sales,demo\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Change player model:\n")
		fmt.Fprintf(os.Stderr, "    %s --id 12345 --model XT1144\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	// Validate required arguments
	if *idFlag == 0 {
		fmt.Fprintf(os.Stderr, "Error: --id must be specified\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Check if at least one update flag is provided
	if *nameFlag == "" && *descriptionFlag == "" && *tagsFlag == "" && *modelFlag == "" && *languageFlag == "" && *statusFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: at least one update option must be specified (--name, --description, --tags, --model, --language, or --status)\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Create client
	var opts []gopurple.Option
	if *timeoutFlag > 0 {
		opts = append(opts, gopurple.WithTimeout(time.Duration(*timeoutFlag)*time.Second))
	}

	// Add network if specified
	if *networkFlag != "" {
		opts = append(opts, gopurple.WithNetwork(*networkFlag))
	}

	if !*jsonFlag && *verboseFlag {
		fmt.Println("Creating BSN.cloud client...")
	}

	client, err := gopurple.New(opts...)
	if err != nil {
		if gopurple.IsConfigurationError(err) {
			log.Fatalf("Configuration error: %v", err)
		}
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Authenticate
	if !*jsonFlag && *verboseFlag {
		fmt.Println("Authenticating with BSN.cloud...")
	}

	if err := client.Authenticate(ctx); err != nil {
		if gopurple.IsAuthenticationError(err) {
			log.Fatalf("Authentication failed: %v", err)
		}
		log.Fatalf("Authentication error: %v", err)
	}

	if !*jsonFlag && *verboseFlag {
		fmt.Println("Authentication successful!")
	}

	// Handle network selection
	if err := handleNetworkSelection(ctx, client, *networkFlag, *jsonFlag && !*verboseFlag); err != nil {
		log.Fatalf("Network selection failed: %v", err)
	}

	// Parse tags
	var tags []string
	if *tagsFlag != "" {
		tags = strings.Split(*tagsFlag, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	// Update presentation
	if err := updatePresentation(ctx, client, *idFlag, *nameFlag, *descriptionFlag, tags, *modelFlag, *languageFlag, *statusFlag, *jsonFlag, *verboseFlag); err != nil {
		log.Fatalf("Failed to update presentation: %v", err)
	}
}

func handleNetworkSelection(ctx context.Context, client *gopurple.Client, requestedNetwork string, jsonMode bool) error {
	// Check if network is already set
	if client.IsNetworkSet() {
		if current, err := client.GetCurrentNetwork(ctx); err == nil {
			if !jsonMode {
				fmt.Printf("Using network: %s (ID: %d)\n", current.Name, current.ID)
			}
			return nil
		}
	}

	// If no network flag was provided, check BS_NETWORK environment variable
	if requestedNetwork == "" {
		if envNetwork := os.Getenv("BS_NETWORK"); envNetwork != "" {
			requestedNetwork = envNetwork
			if !jsonMode {
				fmt.Printf("Using network from BS_NETWORK environment variable\n")
			}
		}
	}

	// Get available networks
	if !jsonMode {
		fmt.Println("Getting available networks...")
	}

	networks, err := client.GetNetworks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get networks: %w", err)
	}

	if len(networks) == 0 {
		return fmt.Errorf("no networks available")
	}

	// If a specific network was requested, try to find it
	if requestedNetwork != "" {
		for _, network := range networks {
			if strings.EqualFold(network.Name, requestedNetwork) {
				if !jsonMode {
					fmt.Printf("Using requested network: %s (ID: %d)\n", network.Name, network.ID)
				}
				return client.SetNetworkByID(ctx, network.ID)
			}
		}

		// Network not found - show error and fall back to interactive selection
		fmt.Printf("Network '%s' not found. Available networks:\n", requestedNetwork)
		for i, network := range networks {
			fmt.Printf("  %d. %s (ID: %d)\n", i+1, network.Name, network.ID)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	// If only one network and no specific network requested, use it automatically
	if len(networks) == 1 && requestedNetwork == "" {
		if !jsonMode {
			fmt.Printf("Using network: %s (ID: %d)\n", networks[0].Name, networks[0].ID)
		}
		return client.SetNetworkByID(ctx, networks[0].ID)
	}

	// Show available networks and let user choose
	if requestedNetwork == "" {
		fmt.Fprintf(os.Stderr, "Available networks:\n")
		for i, network := range networks {
			fmt.Printf("  %d. %s (ID: %d)\n", i+1, network.Name, network.ID)
		}
	}

	// Get user selection
	fmt.Fprint(os.Stderr, "Select network (1-" + strconv.Itoa(len(networks)) + "): ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return fmt.Errorf("failed to read input")
	}

	selection, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || selection < 1 || selection > len(networks) {
		return fmt.Errorf("invalid selection: must be between 1 and %d", len(networks))
	}

	selectedNetwork := networks[selection-1]
	if !jsonMode {
		if !jsonMode {
		fmt.Fprintf(os.Stderr, "Selected network: %s (ID: %d)\n", selectedNetwork.Name, selectedNetwork.ID)
	}
	}

	return client.SetNetworkByID(ctx, selectedNetwork.ID)
}

func updatePresentation(ctx context.Context, client *gopurple.Client, id int, name, description string, tags []string, model, language, status string, jsonMode bool, verbose bool) error {
	// First, get the current presentation
	if !jsonMode && verbose {
		fmt.Printf("\nRetrieving current presentation %d...\n", id)
	}

	current, err := client.Presentations.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !jsonMode {
		fmt.Println("\n=== Current Presentation ===")
		fmt.Printf("ID: %d\n", current.ID)
		fmt.Printf("Name: %s\n", current.Name)
		if current.Type != "" {
			fmt.Printf("Type: %s\n", current.Type)
		}
		if current.Description != "" {
			fmt.Printf("Description: %s\n", current.Description)
		}
	}

	// Build the update request based on current state
	// Start with current values, then override with new values if provided
	request := &gopurple.PresentationCreateRequest{
		ID:                    current.ID,
		Name:                  current.Name,
		CreationDate:          current.CreationDate.Format(time.RFC3339),
		LastModifiedDate:      current.LastModifiedDate.Format(time.RFC3339),
		ProjectFile:           current.ProjectFile,
		DeviceModel:           current.DeviceModel,
		Language:              current.Language,
		Status:                current.Status,
		AutoplayFile:          current.AutoplayFile,
		ResourcesFile:         current.ResourcesFile,
		UserDefinedEventsFile: current.UserDefinedEventsFile,
		ThumbnailFile:         current.ThumbnailFile,
		Files:                 current.Files,
		AutorunPlugins:        current.AutorunPlugins,
		Applications:          current.Applications,
		Dependencies:          current.Dependencies,
		Groups:                current.Groups,
		Permissions:           current.Permissions,
		Tags:                  current.Tags,
	}

	// Copy nested structures if they exist
	if current.Autorun != nil {
		request.Autorun = &gopurple.PresentationAutorun{
			Version:  current.Autorun.Version,
			IsCustom: current.Autorun.IsCustom,
		}
	}

	if current.DeviceWebPage != nil {
		request.DeviceWebPage = &gopurple.PresentationDeviceWebPage{
			ID:   current.DeviceWebPage.ID,
			Name: current.DeviceWebPage.Name,
		}
	}

	if current.ScreenSettings != nil {
		request.ScreenSettings = &gopurple.PresentationScreenSettings{
			VideoMode:       current.ScreenSettings.VideoMode,
			Orientation:     current.ScreenSettings.Orientation,
			Connector:       current.ScreenSettings.Connector,
			BackgroundColor: current.ScreenSettings.BackgroundColor,
			Overscan:        current.ScreenSettings.Overscan,
		}
	}

	// Apply updates
	changesMade := false
	if name != "" && name != current.Name {
		request.Name = name
		changesMade = true
		if !jsonMode {
			fmt.Printf("\n→ Changing name: '%s' → '%s'\n", current.Name, name)
		}
	}

	if description != "" {
		changesMade = true
		if !jsonMode {
			fmt.Printf("→ Setting description: '%s'\n", description)
		}
	}

	if len(tags) > 0 {
		request.Tags = tags
		changesMade = true
		if !jsonMode {
			fmt.Printf("→ Setting tags: %v\n", tags)
		}
	}

	if model != "" && model != current.DeviceModel {
		request.DeviceModel = model
		changesMade = true
		if !jsonMode {
			fmt.Printf("→ Changing device model: '%s' → '%s'\n", current.DeviceModel, model)
		}
	}

	if language != "" && language != current.Language {
		request.Language = language
		changesMade = true
		if !jsonMode {
			fmt.Printf("→ Changing language: '%s' → '%s'\n", current.Language, language)
		}
	}

	if status != "" && status != current.Status {
		request.Status = status
		changesMade = true
		if !jsonMode {
			fmt.Printf("→ Changing status: '%s' → '%s'\n", current.Status, status)
		}
	}

	if !changesMade {
		if !jsonMode {
			fmt.Println("\nNo changes detected (all values are the same as current values).")
		}
		return nil
	}

	// Perform the update
	if !jsonMode && verbose {
		fmt.Printf("\nUpdating presentation %d...\n", id)
	}

	result, err := client.Presentations.Update(ctx, id, request)
	if err != nil {
		return err
	}

	if jsonMode {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	// Display updated presentation information
	fmt.Println("\n=== Presentation Updated Successfully ===")
	fmt.Printf("ID: %d\n", result.ID)
	fmt.Printf("Name: %s\n", result.Name)

	if result.Type != "" {
		fmt.Printf("Type: %s\n", result.Type)
	}

	if result.Description != "" {
		fmt.Printf("Description: %s\n", result.Description)
	}

	if result.PublishState != "" {
		fmt.Printf("Publish State: %s\n", result.PublishState)
	}

	fmt.Printf("Last Modified: %s\n", result.LastModifiedDate.Format("2006-01-02 15:04:05"))

	if verbose {
		fmt.Println("\n=== Additional Details ===")
		if result.DeviceModel != "" {
			fmt.Printf("Device Model: %s\n", result.DeviceModel)
		}
		if result.Language != "" {
			fmt.Printf("Language: %s\n", result.Language)
		}
		if result.Status != "" {
			fmt.Printf("Status: %s\n", result.Status)
		}
		if len(result.Tags) > 0 {
			fmt.Printf("Tags: %v\n", result.Tags)
		}
	}

	return nil
}
