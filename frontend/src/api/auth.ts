import request from '@/utils/request'
import type { LoginRequest, RegisterRequest, AuthResponse, User } from '@/types'

export function login(data: LoginRequest): Promise<AuthResponse> {
  return request.post('/auth/login', data)
}

export function register(data: RegisterRequest): Promise<AuthResponse> {
  return request.post('/auth/register', data)
}

export function getProfile(): Promise<User> {
  return request.get('/auth/profile')
}
