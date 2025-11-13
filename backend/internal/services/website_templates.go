package services

import (
	"fmt"
	"time"
)

// MarketingWebsiteTemplate generates professional marketing website HTML
func MarketingWebsiteTemplate(productName, productDescription, videoURL, productImageURL string) string {
	if productName == "" {
		productName = "Amazing Product"
	}
	if productDescription == "" {
		productDescription = "Discover the future of innovation with our cutting-edge product."
	}

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
            <div class="features-grid">
                <div class="feature-card">
                    <div class="feature-icon">🚀</div>
                    <h3>Lightning Fast</h3>
                    <p>Experience unparalleled speed and efficiency that transforms your workflow.</p>
                </div>
                <div class="feature-card">
                    <div class="feature-icon">💎</div>
                    <h3>Premium Quality</h3>
                    <p>Built with the finest materials and cutting-edge technology.</p>
                </div>
                <div class="feature-card">
                    <div class="feature-icon">🔒</div>
                    <h3>Secure & Reliable</h3>
                    <p>Your data and privacy are our top priorities.</p>
                </div>
                <div class="feature-card">
                    <div class="feature-icon">🎯</div>
                    <h3>Easy to Use</h3>
                    <p>Intuitive design that anyone can master in minutes.</p>
                </div>
            </div>
        </div>
    </section>

    <!-- Video Demo Section -->
    <section id="video" class="video-section">
        <div class="container">
            <h2 class="section-title">See It In Action</h2>
            <p class="section-subtitle">Watch our product demonstration and discover what makes it special</p>
            <div class="video-wrapper">
                %s
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
		productDescription,
		productName,
		productName,
		productName,
		productDescription,
		productImageURL,
		productName,
		productName,
		func() string {
			if videoURL != "" {
				return fmt.Sprintf(`<video controls class="promo-video" poster="%s">
                    <source src="%s" type="video/mp4">
                    Your browser does not support the video tag.
                </video>`, productImageURL, videoURL)
			}
			return `<div class="video-placeholder">
                    <p>Video coming soon...</p>
                </div>`
		}(),
		productName,
		time.Now().Year(),
		productName,
	)
}

