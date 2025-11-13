import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export interface Project {
  id: string
  product_image_path: string
  person_media_path: string
  person_media_type: string
  generated_video_path?: string
  website_path?: string
  website_url?: string
  status: string
  created_at: string
  updated_at: string
}

export interface UploadResponse {
  project_id: string
  status: string
  message: string
}

export interface GenerateResponse {
  project_id: string
  video_path?: string
  website_path?: string
  status: string
}

const api = axios.create({
  baseURL: `${API_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const uploadMedia = async (
  productImage: File,
  personMedia: File
): Promise<UploadResponse> => {
  const formData = new FormData()
  formData.append('product_image', productImage)
  formData.append('person_media', personMedia)

  const response = await api.post<UploadResponse>('/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })

  return response.data
}

export interface VideoGenerationOptions {
  script?: string
  product_video_style?: 'rotation' | 'zoom' | 'pan' | 'reveal' | 'auto'
  layout?: 'product_main' | 'avatar_main'
}

export const generateVideo = async (
  projectId: string, 
  options?: VideoGenerationOptions
): Promise<GenerateResponse> => {
  const response = await api.post<GenerateResponse>(
    `/projects/${projectId}/generate-video`,
    options
  )
  return response.data
}

export const generateWebsite = async (projectId: string): Promise<GenerateResponse> => {
  const response = await api.post<GenerateResponse>(
    `/projects/${projectId}/generate-website`
  )
  return response.data
}

export const getProjects = async (): Promise<Project[]> => {
  const response = await api.get<{ projects: Project[] }>('/projects')
  return response.data.projects
}

export const getProject = async (projectId: string): Promise<Project> => {
  const response = await api.get<Project>(`/projects/${projectId}`)
  return response.data
}


