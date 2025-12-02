import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User, LoginRequest, RegisterRequest } from '@/types'
import * as authApi from '@/api/auth'
import { setToken, setUser, clearAuth, getToken, getUser } from '@/utils/storage'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(getUser())
  const token = ref<string | null>(getToken())

  async function login(data: LoginRequest) {
    const response = await authApi.login(data)
    token.value = response.token
    user.value = response.user
    setToken(response.token)
    setUser(response.user)
  }

  async function register(data: RegisterRequest) {
    const response = await authApi.register(data)
    token.value = response.token
    user.value = response.user
    setToken(response.token)
    setUser(response.user)
  }

  async function fetchProfile() {
    const profile = await authApi.getProfile()
    user.value = profile
    setUser(profile)
  }

  function logout() {
    user.value = null
    token.value = null
    clearAuth()
  }

  const isLoggedIn = () => !!token.value

  return {
    user,
    token,
    login,
    register,
    fetchProfile,
    logout,
    isLoggedIn
  }
})
