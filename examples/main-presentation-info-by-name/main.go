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

	"github.com/brightdevelopers/gopurple"
)

func main() {
	var (
		helpFlag    = flag.Bool("help", false, "Display usage information")
		verboseFlag = flag.Bool("verbose", false, "Show detailed information")
		jsonFlag    = flag.Bool("json", false, "Output as JSON")
		timeoutFlag = flag.Int("timeout", 30, "Request timeout in seconds")
		networkFlag *string
		nameFlag    = flag.String("name", "", "Presentation name to retrieve (required)")
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A command-line tool for retrieving presentation details by name from BSN.cloud.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Get presentation info by name:\n")
		fmt.Fprintf(os.Stderr, "    %s --name \"My Presentation\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Verbose output:\n")
		fmt.Fprintf(os.Stderr, "    %s --name \"My Presentation\" --verbose\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Output as JSON:\n")
		fmt.Fprintf(os.Stderr, "    %s --name \"My Presentation\" --json\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	// Validate required arguments
	if *nameFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: --name must be specified\n\n")
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

	if !*jsonFlag {
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
	if !*jsonFlag {
		fmt.Println("Authenticating with BSN.cloud...")
	}

	if err := client.Authenticate(ctx); err != nil {
		if gopurple.IsAuthenticationError(err) {
			log.Fatalf("Authentication failed: %v", err)
		}
		log.Fatalf("Authentication error: %v", err)
	}

	if !*jsonFlag {
		fmt.Println("Authentication successful!")
	}

	// Handle network selection
	if err := handleNetworkSelection(ctx, client, *networkFlag, *jsonFlag); err != nil {
		log.Fatalf("Network selection failed: %v", err)
	}

	// Get presentation info
	if err := getPresentationInfo(ctx, client, *nameFlag, *jsonFlag, *verboseFlag); err != nil {
		log.Fatalf("Failed to get presentation info: %v", err)
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
	fmt.Fprint(os.Stderr, "Select network (1-"+strconv.Itoa(len(networks))+"): ")
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

func getPresentationInfo(ctx context.Context, client *gopurple.Client, name string, jsonMode bool, verbose bool) error {
	if !jsonMode {
		fmt.Printf("\nRetrieving presentation '%s'...\n", name)
	}

	result, err := client.Presentations.GetByName(ctx, name)
	if err != nil {
		return err
	}

	if jsonMode {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	// Display presentation information
	fmt.Println("\n=== Presentation Details ===")
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

	fmt.Printf("Created: %s\n", result.CreationDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("Last Modified: %s\n", result.LastModifiedDate.Format("2006-01-02 15:04:05"))

	if result.LastPublishedDate != nil {
		fmt.Printf("Last Published: %s\n", result.LastPublishedDate.Format("2006-01-02 15:04:05"))
	}

	if verbose {
		fmt.Println("\n=== Additional Details ===")
		fmt.Printf("Is Simple Playlist: %t\n", result.IsSimplePlaylist)
		fmt.Printf("Zone Count: %d\n", result.ZoneCount)
		fmt.Printf("Content Item Count: %d\n", result.ContentItemCount)

		if result.ThumbnailURL != "" {
			fmt.Printf("Thumbnail URL: %s\n", result.ThumbnailURL)
		}

		if len(result.Tags) > 0 {
			fmt.Printf("Tags: %v\n", result.Tags)
		}

		if result.ScheduleSettings != nil {
			fmt.Println("\n=== Schedule Settings ===")
			if result.ScheduleSettings.StartDate != nil {
				fmt.Printf("Start Date: %s\n", result.ScheduleSettings.StartDate.Format("2006-01-02 15:04:05"))
			}
			if result.ScheduleSettings.EndDate != nil {
				fmt.Printf("End Date: %s\n", result.ScheduleSettings.EndDate.Format("2006-01-02 15:04:05"))
			}
			if result.ScheduleSettings.RecurrenceRule != "" {
				fmt.Printf("Recurrence Rule: %s\n", result.ScheduleSettings.RecurrenceRule)
			}
		}
	}

	return nil
}
