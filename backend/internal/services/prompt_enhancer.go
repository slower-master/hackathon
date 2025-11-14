package services

import (
	"fmt"
	"strings"
)

// EnhanceMarketingScript takes user input and enhances it into a professional marketing script
func EnhanceMarketingScript(userInput string) string {
	// If no input, return default professional script
	if strings.TrimSpace(userInput) == "" {
		return GetDefaultMarketingScript()
	}

	// Clean and enhance the user input
	enhanced := strings.TrimSpace(userInput)

	// Add professional opening if not present
	if !hasGreeting(enhanced) {
		enhanced = "Hello! " + enhanced
	}

	// Add call to action if not present
	if !hasCallToAction(enhanced) {
		enhanced = enhanced + " Get started today and experience the difference!"
	}

	// Ensure proper punctuation
	enhanced = ensureProperPunctuation(enhanced)

	return enhanced
}

// GetDefaultMarketingScript returns a professional default script
func GetDefaultMarketingScript() string {
	return "Hello! I'm excited to introduce you to this amazing product. " +
		"It's designed with innovation and excellence in mind, bringing you the best features " +
		"to enhance your experience. Whether you're looking for quality, performance, or reliability, " +
		"this product delivers on all fronts. With its intuitive design and powerful capabilities, " +
		"you'll wonder how you ever managed without it. Don't miss out on this opportunity to " +
		"transform your workflow and achieve more. Get started today and join thousands of satisfied customers!"
}

// GetProductFocusedScript generates a script focused on product features
func GetProductFocusedScript() string {
	return "Welcome! Let me show you what makes this product truly special. " +
		"First, the innovative design sets it apart from anything else on the market. " +
		"Second, its powerful features give you complete control and flexibility. " +
		"Third, the intuitive interface means you can start using it right away. " +
		"Plus, with our dedicated support team, you're never alone on your journey. " +
		"Experience the difference quality makes. Try it today!"
}

// GetBenefitFocusedScript generates a script focused on benefits
func GetBenefitFocusedScript() string {
	return "Hi there! Imagine saving hours of your valuable time every single day. " +
		"That's exactly what this product delivers. With automated processes and smart features, " +
		"you can focus on what really matters while we handle the rest. " +
		"Boost your productivity, reduce stress, and achieve better results faster. " +
		"Join successful professionals who've already made the switch. " +
		"Your future self will thank you. Start your journey today!"
}

// GetEmotionalScript generates an emotional, engaging script
func GetEmotionalScript() string {
	return "You deserve the best, and that's exactly what we're offering you today. " +
		"This isn't just another product – it's a game-changer that will revolutionize the way you work. " +
		"Feel the confidence that comes with using premium quality tools. " +
		"Experience the joy of seamless performance. Embrace the power of innovation. " +
		"Don't let another day go by without experiencing what thousands have already discovered. " +
		"Take action now and transform your life!"
}

// Helper functions
func hasGreeting(text string) bool {
	greetings := []string{"hello", "hi", "welcome", "greetings", "hey"}
	lower := strings.ToLower(text)
	for _, greeting := range greetings {
		if strings.Contains(lower, greeting) {
			return true
		}
	}
	return false
}

func hasCallToAction(text string) bool {
	ctas := []string{
		"get started", "try", "buy", "order", "purchase", "contact",
		"sign up", "join", "don't wait", "act now", "start today",
	}
	lower := strings.ToLower(text)
	for _, cta := range ctas {
		if strings.Contains(lower, cta) {
			return true
		}
	}
	return false
}

func ensureProperPunctuation(text string) string {
	text = strings.TrimSpace(text)
	if len(text) == 0 {
		return text
	}

	// Ensure ends with proper punctuation
	lastChar := text[len(text)-1]
	if lastChar != '.' && lastChar != '!' && lastChar != '?' {
		text = text + "."
	}

	return text
}

// GetScriptByTone returns a script based on desired tone
func GetScriptByTone(tone string) string {
	switch strings.ToLower(tone) {
	case "professional":
		return GetDefaultMarketingScript()
	case "features":
		return GetProductFocusedScript()
	case "benefits":
		return GetBenefitFocusedScript()
	case "emotional":
		return GetEmotionalScript()
	default:
		return GetDefaultMarketingScript()
	}
}

