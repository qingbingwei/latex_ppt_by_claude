import request from '@/utils/request'
import type { PPTRecord, GeneratePPTRequest } from '@/types'

export function generatePPT(data: GeneratePPTRequest): Promise<PPTRecord> {
  return request.post('/ppt/generate', data)
}

export function getTemplates(): Promise<{ templates: string[] }> {
  return request.get('/ppt/templates')
}

export function compileLaTeX(latexContent: string): Promise<PPTRecord> {
  return request.post('/ppt/compile', { latex_content: latexContent })
}

export function getPPTHistory(): Promise<PPTRecord[]> {
  return request.get('/ppt/history')
}

export function getPPT(id: number): Promise<PPTRecord> {
  return request.get(`/ppt/${id}`)
}

export function downloadPPT(id: number): string {
  return `/api/v1/ppt/${id}/download`
}

export function deletePPT(id: number): Promise<void> {
  return request.delete(`/ppt/${id}`)
}
