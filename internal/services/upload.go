package services

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/brightsign/gopurple/internal/auth"
	"github.com/brightsign/gopurple/internal/config"
	"github.com/brightsign/gopurple/internal/errors"
	httpclient "github.com/brightsign/gopurple/internal/http"
	"github.com/brightsign/gopurple/internal/types"
)

const (
	// Upload API path (appended to the configured API base URL)
	uploadAPIPath   = "/Upload/2019/03/REST"
	chunkSize       = 5 * 1024 * 1024 // 5MB chunks
	maxPollAttempts = 60              // Poll for up to 5 minutes (60 * 5 seconds)
	pollInterval    = 5 * time.Second
)

// UploadService defines the interface for content upload operations
type UploadService interface {
	// Upload uploads a file to BSN.cloud using the session-based Upload API
	Upload(ctx context.Context, filePath string, virtualPath string) (*types.UploadResponse, error)

	// UploadReader uploads file content from an io.Reader
	UploadReader(ctx context.Context, fileName string, fileSize int64, reader io.Reader, virtualPath string) (*types.UploadResponse, error)

	// UploadBytes uploads file content from a byte slice
	UploadBytes(ctx context.Context, fileName string, data []byte) (int, error)
}

type uploadService struct {
	config      *config.Config
	httpClient  *httpclient.HTTPClient
	authManager *auth.AuthManager
}

// getUploadBaseURL returns the base URL for upload operations.
// Uses the configured API base URL to support staging/non-production environments.
func (s *uploadService) getUploadBaseURL() string {
	// Use configured base URL, defaulting to production
	baseURL := s.config.BSNBaseURL
	if baseURL == "" {
		baseURL = "https://api.bsn.cloud"
	}
	return baseURL + uploadAPIPath
}

// NewUploadService creates a new upload service
func NewUploadService(cfg *config.Config, httpClient *httpclient.HTTPClient, authManager *auth.AuthManager) UploadService {
	return &uploadService{
		config:      cfg,
		httpClient:  httpClient,
		authManager: authManager,
	}
}

// Upload uploads a file to BSN.cloud using the session-based Upload API
func (s *uploadService) Upload(ctx context.Context, filePath string, virtualPath string) (*types.UploadResponse, error) {
	// Ensure authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if !s.authManager.IsNetworkSet() {
		return nil, fmt.Errorf("network must be set before uploading files")
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	fileName := filepath.Base(filePath)
	fileSize := fileInfo.Size()

	return s.UploadReader(ctx, fileName, fileSize, file, virtualPath)
}

