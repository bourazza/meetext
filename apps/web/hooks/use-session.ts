'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { getCurrentUser } from '@/services/auth'
import { getWorkspaces } from '@/services/workspace'
import { useAuthStore } from '@/store/auth'

export function useSession() {
  const router = useRouter()
  const { user, workspace, setUser, setWorkspace, clear } = useAuthStore()
  const [loading, setLoading] = useState(!user || !workspace)

  useEffect(() => {
    let active = true
    if (user && workspace) {
      setLoading(false)
      return
    }

    setLoading(true)
    Promise.all([getCurrentUser(), getWorkspaces()])
      .then(([nextUser, workspaces]) => {
        if (!active) return
        setUser(nextUser)
        setWorkspace(workspaces[0] ?? null)
      })
      .catch(() => {
        if (!active) return
        clear()
        router.replace('/login')
      })
      .finally(() => {
        if (active) setLoading(false)
      })

    return () => {
      active = false
    }
  }, [clear, router, setUser, setWorkspace, user, workspace])

  return { user, workspace, loading }
}
