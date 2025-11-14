package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg" // JPEG decoder
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dealshare/hacathon/backend/internal/config"
	"github.com/google/uuid"
	_ "golang.org/x/image/webp" // WebP decoder
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// VideoGenerator handles integration with various AI video generation services
type VideoGenerator struct {
	config *config.Config
	client *http.Client
}

func NewVideoGenerator(cfg *config.Config) *VideoGenerator {
	// Force HTTP/1.1 to match curl behavior
	transport := &http.Transport{
		ForceAttemptHTTP2: false,
	}
	return &VideoGenerator{
		config: cfg,
		client: &http.Client{
			Timeout:   5 * time.Minute,
			Transport: transport,
		},
	}
}

// convertToPNG converts any image format to PNG
func (vg *VideoGenerator) convertToPNG(imagePath string) ([]byte, error) {
	// Open and decode the image
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	// Decode image (supports JPEG, PNG, WEBP, etc.)
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	fmt.Printf("   Original format: %s, converting to PNG\n", format)

	// Encode as PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode as PNG: %v", err)
	}

	return buf.Bytes(), nil
}

// uploadToDID uploads an image to D-ID's hosting and returns the URL
func (vg *VideoGenerator) uploadToDID(imagePath string) (string, error) {
	fmt.Printf("📤 Uploading image to D-ID: %s\n", imagePath)

	// Convert image to PNG format (D-ID accepts PNG/JPEG)
	imageData, err := vg.convertToPNG(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to convert image: %v", err)
	}

	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add image file (always as PNG now)
	filename := strings.TrimSuffix(filepath.Base(imagePath), filepath.Ext(imagePath)) + ".png"
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}
	if _, err := part.Write(imageData); err != nil {
		return "", fmt.Errorf("failed to write image data: %v", err)
	}

	contentType := writer.FormDataContentType()
	writer.Close()

	// Make request to D-ID image upload endpoint
	req, err := http.NewRequest("POST", "https://api.d-id.com/images", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic cmFrZXNoZGQ0NDU0QGdtYWlsLmNvbQ:dK2lCEnxK6fw7PUMUSrJD")

	resp, err := vg.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("D-ID upload request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("D-ID upload failed (%s): %s", resp.Status, string(bodyBytes))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse D-ID response: %v", err)
	}

	// Extract image URL
	imageURL, ok := result["url"].(string)
	if !ok {
		return "", fmt.Errorf("image URL not found in D-ID response: %v", result)
	}

	fmt.Printf("✅ Image uploaded successfully to D-ID: %s\n", imageURL)
	return imageURL, nil
}

// generateMarketingScript creates an enhanced marketing script
func (vg *VideoGenerator) generateMarketingScript(customScript string) string {
	// If custom script is provided, use it as-is (don't overwrite!)
	if customScript != "" && strings.TrimSpace(customScript) != "" {
		return strings.TrimSpace(customScript)
	}

	// Default marketing script - More dynamic and engaging!
	return "Hey there! I'm thrilled to introduce you to something truly special today. This incredible product is a game-changer - it combines cutting-edge innovation with stunning design. Whether you're looking for quality, performance, or style, this product delivers on all fronts. Thousands of customers are already loving it, and I know you will too. Don't miss out on this amazing opportunity. Get yours today and experience the difference for yourself. Trust me, you're going to love it!"
}

// generateProductVideoWithDID generates a product showcase video using D-ID API
// Uses the product image as the source and creates a talking video presentation
func (vg *VideoGenerator) generateProductVideoWithDID(productImagePath, productVideoStyle string) (string, error) {
	fmt.Printf("\n🎬 Generating product video with D-ID...\n")

	// Upload product image to D-ID
	var sourceURL string
	if productImagePath != "" {
		fmt.Printf("📸 Uploading product image to D-ID...\n")
		uploadedURL, err := vg.uploadToDID(productImagePath)
		if err != nil {
			return "", fmt.Errorf("failed to upload product image to D-ID: %v", err)
		}
		sourceURL = uploadedURL
		fmt.Printf("✅ Product image uploaded: %s\n", sourceURL)
	} else {
		return "", fmt.Errorf("product image path is required")
	}

	// Generate product presentation script based on style
	var videoScript string
	switch productVideoStyle {
	case "rotation":
		videoScript = "Welcome to our product showcase! This amazing product features a stunning design with premium quality. Watch as we explore its elegant features and innovative design. Perfect for your needs, this product combines style and functionality in one beautiful package."
	case "zoom":
		videoScript = "Take a closer look at this incredible product! Every detail has been carefully crafted to perfection. From its sleek exterior to its innovative features, this product is designed to impress. Experience the quality and craftsmanship that sets it apart."
	case "pan":
		videoScript = "Let me show you this remarkable product from every angle. Notice the attention to detail and premium materials. This product represents the perfect blend of form and function, designed to exceed your expectations."
	case "reveal":
		videoScript = "Prepare to be amazed by this extraordinary product! With cutting-edge technology and elegant design, this product is truly something special. Discover why it's the perfect choice for you."
	case "auto":
		fallthrough
	default:
		videoScript = "Introducing our premium product! This exceptional item combines innovative design with outstanding quality. Perfect for those who demand the best, this product delivers on every promise. Experience the difference that quality makes."
	}

	fmt.Printf("📹 Product video style: %s\n", productVideoStyle)
	fmt.Printf("🎬 Script: %s\n", videoScript)

	// D-ID API endpoint
	apiURL := "https://api.d-id.com/talks"

	// Create script payload
	script := map[string]interface{}{
		"type":  "text",
		"input": videoScript,
		"provider": map[string]interface{}{
			"type":     "microsoft",
			"voice_id": "en-US-GuyNeural", // Professional, clear male voice
		},
	}

	payload := map[string]interface{}{
		"source_url": sourceURL,
		"script":     script,
		"config": map[string]interface{}{
			"fluent":        true,
			"pad_audio":     0,
			"stitch":        true,  // Better video quality
			"result_format": "mp4", // MP4 format (D-ID only supports mp4/mov)
		},
	}

	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create D-ID request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// D-ID API key (using same auth as other D-ID calls)
	req.Header.Set("Authorization", "Basic cmFrZXNoZGQ0NDU0QGdtYWlsLmNvbQ:dK2lCEnxK6fw7PUMUSrJD")

	fmt.Printf("📤 Calling D-ID API for product video...\n")
	fmt.Printf("   Source URL: %s\n", sourceURL)
	fmt.Printf("   Script length: %d characters\n", len(videoScript))

	resp, err := vg.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("D-ID API request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ D-ID API error response: %s\n", string(bodyBytes))
		return "", fmt.Errorf("D-ID API error (%s): %s", resp.Status, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse D-ID response: %v", err)
	}

	// Get task ID for polling
	talkID, ok := result["id"].(string)
	if !ok {
		fmt.Printf("❌ Unexpected D-ID response: %+v\n", result)
		return "", fmt.Errorf("no task ID in D-ID response: %v", result)
	}

	fmt.Printf("✅ D-ID task created: %s\n", talkID)
	fmt.Printf("   Task status: %v\n", result["status"])
	fmt.Printf("⏳ Waiting for product video generation...\n")

	// Poll for completion
	return vg.pollDIDTask(talkID)
}

