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
	"syscall"
	"time"

	"github.com/brightdevelopers/gopurple"
	"golang.org/x/term"
)

func main() {
	var (
		helpFlag        = flag.Bool("help", false, "Display usage information")
		jsonFlag        = flag.Bool("json", false, "Output as JSON")
		verboseFlag     = flag.Bool("verbose", false, "Show detailed information")
		timeoutFlag     = flag.Int("timeout", 30, "Request timeout in seconds")
		networkFlag     *string
		serialFlag      = flag.String("serial", "", "Device serial number")
		idFlag          = flag.Int("id", 0, "Device ID")
		getFlag         = flag.Bool("get", false, "Get current DWS password info only")
		passwordFlag    = flag.String("password", "", "New DWS password (empty to remove password)")
		prevPassFlag    = flag.String("prev-password", "", "Previous DWS password (empty if none set)")
		interactiveFlag = flag.Bool("interactive", false, "Prompt for passwords interactively")
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A tool to get and set DWS passwords on BrightSign devices via rDWS API.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Get current password info:\n")
		fmt.Fprintf(os.Stderr, "    %s --serial UTD41X000009 --get\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Set password interactively:\n")
		fmt.Fprintf(os.Stderr, "    %s --serial UTD41X000009 --interactive\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Set password via flags:\n")
		fmt.Fprintf(os.Stderr, "    %s --serial UTD41X000009 --password newpass --prev-password oldpass\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Remove password:\n")
		fmt.Fprintf(os.Stderr, "    %s --serial UTD41X000009 --password \"\" --prev-password oldpass\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Output as JSON:\n")
		fmt.Fprintf(os.Stderr, "    %s --serial UTD41X000009 --get --json\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	// Validate input
	if *serialFlag == "" && *idFlag == 0 {
		fmt.Fprintf(os.Stderr, "❌ Error: Must specify either --serial or --id\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if *serialFlag != "" && *idFlag != 0 {
		fmt.Fprintf(os.Stderr, "❌ Error: Cannot specify both --serial and --id\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if *jsonFlag && *interactiveFlag {
		fmt.Fprintf(os.Stderr, "❌ Error: Cannot use --json with --interactive mode\n\n")
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
		fmt.Fprintf(os.Stderr, "🔧 Creating BSN.cloud client...\n")
	}

	client, err := gopurple.New(opts...)
	if err != nil {
		if gopurple.IsConfigurationError(err) {
			log.Fatalf("❌ Configuration error: %v", err)
		}
		log.Fatalf("❌ Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Authenticate
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "🔐 Authenticating with BSN.cloud...\n")
	}

	if err := client.Authenticate(ctx); err != nil {
		if gopurple.IsAuthenticationError(err) {
			log.Fatalf("❌ Authentication failed: %v", err)
		}
		log.Fatalf("❌ Authentication error: %v", err)
	}

	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "✅ Authentication successful!\n")
	}

	// Handle network selection
	if err := handleNetworkSelection(ctx, client, *networkFlag, *verboseFlag, *jsonFlag); err != nil {
		log.Fatalf("❌ Network selection failed: %v", err)
	}

	// Get device info for display
	var deviceInfo string
	if *serialFlag != "" {
		deviceInfo = fmt.Sprintf("serial %s", *serialFlag)
	} else {
		deviceInfo = fmt.Sprintf("ID %d", *idFlag)
	}

	// Get current password info first
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "📋 Getting current DWS password info for device with %s\n", deviceInfo)
	}
	currentInfo, err := getDWSPasswordInfo(ctx, client, *serialFlag, *idFlag, *verboseFlag, *jsonFlag)
	if err != nil {
		log.Fatalf("❌ Failed to get DWS password info: %v", err)
	}

	// Display current status
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "📊 Current DWS Password Status:\n")
		if currentInfo.Password != nil {
			fmt.Fprintf(os.Stderr, "  Valid Info: %t\n", currentInfo.Password.IsResultValid)
			if currentInfo.Password.IsBlank {
				fmt.Fprintf(os.Stderr, "  Status: ✅ No password currently set (blank)\n")
			} else {
				fmt.Fprintf(os.Stderr, "  Status: 🔒 Password is currently set\n")
			}
		} else {
			fmt.Fprintf(os.Stderr, "  Status: ❓ Unable to determine password status\n")
		}
	}

	// If only getting info, exit here
	if *getFlag {
		if *jsonFlag {
			result := map[string]interface{}{
				"serial":        *serialFlag,
				"deviceID":      *idFlag,
				"passwordInfo":  currentInfo,
				"isPasswordSet": currentInfo.Password != nil && !currentInfo.Password.IsBlank,
				"isResultValid": currentInfo.Password != nil && currentInfo.Password.IsResultValid,
			}
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(result); err != nil {
				log.Fatalf("Failed to encode JSON: %v", err)
			}
			return
		}
		fmt.Fprintf(os.Stderr, "✅ DWS password info retrieved successfully!\n")
		return
	}

	// Handle password setting
	var newPassword, prevPassword string

	if *interactiveFlag {
		// Interactive mode - prompt for passwords
		newPassword, prevPassword = promptForPasswords(currentInfo.Password)
	} else {
		// Use flags
		newPassword = *passwordFlag
		prevPassword = *prevPassFlag
	}

	// Set the password
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "🔒 Setting DWS password for device with %s\n", deviceInfo)
	}
	setResponse, err := setDWSPassword(ctx, client, *serialFlag, *idFlag, newPassword, prevPassword, *verboseFlag, *jsonFlag)
	if err != nil {
		log.Fatalf("❌ Failed to set DWS password: %v", err)
	}

	if !*jsonFlag {
		if setResponse.Success {
			if newPassword == "" {
				fmt.Fprintf(os.Stderr, "✅ DWS password removed successfully!\n")
			} else {
				fmt.Fprintf(os.Stderr, "✅ DWS password set successfully!\n")
			}

			if setResponse.Reboot {
				fmt.Fprintf(os.Stderr, "🔄 Device will reboot to apply the password change\n")
			}
		} else {
			fmt.Fprintf(os.Stderr, "❌ Failed to set DWS password\n")
		}
	}

	// Verify the change by getting password info again
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "🔍 Verifying password change...\n")
	}
	verifyInfo, err := getDWSPasswordInfo(ctx, client, *serialFlag, *idFlag, *verboseFlag, *jsonFlag)
	if err != nil {
		log.Printf("⚠️  Warning: Failed to verify password change: %v", err)
		return
	}

	// Check if the change was applied as expected
	expectedBlank := (newPassword == "")
	changeVerified := verifyInfo.Password != nil && verifyInfo.Password.IsBlank == expectedBlank

	// Output as JSON if requested
	if *jsonFlag {
		result := map[string]interface{}{
			"success":         setResponse.Success,
			"serial":          *serialFlag,
			"deviceID":        *idFlag,
			"passwordRemoved": newPassword == "",
			"passwordSet":     newPassword != "",
			"rebootRequired":  setResponse.Reboot,
			"changeVerified":  changeVerified,
			"verifyInfo":      verifyInfo,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			log.Fatalf("Failed to encode JSON: %v", err)
		}
		return
	}

	fmt.Fprintf(os.Stderr, "📊 Updated DWS Password Status:\n")
	if verifyInfo.Password != nil {
		fmt.Fprintf(os.Stderr, "  Valid Info: %t\n", verifyInfo.Password.IsResultValid)
		if verifyInfo.Password.IsBlank {
			fmt.Fprintf(os.Stderr, "  Status: ✅ No password set (blank)\n")
		} else {
			fmt.Fprintf(os.Stderr, "  Status: 🔒 Password is set\n")
		}

		if changeVerified {
			fmt.Fprintf(os.Stderr, "✅ Password change verified successfully!\n")
		} else {
			fmt.Fprintf(os.Stderr, "⚠️  Warning: Password status may not have changed as expected\n")
		}
	} else {
		fmt.Fprintf(os.Stderr, "  Status: ❓ Unable to determine updated password status\n")
	}
}

