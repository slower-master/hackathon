package services

import (
	"fmt"
	"time"
)

// getDefaultFeatures returns default features if Gemini fails
// Updated to be more appropriate for namkeen/snacks
func getDefaultFeatures() []map[string]string {
	return []map[string]string{
		{"icon": "🌾", "title": "100% Natural Ingredients", "description": "Made from premium quality ingredients with no artificial flavors or preservatives for pure, authentic taste."},
		{"icon": "😋", "title": "Irresistibly Delicious", "description": "Perfectly seasoned with traditional spices that create an unforgettable burst of flavor in every bite."},
		{"icon": "✨", "title": "Fresh & Crunchy", "description": "Carefully crafted to maintain optimal crunchiness and freshness in every pack you open."},
		{"icon": "🎯", "title": "Perfect for Every Occasion", "description": "Ideal for tea-time snacking, parties, or whenever you crave something tasty and satisfying."},
	}
}

// MarketingWebsiteTemplate generates professional marketing website HTML
func MarketingWebsiteTemplate(productName, productDescription, videoURL, productImageURL string, features []map[string]string) string {
	if productName == "" {
		productName = "Amazing Product"
	}
	if productDescription == "" {
		productDescription = "Discover the future of innovation with our cutting-edge product."
	}
	if len(features) == 0 {
		features = getDefaultFeatures()
	}

	// Generate features HTML
	featuresHTML := ""
	for _, feature := range features {
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
		featuresHTML += fmt.Sprintf(`
                <div class="feature-card">
                    <div class="feature-icon">%s</div>
                    <h3>%s</h3>
                    <p>%s</p>
                </div>`, icon, title, description)
	}

	// Generate video HTML
	videoHTML := ""
	if videoURL != "" {
		videoHTML = fmt.Sprintf(`
                <video controls class="promo-video" poster="%s">
                    <source src="%s" type="video/mp4">
                    Your browser does not support the video tag.
                </video>`, productImageURL, videoURL)
	} else {
		videoHTML = `
                <div class="video-placeholder">
                    <p>🎬 Video coming soon...</p>
                </div>`
	}

	currentYear := time.Now().Year()

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="%s">
    <title>%s - Official Product Page</title>
    <link rel="stylesheet" href="styles.css">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;600;700;900&display=swap" rel="stylesheet">
</head>
<body>
    <!-- Hero Section -->
    <section class="hero">
        <nav class="navbar">
            <div class="container">
                <div class="logo">%s</div>
                <ul class="nav-menu">
                    <li><a href="#features">Features</a></li>
                    <li><a href="#video">Watch Demo</a></li>
                    <li><a href="#cta">Get Started</a></li>
                </ul>
            </div>
        </nav>
        
        <div class="hero-content container">
            <div class="hero-text">
                <h1 class="hero-title">%s</h1>
                <p class="hero-description">%s</p>
                <div class="hero-cta">
                    <a href="#video" class="btn btn-primary">Watch Demo</a>
                    <a href="#cta" class="btn btn-secondary">Learn More</a>
                </div>
            </div>
            <div class="hero-image">
                <img src="%s" alt="%s" class="product-showcase">
            </div>
        </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="features">
        <div class="container">
            <h2 class="section-title">Why Choose %s?</h2>
            <p class="section-subtitle">Discover what makes our product exceptional</p>
            <div class="features-grid">%s
            </div>
        </div>
    </section>

    <!-- Video Demo Section -->
    <section id="video" class="video-section">
        <div class="container">
            <h2 class="section-title">See It In Action</h2>
            <p class="section-subtitle">Watch our product demonstration and discover what makes it special</p>
            <div class="video-wrapper">%s
            </div>
        </div>
    </section>

    <!-- Call to Action -->
    <section id="cta" class="cta-section">
        <div class="container">
            <div class="cta-content">
                <h2>Ready to Transform Your Experience?</h2>
                <p>Join thousands of satisfied customers who have already made the switch.</p>
                <div class="cta-buttons">
                    <button class="btn btn-large btn-primary">Get Started Now</button>
                    <button class="btn btn-large btn-outline">Contact Sales</button>
                </div>
            </div>
        </div>
    </section>

    <!-- Footer -->
    <footer class="footer">
        <div class="container">
            <div class="footer-content">
                <div class="footer-section">
                    <h4>%s</h4>
                    <p>Innovation at your fingertips</p>
                </div>
                <div class="footer-section">
                    <h4>Product</h4>
                    <ul>
                        <li><a href="#features">Features</a></li>
                        <li><a href="#video">Demo</a></li>
                        <li><a href="#">Pricing</a></li>
                    </ul>
                </div>
                <div class="footer-section">
                    <h4>Company</h4>
                    <ul>
                        <li><a href="#">About</a></li>
                        <li><a href="#">Contact</a></li>
                        <li><a href="#">Blog</a></li>
                    </ul>
                </div>
                <div class="footer-section">
                    <h4>Connect</h4>
                    <div class="social-links">
                        <a href="#">Twitter</a>
                        <a href="#">LinkedIn</a>
                        <a href="#">Instagram</a>
                    </div>
                </div>
            </div>
            <div class="footer-bottom">
                <p>&copy; %d %s. All rights reserved. | Powered by AI Marketing Agent</p>
            </div>
        </div>
    </footer>

    <script src="script.js"></script>
</body>
</html>`,
		productDescription, // meta description
		productName,        // page title
		productName,        // logo
		productName,        // hero title
		productDescription, // hero description
		productImageURL,    // product image src
		productName,        // product image alt
		productName,        // "Why Choose X?"
		featuresHTML,       // features cards
		videoHTML,          // video player
		productName,        // footer brand
		currentYear,        // year
		productName,        // footer copyright
	)
}

// ModernWebsiteCSS generates modern marketing CSS
func ModernWebsiteCSS() string {
	return `/* ===== RESET & BASE ===== */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

:root {
    --primary: #667eea;
    --primary-dark: #5a67d8;
    --secondary: #764ba2;
    --text-dark: #1a202c;
    --text-light: #718096;
    --bg-light: #f7fafc;
    --white: #ffffff;
    --shadow-sm: 0 2px 8px rgba(0,0,0,0.08);
    --shadow-md: 0 4px 20px rgba(0,0,0,0.12);
    --shadow-lg: 0 20px 40px rgba(0,0,0,0.15);
}

body {
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    line-height: 1.6;
    color: var(--text-dark);
    background: var(--white);
    overflow-x: hidden;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
}

/* ===== NAVIGATION ===== */
.navbar {
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(10px);
    box-shadow: var(--shadow-sm);
    position: sticky;
    top: 0;
    z-index: 1000;
    padding: 1rem 0;
}

.navbar .container {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.logo {
    font-size: 1.5rem;
    font-weight: 900;
    background: linear-gradient(135deg, var(--primary) 0%, var(--secondary) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
}

.nav-menu {
    list-style: none;
    display: flex;
    gap: 2rem;
}

.nav-menu a {
    text-decoration: none;
    color: var(--text-dark);
    font-weight: 600;
    transition: color 0.3s;
    font-size: 0.95rem;
}

.nav-menu a:hover {
    color: var(--primary);
}

/* ===== HERO SECTION ===== */
.hero {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 4rem 0 6rem;
    position: relative;
    overflow: hidden;
}

.hero::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: url('data:image/svg+xml,<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg"><circle cx="50" cy="50" r="2" fill="white" opacity="0.1"/></svg>');
    opacity: 0.3;
}

.hero-content {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 4rem;
    align-items: center;
    padding: 4rem 20px;
    position: relative;
    z-index: 1;
}

.hero-text {
    color: var(--white);
}

.hero-title {
    font-size: 3.5rem;
    font-weight: 900;
    line-height: 1.1;
    margin-bottom: 1.5rem;
    text-shadow: 0 4px 20px rgba(0,0,0,0.2);
}

.hero-description {
    font-size: 1.25rem;
    line-height: 1.8;
    margin-bottom: 2rem;
    opacity: 0.95;
}

.hero-cta {
    display: flex;
    gap: 1rem;
    flex-wrap: wrap;
}

.hero-image {
    display: flex;
    justify-content: center;
    align-items: center;
}

.product-showcase {
    width: 100%;
    max-width: 500px;
    height: auto;
    border-radius: 20px;
    box-shadow: 0 30px 60px rgba(0,0,0,0.3);
    animation: float 6s ease-in-out infinite;
    transform-origin: center;
}

@keyframes float {
    0%, 100% { transform: translateY(0) rotate(0deg); }
    25% { transform: translateY(-20px) rotate(1deg); }
    50% { transform: translateY(-10px) rotate(-1deg); }
    75% { transform: translateY(-15px) rotate(0.5deg); }
}

/* ===== BUTTONS ===== */
.btn {
    display: inline-block;
    padding: 0.875rem 2rem;
    font-size: 1rem;
    font-weight: 700;
    text-decoration: none;
    border-radius: 50px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    cursor: pointer;
    border: none;
    text-align: center;
}

.btn-primary {
    background: var(--white);
    color: var(--primary);
    box-shadow: 0 4px 15px rgba(0,0,0,0.2);
}

.btn-primary:hover {
    transform: translateY(-3px);
    box-shadow: 0 10px 30px rgba(0,0,0,0.3);
}

.btn-secondary {
    background: transparent;
    border: 2px solid var(--white);
    color: var(--white);
}

.btn-secondary:hover {
    background: var(--white);
    color: var(--primary);
}

.btn-large {
    padding: 1.125rem 2.5rem;
    font-size: 1.125rem;
}

.btn-outline {
    background: transparent;
    border: 2px solid var(--primary);
    color: var(--primary);
}

.btn-outline:hover {
    background: var(--primary);
    color: var(--white);
}

/* ===== FEATURES SECTION ===== */
.features {
    padding: 6rem 0;
    background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%);
}

.section-title {
    font-size: 2.75rem;
    font-weight: 900;
    text-align: center;
    margin-bottom: 1rem;
    color: var(--text-dark);
    letter-spacing: -0.5px;
}

.section-subtitle {
    text-align: center;
    font-size: 1.15rem;
    color: var(--text-light);
    max-width: 650px;
    margin: 0 auto 3rem;
    line-height: 1.6;
}

.features-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 2rem;
    max-width: 1200px;
    margin: 0 auto;
}

