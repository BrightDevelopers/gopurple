package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/brightsign/gopurple"
)

func main() {
	var (
		helpFlag      = flag.Bool("help", false, "Display usage information")
		jsonFlag      = flag.Bool("json", false, "Output as JSON")
		verboseFlag   = flag.Bool("verbose", false, "Show detailed information")
		timeoutFlag   = flag.Int("timeout", 30, "Request timeout in seconds")
		setupIDFlag   = flag.String("setup-id", "", "ID of the setup record to delete")
		setupNameFlag = flag.String("setup-name", "", "Package name of the setup record to delete")
		forceFlag     = flag.Bool("force", false, "Skip confirmation prompt")
		networkFlag   *string
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A tool to delete B-Deploy setup records from BSN.cloud.\n\n")
		fmt.Fprintf(os.Stderr, "⚠️  WARNING: This operation permanently deletes setup records!\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Delete by setup ID:\n")
		fmt.Fprintf(os.Stderr, "    %s --setup-id \"658f1dbef1d46c829f60a14f\" --network \"My Network\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Delete by setup name (package name):\n")
		fmt.Fprintf(os.Stderr, "    %s --setup-name \"retail-display-v1\" --network \"My Network\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Force delete without confirmation:\n")
		fmt.Fprintf(os.Stderr, "    %s --setup-id \"658f1dbef1d46c829f60a14f\" --force\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Output as JSON:\n")
		fmt.Fprintf(os.Stderr, "    %s --setup-id \"658f1dbef1d46c829f60a14f\" --json --force\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	// Validate required parameter - need either setup-id or setup-name
	if *setupIDFlag == "" && *setupNameFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: Either --setup-id or --setup-name is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Cannot specify both
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

	// Determine network to use
	networkName := getNetworkName(*networkFlag, client, ctx, *verboseFlag)

	// Set network context for B-Deploy operations
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "📡 Setting network context to: %s\n", networkName)
	}
	if err := client.BDeploy.SetNetworkContext(ctx, networkName); err != nil {
		log.Fatalf("❌ Failed to set network context: %v", err)
	}

	if *verboseFlag && !*jsonFlag {
		fmt.Fprintf(os.Stderr, "✅ Network context set successfully!\n")
	}

	// Resolve setup ID (either directly provided or search by name)
	var setupID string
	var setupPackageName string

	if *setupIDFlag != "" {
		// Using setup ID directly
		setupID = *setupIDFlag
		if *verboseFlag && !*jsonFlag {
			fmt.Fprintf(os.Stderr, "🔍 Using setup ID: %s\n", setupID)
		}
	} else {
		// Search by package name
		setupPackageName = *setupNameFlag
		if !*jsonFlag {
			fmt.Fprintf(os.Stderr, "🔍 Searching for setup with package name: %s\n", setupPackageName)
		}

		// Query setup records by package name (API does partial match, so we'll filter for exact match)
		records, err := client.BDeploy.GetSetupRecords(ctx,
			gopurple.WithNetworkName(networkName),
			gopurple.WithPackageName(setupPackageName),
		)
		if err != nil {
			log.Fatalf("❌ Failed to search for setup records: %v", err)
		}

		// Filter for exact matches only (API does substring match)
		var exactMatches []gopurple.BDeployRecord
		for _, record := range records.Items {
			if record.PackageName == setupPackageName {
				exactMatches = append(exactMatches, record)
			}
		}

		if len(exactMatches) == 0 {
			log.Fatalf("❌ No setup record found with exact package name '%s' on network '%s'", setupPackageName, networkName)
		}

		if len(exactMatches) > 1 {
			fmt.Printf("⚠️  Found %d setup records with package name '%s':\n", len(exactMatches), setupPackageName)
			for i, record := range exactMatches {
				fmt.Printf("  %d. ID: %s, Package: %s, Type: %s\n", i+1, record.ID, record.PackageName, record.SetupType)
			}
			log.Fatalf("❌ Multiple setups found. Please use --setup-id to specify which one to delete")
		}

		// Found exactly one exact match
		setupID = exactMatches[0].ID
		setupPackageName = exactMatches[0].PackageName
		if !*jsonFlag {
			fmt.Fprintf(os.Stderr, "✅ Found setup record: ID=%s, Package=%s\n", setupID, setupPackageName)
		}
	}

	// Confirm deletion unless --force is used
	if !*forceFlag && !*jsonFlag {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "⚠️  WARNING: This will permanently delete the B-Deploy setup record!\n")
		fmt.Fprintf(os.Stderr, "Network: %s\n", networkName)
		fmt.Fprintf(os.Stderr, "Setup ID: %s\n", setupID)
		if setupPackageName != "" {
			fmt.Fprintf(os.Stderr, "Package Name: %s\n", setupPackageName)
		}
		fmt.Fprintf(os.Stderr, "\n")

		confirmed := confirmDeletion()
		if !confirmed {
			fmt.Fprintf(os.Stderr, "Operation cancelled.\n")
			return
		}
	}

	// Delete the B-Deploy setup record
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "🗑️  Deleting B-Deploy setup record: %s\n", setupID)
	}

	response, err := client.BDeploy.DeleteSetupRecord(ctx, setupID)
	if err != nil {
		// Check if the error is because the setup is in use
		errMsg := err.Error()
		if strings.Contains(errMsg, "in use") || strings.Contains(errMsg, "is in use") {
			fmt.Fprintf(os.Stderr, "\n")
			fmt.Fprintf(os.Stderr, "To delete this setup, you must first dissociate all devices:\n")
			fmt.Fprintf(os.Stderr, "\n")
			fmt.Fprintf(os.Stderr, "Option 1: List devices using this setup:\n")
			if setupPackageName != "" {
				fmt.Fprintf(os.Stderr, "  ./bin/bdeploy-list-devices --setup-name \"%s\" --network \"%s\"\n", setupPackageName, networkName)
			} else {
				fmt.Fprintf(os.Stderr, "  ./bin/bdeploy-list-devices --setup-id \"%s\" --network \"%s\"\n", setupID, networkName)
			}
			fmt.Fprintf(os.Stderr, "\n")
			fmt.Fprintf(os.Stderr, "Option 2: Dissociate a device from this setup:\n")
			fmt.Fprintf(os.Stderr, "  ./bin/bdeploy-associate --serial <SERIAL> --dissociate --network \"%s\"\n", networkName)
			fmt.Fprintf(os.Stderr, "\n")
			fmt.Fprintf(os.Stderr, "Option 3: Delete devices associated with this setup:\n")
			fmt.Fprintf(os.Stderr, "  ./bin/bdeploy-delete-device --serial <SERIAL> --network \"%s\"\n", networkName)
			fmt.Fprintf(os.Stderr, "\n")
			os.Exit(1)
		}
		log.Fatalf("❌ Failed to delete B-Deploy setup record: %v", err)
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

	// Display results
	fmt.Fprintf(os.Stderr, "✅ B-Deploy setup record deletion completed!\n")

	if *verboseFlag {
		if response.Success {
			fmt.Fprintf(os.Stderr, "🎉 Deletion confirmed successful\n")
		}

		if response.Message != "" {
			fmt.Fprintf(os.Stderr, "📋 Response: %s\n", response.Message)
		}

		if response.Error != "" {
			fmt.Fprintf(os.Stderr, "⚠️  API returned error: %s\n", response.Error)
		}
	}
}

