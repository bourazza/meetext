'use client'

import { Suspense } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm, useWatch } from 'react-hook-form'
import { toast } from 'sonner'
import { z } from 'zod'
import { AuthLink, AuthPanel, Field, PasswordStrength, StatusCard, SubmitButton, inputClass } from '@/components/auth/auth-ui'
import { resetPassword } from '@/services/auth'

const schema = z.object({
  password: z.string().min(8, 'Use at least 8 characters.'),
})

type FormData = z.infer<typeof schema>

export default function ResetPasswordPage() {
  return (
    <Suspense fallback={<ResetFallback />}>
      <ResetPasswordForm />
    </Suspense>
  )
}

function ResetPasswordForm() {
  const router = useRouter()
  const params = useSearchParams()
  const token = params.get('token') ?? ''
  const form = useForm<FormData>({ resolver: zodResolver(schema), defaultValues: { password: '' } })
  const password = useWatch({ control: form.control, name: 'password' }) ?? ''

  const onSubmit = async (values: FormData) => {
    try {
      await resetPassword(token, values.password)
      router.push('/login?reset=true')
    } catch (err: any) {
      toast.error(err?.response?.data?.error?.message ?? 'This reset link is invalid or expired.')
    }
  }

  return (
    <AuthPanel
      eyebrow="Secure reset"
      title="Choose a new password"
      subtitle="Use a password you have not used before. Your old sessions will stop working once this is changed."
      footer={
        <>
          Need a new link? <AuthLink href="/forgot-password">Request one</AuthLink>
        </>
      }
    >
      {!token ? (
        <StatusCard tone="error" title="Missing reset token">
          Open the password reset link from your email.
        </StatusCard>
      ) : (
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <Field label="New password" error={form.formState.errors.password?.message}>
            <input {...form.register('password')} type="password" autoComplete="new-password" placeholder="Create a secure password" className={inputClass} />
          </Field>
          <PasswordStrength password={password} />
          <SubmitButton loading={form.formState.isSubmitting}>Update password</SubmitButton>
        </form>
      )}
    </AuthPanel>
  )
}

function ResetFallback() {
  return (
    <AuthPanel eyebrow="Secure reset" title="Choose a new password" subtitle="Preparing reset form.">
      <div className="h-36 animate-pulse rounded-md bg-zinc-100" />
    </AuthPanel>
  )
}
