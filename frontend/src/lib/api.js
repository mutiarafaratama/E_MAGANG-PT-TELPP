import axios from 'axios'
import { isTokenExpired } from '@/hooks/useAuth'

const baseURL = import.meta.env.VITE_API_BASE_URL ?? ''

const api = axios.create({
  baseURL,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  if (isTokenExpired()) {
    localStorage.removeItem('access_token')
    localStorage.removeItem('user')
    localStorage.removeItem('token_expires_at')
    window.location.href = '/login?expired=1'
    return Promise.reject(new Error('Token kedaluwarsa'))
  }
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const isAuthEndpoint = error.config?.url?.includes('/auth/login')
      if (!isAuthEndpoint) {
        localStorage.removeItem('access_token')
        localStorage.removeItem('user')
        localStorage.removeItem('token_expires_at')
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default api
