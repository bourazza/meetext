'use client'

import type { LucideIcon } from 'lucide-react'
import { AppShell } from '@/components/app/app-shell'

export function EmptyStatePage({
  icon: Icon,
  title,
  description,
}: {
  icon: LucideIcon
  title: string
  description: string
}) {
  return (
    <AppShell>
      <main className="flex min-h-screen items-center justify-center px-5 py-16">
        <section className="w-full max-w-xl text-center">
          <div className="mx-auto mb-6 grid h-14 w-14 place-items-center rounded-xl border border-zinc-200 bg-white text-zinc-700 shadow-sm">
            <Icon className="h-5 w-5" />
          </div>
          <p className="mb-3 text-xs font-semibold uppercase tracking-[0.18em] text-zinc-400">Coming soon</p>
          <h1 className="text-2xl font-semibold tracking-normal text-zinc-950 sm:text-3xl">{title}</h1>
          <p className="mx-auto mt-3 max-w-md text-sm leading-6 text-zinc-500">{description}</p>
          <div className="mx-auto mt-8 h-px max-w-xs bg-gradient-to-r from-transparent via-zinc-200 to-transparent" />
        </section>
      </main>
    </AppShell>
  )
}
