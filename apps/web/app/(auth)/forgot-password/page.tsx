'use client'

import { useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { AuthLink, AuthPanel, Field, MagicMailIcon, StatusCard, SubmitButton, inputClass } from '@/components/auth/auth-ui'
import { forgotPassword } from '@/services/auth'

const schema = z.object({
  email: z.string().email('Enter the email on your account.'),
})

type FormData = z.infer<typeof schema>

export default function ForgotPasswordPage() {
  const [sentTo, setSentTo] = useState<string | null>(null)
  const form = useForm<FormData>({ resolver: zodResolver(schema) })

  const onSubmit = async (values: FormData) => {
    await forgotPassword(values.email)
    setSentTo(values.email)
  }

  return (
    <AuthPanel
      eyebrow="Account recovery"
      title="Reset your password"
      subtitle="Enter your email and we will send a secure reset link if the account exists."
      footer={
        <>
          Remembered it? <AuthLink href="/login">Back to sign in</AuthLink>
        </>
      }
    >
      {sentTo ? (
        <div className="text-center">
          <MagicMailIcon />
          <StatusCard tone="success" title="Check your inbox">
            We sent password reset instructions to {sentTo}. The link expires in 30 minutes.
          </StatusCard>
        </div>
      ) : (
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <Field label="Email" error={form.formState.errors.email?.message}>
            <input {...form.register('email')} type="email" autoComplete="email" placeholder="you@studio.com" className={inputClass} />
          </Field>
          <SubmitButton loading={form.formState.isSubmitting}>Send reset link</SubmitButton>
        </form>
      )}
    </AuthPanel>
  )
}
