package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/brightsign/gopurple"
)

func main() {
	var (
		helpFlag       = flag.Bool("help", false, "Display usage information")
		jsonFlag       = flag.Bool("json", false, "Output as JSON")
		verboseFlag    = flag.Bool("verbose", false, "Show detailed information")
		debugFlag      = flag.Bool("debug", false, "Enable debug logging (shows HTTP requests/responses)")
		timeoutFlag    = flag.Int("timeout", 30, "Request timeout in seconds")
		networkFlag    *string
		virtualPathFlag = flag.String("virtual-path", "", "Virtual path where the file should be placed (e.g., /videos/)")
		filePathFlag   = flag.String("file", "", "Path to the file to upload (required)")
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A command-line tool for uploading content files to a BSN.cloud network.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Upload a video file:\n")
		fmt.Fprintf(os.Stderr, "    %s --file /path/to/video.mp4\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Upload with virtual path:\n")
		fmt.Fprintf(os.Stderr, "    %s --file image.jpg --virtual-path /images/\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Upload and output JSON:\n")
		fmt.Fprintf(os.Stderr, "    %s --file document.pdf --json\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	// Validate required flags
	if *filePathFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: --file flag is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Check if file exists
	if _, err := os.Stat(*filePathFlag); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", *filePathFlag)
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

	// Enable debug mode if requested
	if *debugFlag {
		opts = append(opts, gopurple.WithDebug(true))
	}

	if !*jsonFlag {
		fmt.Println("Creating BSN.cloud client...")
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
	if !*jsonFlag {
		fmt.Println("Authenticating with BSN.cloud...")
	}

	if err := client.Authenticate(ctx); err != nil {
		if gopurple.IsAuthenticationError(err) {
			log.Fatalf("Authentication failed: %v", err)
		}
		log.Fatalf("Authentication error: %v", err)
	}

	if !*jsonFlag {
		fmt.Println("Authentication successful!")
	}

	// Handle network selection
	if err := handleNetworkSelection(ctx, client, *networkFlag, *verboseFlag, *jsonFlag); err != nil {
		log.Fatalf("Network selection failed: %v", err)
	}

	// Upload file
	if err := uploadFile(ctx, client, *filePathFlag, *virtualPathFlag, *jsonFlag, *verboseFlag); err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}
}

func handleNetworkSelection(ctx context.Context, client *gopurple.Client, requestedNetwork string, verbose bool, jsonMode bool) error {
	// Check if network is already set
	if client.IsNetworkSet() {
		if current, err := client.GetCurrentNetwork(ctx); err == nil {
			if !jsonMode {
				fmt.Printf("Using network: %s (ID: %d)\n", current.Name, current.ID)
			}
			return nil
		}
	}

	// If no network flag was provided, check BS_NETWORK environment variable
	if requestedNetwork == "" {
		if envNetwork := os.Getenv("BS_NETWORK"); envNetwork != "" {
			requestedNetwork = envNetwork
			if !jsonMode {
				fmt.Printf("Using network from BS_NETWORK environment variable\n")
			}
		}
	}

	// Get available networks
	if !jsonMode {
		fmt.Println("Getting available networks...")
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
			if network.Name == requestedNetwork {
				if !jsonMode {
					fmt.Printf("Using requested network: %s (ID: %d)\n", network.Name, network.ID)
				}
				return client.SetNetworkByID(ctx, network.ID)
			}
		}
		return fmt.Errorf("network '%s' not found", requestedNetwork)
	}

	// If only one network, use it automatically
	if len(networks) == 1 {
		if !jsonMode {
			fmt.Printf("Using network: %s (ID: %d)\n", networks[0].Name, networks[0].ID)
		}
		return client.SetNetworkByID(ctx, networks[0].ID)
	}

	// Multiple networks available - need user to specify one
	return fmt.Errorf("multiple networks available, please specify one with --network flag")
}

func uploadFile(ctx context.Context, client *gopurple.Client, filePath string, virtualPath string, jsonMode bool, verbose bool) error {
	if !jsonMode {
		fmt.Printf("\nUploading file: %s\n", filepath.Base(filePath))

		// Show file info
		if fileInfo, err := os.Stat(filePath); err == nil {
			fmt.Printf("File size: %s\n", formatFileSize(fileInfo.Size()))
		}

		if virtualPath != "" {
			fmt.Printf("Virtual path: %s\n", virtualPath)
		}

		if verbose {
			fmt.Println("\nUpload process:")
			fmt.Println("  1. Reading file and calculating SHA1 hash...")
			fmt.Println("  2. Initiating upload session...")
			fmt.Println("  3. Uploading file in chunks...")
			fmt.Println("  4. Waiting for server to process chunks...")
			fmt.Println("  5. Completing upload with SHA1 verification...")
			fmt.Fprintf(os.Stderr, "\n")
		}
	}

	// Upload the file
	result, err := client.Upload.Upload(ctx, filePath, virtualPath)
	if err != nil {
		return err
	}

	if jsonMode {
		// Output result as JSON
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	// Display upload result
	fmt.Println("\nUpload successful!")
	fmt.Printf("  Content ID: %d\n", result.ContentID)
	fmt.Printf("  File Name: %s\n", result.FileName)

	if result.FileSize > 0 {
		fmt.Printf("  File Size: %s\n", formatFileSize(result.FileSize))
	}

	if result.VirtualPath != "" {
		fmt.Printf("  Virtual Path: %s\n", result.VirtualPath)
	}

	if result.MediaType != "" {
		fmt.Printf("  Media Type: %s\n", result.MediaType)
	}

	if !result.UploadDate.IsZero() {
		fmt.Printf("  Upload Date: %s\n", result.UploadDate.Format("2006-01-02 15:04:05"))
	}

	fmt.Printf("  Upload Complete: %t\n", result.UploadComplete)

	if verbose && result.FileHash != "" {
		fmt.Printf("  File Hash: %s\n", result.FileHash)
	}

	return nil
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
