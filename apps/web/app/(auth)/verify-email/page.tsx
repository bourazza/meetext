'use client'

import { Suspense, useEffect, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import { AuthPanel, ContinueButton, StatusCard } from '@/components/auth/auth-ui'
import { verifyEmail } from '@/services/auth'

type State = 'loading' | 'success' | 'error'

export default function VerifyEmailPage() {
  return (
    <Suspense fallback={<VerifyShell state="loading" />}>
      <VerifyEmail />
    </Suspense>
  )
}

function VerifyEmail() {
  const params = useSearchParams()
  const token = params.get('token') ?? ''
  const [state, setState] = useState<State>(token ? 'loading' : 'error')

  useEffect(() => {
    if (!token) return
    verifyEmail(token)
      .then(() => setState('success'))
      .catch(() => setState('error'))
  }, [token])

  return <VerifyShell state={state} />
}

function VerifyShell({ state }: { state: State }) {
  return (
    <AuthPanel
      eyebrow="Email verification"
      title={state === 'success' ? 'Email verified' : state === 'error' ? 'Link expired' : 'Verifying your email'}
      subtitle="Verification keeps your workspace secure and helps protect client meeting data."
    >
      {state === 'loading' && <div className="h-24 animate-pulse rounded-md bg-zinc-100" />}
      {state === 'success' && (
        <div className="space-y-5">
          <StatusCard tone="success" title="You are verified">
            Your Meetext account is ready for production work.
          </StatusCard>
          <ContinueButton href="/login?verified=true">Continue</ContinueButton>
        </div>
      )}
      {state === 'error' && (
        <div className="space-y-5">
          <StatusCard tone="error" title="We could not verify this link">
            It may have expired or already been used. Sign in and request a new verification email.
          </StatusCard>
          <ContinueButton href="/login">Back to sign in</ContinueButton>
        </div>
      )}
    </AuthPanel>
  )
}
