import axios, { AxiosError } from 'axios'
import type { TokenPair } from '@/types'

const baseURL = process.env.NEXT_PUBLIC_API_URL ? `${process.env.NEXT_PUBLIC_API_URL}/api/v1` : '/api/v1'

const api = axios.create({
  baseURL,
  headers: { 'Content-Type': 'application/json' },
  withCredentials: true,
  timeout: 15000,
})

let refreshPromise: Promise<TokenPair> | null = null

// On 401, ask the API to rotate the HttpOnly refresh-cookie session once.
api.interceptors.response.use(
  (res) => res,
  async (error: AxiosError) => {
    const original = error.config as (typeof error.config & { _retry?: boolean }) | undefined

    if (error.response?.status === 401 && original && !original._retry && !original.url?.includes('/auth/refresh')) {
      original._retry = true
      try {
        refreshPromise =
          refreshPromise ??
          api
            .post<{ success: boolean; data: TokenPair }>('/auth/refresh', {})
            .then((res) => res.data.data)
            .finally(() => {
              refreshPromise = null
            })
        await refreshPromise
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
  return null
}

export function getRefreshToken(): string | null {
  return null
}

export function setTokens(_access: string, _refresh: string, _remember = true) {
  // Tokens are delivered by the API as HttpOnly cookies. This compatibility
  // hook is intentionally a no-op for older call sites.
}

export function clearTokens() {
  if (typeof window === 'undefined') return
  document.cookie = 'meetext_access=; path=/; max-age=0; SameSite=Lax'
  document.cookie = 'meetext_refresh=; path=/; max-age=0; SameSite=Lax'
}
