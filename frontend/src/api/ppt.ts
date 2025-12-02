import request from '@/utils/request'
import type { PPTRecord, GeneratePPTRequest } from '@/types'
import { getToken } from '@/utils/storage'

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

// 获取 PDF Blob URL（用于预览）
export async function getPPTBlobUrl(id: number): Promise<string> {
  const token = getToken()
  const response = await fetch(`/api/v1/ppt/${id}/download`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  })
  
  if (!response.ok) {
    throw new Error('Failed to get PDF')
  }
  
  const blob = await response.blob()
  return window.URL.createObjectURL(blob)
}

// 下载 PPT PDF - 使用 fetch 携带 Token
export async function downloadPPT(id: number, filename?: string): Promise<void> {
  const url = await getPPTBlobUrl(id)
  const a = document.createElement('a')
  a.href = url
  a.download = filename || `ppt_${id}.pdf`
  document.body.appendChild(a)
  a.click()
  window.URL.revokeObjectURL(url)
  document.body.removeChild(a)
}

export function deletePPT(id: number): Promise<void> {
  return request.delete(`/ppt/${id}`)
}
