export interface User {
  id: number
  username: string
  email: string
  created_at: string
  updated_at: string
}

export interface Document {
  id: number
  user_id: number
  filename: string
  file_type: string
  file_size: number
  file_path: string
  status: 'pending' | 'processing' | 'completed' | 'failed'
  chunk_count: number
  created_at: string
  updated_at: string
}

export interface PPTRecord {
  id: number
  user_id: number
  title: string
  prompt: string
  latex_content: string
  pdf_path: string
  template: string
  status: 'pending' | 'generating' | 'completed' | 'failed'
  error_message?: string
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface GeneratePPTRequest {
  title: string
  prompt: string
  template?: string
  document_ids?: number[]
  use_openai?: boolean
}

export interface SearchResult {
  ChunkID: number
  DocumentID: number
  Content: string
  Score: number
}
