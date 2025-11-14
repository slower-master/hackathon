package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GeminiService handles Google Gemini Pro API integration for script generation
type GeminiService struct {
	apiKey string
	client *http.Client
}

// NewGeminiService creates a new Gemini service
func NewGeminiService(apiKey string) *GeminiService {
	return &GeminiService{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateMarketingScript generates a 15-second marketing script using Gemini Pro
func (g *GeminiService) GenerateMarketingScript(productName, productDescription, productCategory, productPrice string) (string, error) {
	fmt.Printf("\n🤖 Generating script with Google Gemini Pro...\n")

	// Build the prompt
	prompt := g.buildPrompt(productName, productDescription, productCategory, productPrice)

	// Call Gemini API
	script, err := g.callGeminiAPI(prompt)
	if err != nil {
		return "", fmt.Errorf("Gemini API call failed: %v", err)
	}

	// Clean up the response
	script = strings.TrimSpace(script)
	script = strings.Trim(script, "\"")

	fmt.Printf("✅ Script generated successfully!\n")
	fmt.Printf("📝 Script: %s\n", script)

	return script, nil
}

// buildPrompt creates an optimized prompt for Gemini
func (g *GeminiService) buildPrompt(productName, productDescription, productCategory, productPrice string) string {
	prompt := `You are an expert marketing copywriter specializing in short-form video scripts for Instagram Reels and TikTok.

Create a compelling 15-second marketing script (approximately 35-40 words) for the following product:

`

	if productName != "" {
		prompt += fmt.Sprintf("Product Name: %s\n", productName)
	}

	if productDescription != "" {
		prompt += fmt.Sprintf("Description: %s\n", productDescription)
	}

	if productCategory != "" {
		prompt += fmt.Sprintf("Category: %s\n", productCategory)
	}

	if productPrice != "" {
		prompt += fmt.Sprintf("Price: %s\n", productPrice)
	}

	prompt += `
REQUIREMENTS:
1. Exactly 35-40 words (for 15-second video)
2. Start with an attention-grabbing hook (first 3-5 words)
3. Mention the product name
4. Highlight 1-2 key features or benefits
5. Include the price if provided
6. End with a strong call-to-action
7. Use energetic, enthusiastic tone
8. Perfect for short-form video (Instagram Reels/TikTok)
9. Don't use quotation marks in the script
10. Make it conversational and natural

OUTPUT FORMAT:
Return ONLY the script text, no additional commentary, no quotation marks, no explanations.

Example output:
Wait for it! The iPhone 15 Pro features the powerful A17 chip, stunning titanium design, and pro camera system. Perfect for creators and professionals. Only $999! Get yours today and experience the future!

Now generate the script:`

	return prompt
}

// callGeminiAPI makes the actual API call to Gemini Pro
func (g *GeminiService) callGeminiAPI(prompt string) (string, error) {
	// Gemini 2.5 Flash API endpoint (per official docs: https://ai.google.dev/gemini-api/docs)
	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s", g.apiKey)

	// Build request payload
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.7,  // Balanced creativity
			"topK":            40,
			"topP":            0.95,
			"maxOutputTokens": 2048, // Increased to 2048 to prevent truncation
			"stopSequences":   []string{},
		},
		"safetySettings": []map[string]interface{}{
			{
				"category":  "HARM_CATEGORY_HARASSMENT",
				"threshold": "BLOCK_MEDIUM_AND_ABOVE",
			},
			{
				"category":  "HARM_CATEGORY_HATE_SPEECH",
				"threshold": "BLOCK_MEDIUM_AND_ABOVE",
			},
			{
				"category":  "HARM_CATEGORY_SEXUALLY_EXPLICIT",
				"threshold": "BLOCK_MEDIUM_AND_ABOVE",
			},
			{
				"category":  "HARM_CATEGORY_DANGEROUS_CONTENT",
				"threshold": "BLOCK_MEDIUM_AND_ABOVE",
			},
		},
	}

	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("📤 Calling Gemini Pro API...\n")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API error (%s): %s", resp.Status, string(bodyBytes))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Log full response for debugging
	fmt.Printf("📥 Gemini API Response Status: %s\n", resp.Status)
	previewLen := 500
	if len(bodyBytes) < previewLen {
		previewLen = len(bodyBytes)
	}
	fmt.Printf("📋 Response preview (first %d chars): %s\n", previewLen, string(bodyBytes[:previewLen]))

	// Extract generated text
	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		// Check for error in response
		if errorData, hasError := result["error"].(map[string]interface{}); hasError {
			errorMsg := fmt.Sprintf("Gemini API error: %v", errorData)
			if message, ok := errorData["message"].(string); ok {
				errorMsg = fmt.Sprintf("Gemini API error: %s", message)
			}
			return "", fmt.Errorf(errorMsg)
		}
		return "", fmt.Errorf("no candidates in response. Full response: %s", string(bodyBytes))
	}

	candidate := candidates[0].(map[string]interface{})
	
	// Check finishReason - but still try to extract text even if MAX_TOKENS
	finishReason := ""
	if fr, ok := candidate["finishReason"].(string); ok {
		finishReason = fr
		if fr == "SAFETY" {
			return "", fmt.Errorf("content was blocked by safety filters. Try adjusting your product description.")
		}
		if fr != "STOP" && fr != "MAX_TOKENS" {
			fmt.Printf("⚠️  Finish reason: %s\n", fr)
		}
	}

	content, ok := candidate["content"].(map[string]interface{})
	if !ok {
		// Log candidate structure for debugging
		candidateJSON, _ := json.MarshalIndent(candidate, "", "  ")
		return "", fmt.Errorf("no content in candidate. Candidate structure: %s", string(candidateJSON))
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		// Log content structure for debugging
		contentJSON, _ := json.MarshalIndent(content, "", "  ")
		return "", fmt.Errorf("no parts in content. Content structure: %s", string(contentJSON))
	}

	part := parts[0].(map[string]interface{})
	text, ok := part["text"].(string)
	if !ok {
		// Log part structure for debugging
		partJSON, _ := json.MarshalIndent(part, "", "  ")
		return "", fmt.Errorf("no text in part. Part structure: %s", string(partJSON))
	}

	// If MAX_TOKENS but we got text, use it (might be truncated but still usable)
	if finishReason == "MAX_TOKENS" {
		fmt.Printf("⚠️  Response was truncated (MAX_TOKENS), but using generated text\n")
	}

	return text, nil
}

