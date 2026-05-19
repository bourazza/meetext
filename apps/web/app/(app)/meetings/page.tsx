import React from 'react'
import { EmptyState } from '@/components/shared/EmptyState'
import { Video } from 'lucide-react'

export default function MeetingsPage() {
  return (
    <div className="mx-auto max-w-5xl px-8 py-12">
      <h1 className="mb-8 text-2xl font-semibold tracking-tight text-zinc-950">Meetings</h1>
      <EmptyState
        icon={<Video className="h-6 w-6" />}
        title="Meetings library coming soon"
        description="Soon you'll be able to browse, search, and manage all your past meeting recordings and AI-generated transcripts right here."
      />
    </div>
  )
}
