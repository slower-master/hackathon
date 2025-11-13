package services

import (
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

