'use client'

import { Suspense, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { getCurrentUser } from '@/services/auth'
import { getWorkspaces } from '@/services/workspace'
import { useAuthStore } from '@/store/auth'

export default function OAuthCallbackPage() {
  return (
    <Suspense fallback={<CallbackLoading />}>
      <OAuthCallback />
    </Suspense>
  )
}

function OAuthCallback() {
  const router = useRouter()
  const params = useSearchParams()
  const { setUser, setWorkspace } = useAuthStore()

  useEffect(() => {
    const error = params.get('error')

    if (error) {
      router.replace(`/login?error=${error ?? 'oauth_failed'}`)
      return
    }

    Promise.all([getCurrentUser(), getWorkspaces()])
      .then(([user, workspaces]) => {
        setUser(user)
        if (workspaces.length > 0) setWorkspace(workspaces[0])
        router.replace('/dashboard')
      })
      .catch(() => {
        router.replace('/dashboard')
      })
  }, [params, router, setUser, setWorkspace])

  return <CallbackLoading />
}

function CallbackLoading() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-[#f7f8fb] px-4">
      <div className="rounded-lg border border-zinc-200 bg-white p-8 text-center shadow-[0_24px_80px_rgba(15,23,42,0.08)]">
        <div className="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-2 border-zinc-200 border-t-zinc-950" />
        <p className="text-sm font-medium text-zinc-950">Finishing secure sign in</p>
        <p className="mt-1 text-sm text-zinc-500">Your workspace will open in a moment.</p>
      </div>
    </div>
  )
}