.feature-card {
    background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
    padding: 2.5rem 2rem;
    border-radius: 16px;
    text-align: center;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 4px 20px rgba(0,0,0,0.08);
    border: 2px solid transparent;
    position: relative;
    overflow: hidden;
}

.feature-card::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
    transform: scaleX(0);
    transition: transform 0.4s ease;
}

.feature-card:hover::before {
    transform: scaleX(1);
}

.feature-card:hover {
    transform: translateY(-12px);
    box-shadow: 0 20px 40px rgba(0,0,0,0.15);
    border-color: rgba(102, 126, 234, 0.3);
}

.feature-icon {
    font-size: 3.5rem;
    margin-bottom: 1.25rem;
    display: inline-block;
    transition: transform 0.4s ease;
}

.feature-card:hover .feature-icon {
    transform: scale(1.15) rotate(5deg);
}

.feature-card h3 {
    font-size: 1.35rem;
    font-weight: 700;
    margin-bottom: 0.75rem;
    color: var(--text-dark);
    line-height: 1.3;
}

.feature-card p {
    color: var(--text-light);
    font-size: 1rem;
    line-height: 1.7;
}

/* ===== VIDEO SECTION ===== */
.video-section {
    padding: 6rem 0;
    background: var(--white);
}

.video-wrapper {
    max-width: 800px;
    margin: 0 auto;
    border-radius: 24px;
    overflow: hidden;
    box-shadow: 0 25px 60px rgba(0,0,0,0.2);
    border: 4px solid #f8f9fa;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.video-wrapper:hover {
    transform: translateY(-5px);
    box-shadow: 0 30px 70px rgba(0,0,0,0.25);
}

