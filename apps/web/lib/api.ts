import axios, { AxiosError } from 'axios'
import type { TokenPair } from '@/types'

const api = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
  withCredentials: false,
  timeout: 15000,
})

// Inject access token on every request
api.interceptors.request.use((config) => {
  const token = getAccessToken()
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

let refreshPromise: Promise<TokenPair> | null = null

// On 401, try one refresh before sending the user back to sign in.
api.interceptors.response.use(
  (res) => res,
  async (error: AxiosError) => {
    const original = error.config as (typeof error.config & { _retry?: boolean }) | undefined
    const refresh = getRefreshToken()

    if (error.response?.status === 401 && original && !original._retry && refresh && !original.url?.includes('/auth/refresh')) {
      original._retry = true
      try {
        refreshPromise =
          refreshPromise ??
          api
            .post<{ success: boolean; data: TokenPair }>('/auth/refresh', { refresh_token: refresh })
            .then((res) => res.data.data)
            .finally(() => {
              refreshPromise = null
            })
        const tokens = await refreshPromise
        setTokens(tokens.access_token, tokens.refresh_token, true)
        original.headers = original.headers ?? {}
        original.headers.Authorization = `Bearer ${tokens.access_token}`
        return api(original)
      } catch {
        clearTokens()
      }
    }

    if (error.response?.status === 401) {
      clearTokens()
      if (typeof window !== 'undefined') {
        const next = encodeURIComponent(window.location.pathname + window.location.search)
        window.location.href = `/login?next=${next}`
      }
    }
    return Promise.reject(error)
  }
)

export default api

// ── Token helpers ─────────────────────────────────────────────────────────────

export function getAccessToken(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('access_token')
}

export function getRefreshToken(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('refresh_token')
}

export function setTokens(access: string, refresh: string, remember = true) {
  localStorage.setItem('access_token', access)
  localStorage.setItem('refresh_token', refresh)
  // Also set a cookie so Next.js middleware can check auth server-side
  const maxAge = remember ? 7 * 24 * 60 * 60 : 15 * 60
  document.cookie = `meetext_token=${access}; path=/; max-age=${maxAge}; SameSite=Lax`
}

export function clearTokens() {
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
  // Clear the middleware cookie
  document.cookie = 'meetext_token=; path=/; max-age=0; SameSite=Lax'
}
