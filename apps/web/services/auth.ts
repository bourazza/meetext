import api, { setTokens, clearTokens } from '@/lib/api'
import type { AuthResponse, TokenPair, User } from '@/types'

export interface RegisterInput {
  full_name: string
  email: string
  password: string
  workspace_name: string
}

export interface LoginInput {
  email: string
  password: string
  remember_me?: boolean
}

export async function register(input: RegisterInput): Promise<AuthResponse> {
  const { data } = await api.post<{ success: boolean; data: AuthResponse }>('/auth/register', input)
  setTokens(data.data.access_token, data.data.refresh_token)
  return data.data
}

export async function login(input: LoginInput): Promise<AuthResponse> {
  const { data } = await api.post<{ success: boolean; data: AuthResponse }>('/auth/login', input)
  setTokens(data.data.access_token, data.data.refresh_token, input.remember_me)
  return data.data
}

export async function getCurrentUser(): Promise<User> {
  const { data } = await api.get<{ success: boolean; data: User }>('/auth/me')
  return data.data
}

export async function refreshToken(refreshToken: string): Promise<TokenPair> {
  const { data } = await api.post<{ success: boolean; data: TokenPair }>('/auth/refresh', {
    refresh_token: refreshToken,
  })
  setTokens(data.data.access_token, data.data.refresh_token)
  return data.data
}

export async function forgotPassword(email: string): Promise<void> {
  await api.post('/auth/forgot-password', { email })
}

export async function resetPassword(token: string, password: string): Promise<void> {
  await api.post('/auth/reset-password', { token, password })
}

export async function verifyEmail(token: string): Promise<void> {
  await api.post('/auth/verify-email', { token })
}

export async function resendVerification(email: string): Promise<void> {
  await api.post('/auth/resend-verification', { email })
}

export async function logout() {
  try {
    await api.post('/auth/logout')
  } catch {
    // Local cleanup still wins if the network is unavailable.
  }
  clearTokens()
  window.location.href = '/login'
}