.promo-video {
    width: 100%;
    height: auto;
    display: block;
    background: #000;
}

.video-placeholder {
    background: linear-gradient(135deg, #f7fafc 0%, #e2e8f0 100%);
    padding: 6rem 2rem;
    text-align: center;
    color: var(--text-light);
    font-size: 1.25rem;
}

/* ===== CTA SECTION ===== */
.cta-section {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 6rem 0;
    text-align: center;
    color: var(--white);
}

.cta-content h2 {
    font-size: 3rem;
    font-weight: 900;
    margin-bottom: 1.5rem;
}

.cta-content p {
    font-size: 1.25rem;
    margin-bottom: 2.5rem;
    opacity: 0.95;
}

.cta-buttons {
    display: flex;
    gap: 1rem;
    justify-content: center;
    flex-wrap: wrap;
}

/* ===== FOOTER ===== */
.footer {
    background: var(--text-dark);
    color: var(--white);
    padding: 4rem 0 2rem;
}

.footer-content {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 3rem;
    margin-bottom: 3rem;
}

.footer-section h4 {
    font-size: 1.25rem;
    font-weight: 700;
    margin-bottom: 1rem;
}

.footer-section ul {
    list-style: none;
}

.footer-section ul li {
    margin-bottom: 0.75rem;
}

.footer-section a {
    color: rgba(255,255,255,0.7);
    text-decoration: none;
    transition: color 0.3s;
}

.footer-section a:hover {
    color: var(--white);
}

.social-links {
    display: flex;
    gap: 1rem;
    flex-direction: column;
}

.footer-bottom {
    text-align: center;
    padding-top: 2rem;
    border-top: 1px solid rgba(255,255,255,0.1);
    color: rgba(255,255,255,0.6);
    font-size: 0.9rem;
}

/* ===== RESPONSIVE ===== */
@media (max-width: 1024px) {
    .features-grid {
        grid-template-columns: repeat(2, 1fr);
    }
    
    .hero-content {
        grid-template-columns: 1fr;
        text-align: center;
    }
    
    .hero-cta {
        justify-content: center;
    }
}

@media (max-width: 768px) {
    .hero-title {
        font-size: 2.5rem;
    }
    
    .section-title {
        font-size: 2rem;
    }
    
    .footer-content {
        grid-template-columns: repeat(2, 1fr);
    }
}

@media (max-width: 640px) {
    .features-grid {
        grid-template-columns: 1fr;
    }
    
    .footer-content {
        grid-template-columns: 1fr;
    }
    
    .nav-menu {
        gap: 1rem;
    }
    
    .cta-content h2 {
        font-size: 2rem;
    }
}`
}

// ModernWebsiteJS generates interactive JavaScript
func ModernWebsiteJS() string {
	return `// Smooth scrolling for navigation links
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

// Add scroll animation for feature cards
const observerOptions = {
    threshold: 0.1,
    rootMargin: '0px 0px -50px 0px'
};

const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.style.opacity = '1';
            entry.target.style.transform = 'translateY(0)';
        }
    });
}, observerOptions);

// Observe all feature cards
document.querySelectorAll('.feature-card').forEach((card, index) => {
    card.style.opacity = '0';
    card.style.transform = 'translateY(30px)';
    card.style.transition = 'opacity 0.6s ease, transform 0.6s ease';
    card.style.transitionDelay = (index * 0.1) + 's';
    observer.observe(card);
});

// Video autoplay on scroll
const video = document.querySelector('.promo-video');
if (video) {
    const videoObserver = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                video.play().catch(e => console.log('Autoplay prevented'));
            } else {
                video.pause();
            }
        });
    }, { threshold: 0.5 });
    
    videoObserver.observe(video);
}

// Navbar scroll effect
let lastScroll = 0;
const navbar = document.querySelector('.navbar');

window.addEventListener('scroll', () => {
    const currentScroll = window.pageYOffset;
    
    if (currentScroll > 100) {
        navbar.style.padding = '0.5rem 0';
        navbar.style.boxShadow = '0 4px 20px rgba(0,0,0,0.1)';
    } else {
        navbar.style.padding = '1rem 0';
        navbar.style.boxShadow = '0 2px 8px rgba(0,0,0,0.08)';
    }
    
    lastScroll = currentScroll;
});

console.log('🚀 Website loaded successfully!');`
}