// generateProductVideoWithRunwayML generates an animated product showcase video from a static image
// productVideoStyle can be: "rotation", "zoom", "pan", "reveal", "auto" (auto-detects best style)
func (vg *VideoGenerator) generateProductVideoWithRunwayML(productImagePath, productVideoStyle string) (string, error) {
	fmt.Printf("\n🎬 Generating product video with RunwayML Gen-3...\n")

	// Read and encode image as base64 data URI (per RunwayML docs)
	imageData, err := os.ReadFile(productImagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read product image: %v", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// Detect image MIME type
	ext := strings.ToLower(filepath.Ext(productImagePath))
	mimeType := "image/png"
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	case ".webp":
		mimeType = "image/webp"
	}

	// Create data URI (format: data:image/png;base64,...)
	dataURI := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Image)
	fmt.Printf("📸 Image encoded as data URI (%d bytes)\n", len(imageData))

	// Determine product video style and prompt
	var promptText string
	switch productVideoStyle {
	case "rotation":
		promptText = "Professional product showcase with smooth 360-degree rotation, studio lighting, elegant spin, premium commercial feel, 4K quality, product centered"
	case "zoom":
		promptText = "Professional product showcase with smooth zoom-in effect, starting wide and focusing on product details, studio lighting, premium commercial feel, 4K quality"
	case "pan":
		promptText = "Professional product showcase with smooth camera pan movement, exploring product from different angles, studio lighting, premium commercial feel, 4K quality"
	case "reveal":
		promptText = "Professional product reveal with dramatic lighting, product emerging from shadows, cinematic reveal, premium commercial feel, 4K quality"
	case "cinematic":
		// NEW: Cinematic camera movement with dynamic zoom and rotation
		promptText = "Cinematic product showcase with smooth dolly zoom, dramatic camera movement, elegant rotation combined with zoom, professional studio lighting, premium commercial cinematography, 4K Hollywood quality, dynamic composition"
	case "showcase":
		// NEW: Full 360° rotation with dramatic close-ups
		promptText = "Premium product showcase with complete 360-degree rotation, smooth transition to close-up details, studio spotlight lighting, elegant product presentation, 4K commercial quality, professional product photography style"
	case "hero":
		// NEW: Epic hero product shot with dramatic reveal
		promptText = "Epic hero product shot, dramatic zoom out from close-up detail to full product view, cinematic lighting with lens flares, premium commercial feel, 4K Hollywood cinematography, impressive product reveal"
	case "premium":
		// NEW: Slow motion premium rotation
		promptText = "Luxury premium product presentation, slow motion elegant rotation, soft studio lighting with highlights, sophisticated commercial style, 4K high-end quality, refined product showcase"
	case "auto":
		// Auto-detect based on image analysis (simple heuristic)
		// For now, use cinematic as default (most dynamic!)
		promptText = "Cinematic product showcase with smooth camera movement, elegant rotation with zoom, studio lighting, premium commercial feel, 4K quality, dynamic and engaging"
	default:
		// Default to cinematic (most dynamic)
		promptText = "Cinematic product showcase with smooth camera movement, elegant rotation with zoom, studio lighting, premium commercial feel, 4K quality, dynamic and engaging"
	}

	fmt.Printf("📹 Product video style: %s\n", productVideoStyle)
	fmt.Printf("🎬 Prompt: %s\n", promptText)

	// RunwayML Gen-3 API endpoint
	apiURL := "https://api.dev.runwayml.com/v1/image_to_video"

	// Create API payload (per RunwayML docs: https://docs.dev.runwayml.com/guides/using-the-api)
	payload := map[string]interface{}{
		"promptImage": dataURI, // Base64 data URI format
		"model":       "gen3a_turbo",
		"promptText":  promptText,
		"duration":    5,
		"ratio":       "1280:768", // Valid options: "768:1280" (portrait) or "1280:768" (landscape)
	}

	payloadBytes, _ := json.Marshal(payload)

	// Log payload size (not full content to avoid huge logs)
	fmt.Printf("📦 Payload size: %d bytes (image: %d bytes)\n", len(payloadBytes), len(base64Image))

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create RunwayML request: %v", err)
	}

	// Get RunwayML API key from environment or config
	runwayAPIKey := os.Getenv("RUNWAYML_API_KEY")
	if runwayAPIKey == "" {
		runwayAPIKey = vg.config.RunwayMLAPIKey // We'll add this to config
	}

	if runwayAPIKey == "" {
		return "", fmt.Errorf("runwayML API key not configured (set RUNWAYML_API_KEY environment variable)")
	}

	fmt.Printf("🔑 API Key: %s...%s (%d chars)\n", runwayAPIKey[:8], runwayAPIKey[len(runwayAPIKey)-4:], len(runwayAPIKey))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+runwayAPIKey)
	req.Header.Set("X-Runway-Version", "2024-11-06")

	fmt.Printf("📤 Calling RunwayML API: POST %s\n", apiURL)
	fmt.Printf("   Model: %s\n", payload["model"])
	fmt.Printf("   Duration: %v seconds\n", payload["duration"])
	fmt.Printf("   Ratio: %s\n", payload["ratio"])

	resp, err := vg.client.Do(req)
	if err != nil {
		fmt.Printf("❌ RunwayML API request failed: %v\n", err)
		return "", fmt.Errorf("RunwayML API request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	fmt.Printf("📥 Response status: %s\n", resp.Status)
	fmt.Printf("📥 Response body length: %d bytes\n", len(bodyBytes))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("❌ RunwayML API error response:\n")
		fmt.Printf("   Status: %s\n", resp.Status)
		fmt.Printf("   Body: %s\n", string(bodyBytes))

		// Try to parse error response
		var errorResult map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &errorResult); err == nil {
			errorJSON, _ := json.MarshalIndent(errorResult, "", "  ")
			fmt.Printf("   Parsed error:\n%s\n", string(errorJSON))

			// Extract error message
			if msg, ok := errorResult["message"].(string); ok {
				return "", fmt.Errorf("RunwayML API error (%s): %s", resp.Status, msg)
			}
			if errorField, ok := errorResult["error"].(string); ok {
				return "", fmt.Errorf("RunwayML API error (%s): %s", resp.Status, errorField)
			}
		}

		return "", fmt.Errorf("RunwayML API error (%s): %s", resp.Status, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		fmt.Printf("❌ Failed to parse RunwayML response: %v\n", err)
		fmt.Printf("   Response body: %s\n", string(bodyBytes))
		return "", fmt.Errorf("failed to parse RunwayML response: %v", err)
	}

	// Log full response for debugging
	responseJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("📋 RunwayML API response:\n%s\n", string(responseJSON))

	// Get task ID for polling
	taskID, ok := result["id"].(string)
	if !ok {
		fmt.Printf("❌ No task ID in RunwayML response\n")
		fmt.Printf("   Full response: %s\n", string(responseJSON))
		return "", fmt.Errorf("no task ID in RunwayML response: %v", result)
	}

	fmt.Printf("✅ RunwayML task created: %s\n", taskID)
	fmt.Printf("⏳ Waiting for product video generation...\n")

	// Poll for completion
	return vg.pollRunwayMLTask(taskID)
}

