'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import {
  BriefcaseBusiness,
  CheckSquare,
  FileText,
  FolderKanban,
  LayoutDashboard,
  LogOut,
  Menu,
  Settings,
  Users,
  Video,
  X,
} from 'lucide-react'
import { useState } from 'react'
import { cn } from '@/lib/utils'
import { logout } from '@/services/auth'
import { useSession } from '@/hooks/use-session'

const navItems = [
  { href: '/dashboard', label: 'Dashboard', icon: LayoutDashboard },
  { href: '/meetings', label: 'Meetings', icon: Video },
  { href: '/projects', label: 'Projects', icon: FolderKanban },
  { href: '/tasks', label: 'Tasks', icon: CheckSquare },
  { href: '/documents', label: 'Documents', icon: FileText },
  { href: '/clients', label: 'Clients', icon: BriefcaseBusiness },
  { href: '/settings', label: 'Settings', icon: Settings },
]

export function AppShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const router = useRouter()
  const { user, workspace } = useSession()
  const [open, setOpen] = useState(false)

  const signOut = async () => {
    await logout()
    router.replace('/login')
  }

  return (
    <div className="min-h-screen bg-[#f7f8fb] text-zinc-950 lg:grid lg:grid-cols-[260px_1fr]">
      <aside
        className={cn(
          'fixed inset-y-0 left-0 z-40 w-[260px] border-r border-white/10 bg-[#0b1020] text-white transition-transform duration-300 lg:sticky lg:top-0 lg:translate-x-0',
          open ? 'translate-x-0' : '-translate-x-full'
        )}
      >
        <div className="flex h-full flex-col">
          <div className="flex h-20 items-center justify-between px-5">
            <Link href="/dashboard" className="flex items-center gap-3" onClick={() => setOpen(false)}>
              <span className="grid h-9 w-9 place-items-center rounded-lg bg-white text-sm font-bold text-[#0b1020] shadow-sm">
                M
              </span>
              <span>
                <span className="block text-sm font-semibold tracking-normal">Meetext</span>
                <span className="block text-xs text-white/45">Meeting intelligence</span>
              </span>
            </Link>
            <button
              type="button"
              onClick={() => setOpen(false)}
              className="grid h-9 w-9 place-items-center rounded-md text-white/60 hover:bg-white/10 lg:hidden"
              aria-label="Close navigation"
            >
              <X className="h-4 w-4" />
            </button>
          </div>

          <nav className="space-y-1 px-3">
            {navItems.map((item) => {
              const active = pathname === item.href
              const Icon = item.icon
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  onClick={() => setOpen(false)}
                  className={cn(
                    'group flex h-10 items-center gap-3 rounded-md px-3 text-sm text-white/64 transition',
                    active && 'bg-white text-[#0b1020] shadow-sm',
                    !active && 'hover:bg-white/[0.08] hover:text-white'
                  )}
                >
                  <Icon className="h-4 w-4" />
                  <span>{item.label}</span>
                </Link>
              )
            })}
          </nav>

          <div className="mt-auto border-t border-white/10 p-4">
            <div className="mb-3 rounded-lg border border-white/10 bg-white/[0.04] p-3">
              <p className="truncate text-sm font-medium">{workspace?.name ?? 'Workspace'}</p>
              <p className="mt-1 truncate text-xs text-white/45">{user?.email ?? 'Restoring session'}</p>
            </div>
            <button
              type="button"
              onClick={signOut}
              className="flex h-10 w-full items-center justify-center gap-2 rounded-md border border-white/10 text-sm text-white/70 transition hover:bg-white/10 hover:text-white"
            >
              <LogOut className="h-4 w-4" />
              Sign out
            </button>
          </div>
        </div>
      </aside>

      {open && <button className="fixed inset-0 z-30 bg-black/30 lg:hidden" onClick={() => setOpen(false)} aria-label="Close navigation overlay" />}

      <div className="min-w-0">
        <header className="sticky top-0 z-20 flex h-16 items-center justify-between border-b border-zinc-200/80 bg-[#f7f8fb]/90 px-4 backdrop-blur-xl sm:px-6 lg:hidden">
          <button
            type="button"
            onClick={() => setOpen(true)}
            className="grid h-10 w-10 place-items-center rounded-md border border-zinc-200 bg-white text-zinc-700 shadow-sm"
            aria-label="Open navigation"
          >
            <Menu className="h-4 w-4" />
          </button>
          <span className="text-sm font-semibold">Meetext</span>
          <span className="h-10 w-10" />
        </header>
        {children}
      </div>
    </div>
  )
}
