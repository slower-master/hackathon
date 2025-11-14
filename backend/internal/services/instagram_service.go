package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// InstagramService handles Instagram Graph API integration
type InstagramService struct {
	accessToken string
	client      *http.Client
}

// NewInstagramService creates a new Instagram service
func NewInstagramService(accessToken string) *InstagramService {
	return &InstagramService{
		accessToken: accessToken,
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// UploadVideoToInstagram uploads a video to Instagram as a Reel
// Steps:
// 1. Upload video to Instagram's hosting
// 2. Create media container
// 3. Publish the container
func (is *InstagramService) UploadVideoToInstagram(videoPath, caption, instagramUserID string) (string, string, error) {
	fmt.Printf("\n📸 Starting Instagram upload...\n")
	
	// Step 1: Create media container
	containerID, err := is.createMediaContainer(videoPath, caption, instagramUserID)
	if err != nil {
		return "", "", fmt.Errorf("failed to create media container: %v", err)
	}
	
	fmt.Printf("✅ Media container created: %s\n", containerID)
	
	// Step 2: Wait for video processing
	fmt.Printf("⏳ Waiting for video processing...\n")
	time.Sleep(10 * time.Second)
	
	// Step 3: Publish the container
	postID, postURL, err := is.publishMediaContainer(containerID, instagramUserID)
	if err != nil {
		return "", "", fmt.Errorf("failed to publish media: %v", err)
	}
	
	fmt.Printf("🎉 Video published to Instagram!\n")
	fmt.Printf("   Post ID: %s\n", postID)
	fmt.Printf("   Post URL: %s\n", postURL)
	
	return postID, postURL, nil
}

// createMediaContainer creates an Instagram media container for video
func (is *InstagramService) createMediaContainer(videoPath, caption, instagramUserID string) (string, error) {
	// Instagram Graph API endpoint
	apiURL := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/media", instagramUserID)
	
	// Upload video to a public URL first (using temporary hosting service)
	videoURL, err := is.uploadToTemporaryHost(videoPath)
	if err != nil {
		return "", fmt.Errorf("failed to upload to temporary host: %v", err)
	}
	
	fmt.Printf("   Video uploaded to temporary host: %s\n", videoURL)
	
	// Create container payload
	payload := map[string]string{
		"media_type":   "REELS",
		"video_url":    videoURL,
		"caption":      caption,
		"access_token": is.accessToken,
	}
	
	jsonData, _ := json.Marshal(payload)
	
	resp, err := is.client.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()
	
	bodyBytes, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Instagram API error (%s): %s", resp.Status, string(bodyBytes))
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}
	
	containerID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("container ID not found in response")
	}
	
	return containerID, nil
}

// publishMediaContainer publishes the media container to Instagram
func (is *InstagramService) publishMediaContainer(containerID, instagramUserID string) (string, string, error) {
	apiURL := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/media_publish", instagramUserID)
	
	payload := map[string]string{
		"creation_id":  containerID,
		"access_token": is.accessToken,
	}
	
	jsonData, _ := json.Marshal(payload)
	
	resp, err := is.client.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()
	
	bodyBytes, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("Instagram API error (%s): %s", resp.Status, string(bodyBytes))
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", "", fmt.Errorf("failed to parse response: %v", err)
	}
	
	postID, ok := result["id"].(string)
	if !ok {
		return "", "", fmt.Errorf("post ID not found in response")
	}
	
	// Generate Instagram post URL
	postURL := fmt.Sprintf("https://www.instagram.com/p/%s/", postID)
	
	return postID, postURL, nil
}

// uploadToTemporaryHost uploads video to a temporary hosting service
// In production, you should use your own CDN or Instagram's hosting
func (is *InstagramService) uploadToTemporaryHost(videoPath string) (string, error) {
	// Use 0x0.st for temporary hosting
	videoData, err := os.ReadFile(videoPath)
	if err != nil {
		return "", err
	}
	
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	part, err := writer.CreateFormFile("file", "video.mp4")
	if err != nil {
		return "", err
	}
	
	if _, err := part.Write(videoData); err != nil {
		return "", err
	}
	
	writer.Close()
	
	req, err := http.NewRequest("POST", "https://0x0.st", body)
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	resp, err := is.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	bodyBytes, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed (%s): %s", resp.Status, string(bodyBytes))
	}
	
	// 0x0.st returns URL directly
	videoURL := string(bodyBytes)
	
	return videoURL, nil
}

// GenerateInstagramCaption generates an engaging Instagram caption from product details
func GenerateInstagramCaption(productName, productDescription, productPrice string) string {
	caption := ""
	
	// Hook
	hooks := []string{
		"🔥 You NEED to see this!",
		"✨ Game changer alert!",
		"💎 Obsessed with this!",
		"🚀 This is EVERYTHING!",
		"⚡ Wait for it...",
	}
	
	hookIdx := 0
	if len(productName) > 0 {
		hookIdx = len(productName) % len(hooks)
	}
	caption += hooks[hookIdx] + "\n\n"
	
	// Product intro
	if productName != "" {
		caption += fmt.Sprintf("Introducing: %s 🎉\n\n", productName)
	}
	
	// Description (concise)
	if productDescription != "" {
		descWords := len(productDescription)
		if descWords > 100 {
			// Take first 100 characters
			caption += productDescription[:100] + "...\n\n"
		} else {
			caption += productDescription + "\n\n"
		}
	}
	
	// Price
	if productPrice != "" && productPrice != "$0" && productPrice != "0" {
		caption += fmt.Sprintf("💰 Price: %s\n\n", productPrice)
	}
	
	// Hashtags (Instagram best practices: 5-10 relevant hashtags)
	hashtags := []string{
		"#ProductLaunch",
		"#NewProduct",
		"#MustHave",
		"#ShopNow",
		"#Innovation",
		"#TechTok",
		"#ProductReview",
		"#Unboxing",
		"#DealOfTheDay",
		"#TrendingNow",
	}
	
	caption += "\n"
	for i, tag := range hashtags {
		if i >= 8 {
			break // Limit to 8 hashtags
		}
		caption += tag + " "
	}
	
	// Call to action
	caption += "\n\n👉 Link in bio to learn more!"
	
	return caption
}

