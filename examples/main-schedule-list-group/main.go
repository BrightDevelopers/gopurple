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
		jsonFlag    = flag.Bool("json", false, "Output as JSON")
		verboseFlag = flag.Bool("verbose", false, "Show detailed information")
		timeoutFlag = flag.Int("timeout", 30, "Request timeout in seconds")
		networkFlag *string
		groupFlag   = flag.String("group", "", "Group name to list schedules for (required)")
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A command-line tool for listing scheduled presentations for a BSN.cloud group.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  List all schedules for a group:\n")
		fmt.Fprintf(os.Stderr, "    %s --group \"Lobby Displays\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Output as JSON:\n")
		fmt.Fprintf(os.Stderr, "    %s --group \"Lobby Displays\" --json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Show detailed information:\n")
		fmt.Fprintf(os.Stderr, "    %s --group \"Lobby Displays\" --verbose\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	// Validate required arguments
	if *groupFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: --group must be specified\n\n")
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

	if *verboseFlag && !*jsonFlag {
		fmt.Fprintf(os.Stderr, "Creating BSN.cloud client...\n")
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
	if *verboseFlag && !*jsonFlag {
		fmt.Fprintf(os.Stderr, "Authenticating with BSN.cloud...\n")
	}

	if err := client.Authenticate(ctx); err != nil {
		if gopurple.IsAuthenticationError(err) {
			log.Fatalf("Authentication failed: %v", err)
		}
		log.Fatalf("Authentication error: %v", err)
	}

	if *verboseFlag && !*jsonFlag {
		fmt.Fprintf(os.Stderr, "Authentication successful!\n")
	}

	// Handle network selection
	if err := handleNetworkSelection(ctx, client, *networkFlag, *verboseFlag, *jsonFlag); err != nil {
		log.Fatalf("Network selection failed: %v", err)
	}

	// List schedules
	if err := listGroupSchedules(ctx, client, *groupFlag, *verboseFlag, *jsonFlag); err != nil {
		log.Fatalf("Failed to list schedules: %v", err)
	}
}

func handleNetworkSelection(ctx context.Context, client *gopurple.Client, requestedNetwork string, verbose, jsonMode bool) error {
	// Check if network is already set
	if client.IsNetworkSet() {
		if current, err := client.GetCurrentNetwork(ctx); err == nil {
			if verbose && !jsonMode {
				fmt.Fprintf(os.Stderr, "Using network: %s (ID: %d)\n", current.Name, current.ID)
			}
			return nil
		}
	}

	// If no network flag was provided, check BS_NETWORK environment variable
	if requestedNetwork == "" {
		if envNetwork := os.Getenv("BS_NETWORK"); envNetwork != "" {
			requestedNetwork = envNetwork
			if !jsonMode {
				fmt.Fprintf(os.Stderr, "Using network from BS_NETWORK environment variable\n")
			}
		}
	}

	// Get available networks
	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "Getting available networks...\n")
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
				if verbose && !jsonMode {
					fmt.Fprintf(os.Stderr, "Using requested network: %s (ID: %d)\n", network.Name, network.ID)
				}
				return client.SetNetworkByID(ctx, network.ID)
			}
		}

		// Network not found - show error and fall back to interactive selection
		if !jsonMode {
			fmt.Fprintf(os.Stderr, "Network '%s' not found. Available networks:\n", requestedNetwork)
			for i, network := range networks {
				fmt.Fprintf(os.Stderr, "  %d. %s (ID: %d)\n", i+1, network.Name, network.ID)
			}
			fmt.Fprintf(os.Stderr, "\n")
		} else {
			return fmt.Errorf("network '%s' not found", requestedNetwork)
		}
	}

	// If only one network and no specific network requested, use it automatically
	if len(networks) == 1 && requestedNetwork == "" {
		if verbose && !jsonMode {
			fmt.Fprintf(os.Stderr, "Using network: %s (ID: %d)\n", networks[0].Name, networks[0].ID)
		}
		return client.SetNetworkByID(ctx, networks[0].ID)
	}

	// In JSON mode, cannot do interactive selection
	if jsonMode {
		return fmt.Errorf("network selection required: use --network flag or BS_NETWORK environment variable")
	}

	// Show available networks and let user choose
	if requestedNetwork == "" {
		fmt.Fprintf(os.Stderr, "Available networks:\n")
		for i, network := range networks {
			fmt.Fprintf(os.Stderr, "  %d. %s (ID: %d)\n", i+1, network.Name, network.ID)
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
	if verbose {
		fmt.Fprintf(os.Stderr, "Selected network: %s (ID: %d)\n", selectedNetwork.Name, selectedNetwork.ID)
	}

	return client.SetNetworkByID(ctx, selectedNetwork.ID)
}

func listGroupSchedules(ctx context.Context, client *gopurple.Client, groupName string, verbose, jsonMode bool) error {
	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "\nRetrieving schedules for group '%s'...\n", groupName)
	}

	result, err := client.Schedules.GetGroupSchedule(ctx, groupName)
	if err != nil {
		return fmt.Errorf("failed to get schedules: %w", err)
	}

	if jsonMode {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	// Display schedules
	fmt.Fprintf(os.Stderr, "\nScheduled Presentations for Group: %s\n", groupName)
	fmt.Fprintf(os.Stderr, "%s\n", strings.Repeat("=", 80))

	if len(result.Items) == 0 {
		fmt.Fprintf(os.Stderr, "No scheduled presentations found.\n")
		return nil
	}

	for i, schedule := range result.Items {
		scheduleType := "One-time"
		if schedule.IsRecurrent {
			scheduleType = "Recurring"
		}

		fmt.Fprintf(os.Stderr, "\n[%d] Schedule ID: %d (%s)\n", i+1, schedule.ID, scheduleType)
		fmt.Fprintf(os.Stderr, "    Presentation: %s (ID: %d)\n", schedule.PresentationName, schedule.PresentationID)
		fmt.Fprintf(os.Stderr, "    Start Time: %s\n", schedule.StartTime)
		fmt.Fprintf(os.Stderr, "    Duration: %s\n", schedule.Duration)

		if schedule.IsRecurrent {
			if schedule.RecurrenceStartDate != nil {
				fmt.Fprintf(os.Stderr, "    Recurrence Start: %s\n", schedule.RecurrenceStartDate.Format("2006-01-02"))
			}
			if schedule.RecurrenceEndDate != nil {
				fmt.Fprintf(os.Stderr, "    Recurrence End: %s\n", schedule.RecurrenceEndDate.Format("2006-01-02"))
			}
			if schedule.DaysOfWeek != nil {
				fmt.Fprintf(os.Stderr, "    Days of Week: %v\n", schedule.DaysOfWeek)
			}
		} else {
			if schedule.EventDate != nil {
				fmt.Fprintf(os.Stderr, "    Event Date: %s\n", schedule.EventDate.Format("2006-01-02"))
			}
		}

		if verbose {
			if schedule.InterruptScheduling {
				fmt.Fprintf(os.Stderr, "    Interrupt Scheduling: Yes\n")
			}
			if schedule.ExpirationDate != nil {
				fmt.Fprintf(os.Stderr, "    Expiration Date: %s\n", schedule.ExpirationDate.Format("2006-01-02"))
			}
		}
	}

	fmt.Fprintf(os.Stderr, "\nTotal schedules: %d\n", len(result.Items))

	if result.IsTruncated {
		fmt.Fprintf(os.Stderr, "Note: Results are truncated. Total count: %d\n", result.TotalCount)
	}

	return nil
}