// GenerateShortFormScript generates an ultra-short script (10 seconds)
func (g *GeminiService) GenerateShortFormScript(productName, productDescription string) (string, error) {
	prompt := fmt.Sprintf(`Create a 10-second marketing script (25-30 words) for:

Product: %s
Description: %s

Requirements:
- 25-30 words total
- Start with powerful hook
- Ultra-engaging for TikTok/Reels
- Strong call-to-action

Return ONLY the script:`, productName, productDescription)

	return g.callGeminiAPI(prompt)
}

// GenerateInstagramCaption generates an Instagram caption using Gemini
func (g *GeminiService) GenerateInstagramCaption(productName, productDescription, productPrice string) (string, error) {
	prompt := fmt.Sprintf(`Create an engaging Instagram Reels caption for:

Product: %s
Description: %s
Price: %s

Requirements:
- Start with attention-grabbing emoji and hook
- Include product name
- Brief description (2-3 lines)
- Price mention
- 8-10 relevant hashtags
- Call-to-action at the end
- Use emojis strategically
- Perfect for Instagram Reels

Return ONLY the caption:`, productName, productDescription, productPrice)

	return g.callGeminiAPI(prompt)
}

// GenerateWebsiteContent generates website copy using Gemini
func (g *GeminiService) GenerateWebsiteContent(productName, productDescription string) (map[string]string, error) {
	prompt := fmt.Sprintf(`Create website content for:

Product: %s
Description: %s

Generate these sections (return as JSON):
1. hero_title: Compelling hero section title (5-7 words)
2. hero_subtitle: Engaging subtitle (15-20 words)
3. feature_1: First feature benefit (10-15 words)
4. feature_2: Second feature benefit (10-15 words)
5. feature_3: Third feature benefit (10-15 words)
6. feature_4: Fourth feature benefit (10-15 words)
7. cta_text: Call-to-action button text (2-4 words)

Return ONLY valid JSON with these keys:`, productName, productDescription)

	text, err := g.callGeminiAPI(prompt)
	if err != nil {
		return nil, err
	}

	// Try to parse JSON
	var content map[string]string
	if err := json.Unmarshal([]byte(text), &content); err != nil {
		// If not valid JSON, create default structure
		return map[string]string{
			"hero_title":    productName,
			"hero_subtitle": productDescription,
			"feature_1":     "Premium Quality",
			"feature_2":     "Fast Delivery",
			"feature_3":     "Best Price",
			"feature_4":     "Customer Support",
			"cta_text":      "Get Started",
		}, nil
	}

	return content, nil
}

