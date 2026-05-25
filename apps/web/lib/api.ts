import axios, { AxiosError, AxiosRequestConfig } from 'axios'
import type { TokenPair } from '@/types'

const baseURL = process.env.NEXT_PUBLIC_API_URL ? `${process.env.NEXT_PUBLIC_API_URL}/api/v1` : '/api/v1'

const DEBUG_AUTH = process.env.NODE_ENV === 'development'

const api = axios.create({
  baseURL,
  headers: { 'Content-Type': 'application/json' },
  withCredentials: true, // CRITICAL: Send cookies with every request
  timeout: 15000,
})

let refreshPromise: Promise<TokenPair> | null = null
let isRefreshing = false

// Debug logger
function debugLog(message: string, data?: any) {
  if (DEBUG_AUTH) {
    console.log(`[API Auth] ${message}`, data || '')
  }
}

// Request interceptor - add auth headers and debug logging
api.interceptors.request.use(
  (config) => {
    // NOTE: meetext_access and meetext_refresh are HttpOnly cookies — they are NOT
    // readable via document.cookie by design (browser security). Do not try to check
    // for their presence here; the browser sends them automatically on every request
    // because withCredentials: true is set. hasCookies will always be false for HttpOnly cookies.
    debugLog(`Request: ${config.method?.toUpperCase()} ${config.url}`, {
      withCredentials: config.withCredentials,
      headers: config.headers,
    })

    return config
  },
  (error) => {
    debugLog('Request error', error)
    return Promise.reject(error)
  }
)

// Response interceptor - handle 401 with automatic token refresh
api.interceptors.response.use(
  (res) => {
    debugLog(`Response: ${res.config.method?.toUpperCase()} ${res.config.url}`, {
      status: res.status,
    })
    return res
  },
  async (error: AxiosError) => {
    const original = error.config as (typeof error.config & { _retry?: boolean; _retryCount?: number }) | undefined

    if (!original) {
      return Promise.reject(error)
    }

    // Initialize retry count
    if (!original._retryCount) {
      original._retryCount = 0
    }

    debugLog(`Response error: ${original.method?.toUpperCase()} ${original.url}`, {
      status: error.response?.status,
      retryCount: original._retryCount,
      isRefreshing,
    })

    // Handle 401 Unauthorized - attempt token refresh
    if (error.response?.status === 401 && !original._retry && !original.url?.includes('/auth/refresh')) {
      original._retry = true

      // Prevent multiple simultaneous refresh attempts
      if (isRefreshing) {
        debugLog('Already refreshing, waiting...')
        await refreshPromise
        debugLog('Refresh complete, retrying original request')
        return api(original)
      }

      try {
        isRefreshing = true
        debugLog('Attempting token refresh...')

        refreshPromise = api
          .post<{ success: boolean; data: TokenPair }>('/auth/refresh', {})
          .then((res) => {
            debugLog('Token refresh successful')
            return res.data.data
          })
          .finally(() => {
            isRefreshing = false
            refreshPromise = null
          })

        await refreshPromise
        
        // Retry the original request
        debugLog('Retrying original request after refresh')
        return api(original)
      } catch (refreshError) {
        debugLog('Token refresh failed', refreshError)
        isRefreshing = false
        refreshPromise = null
        
        // Only clear tokens and redirect if refresh truly failed
        clearTokens()
        if (typeof window !== 'undefined' && !window.location.pathname.includes('/login')) {
          const next = encodeURIComponent(window.location.pathname + window.location.search)
          debugLog('Redirecting to login', { next })
          window.location.href = `/login?next=${next}`
        }
        return Promise.reject(refreshError)
      }
    }

    // For polling requests, implement retry with exponential backoff
    if (original.url?.includes('/status') && original._retryCount < 3) {
      original._retryCount++
      const delay = Math.min(1000 * Math.pow(2, original._retryCount - 1), 5000)
      
      debugLog(`Retrying status poll in ${delay}ms (attempt ${original._retryCount}/3)`)
      
      await new Promise(resolve => setTimeout(resolve, delay))
      return api(original)
    }

    return Promise.reject(error)
  }
)

// Create a special API instance for long-polling with extended timeout
export const pollingApi = axios.create({
  baseURL,
  headers: { 'Content-Type': 'application/json' },
  withCredentials: true,
  timeout: 30000, // 30s timeout for polling
})

// Apply same interceptors to polling API
pollingApi.interceptors.request.use(
  (config) => {
    const hasCookies = typeof document !== 'undefined' && document.cookie.includes('meetext_access')
    debugLog(`Polling Request: ${config.method?.toUpperCase()} ${config.url}`, { hasCookies })
    return config
  },
  (error) => Promise.reject(error)
)

pollingApi.interceptors.response.use(
  (res) => res,
  async (error: AxiosError) => {
    const original = error.config as (typeof error.config & { _retry?: boolean; _retryCount?: number }) | undefined
    if (!original) return Promise.reject(error)

    // For polling, implement same retry logic
    if (!original._retryCount) original._retryCount = 0

    if (error.response?.status === 401 && !original._retry && !original.url?.includes('/auth/refresh')) {
      original._retry = true
      if (isRefreshing) {
        await refreshPromise
        return pollingApi(original)
      }
      try {
        isRefreshing = true
        refreshPromise = api.post<{ success: boolean; data: TokenPair }>('/auth/refresh', {})
          .then((res) => res.data.data)
          .finally(() => { isRefreshing = false; refreshPromise = null })
        await refreshPromise
        return pollingApi(original)
      } catch {
        isRefreshing = false
        return Promise.reject(error)
      }
    }

    if (original.url?.includes('/status') && original._retryCount < 3) {
      original._retryCount++
      const delay = Math.min(1000 * Math.pow(2, original._retryCount - 1), 5000)
      await new Promise(resolve => setTimeout(resolve, delay))
      return pollingApi(original)
    }

    return Promise.reject(error)
  }
)

export default api

// ── Token helpers ─────────────────────────────────────────────────────────────

export function getAccessToken(): string | null {
  if (typeof document === 'undefined') return null
  const match = document.cookie.match(/meetext_access=([^;]+)/)
  return match ? match[1] : null
}

export function getRefreshToken(): string | null {
  if (typeof document === 'undefined') return null
  const match = document.cookie.match(/meetext_refresh=([^;]+)/)
  return match ? match[1] : null
}

export function hasValidSession(): boolean {
  return !!getAccessToken()
}

export function setTokens(_access: string, _refresh: string, _remember = true) {
  // Tokens are delivered by the API as HttpOnly cookies. This compatibility
  // hook is intentionally a no-op for older call sites.
  debugLog('Tokens set via HttpOnly cookies')
}

export function clearTokens() {
  if (typeof window === 'undefined') return
  debugLog('Clearing tokens')
  document.cookie = 'meetext_access=; path=/; max-age=0; SameSite=Lax'
  document.cookie = 'meetext_refresh=; path=/; max-age=0; SameSite=Lax'
}