// UploadReader uploads file content from an io.Reader using the session-based Upload API
func (s *uploadService) UploadReader(ctx context.Context, fileName string, fileSize int64, reader io.Reader, virtualPath string) (*types.UploadResponse, error) {
	// Ensure authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if !s.authManager.IsNetworkSet() {
		return nil, fmt.Errorf("network must be set before uploading files")
	}

	// Read file content and calculate SHA1 hash
	if s.config.Debug {
		fmt.Printf("[DEBUG] Reading file and calculating SHA1 hash...\n")
	}
	fmt.Printf("Reading file (%s)...\n", formatBytes(fileSize))
	fileData, sha1Hash, err := readAndHashFile(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	fmt.Printf("SHA1 hash: %s\n", sha1Hash)
	if s.config.Debug {
		fmt.Printf("[DEBUG] File size: %d bytes, SHA1: %s\n", len(fileData), sha1Hash)
	}

	// Step 1: Create upload session
	if s.config.Debug {
		fmt.Printf("[DEBUG] Creating upload session...\n")
	}
	fmt.Printf("Initiating upload session...\n")
	sessionToken, uploadToken, err := s.createUploadSession(ctx, fileName, fileSize, sha1Hash, virtualPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload session: %w", err)
	}
	if s.config.Debug {
		fmt.Printf("[DEBUG] Session created - SessionToken: %s, UploadToken: %s\n", sessionToken, uploadToken)
	}

	// Step 2: Upload file chunks
	if s.config.Debug {
		fmt.Printf("[DEBUG] Uploading file chunks...\n")
	}
	fmt.Printf("Uploading %s in chunks...\n", formatBytes(int64(len(fileData))))
	if err := s.uploadFileChunks(ctx, sessionToken, uploadToken, fileData); err != nil {
		return nil, fmt.Errorf("failed to upload file chunks: %w", err)
	}
	fmt.Printf("All chunks uploaded successfully\n")
	if s.config.Debug {
		fmt.Printf("[DEBUG] All chunks uploaded successfully\n")
	}

	// Step 3: Wait for chunks to be processed (poll until "Uploaded" or "Verified")
	if s.config.Debug {
		fmt.Printf("[DEBUG] Waiting for server to process chunks...\n")
	}
	fmt.Printf("Waiting for server to process chunks...\n")
	if err := s.waitForChunkProcessing(ctx, sessionToken, uploadToken); err != nil {
		return nil, fmt.Errorf("failed waiting for chunk processing: %w", err)
	}
	if s.config.Debug {
		fmt.Printf("[DEBUG] Chunks processed successfully\n")
	}

	// Step 4: Complete the upload with SHA1 hash verification
	if s.config.Debug {
		fmt.Printf("[DEBUG] Completing upload with SHA1 verification...\n")
	}
	fmt.Printf("Completing upload with SHA1 verification...\n")
	status, err := s.completeUpload(ctx, sessionToken, uploadToken, fileName, fileSize, sha1Hash, virtualPath)
	if err != nil {
		return nil, fmt.Errorf("failed to complete upload: %w", err)
	}
	fmt.Printf("Upload completed: ContentID %d\n", status.ContentID)
	if s.config.Debug {
		fmt.Printf("[DEBUG] Upload completed with state: %s, ContentID: %d\n", status.State, status.ContentID)
	}

	// Convert status to UploadResponse
	response := &types.UploadResponse{
		ContentID:      status.ContentID,
		FileName:       status.FileName,
		FileSize:       status.FileSize,
		VirtualPath:    virtualPath,
		FileHash:       status.SHA1Hash,
		UploadComplete: status.State == "Completed" || status.State == "Verified",
		UploadDate:     status.EndTime,
	}

	return response, nil
}

