import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/lib/api'

const user = ref(null)
const loading = ref(false)
const error = ref(null)

const TOKEN_DURATION_MS = 24 * 60 * 60 * 1000

let expiryTimer = null

const storedUser = localStorage.getItem('user')
if (storedUser) {
  try {
    user.value = JSON.parse(storedUser)
  } catch {
    localStorage.removeItem('user')
  }
}

function dashboardRouteForRole(role) {
  if (role === 'admin') return '/admin'
  if (role === 'hrd') return '/staff'
  return '/dashboard'
}

function clearExpiryTimer() {
  if (expiryTimer) {
    clearInterval(expiryTimer)
    expiryTimer = null
  }
}

export function isTokenExpired() {
  const expiresAt = localStorage.getItem('token_expires_at')
  if (!expiresAt) return false
  return Date.now() >= parseInt(expiresAt, 10)
}

export function startExpiryTimer() {
  clearExpiryTimer()
  expiryTimer = setInterval(() => {
    if (isTokenExpired()) {
      clearExpiryTimer()
      user.value = null
      localStorage.removeItem('access_token')
      localStorage.removeItem('user')
      localStorage.removeItem('token_expires_at')
      window.location.href = '/login?expired=1'
    }
  }, 60_000)
}

export function useAuth() {
  const router = useRouter()

  const isLoggedIn = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isHRD = computed(() => user.value?.role === 'hrd')
  const isPeserta = computed(() => user.value?.role === 'peserta')

  async function login(email, password) {
    loading.value = true
    error.value = null
    try {
      const res = await api.post('/api/auth/login', { email, password })
      const data = res.data.data
      const expiresAt = Date.now() + TOKEN_DURATION_MS
      localStorage.setItem('access_token', data.access_token)
      localStorage.setItem('token_expires_at', String(expiresAt))
      user.value = data.user
      localStorage.setItem('user', JSON.stringify(data.user))
      startExpiryTimer()
      router.push(dashboardRouteForRole(data.user.role))
    } catch (err) {
      error.value = err.response?.data?.message || 'Login gagal. Periksa email dan password.'
    } finally {
      loading.value = false
    }
  }

  async function register(payload) {
    loading.value = true
    error.value = null
    try {
      await api.post('/api/auth/register', payload)
      router.push('/login?registered=1')
    } catch (err) {
      error.value = err.response?.data?.message || 'Pendaftaran gagal. Coba lagi.'
    } finally {
      loading.value = false
    }
  }

  function logout() {
    clearExpiryTimer()
    user.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('user')
    localStorage.removeItem('token_expires_at')
    router.push('/login')
  }

  return {
    user,
    loading,
    error,
    isLoggedIn,
    isAdmin,
    isHRD,
    isPeserta,
    login,
    register,
    logout,
    dashboardRouteForRole,
  }
}
