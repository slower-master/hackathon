package services

import "fmt"

// SIMPLIFIED Shotstack compositing for product_main layout
// Product fullscreen + Person bottom-right WITH LOOPING
func simpleProductMainTracks(productVideoURL, avatarVideoURL string, avatarDuration float64) []interface{} {
	// Create avatar clips to loop for 15 seconds
	avatarClips := []interface{}{}
	currentTime := 0.0
	targetDuration := 15.0
	
	fmt.Printf("\n🔄 Creating person video loops for OVERLAY:\n")
	for currentTime < targetDuration {
		clipLength := avatarDuration
		if currentTime + clipLength > targetDuration {
			clipLength = targetDuration - currentTime
		}
		
		fmt.Printf("   • Clip at %.1fs for %.1fs\n", currentTime, clipLength)
		
		// CRITICAL: Person video must be VERY visible!
		// In Shotstack API: offset is from CENTER (0,0)
		// Positive X = right, Positive Y = down
		// For bottom-right: both positive values
		avatarClips = append(avatarClips, map[string]interface{}{
			"asset": map[string]interface{}{
				"type": "video",
				"src":  avatarVideoURL,
			},
			"start":  currentTime,
			"length": clipLength,
			"offset": map[string]interface{}{
				"x": 0.35,  // Right side (35% from center)
				"y": 0.35,  // Bottom side (35% from center)
			},
			"scale":   0.45,  // 45% size - BIGGER for clear visibility
			"opacity": 1.0,   // 100% opacity - FULLY visible, NOT transparent
			"fit":     "contain",  // Preserve aspect ratio
		})
		
		currentTime += clipLength
	}
	fmt.Printf("✅ Created %d person clips - PERSON WILL BE IN BOTTOM-RIGHT CORNER\n", len(avatarClips))
	fmt.Printf("📌 Position: Bottom-Right (offset x=0.35, y=0.35 from center)\n")
	fmt.Printf("📏 Size: 45%% of screen (scale=0.45) - LARGE SIZE FOR VISIBILITY\n")
	fmt.Printf("👁️  Opacity: 100%% (fully visible, NOT transparent)\n")
	fmt.Printf("🎬 Product video URL: %s\n", productVideoURL)
	fmt.Printf("🎥 Person video URL: %s\n", avatarVideoURL)
	fmt.Printf("🎯 REVERSED TRACK ORDER: Track 0 = Person (TOP), Track 1 = Product (BOTTOM)\n")
	fmt.Printf("💡 Testing reversed order - person video should be VISIBLE on top\n\n")
	
	// REVERSED ORDER: Try putting person video FIRST (Track 0) and product SECOND (Track 1)
	// Sometimes video APIs render in reverse order
	return []interface{}{
		// Track 0: Person video (TRY AS FIRST TRACK - may render on TOP)
		map[string]interface{}{
			"clips": avatarClips,
		},
		// Track 1: Product video (TRY AS SECOND TRACK - may render as BACKGROUND)
		map[string]interface{}{
			"clips": []interface{}{
				map[string]interface{}{
					"asset": map[string]interface{}{
						"type": "video",
						"src":  productVideoURL,
					},
					"start":  0.0,
					"length": 15.0,
					"fit":    "cover",  // Fill entire frame
				},
			},
		},
	}
}

