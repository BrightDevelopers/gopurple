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
		helpFlag    = flag.Bool("help", false, "Display usage information")
		jsonFlag    = flag.Bool("json", false, "Output as JSON")
		timeoutFlag = flag.Int("timeout", 30, "Request timeout in seconds")
		networkFlag *string
		filterFlag  = flag.String("filter", "", "Filter expression to select presentations to delete (required)")
		confirmFlag = flag.Bool("yes", false, "Skip confirmation prompt")
		dryRunFlag  = flag.Bool("dry-run", false, "Preview presentations that would be deleted without actually deleting them")
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A command-line tool for deleting presentations from a BSN.cloud network using filters.\n\n")
		fmt.Fprintf(os.Stderr, "WARNING: This operation is destructive and cannot be undone!\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Preview presentations that would be deleted:\n")
		fmt.Fprintf(os.Stderr, "    %s --filter \"name contains 'test'\" --dry-run\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Delete presentations with confirmation:\n")
		fmt.Fprintf(os.Stderr, "    %s --filter \"publishState eq 'Draft'\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Delete presentations without confirmation:\n")
		fmt.Fprintf(os.Stderr, "    %s --filter \"name startsWith 'temp_'\" --yes\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	// Validate required flags
	if *filterFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: --filter is required\n\n")
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

	// Delete presentations
	if err := deletePresentations(ctx, client, *filterFlag, *jsonFlag, *confirmFlag, *dryRunFlag); err != nil {
		log.Fatalf("Failed to delete presentations: %v", err)
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

func deletePresentations(ctx context.Context, client *gopurple.Client, filter string, jsonMode bool, skipConfirm bool, dryRun bool) error {
	// First, get the list of presentations that match the filter to show what we're deleting
	if !jsonMode {
		fmt.Printf("\nFinding presentations matching filter: %s\n", filter)
	}

	listOpts := []gopurple.ListOption{
		gopurple.WithFilter(filter),
		gopurple.WithPageSize(100),
	}

	result, err := client.Presentations.List(ctx, listOpts...)
	if err != nil {
		return fmt.Errorf("failed to list presentations: %w", err)
	}

	if len(result.Items) == 0 {
		fmt.Println("\nNo presentations match the specified filter.")
		return nil
	}

	// Show what we're about to delete
	fmt.Printf("\n=== Presentations Matching Filter (%d found) ===\n", len(result.Items))
	for i, pres := range result.Items {
		fmt.Printf("%d. ID: %d, Name: %s", i+1, pres.ID, pres.Name)
		if pres.Type != "" {
			fmt.Printf(", Type: %s", pres.Type)
		}
		if pres.PublishState != "" {
			fmt.Printf(", State: %s", pres.PublishState)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	// If dry-run, stop here
	if dryRun {
		fmt.Println("\nDRY RUN: No presentations were deleted.")
		return nil
	}

	// Confirm deletion unless --yes flag is used
	if !skipConfirm && !jsonMode {
		fmt.Printf("\nWARNING: This will delete %d presentation(s). This action cannot be undone!\n", len(result.Items))
		fmt.Print("Are you sure you want to continue? (yes/no): ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}

		response := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if response != "yes" && response != "y" {
			return fmt.Errorf("deletion cancelled")
		}
	}

	if !jsonMode {
		fmt.Printf("\nDeleting presentations...\n")
	}

	// Perform the deletion
	deleteResult, err := client.Presentations.DeleteByFilter(ctx, filter)
	if err != nil {
		return err
	}

	if jsonMode {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(deleteResult)
	}

	// Display results
	fmt.Println("\n=== Deletion Complete ===")
	fmt.Printf("Presentations deleted: %d\n", deleteResult.DeletedCount)

	if len(deleteResult.DeletedIDs) > 0 {
		fmt.Printf("Deleted IDs: %v\n", deleteResult.DeletedIDs)
	}

	if len(deleteResult.Errors) > 0 {
		fmt.Println("\nErrors encountered:")
		for _, errMsg := range deleteResult.Errors {
			fmt.Printf("  - %s\n", errMsg)
		}
	}

	return nil
}
