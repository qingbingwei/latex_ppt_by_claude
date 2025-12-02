import request from '@/utils/request'
import type { Document, SearchResult } from '@/types'

export function uploadDocument(file: File): Promise<Document> {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/knowledge/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export function getDocuments(): Promise<Document[]> {
  return request.get('/knowledge/list')
}

export function getDocument(id: number): Promise<Document> {
  return request.get(`/knowledge/${id}`)
}

export function deleteDocument(id: number): Promise<void> {
  return request.delete(`/knowledge/${id}`)
}

export function searchKnowledge(query: string, topK = 5): Promise<SearchResult[]> {
  return request.post('/knowledge/search', { query, top_k: topK })
}
