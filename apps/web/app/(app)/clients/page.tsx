import React from 'react'
import { EmptyState } from '@/components/shared/EmptyState'
import { Users } from 'lucide-react'

export default function ClientsPage() {
  return (
    <div className="mx-auto max-w-5xl px-8 py-12">
      <h1 className="mb-8 text-2xl font-semibold tracking-tight text-zinc-950">Clients</h1>
      <EmptyState
        icon={<Users className="h-6 w-6" />}
        title="Client management coming soon"
        description="Keep track of client preferences, key stakeholders, and project history across all your client interactions."
      />
    </div>
  )
}
