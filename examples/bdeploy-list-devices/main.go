package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/brightdevelopers/gopurple"
)

func main() {
	var (
		helpFlag      = flag.Bool("help", false, "Display usage information")
		jsonFlag      = flag.Bool("json", false, "Output as JSON")
		verboseFlag   = flag.Bool("verbose", false, "Show detailed information")
		timeoutFlag   = flag.Int("timeout", 30, "Request timeout in seconds")
		summaryFlag   = flag.Bool("summary", false, "Show only summary count")
		debugFlag     = flag.Bool("debug", false, "Show raw API request and response details")
		detailedFlag  = flag.Bool("detailed", false, "Fetch full details for each device (slower, but shows setup info)")
		setupIDFlag   = flag.String("setup-id", "", "Filter by setup ID")
		setupNameFlag = flag.String("setup-name", "", "Filter by setup package name")
		networkFlag   *string
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A tool to list all B-Deploy devices on the network from BSN.cloud.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  List all devices:\n")
		fmt.Fprintf(os.Stderr, "    %s --network \"My Network\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  List devices with setup information (slower):\n")
		fmt.Fprintf(os.Stderr, "    %s --network \"My Network\" --detailed\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  List devices using a specific setup ID:\n")
		fmt.Fprintf(os.Stderr, "    %s --setup-id \"658f1dbef1d46c829f60a14f\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  List devices using a specific setup name:\n")
		fmt.Fprintf(os.Stderr, "    %s --setup-name \"my-setup\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Show summary only:\n")
		fmt.Fprintf(os.Stderr, "    %s --summary\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Show verbose output:\n")
		fmt.Fprintf(os.Stderr, "    %s --verbose\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Output as JSON:\n")
		fmt.Fprintf(os.Stderr, "    %s --network \"My Network\" --json\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	// Validate that both setup filters aren't specified
	if *setupIDFlag != "" && *setupNameFlag != "" {
		fmt.Fprintf(os.Stderr, "Error: Cannot specify both --setup-id and --setup-name\n\n")
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

	// Enable debug logging if requested
	if *debugFlag {
		opts = append(opts, gopurple.WithDebug(true))
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
			log.Fatalf("❌ Authentication error: %v", err)
		}
		log.Fatalf("❌ Failed to authenticate: %v", err)
	}

	// Get network name from command line or environment
	networkName := *networkFlag
	if networkName == "" {
		networkName = os.Getenv("BS_NETWORK")
	}

	// Set network context if provided
	if networkName != "" {
		if !*jsonFlag {
			fmt.Fprintf(os.Stderr, "🌐 Setting network context to: %s\n", networkName)
		}
		err = client.BDeploy.SetNetworkContext(ctx, networkName)
		if err != nil {
			log.Fatalf("❌ Failed to set network context: %v", err)
		}
	}

	// Build device query options
	var deviceOpts []gopurple.BDeployDeviceListOption

	// Use setup-name filter via API query parameter if provided
	if *setupNameFlag != "" {
		if !*jsonFlag {
			fmt.Fprintf(os.Stderr, "🔍 Filtering devices by setup name: %s\n", *setupNameFlag)
		}
		deviceOpts = append(deviceOpts, gopurple.WithSetupName(*setupNameFlag))
	}

	// Debug: Show the request that will be made
	if *debugFlag && !*jsonFlag {
		fmt.Fprintf(os.Stderr, "\n🔍 DEBUG: Request Parameters:\n")
		fmt.Fprintf(os.Stderr, "   Endpoint: https://provision.bsn.cloud/rest-device/v2/device/\n")
		if networkName != "" {
			fmt.Fprintf(os.Stderr, "   NetworkName: %s\n", networkName)
		}
		if *setupNameFlag != "" {
			fmt.Fprintf(os.Stderr, "   query[setupName]: %s\n", *setupNameFlag)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Get all devices (with optional setup name filter)
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "📱 Fetching B-Deploy devices...\n")
	}

	response, err := client.BDeploy.GetAllDevices(ctx, deviceOpts...)
	if err != nil {
		log.Fatalf("❌ Failed to get devices: %v", err)
	}

	// Client-side filtering by setup ID if specified
	var filterSetupID string
	if *setupIDFlag != "" {
		filterSetupID = *setupIDFlag
		if !*jsonFlag {
			fmt.Fprintf(os.Stderr, "🔍 Filtering devices by setup ID: %s\n", filterSetupID)
		}
	}

	// Debug: Show what GetAllDevices initially returns
	if *debugFlag && !*jsonFlag && len(response.Players) > 0 {
		fmt.Fprintf(os.Stderr, "\n🔍 DEBUG: Initial GetAllDevices() response for first device:\n")
		fmt.Fprintf(os.Stderr, "   Serial: %s\n", response.Players[0].Serial)
		fmt.Fprintf(os.Stderr, "   SetupName: '%s'\n", response.Players[0].SetupName)
		fmt.Fprintf(os.Stderr, "   SetupID: '%s'\n", response.Players[0].SetupID)
		fmt.Fprintf(os.Stderr, "   URL: '%s'\n", response.Players[0].URL)
		initialJSON, _ := json.MarshalIndent(response.Players[0], "   ", "  ")
		fmt.Fprintf(os.Stderr, "\n   Full JSON:\n   %s\n\n", string(initialJSON))
	}

	// Fetch detailed information for each device if requested
	if *detailedFlag {
		if !*jsonFlag {
			fmt.Fprintf(os.Stderr, "🔍 Fetching detailed information for %d devices...\n", len(response.Players))
		}

		for i := range response.Players {
			if !*jsonFlag && i > 0 && i%10 == 0 {
				fmt.Fprintf(os.Stderr, "   Progress: %d/%d devices...\n", i, len(response.Players))
			}

			// Fetch individual device details by serial number
			detailResp, err := client.BDeploy.GetDeviceBySerial(ctx, response.Players[i].Serial)
			if err != nil {
				if !*jsonFlag {
					fmt.Fprintf(os.Stderr, "⚠️  Failed to get details for device %s: %v\n", response.Players[i].Serial, err)
				}
				continue
			}

			// Update the device with detailed information if available
			if len(detailResp.Result.Players) > 0 {
				detailedDevice := detailResp.Result.Players[0]

				// Debug: Show what we got from the API
				if *debugFlag && !*jsonFlag && i == 0 {
					fmt.Fprintf(os.Stderr, "\n🔍 DEBUG: First detailed device response:\n")
					fmt.Fprintf(os.Stderr, "   Serial: %s\n", detailedDevice.Serial)
					fmt.Fprintf(os.Stderr, "   SetupName: '%s'\n", detailedDevice.SetupName)
					fmt.Fprintf(os.Stderr, "   SetupID: '%s'\n", detailedDevice.SetupID)
					fmt.Fprintf(os.Stderr, "   URL: '%s'\n", detailedDevice.URL)
					fmt.Fprintf(os.Stderr, "   Name: '%s'\n", detailedDevice.Name)
					detailJSON, _ := json.MarshalIndent(detailedDevice, "   ", "  ")
					fmt.Fprintf(os.Stderr, "\n   Full JSON:\n   %s\n\n", string(detailJSON))
				}

				// Copy over the setup information, but only if not empty
				// Note: GetDeviceBySerial may return setupName but not setupId/url
				// GetAllDevices may return setupId/url but not setupName
				// So we only update fields that have values to avoid overwriting
				if detailedDevice.SetupName != "" {
					response.Players[i].SetupName = detailedDevice.SetupName
				}
				if detailedDevice.SetupID != "" {
					response.Players[i].SetupID = detailedDevice.SetupID
				}
				if detailedDevice.URL != "" {
					response.Players[i].URL = detailedDevice.URL
				}
			}
		}

		if !*jsonFlag {
			fmt.Fprintf(os.Stderr, "✅ Detailed information fetched for all devices\n")
		}
	}

	// Filter devices by setup ID if specified
	var filteredDevices []gopurple.BDeployDevice
	if filterSetupID != "" {
		if *debugFlag && !*jsonFlag {
			fmt.Fprintf(os.Stderr, "🔍 Filtering for setup ID: %s\n", filterSetupID)
			fmt.Fprintf(os.Stderr, "🔍 Total devices before filtering: %d\n", len(response.Players))
		}

		for _, device := range response.Players {
			if *debugFlag && !*jsonFlag {
				fmt.Fprintf(os.Stderr, "   Device %s (serial: %s) has SetupID: '%s'\n", device.ID, device.Serial, device.SetupID)
			}
			if device.SetupID == filterSetupID {
				filteredDevices = append(filteredDevices, device)
			}
		}
		response.Players = filteredDevices
		response.Matched = len(filteredDevices)

		if len(filteredDevices) == 0 {
			if !*jsonFlag {
				fmt.Fprintf(os.Stderr, "\n❌ No devices found using setup ID: %s\n", filterSetupID)
				fmt.Fprintf(os.Stderr, "\nThis setup may not be associated with any devices.\n")
				fmt.Fprintf(os.Stderr, "\n💡 Tip: Use --debug flag to see all devices and their setup IDs\n")
			}
			os.Exit(0)
		}

		if !*jsonFlag {
			fmt.Fprintf(os.Stderr, "✅ Found %d device(s) using setup ID: %s\n", len(filteredDevices), filterSetupID)
		}
	} else if *setupNameFlag != "" {
		// Setup name was filtered via API - check if any devices were returned
		if len(response.Players) == 0 {
			if !*jsonFlag {
				fmt.Fprintf(os.Stderr, "\n❌ No devices found using setup name: %s\n", *setupNameFlag)
				fmt.Fprintf(os.Stderr, "\nThis setup may not be associated with any devices, or the setup name may not exist.\n")
			}
			os.Exit(0)
		}

		if !*jsonFlag {
			fmt.Fprintf(os.Stderr, "✅ Found %d device(s) using setup name: %s\n", len(response.Players), *setupNameFlag)
		}
	}

	// Debug output
	if *verboseFlag && !*jsonFlag {
		fmt.Fprintf(os.Stderr, "🔍 Debug - API Response: Total=%d, Matched=%d, Players count=%d\n",
			response.Total, response.Matched, len(response.Players))
	}

	// Raw response output for debugging
	if *debugFlag && !*jsonFlag {
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Failed to marshal response: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "🔍 Raw API Response:\n%s\n", string(jsonData))
		}

		// Also show first device raw JSON for detailed inspection
		if len(response.Players) > 0 {
			fmt.Fprintf(os.Stderr, "\n🔍 First Device Raw JSON:\n")
			devJSON, _ := json.MarshalIndent(response.Players[0], "", "  ")
			fmt.Fprintf(os.Stderr, "%s\n", string(devJSON))
		}
	}

	// Check if any devices were found
	if len(response.Players) == 0 {
		if *jsonFlag {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(response); err != nil {
				log.Fatalf("Failed to encode JSON: %v", err)
			}
			return
		}
		fmt.Fprintf(os.Stderr, "❌ No devices found (Total: %d, Matched: %d)\n", response.Total, response.Matched)
		if response.Total > 0 || response.Matched > 0 {
			fmt.Fprintf(os.Stderr, "⚠️  API indicates devices exist but Players array is empty - possible API format issue\n")
		}
		os.Exit(1)
	}

	// Output as JSON if requested
	if *jsonFlag {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(response); err != nil {
			log.Fatalf("Failed to encode JSON: %v", err)
		}
		return
	}

	// Display summary if requested
	if *summaryFlag {
		fmt.Fprintf(os.Stderr, "\n📊 Summary:\n")
		fmt.Fprintf(os.Stderr, "   Total devices:   %d\n", response.Total)
		fmt.Fprintf(os.Stderr, "   Matched devices: %d\n", response.Matched)
		fmt.Fprintf(os.Stderr, "   Listed devices:  %d\n", len(response.Players))
		fmt.Fprintf(os.Stderr, "\n✅ Device summary completed successfully!\n")
		return
	}

	// Display the devices in a table format
	fmt.Fprintf(os.Stderr, "\n📱 B-Deploy Devices (Total: %d, Matched: %d):\n\n", response.Total, response.Matched)

	// Table headers
	fmt.Fprintf(os.Stderr, "%-20s %-14s %-20s %-12s %-25s %-20s %-20s\n",
		"ID", "SERIAL", "NAME", "MODEL", "DESCRIPTION", "SETUP NAME", "SETUP ID")
	fmt.Fprintf(os.Stderr, "%s\n", strings.Repeat("-", 133))

	// Table rows
	for _, device := range response.Players {
		// Truncate long fields to fit in table
		id := truncateString(device.ID, 20)
		serial := truncateString(device.Serial, 14)
		name := truncateString(device.Name, 20)
		model := truncateString(device.Model, 12)
		desc := truncateString(device.Desc, 25)
		setupName := truncateString(device.SetupName, 20)
		setupID := truncateString(device.SetupID, 20)

		fmt.Fprintf(os.Stderr, "%-20s %-14s %-20s %-12s %-25s %-20s %-20s\n",
			id, serial, name, model, desc, setupName, setupID)
	}

	if *verboseFlag {
		fmt.Fprintf(os.Stderr, "\n🔍 Detailed Information:\n")
		for i, device := range response.Players {
			fmt.Fprintf(os.Stderr, "\n[%d] Device Details:\n", i+1)
			fmt.Fprintf(os.Stderr, "   ID:          %s\n", device.ID)
			fmt.Fprintf(os.Stderr, "   Serial:      %s\n", device.Serial)
			fmt.Fprintf(os.Stderr, "   Name:        %s\n", device.Name)
			fmt.Fprintf(os.Stderr, "   Model:       %s\n", device.Model)
			fmt.Fprintf(os.Stderr, "   Description: %s\n", device.Desc)
			fmt.Fprintf(os.Stderr, "   Setup ID:    %s\n", device.SetupID)
			fmt.Fprintf(os.Stderr, "   Setup Name:  %s\n", device.SetupName)
			fmt.Fprintf(os.Stderr, "   Username:    %s\n", device.Username)
			fmt.Fprintf(os.Stderr, "   Network:     %s\n", device.NetworkName)
			fmt.Fprintf(os.Stderr, "   Client:      %s\n", device.Client)
			fmt.Fprintf(os.Stderr, "   URL:         %s\n", device.URL)
			fmt.Fprintf(os.Stderr, "   Created:     %s\n", device.CreatedAt.Format(time.RFC3339))
			fmt.Fprintf(os.Stderr, "   Updated:     %s\n", device.UpdatedAt.Format(time.RFC3339))
			fmt.Fprintf(os.Stderr, "   Version:     %d\n", device.Version)
		}
	}

	fmt.Fprintf(os.Stderr, "\n✅ Device listing completed successfully!\n")
}

// truncateString truncates a string to the specified length with ellipsis if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