func handleNetworkSelection(ctx context.Context, client *gopurple.Client, requestedNetwork string, verbose, jsonMode bool) error {
	// Check if network is already set
	if client.IsNetworkSet() {
		if current, err := client.GetCurrentNetwork(ctx); err == nil {
			if !jsonMode {
				fmt.Fprintf(os.Stderr, "📡 Using network: %s (ID: %d)\n", current.Name, current.ID)
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
	if !jsonMode {
		fmt.Fprintf(os.Stderr, "📡 Getting available networks...\n")
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
					fmt.Fprintf(os.Stderr, "📡 Using requested network: %s (ID: %d)\n", network.Name, network.ID)
				}
				return client.SetNetworkByID(ctx, network.ID)
			}
		}

		// Network not found - show error and fall back to interactive selection
		if !jsonMode {
			fmt.Fprintf(os.Stderr, "❌ Network '%s' not found. Available networks:\n", requestedNetwork)
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
		if !jsonMode {
			fmt.Fprintf(os.Stderr, "📡 Using network: %s (ID: %d)\n", networks[0].Name, networks[0].ID)
		}
		return client.SetNetworkByID(ctx, networks[0].ID)
	}

	// In JSON mode, cannot do interactive selection
	if jsonMode {
		return fmt.Errorf("network selection required: use --network flag or BS_NETWORK environment variable")
	}

	// Show available networks and let user choose
	if requestedNetwork == "" {
		fmt.Fprintf(os.Stderr, "📡 Available networks:\n")
		for i, network := range networks {
			fmt.Fprintf(os.Stderr, "  %d. %s (ID: %d)\n", i+1, network.Name, network.ID)
			if verbose {
				fmt.Fprintf(os.Stderr, "     Created: %s, Modified: %s\n",
					network.CreationDate.Format("2006-01-02"),
					network.LastModifiedDate.Format("2006-01-02"))
			}
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
	fmt.Fprintf(os.Stderr, "📡 Selected network: %s (ID: %d)\n", selectedNetwork.Name, selectedNetwork.ID)

	return client.SetNetworkByID(ctx, selectedNetwork.ID)
}

func getDWSPasswordInfo(ctx context.Context, client *gopurple.Client, serial string, deviceID int, verbose, jsonMode bool) (*gopurple.DWSPasswordGetResponse, error) {
	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "📋 Making DWS password info API call...\n")
	}

	// Make the API call using the device service
	var response *gopurple.DWSPasswordGetResponse
	var err error

	if serial != "" {
		response, err = client.Devices.GetDWSPasswordBySerial(ctx, serial)
	} else {
		response, err = client.Devices.GetDWSPassword(ctx, deviceID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get DWS password info: %w", err)
	}

	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "✅ API call successful\n")
	}

	return response, nil
}