// createUploadSession initiates a new upload and returns the upload token
// This uses the /Sessions/None/Uploads/ endpoint with null tokens
func (s *uploadService) createUploadSession(ctx context.Context, fileName string, fileSize int64, sha1Hash string, virtualPath string) (string, string, error) {
	token, err := s.authManager.GetToken()
	if err != nil {
		return "", "", err
	}

	// POST to /Sessions/None/Uploads/ to initiate upload
	url := fmt.Sprintf("%s/Sessions/None/Uploads/", s.getUploadBaseURL())

	// Prepare upload arguments - only fileName and fileSize are required
	// Other fields are optional but help with deduplication
	args := map[string]interface{}{
		"sessionToken":         nil,
		"uploadToken":          nil,
		"fileName":             fileName,
		"fileSize":             fileSize,
		"virtualPath":          virtualPath,
		"mediaType":            "Auto",
		"fileLastModifiedDate": "0001-01-01T00:00:00",
		"sha1Hash":             "", // Empty initially, provided in final PUT
		"fileThumb":            nil,
		"tags":                 map[string]string{},
	}

	if s.config.Debug {
		fmt.Printf("[DEBUG] POST %s\n", url)
		fmt.Printf("[DEBUG] Request body: %+v\n", args)
	}

	var result types.ContentUploadStatus
	request := s.httpClient.GetClient().R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Content-Type", "application/vnd.bsn.start.content.upload.arguments.201903+json").
		SetHeader("Accept", "application/vnd.bsn.content.upload.status.201903+json,application/vnd.bsn.content.upload.negotiation.status.201903+json,application/vnd.bsn.error+json").
		SetBody(args).
		SetResult(&result)

	resp, err := request.Post(url)
	if err != nil {
		return "", "", errors.NewNetworkError(fmt.Sprintf("POST %s", url), err)
	}

	if s.config.Debug {
		fmt.Printf("[DEBUG] Response status: %d\n", resp.StatusCode())
		fmt.Printf("[DEBUG] Response body: %s\n", string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return "", "", errors.NewAPIError(resp.StatusCode(), "upload_create_failed",
			fmt.Sprintf("Failed to create upload for '%s'", fileName),
			string(resp.Body()))
	}

	if result.UploadToken == "" {
		return "", "", errors.NewAPIError(0, "invalid_response",
			"Server did not return upload token", "")
	}

	// Return "None" as session token (used in URLs) and the actual upload token
	return "None", result.UploadToken, nil
}

// uploadFileChunks uploads file content in chunks using query parameter offset
func (s *uploadService) uploadFileChunks(ctx context.Context, sessionToken, uploadToken string, fileData []byte) error {
	token, err := s.authManager.GetToken()
	if err != nil {
		return err
	}

	totalSize := int64(len(fileData))
	offset := int64(0)

	for offset < totalSize {
		// Calculate chunk size
		currentChunkSize := chunkSize
		if offset+int64(currentChunkSize) > totalSize {
			currentChunkSize = int(totalSize - offset)
		}

		// Extract chunk
		chunk := fileData[offset : offset+int64(currentChunkSize)]

		// Upload chunk using ?offset= query parameter
		url := fmt.Sprintf("%s/Sessions/%s/Uploads/%s/Chunks/?offset=%d", s.getUploadBaseURL(), sessionToken, uploadToken, offset)

		if s.config.Debug && offset == 0 {
			fmt.Printf("[DEBUG] Uploading chunks to: %s\n", url)
		}

		request := s.httpClient.GetClient().R().
			SetContext(ctx).
			SetAuthToken(token).
			SetHeader("Content-Type", "application/octet-stream").
			SetBody(chunk)

		resp, err := request.Post(url)
		if err != nil {
			return errors.NewNetworkError(fmt.Sprintf("POST %s", url), err)
		}

		if !resp.IsSuccess() {
			if s.config.Debug {
				fmt.Printf("[DEBUG] Chunk upload failed at offset %d: %s\n", offset, string(resp.Body()))
			}
			return errors.NewAPIError(resp.StatusCode(), "chunk_upload_failed",
				fmt.Sprintf("Failed to upload chunk at offset %d", offset),
				string(resp.Body()))
		}

		// Show progress for every chunk (not just debug mode)
		progress := float64(offset+int64(currentChunkSize)) / float64(totalSize) * 100
		fmt.Printf("  Progress: %.1f%% (%s / %s)\n",
			progress,
			formatBytes(offset+int64(currentChunkSize)),
			formatBytes(totalSize))

		if s.config.Debug {
			fmt.Printf("[DEBUG] Uploaded chunk: offset=%d size=%d (%.1f%%)\n",
				offset, currentChunkSize, progress)
		}

		offset += int64(currentChunkSize)
	}

	return nil
}

// completeUpload sends a PUT request to finalize the upload with SHA1 hash verification
// Returns the final upload status with ContentID
func (s *uploadService) completeUpload(ctx context.Context, sessionToken, uploadToken, fileName string, fileSize int64, sha1Hash, virtualPath string) (*types.ContentUploadStatus, error) {
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/Sessions/%s/Uploads/%s/", s.getUploadBaseURL(), sessionToken, uploadToken)

	// Send completion arguments with SHA1 hash for verification
	args := map[string]interface{}{
		"sessionToken":         nil,
		"uploadToken":          nil,
		"fileName":             fileName,
		"fileSize":             fileSize,
		"virtualPath":          virtualPath,
		"mediaType":            "Auto",
		"fileLastModifiedDate": "0001-01-01T00:00:00",
		"sha1Hash":             sha1Hash, // Include hash for verification
		"fileThumb":            nil,
		"tags":                 map[string]string{},
	}

	if s.config.Debug {
		fmt.Printf("[DEBUG] PUT %s\n", url)
		fmt.Printf("[DEBUG] Request body: %+v\n", args)
	}

	var result types.ContentUploadStatus
	request := s.httpClient.GetClient().R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Content-Type", "application/vnd.bsn.complete.content.upload.arguments.201903+json").
		SetHeader("Accept", "application/vnd.bsn.content.upload.status.201903+json,application/vnd.bsn.content.upload.negotiation.status.201903+json,application/vnd.bsn.error+json").
		SetBody(args).
		SetResult(&result)

	resp, err := request.Put(url)
	if err != nil {
		return nil, errors.NewNetworkError(fmt.Sprintf("PUT %s", url), err)
	}

	if s.config.Debug {
		fmt.Printf("[DEBUG] Response status: %d\n", resp.StatusCode())
		fmt.Printf("[DEBUG] Response body: %s\n", string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return nil, errors.NewAPIError(resp.StatusCode(), "upload_complete_failed",
			fmt.Sprintf("Failed to complete upload for '%s'", fileName),
			string(resp.Body()))
	}

	// Check for corrupted state (hash mismatch)
	if result.State == "Corrupted" {
		return nil, errors.NewAPIError(0, "hash_verification_failed",
			"File hash verification failed - uploaded file is corrupted", "")
	}

	return &result, nil
}

// waitForChunkProcessing polls the upload status until chunks are processed
// Waits for state to be "Uploaded" or "Verified" before completing the upload
func (s *uploadService) waitForChunkProcessing(ctx context.Context, sessionToken, uploadToken string) error {
	token, err := s.authManager.GetToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/Sessions/%s/Uploads/%s/", s.getUploadBaseURL(), sessionToken, uploadToken)

	for attempt := 0; attempt < maxPollAttempts; attempt++ {
		var status types.ContentUploadStatus
		request := s.httpClient.GetClient().R().
			SetContext(ctx).
			SetAuthToken(token).
			SetHeader("Accept", "application/vnd.bsn.content.upload.status.201903+json").
			SetResult(&status)

		resp, err := request.Get(url)
		if err != nil {
			return errors.NewNetworkError(fmt.Sprintf("GET %s", url), err)
		}

		if !resp.IsSuccess() {
			return errors.NewAPIError(resp.StatusCode(), "status_check_failed",
				"Failed to check upload status",
				string(resp.Body()))
		}

		if s.config.Debug {
			fmt.Printf("[DEBUG] Chunk processing status: %s (attempt %d/%d)\n", status.State, attempt+1, maxPollAttempts)
		}

		// Check if chunks are processed and ready for completion
		switch status.State {
		case "Uploaded", "Verified":
			// Chunks are processed, ready to complete
			return nil
		case "Corrupted", "Cancelled", "Terminated":
			return errors.NewAPIError(0, "upload_failed",
				fmt.Sprintf("Upload failed with state: %s", status.State), "")
		case "Uploading", "Started":
			// Still processing chunks, continue polling
		default:
			// Unknown state
			return errors.NewAPIError(0, "unknown_state",
				fmt.Sprintf("Upload in unknown state: %s", status.State), "")
		}

		// Wait before next poll
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(pollInterval):
			// Continue polling
		}
	}

	return errors.NewAPIError(0, "upload_timeout",
		"Chunk processing did not complete within timeout period", "")
}

// readAndHashFile reads file content and calculates SHA1 hash
func readAndHashFile(reader io.Reader) ([]byte, string, error) {
	hash := sha1.New()
	var buf bytes.Buffer

	// Read and hash simultaneously
	teeReader := io.TeeReader(reader, hash)
	if _, err := io.Copy(&buf, teeReader); err != nil {
		return nil, "", err
	}

	sha1Hash := hex.EncodeToString(hash.Sum(nil))
	return buf.Bytes(), sha1Hash, nil
}

// formatBytes formats byte count as human-readable string
func formatBytes(bytes int64) string {
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

// UploadBytes uploads file content from a byte slice and returns the content ID.
// This is a convenience method that wraps UploadReader for in-memory data.
func (s *uploadService) UploadBytes(ctx context.Context, fileName string, data []byte) (int, error) {
	reader := bytes.NewReader(data)
	resp, err := s.UploadReader(ctx, fileName, int64(len(data)), reader, "")
	if err != nil {
		return 0, err
	}
	return resp.ContentID, nil
}