// GenerateWebsiteFeatures generates 4 product features/benefits using Gemini
func (g *GeminiService) GenerateWebsiteFeatures(productName, productDescription, productCategory, productPrice string) ([]map[string]string, error) {
	fmt.Printf("\n🤖 Generating website features with Gemini...\n")
	fmt.Printf("🔍 Input: Name='%s', Desc='%s', Cat='%s', Price='%s'\n", productName, productDescription, productCategory, productPrice)

	prompt := fmt.Sprintf(`Generate 4 product features for: %s (%s). Category: %s, Price: %s

Return ONLY valid JSON:
{
  "features": [
    {"icon": "emoji", "title": "2-4 words", "description": "15-25 words benefit-focused"},
    {"icon": "emoji", "title": "2-4 words", "description": "15-25 words benefit-focused"},
    {"icon": "emoji", "title": "2-4 words", "description": "15-25 words benefit-focused"},
    {"icon": "emoji", "title": "2-4 words", "description": "15-25 words benefit-focused"}
  ]
}

Make features product-specific, use varied emojis (🚀💎🔒⚡🎯✨🌟💪🎨🔥), compelling titles, benefit-focused descriptions.`, productName, productDescription, productCategory, productPrice)

	response, err := g.callGeminiAPI(prompt)
	if err != nil {
		fmt.Printf("❌ Gemini API call failed: %v\n", err)
		return nil, fmt.Errorf("failed to generate features: %v", err)
	}

	fmt.Printf("📥 Raw Gemini response (first 200 chars): %s...\n", response[:min(200, len(response))])

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		fmt.Printf("⚠️  Initial JSON parse failed, trying to clean response...\n")
		// Try to extract JSON from markdown code blocks
		response = strings.TrimSpace(response)
		if strings.HasPrefix(response, "```json") {
			response = strings.TrimPrefix(response, "```json")
			response = strings.TrimSuffix(response, "```")
			response = strings.TrimSpace(response)
		} else if strings.HasPrefix(response, "```") {
			response = strings.TrimPrefix(response, "```")
			response = strings.TrimSuffix(response, "```")
			response = strings.TrimSpace(response)
		}
		
		if err := json.Unmarshal([]byte(response), &result); err != nil {
			fmt.Printf("❌ JSON parse failed after cleaning: %v\n", err)
			fmt.Printf("📄 Cleaned response: %s\n", response)
			return nil, fmt.Errorf("failed to parse features JSON: %v", err)
		}
		fmt.Printf("✅ Successfully cleaned and parsed JSON\n")
	}

	featuresArray, ok := result["features"].([]interface{})
	if !ok {
		fmt.Printf("❌ No 'features' key in response. Keys found: %v\n", getKeys(result))
		return nil, fmt.Errorf("invalid features format in response")
	}
	
	fmt.Printf("✅ Found %d features in response\n", len(featuresArray))

	features := make([]map[string]string, 0, 4)
	for i, f := range featuresArray {
		if i >= 4 {
			break
		}
		featureMap, ok := f.(map[string]interface{})
		if !ok {
			continue
		}

		icon, _ := featureMap["icon"].(string)
		title, _ := featureMap["title"].(string)
		description, _ := featureMap["description"].(string)

		if icon == "" {
			icon = "✨"
		}
		if title == "" {
			title = "Feature"
		}
		if description == "" {
			description = "Experience the difference."
		}

		features = append(features, map[string]string{
			"icon":        icon,
			"title":       title,
			"description": description,
		})
	}

	// Ensure we have exactly 4 features
	if len(features) < 4 {
		defaults := getDefaultFeatures()
		for len(features) < 4 && len(features) < len(defaults) {
			features = append(features, defaults[len(features)])
		}
	}

	fmt.Printf("✅ Generated %d features\n", len(features))
	fmt.Printf("📋 Features details:\n")
	for i, f := range features {
		fmt.Printf("   %d. [%s] %s - %s\n", i+1, f["icon"], f["title"], f["description"][:50])
	}
	return features[:4], nil
}

// Helper function to get map keys for debugging
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

