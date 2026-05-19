'use client'

import React from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { motion } from 'framer-motion'
import {
  LayoutDashboard,
  Video,
  FolderKanban,
  CheckSquare,
  FileText,
  Users,
  Settings,
  LogOut,
  UploadCloud
} from 'lucide-react'
import { useAuthStore } from '@/store/auth'
import { logout } from '@/services/auth'

const navigation = [
  { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
  { name: 'Meetings', href: '/meetings', icon: Video },
  { name: 'Projects', href: '/projects', icon: FolderKanban },
  { name: 'Tasks', href: '/tasks', icon: CheckSquare },
  { name: 'Documents', href: '/documents', icon: FileText },
  { name: 'Clients', href: '/clients', icon: Users },
  { name: 'Settings', href: '/settings', icon: Settings },
]

export function Sidebar() {
  const pathname = usePathname()
  const { clear } = useAuthStore()

  const handleLogout = async () => {
    try {
      await logout()
    } finally {
      clear()
      window.location.href = '/login'
    }
  }

  return (
    <div className="flex w-64 flex-col border-r border-zinc-200 bg-zinc-50/50 px-4 py-6 text-sm font-medium text-zinc-600">
      <div className="mb-8 flex items-center px-3">
        <div className="flex h-8 w-8 items-center justify-center rounded bg-zinc-900 text-white">
          <UploadCloud className="h-5 w-5" />
        </div>
        <span className="ml-3 text-lg font-semibold tracking-tight text-zinc-950">Meetext</span>
      </div>

      <nav className="flex flex-1 flex-col space-y-1">
        {navigation.map((item) => {
          const isActive = pathname === item.href
          return (
            <Link
              key={item.name}
              href={item.href}
              className={`relative flex items-center gap-3 rounded-md px-3 py-2 transition-colors ${
                isActive ? 'text-zinc-950' : 'hover:bg-zinc-100 hover:text-zinc-900'
              }`}
            >
              {isActive && (
                <motion.div
                  layoutId="sidebar-active"
                  className="absolute inset-0 rounded-md bg-white shadow-sm ring-1 ring-zinc-200/50"
                  initial={false}
                  transition={{ type: 'spring', stiffness: 300, damping: 30 }}
                />
              )}
              <item.icon className="relative z-10 h-4 w-4" />
              <span className="relative z-10">{item.name}</span>
            </Link>
          )
        })}
      </nav>

      <div className="mt-auto">
        <button
          onClick={handleLogout}
          className="flex w-full items-center gap-3 rounded-md px-3 py-2 text-zinc-500 transition-colors hover:bg-zinc-100 hover:text-zinc-900"
        >
          <LogOut className="h-4 w-4" />
          <span>Sign out</span>
        </button>
      </div>
    </div>
  )
}
