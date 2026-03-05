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
		helpFlag       = flag.Bool("help", false, "Display usage information")
		jsonFlag       = flag.Bool("json", false, "Output as JSON")
		verboseFlag    = flag.Bool("verbose", false, "Show detailed information")
		timeoutFlag    = flag.Int("timeout", 30, "Request timeout in seconds")
		networkFlag    *string
		groupFlag      = flag.String("group", "", "Group name containing the schedule (required)")
		scheduleIDFlag = flag.Int("schedule-id", 0, "Schedule ID to update (required)")
		startTimeFlag  = flag.String("start-time", "", "New start time in HH:MM:SS format (optional)")
		durationFlag   = flag.String("duration", "", "New duration in HH:MM:SS format (optional)")
		daysFlag       = flag.String("days", "", "New days of week for recurring schedules (comma-separated, optional)")
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A command-line tool for updating an existing scheduled presentation.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Update start time:\n")
		fmt.Fprintf(os.Stderr, "    %s --group \"Lobby Displays\" --schedule-id 456 --start-time 14:00:00\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Update duration:\n")
		fmt.Fprintf(os.Stderr, "    %s --group \"Lobby Displays\" --schedule-id 456 --duration 02:00:00\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Update multiple fields:\n")
		fmt.Fprintf(os.Stderr, "    %s --group \"Lobby Displays\" --schedule-id 456 --start-time 14:00:00 --duration 02:00:00 --days \"Monday,Wednesday,Friday\"\n", os.Args[0])
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
	if *scheduleIDFlag == 0 {
		fmt.Fprintf(os.Stderr, "Error: --schedule-id must be specified\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Check that at least one update field is specified
	if *startTimeFlag == "" && *durationFlag == "" && *daysFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: at least one of --start-time, --duration, or --days must be specified\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Validate time formats if provided
	if *startTimeFlag != "" && !isValidTimeFormat(*startTimeFlag) {
		log.Fatalf("Invalid start-time format (use HH:MM:SS): %s", *startTimeFlag)
	}
	if *durationFlag != "" && !isValidTimeFormat(*durationFlag) {
		log.Fatalf("Invalid duration format (use HH:MM:SS): %s", *durationFlag)
	}

	// Validate days of week if provided
	if *daysFlag != "" {
		validDays := map[string]bool{
			"Monday": true, "Tuesday": true, "Wednesday": true, "Thursday": true,
			"Friday": true, "Saturday": true, "Sunday": true, "EveryDay": true,
		}
		daysList := strings.Split(*daysFlag, ",")
		for _, day := range daysList {
			day = strings.TrimSpace(day)
			if !validDays[day] {
				log.Fatalf("Invalid day of week: %s (use Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday, or EveryDay)", day)
			}
		}
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

	// Update schedule
	before, after, err := updateSchedule(ctx, client, *groupFlag, *scheduleIDFlag, *startTimeFlag, *durationFlag, *daysFlag, *verboseFlag, *jsonFlag)
	if err != nil {
		log.Fatalf("Failed to update schedule: %v", err)
	}

	// Output as JSON if requested
	if *jsonFlag {
		result := map[string]interface{}{
			"success":    true,
			"scheduleID": *scheduleIDFlag,
			"before":     before,
			"after":      after,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			log.Fatalf("Failed to encode JSON: %v", err)
		}
		return
	}

	fmt.Fprintf(os.Stderr, "\nSchedule %d updated successfully in group '%s'\n", *scheduleIDFlag, *groupFlag)
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

func updateSchedule(ctx context.Context, client *gopurple.Client, groupName string, scheduleID int, newStartTime, newDuration, newDays string, verbose bool, jsonMode bool) (*gopurple.ScheduledPresentation, *gopurple.ScheduledPresentation, error) {
	// First, get the current schedule
	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "\nRetrieving current schedule %d from group '%s'...\n", scheduleID, groupName)
	}

	currentSchedule, err := client.Schedules.GetScheduledPresentation(ctx, groupName, scheduleID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current schedule: %w", err)
	}

	// Show current schedule
	if !jsonMode {
		scheduleType := "One-time"
		if currentSchedule.IsRecurrent {
			scheduleType = "Recurring"
		}

		fmt.Fprintf(os.Stderr, "\n=== Current Schedule ===\n")
		fmt.Fprintf(os.Stderr, "Schedule ID: %d (%s)\n", currentSchedule.ID, scheduleType)
		fmt.Fprintf(os.Stderr, "Presentation: %s (ID: %d)\n", currentSchedule.PresentationName, currentSchedule.PresentationID)
		fmt.Fprintf(os.Stderr, "Start Time: %s\n", currentSchedule.StartTime)
		fmt.Fprintf(os.Stderr, "Duration: %s\n", currentSchedule.Duration)
		if currentSchedule.IsRecurrent && currentSchedule.DaysOfWeek != nil {
			fmt.Fprintf(os.Stderr, "Days of Week: %v\n", currentSchedule.DaysOfWeek)
		}
	}

	// Apply updates to the schedule object
	updatedSchedule := currentSchedule

	if newStartTime != "" {
		updatedSchedule.StartTime = newStartTime
	}
	if newDuration != "" {
		updatedSchedule.Duration = newDuration
	}
	if newDays != "" {
		updatedSchedule.DaysOfWeek = newDays
	}

	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "\n=== Changes to Apply ===\n")
		if newStartTime != "" {
			fmt.Fprintf(os.Stderr, "Start Time: %s -> %s\n", currentSchedule.StartTime, newStartTime)
		}
		if newDuration != "" {
			fmt.Fprintf(os.Stderr, "Duration: %s -> %s\n", currentSchedule.Duration, newDuration)
		}
		if newDays != "" {
			fmt.Fprintf(os.Stderr, "Days of Week: %v -> %s\n", currentSchedule.DaysOfWeek, newDays)
		}
	}

	// Update the schedule
	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "\nUpdating schedule...\n")
	}

	result, err := client.Schedules.UpdateScheduledPresentation(ctx, groupName, scheduleID, updatedSchedule)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	// Show updated schedule
	if !jsonMode {
		fmt.Fprintf(os.Stderr, "\n=== Updated Schedule ===\n")
		fmt.Fprintf(os.Stderr, "Schedule ID: %d\n", result.ID)
		fmt.Fprintf(os.Stderr, "Presentation: %s (ID: %d)\n", result.PresentationName, result.PresentationID)
		fmt.Fprintf(os.Stderr, "Start Time: %s\n", result.StartTime)
		fmt.Fprintf(os.Stderr, "Duration: %s\n", result.Duration)
		if result.IsRecurrent && result.DaysOfWeek != nil {
			fmt.Fprintf(os.Stderr, "Days of Week: %v\n", result.DaysOfWeek)
		}
	}

	return currentSchedule, result, nil
}

func isValidTimeFormat(timeStr string) bool {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return false
	}
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return false
		}
	}
	return true
}