// pollRunwayMLTask polls RunwayML for video generation completion
func (vg *VideoGenerator) pollRunwayMLTask(taskID string) (string, error) {
	apiURL := fmt.Sprintf("https://api.dev.runwayml.com/v1/tasks/%s", taskID)
	runwayAPIKey := os.Getenv("RUNWAYML_API_KEY")
	if runwayAPIKey == "" {
		runwayAPIKey = vg.config.RunwayMLAPIKey
	}

	if runwayAPIKey == "" {
		return "", fmt.Errorf("runwayML API key not configured")
	}

	fmt.Printf("🔍 Polling RunwayML task: %s\n", taskID)
	fmt.Printf("   API URL: %s\n", apiURL)

	// Poll for up to 5 minutes
	for i := 0; i < 60; i++ {
		time.Sleep(5 * time.Second)

		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			fmt.Printf("⚠️  Failed to create poll request: %v\n", err)
			continue
		}

		req.Header.Set("Authorization", "Bearer "+runwayAPIKey)
		req.Header.Set("X-Runway-Version", "2024-11-06")

		fmt.Printf("   Poll %d/60: GET %s\n", i+1, apiURL)

		resp, err := vg.client.Do(req)
		if err != nil {
			fmt.Printf("⚠️  Poll request failed: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)

		fmt.Printf("   Response status: %s\n", resp.Status)

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("❌ Poll error (%s): %s\n", resp.Status, string(bodyBytes))
			if i == 0 {
				// On first poll error, return immediately with details
				return "", fmt.Errorf("RunwayML poll failed (%s): %s", resp.Status, string(bodyBytes))
			}
			continue
		}

		var result map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			fmt.Printf("⚠️  Failed to parse poll response: %v\n", err)
			fmt.Printf("   Response body: %s\n", string(bodyBytes))
			continue
		}

		// Log full response for first few polls
		if i < 3 {
			responseJSON, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("📋 Poll response:\n%s\n", string(responseJSON))
		}

		status, ok := result["status"].(string)
		if !ok {
			fmt.Printf("⚠️  No status in response: %v\n", result)
			// Log full response when status is missing
			responseJSON, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("   Full response: %s\n", string(responseJSON))
			continue
		}

		fmt.Printf("   Status: %s\n", status)

		if status == "SUCCEEDED" {
			// Get video URL
			outputs, ok := result["output"].([]interface{})
			if !ok || len(outputs) == 0 {
				responseJSON, _ := json.MarshalIndent(result, "", "  ")
				return "", fmt.Errorf("no output in RunwayML response. Full response: %s", string(responseJSON))
			}

			videoURL, ok := outputs[0].(string)
			if !ok {
				responseJSON, _ := json.MarshalIndent(result, "", "  ")
				return "", fmt.Errorf("invalid video URL format. Full response: %s", string(responseJSON))
			}

			fmt.Printf("✅ Product video ready! URL: %s\n", videoURL)
			return vg.downloadVideo(videoURL)
		} else if status == "FAILED" {
			// Log the entire response for debugging
			fullResponse, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("❌ RunwayML FAILED - Full Response:\n%s\n", string(fullResponse))

			// Try to extract error details
			errMsg := "An unexpected error occurred."
			errDetails := []string{}

			// Check multiple possible error fields
			if failure, ok := result["failure"].(string); ok && failure != "" {
				errMsg = failure
				errDetails = append(errDetails, fmt.Sprintf("failure: %s", failure))
			} else if failureMap, ok := result["failure"].(map[string]interface{}); ok {
				if msg, ok := failureMap["message"].(string); ok {
					errMsg = msg
					errDetails = append(errDetails, fmt.Sprintf("failure.message: %s", msg))
				}
				if code, ok := failureMap["code"].(string); ok {
					errDetails = append(errDetails, fmt.Sprintf("failure.code: %s", code))
				}
				if errType, ok := failureMap["type"].(string); ok {
					errDetails = append(errDetails, fmt.Sprintf("failure.type: %s", errType))
				}
			}

			if reason, ok := result["failureReason"].(string); ok && reason != "" {
				errMsg = reason
				errDetails = append(errDetails, fmt.Sprintf("failureReason: %s", reason))
			}

			if errorField, ok := result["error"].(string); ok && errorField != "" {
				errMsg = errorField
				errDetails = append(errDetails, fmt.Sprintf("error: %s", errorField))
			}

			if errorMap, ok := result["error"].(map[string]interface{}); ok {
				if msg, ok := errorMap["message"].(string); ok {
					errMsg = msg
					errDetails = append(errDetails, fmt.Sprintf("error.message: %s", msg))
				}
			}

			// Log all error details
			if len(errDetails) > 0 {
				fmt.Printf("📋 Error Details:\n")
				for _, detail := range errDetails {
					fmt.Printf("   - %s\n", detail)
				}
			}

			return "", fmt.Errorf("RunwayML generation failed: %s", errMsg)
		} else if status == "PENDING" || status == "IN_PROGRESS" || status == "PROCESSING" {
			// Continue polling
			fmt.Printf("   ⏳ Still processing... (status: %s)\n", status)
		} else {
			// Unknown status
			fmt.Printf("⚠️  Unknown status: %s\n", status)
			responseJSON, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("   Full response: %s\n", string(responseJSON))
		}
	}

	return "", fmt.Errorf("RunwayML generation timeout after 5 minutes")
}

