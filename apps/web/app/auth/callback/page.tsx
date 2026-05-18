'use client'

import { useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { setTokens } from '@/lib/api'
import { getWorkspaces } from '@/services/workspace'
import { useAuthStore } from '@/store/auth'

export default function OAuthCallbackPage() {
  const router = useRouter()
  const params = useSearchParams()
  const { setWorkspace } = useAuthStore()

  useEffect(() => {
    const accessToken = params.get('access_token')
    const refreshToken = params.get('refresh_token')
    const error = params.get('error')

    if (error || !accessToken || !refreshToken) {
      router.replace(`/login?error=${error ?? 'oauth_failed'}`)
      return
    }

    setTokens(accessToken, refreshToken)

    // Fetch workspace so the store is populated before hitting the dashboard
    getWorkspaces()
      .then((workspaces) => {
        if (workspaces.length > 0) setWorkspace(workspaces[0])
        router.replace('/dashboard')
      })
      .catch(() => {
        router.replace('/dashboard')
      })
  }, [params, router, setWorkspace])

  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <p className="text-sm text-muted-foreground">Signing you in…</p>
    </div>
  )
}
