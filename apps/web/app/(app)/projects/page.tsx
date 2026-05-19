import React from 'react'
import { EmptyState } from '@/components/shared/EmptyState'
import { FolderKanban } from 'lucide-react'

export default function ProjectsPage() {
  return (
    <div className="mx-auto max-w-5xl px-8 py-12">
      <h1 className="mb-8 text-2xl font-semibold tracking-tight text-zinc-950">Projects</h1>
      <EmptyState
        icon={<FolderKanban className="h-6 w-6" />}
        title="Project organization coming soon"
        description="Organize your meetings, tasks, and documents by project. Let AI automatically categorize your intelligence."
      />
    </div>
  )
}