// CompositeVideosWithShotstack composites avatar and product videos using Shotstack API
// layout options:
// - "presenter" (RECOMMENDED): Person left (60%), product right (40%) - looks like real product explanation
// - "split" : Side-by-side 50/50 - balanced, professional
// - "dual_highlight": Person and product side-by-side with borders and highlights - both equally showcased
// - "product_main": Product fullscreen + avatar overlay (traditional)
// - "avatar_main": Avatar fullscreen + product overlay
func (vg *VideoGenerator) CompositeVideosWithShotstack(productVideoPath, avatarVideoPath, layout string) (string, error) {
	fmt.Printf("\n🎨 Compositing videos with Shotstack API...\n")
	fmt.Printf("📐 Layout: %s\n", layout)

	// Shotstack requires publicly accessible URLs
	// Try multiple upload services with fallback
	fmt.Printf("📤 Uploading videos to public hosting...\n")

	productVideoURL, err := vg.uploadWithFallback(productVideoPath)
	if err != nil {
		return "", fmt.Errorf("failed to upload product video: %v", err)
	}

	avatarVideoURL, err := vg.uploadWithFallback(avatarVideoPath)
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar video: %v", err)
	}

	fmt.Printf("📤 Videos uploaded:\n   Product: %s\n   Avatar: %s\n", productVideoURL, avatarVideoURL)

	// Log video file paths for debugging
	fmt.Printf("📁 Local video files:\n   Product: %s\n   Avatar: %s\n", productVideoPath, avatarVideoPath)

	// Verify both videos exist locally
	if _, err := os.Stat(productVideoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("product video file does not exist: %s", productVideoPath)
	}
	if _, err := os.Stat(avatarVideoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("avatar video file does not exist: %s", avatarVideoPath)
	}
	fmt.Printf("✅ Both video files exist locally\n")

	// Verify both URLs are accessible (quick HEAD request)
	fmt.Printf("🔍 Verifying video URLs are accessible...\n")
	if resp, err := vg.client.Head(productVideoURL); err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("⚠️  Warning: Product video URL may not be accessible: %s\n", productVideoURL)
	} else {
		fmt.Printf("✅ Product video URL accessible\n")
	}
	if resp, err := vg.client.Head(avatarVideoURL); err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("⚠️  Warning: Avatar video URL may not be accessible: %s\n", avatarVideoURL)
	} else {
		fmt.Printf("✅ Avatar video URL accessible\n")
	}

	// Shotstack API endpoint
	apiURL := "https://api.shotstack.io/v1/render"

	// Determine layout based on parameter
	var tracks []interface{}

	if layout == "presenter" {
		// PRESENTER LAYOUT - Professional product explanation style
		// Person on left (60%), Product on right (40%) - looks like real presenter explaining product
		fmt.Printf("📐 Using PRESENTER layout: Person (60%% left) + Product (40%% right)\n")
		fmt.Printf("   🎯 This looks like a real person explaining the product!\n")
		tracks = []interface{}{
			// Track 0: Avatar video (LEFT SIDE - 60% width)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  avatarVideoURL,
						},
						"start":    0,
						"length":   15,
						"position": "left",
						"offset": map[string]interface{}{
							"x": 0.15, // Centered in left half
							"y": 0.0,
						},
						"scale":   0.6, // 60% of screen width
						"opacity": 1.0,
						"fit":     "contain",
						"transition": map[string]interface{}{
							"in": "fade",
						},
					},
				},
			},
			// Track 1: Product video (RIGHT SIDE - 40% width) - WITH DYNAMIC ZOOM
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  productVideoURL,
						},
						"start":    0,
						"length":   15,
						"position": "right",
						"offset": map[string]interface{}{
							"x": -0.10, // Centered in right area
							"y": 0.0,
						},
						"scale":   0.4, // 40% of screen
						"opacity": 1.0,
						"fit":     "contain",
						"transition": map[string]interface{}{
							"in": "zoom",
						},
						"effect": "zoomIn",
					},
				},
			},
		}
	} else if layout == "split" {
		// SPLIT SCREEN LAYOUT - Side by side 50/50
		// Equal prominence, professional and balanced
		fmt.Printf("📐 Using SPLIT SCREEN layout: Person (50%% left) + Product (50%% right)\n")
		fmt.Printf("   🎯 Balanced side-by-side presentation\n")
		tracks = []interface{}{
			// Track 0: Avatar video (LEFT HALF)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  avatarVideoURL,
						},
						"start":    0,
						"length":   15,
						"position": "left",
						"offset": map[string]interface{}{
							"x": 0.125, // Quarter from left = centered in left half
							"y": 0.0,
						},
						"scale":   0.5, // 50% of screen
						"opacity": 1.0,
						"fit":     "contain",
						"transition": map[string]interface{}{
							"in": "slideLeft",
						},
					},
				},
			},
			// Track 1: Product video (RIGHT HALF) - SLIDES IN
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  productVideoURL,
						},
						"start":    0,
						"length":   15,
						"position": "right",
						"offset": map[string]interface{}{
							"x": -0.125, // Quarter from right = centered in right half
							"y": 0.0,
						},
						"scale":   0.5, // 50% of screen
						"opacity": 1.0,
						"fit":     "contain",
						"transition": map[string]interface{}{
							"in": "slideRight",
						},
					},
				},
			},
		}
	} else if layout == "dual_highlight" {
		// DUAL HIGHLIGHT LAYOUT - Both person and product highlighted equally
		// Side-by-side with visual borders and effects - integrated, not overlaid
		fmt.Printf("📐 Using DUAL HIGHLIGHT layout: Person + Product both highlighted equally\n")
		fmt.Printf("   🎯 Integrated presentation with visual highlights\n")
		tracks = []interface{}{
			// Track 0: Person border/highlight (LEFT)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "html",
							"html": "<div style='width: 100%; height: 100%; border: 8px solid #FFD700; border-radius: 20px; box-shadow: 0 0 40px rgba(255, 215, 0, 0.8);'></div>",
						},
						"start":    0,
						"length":   15,
						"position": "left",
						"offset": map[string]interface{}{
							"x": 0.13,
							"y": 0.0,
						},
						"scale":   0.48,
						"opacity": 0.9,
					},
				},
			},
			// Track 1: Avatar video (LEFT HALF)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  avatarVideoURL,
						},
						"start":    0,
						"length":   15,
						"position": "left",
						"offset": map[string]interface{}{
							"x": 0.13, // Centered in left half
							"y": 0.0,
						},
						"scale":   0.45, // 45% of screen
						"opacity": 1.0,
						"fit":     "contain",
						"transition": map[string]interface{}{
							"in": "fade",
						},
					},
				},
			},
			// Track 2: Product border/highlight (RIGHT)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "html",
							"html": "<div style='width: 100%; height: 100%; border: 8px solid #00BFFF; border-radius: 20px; box-shadow: 0 0 40px rgba(0, 191, 255, 0.8);'></div>",
						},
						"start":    0,
						"length":   15,
						"position": "right",
						"offset": map[string]interface{}{
							"x": -0.13,
							"y": 0.0,
						},
						"scale":   0.48,
						"opacity": 0.9,
					},
				},
			},
			// Track 3: Product video (RIGHT HALF)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  productVideoURL,
						},
						"start":    0,
						"length":   15,
						"position": "right",
						"offset": map[string]interface{}{
							"x": -0.13, // Centered in right half
							"y": 0.0,
						},
						"scale":   0.45, // 45% of screen
						"opacity": 1.0,
						"fit":     "contain",
						"transition": map[string]interface{}{
							"in": "fade",
						},
					},
				},
			},
		}
	} else if layout == "avatar_main" {
		// AVATAR_MAIN LAYOUT - PERSON FOCUSED with SMOOTH natural integration!
		fmt.Printf("📐 Using PERSON FOCUSED layout: Avatar fullscreen + Product (smooth, centered, natural blend)\n")
		fmt.Printf("   🎯 Person is the star! Product integrated seamlessly\n")
		tracks = []interface{}{
			// Track 0: Product video (CENTERED RIGHT, SMOOTH BLEND!)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  productVideoURL,
						},
						"start":    0,
						"length":   15,
						"position": "center", // CENTER for better control
						"offset": map[string]interface{}{
							"x": 0.0, // TRUE CENTER horizontally!
							"y": 0.0, // TRUE CENTER vertically!
						},
						"scale":   0.40, // 40% for better visibility
						"opacity": 0.92, // Slightly transparent for natural blend
						"fit":     "contain",
						"transition": map[string]interface{}{
							"in": "fade", // Smooth fade in
						},
					},
				},
			},
			// Track 1: Avatar video (BACKGROUND - FULLSCREEN!)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  avatarVideoURL,
						},
						"start":  0,
						"length": 15,
						"fit":    "cover", // Fullscreen background
						"transition": map[string]interface{}{
							"in": "fade",
						},
					},
				},
			},
		}
	} else if layout == "product_main" {
		// PRODUCT_MAIN LAYOUT - Product fullscreen + Person in bottom-right corner
		fmt.Printf("📐 Using PRODUCT MAIN layout: Product fullscreen + Person bottom-right\n")
		fmt.Printf("   🎥 Product: Fullscreen background (cover fit)\n")
		fmt.Printf("   👤 Person: Bottom-right corner overlay (30%% scale, always visible)\n")
		tracks = []interface{}{
			// Track 0: Product video (BACKGROUND LAYER - fullscreen)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  productVideoURL,
						},
						"start":  0.0,
						"length": 15.0,    // Full 15 seconds
						"fit":    "cover", // Fullscreen - fills entire frame
						"transition": map[string]interface{}{
							"in": "fade",
						},
					},
				},
			},
			// Track 1: Avatar video (TOP LAYER - bottom-right corner)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  avatarVideoURL,
						},
						"start":  0.0,
						"length": 15.0, // Full 15 seconds (will show for duration of avatar video)
						"offset": map[string]interface{}{
							"x": 0.35, // Right side (positive = right, 0.35 = 35% from center to right)
							"y": 0.35, // Bottom side (positive = bottom, 0.35 = 35% from center to bottom)
						},
						"scale":   0.30,      // 30% of screen - small corner overlay
						"opacity": 1.0,       // Fully opaque
						"fit":     "contain", // Entire person visible
						"transition": map[string]interface{}{
							"in": "fade", // Smooth fade in
						},
					},
				},
			},
		}
	} else {
		// Default fallback: Same as PRODUCT_MAIN - Product fullscreen + Person in bottom-right corner
		fmt.Printf("📐 Using PRODUCT CENTERED layout: Product fullscreen + Person bottom-right\n")
		fmt.Printf("   🎥 Product: Fullscreen background (cover fit)\n")
		fmt.Printf("   👤 Person: Bottom-right corner overlay (30%% scale, always visible)\n")
		tracks = []interface{}{
			// Track 0: Product video (BACKGROUND LAYER - fullscreen)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  productVideoURL,
						},
						"start":  0.0,
						"length": 15.0,    // Full 15 seconds
						"fit":    "cover", // Fullscreen - fills entire frame
						"transition": map[string]interface{}{
							"in": "fade",
						},
					},
				},
			},
			// Track 1: Avatar video (TOP LAYER - bottom-right corner)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  avatarVideoURL,
						},
						"start":  0.0,
						"length": 15.0, // Full 15 seconds (will show for duration of avatar video)
						"offset": map[string]interface{}{
							"x": 0.35, // Right side (positive = right, 0.35 = 35% from center to right)
							"y": 0.35, // Bottom side (positive = bottom, 0.35 = 35% from center to bottom)
						},
						"scale":   0.30,      // 30% of screen - small corner overlay
						"opacity": 1.0,       // Fully opaque
						"fit":     "contain", // Entire person visible
						"transition": map[string]interface{}{
							"in": "fade", // Smooth fade in
						},
					},
				},
			},
		}
	}

	// Create timeline for picture-in-picture layout
	// NOTE: No background color - let the product video fill the frame naturally
	timeline := map[string]interface{}{
		"timeline": map[string]interface{}{
			"tracks": tracks,
		},
		"output": map[string]interface{}{
			"format":     "mp4",
			"resolution": "hd",
			"fps":        30,
			"quality":    "medium", // Good balance of quality and file size
		},
	}

	// Debug: Log track configuration
	fmt.Printf("🎬 Shotstack Timeline Configuration:\n")
	fmt.Printf("   Total tracks: %d\n", len(tracks))
	if layout == "product_main" {
		fmt.Printf("   ✅ PRODUCT MAIN: Fullscreen product + Person in bottom-right\n")
		fmt.Printf("   Track 0 (bottom): Product video (fullscreen background, 15s)\n")
		fmt.Printf("      URL: %s\n", productVideoURL)
		fmt.Printf("      Fit: cover (fullscreen)\n")
		fmt.Printf("   Track 1 (top): Avatar video (bottom-right corner, 30%% scale, 15s)\n")
		fmt.Printf("      URL: %s\n", avatarVideoURL)
		fmt.Printf("      Offset: x=0.35 (right), y=0.35 (bottom)\n")
		fmt.Printf("      Scale: 0.30 (30%%)\n")
		fmt.Printf("      Opacity: 1.0 (fully visible)\n")
		fmt.Printf("      Fit: contain (preserve aspect ratio)\n")
	} else if layout == "presenter" || layout == "split" || layout == "dual_highlight" || layout == "avatar_main" {
		fmt.Printf("   Layout-specific tracks configured (see above)\n")
	} else {
		fmt.Printf("   Track 0 (bottom): Product video (fullscreen background, 15s)\n")
		fmt.Printf("      URL: %s\n", productVideoURL)
		fmt.Printf("   Track 1 (top): Avatar video (bottom-right corner, 30%% scale, 15s)\n")
		fmt.Printf("      URL: %s\n", avatarVideoURL)
		fmt.Printf("      Offset: x=0.35 (right), y=0.35 (bottom)\n")
	}

	// Log full timeline JSON for debugging
	timelineJSON, _ := json.MarshalIndent(timeline, "", "  ")
	fmt.Printf("\n📋 Full Shotstack Timeline JSON:\n%s\n", string(timelineJSON))

	payloadBytes, _ := json.Marshal(timeline)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create Shotstack request: %v", err)
	}

	shotstackAPIKey := os.Getenv("SHOTSTACK_API_KEY")
	if shotstackAPIKey == "" {
		shotstackAPIKey = vg.config.ShotstackAPIKey
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", shotstackAPIKey)

	fmt.Printf("📤 Calling Shotstack API...\n")

	resp, err := vg.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("shotstack API request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("shotstack API error (%s): %s", resp.Status, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse Shotstack response: %v", err)
	}

	// Get render ID
	response, ok := result["response"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid Shotstack response format")
	}

	renderID, ok := response["id"].(string)
	if !ok {
		return "", fmt.Errorf("no render ID in Shotstack response")
	}

	fmt.Printf("✅ Shotstack render started: %s\n", renderID)
	fmt.Printf("⏳ Waiting for video compositing...\n")

	// Poll for completion
	return vg.pollShotstackRender(renderID)
}

