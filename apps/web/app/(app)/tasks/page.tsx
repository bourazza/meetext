import React from 'react'
import { EmptyState } from '@/components/shared/EmptyState'
import { CheckSquare } from 'lucide-react'

export default function TasksPage() {
  return (
    <div className="mx-auto max-w-5xl px-8 py-12">
      <h1 className="mb-8 text-2xl font-semibold tracking-tight text-zinc-950">Tasks</h1>
      <EmptyState
        icon={<CheckSquare className="h-6 w-6" />}
        title="Task extraction coming soon"
        description="Meetext will automatically identify and extract actionable tasks from your meetings, assigned and prioritized."
      />
    </div>
  )
}
