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
		helpFlag                = flag.Bool("help", false, "Display usage information")
		jsonFlag                = flag.Bool("json", false, "Output as JSON")
		verboseFlag             = flag.Bool("verbose", false, "Show detailed information")
		timeoutFlag             = flag.Int("timeout", 30, "Request timeout in seconds")
		networkFlag             *string
		groupFlag               = flag.String("group", "", "Group name to add schedules to (required)")
		startDateFlag           = flag.String("start-date", "", "Start date for recurrence (YYYY-MM-DD format) (required)")
		endDateFlag             = flag.String("end-date", "", "End date for recurrence (YYYY-MM-DD format) (required)")
		morningPresentationFlag = flag.Int("morning-presentation", 0, "Presentation ID for morning slot (required)")
		afternoonPresentationFlag = flag.Int("afternoon-presentation", 0, "Presentation ID for afternoon slot (required)")
		eveningPresentationFlag = flag.Int("evening-presentation", 0, "Presentation ID for evening slot (required)")
		morningStartFlag        = flag.String("morning-start", "08:00:00", "Morning start time (HH:MM:SS)")
		afternoonStartFlag      = flag.String("afternoon-start", "12:00:00", "Afternoon start time (HH:MM:SS)")
		eveningStartFlag        = flag.String("evening-start", "17:00:00", "Evening start time (HH:MM:SS)")
		durationFlag            = flag.String("duration", "04:00:00", "Duration for each slot (HH:MM:SS)")
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A command-line tool for day-parting: scheduling different presentations at different times of day.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Schedule three different presentations for morning, afternoon, and evening:\n")
		fmt.Fprintf(os.Stderr, "    %s --group \"Lobby Displays\" --morning-presentation 100 --afternoon-presentation 200 --evening-presentation 300 --start-date 2026-04-01 --end-date 2026-04-30\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  With custom times:\n")
		fmt.Fprintf(os.Stderr, "    %s --group \"Lobby Displays\" --morning-presentation 100 --afternoon-presentation 200 --evening-presentation 300 --start-date 2026-04-01 --end-date 2026-04-30 --morning-start 07:00:00 --afternoon-start 13:00:00 --evening-start 18:00:00\n", os.Args[0])
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
	if *morningPresentationFlag == 0 {
		fmt.Fprintf(os.Stderr, "Error: --morning-presentation must be specified\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *afternoonPresentationFlag == 0 {
		fmt.Fprintf(os.Stderr, "Error: --afternoon-presentation must be specified\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *eveningPresentationFlag == 0 {
		fmt.Fprintf(os.Stderr, "Error: --evening-presentation must be specified\n\n")
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
	if !isValidTimeFormat(*morningStartFlag) {
		log.Fatalf("Invalid morning-start format (use HH:MM:SS): %s", *morningStartFlag)
	}
	if !isValidTimeFormat(*afternoonStartFlag) {
		log.Fatalf("Invalid afternoon-start format (use HH:MM:SS): %s", *afternoonStartFlag)
	}
	if !isValidTimeFormat(*eveningStartFlag) {
		log.Fatalf("Invalid evening-start format (use HH:MM:SS): %s", *eveningStartFlag)
	}
	if !isValidTimeFormat(*durationFlag) {
		log.Fatalf("Invalid duration format (use HH:MM:SS): %s", *durationFlag)
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

	// Create day-parting schedules
	schedules, err := createDayPartingSchedules(ctx, client, *groupFlag, startDate, endDate,
		*morningPresentationFlag, *afternoonPresentationFlag, *eveningPresentationFlag,
		*morningStartFlag, *afternoonStartFlag, *eveningStartFlag, *durationFlag,
		*verboseFlag, *jsonFlag)
	if err != nil {
		log.Fatalf("Failed to create day-parting schedules: %v", err)
	}

	// Output as JSON if requested
	if *jsonFlag {
		result := map[string]interface{}{
			"success":   true,
			"schedules": schedules,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			log.Fatalf("Failed to encode JSON: %v", err)
		}
		return
	}

	fmt.Fprintf(os.Stderr, "\nDay-parting schedules created successfully:\n")
	for i, schedule := range schedules {
		slot := []string{"Morning", "Afternoon", "Evening"}[i]
		fmt.Fprintf(os.Stderr, "\n%s Slot:\n", slot)
		fmt.Fprintf(os.Stderr, "  Schedule ID: %d\n", schedule.ID)
		fmt.Fprintf(os.Stderr, "  Presentation: %s (ID: %d)\n", schedule.PresentationName, schedule.PresentationID)
		fmt.Fprintf(os.Stderr, "  Start Time: %s\n", schedule.StartTime)
		fmt.Fprintf(os.Stderr, "  Duration: %s\n", schedule.Duration)
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

func createDayPartingSchedules(ctx context.Context, client *gopurple.Client, groupName string,
	startDate, endDate time.Time, morningID, afternoonID, eveningID int,
	morningStart, afternoonStart, eveningStart, duration string,
	verbose, jsonMode bool) ([]*gopurple.ScheduledPresentation, error) {

	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "\nCreating day-parting schedules...\n")
		fmt.Fprintf(os.Stderr, "  Group: %s\n", groupName)
		fmt.Fprintf(os.Stderr, "  Date Range: %s to %s\n", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
		fmt.Fprintf(os.Stderr, "  Morning: Presentation %d at %s\n", morningID, morningStart)
		fmt.Fprintf(os.Stderr, "  Afternoon: Presentation %d at %s\n", afternoonID, afternoonStart)
		fmt.Fprintf(os.Stderr, "  Evening: Presentation %d at %s\n", eveningID, eveningStart)
	}

	var schedules []*gopurple.ScheduledPresentation

	// Create morning schedule
	morningSchedule := &gopurple.ScheduledPresentation{
		PresentationID:      morningID,
		IsRecurrent:         true,
		RecurrenceStartDate: &startDate,
		RecurrenceEndDate:   &endDate,
		DaysOfWeek:          "EveryDay",
		StartTime:           morningStart,
		Duration:            duration,
	}

	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "\nCreating morning schedule...\n")
	}

	morning, err := client.Schedules.AddScheduledPresentation(ctx, groupName, morningSchedule)
	if err != nil {
		return nil, fmt.Errorf("failed to create morning schedule: %w", err)
	}
	schedules = append(schedules, morning)

	// Create afternoon schedule
	afternoonSchedule := &gopurple.ScheduledPresentation{
		PresentationID:      afternoonID,
		IsRecurrent:         true,
		RecurrenceStartDate: &startDate,
		RecurrenceEndDate:   &endDate,
		DaysOfWeek:          "EveryDay",
		StartTime:           afternoonStart,
		Duration:            duration,
	}

	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "Creating afternoon schedule...\n")
	}

	afternoon, err := client.Schedules.AddScheduledPresentation(ctx, groupName, afternoonSchedule)
	if err != nil {
		return nil, fmt.Errorf("failed to create afternoon schedule: %w", err)
	}
	schedules = append(schedules, afternoon)

	// Create evening schedule
	eveningSchedule := &gopurple.ScheduledPresentation{
		PresentationID:      eveningID,
		IsRecurrent:         true,
		RecurrenceStartDate: &startDate,
		RecurrenceEndDate:   &endDate,
		DaysOfWeek:          "EveryDay",
		StartTime:           eveningStart,
		Duration:            duration,
	}

	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "Creating evening schedule...\n")
	}

	evening, err := client.Schedules.AddScheduledPresentation(ctx, groupName, eveningSchedule)
	if err != nil {
		return nil, fmt.Errorf("failed to create evening schedule: %w", err)
	}
	schedules = append(schedules, evening)

	return schedules, nil
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
