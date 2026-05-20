'use client'

import React, { useEffect, useState } from 'react'
import { Sidebar } from '@/components/layout/Sidebar'
import { useAuthStore } from '@/store/auth'
import { getCurrentUser } from '@/services/auth'

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const { user, setUser } = useAuthStore()
  const [loading, setLoading] = useState(!user)

  useEffect(() => {
    if (!user) {
      getCurrentUser()
        .then((fetchedUser) => {
          setUser(fetchedUser)
        })
        .catch(() => {
          // If fetching fails but they have a session cookie, they'll eventually get redirected or it'll retry.
        })
        .finally(() => {
          setLoading(false)
        })
    } else {
      setLoading(false)
    }
  }, [user, setUser])

  if (loading) {
    return (
      <div className="flex h-screen w-full items-center justify-center bg-white">
        <div className="h-8 w-8 animate-spin rounded-full border-2 border-zinc-200 border-t-zinc-950" />
      </div>
    )
  }

  return (
    <div className="flex flex-col lg:flex-row h-screen w-full bg-zinc-50/40 text-zinc-950">
      <Sidebar />
      <main className="flex-1 overflow-y-auto bg-white lg:rounded-l-2xl border-l border-zinc-200/60 shadow-sm relative">
        {children}
      </main>
    </div>
  )
}