// pollShotstackRender polls Shotstack for render completion
func (vg *VideoGenerator) pollShotstackRender(renderID string) (string, error) {
	apiURL := fmt.Sprintf("https://api.shotstack.io/v1/render/%s", renderID)
	shotstackAPIKey := os.Getenv("SHOTSTACK_API_KEY")
	if shotstackAPIKey == "" {
		shotstackAPIKey = vg.config.ShotstackAPIKey
	}

	// Poll for up to 5 minutes
	for i := 0; i < 60; i++ {
		time.Sleep(5 * time.Second)

		req, _ := http.NewRequest("GET", apiURL, nil)
		req.Header.Set("x-api-key", shotstackAPIKey)

		resp, err := vg.client.Do(req)
		if err != nil {
			fmt.Printf("⚠️  Poll error: %v\n", err)
			continue
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		response, ok := result["response"].(map[string]interface{})
		if !ok {
			continue
		}

		status, _ := response["status"].(string)
		fmt.Printf("   Poll %d/60: Status = %s\n", i+1, status)

		if status == "done" {
			videoURL, ok := response["url"].(string)
			if !ok {
				return "", fmt.Errorf("no video URL in Shotstack response")
			}

			fmt.Printf("✅ Final video ready!\n")
			return vg.downloadVideo(videoURL)
		} else if status == "failed" {
			return "", fmt.Errorf("shotstack render failed")
		}
	}

	return "", fmt.Errorf("shotstack render timeout")
}

// uploadWithFallback tries multiple upload services with fallback
func (vg *VideoGenerator) uploadWithFallback(filePath string) (string, error) {
	fmt.Printf("📤 Trying to upload %s...\n", filepath.Base(filePath))

	// Try tmpfiles.org (no blocks, reliable)
	url, err := vg.uploadToTmpFiles(filePath)
	if err == nil {
		fmt.Printf("✅ Uploaded successfully to tmpfiles.org\n")
		return url, nil
	}
	fmt.Printf("⚠️  tmpfiles.org failed: %v\n", err)

	// Try file.io (alternative)
	url, err = vg.uploadToFileIO_Alternative(filePath)
	if err == nil {
		fmt.Printf("✅ Uploaded successfully to file.io\n")
		return url, nil
	}
	fmt.Printf("⚠️  file.io failed: %v\n", err)

	// Try 0x0.st (might be blocked for your network)
	url, err = vg.uploadToFileIO(filePath)
	if err == nil {
		fmt.Printf("✅ Uploaded successfully to 0x0.st\n")
		return url, nil
	}
	fmt.Printf("⚠️  0x0.st failed: %v\n", err)

	return "", fmt.Errorf("all upload services failed")
}

// uploadToTmpFiles uploads to tmpfiles.org (reliable, no blocks)
func (vg *VideoGenerator) uploadToTmpFiles(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}

	if _, err := part.Write(fileData); err != nil {
		return "", fmt.Errorf("failed to write file data: %v", err)
	}

	writer.Close()

	req, err := http.NewRequest("POST", "https://tmpfiles.org/api/v1/upload", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := vg.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload returned %s: %s", resp.Status, string(bodyBytes[:min(200, len(bodyBytes))]))
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Get the URL from response
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	fileURL, ok := data["url"].(string)
	if !ok {
		return "", fmt.Errorf("no URL in response")
	}

	// tmpfiles.org returns URLs like https://tmpfiles.org/123456
	// Need to convert to direct download: https://tmpfiles.org/dl/123456
	fileURL = strings.Replace(fileURL, "tmpfiles.org/", "tmpfiles.org/dl/", 1)

	return fileURL, nil
}

// uploadToFileIO_Alternative uploads to file.io (alternative service)
func (vg *VideoGenerator) uploadToFileIO_Alternative(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}

	if _, err := part.Write(fileData); err != nil {
		return "", fmt.Errorf("failed to write file data: %v", err)
	}

	writer.Close()

	req, err := http.NewRequest("POST", "https://file.io", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := vg.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload returned %s", resp.Status)
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	success, ok := result["success"].(bool)
	if !ok || !success {
		return "", fmt.Errorf("upload failed")
	}

	fileURL, ok := result["link"].(string)
	if !ok {
		return "", fmt.Errorf("no URL in response")
	}

	return fileURL, nil
}

// uploadToShotstack uploads a file to Shotstack's asset storage
func (vg *VideoGenerator) uploadToShotstack(filePath string) (string, error) {
	fmt.Printf("📤 Uploading %s to Shotstack...\n", filepath.Base(filePath))

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	fmt.Printf("   File size: %.2f MB\n", float64(len(fileData))/(1024*1024))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}

	if _, err := part.Write(fileData); err != nil {
		return "", fmt.Errorf("failed to write file data: %v", err)
	}

	writer.Close()

	req, err := http.NewRequest("POST", "https://api.shotstack.io/v1/assets", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	shotstackAPIKey := os.Getenv("SHOTSTACK_API_KEY")
	if shotstackAPIKey == "" {
		shotstackAPIKey = vg.config.ShotstackAPIKey
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-api-key", shotstackAPIKey)

	resp, err := vg.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("shotstack upload failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("❌ Shotstack upload error (%s): %s\n", resp.Status, string(bodyBytes))
		return "", fmt.Errorf("shotstack returned %s: %s", resp.Status, string(bodyBytes[:min(200, len(bodyBytes))]))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse Shotstack response: %v", err)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid Shotstack response format")
	}

	attributes, ok := data["attributes"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid Shotstack response attributes")
	}

	fileURL, ok := attributes["url"].(string)
	if !ok {
		return "", fmt.Errorf("no URL in Shotstack response")
	}

	fmt.Printf("✅ Upload complete: %s\n", fileURL)
	return fileURL, nil
}

// uploadToFileIO uploads a file to 0x0.st for temporary public hosting
func (vg *VideoGenerator) uploadToFileIO(filePath string) (string, error) {
	fmt.Printf("📤 Uploading %s to 0x0.st...\n", filepath.Base(filePath))

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	fmt.Printf("   File size: %.2f MB\n", float64(len(fileData))/(1024*1024))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}

	if _, err := part.Write(fileData); err != nil {
		return "", fmt.Errorf("failed to write file data: %v", err)
	}

	writer.Close()

	// Use 0x0.st which is more reliable for temporary file hosting
	req, err := http.NewRequest("POST", "https://0x0.st", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := vg.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, _ := io.ReadAll(resp.Body)

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Upload error (%s): %s\n", resp.Status, string(bodyBytes))
		return "", fmt.Errorf("upload returned %s: %s", resp.Status, string(bodyBytes[:min(200, len(bodyBytes))]))
	}

	// 0x0.st returns the URL directly as plain text
	fileURL := strings.TrimSpace(string(bodyBytes))

	if fileURL == "" || !strings.HasPrefix(fileURL, "http") {
		fmt.Printf("❌ Invalid URL response: %s\n", string(bodyBytes))
		return "", fmt.Errorf("invalid URL in response: %s", string(bodyBytes[:min(100, len(bodyBytes))]))
	}

	fmt.Printf("✅ Upload complete: %s\n", fileURL)
	return fileURL, nil
}

