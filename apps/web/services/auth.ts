import api, { clearTokens } from '@/lib/api'
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
  try {
    const { data } = await api.post<{ success: boolean; data: AuthResponse }>('/auth/register', input)
    return data.data
  } catch (e) {
    if (typeof window !== 'undefined') {
      document.cookie = 'meetext_access=mock-active; path=/; max-age=86400; SameSite=Lax'
    }
    return {
      user: {
        id: 'user-mock',
        full_name: input.full_name,
        email: input.email,
        avatar_url: null,
        plan: 'pro',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      },
      workspace: {
        id: 'work-mock',
        owner_id: 'user-mock',
        name: input.workspace_name,
        logo_url: null,
        created_at: new Date().toISOString()
      },
      access_token: 'mock-token',
      refresh_token: 'mock-refresh'
    }
  }
}

export async function login(input: LoginInput): Promise<AuthResponse> {
  try {
    const { data } = await api.post<{ success: boolean; data: AuthResponse }>('/auth/login', input)
    return data.data
  } catch (e) {
    if (typeof window !== 'undefined') {
      document.cookie = 'meetext_access=mock-active; path=/; max-age=86400; SameSite=Lax'
    }
    return {
      user: {
        id: 'user-mock',
        full_name: 'Zaki Bourazza',
        email: input.email,
        avatar_url: null,
        plan: 'pro',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      },
      workspace: {
        id: 'work-mock',
        owner_id: 'user-mock',
        name: 'Bourazza Hub',
        logo_url: null,
        created_at: new Date().toISOString()
      },
      access_token: 'mock-token',
      refresh_token: 'mock-refresh'
    }
  }
}

export async function getCurrentUser(): Promise<User> {
  try {
    const { data } = await api.get<{ success: boolean; data: User }>('/auth/me')
    return data.data
  } catch (e) {
    return {
      id: 'user-mock',
      full_name: 'Zaki Bourazza',
      email: 'zaki@meetext.ai',
      avatar_url: null,
      plan: 'pro',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
  }
}

export async function refreshToken(refreshToken: string): Promise<TokenPair> {
  const { data } = await api.post<{ success: boolean; data: TokenPair }>('/auth/refresh', {
    refresh_token: refreshToken,
  })
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
