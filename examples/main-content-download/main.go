package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
		idFlag      = flag.Int("id", 0, "Content file ID to download (required)")
		outputFlag  = flag.String("output", "", "Output file path (default: use original filename)")
		infoOnlyFlag = flag.Bool("info-only", false, "Show file info without downloading")
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A command-line tool for downloading content files from BSN.cloud.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Show file info:\n")
		fmt.Fprintf(os.Stderr, "    %s --id 12345 --info-only\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Download file:\n")
		fmt.Fprintf(os.Stderr, "    %s --id 12345\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Download to specific path:\n")
		fmt.Fprintf(os.Stderr, "    %s --id 12345 --output /path/to/file.mp4\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	// Validate required arguments
	if *idFlag == 0 {
		fmt.Fprintf(os.Stderr, "Error: --id must be specified\n\n")
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

	fmt.Println("Creating BSN.cloud client...")

	client, err := gopurple.New(opts...)
	if err != nil {
		if gopurple.IsConfigurationError(err) {
			log.Fatalf("Configuration error: %v", err)
		}
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Authenticate
	fmt.Println("Authenticating with BSN.cloud...")

	if err := client.Authenticate(ctx); err != nil {
		if gopurple.IsAuthenticationError(err) {
			log.Fatalf("Authentication failed: %v", err)
		}
		log.Fatalf("Authentication error: %v", err)
	}

	fmt.Println("Authentication successful!")

	// Handle network selection
	if err := handleNetworkSelection(ctx, client, *networkFlag); err != nil {
		log.Fatalf("Network selection failed: %v", err)
	}

	// Download or show info
	if err := downloadContent(ctx, client, *idFlag, *outputFlag, *infoOnlyFlag, *verboseFlag); err != nil {
		log.Fatalf("Failed: %v", err)
	}
}

func handleNetworkSelection(ctx context.Context, client *gopurple.Client, requestedNetwork string) error {
	// Check if network is already set
	if client.IsNetworkSet() {
		if current, err := client.GetCurrentNetwork(ctx); err == nil {
			fmt.Printf("Using network: %s (ID: %d)\n", current.Name, current.ID)
			return nil
		}
	}

	// If no network flag was provided, check BS_NETWORK environment variable
	if requestedNetwork == "" {
		if envNetwork := os.Getenv("BS_NETWORK"); envNetwork != "" {
			requestedNetwork = envNetwork
			fmt.Printf("Using network from BS_NETWORK environment variable\n")
		}
	}

	// Get available networks
	fmt.Println("Getting available networks...")

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
				fmt.Printf("Using requested network: %s (ID: %d)\n", network.Name, network.ID)
				return client.SetNetworkByID(ctx, network.ID)
			}
		}

		// Network not found - show error and fall back to interactive selection
		fmt.Printf("Network '%s' not found. Available networks:\n", requestedNetwork)
		for i, network := range networks {
			fmt.Printf("  %d. %s (ID: %d)\n", i+1, network.Name, network.ID)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	// If only one network and no specific network requested, use it automatically
	if len(networks) == 1 && requestedNetwork == "" {
		fmt.Printf("Using network: %s (ID: %d)\n", networks[0].Name, networks[0].ID)
		return client.SetNetworkByID(ctx, networks[0].ID)
	}

	// Show available networks and let user choose
	if requestedNetwork == "" {
		fmt.Fprintf(os.Stderr, "Available networks:\n")
		for i, network := range networks {
			fmt.Printf("  %d. %s (ID: %d)\n", i+1, network.Name, network.ID)
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
	fmt.Fprintf(os.Stderr, "Selected network: %s (ID: %d)\n", selectedNetwork.Name, selectedNetwork.ID)

	return client.SetNetworkByID(ctx, selectedNetwork.ID)
}

func downloadContent(ctx context.Context, client *gopurple.Client, id int, outputPath string, infoOnly bool, verbose bool) error {
	// First, get file metadata
	fmt.Printf("\nRetrieving content file %d metadata...\n", id)

	fileInfo, err := client.Content.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Display file information
	fmt.Println("\n=== Content File Details ===")
	fmt.Printf("ID: %d\n", fileInfo.ID)
	fmt.Printf("Name: %s\n", fileInfo.Name)
	fmt.Printf("Type: %s\n", fileInfo.Type)

	if fileInfo.MediaType != "" {
		fmt.Printf("Media Type: %s\n", fileInfo.MediaType)
	}

	if fileInfo.FileSize > 0 {
		fmt.Printf("File Size: %s\n", formatFileSize(fileInfo.FileSize))
	}

	if fileInfo.VirtualPath != "" {
		fmt.Printf("Virtual Path: %s\n", fileInfo.VirtualPath)
	}

	if fileInfo.MimeType != "" {
		fmt.Printf("MIME Type: %s\n", fileInfo.MimeType)
	}

	fmt.Printf("Created: %s\n", fileInfo.CreationDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("Last Modified: %s\n", fileInfo.LastModifiedDate.Format("2006-01-02 15:04:05"))

	if verbose {
		if fileInfo.Hash != "" {
			fmt.Printf("Hash: %s\n", fileInfo.Hash)
		}
		if fileInfo.StorageProvider != "" {
			fmt.Printf("Storage Provider: %s\n", fileInfo.StorageProvider)
		}
		fmt.Printf("Upload Complete: %t\n", fileInfo.UploadComplete)
	}

	// If info-only mode, return here
	if infoOnly {
		return nil
	}

	// Check if it's a folder
	if fileInfo.Type == "Folder" {
		return fmt.Errorf("cannot download folders (ID %d is a folder)", id)
	}

	// Determine output path
	if outputPath == "" {
		outputPath = fileInfo.Name
	}

	// Check if file already exists
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("\nFile '%s' already exists. Overwrite? (y/n): ", outputPath)
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}
		response := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if response != "y" && response != "yes" {
			fmt.Println("Download cancelled.")
			return nil
		}
	}

	// Download the file
	fmt.Printf("\nDownloading '%s'...\n", fileInfo.Name)

	data, err := client.Content.Download(ctx, id)
	if err != nil {
		return err
	}

	// Create output directory if needed
	outputDir := filepath.Dir(outputPath)
	if outputDir != "." && outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("\n✓ Downloaded successfully to: %s\n", outputPath)
	fmt.Printf("  Size: %s (%d bytes)\n", formatFileSize(int64(len(data))), len(data))

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