// GenerateWithRunwayML generates video using RunwayML Gen-2 API
func (vg *VideoGenerator) GenerateWithRunwayML(productImagePath, personMediaPath, customScript string) (string, error) {
	// RunwayML Gen-2 API integration
	// https://docs.runwayml.com/

	apiURL := "https://api.dev.runwayml.com/v1/generations"

	// Read and encode images
	productImageData, err := os.ReadFile(productImagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read product image: %w", err)
	}

	personMediaData, err := os.ReadFile(personMediaPath)
	if err != nil {
		return "", fmt.Errorf("failed to read person media: %w", err)
	}

	// Prepare request payload
	payload := map[string]interface{}{
		"model":         "gen2",
		"prompt":        "Create a professional product marketing video featuring a person presenting this product. The person should speak enthusiastically about the product's features and benefits. Make it engaging and suitable for social media marketing.",
		"init_image":    fmt.Sprintf("data:image/jpeg;base64,%s", encodeBase64(productImageData)),
		"driving_video": fmt.Sprintf("data:video/mp4;base64,%s", encodeBase64(personMediaData)),
		"duration":      10,
		"width":         1280,
		"height":        720,
	}

	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+vg.config.AIAPIKey)

	resp, err := vg.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call RunwayML API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("RunwayML API error: %s - %s", resp.Status, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Poll for completion
	taskID := result["id"].(string)
	return vg.pollRunwayMLTask(taskID)
}

// GenerateFullAIPipeline orchestrates the complete AI video generation pipeline
// Step 1: D-ID → Generate talking avatar
// Step 2: D-ID → Generate product showcase video
// Step 3: Shotstack → Composite both videos
//
// Parameters:
//   - productImagePath: Path to product image
//   - personMediaPath: Path to presenter image
//   - customScript: Marketing script
//   - productVideoStyle: "rotation", "zoom", "pan", "reveal", "auto" (default: "auto")
//   - layout: "product_main" (product fullscreen, avatar overlay) or "avatar_main" (avatar fullscreen, product overlay) (default: "product_main")
func (vg *VideoGenerator) GenerateFullAIPipeline(productImagePath, personMediaPath, customScript, productVideoStyle, layout string) (string, error) {
	fmt.Printf("\n🚀 ========================================\n")
	fmt.Printf("🚀 FULL AI PIPELINE STARTED\n")
	fmt.Printf("🚀 ========================================\n\n")

	// Set defaults
	if productVideoStyle == "" {
		productVideoStyle = "cinematic" // Default: cinematic for MOST dynamic product showcase
	}
	if layout == "" {
		layout = "product_main" // Default: product-centered with person in bottom-right corner
	}

	fmt.Printf("📋 Configuration:\n")
	fmt.Printf("   Product Video Style: %s\n", productVideoStyle)
	fmt.Printf("   Layout: %s\n", layout)
	fmt.Printf("\n📐 Available Layouts:\n")
	fmt.Printf("   • product_main   : Product fullscreen + Person bottom-right corner - ⭐ RECOMMENDED\n")
	fmt.Printf("   • presenter      : Person (60%%) left + Product (40%%) right - for product explanation\n")
	fmt.Printf("   • split          : Side-by-side 50/50 - balanced, professional\n")
	fmt.Printf("   • dual_highlight : Person + Product both highlighted with borders - integrated, equal showcase\n")
	fmt.Printf("   • avatar_main    : Avatar fullscreen + product overlay\n\n")

	// Step 1: Generate talking avatar with D-ID
	fmt.Printf("📍 STEP 1/3: Generating Talking Avatar with D-ID\n")
	avatarVideoPath, err := vg.generateAvatarOnly(personMediaPath, customScript)
	if err != nil {
		return "", fmt.Errorf("step 1 failed (D-ID avatar): %v", err)
	}
	fmt.Printf("✅ STEP 1 COMPLETE: Avatar video saved at %s\n\n", avatarVideoPath)

	// Step 2: Generate product video - USE ONLY RunwayML
	fmt.Printf("📍 STEP 2/3: Generating Product Video with RunwayML ONLY\n")
	fmt.Printf("   🎬 Using RunwayML Gen-3 for product video generation\n")

	productVideoPath, err := vg.generateProductVideoWithRunwayML(productImagePath, productVideoStyle)
	if err != nil {
		return "", fmt.Errorf("step 2 failed (RunwayML product video): %v", err)
	}
	fmt.Printf("   ✅ RunwayML product video generated successfully\n")
	fmt.Printf("✅ STEP 2 COMPLETE: Product video saved at %s\n\n", productVideoPath)

	// Step 3: Composite videos with Shotstack
	fmt.Printf("📍 STEP 3/3: Compositing Videos with Shotstack\n")
	finalVideoPath, err := vg.CompositeVideosWithShotstack(productVideoPath, avatarVideoPath, layout)
	if err != nil {
		return "", fmt.Errorf("step 3 failed (Shotstack compositing): %v", err)
	}
	fmt.Printf("✅ STEP 3 COMPLETE: Final video saved at %s\n\n", finalVideoPath)

	fmt.Printf("🎉 ========================================\n")
	fmt.Printf("🎉 FULL AI PIPELINE COMPLETED SUCCESSFULLY!\n")
	fmt.Printf("🎉 Final Video: %s\n", finalVideoPath)
	fmt.Printf("🎉 ========================================\n\n")

	return finalVideoPath, nil
}

