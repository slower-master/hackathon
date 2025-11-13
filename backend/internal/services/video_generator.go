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
	req.Header.Set("Authorization", "Basic cmFrZXNoZGQ0NDU0QGdtYWlsLmNvbQ:DEGE6f5zBPjimAmsqg0oL")

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
	// If custom script is provided, enhance it
	if customScript != "" {
		// Add engaging introduction and call-to-action
		enhancedScript := fmt.Sprintf(
			"Hello! I'm excited to share something amazing with you today. %s This product is designed to make your life easier and better. Don't miss out on this incredible opportunity - get yours today and experience the difference!",
			customScript,
		)
		return enhancedScript
	}

	// Default marketing script
	return "Hello! Welcome to our amazing product showcase. This innovative solution is designed specifically for you, combining quality, style, and functionality. It's perfect for anyone looking to upgrade their experience. Join thousands of satisfied customers who have already made the smart choice. Order now and transform the way you live. Don't wait - this is your chance to experience excellence!"
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
			"voice_id": "en-US-JennyNeural", // Professional, clear voice
		},
	}

	payload := map[string]interface{}{
		"source_url": sourceURL,
		"script":     script,
		"config": map[string]interface{}{
			"fluent":    true,
			"pad_audio": 0,
			"stitch":    true, // Better video quality
		},
	}

	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create D-ID request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// D-ID API key (using same auth as other D-ID calls)
	req.Header.Set("Authorization", "Basic cmFrZXNoZGQ0NDU0QGdtYWlsLmNvbQ:DEGE6f5zBPjimAmsqg0oL")

	fmt.Printf("📤 Calling D-ID API for product video...\n")

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
	case "auto":
		// Auto-detect based on image analysis (simple heuristic)
		// For now, use rotation as default (works well for most products)
		promptText = "Professional product showcase with smooth camera movement, elegant rotation, studio lighting, premium commercial feel, 4K quality, product centered"
	default:
		// Default to rotation
		promptText = "Professional product showcase with smooth 360-degree rotation, studio lighting, elegant spin, premium commercial feel, 4K quality, product centered"
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

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create RunwayML request: %v", err)
	}

	// Get RunwayML API key from environment or config
	runwayAPIKey := os.Getenv("RUNWAYML_API_KEY")
	if runwayAPIKey == "" {
		runwayAPIKey = vg.config.RunwayMLAPIKey // We'll add this to config
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+runwayAPIKey)
	req.Header.Set("X-Runway-Version", "2024-11-06")

	fmt.Printf("📤 Calling RunwayML API...\n")

	resp, err := vg.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("RunwayML API request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("RunwayML API error (%s): %s", resp.Status, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse RunwayML response: %v", err)
	}

	// Get task ID for polling
	taskID, ok := result["id"].(string)
	if !ok {
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

	// Poll for up to 5 minutes
	for i := 0; i < 60; i++ {
		time.Sleep(5 * time.Second)

		req, _ := http.NewRequest("GET", apiURL, nil)
		req.Header.Set("Authorization", "Bearer "+runwayAPIKey)
		req.Header.Set("X-Runway-Version", "2024-11-06")

		resp, err := vg.client.Do(req)
		if err != nil {
			fmt.Printf("⚠️  Poll error: %v\n", err)
			continue
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		status, _ := result["status"].(string)
		fmt.Printf("   Poll %d/60: Status = %s\n", i+1, status)

		if status == "SUCCEEDED" {
			// Get video URL
			outputs, ok := result["output"].([]interface{})
			if !ok || len(outputs) == 0 {
				return "", fmt.Errorf("no output in RunwayML response")
			}

			videoURL, ok := outputs[0].(string)
			if !ok {
				return "", fmt.Errorf("invalid video URL format")
			}

			fmt.Printf("✅ Product video ready!\n")
			return vg.downloadVideo(videoURL)
		} else if status == "FAILED" {
			errMsg, _ := result["failure"].(string)
			return "", fmt.Errorf("RunwayML generation failed: %s", errMsg)
		}
	}

	return "", fmt.Errorf("RunwayML generation timeout")
}

// compositeVideosWithShotstack composites avatar and product videos using Shotstack API
// layout: "product_main" (product fullscreen + avatar overlay) or "avatar_main" (avatar fullscreen + product overlay)
func (vg *VideoGenerator) compositeVideosWithShotstack(productVideoPath, avatarVideoPath, layout string) (string, error) {
	fmt.Printf("\n🎨 Compositing videos with Shotstack API...\n")
	fmt.Printf("📐 Layout: %s\n", layout)

	// Upload videos to public URLs (we'll use file.io for temporary hosting)
	productVideoURL, err := vg.uploadToFileIO(productVideoPath)
	if err != nil {
		return "", fmt.Errorf("failed to upload product video: %v", err)
	}

	avatarVideoURL, err := vg.uploadToFileIO(avatarVideoPath)
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar video: %v", err)
	}

	fmt.Printf("📤 Videos uploaded:\n   Product: %s\n   Avatar: %s\n", productVideoURL, avatarVideoURL)

	// Shotstack API endpoint
	apiURL := "https://api.shotstack.io/v1/render"

	// Determine layout based on parameter
	var tracks []interface{}

	if layout == "avatar_main" {
		// Avatar as main (fullscreen), Product as overlay
		fmt.Printf("📐 Using layout: Avatar fullscreen + Product overlay\n")
		tracks = []interface{}{
			// Track 1: Avatar video (main/background - fullscreen)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  avatarVideoURL,
						},
						"start":  0,
						"length": 15,
						"fit":    "cover",
					},
				},
			},
			// Track 2: Product video (overlay in corner)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  productVideoURL,
						},
						"start":    0,
						"length":   15,
						"position": "bottomRight",
						"offset": map[string]interface{}{
							"x": -0.02,
							"y": -0.02,
						},
						"scale": 0.3, // Product takes 30% of screen
					},
				},
			},
		}
	} else {
		// Default: Product as main (fullscreen), Avatar as overlay
		fmt.Printf("📐 Using layout: Product fullscreen + Avatar overlay\n")
		tracks = []interface{}{
			// Track 1: Product video (main/background - fullscreen)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  productVideoURL,
						},
						"start":  0,
						"length": 15,
						"fit":    "cover",
					},
				},
			},
			// Track 2: Avatar video (overlay in corner)
			map[string]interface{}{
				"clips": []interface{}{
					map[string]interface{}{
						"asset": map[string]interface{}{
							"type": "video",
							"src":  avatarVideoURL,
						},
						"start":    0,
						"length":   15,
						"position": "bottomRight",
						"offset": map[string]interface{}{
							"x": -0.02,
							"y": -0.02,
						},
						"scale": 0.25, // Avatar takes 25% of screen
					},
				},
			},
		}
	}

	// Create timeline for picture-in-picture layout
	timeline := map[string]interface{}{
		"timeline": map[string]interface{}{
			"background": "#000000",
			"tracks":     tracks,
		},
		"output": map[string]interface{}{
			"format":     "mp4",
			"resolution": "hd",
			"fps":        30,
		},
	}

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
		return "", fmt.Errorf("Shotstack API request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("Shotstack API error (%s): %s", resp.Status, string(bodyBytes))
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
			return "", fmt.Errorf("Shotstack render failed")
		}
	}

	return "", fmt.Errorf("Shotstack render timeout")
}

