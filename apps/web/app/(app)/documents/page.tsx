import React from 'react'
import { EmptyState } from '@/components/shared/EmptyState'
import { FileText } from 'lucide-react'

export default function DocumentsPage() {
  return (
    <div className="mx-auto max-w-5xl px-8 py-12">
      <h1 className="mb-8 text-2xl font-semibold tracking-tight text-zinc-950">Documents</h1>
      <EmptyState
        icon={<FileText className="h-6 w-6" />}
        title="AI documentation coming soon"
        description="Review, edit, and export automatically generated technical specs, PRDs, and summaries from your meetings."
      />
    </div>
  )
}