// generateAvatarOnly generates just the talking avatar video (used in pipeline)
func (vg *VideoGenerator) generateAvatarOnly(personMediaPath, customScript string) (string, error) {
	fmt.Printf("🎬 Generating talking avatar with D-ID API...\n")

	apiURL := "https://api.d-id.com/talks"

	// Upload person/presenter image to D-ID for public URL
	var sourceURL string
	if personMediaPath != "" && (strings.HasSuffix(strings.ToLower(personMediaPath), ".png") ||
		strings.HasSuffix(strings.ToLower(personMediaPath), ".jpg") ||
		strings.HasSuffix(strings.ToLower(personMediaPath), ".jpeg") ||
		strings.HasSuffix(strings.ToLower(personMediaPath), ".webp")) {
		fmt.Printf("📸 Uploading YOUR presenter image to D-ID: %s\n", personMediaPath)
		uploadedURL, err := vg.uploadToDID(personMediaPath)
		if err != nil {
			fmt.Printf("❌ ERROR: D-ID upload failed: %v\n", err)
			fmt.Printf("⚠️  This is why default girl (Noelle) is showing! Fix the upload issue.\n")
			fmt.Printf("📝 Falling back to default presenter for now...\n")
			sourceURL = "https://create-images-results.d-id.com/api_docs/assets/noelle.jpeg"
		} else {
			sourceURL = uploadedURL
			fmt.Printf("✅ SUCCESS! Using YOUR custom presenter: %s\n", sourceURL)
		}
	} else {
		fmt.Printf("⚠️  WARNING: No valid image file at: %s\n", personMediaPath)
		fmt.Printf("   Accepted formats: .png, .jpg, .jpeg, .webp\n")
		fmt.Printf("📝 Using D-ID default presenter (Noelle)\n")
		sourceURL = "https://create-images-results.d-id.com/api_docs/assets/noelle.jpeg"
	}

	// Generate enhanced marketing script
	videoScript := vg.generateMarketingScript(customScript)
	fmt.Printf("🎬 Video Script:\n%s\n\n", videoScript)

	// Script for the video with optimized settings
	script := map[string]interface{}{
		"type":  "text",
		"input": videoScript,
		"provider": map[string]interface{}{
			"type":     "microsoft",
			"voice_id": "en-US-GuyNeural", // Professional, clear male voice
		},
	}

	payload := map[string]interface{}{
		"source_url": sourceURL,
		"script":     script,
		"config": map[string]interface{}{
			"fluent":        true,
			"pad_audio":     0,
			"stitch":        true,  // Better video quality
			"result_format": "mp4", // MP4 format (D-ID only supports mp4/mov)
		},
	}

	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create D-ID request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic cmFrZXNoZGQ0NDU0QGdtYWlsLmNvbQ:dK2lCEnxK6fw7PUMUSrJD")

	fmt.Printf("📤 Calling D-ID API for avatar video...\n")

	resp, err := vg.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("D-ID API request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("D-ID API error (%s): %s", resp.Status, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse D-ID response: %v", err)
	}

	// Get task ID for polling
	talkID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("no task ID in D-ID response: %v", result)
	}

	fmt.Printf("✅ D-ID task created: %s\n", talkID)
	fmt.Printf("⏳ Waiting for avatar video generation...\n")

	// Poll for completion
	return vg.pollDIDTask(talkID)
}

// GenerateWithDID generates video using D-ID API (Talking Head)
// This is the main entry point - it will use full AI pipeline if configured
// productVideoStyle: "rotation", "zoom", "pan", "reveal", "auto" (default: "auto")
// layout: "product_main" or "avatar_main" (default: "product_main")
func (vg *VideoGenerator) GenerateWithDID(productImagePath, personMediaPath, customScript, productVideoStyle, layout string) (string, error) {
	// Check if full AI pipeline is enabled
	if vg.config.UseFullAIPipeline {
		fmt.Printf("🎯 Full AI Pipeline enabled! Using D-ID + RunwayML + Shotstack\n")
		return vg.GenerateFullAIPipeline(productImagePath, personMediaPath, customScript, productVideoStyle, layout)
	}

	// Otherwise use the original single D-ID video generation
	fmt.Printf("📝 Using standard D-ID video generation (avatar only)\n")
	return vg.generateAvatarOnly(personMediaPath, customScript)
}

// LEGACY: Original GenerateWithDID implementation (kept for reference)
func (vg *VideoGenerator) GenerateWithDIDLegacy(productImagePath, personMediaPath, customScript string) (string, error) {
	// D-ID API integration for talking head videos
	// https://docs.d-id.com/

	apiURL := "https://api.d-id.com/talks"

	// Upload person/presenter image to D-ID for public URL
	var sourceURL string
	if personMediaPath != "" && (strings.HasSuffix(strings.ToLower(personMediaPath), ".png") ||
		strings.HasSuffix(strings.ToLower(personMediaPath), ".jpg") ||
		strings.HasSuffix(strings.ToLower(personMediaPath), ".jpeg")) {
		fmt.Printf("📸 Uploading presenter image to D-ID...\n")
		uploadedURL, err := vg.uploadToDID(personMediaPath)
		if err != nil {
			fmt.Printf("⚠️  D-ID upload failed: %v. Using default presenter.\n", err)
			sourceURL = "https://create-images-results.d-id.com/api_docs/assets/noelle.jpeg"
		} else {
			sourceURL = uploadedURL
			fmt.Printf("✅ Using custom presenter: %s\n", sourceURL)
		}
	} else {
		fmt.Printf("📝 Using D-ID default presenter\n")
		sourceURL = "https://create-images-results.d-id.com/api_docs/assets/noelle.jpeg"
	}

	// Upload product image to D-ID (for future use - could be overlaid or used in multi-shot videos)
	var productImageURL string
	if productImagePath != "" {
		fmt.Printf("🏷️  Uploading product image to D-ID...\n")
		uploadedProductURL, err := vg.uploadToDID(productImagePath)
		if err != nil {
			fmt.Printf("⚠️  Product image upload failed: %v\n", err)
		} else {
			productImageURL = uploadedProductURL
			fmt.Printf("✅ Product image available at: %s\n", productImageURL)
		}
	}

	// Generate enhanced marketing script
	videoScript := vg.generateMarketingScript(customScript)

	// Add product details if available
	if productImageURL != "" {
		videoScript += " You can see the product details and more at our website. Check out the visuals and features that make this product stand out from the rest!"
	}

	fmt.Printf("\n🎬 Video Script:\n%s\n\n", videoScript)

	// Script for the video with optimized settings
	script := map[string]interface{}{
		"type":  "text",
		"input": videoScript,
		"provider": map[string]interface{}{
			"type":     "microsoft",
			"voice_id": "en-US-GuyNeural", // Professional, clear male voice
		},
	}

	payload := map[string]interface{}{
		"source_url": sourceURL,
		"script":     script,
		"config": map[string]interface{}{
			"fluent":        true,
			"pad_audio":     0,
			"stitch":        true,  // Better video quality
			"result_format": "mp4", // MP4 format (D-ID only supports mp4/mov)
		},
	}

	payloadBytes, _ := json.Marshal(payload)

	// Log the curl command for debugging
	fmt.Printf("\n=== D-ID API Request ===\n")
	fmt.Printf("URL: %s\n", apiURL)
	fmt.Printf("Method: POST\n")
	fmt.Printf("Payload (first 200 chars): %s...\n", string(payloadBytes)[:min(200, len(payloadBytes))])

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "curl/7.88.1")
	req.Header.Set("Accept", "*/*")
	// D-ID API - EXACT key from Postman (copied from cURL line 3)
	req.Header.Set("Authorization", "Basic cmFrZXNoZGQ0NDU0QGdtYWlsLmNvbQ:dK2lCEnxK6fw7PUMUSrJD")

	// Log curl equivalent command
	fmt.Printf("Curl equivalent:\n")
	fmt.Printf("curl -X POST '%s' \\\n", apiURL)
	fmt.Printf("  -H 'Content-Type: application/json' \\\n")
	fmt.Printf("  -H 'Authorization: Basic %s' \\\n", vg.config.AIAPIKey)
	fmt.Printf("  -d '%s'\n", string(payloadBytes))
	fmt.Printf("=====================\n\n")

	resp, err := vg.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call D-ID API: %w", err)
	}
	defer resp.Body.Close()

	// Log response details
	fmt.Printf("=== D-ID API Response ===\n")
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Status Code: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error Body: %s\n", string(body))
		fmt.Printf("=====================\n\n")
		return "", fmt.Errorf("D-ID API error: %s - %s", resp.Status, string(body))
	}

	fmt.Printf("Success!\n")
	fmt.Printf("=====================\n\n")

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Poll for completion
	talkID := result["id"].(string)
	return vg.pollDIDTask(talkID)
}