// ModernWebsiteCSS generates modern marketing CSS
func ModernWebsiteCSS() string {
	return `* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

:root {
    --primary-color: #6366f1;
    --primary-dark: #4f46e5;
    --secondary-color: #8b5cf6;
    --text-dark: #1f2937;
    --text-light: #6b7280;
    --bg-light: #f9fafb;
    --white: #ffffff;
    --gradient: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

body {
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    line-height: 1.6;
    color: var(--text-dark);
    background: var(--white);
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
}

/* Navigation */
.navbar {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    padding: 1rem 0;
    position: sticky;
    top: 0;
    z-index: 1000;
    box-shadow: 0 2px 10px rgba(0,0,0,0.05);
}

.navbar .container {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.logo {
    font-size: 1.5rem;
    font-weight: 900;
    background: var(--gradient);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
}

.nav-menu {
    display: flex;
    list-style: none;
    gap: 2rem;
}

.nav-menu a {
    text-decoration: none;
    color: var(--text-dark);
    font-weight: 500;
    transition: color 0.3s;
}

.nav-menu a:hover {
    color: var(--primary-color);
}

/* Hero Section */
.hero {
    background: var(--gradient);
    padding: 4rem 0 6rem;
    color: var(--white);
}

.hero-content {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 4rem;
    align-items: center;
    padding: 4rem 20px;
}

.hero-title {
    font-size: 3.5rem;
    font-weight: 900;
    line-height: 1.1;
    margin-bottom: 1.5rem;
}

.hero-description {
    font-size: 1.25rem;
    margin-bottom: 2rem;
    opacity: 0.95;
}

.hero-cta {
    display: flex;
    gap: 1rem;
}

.btn {
    padding: 0.875rem 2rem;
    border-radius: 8px;
    text-decoration: none;
    font-weight: 600;
    transition: all 0.3s;
    display: inline-block;
    border: none;
    cursor: pointer;
    font-size: 1rem;
}

.btn-primary {
    background: var(--white);
    color: var(--primary-color);
}

.btn-primary:hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 25px rgba(0,0,0,0.2);
}

.btn-secondary {
    background: transparent;
    color: var(--white);
    border: 2px solid var(--white);
}

.btn-secondary:hover {
    background: var(--white);
    color: var(--primary-color);
}

.btn-large {
    padding: 1.25rem 3rem;
    font-size: 1.125rem;
}

.btn-outline {
    background: transparent;
    border: 2px solid var(--white);
    color: var(--white);
}

.product-showcase {
    width: 100%;
    height: auto;
    border-radius: 20px;
    box-shadow: 0 20px 60px rgba(0,0,0,0.3);
    animation: float 3s ease-in-out infinite;
}

@keyframes float {
    0%, 100% { transform: translateY(0px); }
    50% { transform: translateY(-20px); }
}

/* Features Section */
.features {
    padding: 6rem 0;
    background: var(--bg-light);
}

.section-title {
    font-size: 2.5rem;
    font-weight: 900;
    text-align: center;
    margin-bottom: 3rem;
    color: var(--text-dark);
}

.section-subtitle {
    text-align: center;
    font-size: 1.25rem;
    color: var(--text-light);
    margin-bottom: 3rem;
    max-width: 600px;
    margin-left: auto;
    margin-right: auto;
}

.features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 2rem;
}

.feature-card {
    background: var(--white);
    padding: 2rem;
    border-radius: 15px;
    text-align: center;
    transition: transform 0.3s, box-shadow 0.3s;
    box-shadow: 0 5px 15px rgba(0,0,0,0.08);
}

.feature-card:hover {
    transform: translateY(-10px);
    box-shadow: 0 15px 40px rgba(0,0,0,0.15);
}

.feature-icon {
    font-size: 3rem;
    margin-bottom: 1rem;
}

.feature-card h3 {
    font-size: 1.5rem;
    margin-bottom: 1rem;
    color: var(--text-dark);
}

.feature-card p {
    color: var(--text-light);
}

/* Video Section */
.video-section {
    padding: 6rem 0;
    background: var(--white);
}

.video-wrapper {
    max-width: 900px;
    margin: 0 auto;
    border-radius: 20px;
    overflow: hidden;
    box-shadow: 0 20px 60px rgba(0,0,0,0.15);
}

.promo-video {
    width: 100%;
    height: auto;
    display: block;
}

.video-placeholder {
    background: var(--bg-light);
    padding: 4rem;
    text-align: center;
    color: var(--text-light);
    font-size: 1.25rem;
}

/* CTA Section */
.cta-section {
    padding: 6rem 0;
    background: var(--gradient);
    color: var(--white);
    text-align: center;
}

.cta-content h2 {
    font-size: 2.5rem;
    font-weight: 900;
    margin-bottom: 1rem;
}

.cta-content p {
    font-size: 1.25rem;
    margin-bottom: 2rem;
    opacity: 0.95;
}

.cta-buttons {
    display: flex;
    gap: 1rem;
    justify-content: center;
    flex-wrap: wrap;
}

/* Footer */
.footer {
    background: var(--text-dark);
    color: var(--white);
    padding: 4rem 0 2rem;
}

.footer-content {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 2rem;
    margin-bottom: 2rem;
}

.footer-section h4 {
    margin-bottom: 1rem;
    font-size: 1.125rem;
}

.footer-section ul {
    list-style: none;
}

.footer-section a {
    color: rgba(255,255,255,0.8);
    text-decoration: none;
    display: block;
    margin-bottom: 0.5rem;
    transition: color 0.3s;
}

.footer-section a:hover {
    color: var(--white);
}

.social-links {
    display: flex;
    gap: 1rem;
}

.footer-bottom {
    text-align: center;
    padding-top: 2rem;
    border-top: 1px solid rgba(255,255,255,0.1);
    color: rgba(255,255,255,0.6);
}

/* Responsive */
@media (max-width: 768px) {
    .hero-content {
        grid-template-columns: 1fr;
        text-align: center;
    }
    
    .hero-title {
        font-size: 2.5rem;
    }
    
    .hero-cta {
        justify-content: center;
        flex-direction: column;
    }
    
    .features-grid {
        grid-template-columns: 1fr;
    }
    
    .nav-menu {
        display: none;
    }
    
    .cta-buttons {
        flex-direction: column;
        align-items: center;
    }
}
`
}

// InteractiveJS generates enhanced JavaScript for the website
func InteractiveJS() string {
	return `// Marketing Website Interactive Features
document.addEventListener('DOMContentLoaded', function() {
    console.log('AI Marketing Website Loaded');
    
    // Smooth scrolling for navigation links
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
    
    // Video play tracking
    const video = document.querySelector('.promo-video');
    if (video) {
        video.addEventListener('play', function() {
            console.log('Video started playing');
            // Track video engagement (integrate with analytics)
        });
        
        video.addEventListener('ended', function() {
            console.log('Video finished');
            // Track completion
        });
    }
    
    // Intersection Observer for animations
    const observerOptions = {
        threshold: 0.1,
        rootMargin: '0px 0px -100px 0px'
    };
    
    const observer = new IntersectionObserver(function(entries) {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.style.opacity = '1';
                entry.target.style.transform = 'translateY(0)';
            }
        });
    }, observerOptions);
    
    // Observe feature cards
    document.querySelectorAll('.feature-card').forEach(card => {
        card.style.opacity = '0';
        card.style.transform = 'translateY(20px)';
        card.style.transition = 'opacity 0.6s ease, transform 0.6s ease';
        observer.observe(card);
    });
    
    // CTA button click tracking
    document.querySelectorAll('.btn').forEach(button => {
        button.addEventListener('click', function() {
            const buttonText = this.textContent;
            console.log('Button clicked:', buttonText);
            // Track button clicks
        });
    });
    
    // Add parallax effect to hero
    window.addEventListener('scroll', function() {
        const scrolled = window.pageYOffset;
        const hero = document.querySelector('.hero');
        if (hero) {
            hero.style.transform = 'translateY(' + (scrolled * 0.5) + 'px)';
        }
    });
});
`
}

