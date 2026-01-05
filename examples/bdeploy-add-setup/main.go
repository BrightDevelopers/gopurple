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

	"github.com/brightsign/gopurple"
)

// SetupConfig represents the JSON configuration file with optional timeout field
type SetupConfig struct {
	gopurple.BDeploySetupRecord
	Timeout int `json:"timeout,omitempty"` // HTTP timeout (not part of B-Deploy API)
}

func main() {
	var (
		helpFlag    = flag.Bool("help", false, "Display usage information")
		verboseFlag = flag.Bool("verbose", false, "Show detailed information")
		jsonFlag    = flag.Bool("json", false, "Output as JSON")
		timeoutFlag = flag.Int("timeout", 30, "Request timeout in seconds (overrides config file)")
	)

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <config.json>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A tool to create a B-Deploy setup record using a JSON configuration file.\n\n")
		fmt.Fprintf(os.Stderr, "This program performs the following workflow:\n")
		fmt.Fprintf(os.Stderr, "  1. Authenticate with BSN.cloud\n")
		fmt.Fprintf(os.Stderr, "  2. Select and set the network context\n")
		fmt.Fprintf(os.Stderr, "  3. Generate a device registration token\n")
		fmt.Fprintf(os.Stderr, "  4. Create a B-Deploy setup record\n")
		fmt.Fprintf(os.Stderr, "  5. Output setup-id for use with player association\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Create setup using default config:\n")
		fmt.Fprintf(os.Stderr, "    %s config.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Create setup with verbose output:\n")
		fmt.Fprintf(os.Stderr, "    %s --verbose config.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Use custom timeout:\n")
		fmt.Fprintf(os.Stderr, "    %s --timeout 60 config.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Output as JSON:\n")
		fmt.Fprintf(os.Stderr, "    %s --json config.json\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	// Validate command line arguments
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Error: config file is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	configFile := flag.Arg(0)

	// Load configuration
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "📋 Loading configuration from: %s\n", configFile)
	}
	setupConfig, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Validate configuration
	if err := validateConfig(&setupConfig.BDeploySetupRecord); err != nil {
		log.Fatalf("❌ Invalid configuration: %v", err)
	}

	if *verboseFlag && !*jsonFlag {
		fmt.Fprintf(os.Stderr, "✅ Configuration loaded successfully\n")
		fmt.Fprintf(os.Stderr, "   Network: %s\n", setupConfig.BDeploy.NetworkName)
		fmt.Fprintf(os.Stderr, "   Username: %s\n", setupConfig.BDeploy.Username)
		fmt.Fprintf(os.Stderr, "   Package: %s\n", setupConfig.BDeploy.PackageName)
		fmt.Fprintf(os.Stderr, "   Setup Type: %s\n", setupConfig.SetupType)
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Use timeout from flag if specified, otherwise from config
	timeout := *timeoutFlag
	if timeout == 30 && setupConfig.Timeout > 0 {
		timeout = setupConfig.Timeout
	}

	// Create client
	var opts []gopurple.Option
	if timeout > 0 {
		opts = append(opts, gopurple.WithTimeout(time.Duration(timeout)*time.Second))
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

	// Step 1: Authenticate
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

	// Step 2: Verify network exists
	if *verboseFlag && !*jsonFlag {
		fmt.Fprintf(os.Stderr, "📡 Verifying network access...\n")
		networks, err := client.GetNetworks(ctx)
		if err != nil {
			log.Fatalf("❌ Failed to get networks: %v", err)
		}

		found := false
		for _, net := range networks {
			if net.Name == setupConfig.BDeploy.NetworkName {
				found = true
				fmt.Fprintf(os.Stderr, "✅ Network found: %s (ID: %d)\n", net.Name, net.ID)
				break
			}
		}

		if !found {
			fmt.Fprintf(os.Stderr, "⚠️  Warning: Network '%s' not found in available networks\n", setupConfig.BDeploy.NetworkName)
			fmt.Fprintf(os.Stderr, "   Available networks:\n")
			for _, net := range networks {
				fmt.Fprintf(os.Stderr, "     - %s (ID: %d)\n", net.Name, net.ID)
			}
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Step 3: Set network context
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "📡 Setting network context to: %s\n", setupConfig.BDeploy.NetworkName)
	}
	if err := client.BDeploy.SetNetworkContext(ctx, setupConfig.BDeploy.NetworkName); err != nil {
		log.Fatalf("❌ Failed to set network context: %v", err)
	}
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "✅ Network context set successfully!\n")
	}

	// Step 3.5: Generate device registration token (if not already provided)
	if setupConfig.BSNDeviceRegistrationTokenEntity == nil {
		if !*jsonFlag {
			fmt.Fprintf(os.Stderr, "🔑 Generating device registration token...\n")
		}
		deviceToken, err := client.Provisioning.GenerateDeviceToken(ctx)
		if err != nil {
			log.Fatalf("❌ Failed to generate device token: %v", err)
		}
		setupConfig.BSNDeviceRegistrationTokenEntity = deviceToken
		if !*jsonFlag {
			fmt.Fprintf(os.Stderr, "✅ Device registration token generated successfully!\n")
		}
		if *verboseFlag && !*jsonFlag {
			fmt.Fprintf(os.Stderr, "   Token: %s...\n", deviceToken.Token[:32])
			fmt.Fprintf(os.Stderr, "   Scope: %s\n", deviceToken.Scope)
			fmt.Fprintf(os.Stderr, "   Valid from: %s\n", deviceToken.ValidFrom)
			fmt.Fprintf(os.Stderr, "   Valid to: %s\n", deviceToken.ValidTo)
		}
	} else {
		if !*jsonFlag {
			fmt.Fprintf(os.Stderr, "✅ Using provided device registration token\n")
		}
		if *verboseFlag && !*jsonFlag {
			fmt.Fprintf(os.Stderr, "   Token: %s...\n", setupConfig.BSNDeviceRegistrationTokenEntity.Token[:32])
			fmt.Fprintf(os.Stderr, "   Scope: %s\n", setupConfig.BSNDeviceRegistrationTokenEntity.Scope)
		}
	}
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Set version if not specified
	if setupConfig.Version == "" {
		setupConfig.Version = "3.0.0"
	}

	// Set defaults for optional fields
	if setupConfig.TimeZone == "" {
		setupConfig.TimeZone = "America/New_York"
	}
	if setupConfig.BSNGroupName == "" {
		setupConfig.BSNGroupName = "Default"
	}

	// Step 4: Create B-Deploy setup record
	if !*jsonFlag {
		fmt.Fprintf(os.Stderr, "📝 Creating B-Deploy setup record...\n")
	}

	if *verboseFlag && !*jsonFlag {
		fmt.Fprintf(os.Stderr, "   Setup details:\n")
		fmt.Fprintf(os.Stderr, "     Version: %s\n", setupConfig.Version)
		fmt.Fprintf(os.Stderr, "     Username: %s\n", setupConfig.BDeploy.Username)
		fmt.Fprintf(os.Stderr, "     Network: %s\n", setupConfig.BDeploy.NetworkName)
		fmt.Fprintf(os.Stderr, "     Package: %s\n", setupConfig.BDeploy.PackageName)
		fmt.Fprintf(os.Stderr, "     Type: %s\n", setupConfig.SetupType)
		fmt.Fprintf(os.Stderr, "     Time Zone: %s\n", setupConfig.TimeZone)
		fmt.Fprintf(os.Stderr, "     BSN Group: %s\n", setupConfig.BSNGroupName)
		if setupConfig.BSNDeviceRegistrationTokenEntity != nil {
			fmt.Fprintf(os.Stderr, "     Token Scope: %s\n", setupConfig.BSNDeviceRegistrationTokenEntity.Scope)
		}
		if setupConfig.Network != nil {
			fmt.Fprintf(os.Stderr, "     Interfaces: %d\n", len(setupConfig.Network.Interfaces))
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	response, err := client.BDeploy.AddSetupRecord(ctx, &setupConfig.BDeploySetupRecord)
	if err != nil {
		log.Fatalf("❌ Failed to create B-Deploy setup record: %v", err)
	}

	// Step 5: Display results
	if *jsonFlag {
		// Output as JSON
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(response); err != nil {
			log.Fatalf("Failed to encode JSON: %v", err)
		}
		return
	}

	// Normal output
	fmt.Fprintf(os.Stderr, "✅ B-Deploy setup record created successfully!\n")

	// Debug: Show what we got back
	if *verboseFlag {
		fmt.Fprintf(os.Stderr, "\n🔍 Response Debug:\n")
		fmt.Fprintf(os.Stderr, "   ID: '%s'\n", response.ID)
		fmt.Fprintf(os.Stderr, "   Success: %v\n", response.Success)
		fmt.Fprintf(os.Stderr, "   Error: '%s'\n", response.Error)
	}

	if response.Error != "" {
		fmt.Fprintf(os.Stderr, "   ⚠️  API Error: %s\n", response.Error)
	}

	// Prominently display the setup-id
	fmt.Fprintf(os.Stderr, "\n%s\n", strings.Repeat("=", 70))
	fmt.Fprintf(os.Stderr, "SETUP-ID: %s\n", response.ID)
	fmt.Fprintf(os.Stderr, "%s\n", strings.Repeat("=", 70))
	fmt.Fprintf(os.Stderr, "Save this setup-id - you'll need it to associate players with this setup\n")

	// Generate setup URLs
	fmt.Fprintf(os.Stderr, "\n🔗 Setup URLs:\n")
	fmt.Fprintf(os.Stderr, "   Web: https://provision.bsn.cloud/setup/%s\n", response.ID)
	fmt.Fprintf(os.Stderr, "   API: https://provision.bsn.cloud/rest-setup/v3/setup/%s\n", response.ID)

	fmt.Fprintf(os.Stderr, "\n📱 Next Steps:\n")
	fmt.Fprintf(os.Stderr, "   Option 1: Associate existing players with this setup\n")
	fmt.Fprintf(os.Stderr, "     Use the associate-player-with-setup example with the setup-id above\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "   Option 2: Provision new players directly\n")
	fmt.Fprintf(os.Stderr, "     - Enter the setup URL in the player's web interface\n")
	fmt.Fprintf(os.Stderr, "     - Generate a QR code from the URL and scan it with the player\n")
	fmt.Fprintf(os.Stderr, "     - Load the setup configuration via USB storage\n")

	fmt.Fprintf(os.Stderr, "\n🎉 Setup record created successfully!\n")
}

// loadConfig reads and parses the JSON configuration file
func loadConfig(filename string) (*SetupConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config SetupConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &config, nil
}

// validateConfig checks that required configuration fields are present
func validateConfig(config *gopurple.BDeploySetupRecord) error {
	if config.BDeploy.NetworkName == "" {
		return fmt.Errorf("bDeploy.networkName is required")
	}
	if config.BDeploy.Username == "" {
		return fmt.Errorf("bDeploy.username is required")
	}
	if config.BDeploy.PackageName == "" {
		return fmt.Errorf("bDeploy.packageName is required")
	}
	if config.SetupType == "" {
		return fmt.Errorf("setupType is required")
	}
	// Token is optional - will be auto-generated if not provided
	// Network interfaces are optional for some setup types
	return nil
}