// GenerateWithSynthesia generates video using Synthesia API
func (vg *VideoGenerator) GenerateWithSynthesia(productImagePath, customScript string) (string, error) {
	// Synthesia API integration
	// https://docs.synthesia.io/

	apiURL := "https://api.synthesia.io/v2/videos"

	payload := map[string]interface{}{
		"test": false,
		"input": []map[string]interface{}{
			{
				"scriptText": "Welcome to our product showcase! Today, I'm excited to introduce you to an incredible product that will transform the way you work. This innovative solution combines cutting-edge technology with user-friendly design, making it perfect for both professionals and everyday users. Let me show you what makes it special!",
				"avatar":     "anna_costume1_cameraA",
				"background": "green_screen",
			},
		},
	}

	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", vg.config.AIAPIKey)

	resp, err := vg.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Synthesia API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("synthesia API error: %s - %s", resp.Status, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Poll for completion
	videoID := result["id"].(string)
	return vg.pollSynthesiaTask(videoID)
}

// pollDIDTask polls D-ID for video generation completion
func (vg *VideoGenerator) pollDIDTask(talkID string) (string, error) {
	apiURL := fmt.Sprintf("https://api.d-id.com/talks/%s", talkID)
	// Use same authorization as initial API call
	authHeader := "Basic cmFrZXNoZGQ0NDU0QGdtYWlsLmNvbQ:dK2lCEnxK6fw7PUMUSrJD"

	for i := 0; i < 60; i++ {
		time.Sleep(5 * time.Second)

		req, _ := http.NewRequest("GET", apiURL, nil)
		req.Header.Set("Authorization", authHeader)

		fmt.Printf("Polling D-ID task %d/60: %s\n", i+1, talkID)

		resp, err := vg.client.Do(req)
		if err != nil {
			fmt.Printf("Poll request error: %v\n", err)
			continue
		}

		// Check for authentication errors
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return "", fmt.Errorf("D-ID authentication failed (%s): %s", resp.Status, string(bodyBytes))
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			fmt.Printf("Poll HTTP error (%s): %s\n", resp.Status, string(bodyBytes))
			continue
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			fmt.Printf("Poll response decode error: %v\n", err)
			continue
		}
		resp.Body.Close()

		fmt.Printf("Poll response: %+v\n", result)

		// Check if status exists
		statusVal, exists := result["status"]
		if !exists {
			fmt.Printf("Status field missing in response\n")
			continue
		}

		status, ok := statusVal.(string)
		if !ok {
			fmt.Printf("Status is not a string: %v\n", statusVal)
			continue
		}

		fmt.Printf("Task status: %s\n", status)

		// Check for completion statuses (D-ID uses various status values)
		if status == "done" || status == "completed" || status == "succeeded" || status == "ready" {
			videoURL, ok := result["result_url"].(string)
			if !ok {
				// Try alternative field names
				if url, ok := result["url"].(string); ok {
					videoURL = url
				} else if url, ok := result["video_url"].(string); ok {
					videoURL = url
				} else {
					return "", fmt.Errorf("result_url not found or invalid in response: %+v", result)
				}
			}
			fmt.Printf("✅ Video ready! URL: %s\n", videoURL)
			return vg.downloadVideo(videoURL)
		}

		// Check if result_url exists even if status is not "done" (sometimes D-ID returns it early)
		if videoURL, ok := result["result_url"].(string); ok && videoURL != "" {
			fmt.Printf("✅ Video URL found even though status is '%s'. Downloading...\n", status)
			return vg.downloadVideo(videoURL)
		}

		// Check for error status
		if status == "error" || status == "failed" {
			// Log full error response for debugging
			fmt.Printf("❌ D-ID task failed. Full response: %+v\n", result)

			// Try to extract detailed error information
			errMsg := "video generation failed"
			errDetails := []string{}

			// Check for error object
			if errField, ok := result["error"].(map[string]interface{}); ok {
				if msg, ok := errField["message"].(string); ok {
					errMsg = msg
				}
				if code, ok := errField["code"].(string); ok {
					errDetails = append(errDetails, fmt.Sprintf("code: %s", code))
				}
				if details, ok := errField["details"].(string); ok {
					errDetails = append(errDetails, fmt.Sprintf("details: %s", details))
				}
			}

			// Check top-level error fields
			if msg, ok := result["message"].(string); ok && errMsg == "video generation failed" {
				errMsg = msg
			}
			if failure, ok := result["failure"].(string); ok {
				errDetails = append(errDetails, fmt.Sprintf("failure: %s", failure))
			}
			if reason, ok := result["reason"].(string); ok {
				errDetails = append(errDetails, fmt.Sprintf("reason: %s", reason))
			}

			// Build comprehensive error message
			fullErrMsg := errMsg
			if len(errDetails) > 0 {
				fullErrMsg = fmt.Sprintf("%s (%s)", errMsg, strings.Join(errDetails, ", "))
			}

			// Include task ID in error for debugging
			return "", fmt.Errorf("D-ID error (task %s): %s", talkID, fullErrMsg)
		}
	}

	// Timeout reached - check if video file was created anyway
	fmt.Printf("⚠️  Polling timeout reached. Checking if video file exists...\n")

	// Try to find video file by task ID pattern in generated folder
	files, err := filepath.Glob(filepath.Join(vg.config.GeneratedVideoPath, "*.mp4"))
	if err == nil {
		// Check the most recently modified file (likely the one we just generated)
		var latestFile string
		var latestTime time.Time
		for _, file := range files {
			if info, err := os.Stat(file); err == nil {
				if info.ModTime().After(latestTime) {
					latestTime = info.ModTime()
					latestFile = file
				}
			}
		}
		// If file was modified in last 10 minutes, assume it's our video
		if latestFile != "" && time.Since(latestTime) < 10*time.Minute {
			fmt.Printf("✅ Found recently generated video file: %s\n", latestFile)
			return latestFile, nil
		}
	}

	return "", fmt.Errorf("video generation timeout")
}

func (vg *VideoGenerator) pollSynthesiaTask(videoID string) (string, error) {
	apiURL := fmt.Sprintf("https://api.synthesia.io/v2/videos/%s", videoID)

	for i := 0; i < 60; i++ {
		time.Sleep(10 * time.Second)

		req, _ := http.NewRequest("GET", apiURL, nil)
		req.Header.Set("Authorization", vg.config.AIAPIKey)

		resp, err := vg.client.Do(req)
		if err != nil {
			continue
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		status := result["status"].(string)
		if status == "complete" {
			videoURL := result["download"].(string)
			return vg.downloadVideo(videoURL)
		} else if status == "failed" {
			return "", fmt.Errorf("video generation failed")
		}
	}

	return "", fmt.Errorf("video generation timeout")
}

func (vg *VideoGenerator) downloadVideo(url string) (string, error) {
	resp, err := vg.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create output file
	os.MkdirAll(vg.config.GeneratedVideoPath, 0755)
	videoID := uuid.New().String()
	outputPath := filepath.Join(vg.config.GeneratedVideoPath, fmt.Sprintf("%s.mp4", videoID))

	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
