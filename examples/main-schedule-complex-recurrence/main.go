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
		helpFlag          = flag.Bool("help", false, "Display usage information")
		jsonFlag          = flag.Bool("json", false, "Output as JSON")
		verboseFlag       = flag.Bool("verbose", false, "Show detailed information")
		timeoutFlag       = flag.Int("timeout", 30, "Request timeout in seconds")
		networkFlag       *string
		presentationIDFlag = flag.Int("presentation-id", 0, "Presentation ID to schedule (required)")
		groupFlag         = flag.String("group", "", "Group name to add schedule to (required)")
		startDateFlag     = flag.String("start-date", "", "Start date for recurrence (YYYY-MM-DD format) (required)")
		endDateFlag       = flag.String("end-date", "", "End date for recurrence (YYYY-MM-DD format) (required)")
		daysFlag          = flag.String("days", "Saturday,Sunday", "Days of week (comma-separated: Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday)")
		startTimeFlag     = flag.String("start-time", "09:00:00", "Start time in HH:MM:SS format")
		durationFlag      = flag.String("duration", "01:00:00", "Duration in HH:MM:SS format")
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A command-line tool for scheduling presentations with complex recurring patterns.\n")
		fmt.Fprintf(os.Stderr, "Supports custom combinations of days (e.g., weekends only, specific days, etc.).\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Schedule for weekends only:\n")
		fmt.Fprintf(os.Stderr, "    %s --presentation-id 123 --group \"Lobby Displays\" --start-date 2026-04-01 --end-date 2026-04-30 --days \"Saturday,Sunday\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Schedule for Monday, Wednesday, Friday:\n")
		fmt.Fprintf(os.Stderr, "    %s --presentation-id 123 --group \"Lobby Displays\" --start-date 2026-04-01 --end-date 2026-04-30 --days \"Monday,Wednesday,Friday\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Schedule for Tuesday and Thursday with custom time:\n")
		fmt.Fprintf(os.Stderr, "    %s --presentation-id 123 --group \"Lobby Displays\" --start-date 2026-04-01 --end-date 2026-04-30 --days \"Tuesday,Thursday\" --start-time 14:00:00\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	// Validate required arguments
	if *presentationIDFlag == 0 {
		fmt.Fprintf(os.Stderr, "Error: --presentation-id must be specified\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *groupFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: --group must be specified\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *startDateFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: --start-date must be specified\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *endDateFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: --end-date must be specified\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Parse and validate dates
	startDate, err := time.Parse("2006-01-02", *startDateFlag)
	if err != nil {
		log.Fatalf("Invalid start-date format (use YYYY-MM-DD): %v", err)
	}
	endDate, err := time.Parse("2006-01-02", *endDateFlag)
	if err != nil {
		log.Fatalf("Invalid end-date format (use YYYY-MM-DD): %v", err)
	}
	if endDate.Before(startDate) {
		log.Fatalf("End date must be after start date")
	}

	// Validate time formats
	if !isValidTimeFormat(*startTimeFlag) {
		log.Fatalf("Invalid start-time format (use HH:MM:SS): %s", *startTimeFlag)
	}
	if !isValidTimeFormat(*durationFlag) {
		log.Fatalf("Invalid duration format (use HH:MM:SS): %s", *durationFlag)
	}

	// Validate days of week
	validDays := map[string]bool{
		"Monday": true, "Tuesday": true, "Wednesday": true, "Thursday": true,
		"Friday": true, "Saturday": true, "Sunday": true,
	}
	daysList := strings.Split(*daysFlag, ",")
	for _, day := range daysList {
		day = strings.TrimSpace(day)
		if !validDays[day] {
			log.Fatalf("Invalid day of week: %s (use Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday)", day)
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

	// Add complex recurring schedule
	schedule, err := addComplexRecurringSchedule(ctx, client, *presentationIDFlag, *groupFlag, startDate, endDate, *daysFlag, *startTimeFlag, *durationFlag, *verboseFlag, *jsonFlag)
	if err != nil {
		log.Fatalf("Failed to add scheduled presentation: %v", err)
	}

	// Output as JSON if requested
	if *jsonFlag {
		result := map[string]interface{}{
			"success":  true,
			"schedule": schedule,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			log.Fatalf("Failed to encode JSON: %v", err)
		}
		return
	}

	fmt.Fprintf(os.Stderr, "\nComplex recurring schedule created successfully:\n")
	fmt.Fprintf(os.Stderr, "  Schedule ID: %d\n", schedule.ID)
	fmt.Fprintf(os.Stderr, "  Presentation: %s (ID: %d)\n", schedule.PresentationName, schedule.PresentationID)
	fmt.Fprintf(os.Stderr, "  Group: %s\n", *groupFlag)
	fmt.Fprintf(os.Stderr, "  Days: %s\n", *daysFlag)
	fmt.Fprintf(os.Stderr, "  Start Date: %s\n", startDate.Format("2006-01-02"))
	fmt.Fprintf(os.Stderr, "  End Date: %s\n", endDate.Format("2006-01-02"))
	fmt.Fprintf(os.Stderr, "  Start Time: %s\n", schedule.StartTime)
	fmt.Fprintf(os.Stderr, "  Duration: %s\n", schedule.Duration)
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

func addComplexRecurringSchedule(ctx context.Context, client *gopurple.Client, presentationID int, groupName string, startDate, endDate time.Time, days, startTime, duration string, verbose, jsonMode bool) (*gopurple.ScheduledPresentation, error) {
	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "\nCreating complex recurring schedule...\n")
		fmt.Fprintf(os.Stderr, "  Presentation ID: %d\n", presentationID)
		fmt.Fprintf(os.Stderr, "  Group: %s\n", groupName)
		fmt.Fprintf(os.Stderr, "  Days: %s\n", days)
		fmt.Fprintf(os.Stderr, "  Start Date: %s\n", startDate.Format("2006-01-02"))
		fmt.Fprintf(os.Stderr, "  End Date: %s\n", endDate.Format("2006-01-02"))
		fmt.Fprintf(os.Stderr, "  Start Time: %s\n", startTime)
		fmt.Fprintf(os.Stderr, "  Duration: %s\n", duration)
	}

	// Create schedule object with custom day pattern
	schedule := &gopurple.ScheduledPresentation{
		PresentationID:      presentationID,
		IsRecurrent:         true,
		RecurrenceStartDate: &startDate,
		RecurrenceEndDate:   &endDate,
		DaysOfWeek:          days,
		StartTime:           startTime,
		Duration:            duration,
	}

	result, err := client.Schedules.AddScheduledPresentation(ctx, groupName, schedule)
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return result, nil
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
