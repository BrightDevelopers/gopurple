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
		jsonFlag    = flag.Bool("json", false, "Output as JSON")
		verboseFlag = flag.Bool("verbose", false, "Show detailed information")
		timeoutFlag = flag.Int("timeout", 30, "Request timeout in seconds")
		networkFlag *string
		pageSizeFlag = flag.Int("page-size", 100, "Number of items per page")
		filterFlag   = flag.String("filter", "", "Filter expression (e.g., \"type eq 'Video'\")")
		sortFlag     = flag.String("sort", "", "Sort expression (e.g., \"name asc\")")
		allFlag      = flag.Bool("all", false, "Retrieve all content files (paginate through all results)")
	)

	// Set up network flags to point to the same variable
	networkFlag = flag.String("network", "", "Network name to use (overrides BS_NETWORK)")
	flag.StringVar(networkFlag, "n", "", "Network name to use (overrides BS_NETWORK) [alias for --network]")

	// Custom usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A command-line tool for listing content files on a BSN.cloud network.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  BS_CLIENT_ID        BSN.cloud API client ID (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_SECRET          BSN.cloud API client secret (required)\n")
		fmt.Fprintf(os.Stderr, "  BS_NETWORK         BSN.cloud network name (optional)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  List all content files:\n")
		fmt.Fprintf(os.Stderr, "    %s --all\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  List video files only:\n")
		fmt.Fprintf(os.Stderr, "    %s --filter \"mediaType eq 'Video'\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  List files sorted by name:\n")
		fmt.Fprintf(os.Stderr, "    %s --sort \"name asc\" --all\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Output as JSON:\n")
		fmt.Fprintf(os.Stderr, "    %s --json\n", os.Args[0])
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
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

	// List content files
	if err := listContentFiles(ctx, client, *jsonFlag, *verboseFlag, *pageSizeFlag, *filterFlag, *sortFlag, *allFlag); err != nil {
		log.Fatalf("Failed to list content files: %v", err)
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
			if strings.EqualFold(network.Name, requestedNetwork) {
				if !jsonMode {
					fmt.Printf("Using requested network: %s (ID: %d)\n", network.Name, network.ID)
				}
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
		if !jsonMode {
			fmt.Printf("Using network: %s (ID: %d)\n", networks[0].Name, networks[0].ID)
		}
		return client.SetNetworkByID(ctx, networks[0].ID)
	}

	// Show available networks and let user choose
	if requestedNetwork == "" {
		fmt.Fprintf(os.Stderr, "Available networks:\n")
		for i, network := range networks {
			fmt.Printf("  %d. %s (ID: %d)\n", i+1, network.Name, network.ID)
			if verbose {
				fmt.Printf("     Created: %s, Modified: %s\n",
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
	if !jsonMode {
		if !jsonMode {
		fmt.Fprintf(os.Stderr, "Selected network: %s (ID: %d)\n", selectedNetwork.Name, selectedNetwork.ID)
	}
	}

	return client.SetNetworkByID(ctx, selectedNetwork.ID)
}

func listContentFiles(ctx context.Context, client *gopurple.Client, jsonMode bool, verbose bool, pageSize int, filter string, sort string, allPages bool) error {
	if !jsonMode {
		fmt.Println("\nRetrieving content files...")
	}

	// Build list options
	var listOpts []gopurple.ListOption
	if pageSize > 0 {
		listOpts = append(listOpts, gopurple.WithPageSize(pageSize))
	}
	if filter != "" {
		listOpts = append(listOpts, gopurple.WithFilter(filter))
	}
	if sort != "" {
		listOpts = append(listOpts, gopurple.WithSort(sort))
	}

	var allFiles []gopurple.ContentFile
	var totalRetrieved int

	for {
		result, err := client.Content.List(ctx, listOpts...)
		if err != nil {
			return err
		}

		allFiles = append(allFiles, result.Items...)
		totalRetrieved += len(result.Items)

		if !jsonMode && !allPages {
			fmt.Printf("Retrieved %d content files\n", len(result.Items))
		}

		// If not fetching all pages or no more pages, break
		if !allPages || !result.IsTruncated {
			break
		}

		// Update marker for next page
		if !jsonMode {
			fmt.Printf("Retrieved %d files so far, fetching next page...\n", totalRetrieved)
		}
		listOpts = updateMarker(listOpts, result.NextMarker)
	}

	if jsonMode {
		// Output all files as JSON
		output := gopurple.ContentFileList{
			Items:       allFiles,
			TotalCount:  len(allFiles),
			IsTruncated: false,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(output)
	}

	// Display content files
	fmt.Println("\nContent Files:")
	fmt.Println(strings.Repeat("=", 80))

	if len(allFiles) == 0 {
		fmt.Println("No content files found.")
		return nil
	}

	for i, file := range allFiles {
		// Build display name from available information
		displayName := file.Name
		if displayName == "" {
			// Use media type and ID if name is not available
			if file.MediaType != "" {
				displayName = fmt.Sprintf("%s File (ID: %d)", file.MediaType, file.ID)
			} else {
				displayName = fmt.Sprintf("File ID: %d", file.ID)
			}
		}

		fmt.Printf("\n[%d] %s\n", i+1, displayName)

		// Show type and media type if available
		if file.Type != "" || file.MediaType != "" {
			fmt.Printf("    Type: %s", file.Type)
			if file.MediaType != "" && file.Type != file.MediaType {
				fmt.Printf(" / Media: %s", file.MediaType)
			}
			fmt.Fprintf(os.Stderr, "\n")
		}

		if file.FileSize > 0 {
			fmt.Printf("    Size: %s\n", formatFileSize(file.FileSize))
		}

		// Always show virtual path if available (not just in verbose mode)
		if file.VirtualPath != "" && file.VirtualPath != "\\" {
			fmt.Printf("    Path: %s\n", file.VirtualPath)
		}

		if verbose {
			if file.Path != "" {
				fmt.Printf("    Full Path: %s\n", file.Path)
			}
			if file.FileExtension != "" {
				fmt.Printf("    Extension: %s\n", file.FileExtension)
			}
			if file.MimeType != "" {
				fmt.Printf("    MIME Type: %s\n", file.MimeType)
			}
			if file.Hash != "" {
				fmt.Printf("    Hash: %s\n", file.Hash)
			}
			// Only show creation date if it's not zero
			if !file.CreationDate.IsZero() && file.CreationDate.Year() > 1 {
				fmt.Printf("    Created: %s\n", file.CreationDate.Format("2006-01-02 15:04:05"))
			}
			if !file.LastModifiedDate.IsZero() {
				fmt.Printf("    Modified: %s\n", file.LastModifiedDate.Format("2006-01-02 15:04:05"))
			}
			if file.LastUploadedDate != nil && !file.LastUploadedDate.IsZero() {
				fmt.Printf("    Last Uploaded: %s\n", file.LastUploadedDate.Format("2006-01-02 15:04:05"))
			}
			fmt.Printf("    Upload Complete: %t\n", file.UploadComplete)
			if file.ExternalID != "" {
				fmt.Printf("    External ID: %s\n", file.ExternalID)
			}
		}
	}

	fmt.Printf("\nTotal content files: %d\n", len(allFiles))

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

func updateMarker(opts []gopurple.ListOption, newMarker string) []gopurple.ListOption {
	// Remove old marker if present and add new one
	var newOpts []gopurple.ListOption
	for _, opt := range opts {
		// We can't easily detect if it's a marker option, so we'll just add all and append the new marker
		newOpts = append(newOpts, opt)
	}
	newOpts = append(newOpts, gopurple.WithMarker(newMarker))
	return newOpts
}
