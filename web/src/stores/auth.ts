import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '../api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('accessToken') || '')
  const userInfo = ref<{ username: string } | null>(null)

  async function login(username: string, password: string) {
    const data = await authApi.login({ username, password })
    token.value = data.token
    localStorage.setItem('accessToken', data.token)
    if (data.refresh_token) {
      localStorage.setItem('refreshToken', data.refresh_token)
    }
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
  }

  return { token, userInfo, login, logout }
})
