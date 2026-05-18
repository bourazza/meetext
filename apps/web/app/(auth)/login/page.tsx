'use client'

import { Suspense, useEffect } from 'react'
import Link from 'next/link'
import { useRouter, useSearchParams } from 'next/navigation'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { toast } from 'sonner'
import { z } from 'zod'
import {
  AuthLink,
  AuthPanel,
  Divider,
  Field,
  OAuthButtons,
  StatusCard,
  SubmitButton,
  inputClass,
} from '@/components/auth/auth-ui'
import { login } from '@/services/auth'
import { useAuthStore } from '@/store/auth'

const schema = z.object({
  email: z.string().email('Enter a work email address.'),
  password: z.string().min(1, 'Enter your password.'),
  remember_me: z.boolean().default(true),
})

type FormData = z.infer<typeof schema>

export default function LoginPage() {
  return (
    <Suspense fallback={<LoginFallback />}>
      <LoginCard />
    </Suspense>
  )
}

function LoginCard() {
  const router = useRouter()
  const params = useSearchParams()
  const { setUser, setWorkspace } = useAuthStore()
  const redirectTo = params.get('next') || '/dashboard'

  useEffect(() => {
    const error = params.get('error')
    if (error) toast.error(error.replace(/_/g, ' '))
  }, [params])

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { remember_me: true },
  })

  const onSubmit = async (values: FormData) => {
    try {
      const res = await login(values)
      setUser(res.user)
      if (res.workspace) setWorkspace(res.workspace)
      router.push(redirectTo)
    } catch (err: any) {
      toast.error(err?.response?.data?.error?.message ?? 'Could not sign you in.')
    }
  }

  return (
    <AuthPanel
      eyebrow="Welcome back"
      title="Sign in to Meetext"
      subtitle="Continue to your meeting workspace and pick up exactly where the client conversation left off."
      footer={
        <>
          New to Meetext? <AuthLink href="/register">Create an account</AuthLink>
        </>
      }
    >
      {params.get('verified') === 'true' && (
        <div className="mb-5">
          <StatusCard tone="success" title="Email verified">
            You are all set. Sign in to continue.
          </StatusCard>
        </div>
      )}
      {params.get('reset') === 'true' && (
        <div className="mb-5">
          <StatusCard tone="success" title="Password updated">
            Your new password is ready to use.
          </StatusCard>
        </div>
      )}

      <OAuthButtons />
      <Divider />

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <Field label="Email" error={errors.email?.message}>
          <input {...register('email')} type="email" autoComplete="email" placeholder="you@studio.com" className={inputClass} />
        </Field>

        <Field label="Password" error={errors.password?.message}>
          <input {...register('password')} type="password" autoComplete="current-password" placeholder="Enter your password" className={inputClass} />
        </Field>

        <div className="flex items-center justify-between gap-3">
          <label className="inline-flex items-center gap-2 text-sm text-zinc-600">
            <input {...register('remember_me')} type="checkbox" className="h-4 w-4 rounded border-zinc-300 text-zinc-950" />
            Remember me
          </label>
          <Link href="/forgot-password" className="text-sm font-medium text-zinc-950 underline-offset-4 hover:underline">
            Forgot password?
          </Link>
        </div>

        <SubmitButton loading={isSubmitting}>Sign in</SubmitButton>
      </form>
    </AuthPanel>
  )
}

function LoginFallback() {
  return (
    <AuthPanel eyebrow="Welcome back" title="Sign in to Meetext" subtitle="Preparing your secure sign-in experience.">
      <div className="h-56 animate-pulse rounded-md bg-zinc-100" />
    </AuthPanel>
  )
}
