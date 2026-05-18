'use client'

import { CheckCircle2, FileText, Mic, Users } from 'lucide-react'
import { AuthPanel, ContinueButton, StatusCard } from '@/components/auth/auth-ui'
import { useAuthStore } from '@/store/auth'

export default function OnboardingPage() {
  const user = useAuthStore((state) => state.user)

  return (
    <AuthPanel
      eyebrow="Workspace ready"
      title={`Welcome${user?.full_name ? `, ${user.full_name.split(' ')[0]}` : ''}`}
      subtitle="Your Meetext workspace is set up. Here is what is ready for your next client call."
    >
      <div className="space-y-3">
        {[
          { icon: Mic, title: 'Record and upload meetings', text: 'Capture calls, audio, video, and client notes.' },
          { icon: FileText, title: 'Generate polished documentation', text: 'Turn transcripts into summaries, decisions, and reports.' },
          { icon: Users, title: 'Keep projects organized', text: 'Attach outcomes to workspaces, clients, and projects.' },
        ].map((item) => (
          <div key={item.title} className="flex gap-3 rounded-md border border-zinc-200 bg-zinc-50 p-3">
            <item.icon className="mt-0.5 h-4 w-4 text-zinc-700" />
            <div>
              <div className="text-sm font-medium text-zinc-950">{item.title}</div>
              <p className="mt-1 text-sm leading-5 text-zinc-600">{item.text}</p>
            </div>
          </div>
        ))}
      </div>
      <div className="mt-5">
        <StatusCard tone="neutral" title="One small thing">
          We sent a verification link to your email. You can start now and verify when convenient.
        </StatusCard>
      </div>
      <div className="mt-5">
        <ContinueButton href="/dashboard">
          <CheckCircle2 className="h-4 w-4" />
          Go to dashboard
        </ContinueButton>
      </div>
    </AuthPanel>
  )
}