// uploadToFileIO uploads a file to file.io for temporary public hosting
func (vg *VideoGenerator) uploadToFileIO(filePath string) (string, error) {
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
		return "", fmt.Errorf("file.io upload failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse file.io response: %v", err)
	}

	if success, ok := result["success"].(bool); !ok || !success {
		return "", fmt.Errorf("file.io upload failed: %v", result)
	}

	fileURL, ok := result["link"].(string)
	if !ok {
		return "", fmt.Errorf("no link in file.io response")
	}

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
		productVideoStyle = "auto"
	}
	if layout == "" {
		layout = "product_main" // Default: product as main, avatar as overlay
	}

	fmt.Printf("📋 Configuration:\n")
	fmt.Printf("   Product Video Style: %s\n", productVideoStyle)
	fmt.Printf("   Layout: %s\n", layout)
	fmt.Printf("   (product_main = product fullscreen + avatar overlay)\n")
	fmt.Printf("   (avatar_main = avatar fullscreen + product overlay)\n\n")

	// Step 1: Generate talking avatar with D-ID
	fmt.Printf("📍 STEP 1/3: Generating Talking Avatar with D-ID\n")
	avatarVideoPath, err := vg.generateAvatarOnly(personMediaPath, customScript)
	if err != nil {
		return "", fmt.Errorf("step 1 failed (D-ID avatar): %v", err)
	}
	fmt.Printf("✅ STEP 1 COMPLETE: Avatar video saved at %s\n\n", avatarVideoPath)

	// Step 2: Generate product video with D-ID
	fmt.Printf("📍 STEP 2/3: Generating Product Video with D-ID\n")
	productVideoPath, err := vg.generateProductVideoWithDID(productImagePath, productVideoStyle)
	if err != nil {
		return "", fmt.Errorf("step 2 failed (D-ID product video): %v", err)
	}
	fmt.Printf("✅ STEP 2 COMPLETE: Product video saved at %s\n\n", productVideoPath)

	// Step 3: Composite videos with Shotstack
	fmt.Printf("📍 STEP 3/3: Compositing Videos with Shotstack\n")
	finalVideoPath, err := vg.compositeVideosWithShotstack(productVideoPath, avatarVideoPath, layout)
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

	// Generate enhanced marketing script
	videoScript := vg.generateMarketingScript(customScript)
	fmt.Printf("🎬 Video Script:\n%s\n\n", videoScript)

	// Script for the video with optimized settings
	script := map[string]interface{}{
		"type":  "text",
		"input": videoScript,
		"provider": map[string]interface{}{
			"type":     "microsoft",
			"voice_id": "en-US-JennyNeural", // Professional, clear voice
		},
	}

	payload := map[string]interface{}{
		"source_url": sourceURL,
		"script":     script,
		"config": map[string]interface{}{
			"fluent":    true,
			"pad_audio": 0,
			"stitch":    true, // Better video quality
		},
	}

	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create D-ID request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic cmFrZXNoZGQ0NDU0QGdtYWlsLmNvbQ:DEGE6f5zBPjimAmsqg0oL")

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
			"voice_id": "en-US-JennyNeural", // Professional, clear voice
		},
	}

	payload := map[string]interface{}{
		"source_url": sourceURL,
		"script":     script,
		"config": map[string]interface{}{
			"fluent":    true,
			"pad_audio": 0,
			"stitch":    true, // Better video quality
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
	req.Header.Set("Authorization", "Basic cmFrZXNoZGQ0NDU0QGdtYWlsLmNvbQ:DEGE6f5zBPjimAmsqg0oL")

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

	for i := 0; i < 60; i++ {
		time.Sleep(5 * time.Second)

		req, _ := http.NewRequest("GET", apiURL, nil)
		// D-ID API - EXACT key from Postman (copied from cURL line 3)
		req.Header.Set("Authorization", "Basic cmFrZXNoZGQ0NDU0QGdtYWlsLmNvbQ:DEGE6f5zBPjimAmsqg0oL")

		fmt.Printf("Polling D-ID task %d/60: %s\n", i+1, talkID)

		resp, err := vg.client.Do(req)
		if err != nil {
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

		if status == "done" {
			videoURL, ok := result["result_url"].(string)
			if !ok {
				return "", fmt.Errorf("result_url not found or invalid")
			}
			return vg.downloadVideo(videoURL)
		} else if status == "error" {
			errMsg := "video generation failed"
			if errField, ok := result["error"].(map[string]interface{}); ok {
				if msg, ok := errField["message"].(string); ok {
					errMsg = msg
				}
			}
			return "", fmt.Errorf("D-ID error: %s", errMsg)
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
