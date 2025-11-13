'use client'

import { useState, useCallback } from 'react'
import { useDropzone } from 'react-dropzone'
import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

interface Project {
  id: string
  status: string
  product_image_path: string
  person_media_path: string
  generated_video_path?: string
  website_path?: string
  website_url?: string
}

export default function Home() {
  const [productImage, setProductImage] = useState<File | null>(null)
  const [personMedia, setPersonMedia] = useState<File | null>(null)
  const [videoScript, setVideoScript] = useState<string>('')
  const [productVideoStyle, setProductVideoStyle] = useState<string>('auto')
  const [layout, setLayout] = useState<string>('product_main')
  const [project, setProject] = useState<Project | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const onProductImageDrop = useCallback((acceptedFiles: File[]) => {
    if (acceptedFiles.length > 0) {
      setProductImage(acceptedFiles[0])
    }
  }, [])

  const onPersonMediaDrop = useCallback((acceptedFiles: File[]) => {
    if (acceptedFiles.length > 0) {
      setPersonMedia(acceptedFiles[0])
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
  })

  const {
    getRootProps: getPersonRootProps,
    getInputProps: getPersonInputProps,
    isDragActive: isPersonDragActive,
  } = useDropzone({
    onDrop: onPersonMediaDrop,
    accept: {
      'image/*': ['.png', '.jpg', '.jpeg', '.gif', '.webp'],
      'video/*': ['.mp4', '.mov', '.avi'],
    },
  })

  const handleUpload = async () => {
    if (!productImage || !personMedia) {
      setError('Please upload both product image and person media')
      return
    }

    setLoading(true)
    setError(null)

    const formData = new FormData()
    formData.append('product_image', productImage)
    formData.append('person_media', personMedia)

    try {
      const response = await axios.post(`${API_URL}/api/v1/upload`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })

      setProject(response.data)
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
          script: videoScript || undefined,
          product_video_style: productVideoStyle,
          layout: layout
        }
      )
      setProject({ ...project, ...response.data })
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
      setProject({ ...project, ...response.data })
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to generate website')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-6xl mx-auto">
        <h1 className="text-4xl font-bold text-center mb-8 text-primary-700">
          AI Product Marketing Agent
        </h1>

        <div className="bg-white rounded-lg shadow-lg p-8 mb-8">
          <h2 className="text-2xl font-semibold mb-6">Upload Media</h2>

          <div className="grid md:grid-cols-2 gap-6 mb-6">
            {/* Product Image Upload */}
            <div>
              <label className="block text-sm font-medium mb-2">
                Product Image
              </label>
              <div
                {...getProductRootProps()}
                className={`border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors ${
                  isProductDragActive
                    ? 'border-primary-500 bg-primary-50'
                    : 'border-gray-300 hover:border-primary-400'
                }`}
              >
                <input {...getProductInputProps()} />
                {productImage ? (
                  <div>
                    <p className="text-green-600 font-medium">
                      {productImage.name}
                    </p>
                    <p className="text-sm text-gray-500 mt-2">
                      Click or drag to replace
                    </p>
                  </div>
                ) : (
                  <div>
                    <p className="text-gray-600">
                      {isProductDragActive
                        ? 'Drop the image here'
                        : 'Drag & drop product image, or click to select'}
                    </p>
                    <p className="text-sm text-gray-500 mt-2">
                      PNG, JPG, GIF, WEBP
                    </p>
                  </div>
                )}
              </div>
            </div>

            {/* Person Media Upload */}
            <div>
              <label className="block text-sm font-medium mb-2">
                Person Photo/Video
              </label>
              <div
                {...getPersonRootProps()}
                className={`border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors ${
                  isPersonDragActive
                    ? 'border-primary-500 bg-primary-50'
                    : 'border-gray-300 hover:border-primary-400'
                }`}
              >
                <input {...getPersonInputProps()} />
                {personMedia ? (
                  <div>
                    <p className="text-green-600 font-medium">
                      {personMedia.name}
                    </p>
                    <p className="text-sm text-gray-500 mt-2">
                      Click or drag to replace
                    </p>
                  </div>
                ) : (
                  <div>
                    <p className="text-gray-600">
                      {isPersonDragActive
                        ? 'Drop the media here'
                        : 'Drag & drop person photo/video, or click to select'}
                    </p>
                    <p className="text-sm text-gray-500 mt-2">
                      Image or Video (MP4, MOV, AVI)
                    </p>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Video Script Input */}
          <div className="mb-6">
            <label className="block text-sm font-medium mb-2">
              Video Script (Optional)
            </label>
            <textarea
              value={videoScript}
              onChange={(e) => setVideoScript(e.target.value)}
              placeholder="Describe what you want to say in the video... (Leave empty for AI-generated script)"
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent resize-none"
              rows={4}
            />
            <p className="text-xs text-gray-500 mt-2">
              💡 Tip: Mention product features, benefits, or leave empty for AI to generate a professional script
            </p>
          </div>

          {/* NEW: Video Options */}
          <div className="grid md:grid-cols-2 gap-6 mb-6">
            {/* Product Video Style */}
            <div>
              <label className="block text-sm font-medium mb-2">
                Product Animation Style
              </label>
              <select
                value={productVideoStyle}
                onChange={(e) => setProductVideoStyle(e.target.value)}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              >
                <option value="auto">🤖 Auto (Let AI decide)</option>
                <option value="rotation">🔄 360° Rotation (Best for gadgets)</option>
                <option value="zoom">🔍 Zoom In (Best for details)</option>
                <option value="pan">📷 Pan Around (Best for large items)</option>
                <option value="reveal">✨ Dramatic Reveal (Best for luxury)</option>
              </select>
              <p className="text-xs text-gray-500 mt-2">
                How should your product be animated in the video?
              </p>
            </div>

            {/* Video Layout */}
            <div>
              <label className="block text-sm font-medium mb-2">
                Video Layout
              </label>
              <select
                value={layout}
                onChange={(e) => setLayout(e.target.value)}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              >
                <option value="product_main">📦 Product Focus (Product fullscreen + Avatar overlay)</option>
                <option value="avatar_main">👤 Presenter Focus (Avatar fullscreen + Product overlay)</option>
              </select>
              <p className="text-xs text-gray-500 mt-2">
                Which element should be the main focus?
              </p>
            </div>
          </div>

          {error && (
            <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
              {error}
            </div>
          )}

          <button
            onClick={handleUpload}
            disabled={loading || !productImage || !personMedia}
            className="w-full bg-primary-600 text-white py-3 px-6 rounded-lg font-semibold hover:bg-primary-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            {loading ? 'Uploading...' : 'Upload & Create Project'}
          </button>
        </div>

        {project && (
          <div className="bg-white rounded-lg shadow-lg p-8">
            <h2 className="text-2xl font-semibold mb-4">Project Status</h2>
            <div className="space-y-4">
              <div>
                <p className="text-sm text-gray-600">Project ID</p>
                <p className="font-mono text-sm">{project.project_id || project.id}</p>
              </div>
              <div>
                <p className="text-sm text-gray-600">Status</p>
                <p className="font-semibold capitalize">{project.status}</p>
              </div>

              <div className="flex gap-4 mt-6">
                {project.status === 'uploaded' && (
                  <button
                    onClick={handleGenerateVideo}
                    disabled={loading}
                    className="bg-primary-600 text-white py-2 px-6 rounded-lg font-semibold hover:bg-primary-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
                  >
                    Generate Video
                  </button>
                )}

                {project.status === 'video_complete' && (
                  <button
                    onClick={handleGenerateWebsite}
                    disabled={loading}
                    className="bg-primary-600 text-white py-2 px-6 rounded-lg font-semibold hover:bg-primary-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
                  >
                    Generate Website
                  </button>
                )}

                {project.generated_video_path && (
                  <a
                    href={`${API_URL}/static/generated/videos/${project.generated_video_path.split('/').pop()}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="bg-green-600 text-white py-2 px-6 rounded-lg font-semibold hover:bg-green-700 transition-colors inline-block"
                  >
                    View Video
                  </a>
                )}

                {project.website_path && (
                  <a
                    href={`${API_URL}/static/generated/websites/${project.website_path.split('/').pop()}/index.html`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="bg-blue-600 text-white py-2 px-6 rounded-lg font-semibold hover:bg-blue-700 transition-colors inline-block"
                  >
                    View Website
                  </a>
                )}
              </div>
            </div>
          </div>
        )}
      </div>
    </main>
  )
}

