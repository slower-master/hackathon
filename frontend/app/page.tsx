'use client'

import { useState, useCallback } from 'react'
import { useDropzone } from 'react-dropzone'
import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

interface Project {
  id: string
  status: string
  product_name: string
  product_description: string
  product_category: string
  product_price: string
  generated_script: string
  product_image_path: string
  person_media_path: string
  generated_video_path?: string
  website_path?: string
  website_url?: string
  instagram_post_url?: string
}

export default function Home() {
  const [productImage, setProductImage] = useState<File | null>(null)
  const [personMedia, setPersonMedia] = useState<File | null>(null)
  const [productName, setProductName] = useState<string>('')
  const [productDescription, setProductDescription] = useState<string>('')
  const [productCategory, setProductCategory] = useState<string>('electronics')
  const [productPrice, setProductPrice] = useState<string>('')
  const [productVideoStyle, setProductVideoStyle] = useState<string>('cinematic')
  const [layout, setLayout] = useState<string>('product_main')
  const [project, setProject] = useState<Project | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [activeStep, setActiveStep] = useState(1)

  // Image preview URLs
  const [productImagePreview, setProductImagePreview] = useState<string | null>(null)
  const [personMediaPreview, setPersonMediaPreview] = useState<string | null>(null)

  const onProductImageDrop = useCallback((acceptedFiles: File[]) => {
    if (acceptedFiles.length > 0) {
      const file = acceptedFiles[0]
      setProductImage(file)
      setProductImagePreview(URL.createObjectURL(file))
    }
  }, [])

  const onPersonMediaDrop = useCallback((acceptedFiles: File[]) => {
    if (acceptedFiles.length > 0) {
      const file = acceptedFiles[0]
      setPersonMedia(file)
      setPersonMediaPreview(URL.createObjectURL(file))
    }
  }, [])

  const {
    getRootProps: getProductRootProps,
    getInputProps: getProductInputProps,
    isDragActive: isProductDragActive,
  } = useDropzone({
    onDrop: onProductImageDrop,
    accept: {
      'image/*': ['.png', '.jpg', '.jpeg', '.gif', '.webp'],
    },
    maxFiles: 1,
  })

  const {
    getRootProps: getPersonRootProps,
    getInputProps: getPersonInputProps,
    isDragActive: isPersonDragActive,
  } = useDropzone({
    onDrop: onPersonMediaDrop,
    accept: {
      'image/*': ['.png', '.jpg', '.jpeg', '.gif', '.webp'],
    },
    maxFiles: 1,
  })

  const handleUpload = async () => {
    if (!productImage || !personMedia) {
      setError('Please upload both product image and person photo')
      return
    }

    if (!productDescription.trim()) {
      setError('Please provide a product description')
      return
    }

    setLoading(true)
    setError(null)

    const formData = new FormData()
    formData.append('product_image', productImage)
    formData.append('person_media', personMedia)
    formData.append('product_name', productName)
    formData.append('product_description', productDescription)
    formData.append('product_category', productCategory)
    formData.append('product_price', productPrice)

    try {
      const response = await axios.post(`${API_URL}/api/v1/upload`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })

      setProject(response.data)
      setActiveStep(2)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to upload files')
    } finally {
      setLoading(false)
    }
  }

  const handleGenerateVideo = async () => {
    if (!project) return

    setLoading(true)
    setError(null)

    try {
      const response = await axios.post(
        `${API_URL}/api/v1/projects/${project.project_id || project.id}/generate-video`,
        {
          product_video_style: productVideoStyle,
          layout: layout
        }
      )
      setProject({ ...project, ...response.data })
      setActiveStep(3)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to generate video')
    } finally {
      setLoading(false)
    }
  }

  const handleGenerateWebsite = async () => {
    if (!project) return

    setLoading(true)
    setError(null)

    try {
      const response = await axios.post(
        `${API_URL}/api/v1/projects/${project.project_id || project.id}/generate-website`
      )
      
      console.log('✅ Website generated:', response.data)
      
      // Update project with new data
      const updatedProject = { 
        ...project, 
        website_path: response.data.website_path,
        status: response.data.status 
      }
      
      console.log('📦 Updated project:', updatedProject)
      setProject(updatedProject)
      setActiveStep(4)
    } catch (err: any) {
      console.error('❌ Website generation error:', err)
      setError(err.response?.data?.error || 'Failed to generate website')
    } finally {
      setLoading(false)
    }
  }

  const handleUploadToInstagram = async () => {
    if (!project) return

    setLoading(true)
    setError(null)

    try {
      const response = await axios.post(
        `${API_URL}/api/v1/projects/${project.project_id || project.id}/upload-to-instagram`,
        {
          // Instagram credentials can be set in environment variables
          // or provided here
        }
      )
      setProject({ ...project, ...response.data })
      alert('Video posted to Instagram successfully! 🎉')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to upload to Instagram')
    } finally {
      setLoading(false)
    }
  }

  const resetForm = () => {
    setProductImage(null)
    setPersonMedia(null)
    setProductName('')
    setProductDescription('')
    setProductCategory('electronics')
    setProductPrice('')
    setProductImagePreview(null)
    setPersonMediaPreview(null)
    setProject(null)
    setActiveStep(1)
    setError(null)
  }

  return (
    <main className="min-h-screen bg-gradient-to-br from-purple-50 via-white to-blue-50">
      {/* Header */}
      <header className="bg-white/80 backdrop-blur-lg border-b border-gray-200 sticky top-0 z-50 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="w-12 h-12 bg-gradient-to-br from-purple-600 to-blue-600 rounded-xl flex items-center justify-center shadow-lg">
                <span className="text-2xl">🚀</span>
              </div>
              <div>
                <h1 className="text-3xl font-extrabold bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
                  AI Product Marketing Studio
        </h1>
                <p className="text-sm text-gray-600">Create stunning videos & websites in minutes</p>
              </div>
            </div>
            {project && (
              <button
                onClick={resetForm}
                className="px-6 py-2 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-lg font-semibold hover:from-purple-700 hover:to-blue-700 transition-all shadow-md hover:shadow-lg"
              >
                New Project
              </button>
            )}
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Progress Steps */}
        <div className="mb-12">
          <div className="flex items-center justify-center space-x-4">
            {[
              { step: 1, label: 'Upload & Details', icon: '📤' },
              { step: 2, label: 'Generate Video', icon: '🎬' },
              { step: 3, label: 'Generate & Share', icon: '🚀' },
              { step: 4, label: 'View Website', icon: '🌐' },
            ].map((item) => (
              <div key={item.step} className="flex items-center">
                <div className={`flex flex-col items-center ${item.step <= activeStep ? 'opacity-100' : 'opacity-40'}`}>
                  <div className={`w-16 h-16 rounded-full flex items-center justify-center text-2xl transition-all ${
                    item.step === activeStep 
                      ? 'bg-gradient-to-br from-purple-600 to-blue-600 shadow-xl scale-110' 
                      : item.step < activeStep
                      ? 'bg-green-500 shadow-lg'
                      : 'bg-gray-200'
                  }`}>
                    {item.step < activeStep ? '✓' : item.icon}
                  </div>
                  <span className={`mt-2 text-sm font-medium ${
                    item.step === activeStep ? 'text-purple-600' : 'text-gray-600'
                  }`}>
                    {item.label}
                  </span>
                </div>
                {item.step < 4 && (
                  <div className={`w-16 h-1 mx-2 rounded-full ${
                    item.step < activeStep ? 'bg-green-500' : 'bg-gray-200'
                  }`} />
                )}
              </div>
            ))}
          </div>
        </div>

        {error && (
          <div className="mb-8 bg-red-50 border-l-4 border-red-500 p-6 rounded-lg shadow-md animate-shake">
            <div className="flex items-center">
              <span className="text-2xl mr-3">⚠️</span>
              <div>
                <h3 className="text-red-800 font-semibold">Error</h3>
                <p className="text-red-700">{error}</p>
              </div>
            </div>
          </div>
        )}

        {/* Step 1: Upload & Details */}
        {activeStep === 1 && (
          <div className="bg-white rounded-2xl shadow-xl p-8 border border-gray-100">
            <div className="flex items-center mb-8">
              <span className="text-4xl mr-4">📤</span>
              <div>
                <h2 className="text-3xl font-bold text-gray-900">Upload & Product Details</h2>
                <p className="text-gray-600">Let's start by uploading your media and describing your product</p>
              </div>
            </div>

            <div className="grid md:grid-cols-2 gap-8 mb-8">
            {/* Product Image Upload */}
              <div className="space-y-4">
                <label className="block text-lg font-semibold text-gray-900">
                  Product Image <span className="text-red-500">*</span>
              </label>
              <div
                {...getProductRootProps()}
                  className={`border-3 border-dashed rounded-2xl p-8 text-center cursor-pointer transition-all hover:scale-105 ${
                  isProductDragActive
                      ? 'border-purple-500 bg-purple-50 shadow-lg'
                      : productImage
                      ? 'border-green-500 bg-green-50'
                      : 'border-gray-300 hover:border-purple-400 hover:bg-purple-50'
                }`}
              >
                <input {...getProductInputProps()} />
                  {productImagePreview ? (
                    <div className="space-y-4">
                      <img src={productImagePreview} alt="Product" className="w-full h-48 object-contain rounded-lg" />
                      <p className="text-green-600 font-semibold">✓ {productImage?.name}</p>
                      <p className="text-sm text-gray-500">Click or drag to replace</p>
                  </div>
                ) : (
                    <div className="space-y-4">
                      <div className="text-6xl">📦</div>
                      <p className="text-gray-700 font-medium">
                        {isProductDragActive ? 'Drop it here!' : 'Drag & drop product image'}
                      </p>
                      <p className="text-sm text-gray-500">PNG, JPG, WEBP • Max 10MB</p>
                  </div>
                )}
              </div>
            </div>

              {/* Person Photo Upload */}
              <div className="space-y-4">
                <label className="block text-lg font-semibold text-gray-900">
                  Your Photo <span className="text-red-500">*</span>
              </label>
              <div
                {...getPersonRootProps()}
                  className={`border-3 border-dashed rounded-2xl p-8 text-center cursor-pointer transition-all hover:scale-105 ${
                  isPersonDragActive
                      ? 'border-blue-500 bg-blue-50 shadow-lg'
                      : personMedia
                      ? 'border-green-500 bg-green-50'
                      : 'border-gray-300 hover:border-blue-400 hover:bg-blue-50'
                }`}
              >
                <input {...getPersonInputProps()} />
                  {personMediaPreview ? (
                    <div className="space-y-4">
                      <img src={personMediaPreview} alt="Person" className="w-full h-48 object-contain rounded-lg" />
                      <p className="text-green-600 font-semibold">✓ {personMedia?.name}</p>
                      <p className="text-sm text-gray-500">Click or drag to replace</p>
                  </div>
                ) : (
                    <div className="space-y-4">
                      <div className="text-6xl">👤</div>
                      <p className="text-gray-700 font-medium">
                        {isPersonDragActive ? 'Drop it here!' : 'Drag & drop your photo'}
                      </p>
                      <p className="text-sm text-gray-500">PNG, JPG, WEBP • Max 10MB</p>
                  </div>
                )}
              </div>
            </div>
          </div>

            {/* Product Details Form */}
            <div className="space-y-6">
              <div className="grid md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-semibold text-gray-700 mb-2">
                    Product Name
                  </label>
                  <input
                    type="text"
                    value={productName}
                    onChange={(e) => setProductName(e.target.value)}
                    placeholder="e.g., iPhone 15 Pro Max"
                    className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all"
                  />
                </div>

                <div>
                  <label className="block text-sm font-semibold text-gray-700 mb-2">
                    Price
                  </label>
                  <input
                    type="text"
                    value={productPrice}
                    onChange={(e) => setProductPrice(e.target.value)}
                    placeholder="e.g., $999"
                    className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Category
                </label>
                <select
                  value={productCategory}
                  onChange={(e) => setProductCategory(e.target.value)}
                  className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all"
                >
                  <option value="electronics">📱 Electronics</option>
                  <option value="fashion">👗 Fashion</option>
                  <option value="home">🏠 Home & Living</option>
                  <option value="beauty">💄 Beauty & Personal Care</option>
                  <option value="sports">⚽ Sports & Outdoors</option>
                  <option value="books">📚 Books & Media</option>
                  <option value="toys">🎮 Toys & Games</option>
                  <option value="food">🍔 Food & Beverages</option>
                  <option value="other">🎯 Other</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Product Description <span className="text-red-500">*</span>
            </label>
            <textarea
                  value={productDescription}
                  onChange={(e) => setProductDescription(e.target.value)}
                  placeholder="Describe your product... What makes it special? What features does it have? Who is it for?&#10;&#10;Example: This premium smartphone features a stunning 6.7-inch display, powerful A17 Pro chip, professional camera system with 5x optical zoom, and all-day battery life. Perfect for professionals and content creators who demand the best performance."
                  className="w-full px-4 py-4 border-2 border-gray-200 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent resize-none transition-all"
                  rows={6}
                />
                <div className="mt-2 flex items-center justify-between">
                  <p className="text-xs text-gray-500">
                    💡 AI will generate a compelling 15-second video script from this description
                  </p>
                  <p className="text-xs text-gray-400">{productDescription.length} characters</p>
                </div>
          </div>

              {/* Video Options */}
              <div className="bg-gradient-to-br from-purple-50 to-blue-50 rounded-xl p-6 space-y-4">
                <h3 className="text-lg font-bold text-gray-900 flex items-center">
                  <span className="mr-2">🎨</span>
                  Video Customization
                </h3>
                
                <div className="grid md:grid-cols-2 gap-4">
            <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">
                Product Animation Style
              </label>
              <select
                value={productVideoStyle}
                onChange={(e) => setProductVideoStyle(e.target.value)}
                      className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all bg-white"
                    >
                      <option value="cinematic">🎬 Cinematic (Dynamic & Engaging)</option>
                      <option value="rotation">🔄 360° Rotation</option>
                      <option value="zoom">🔍 Zoom In (Show Details)</option>
                      <option value="pan">📷 Pan Around</option>
                      <option value="reveal">✨ Dramatic Reveal</option>
                      <option value="auto">🤖 Auto (AI Decides)</option>
              </select>
            </div>

            <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">
                Video Layout
              </label>
              <div className="w-full px-4 py-3 border-2 border-purple-200 rounded-xl bg-purple-50 flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <span className="text-2xl">📦</span>
                  <div>
                    <div className="font-bold text-purple-900">Product Centered</div>
                    <div className="text-sm text-purple-700">Fullscreen product + Person in corner</div>
                  </div>
                </div>
                <span className="text-purple-500 font-semibold">✓ Active</span>
              </div>
              <input type="hidden" value={layout} />
                  </div>
            </div>
          </div>
            </div>

          <button
            onClick={handleUpload}
              disabled={loading || !productImage || !personMedia || !productDescription.trim()}
              className="mt-8 w-full bg-gradient-to-r from-purple-600 to-blue-600 text-white py-4 px-8 rounded-xl font-bold text-lg hover:from-purple-700 hover:to-blue-700 disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed transition-all shadow-lg hover:shadow-xl transform hover:scale-105"
            >
              {loading ? (
                <span className="flex items-center justify-center">
                  <svg className="animate-spin h-6 w-6 mr-3" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                  </svg>
                  Processing...
                </span>
              ) : (
                'Continue to Video Generation →'
              )}
          </button>
        </div>
        )}

        {/* Step 2: Generate Video */}
        {project && activeStep >= 2 && (
          <div className="bg-white rounded-2xl shadow-xl p-8 border border-gray-100 mb-8">
            <div className="flex items-center mb-8">
              <span className="text-4xl mr-4">🎬</span>
              <div>
                <h2 className="text-3xl font-bold text-gray-900">AI Video Generation</h2>
                <p className="text-gray-600">Create a stunning marketing video with AI</p>
              </div>
            </div>

            {project.generated_script && (
              <div className="bg-gradient-to-br from-green-50 to-emerald-50 rounded-xl p-6 mb-6 border-2 border-green-200">
                <h3 className="text-lg font-bold text-green-900 mb-3 flex items-center">
                  <span className="mr-2">✨</span>
                  AI-Generated Script (15 seconds)
                </h3>
                <p className="text-green-800 leading-relaxed italic">"{project.generated_script}"</p>
              </div>
            )}

                {project.status === 'uploaded' && (
                  <button
                    onClick={handleGenerateVideo}
                    disabled={loading}
                className="w-full bg-gradient-to-r from-green-600 to-emerald-600 text-white py-4 px-8 rounded-xl font-bold text-lg hover:from-green-700 hover:to-emerald-700 disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed transition-all shadow-lg hover:shadow-xl transform hover:scale-105"
              >
                {loading ? (
                  <span className="flex items-center justify-center">
                    <svg className="animate-spin h-6 w-6 mr-3" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    Generating Video... (This may take 2-3 minutes)
                  </span>
                ) : (
                  '🚀 Generate Marketing Video'
                )}
                  </button>
                )}

            {project.status === 'video_generating' && (
              <div className="text-center py-12">
                <div className="inline-block animate-spin rounded-full h-16 w-16 border-b-4 border-purple-600 mb-4"></div>
                <p className="text-xl font-semibold text-gray-700">Creating your amazing video...</p>
                <p className="text-gray-500 mt-2">Using AI to generate a professional marketing video</p>
              </div>
                )}

                {project.generated_video_path && (
              <div className="space-y-4">
                <div className="bg-gradient-to-br from-green-50 to-emerald-50 rounded-xl p-6 border-2 border-green-200">
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="text-lg font-bold text-green-900 flex items-center">
                      <span className="mr-2">✅</span>
                      Video Generated Successfully!
                    </h3>
                  <a
                    href={`${API_URL}/static/generated/videos/${project.generated_video_path.split('/').pop()}`}
                    target="_blank"
                    rel="noopener noreferrer"
                      className="px-6 py-2 bg-green-600 text-white rounded-lg font-semibold hover:bg-green-700 transition-all shadow-md hover:shadow-lg"
                    >
                      📥 Download Video
                    </a>
                  </div>
                  <video
                    controls
                    className="w-full rounded-lg shadow-lg"
                    poster={productImagePreview || undefined}
                  >
                    <source src={`${API_URL}/static/generated/videos/${project.generated_video_path.split('/').pop()}`} type="video/mp4" />
                  </video>
                </div>
                
                {/* Next Steps Banner */}
                <div className="bg-gradient-to-r from-purple-100 via-blue-100 to-pink-100 rounded-xl p-6 border-2 border-purple-300 animate-pulse">
                  <h3 className="text-lg font-bold text-purple-900 mb-2 flex items-center">
                    <span className="mr-2">👇</span>
                    Scroll Down for More Options!
                  </h3>
                  <div className="flex flex-wrap gap-3 mb-3">
                    <div className="flex items-center text-purple-800">
                      <span className="mr-2">🌐</span>
                      <span className="font-semibold">Generate Website</span>
                    </div>
                    <span className="text-purple-400">•</span>
                    <div className="flex items-center text-purple-800">
                      <span className="mr-2">📸</span>
                      <span className="font-semibold">Share to Instagram</span>
                    </div>
                  </div>
                  <p className="text-purple-700 text-sm">
                    Both options are available below ⬇️
                  </p>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Step 3: Generate Website & Share to Instagram */}
        {project && project.status === 'video_complete' && activeStep >= 3 && (
          <div className="grid md:grid-cols-2 gap-8 mb-8">
            {/* Website Generation */}
            <div className="bg-white rounded-2xl shadow-xl p-8 border border-gray-100">
              <div className="flex items-center mb-6">
                <span className="text-4xl mr-4">🌐</span>
                <div>
                  <h2 className="text-2xl font-bold text-gray-900">Website Generation</h2>
                  <p className="text-gray-600 text-sm">Create a professional landing page</p>
                </div>
              </div>

              <button
                onClick={handleGenerateWebsite}
                disabled={loading || !!project.website_path}
                className="w-full bg-gradient-to-r from-blue-600 to-indigo-600 text-white py-4 px-6 rounded-xl font-bold hover:from-blue-700 hover:to-indigo-700 disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed transition-all shadow-lg hover:shadow-xl transform hover:scale-105"
              >
                {loading ? (
                  <span className="flex items-center justify-center">
                    <svg className="animate-spin h-5 w-5 mr-3" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    Generating...
                  </span>
                ) : project.website_path ? (
                  '✅ Website Generated!'
                ) : (
                  '🌟 Generate Website'
                )}
              </button>

                {project.website_path && (
                <div className="mt-6 bg-gradient-to-br from-green-50 to-emerald-50 rounded-xl p-4 border-2 border-green-300">
                  <div className="flex items-center justify-between">
                    <div>
                      <h3 className="text-base font-bold text-green-900 flex items-center">
                        <span className="mr-2">✅</span>
                        Website Ready!
                      </h3>
                      <p className="text-sm text-green-700 mt-1">Scroll down to view your website</p>
                    </div>
                    <button
                      onClick={() => setActiveStep(4)}
                      className="px-4 py-2 bg-green-600 text-white rounded-lg font-semibold hover:bg-green-700 transition-all shadow-md hover:shadow-lg"
                    >
                      View →
                    </button>
                  </div>
                </div>
              )}
            </div>

            {/* Instagram Upload */}
            <div className="bg-white rounded-2xl shadow-xl p-8 border border-gray-100">
              <div className="flex items-center mb-6">
                <span className="text-4xl mr-4">📸</span>
                <div>
                  <h2 className="text-2xl font-bold text-gray-900">Share to Instagram</h2>
                  <p className="text-gray-600 text-sm">Post your video directly</p>
                </div>
              </div>

              <div className="bg-gradient-to-br from-pink-50 to-rose-50 rounded-xl p-4 mb-4 border-2 border-pink-200">
                <p className="text-xs text-pink-800">
                  📝 Setup Instagram credentials in environment variables (Business/Creator account required)
                </p>
              </div>

              <button
                onClick={handleUploadToInstagram}
                disabled={loading}
                className="w-full bg-gradient-to-r from-pink-600 to-rose-600 text-white py-4 px-6 rounded-xl font-bold hover:from-pink-700 hover:to-rose-700 disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed transition-all shadow-lg hover:shadow-xl transform hover:scale-105"
              >
                {loading ? (
                  <span className="flex items-center justify-center">
                    <svg className="animate-spin h-5 w-5 mr-3" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    Uploading...
                  </span>
                ) : (
                  '📸 Upload to Instagram'
                )}
              </button>

              {project.instagram_post_url && (
                <div className="mt-6 bg-gradient-to-br from-pink-50 to-rose-50 rounded-xl p-4 border-2 border-pink-200">
                  <h3 className="text-base font-bold text-pink-900 flex items-center mb-2">
                    <span className="mr-2">✅</span>
                    Posted Successfully!
                  </h3>
                  <a
                    href={project.instagram_post_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="w-full mt-3 px-4 py-2 bg-pink-600 text-white rounded-lg font-semibold hover:bg-pink-700 transition-all shadow-md hover:shadow-lg flex items-center justify-center space-x-2"
                  >
                    <span>📱</span>
                    <span>View on Instagram</span>
                  </a>
                </div>
                )}
            </div>
          </div>
        )}

        {/* Step 4: View Website */}
        {project && project.website_path && activeStep >= 4 && (
          <div className="bg-white rounded-2xl shadow-xl p-8 border border-gray-100 mb-8">
            <div className="flex items-center justify-between mb-6">
              <div className="flex items-center">
                <span className="text-4xl mr-4">🌐</span>
                <div>
                  <h2 className="text-2xl font-bold text-gray-900">Your Website</h2>
                  <p className="text-gray-600 text-sm">Professional landing page for {project.product_name}</p>
                </div>
              </div>
              <a
                href={`${API_URL}/static/${project.website_path}/index.html`}
                target="_blank"
                rel="noopener noreferrer"
                className="px-6 py-3 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-xl font-bold hover:from-blue-700 hover:to-indigo-700 transition-all shadow-lg hover:shadow-xl transform hover:scale-105 flex items-center space-x-2"
              >
                <span>🔗</span>
                <span>Open in New Tab</span>
              </a>
            </div>

            {/* Website Preview */}
            <div className="bg-gray-50 rounded-xl p-4 border-2 border-gray-200">
              <div className="bg-white rounded-lg shadow-lg overflow-hidden">
                <div className="bg-gray-800 px-4 py-2 flex items-center space-x-2">
                  <div className="flex space-x-2">
                    <div className="w-3 h-3 rounded-full bg-red-500"></div>
                    <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                    <div className="w-3 h-3 rounded-full bg-green-500"></div>
                  </div>
                  <div className="flex-1 bg-gray-700 rounded px-3 py-1 text-xs text-gray-300 font-mono">
                    {API_URL}/static/{project.website_path}/index.html
                  </div>
                </div>
                <iframe
                  src={`${API_URL}/static/${project.website_path}/index.html`}
                  className="w-full"
                  style={{ height: '600px', border: 'none' }}
                  title="Website Preview"
                />
              </div>
              
              <div className="mt-4 flex items-center justify-between bg-blue-50 rounded-lg p-4 border border-blue-200">
                <div className="flex items-center space-x-3">
                  <span className="text-2xl">✨</span>
                  <div>
                    <p className="font-semibold text-blue-900">Your website is live!</p>
                    <p className="text-sm text-blue-700">Share this URL with your customers</p>
                  </div>
                </div>
                <button
                  onClick={() => {
                    const url = `${API_URL}/static/${project.website_path}/index.html`
                    navigator.clipboard.writeText(url)
                    alert('Website URL copied to clipboard!')
                  }}
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-all shadow-md hover:shadow-lg flex items-center space-x-2"
                >
                  <span>📋</span>
                  <span>Copy URL</span>
                </button>
              </div>
            </div>
          </div>
        )}

      </div>

      {/* Footer */}
      <footer className="bg-gray-900 text-white py-12 mt-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <p className="text-lg mb-4">🚀 AI Product Marketing Studio</p>
          <p className="text-gray-400">Create, Generate, and Share with AI • Built with ❤️ using Next.js, Go, D-ID, RunwayML & Shotstack</p>
        </div>
      </footer>
    </main>
  )
}
