import api, { setTokens, clearTokens } from '@/lib/api'
import type { AuthResponse, TokenPair } from '@/types'

export interface RegisterInput {
  full_name: string
  email: string
  password: string
  workspace_name: string
}

export interface LoginInput {
  email: string
  password: string
}

export async function register(input: RegisterInput): Promise<AuthResponse> {
  const { data } = await api.post<{ success: boolean; data: AuthResponse }>('/auth/register', input)
  setTokens(data.data.access_token, data.data.refresh_token)
  return data.data
}

export async function login(input: LoginInput): Promise<AuthResponse> {
  const { data } = await api.post<{ success: boolean; data: AuthResponse }>('/auth/login', input)
  setTokens(data.data.access_token, data.data.refresh_token)
  return data.data
}

export async function refreshToken(refreshToken: string): Promise<TokenPair> {
  const { data } = await api.post<{ success: boolean; data: TokenPair }>('/auth/refresh', {
    refresh_token: refreshToken,
  })
  setTokens(data.data.access_token, data.data.refresh_token)
  return data.data
}

export function logout() {
  clearTokens()
  window.location.href = '/login'
}
