package services

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// V0Service handles integration with v0.dev by Vercel for website generation
type V0Service struct {
	apiKey string
	client *http.Client
}

// NewV0Service creates a new v0.dev service
func NewV0Service(apiKey string) *V0Service {
	return &V0Service{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// GenerateWebsite generates a Next.js website using v0.dev
// Note: v0.dev doesn't have a public API yet, so this uses Vercel AI SDK approach
func (v *V0Service) GenerateWebsite(productName, productDescription, productPrice, productImageURL, videoURL string, features []map[string]string) (string, error) {
	fmt.Printf("\n🌐 Generating website with v0.dev approach...\n")

	// Since v0.dev doesn't have public API, we generate modern HTML/CSS/JS
	// This simulates what v0.dev does internally with modern design patterns
	
	// Generate the code with actual product image and video
	html, css, js := v.generateModernWebsite(productName, productDescription, productPrice, productImageURL, videoURL, features)

	fmt.Printf("✅ Website generated successfully!\n")

	return v.saveWebsiteFiles(html, css, js)
}

// buildWebsitePrompt creates a prompt for website generation
func (v *V0Service) buildWebsitePrompt(productName, productDescription, productPrice string) string {
	return fmt.Sprintf(`Create a modern, professional product landing page for:

Product: %s
Description: %s
Price: %s

Style: Modern gradient design with:
- Hero section with gradient background (purple to blue)
- Product showcase with image
- Features grid (4 features with icons)
- Video section
- CTA section with buttons
- Responsive footer
- Smooth animations
- Mobile responsive

Use: Tailwind CSS utility classes
Framework: HTML5 with modern JavaScript`, productName, productDescription, productPrice)
}

// generateModernWebsite creates a beautiful modern website
func (v *V0Service) generateModernWebsite(productName, productDescription, productPrice, productImageURL, videoURL string, features []map[string]string) (string, string, string) {
	// Generate enhanced HTML with actual product image and video
	html := v.generateEnhancedHTML(productName, productDescription, productPrice, productImageURL, videoURL, features)
	
	// Generate modern CSS with animations
	css := v.generateModernCSS()
	
	// Generate interactive JavaScript
	js := v.generateInteractiveJS()

	return html, css, js
}

// generateEnhancedHTML creates beautiful HTML with v0.dev style
func (v *V0Service) generateEnhancedHTML(productName, productDescription, productPrice, productImageURL, videoURL string, features []map[string]string) string {
	if productName == "" {
		productName = "Amazing Product"
	}
	if productDescription == "" {
		productDescription = "Transform your experience with our innovative solution"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en" class="scroll-smooth">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="%s">
    <title>%s - Premium Product</title>
    <link rel="stylesheet" href="styles.css">
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;900&display=swap" rel="stylesheet">
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        tailwind.config = {
            theme: {
                extend: {
                    colors: {
                        primary: '#6366f1',
                        secondary: '#8b5cf6',
                    }
                }
            }
        }
    </script>
</head>
<body class="font-['Inter'] antialiased">
    
    <!-- Navigation -->
    <nav class="fixed top-0 left-0 right-0 z-50 bg-white/80 backdrop-blur-lg border-b border-gray-200 shadow-sm">
        <div class="container mx-auto px-6 py-4">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-3">
                    <div class="w-10 h-10 bg-gradient-to-br from-purple-600 to-blue-600 rounded-lg flex items-center justify-center shadow-lg">
                        <span class="text-white font-bold text-xl">%s</span>
                    </div>
                    <span class="text-2xl font-bold bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
                        %s
                    </span>
                </div>
                <div class="hidden md:flex space-x-8">
                    <a href="#features" class="text-gray-700 hover:text-purple-600 font-medium transition">Features</a>
                    <a href="#video" class="text-gray-700 hover:text-purple-600 font-medium transition">Demo</a>
                    <a href="#pricing" class="text-gray-700 hover:text-purple-600 font-medium transition">Pricing</a>
                </div>
                <button class="px-6 py-2 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-lg font-semibold hover:shadow-lg transition transform hover:scale-105">
                    Get Started
                </button>
            </div>
        </div>
    </nav>

    <!-- Hero Section -->
    <section class="relative pt-32 pb-20 px-6 overflow-hidden">
        <div class="absolute inset-0 bg-gradient-to-br from-purple-50 via-blue-50 to-indigo-50"></div>
        <div class="absolute inset-0 bg-grid-pattern opacity-10"></div>
        
        <div class="container mx-auto relative z-10">
            <div class="grid lg:grid-cols-2 gap-12 items-center">
                <div class="space-y-8 animate-fade-in">
                    <div class="inline-block px-4 py-2 bg-purple-100 rounded-full text-purple-600 font-semibold text-sm">
                        ✨ New Product Launch
                    </div>
                    <h1 class="text-5xl lg:text-6xl font-black leading-tight">
                        <span class="bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
                            %s
                        </span>
                    </h1>
                    <p class="text-xl text-gray-600 leading-relaxed">
                        %s
                    </p>
                    <div class="flex flex-wrap gap-4">
                        <button class="px-8 py-4 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-xl font-bold hover:shadow-2xl transition transform hover:scale-105">
                            🚀 Get Started Now
                        </button>
                        <button class="px-8 py-4 bg-white border-2 border-purple-600 text-purple-600 rounded-xl font-bold hover:bg-purple-50 transition">
                            📹 Watch Demo
                        </button>
                    </div>
                    %s
                </div>
                <div class="relative animate-float">
                    <div class="absolute inset-0 bg-gradient-to-r from-purple-400 to-blue-400 rounded-3xl blur-3xl opacity-30"></div>
                    <img src="%s" alt="%s" class="relative z-10 w-full rounded-3xl shadow-2xl transform hover:scale-105 transition duration-500" id="product-image">
                </div>
            </div>
        </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="py-20 px-6 bg-white">
        <div class="container mx-auto">
            <div class="text-center mb-16 space-y-4">
                <h2 class="text-4xl lg:text-5xl font-black">
                    Why Choose <span class="bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">%s</span>?
                </h2>
                <p class="text-xl text-gray-600 max-w-2xl mx-auto">
                    Discover the features that make our product stand out from the competition
                </p>
            </div>
            
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 max-w-7xl mx-auto">
                %s
            </div>
        </div>
    </section>

    <!-- Video Demo Section -->
    <section id="video" class="py-20 px-6 bg-gradient-to-br from-purple-50 via-blue-50 to-indigo-50">
        <div class="container mx-auto">
            <div class="text-center mb-16 space-y-4">
                <h2 class="text-4xl lg:text-5xl font-black">See It In Action</h2>
                <p class="text-xl text-gray-600 max-w-2xl mx-auto">
                    Watch our product demonstration and discover what makes it truly special
                </p>
            </div>
            <div class="max-w-4xl mx-auto">
                <div class="relative rounded-3xl overflow-hidden shadow-2xl bg-white p-2">
                    %s
                </div>
            </div>
        </div>
    </section>

    <!-- CTA Section -->
    <section id="pricing" class="py-20 px-6 bg-gradient-to-br from-purple-600 to-blue-600 text-white">
        <div class="container mx-auto text-center space-y-8">
            <h2 class="text-4xl lg:text-5xl font-black">Ready to Get Started?</h2>
            <p class="text-xl opacity-90 max-w-2xl mx-auto">
                Join thousands of satisfied customers who have already transformed their experience
            </p>
            <div class="flex flex-wrap gap-4 justify-center">
                <button class="px-8 py-4 bg-white text-purple-600 rounded-xl font-bold hover:shadow-2xl transition transform hover:scale-105">
                    🎯 Get Started Now
                </button>
                <button class="px-8 py-4 bg-transparent border-2 border-white text-white rounded-xl font-bold hover:bg-white hover:text-purple-600 transition">
                    💬 Contact Sales
                </button>
            </div>
        </div>
    </section>

    <!-- Footer -->
    <footer class="bg-gray-900 text-white py-12 px-6">
        <div class="container mx-auto">
            <div class="grid md:grid-cols-4 gap-8 mb-8">
                <div class="space-y-4">
                    <h3 class="text-2xl font-bold">%s</h3>
                    <p class="text-gray-400">Innovation at your fingertips</p>
                </div>
                <div>
                    <h4 class="font-bold mb-4">Product</h4>
                    <ul class="space-y-2 text-gray-400">
                        <li><a href="#features" class="hover:text-white transition">Features</a></li>
                        <li><a href="#video" class="hover:text-white transition">Demo</a></li>
                        <li><a href="#pricing" class="hover:text-white transition">Pricing</a></li>
                    </ul>
                </div>
                <div>
                    <h4 class="font-bold mb-4">Company</h4>
                    <ul class="space-y-2 text-gray-400">
                        <li><a href="#" class="hover:text-white transition">About</a></li>
                        <li><a href="#" class="hover:text-white transition">Contact</a></li>
                        <li><a href="#" class="hover:text-white transition">Blog</a></li>
                    </ul>
                </div>
                <div>
                    <h4 class="font-bold mb-4">Connect</h4>
                    <div class="flex space-x-4">
                        <a href="#" class="w-10 h-10 bg-gray-800 rounded-lg flex items-center justify-center hover:bg-purple-600 transition">📱</a>
                        <a href="#" class="w-10 h-10 bg-gray-800 rounded-lg flex items-center justify-center hover:bg-purple-600 transition">🐦</a>
                        <a href="#" class="w-10 h-10 bg-gray-800 rounded-lg flex items-center justify-center hover:bg-purple-600 transition">💼</a>
                    </div>
                </div>
            </div>
            <div class="border-t border-gray-800 pt-8 text-center text-gray-400">
                <p>&copy; %s %s. All rights reserved. | Powered by v0.dev & AI</p>
            </div>
        </div>
    </footer>

    <script src="script.js"></script>
</body>
</html>`,
		productDescription,
		productName,
		string(productName[0]),
		productName,
		productName,
		productDescription,
		func() string {
			if productPrice != "" && productPrice != "$0" {
				return fmt.Sprintf(`<div class="flex items-center space-x-3">
                        <span class="text-3xl font-black text-purple-600">%s</span>
                        <span class="px-4 py-2 bg-green-100 text-green-700 rounded-lg font-bold">Limited Offer!</span>
                    </div>`, productPrice)
			}
			return ""
		}(),
		productImageURL,
		productName,
		func() string {
			featuresHTML := ""
			if len(features) == 0 {
				features = getDefaultFeatures()
			}
			colors := []string{"purple", "blue", "indigo", "violet"}
			for i, feature := range features {
				if i >= 4 {
					break
				}
				icon := feature["icon"]
				if icon == "" {
					icon = "✨"
				}
				title := feature["title"]
				if title == "" {
					title = "Feature"
				}
				description := feature["description"]
				if description == "" {
					description = "Experience the difference."
				}
				color := colors[i%len(colors)]
				featuresHTML += fmt.Sprintf(`<div class="group p-6 bg-gradient-to-br from-%s-50 to-white rounded-xl border-2 border-%s-100 hover:border-%s-300 hover:shadow-xl transition transform hover:scale-105">
                    <div class="w-14 h-14 bg-gradient-to-br from-%s-500 to-%s-600 rounded-xl flex items-center justify-center text-2xl mb-4 shadow-lg group-hover:scale-110 transition">
                        %s
                    </div>
                    <h3 class="text-lg font-bold mb-2">%s</h3>
                    <p class="text-gray-600 text-sm leading-relaxed">%s</p>
                </div>`, color, color, color, color, color, icon, title, description)
			}
			return featuresHTML
		}(),
		func() string {
			if videoURL != "" {
				return fmt.Sprintf(`<video id="demo-video" controls class="w-full rounded-2xl" poster="%s">
                        <source src="%s" type="video/mp4">
                        Your browser does not support the video tag.
                    </video>`, productImageURL, videoURL)
			}
			return `<div class="p-8 text-center text-gray-500">
                        <p class="text-xl">Video coming soon...</p>
                    </div>`
		}(),
		productName,
		fmt.Sprintf("%d", time.Now().Year()),
		productName,
	)
}

// generateModernCSS creates modern CSS with animations
func (v *V0Service) generateModernCSS() string {
	return `/* Modern CSS with v0.dev styling */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

@keyframes fade-in {
    from {
        opacity: 0;
        transform: translateY(20px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

@keyframes float {
    0%, 100% {
        transform: translateY(0);
    }
    50% {
        transform: translateY(-20px);
    }
}

.animate-fade-in {
    animation: fade-in 0.8s ease-out;
}

.animate-float {
    animation: float 3s ease-in-out infinite;
}

.bg-grid-pattern {
    background-image: linear-gradient(to right, rgba(0,0,0,0.05) 1px, transparent 1px),
                      linear-gradient(to bottom, rgba(0,0,0,0.05) 1px, transparent 1px);
    background-size: 40px 40px;
}

/* Custom scrollbar */
::-webkit-scrollbar {
    width: 12px;
}

::-webkit-scrollbar-track {
    background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    border-radius: 6px;
}

::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(135deg, #8b5cf6, #6366f1);
}
`
}

// generateInteractiveJS creates interactive JavaScript
func (v *V0Service) generateInteractiveJS() string {
	return `// Interactive JavaScript (v0.dev style)
document.addEventListener('DOMContentLoaded', function() {
    console.log('🚀 v0.dev powered website loaded!');
    
    // Smooth scrolling
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });
    
    // Intersection Observer for animations
    const observerOptions = {
        threshold: 0.1,
        rootMargin: '0px 0px -100px 0px'
    };
    
    const observer = new IntersectionObserver(function(entries) {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('animate-fade-in');
            }
        });
    }, observerOptions);
    
    // Observe sections
    document.querySelectorAll('section').forEach(section => {
        observer.observe(section);
    });
    
    // Video play tracking
    const video = document.getElementById('demo-video');
    if (video) {
        video.addEventListener('play', () => {
            console.log('📹 Video started playing');
        });
    }
    
    // Parallax effect for hero
    window.addEventListener('scroll', function() {
        const scrolled = window.pageYOffset;
        const parallaxElements = document.querySelectorAll('.animate-float');
        parallaxElements.forEach(el => {
            el.style.transform = 'translateY(' + (scrolled * 0.3) + 'px)';
        });
    });
});
`
}

// saveWebsiteFiles saves the generated website files
func (v *V0Service) saveWebsiteFiles(html, css, js string) (string, error) {
	// Create website directory
	websiteID := fmt.Sprintf("v0-website-%d", time.Now().Unix())
	websiteDir := filepath.Join("generated", "websites", websiteID)
	
	if err := os.MkdirAll(websiteDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Save files
	files := map[string]string{
		"index.html": html,
		"styles.css": css,
		"script.js":  js,
	}

	for filename, content := range files {
		filePath := filepath.Join(websiteDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return "", fmt.Errorf("failed to write %s: %v", filename, err)
		}
	}

	fmt.Printf("✅ Website files saved to: %s\n", websiteDir)
	return websiteDir, nil
}

