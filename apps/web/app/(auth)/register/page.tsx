'use client'

import { useRouter } from 'next/navigation'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm, useWatch } from 'react-hook-form'
import { toast } from 'sonner'
import { z } from 'zod'
import {
  AuthLink,
  AuthPanel,
  Divider,
  Field,
  OAuthButtons,
  PasswordStrength,
  SubmitButton,
  inputClass,
} from '@/components/auth/auth-ui'
import { register as registerUser } from '@/services/auth'
import { useAuthStore } from '@/store/auth'

const schema = z.object({
  full_name: z.string().min(2, 'Enter your full name.').max(100, 'Name is too long.'),
  email: z.string().email('Enter a valid email address.'),
  password: z.string().min(8, 'Use at least 8 characters.'),
})

type FormData = z.infer<typeof schema>

export default function RegisterPage() {
  const router = useRouter()
  const { setUser, setWorkspace } = useAuthStore()
  const form = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { full_name: '', email: '', password: '' },
  })
  const password = useWatch({ control: form.control, name: 'password' }) ?? ''
  const fullName = useWatch({ control: form.control, name: 'full_name' }) ?? ''
  const firstName = fullName.trim().split(/\s+/)[0]
  const workspaceName = firstName ? `${firstName}'s Workspace` : 'My Workspace'

  const onSubmit = async (values: FormData) => {
    try {
      const res = await registerUser({ ...values, workspace_name: workspaceName })
      setUser(res.user)
      if (res.workspace) setWorkspace(res.workspace)
      toast.success('Account created. Check your inbox to verify your email.')
      router.push('/dashboard')
    } catch (err: any) {
      toast.error(err?.response?.data?.error?.message ?? 'Could not create your account.')
    }
  }

  return (
    <AuthPanel
      eyebrow="Start free"
      title="Create your Meetext account"
      subtitle="Set up your workspace in less than a minute. No credit card, no setup maze."
      footer={
        <>
          Already have an account? <AuthLink href="/login">Sign in</AuthLink>
        </>
      }
    >
      <OAuthButtons mode="signup" />
      <Divider />

      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <Field label="Full name" error={form.formState.errors.full_name?.message}>
          <input {...form.register('full_name')} autoComplete="name" placeholder="Zaki Nirvana" className={inputClass} />
        </Field>

        <Field label="Email" error={form.formState.errors.email?.message}>
          <input {...form.register('email')} type="email" autoComplete="email" placeholder="you@studio.com" className={inputClass} />
        </Field>

        <Field label="Password" error={form.formState.errors.password?.message}>
          <input {...form.register('password')} type="password" autoComplete="new-password" placeholder="Create a secure password" className={inputClass} />
        </Field>

        <PasswordStrength password={password} />

        <SubmitButton loading={form.formState.isSubmitting}>Create account</SubmitButton>
      </form>
    </AuthPanel>
  )
}
