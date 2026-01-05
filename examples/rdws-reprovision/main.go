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
		verboseFlag = flag.Bool("verbose", false, "Show detailed information")
		timeoutFlag = flag.Int("timeout", 30, "Request timeout in seconds")
		networkFlag *string
		serialFlag  = flag.String("serial", "", "Device serial number")
		idFlag      = flag.Int("id", 0, "Device ID")
		confirmFlag = flag.Bool("y", false, "Skip confirmation prompt")
		jsonFlag    = flag.Bool("json", false, "Output as JSON")
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A tool to re-provision BrightSign devices via rDWS API.\n\n")
		fmt.Fprintf(os.Stderr, "⚠️  WARNING: Re-provisioning will:\n")
		fmt.Fprintf(os.Stderr, "  • Reset the device to factory defaults\n")
		fmt.Fprintf(os.Stderr, "  • Clear all content and settings\n")
		fmt.Fprintf(os.Stderr, "  • Reboot the device and run B-Deploy setup again\n")
		fmt.Fprintf(os.Stderr, "  • This action is IRREVERSIBLE\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Re-provision with confirmation:\n")
		fmt.Fprintf(os.Stderr, "    %s --serial UTD41X000009\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Re-provision without confirmation (automation):\n")
		fmt.Fprintf(os.Stderr, "    %s --serial UTD41X000009 -y\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Using device ID instead of serial:\n")
		fmt.Fprintf(os.Stderr, "    %s --id 12345\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Get reprovision result as JSON:\n")
		fmt.Fprintf(os.Stderr, "    %s --serial UTD41X000009 -y --json\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	// Require -y flag when using --json (cannot prompt for confirmation)
	if *jsonFlag && !*confirmFlag {
		fmt.Fprintf(os.Stderr, "Error: -y flag is required when using --json (cannot prompt for confirmation)\n\n")
		flag.Usage()
		os.Exit(1)
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

	// Confirmation prompt unless -y flag is used
	if !*confirmFlag {
		fmt.Printf("🚨 CRITICAL WARNING: Re-provisioning device with %s\n\n", deviceInfo)
		fmt.Println("This operation will:")
		fmt.Println("  ❌ PERMANENTLY DELETE all content and settings")
		fmt.Println("  ❌ Reset device to factory defaults")
		fmt.Println("  ❌ Clear all networking configuration (except essential setup keys)")
		fmt.Println("  🔄 Reboot and restart B-Deploy provisioning process")
		fmt.Println("  ⏱️  Device will be offline during re-provisioning")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Println("📋 Re-provisioning process:")
		fmt.Println("  1. Clear device registry (keeping essential network keys)")
		fmt.Println("  2. Format storage device (SD card)")
		fmt.Println("  3. Reboot device")
		fmt.Println("  4. Device fetches setup package from B-Deploy")
		fmt.Println("  5. Device provisions itself with new configuration")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Printf("⚠️  THIS ACTION CANNOT BE UNDONE!\n\n")

		fmt.Print("Are you absolutely sure you want to re-provision this device? Type 'yes' to continue: ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			fmt.Println("❌ Aborted")
			os.Exit(1)
		}

		response := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if response != "yes" {
			fmt.Println("❌ Re-provisioning cancelled. Only 'yes' confirms this action.")
			os.Exit(0)
		}
	}

	// Perform re-provisioning
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "🔄 Initiating re-provisioning for device with %s\n", deviceInfo)
	}
	success, err := handleReprovision(ctx, client, *serialFlag, *idFlag, *verboseFlag, *jsonFlag)

	if err != nil {
		log.Fatalf("❌ Failed to re-provision device: %v", err)
	}

	// Output as JSON if requested
	if *jsonFlag {
		result := map[string]interface{}{
			"success": success,
			"message": "Re-provisioning initiated successfully",
		}
		if *serialFlag != "" {
			result["serial"] = *serialFlag
		}
		if *idFlag != 0 {
			result["device_id"] = *idFlag
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			log.Fatalf("Failed to encode JSON: %v", err)
		}
		return
	}

	if success {
		fmt.Fprintf(os.Stderr, "✅ Re-provisioning initiated successfully!\n")
		fmt.Fprintf(os.Stderr, "🔄 Device with %s is now rebooting and will re-provision itself\n", deviceInfo)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "📋 Next steps:\n")
		fmt.Fprintf(os.Stderr, "  1. Device is rebooting (this may take 1-2 minutes)\n")
		fmt.Fprintf(os.Stderr, "  2. Device will fetch B-Deploy setup package\n")
		fmt.Fprintf(os.Stderr, "  3. Device will run through provisioning setup\n")
		fmt.Fprintf(os.Stderr, "  4. Monitor device status in BSN.cloud dashboard\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "💡 The device will appear offline until re-provisioning completes.\n")
	} else {
		fmt.Fprintf(os.Stderr, "❌ Re-provisioning request failed\n")
	}
}

func handleNetworkSelection(ctx context.Context, client *gopurple.Client, requestedNetwork string, verbose bool, jsonMode bool) error {
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
				fmt.Fprintf(os.Stderr, "📡 Using network from BS_NETWORK environment variable\n")
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
	fmt.Fprintf(os.Stderr, "📡 Selected network: %s (ID: %d)\n", selectedNetwork.Name, selectedNetwork.ID)

	return client.SetNetworkByID(ctx, selectedNetwork.ID)
}

func handleReprovision(ctx context.Context, client *gopurple.Client, serial string, deviceID int, verbose bool, jsonMode bool) (bool, error) {
	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "📋 Making re-provision API call...\n")
		fmt.Fprintf(os.Stderr, "📋 What happens during re-provisioning:\n")
		fmt.Fprintf(os.Stderr, "  1. Device preserves essential networking keys:\n")
		fmt.Fprintf(os.Stderr, "     • Network interface settings (DHCP, static IP, WiFi)\n")
		fmt.Fprintf(os.Stderr, "     • Proxy settings and time servers\n")
		fmt.Fprintf(os.Stderr, "     • Rate limiting configurations\n")
		fmt.Fprintf(os.Stderr, "     • VLAN settings if applicable\n")
		fmt.Fprintf(os.Stderr, "  2. Device clears all other registry entries\n")
		fmt.Fprintf(os.Stderr, "  3. Device formats the SD card storage\n")
		fmt.Fprintf(os.Stderr, "  4. Device reboots and fetches B-Deploy setup package\n")
		fmt.Fprintf(os.Stderr, "  5. Device runs through complete provisioning process\n")
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Make the API call using the device service
	var response *gopurple.ReprovisionResponse
	var err error

	if serial != "" {
		response, err = client.Devices.ReprovisionBySerial(ctx, serial)
	} else {
		response, err = client.Devices.Reprovision(ctx, deviceID)
	}

	if err != nil {
		return false, fmt.Errorf("failed to re-provision device: %w", err)
	}

	if verbose && !jsonMode {
		fmt.Fprintf(os.Stderr, "✅ API call successful\n")
		if response.Message != "" {
			fmt.Fprintf(os.Stderr, "  Message: %s\n", response.Message)
		}
	}

	return response.Success, nil
}