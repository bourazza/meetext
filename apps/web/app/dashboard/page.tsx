'use client'

import { useQuery } from '@tanstack/react-query'
import { useAuthStore } from '@/store/auth'
import { getMeetings } from '@/services/meetings'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { getAccessToken } from '@/lib/api'
import Link from 'next/link'

export default function DashboardPage() {
  const router = useRouter()
  const { user, workspace } = useAuthStore()

  useEffect(() => {
    if (!getAccessToken()) router.push('/login')
  }, [router])

  const { data: meetings, isLoading } = useQuery({
    queryKey: ['meetings', workspace?.id],
    queryFn: () => getMeetings(workspace!.id),
    enabled: !!workspace?.id,
  })

  if (!user || !workspace) return null

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <span className="text-xl font-bold text-primary">Meetext</span>
          <span className="text-muted-foreground text-sm">/ {workspace.name}</span>
        </div>
        <div className="flex items-center gap-4">
          <span className="text-sm text-muted-foreground">{user.full_name}</span>
          <button
            onClick={() => { localStorage.clear(); router.push('/login') }}
            className="text-sm text-muted-foreground hover:text-foreground transition"
          >
            Sign out
          </button>
        </div>
      </header>

      <main className="max-w-5xl mx-auto px-6 py-10">
        {/* Welcome */}
        <div className="mb-10">
          <h1 className="text-2xl font-semibold">Welcome back, {user.full_name.split(' ')[0]} 👋</h1>
          <p className="text-muted-foreground mt-1 text-sm">Here&apos;s what&apos;s happening in your workspace.</p>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-3 gap-4 mb-10">
          {[
            { label: 'Total Meetings', value: meetings?.length ?? 0 },
            { label: 'Processing', value: meetings?.filter(m => m.status === 'processing').length ?? 0 },
            { label: 'Completed', value: meetings?.filter(m => m.status === 'completed').length ?? 0 },
          ].map((stat) => (
            <div key={stat.label} className="border rounded-xl p-5 bg-card">
              <p className="text-sm text-muted-foreground">{stat.label}</p>
              <p className="text-3xl font-bold mt-1">{stat.value}</p>
            </div>
          ))}
        </div>

        {/* Recent Meetings */}
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold">Recent Meetings</h2>
          <Link href="/meetings" className="text-sm text-primary hover:underline">View all</Link>
        </div>

        {isLoading ? (
          <div className="text-sm text-muted-foreground">Loading meetings...</div>
        ) : !meetings?.length ? (
          <div className="border rounded-xl p-10 text-center text-muted-foreground text-sm">
            No meetings yet.{' '}
            <Link href="/meetings" className="text-primary hover:underline">Upload your first meeting</Link>
          </div>
        ) : (
          <div className="space-y-3">
            {meetings.slice(0, 5).map((m) => (
              <div key={m.id} className="border rounded-xl px-5 py-4 bg-card flex items-center justify-between">
                <div>
                  <p className="font-medium text-sm">{m.title}</p>
                  <p className="text-xs text-muted-foreground mt-0.5">
                    {m.upload_type.toUpperCase()} · {new Date(m.created_at).toLocaleDateString()}
                  </p>
                </div>
                <span className={`text-xs px-2 py-1 rounded-full font-medium ${
                  m.status === 'completed' ? 'bg-green-100 text-green-700' :
                  m.status === 'processing' ? 'bg-yellow-100 text-yellow-700' :
                  m.status === 'failed' ? 'bg-red-100 text-red-700' :
                  'bg-muted text-muted-foreground'
                }`}>
                  {m.status}
                </span>
              </div>
            ))}
          </div>
        )}
      </main>
    </div>
  )
}
