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
		expirationFlag    = flag.String("expiration", "", "Expiration date for schedule (YYYY-MM-DD format) (required)")
		startTimeFlag     = flag.String("start-time", "09:00:00", "Start time in HH:MM:SS format")
		durationFlag      = flag.String("duration", "01:00:00", "Duration in HH:MM:SS format")
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A command-line tool for scheduling a presentation with an expiration date.\n")
		fmt.Fprintf(os.Stderr, "The schedule will automatically be removed after the expiration date.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Schedule a presentation with expiration:\n")
		fmt.Fprintf(os.Stderr, "    %s --presentation-id 123 --group \"Lobby Displays\" --start-date 2026-04-01 --end-date 2026-04-30 --expiration 2026-05-01\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Schedule with custom times:\n")
		fmt.Fprintf(os.Stderr, "    %s --presentation-id 123 --group \"Lobby Displays\" --start-date 2026-04-01 --end-date 2026-04-30 --expiration 2026-05-01 --start-time 14:00:00 --duration 02:00:00\n", os.Args[0])
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
	if *expirationFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: --expiration must be specified\n\n")
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
	expirationDate, err := time.Parse("2006-01-02", *expirationFlag)
	if err != nil {
		log.Fatalf("Invalid expiration format (use YYYY-MM-DD): %v", err)
	}

	if endDate.Before(startDate) {
		log.Fatalf("End date must be after start date")
	}
	if expirationDate.Before(endDate) {
		log.Fatalf("Expiration date should be after end date")
	}

	// Validate time formats
	if !isValidTimeFormat(*startTimeFlag) {
		log.Fatalf("Invalid start-time format (use HH:MM:SS): %s", *startTimeFlag)
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

	// Add scheduled presentation with expiration
	schedule, err := addScheduleWithExpiration(ctx, client, *presentationIDFlag, *groupFlag, startDate, endDate, expirationDate, *startTimeFlag, *durationFlag, *verboseFlag, *jsonFlag)
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

	fmt.Fprintf(os.Stderr, "\nSchedule with expiration created successfully:\n")
	fmt.Fprintf(os.Stderr, "  Schedule ID: %d\n", schedule.ID)
	fmt.Fprintf(os.Stderr, "  Presentation: %s (ID: %d)\n", schedule.PresentationName, schedule.PresentationID)
	fmt.Fprintf(os.Stderr, "  Group: %s\n", *groupFlag)
	fmt.Fprintf(os.Stderr, "  Start Date: %s\n", startDate.Format("2006-01-02"))
	fmt.Fprintf(os.Stderr, "  End Date: %s\n", endDate.Format("2006-01-02"))
	fmt.Fprintf(os.Stderr, "  Expiration Date: %s\n", expirationDate.Format("2006-01-02"))
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

func addScheduleWithExpiration(ctx context.Context, client *gopurple.Client, presentationID int, groupName string, startDate, endDate, expirationDate time.Time, startTime, duration string, verbose, jsonMode bool) (*gopurple.ScheduledPresentation, error) {
	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "\nCreating schedule with expiration...\n")
		fmt.Fprintf(os.Stderr, "  Presentation ID: %d\n", presentationID)
		fmt.Fprintf(os.Stderr, "  Group: %s\n", groupName)
		fmt.Fprintf(os.Stderr, "  Start Date: %s\n", startDate.Format("2006-01-02"))
		fmt.Fprintf(os.Stderr, "  End Date: %s\n", endDate.Format("2006-01-02"))
		fmt.Fprintf(os.Stderr, "  Expiration Date: %s\n", expirationDate.Format("2006-01-02"))
		fmt.Fprintf(os.Stderr, "  Start Time: %s\n", startTime)
		fmt.Fprintf(os.Stderr, "  Duration: %s\n", duration)
	}

	// Create schedule object with expiration date
	schedule := &gopurple.ScheduledPresentation{
		PresentationID:      presentationID,
		IsRecurrent:         true,
		RecurrenceStartDate: &startDate,
		RecurrenceEndDate:   &endDate,
		DaysOfWeek:          "EveryDay",
		StartTime:           startTime,
		Duration:            duration,
		ExpirationDate:      &expirationDate,
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