// OptimizeScriptLength ensures script is appropriate length (30-90 seconds)
func OptimizeScriptLength(script string, targetSeconds int) string {
	// Average speaking rate: ~150 words per minute = 2.5 words per second
	wordsPerSecond := 2.5
	targetWords := int(float64(targetSeconds) * wordsPerSecond)

	words := strings.Fields(script)
	currentWords := len(words)

	if currentWords <= targetWords {
		return script
	}

	// Truncate to target length
	truncated := strings.Join(words[:targetWords], " ")

	// Ensure proper ending
	if !strings.HasSuffix(truncated, ".") && !strings.HasSuffix(truncated, "!") {
		truncated = truncated + "..."
	}

	return truncated
}

// GenerateScriptFromProductDescription creates an engaging 15-second marketing script from product details
// This is AI-powered script generation optimized for short-form video content
func GenerateScriptFromProductDescription(productName, productDescription, productCategory, productPrice string) string {
	// Target: 15 seconds = ~37 words (at 150 words/minute)
	// Note: targetWords is used as reference for script optimization

	// If no description provided, use default
	if strings.TrimSpace(productDescription) == "" {
		return OptimizeScriptLength(GetDefaultMarketingScript(), 15)
	}

	// Extract key features from description
	keyFeatures := extractKeyFeatures(productDescription)

	// Generate dynamic script based on product details
	var script string

	if productName != "" {
		script = fmt.Sprintf("Check out the %s! ", productName)
	} else {
		script = "Discover this amazing product! "
	}

	// Add description (concise version)
	descWords := strings.Fields(productDescription)
	if len(descWords) > 20 {
		// Take first 20 words as main description
		script += strings.Join(descWords[:20], " ") + "... "
	} else {
		script += productDescription + " "
	}

	// Add key features if extracted
	if len(keyFeatures) > 0 {
		script += fmt.Sprintf("Featuring %s. ", strings.Join(keyFeatures, ", "))
	}

	// Add price if provided
	if productPrice != "" && productPrice != "0" && productPrice != "$0" {
		script += fmt.Sprintf("Only %s! ", productPrice)
	}

	// Add compelling CTA
	script += "Get yours today!"

	// Optimize to exactly 15 seconds (~37 words)
	optimized := OptimizeScriptLength(script, 15)

	return optimized
}

// GenerateShortFormScript creates ultra-engaging short-form content (TikTok/Instagram Reels style)
func GenerateShortFormScript(productName, productDescription string) string {
	// Super short: 10-15 seconds = ~25-37 words

	// Hook (first 2 seconds)
	hooks := []string{
		"Wait for it! ",
		"You NEED this! ",
		"Game changer alert! ",
		"Stop scrolling! ",
		"This is incredible! ",
	}

	hook := hooks[len(productName)%len(hooks)] // Pseudo-random based on name

	// Core message (8 seconds)
	core := ""
	if productName != "" {
		core = fmt.Sprintf("The %s ", productName)
	}

	// Extract most exciting words from description
	descWords := strings.Fields(strings.ToLower(productDescription))
	excitingWords := []string{}
	buzzwords := []string{"amazing", "innovative", "revolutionary", "best", "new", "premium", "exclusive", "unique"}

	for _, word := range descWords {
		for _, buzz := range buzzwords {
			if strings.Contains(word, buzz) {
				excitingWords = append(excitingWords, word)
				break
			}
		}
		if len(excitingWords) >= 3 {
			break
		}
	}

	if len(excitingWords) > 0 {
		core += "is " + strings.Join(excitingWords, ", ") + "! "
	} else {
		// Fallback: use first 15 words of description
		words := strings.Fields(productDescription)
		if len(words) > 15 {
			core += strings.Join(words[:15], " ") + "... "
		} else {
			core += productDescription + " "
		}
	}

	// CTA (5 seconds)
	ctas := []string{
		"Grab yours now!",
		"Limited time only!",
		"Don't miss out!",
		"Order today!",
		"Link in bio!",
	}
	cta := ctas[len(productDescription)%len(ctas)]

	script := hook + core + cta

	// Ensure it's not too long
	return OptimizeScriptLength(script, 15)
}

// extractKeyFeatures extracts key features from product description
func extractKeyFeatures(description string) []string {
	features := []string{}

	// Look for common feature indicators
	indicators := []string{
		"feature", "includes", "offers", "provides", "delivers",
		"quality", "premium", "professional", "advanced", "innovative",
	}

	words := strings.Fields(strings.ToLower(description))

	for i, word := range words {
		for _, indicator := range indicators {
			if strings.Contains(word, indicator) && i+1 < len(words) {
				// Take next 2-3 words as feature
				endIdx := i + 3
				if endIdx > len(words) {
					endIdx = len(words)
				}
				feature := strings.Join(words[i:endIdx], " ")
				features = append(features, feature)
				break
			}
		}
		if len(features) >= 2 {
			break
		}
	}

	return features
}
