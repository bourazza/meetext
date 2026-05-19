import React from 'react'
import { EmptyState } from '@/components/shared/EmptyState'
import { Settings as SettingsIcon } from 'lucide-react'

export default function SettingsPage() {
  return (
    <div className="mx-auto max-w-5xl px-8 py-12">
      <h1 className="mb-8 text-2xl font-semibold tracking-tight text-zinc-950">Settings</h1>
      <EmptyState
        icon={<SettingsIcon className="h-6 w-6" />}
        title="Workspace settings coming soon"
        description="Configure your workspace preferences, manage billing, and connect external integrations like Notion and Jira."
      />
    </div>
  )
}