func setDWSPassword(ctx context.Context, client *gopurple.Client, serial string, deviceID int, newPassword, prevPassword string, verbose, jsonMode bool) (*gopurple.DWSPasswordSetResponse, error) {
	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "📋 Making DWS password set API call...\n")
		if newPassword == "" {
			fmt.Fprintf(os.Stderr, "  Action: Removing password\n")
		} else {
			fmt.Fprintf(os.Stderr, "  Action: Setting new password\n")
		}
	}

	// Build the request
	request := &gopurple.DWSPasswordRequest{
		Password:         newPassword,
		PreviousPassword: prevPassword,
	}

	// Make the API call using the device service
	var response *gopurple.DWSPasswordSetResponse
	var err error

	if serial != "" {
		response, err = client.Devices.SetDWSPasswordBySerial(ctx, serial, request)
	} else {
		response, err = client.Devices.SetDWSPassword(ctx, deviceID, request)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to set DWS password: %w", err)
	}

	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "✅ API call successful\n")
		if response.Reboot {
			fmt.Fprintf(os.Stderr, "  Device will reboot to apply changes\n")
		}
	}

	return response, nil
}

func promptForPasswords(currentInfo *gopurple.DWSPasswordInfo) (string, string) {
	// Prompt for previous password
	var prevPassword string
	if currentInfo != nil && !currentInfo.IsBlank {
		fmt.Print("Enter current DWS password: ")
		prevPassword = promptForPassword()
	} else {
		fmt.Println("No current password is set")
		prevPassword = ""
	}

	// Prompt for new password
	fmt.Print("Enter new DWS password (or press Enter to remove password): ")
	newPassword := promptForPassword()

	return newPassword, prevPassword
}

func promptForPassword() string {
	fmt.Print("Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintf(os.Stderr, "\n") // Print newline after hidden input
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	return string(bytePassword)
}
