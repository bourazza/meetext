'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuthStore } from '@/store/auth'
import { getMeetings, uploadMeeting, deleteMeeting } from '@/services/meetings'
import { useRef, useState } from 'react'
import { toast } from 'sonner'
import Link from 'next/link'

export default function MeetingsPage() {
  const { workspace } = useAuthStore()
  const qc = useQueryClient()
  const fileRef = useRef<HTMLInputElement>(null)
  const [title, setTitle] = useState('')

  const { data: meetings, isLoading } = useQuery({
    queryKey: ['meetings', workspace?.id],
    queryFn: () => getMeetings(workspace!.id),
    enabled: !!workspace?.id,
  })

  const upload = useMutation({
    mutationFn: (file: File) =>
      uploadMeeting({ workspaceId: workspace!.id, file, title: title || undefined }),
    onSuccess: () => {
      toast.success('Meeting uploaded successfully')
      setTitle('')
      if (fileRef.current) fileRef.current.value = ''
      qc.invalidateQueries({ queryKey: ['meetings', workspace?.id] })
    },
    onError: (err: any) => {
      toast.error(err?.response?.data?.error?.message ?? 'Upload failed')
    },
  })

  const remove = useMutation({
    mutationFn: (meetingId: string) => deleteMeeting(workspace!.id, meetingId),
    onSuccess: () => {
      toast.success('Meeting deleted')
      qc.invalidateQueries({ queryKey: ['meetings', workspace?.id] })
    },
  })

  const handleUpload = () => {
    const file = fileRef.current?.files?.[0]
    if (!file) { toast.error('Select a file first'); return }
    upload.mutate(file)
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b px-6 py-4 flex items-center gap-4">
        <Link href="/dashboard" className="text-muted-foreground text-sm hover:text-foreground">← Dashboard</Link>
        <span className="text-lg font-semibold">Meetings</span>
      </header>

      <main className="max-w-4xl mx-auto px-6 py-10 space-y-10">

        {/* Upload */}
        <div className="border rounded-xl p-6 bg-card space-y-4">
          <h2 className="font-semibold">Upload a Meeting</h2>
          <div className="flex flex-col gap-3">
            <input
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Title (optional — defaults to filename)"
              className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
            />
            <input
              ref={fileRef}
              type="file"
              accept="audio/mpeg,audio/wav,video/mp4,application/pdf"
              className="text-sm text-muted-foreground file:mr-3 file:rounded-md file:border-0 file:bg-primary file:text-primary-foreground file:px-3 file:py-1.5 file:text-xs file:font-medium"
            />
            <button
              onClick={handleUpload}
              disabled={upload.isPending}
              className="self-start bg-primary text-primary-foreground rounded-md px-4 py-2 text-sm font-medium hover:opacity-90 disabled:opacity-50 transition"
            >
              {upload.isPending ? 'Uploading...' : 'Upload'}
            </button>
          </div>
          <p className="text-xs text-muted-foreground">Supported: MP3, WAV, MP4, PDF · Max 500 MB</p>
        </div>

        {/* List */}
        <div>
          <h2 className="font-semibold mb-4">All Meetings</h2>
          {isLoading ? (
            <p className="text-sm text-muted-foreground">Loading...</p>
          ) : !meetings?.length ? (
            <p className="text-sm text-muted-foreground">No meetings yet.</p>
          ) : (
            <div className="space-y-3">
              {meetings.map((m) => (
                <div key={m.id} className="border rounded-xl px-5 py-4 bg-card flex items-center justify-between">
                  <div>
                    <p className="font-medium text-sm">{m.title}</p>
                    <p className="text-xs text-muted-foreground mt-0.5">
                      {m.upload_type.toUpperCase()} · {new Date(m.created_at).toLocaleDateString()}
                    </p>
                    {m.ai_summary && (
                      <p className="text-xs text-muted-foreground mt-1 line-clamp-1">{m.ai_summary}</p>
                    )}
                  </div>
                  <div className="flex items-center gap-3">
                    <span className={`text-xs px-2 py-1 rounded-full font-medium ${
                      m.status === 'completed' ? 'bg-green-100 text-green-700' :
                      m.status === 'processing' ? 'bg-yellow-100 text-yellow-700' :
                      m.status === 'failed' ? 'bg-red-100 text-red-700' :
                      'bg-muted text-muted-foreground'
                    }`}>
                      {m.status}
                    </span>
                    <button
                      onClick={() => remove.mutate(m.id)}
                      disabled={remove.isPending}
                      className="text-xs text-destructive hover:underline disabled:opacity-50"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  )
}
