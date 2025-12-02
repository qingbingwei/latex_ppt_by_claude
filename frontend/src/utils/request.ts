import axios, { AxiosError, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { getToken, removeToken } from './storage'
import router from '@/router'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 60000
})

// Request interceptor
request.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

// Response interceptor
request.interceptors.response.use(
  (response: AxiosResponse) => {
    return response.data
  },
  (error: AxiosError) => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          ElMessage.error('Unauthorized. Please login.')
          removeToken()
          router.push('/login')
          break
        case 403:
          ElMessage.error('Access denied.')
          break
        case 404:
          ElMessage.error('Resource not found.')
          break
        case 500:
          ElMessage.error('Server error. Please try again later.')
          break
        default:
          ElMessage.error(
            (error.response.data as any)?.error || 'An error occurred'
          )
      }
    } else if (error.request) {
      ElMessage.error('Network error. Please check your connection.')
    } else {
      ElMessage.error('An error occurred.')
    }
    return Promise.reject(error)
  }
)

export default request
