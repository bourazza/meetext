'use client'

import { useRef, useState } from 'react'
import { AudioLines, CheckCircle2, FileText, Link2, Loader2, UploadCloud, Video, X } from 'lucide-react'
import { toast } from 'sonner'
import { cn } from '@/lib/utils'
import { uploadMeeting } from '@/services/meetings'

type UploadPhase = 'idle' | 'uploading' | 'processing' | 'complete'

const supportedTypes = [
  { label: 'Audio', icon: AudioLines },
  { label: 'Video', icon: Video },
  { label: 'PDF', icon: FileText },
  { label: 'Meeting URL', icon: Link2 },
]

export function UploadMeetingCard({ workspaceId }: { workspaceId: string }) {
  const inputRef = useRef<HTMLInputElement>(null)
  const [dragging, setDragging] = useState(false)
  const [file, setFile] = useState<File | null>(null)
  const [meetingURL, setMeetingURL] = useState('')
  const [progress, setProgress] = useState(0)
  const [phase, setPhase] = useState<UploadPhase>('idle')

  const disabled = phase === 'uploading' || phase === 'processing'

  const chooseFile = (nextFile: File | null) => {
    if (!nextFile || disabled) return
    setFile(nextFile)
    setMeetingURL('')
    setProgress(0)
    setPhase('idle')
  }

  const startUpload = async () => {
    if (!file && !meetingURL.trim()) {
      toast.error('Add a file or paste a meeting URL first.')
      return
    }

    if (meetingURL.trim() && !file) {
      setPhase('processing')
      setProgress(100)
      window.setTimeout(() => {
        setPhase('complete')
        toast.success('Meeting URL captured. Link processing will be available soon.')
      }, 1100)
      return
    }

    if (!file) return

    try {
      setPhase('uploading')
      setProgress(8)
      await uploadMeeting({
        workspaceId,
        file,
        title: file.name.replace(/\.[^/.]+$/, ''),
        onProgress: setProgress,
      })
      setPhase('processing')
      setProgress(100)
      window.setTimeout(() => {
        setPhase('complete')
        toast.success('Meeting uploaded. Meetext is preparing the documentation workspace.')
      }, 1300)
    } catch (err: any) {
      setPhase('idle')
      toast.error(err?.response?.data?.error?.message ?? 'Upload failed. Try another file.')
    }
  }

  const reset = () => {
    setFile(null)
    setMeetingURL('')
    setProgress(0)
    setPhase('idle')
    if (inputRef.current) inputRef.current.value = ''
  }

  return (
    <section className="mx-auto w-full max-w-5xl">
      <div
        className={cn(
          'relative overflow-hidden rounded-2xl border border-zinc-200 bg-white shadow-[0_24px_80px_rgba(15,23,42,0.08)] transition duration-300',
          dragging && 'border-zinc-400 shadow-[0_28px_90px_rgba(15,23,42,0.12)]'
        )}
      >
        <div className="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-zinc-950 via-violet-600 to-cyan-500" />

        <div className="grid gap-0 lg:grid-cols-[1fr_340px]">
          <div className="p-5 sm:p-8">
            <div
              role="button"
              tabIndex={0}
              onClick={() => inputRef.current?.click()}
              onKeyDown={(event) => {
                if (event.key === 'Enter' || event.key === ' ') inputRef.current?.click()
              }}
              onDragOver={(event) => {
                event.preventDefault()
                setDragging(true)
              }}
              onDragLeave={() => setDragging(false)}
              onDrop={(event) => {
                event.preventDefault()
                setDragging(false)
                chooseFile(event.dataTransfer.files?.[0] ?? null)
              }}
              className={cn(
                'group flex min-h-[360px] cursor-pointer flex-col items-center justify-center rounded-xl border border-dashed border-zinc-300 bg-[#fbfbfd] px-5 py-10 text-center outline-none transition duration-300',
                'hover:border-zinc-500 hover:bg-white focus:ring-4 focus:ring-zinc-950/5',
                dragging && 'border-zinc-600 bg-white'
              )}
            >
              <input
                ref={inputRef}
                type="file"
                accept="audio/mpeg,audio/wav,video/mp4,application/pdf"
                className="hidden"
                onChange={(event) => chooseFile(event.target.files?.[0] ?? null)}
              />

              <div className="mb-6 grid h-16 w-16 place-items-center rounded-2xl bg-zinc-950 text-white shadow-lg shadow-zinc-950/15 transition duration-300 group-hover:-translate-y-1">
                {phase === 'complete' ? <CheckCircle2 className="h-7 w-7" /> : <UploadCloud className="h-7 w-7" />}
              </div>

              <h2 className="text-xl font-semibold tracking-normal text-zinc-950 sm:text-2xl">
                {file ? file.name : phase === 'complete' ? 'Meeting received' : 'Drop your meeting here'}
              </h2>
              <p className="mt-3 max-w-md text-sm leading-6 text-zinc-500">
                Upload audio, video, or a PDF. Meetext will turn it into clean documentation, decisions, and next steps.
              </p>

              <div className="mt-7 flex flex-wrap justify-center gap-2">
                {supportedTypes.map((item) => {
                  const Icon = item.icon
                  return (
                    <span key={item.label} className="inline-flex h-9 items-center gap-2 rounded-md border border-zinc-200 bg-white px-3 text-xs font-medium text-zinc-600 shadow-sm">
                      <Icon className="h-3.5 w-3.5" />
                      {item.label}
                    </span>
                  )
                })}
              </div>

              <button
                type="button"
                className="mt-8 inline-flex h-11 items-center justify-center rounded-md bg-zinc-950 px-5 text-sm font-medium text-white transition hover:bg-zinc-800"
              >
                Browse files
              </button>
            </div>
          </div>

          <aside className="border-t border-zinc-200 bg-zinc-50/70 p-5 sm:p-8 lg:border-l lg:border-t-0">
            <div className="space-y-6">
              <div>
                <p className="text-xs font-semibold uppercase tracking-[0.16em] text-zinc-400">Upload</p>
                <h3 className="mt-2 text-lg font-semibold text-zinc-950">Start with one meeting</h3>
                <p className="mt-2 text-sm leading-6 text-zinc-500">
                  Keep the MVP simple: upload the source, then let Meetext prepare the workspace.
                </p>
              </div>

              <label className="block space-y-2">
                <span className="text-sm font-medium text-zinc-700">Meeting URL</span>
                <div className="flex h-11 items-center gap-2 rounded-md border border-zinc-200 bg-white px-3 shadow-sm focus-within:ring-4 focus-within:ring-zinc-950/5">
                  <Link2 className="h-4 w-4 text-zinc-400" />
                  <input
                    value={meetingURL}
                    onChange={(event) => {
                      setMeetingURL(event.target.value)
                      if (event.target.value) setFile(null)
                    }}
                    disabled={disabled}
                    placeholder="Paste Zoom, Meet, or Loom link"
                    className="min-w-0 flex-1 bg-transparent text-sm outline-none placeholder:text-zinc-400"
                  />
                </div>
              </label>

              {(file || meetingURL || phase !== 'idle') && (
                <div className="rounded-xl border border-zinc-200 bg-white p-4 shadow-sm">
                  <div className="mb-4 flex items-start justify-between gap-3">
                    <div className="min-w-0">
                      <p className="truncate text-sm font-medium text-zinc-950">{file?.name ?? meetingURL}</p>
                      <p className="mt-1 text-xs text-zinc-500">{statusText(phase)}</p>
                    </div>
                    {!disabled && (
                      <button type="button" onClick={reset} className="grid h-7 w-7 place-items-center rounded-md text-zinc-400 hover:bg-zinc-100 hover:text-zinc-700" aria-label="Clear upload">
                        <X className="h-4 w-4" />
                      </button>
                    )}
                  </div>
                  <div className="h-2 overflow-hidden rounded-full bg-zinc-100">
                    <div className="h-full rounded-full bg-zinc-950 transition-all duration-500" style={{ width: `${phase === 'idle' ? 0 : progress}%` }} />
                  </div>
                </div>
              )}

              {phase === 'processing' && (
                <div className="rounded-xl border border-zinc-200 bg-white p-4 shadow-sm">
                  <div className="flex items-center gap-3 text-sm font-medium text-zinc-950">
                    <Loader2 className="h-4 w-4 animate-spin" />
                    AI is preparing your documentation
                  </div>
                  <div className="mt-4 space-y-3 text-xs text-zinc-500">
                    <ProcessingStep label="Transcribing meeting" active />
                    <ProcessingStep label="Extracting decisions" active />
                    <ProcessingStep label="Generating summary" />
                  </div>
                </div>
              )}

              <button
                type="button"
                onClick={startUpload}
                disabled={disabled}
                className="inline-flex h-11 w-full items-center justify-center gap-2 rounded-md bg-zinc-950 px-4 text-sm font-medium text-white shadow-sm transition hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {disabled && <Loader2 className="h-4 w-4 animate-spin" />}
                {phase === 'uploading' ? 'Uploading' : phase === 'processing' ? 'Processing' : 'Upload new meeting'}
              </button>
            </div>
          </aside>
        </div>
      </div>
    </section>
  )
}

function ProcessingStep({ label, active }: { label: string; active?: boolean }) {
  return (
    <div className="flex items-center gap-2">
      <span className={cn('h-1.5 w-1.5 rounded-full', active ? 'animate-pulse bg-emerald-500' : 'bg-zinc-300')} />
      {label}
    </div>
  )
}

function statusText(phase: UploadPhase) {
  if (phase === 'uploading') return 'Uploading securely'
  if (phase === 'processing') return 'Processing with AI'
  if (phase === 'complete') return 'Ready for documentation'
  return 'Ready to upload'
}