func getNetworkName(requestedNetwork string, client *gopurple.Client, ctx context.Context, verbose bool) string {
	// If network was specified via flag, use it
	if requestedNetwork != "" {
		return requestedNetwork
	}

	// Check if network is already set in client
	if client.IsNetworkSet() {
		if current, err := client.GetCurrentNetwork(ctx); err == nil {
			if verbose {
				fmt.Printf("📡 Using current network: %s (ID: %d)\n", current.Name, current.ID)
			}
			return current.Name
		}
	}

	// Check environment variable
	if envNetwork := os.Getenv("BS_NETWORK"); envNetwork != "" {
		if verbose {
			fmt.Printf("📡 Using network from BS_NETWORK: %s\n", envNetwork)
		}
		return envNetwork
	}

	// Need to select a network
	fmt.Println("📡 Getting available networks...")

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
		if verbose {
			fmt.Printf("📡 Using network: %s (ID: %d)\n", networkName, networks[0].ID)
		}
		return networkName
	}

	// Multiple networks - need user to specify
	fmt.Println("❌ Multiple networks available. Please specify --network or set BS_NETWORK:")
	for i, network := range networks {
		fmt.Printf("  %d. %s (ID: %d)\n", i+1, network.Name, network.ID)
	}
	os.Exit(1)
	return ""
}

func confirmDeletion() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Are you sure you want to delete this record? (y/N): ")
	
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("❌ Failed to read input: %v", err)
	}
	
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}