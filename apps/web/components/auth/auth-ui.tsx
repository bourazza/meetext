'use client'

import Link from 'next/link'
import { ArrowRight, Github, Loader2, Mail } from 'lucide-react'
import { cn } from '@/lib/utils'

export function AuthPanel({
  eyebrow,
  title,
  subtitle,
  children,
  footer,
}: {
  eyebrow: string
  title: string
  subtitle: string
  children: React.ReactNode
  footer?: React.ReactNode
}) {
  return (
    <section className="w-full rounded-2xl border border-zinc-200 bg-white p-10 shadow-sm">
      <div className="mb-10 space-y-2 text-center">
        <h1 className="text-3xl font-semibold tracking-tight text-zinc-950">{title}</h1>
        <p className="text-base text-zinc-500">{subtitle}</p>
      </div>
      {children}
      {footer && <div className="mt-10 text-center text-base text-zinc-600">{footer}</div>}
    </section>
  )
}

export function OAuthButtons({ mode = 'signin' }: { mode?: 'signin' | 'signup' }) {
  const action = mode === 'signup' ? 'Continue' : 'Sign in'

  return (
    <div className="grid gap-3 sm:grid-cols-2">
      <a
        href="/api/v1/auth/oauth/google"
        className="inline-flex h-11 items-center justify-center gap-2 rounded-md border border-zinc-200 bg-white px-3 text-sm font-medium text-zinc-900 transition hover:bg-zinc-50 focus:outline-none focus:ring-2 focus:ring-zinc-900/10"
      >
        <GoogleMark />
        {action} with Google
      </a>
      <a
        href="/api/v1/auth/oauth/github"
        className="inline-flex h-11 items-center justify-center gap-2 rounded-md border border-zinc-200 bg-white px-3 text-sm font-medium text-zinc-900 transition hover:bg-zinc-50 focus:outline-none focus:ring-2 focus:ring-zinc-900/10"
      >
        <Github className="h-4 w-4" />
        GitHub
      </a>
    </div>
  )
}

export function Divider({ label = 'or continue with email' }: { label?: string }) {
  return (
    <div className="relative my-6">
      <div className="absolute inset-0 flex items-center">
        <span className="w-full border-t border-zinc-200" />
      </div>
      <div className="relative flex justify-center text-xs">
        <span className="bg-white px-3 text-zinc-500">{label}</span>
      </div>
    </div>
  )
}

export function Field({
  label,
  error,
  children,
}: {
  label: string
  error?: string
  children: React.ReactNode
}) {
  return (
    <label className="block space-y-2">
      <span className="text-sm font-medium text-zinc-800">{label}</span>
      {children}
      {error && <span className="block text-xs text-red-600">{error}</span>}
    </label>
  )
}

export const inputClass =
  'h-11 w-full rounded-md border border-zinc-200 bg-white px-3 text-sm text-zinc-950 outline-none transition placeholder:text-zinc-400 focus:border-zinc-400 focus:ring-4 focus:ring-zinc-900/5 disabled:cursor-not-allowed disabled:bg-zinc-50'

export function SubmitButton({
  children,
  loading,
  disabled,
}: {
  children: React.ReactNode
  loading?: boolean
  disabled?: boolean
}) {
  return (
    <button
      type="submit"
      disabled={disabled || loading}
      className="inline-flex h-11 w-full items-center justify-center gap-2 rounded-md bg-zinc-950 px-4 text-sm font-medium text-white shadow-sm transition hover:bg-zinc-800 focus:outline-none focus:ring-4 focus:ring-zinc-900/15 disabled:cursor-not-allowed disabled:opacity-60"
    >
      {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
      {children}
    </button>
  )
}

export function AuthLink({ href, children }: { href: string; children: React.ReactNode }) {
  return (
    <Link href={href} className="font-medium text-zinc-950 underline-offset-4 hover:underline">
      {children}
    </Link>
  )
}

export function StatusCard({
  tone = 'neutral',
  title,
  children,
}: {
  tone?: 'neutral' | 'success' | 'error'
  title: string
  children: React.ReactNode
}) {
  return (
    <div
      className={cn(
        'rounded-md border p-4 text-sm',
        tone === 'success' && 'border-emerald-200 bg-emerald-50 text-emerald-950',
        tone === 'error' && 'border-red-200 bg-red-50 text-red-950',
        tone === 'neutral' && 'border-zinc-200 bg-zinc-50 text-zinc-700'
      )}
    >
      <div className="mb-1 font-medium">{title}</div>
      <div className="leading-6">{children}</div>
    </div>
  )
}

export function PasswordStrength({ password }: { password: string }) {
  const checks = [
    password.length >= 8,
    /[A-Z]/.test(password),
    /[0-9]/.test(password),
    /[^A-Za-z0-9]/.test(password),
  ]
  const score = checks.filter(Boolean).length
  const labels = ['Too weak', 'Getting there', 'Good', 'Strong', 'Excellent']

  return (
    <div className="space-y-2">
      <div className="grid grid-cols-4 gap-1">
        {checks.map((_, index) => (
          <span
            key={index}
            className={cn(
              'h-1 rounded-full bg-zinc-200 transition',
              index < score && score < 3 && 'bg-amber-500',
              index < score && score >= 3 && 'bg-emerald-500'
            )}
          />
        ))}
      </div>
      <p className="text-xs text-zinc-500">{labels[score]}</p>
    </div>
  )
}

export function MagicMailIcon() {
  return (
    <div className="mb-5 inline-flex h-11 w-11 items-center justify-center rounded-md border border-zinc-200 bg-white shadow-sm">
      <Mail className="h-5 w-5 text-zinc-800" />
    </div>
  )
}

export function ContinueButton({ href, children }: { href: string; children: React.ReactNode }) {
  return (
    <Link
      href={href}
      className="inline-flex h-11 w-full items-center justify-center gap-2 rounded-md bg-zinc-950 px-4 text-sm font-medium text-white transition hover:bg-zinc-800"
    >
      {children}
      <ArrowRight className="h-4 w-4" />
    </Link>
  )
}

function GoogleMark() {
  return (
    <svg className="h-4 w-4" viewBox="0 0 18 18" aria-hidden="true">
      <path d="M17.64 9.2c0-.637-.057-1.251-.164-1.84H9v3.481h4.844c-.209 1.125-.843 2.078-1.796 2.717v2.258h2.908c1.702-1.567 2.684-3.875 2.684-6.615z" fill="#4285F4" />
      <path d="M9 18c2.43 0 4.467-.806 5.956-2.18l-2.908-2.259c-.806.54-1.837.86-3.048.86-2.344 0-4.328-1.584-5.036-3.711H.957v2.332A8.997 8.997 0 0 0 9 18z" fill="#34A853" />
      <path d="M3.964 10.71A5.41 5.41 0 0 1 3.682 9c0-.593.102-1.17.282-1.71V4.958H.957A8.996 8.996 0 0 0 0 9c0 1.452.348 2.827.957 4.042l3.007-2.332z" fill="#FBBC05" />
      <path d="M9 3.58c1.321 0 2.508.454 3.44 1.345l2.582-2.58C13.463.891 11.426 0 9 0A8.997 8.997 0 0 0 .957 4.958L3.964 7.29C4.672 5.163 6.656 3.58 9 3.58z" fill="#EA4335" />
    </svg>
  )
}
