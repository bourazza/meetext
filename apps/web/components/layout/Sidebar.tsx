'use client'

import React, { useState } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { motion, AnimatePresence } from 'framer-motion'
import {
  LayoutDashboard,
  Video,
  FolderKanban,
  CheckSquare,
  FileText,
  Users,
  Settings,
  LogOut,
  Sparkles,
  Menu,
  X,
  Compass
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
  const { user, clear } = useAuthStore()
  const [mounted, setMounted] = useState(false)
  const [mobileOpen, setMobileOpen] = useState(false)

  React.useEffect(() => {
    setMounted(true)
  }, [])

  const handleLogout = async () => {
    try {
      await logout()
    } finally {
      clear()
      window.location.href = '/login'
    }
  }

  const displayName = user?.full_name || 'Zaki Bourazza'
  const displayEmail = user?.email || 'zaki@meetext.ai'
  const displayPlan = user?.plan || 'pro'
  const initial = displayName.charAt(0).toUpperCase()

  const userInitial = mounted ? initial : 'Z'
  const userName = mounted ? displayName : 'Zaki Bourazza'
  const userEmail = mounted ? displayEmail : 'zaki@meetext.ai'

  const sidebarContent = (
    <div className="flex h-full w-full flex-col bg-zinc-950 px-4 py-6 text-sm font-medium text-zinc-400 border-r border-zinc-900/50">
      {/* Brand logo & mobile close */}
      <div className="mb-8 flex items-center justify-between px-3">
        <Link href="/dashboard" className="flex items-center gap-3">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-indigo-600 text-white shadow-lg shadow-indigo-600/30">
            <Compass className="h-4 w-4 animate-pulse" />
          </div>
          <span className="text-lg font-bold tracking-tight text-white">Meetext</span>
          <span className="rounded bg-indigo-950 px-1.5 py-0.5 text-[10px] font-semibold text-indigo-400 border border-indigo-900">
            MVP
          </span>
        </Link>
        <button
          onClick={() => setMobileOpen(false)}
          className="lg:hidden text-zinc-400 hover:text-white p-1"
        >
          <X className="h-5 w-5" />
        </button>
      </div>

      {/* Nav Link items */}
      <nav className="flex-1 space-y-1 px-1">
        {navigation.map((item) => {
          const isActive = pathname === item.href || pathname.startsWith(item.href + '/')
          return (
            <Link
              key={item.name}
              href={item.href}
              onClick={() => setMobileOpen(false)}
              className={`relative flex items-center gap-3 rounded-lg px-3 py-2.5 transition-all duration-200 ${
                isActive 
                  ? 'text-white font-semibold' 
                  : 'hover:bg-zinc-900 hover:text-zinc-200'
              }`}
            >
              {isActive && (
                <motion.div
                  layoutId="sidebar-active-pill"
                  className="absolute inset-0 rounded-lg bg-zinc-900 border-l-2 border-indigo-500"
                  initial={false}
                  transition={{ type: 'spring', stiffness: 350, damping: 30 }}
                />
              )}
              <item.icon className={`relative z-10 h-4 w-4 ${isActive ? 'text-indigo-400' : 'text-zinc-400 group-hover:text-zinc-200'}`} />
              <span className="relative z-10">{item.name}</span>
            </Link>
          )
        })}
      </nav>

      {/* Sidebar bottom panel */}
      <div className="mt-8 flex flex-col space-y-4 border-t border-zinc-900 pt-6">
        {/* User Card */}
        <div className="flex items-center gap-3 px-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-indigo-900/40 text-indigo-300 font-semibold border border-indigo-500/20">
            {userInitial}
          </div>
          <div className="flex flex-col overflow-hidden">
            <span className="truncate font-semibold text-white leading-tight">{userName}</span>
            <span className="truncate text-xs text-zinc-500 leading-tight">{userEmail}</span>
          </div>
        </div>

        {/* Plan Upgrade Banner / Badge */}
        <div className="relative overflow-hidden rounded-xl bg-gradient-to-r from-zinc-900 to-zinc-900 border border-zinc-800 p-4 shadow-xl">
          <div className="relative z-10 flex flex-col gap-2">
            <div className="flex items-center justify-between">
              <span className="text-xs text-zinc-400">Current Workspace Plan</span>
              <span className="rounded-full bg-indigo-500/10 px-2 py-0.5 text-[10px] font-bold text-indigo-400 border border-indigo-500/20 capitalize">
                {displayPlan}
              </span>
            </div>
            <p className="text-xs text-zinc-500">Unlock Notion & Jira automated exports</p>
            <button className="flex w-full items-center justify-center gap-2 rounded-lg bg-indigo-600 px-3 py-2 text-xs font-semibold text-white transition-all hover:bg-indigo-500 shadow-md shadow-indigo-600/20 mt-1">
              <Sparkles className="h-3 w-3" />
              <span>Upgrade Plan</span>
            </button>
          </div>
        </div>

        <button
          onClick={handleLogout}
          className="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-zinc-500 hover:bg-zinc-900 hover:text-zinc-200 transition-colors"
        >
          <LogOut className="h-4 w-4" />
          <span>Sign out</span>
        </button>
      </div>
    </div>
  )

  return (
    <>
      {/* Mobile top navigation header */}
      <header className="lg:hidden flex items-center justify-between bg-zinc-950 border-b border-zinc-900 px-4 py-4 w-full text-white">
        <div className="flex items-center gap-3">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-indigo-600 text-white">
            <Compass className="h-4 w-4" />
          </div>
          <span className="text-lg font-bold tracking-tight">Meetext</span>
        </div>
        <button
          onClick={() => setMobileOpen(true)}
          className="text-zinc-400 hover:text-white p-1 rounded hover:bg-zinc-900"
        >
          <Menu className="h-6 w-6" />
        </button>
      </header>

      {/* Desktop sidebar */}
      <aside className="hidden lg:flex w-64 h-screen flex-shrink-0">
        {sidebarContent}
      </aside>

      {/* Mobile sidebar modal */}
      <AnimatePresence>
        {mobileOpen && (
          <div className="fixed inset-0 z-50 flex lg:hidden">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setMobileOpen(false)}
              className="fixed inset-0 bg-black/60 backdrop-blur-sm"
            />
            <motion.div
              initial={{ x: '-100%' }}
              animate={{ x: 0 }}
              exit={{ x: '-100%' }}
              transition={{ type: 'spring', damping: 25, stiffness: 200 }}
              className="relative flex w-64 max-w-xs h-full flex-col"
            >
              {sidebarContent}
            </motion.div>
          </div>
        )}
      </AnimatePresence>
    </>
  )
}

