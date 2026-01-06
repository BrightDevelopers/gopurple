package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brightdevelopers/gopurple"
)

func main() {
	var (
		helpFlag      = flag.Bool("help", false, "Display usage information")
		timeoutFlag   = flag.Int("timeout", 30, "Request timeout in seconds")
		setupIDFlag   = flag.String("setup-id", "", "ID of the setup record to retrieve")
		setupNameFlag = flag.String("setup-name", "", "Package name of the setup record to retrieve")
		jsonFlag      = flag.Bool("json", false, "Output raw JSON (default shows formatted structure)")
		networkFlag   *string
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A tool to retrieve a specific B-Deploy setup record from BSN.cloud.\n\n")
		fmt.Fprintf(os.Stderr, "By default, displays the setup record structure in formatted JSON.\n")
		fmt.Fprintf(os.Stderr, "Use --json flag to output raw JSON only (no status messages).\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Get a setup record by ID (formatted with status messages):\n")
		fmt.Fprintf(os.Stderr, "    %s --setup-id \"658f1dbef1d46c829f60a14f\" --network \"My Network\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Get a setup record by package name:\n")
		fmt.Fprintf(os.Stderr, "    %s --setup-name \"retail-display-v1\" --network \"My Network\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Get raw JSON output only:\n")
		fmt.Fprintf(os.Stderr, "    %s --setup-id \"658f1dbef1d46c829f60a14f\" --network \"My Network\" --json\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	// Validate required parameters - exactly one of setup-id or setup-name must be provided
	setupID := *setupIDFlag
	setupName := *setupNameFlag

	if setupID == "" && setupName == "" {
		fmt.Fprintf(os.Stderr, "Error: either --setup-id or --setup-name is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if setupID != "" && setupName != "" {
		fmt.Fprintf(os.Stderr, "Error: cannot specify both --setup-id and --setup-name\n\n")
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
		fmt.Println("🔧 Creating BSN.cloud client...")
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
		fmt.Println("🔐 Authenticating with BSN.cloud...")
	}

	if err := client.Authenticate(ctx); err != nil {
		if gopurple.IsAuthenticationError(err) {
			log.Fatalf("❌ Authentication failed: %v", err)
		}
		log.Fatalf("❌ Authentication error: %v", err)
	}

	if !*jsonFlag {
		fmt.Println("✅ Authentication successful!")
	}

	// Determine network to use
	networkName := getNetworkName(*networkFlag, client, ctx, *jsonFlag)

	// Set network context for B-Deploy operations
	if !*jsonFlag {
		fmt.Printf("📡 Setting network context to: %s\n", networkName)
	}
	if err := client.BDeploy.SetNetworkContext(ctx, networkName); err != nil {
		log.Fatalf("❌ Failed to set network context: %v", err)
	}

	var record *gopurple.BDeploySetupRecord

	// Get the B-Deploy setup record
	if setupID != "" {
		// Get by ID
		if !*jsonFlag {
			fmt.Printf("📋 Fetching B-Deploy setup record: %s\n", setupID)
		}

		var err error
		record, err = client.BDeploy.GetSetupRecord(ctx, setupID)
		if err != nil {
			log.Fatalf("❌ Failed to get B-Deploy setup record: %v", err)
		}
	} else {
		// Get by package name (setup name)
		if !*jsonFlag {
			fmt.Printf("🔍 Searching for setup with package name: %s\n", setupName)
		}

		opts := []gopurple.BDeployListOption{
			gopurple.WithNetworkName(networkName),
			gopurple.WithPackageName(setupName),
		}

		records, err := client.BDeploy.GetSetupRecords(ctx, opts...)
		if err != nil {
			log.Fatalf("❌ Failed to search for setup records: %v", err)
		}

		if records.TotalCount == 0 {
			log.Fatalf("❌ No setup found with package name: %s", setupName)
		}

		if records.TotalCount > 1 {
			if !*jsonFlag {
				fmt.Fprintf(os.Stderr, "⚠️  Found %d setups with package name '%s':\n", records.TotalCount, setupName)
				for i, rec := range records.Items {
					fmt.Fprintf(os.Stderr, "  %d. ID: %s, Type: %s, Group: %s\n",
						i+1, rec.ID, rec.SetupType, rec.BSNGroupName)
				}
				fmt.Fprintf(os.Stderr, "\nDisplaying first matching setup:\n")
			}
		}

		// Use the first matching record
		if !*jsonFlag {
			fmt.Printf("✅ Found setup with ID: %s\n", records.Items[0].ID)
		}

		// Convert BDeployRecord to BDeploySetupRecord by fetching full details
		record, err = client.BDeploy.GetSetupRecord(ctx, records.Items[0].ID)
		if err != nil {
			log.Fatalf("❌ Failed to get full setup record details: %v", err)
		}
	}

	// Display results
	if *jsonFlag {
		jsonOutput, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			log.Fatalf("❌ Failed to marshal JSON: %v", err)
		}
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("\n✅ B-Deploy Setup Record:")
		printSetupRecord(record)
	}
}

func getNetworkName(requestedNetwork string, client *gopurple.Client, ctx context.Context, jsonOutput bool) string {
	// If network was specified via flag, use it
	if requestedNetwork != "" {
		return requestedNetwork
	}

	// Check if network is already set in client
	if client.IsNetworkSet() {
		if current, err := client.GetCurrentNetwork(ctx); err == nil {
			return current.Name
		}
	}

	// Check environment variable
	if envNetwork := os.Getenv("BS_NETWORK"); envNetwork != "" {
		return envNetwork
	}

	// Need to select a network
	if !jsonOutput {
		fmt.Println("📡 Getting available networks...")
	}

	networks, err := client.GetNetworks(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to get networks: %v", err)
	}

	if len(networks) == 0 {
		log.Fatalf("❌ No networks available")
	}

	// If only one network, use it automatically
	if len(networks) == 1 {
		networkName := networks[0].Name
		return networkName
	}

	// Multiple networks - need user to specify
	if !jsonOutput {
		fmt.Println("❌ Multiple networks available. Please specify --network or set BS_NETWORK:")
		for i, network := range networks {
			fmt.Printf("  %d. %s (ID: %d)\n", i+1, network.Name, network.ID)
		}
	}
	os.Exit(1)
	return ""
}

func printSetupRecord(record *gopurple.BDeploySetupRecord) {
	// Marshal to JSON for pretty printing with indentation
	jsonOutput, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		log.Printf("⚠️  Failed to format setup record: %v", err)
		// Fallback to basic output
		fmt.Fprintf(os.Stderr, "ID:           %s\n", record.ID)
		fmt.Fprintf(os.Stderr, "Setup Type:   %s\n", record.SetupType)
		fmt.Fprintf(os.Stderr, "Version:      %s\n", record.Version)
		return
	}
	fmt.Println(string(jsonOutput))
}
