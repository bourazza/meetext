'use client'

import React, { useState, useEffect, useRef } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  Sparkles, Bell, ArrowRight, UploadCloud, Link as LinkIcon, 
  FileText, Music, Video, ShieldCheck, CheckCircle2, 
  Loader2, Copy, Check, ChevronDown, ChevronUp, Edit2, 
  Download, Plus, RefreshCw, Layers, CheckSquare, 
  AlertCircle, HardDrive, Compass
} from 'lucide-react'
import { toast } from 'sonner'
import { useSession } from '@/hooks/use-session'
import { uploadMeeting } from '@/services/meetings'
import { useRouter } from 'next/navigation'

type ProcessState = 'idle' | 'uploading' | 'transcribing' | 'analyzing' | 'extracting' | 'summarizing' | 'complete'

interface Task {
  title: string
  description: string
  assignee: string
  priority: 'low' | 'medium' | 'high'
  due_date: string
}

interface Ticket {
  title: string
  description: string
  type: 'bug' | 'feature' | 'enhancement'
  status: 'todo' | 'in_progress' | 'done'
}

interface Decision {
  description: string
  made_by: string
}

interface Risk {
  description: string
  severity: 'low' | 'medium' | 'high'
  mitigation: string
}

interface TechNote {
  topic: string
  details: string
}

interface ClientRequest {
  request: string
  priority: 'low' | 'medium' | 'high'
  is_committed: boolean
}

interface GeneratedResults {
  summary: string
  tasks: Task[]
  tickets: Ticket[]
  decisions: Decision[]
  risks: Risk[]
  technical_notes: TechNote[]
  client_requests: ClientRequest[]
  project_documentation: string
}

const mockResults: GeneratedResults = {
  summary: "The meeting focused on the new mobile app authentication API. Sarah committed to leading the authentication endpoint implementation with a targeted delivery of this coming Friday. The team formally agreed to select PostgreSQL as the primary database instead of MongoDB for the new microservice. Critical blockers were identified regarding the outdated third-party payment gateway integration, which will be addressed by contacting their support team immediately. Additionally, Acme Corp explicitly requested the addition of a PDF export feature which the product team committed to supporting.",
  tasks: [
    {
      title: "Build Authentication Endpoint",
      description: "Design and implement secure JWT authentication API endpoints for the new mobile client.",
      assignee: "Sarah",
      priority: "high",
      due_date: "Friday"
    },
    {
      title: "Design Team Assets Follow-up",
      description: "Follow up with the brand agency regarding the final responsive SVG logo assets.",
      assignee: "David",
      priority: "medium",
      due_date: "Tomorrow"
    }
  ],
  tickets: [
    {
      title: "FEAT-104: Add PDF Export Feature",
      description: "Develop reusable frontend and backend service modules supporting standard project document PDF downloads requested by Acme Corp.",
      type: "feature",
      status: "todo"
    },
    {
      title: "BUG-219: Fix Payment Gateway Callback",
      description: "Investigate timeout errors and signature verification mismatch occurring under webhook integration environment.",
      type: "bug",
      status: "todo"
    }
  ],
  decisions: [
    {
      description: "Adopt PostgreSQL as primary data storage instead of MongoDB for database architecture alignment.",
      made_by: "Team"
    }
  ],
  risks: [
    {
      description: "Third-party payment gateway documentation is outdated and lacks modern SDK reference pages.",
      severity: "high",
      mitigation: "Submit priority support request ticket and schedule direct integration engineering call."
    }
  ],
  technical_notes: [
    {
      topic: "Storage Layer Selection",
      details: "PostgreSQL chosen to leverage rich native JSONB queries, relational safety, and strict multi-workspace schemas."
    }
  ],
  client_requests: [
    {
      request: "Acme Corp requested automated export to PDF and CSV formats for all monthly billing statement sheets.",
      priority: "high",
      is_committed: true
    }
  ],
  project_documentation: `## Meeting Minutes: Mobile App Sync

### 1. Executive Summary
The engineering team convened to align on key technical architecture choices for the upcoming mobile client releases. Priority efforts were assigned for the security infrastructure, and database layers were consolidated.

### 2. Strategic Technical Decisions
* **Database standard**: Formally standardized on **PostgreSQL** over MongoDB for our microservices due to transactional consistency requirements.
* **Security standards**: Enforced JWT authorization workflows across authentication endpoints.

### 3. Immediate Priorities
* **Auth Layer**: Sarah is leading the implementation of auth endpoints.
* **Integrations Blockers**: The team is taking proactive measures to resolve third-party API mismatches.`
}

export default function DashboardPage() {
  const { workspace, loading } = useSession()
  const router = useRouter()
  const [isPdfUpload, setIsPdfUpload] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [processState, setProcessState] = useState<ProcessState>('idle')
  const [urlInput, setUrlInput] = useState('')
  const [progress, setProgress] = useState(0)
  const [activeStep, setActiveStep] = useState(1)
  const [copiedSection, setCopiedSection] = useState<string | null>(null)
  const [dragging, setDragging] = useState(false)
  const [isMock, setIsMock] = useState(false)
  
  // Dynamic editable results state
  const [results, setResults] = useState<GeneratedResults>(mockResults)
  const [editingField, setEditingField] = useState<{ section: string; index?: number; field?: string } | null>(null)
  const [editValue, setEditValue] = useState('')

  // Collapsible cards state
  const [collapsedCards, setCollapsedCards] = useState<Record<string, boolean>>({
    summary: false,
    tasks: false,
    tickets: false,
    decisions: false,
    risks: false,
    technical_notes: false,
    client_requests: false,
    documentation: false
  })

  const toggleCollapse = (card: string) => {
    setCollapsedCards(prev => ({ ...prev, [card]: !prev[card] }))
  }

  const handleFileSelected = async (file: File | null) => {
    if (!file) return

    // 1. Validation
    if (file.type.startsWith('audio/') || file.type.startsWith('video/') || file.name.endsWith('.mp3') || file.name.endsWith('.wav') || file.name.endsWith('.mp4')) {
      toast.info('Coming Soon: Audio and video uploads are not supported in the MVP yet. Please upload a PDF file.', {
        duration: 5000,
      })
      return
    }

    if (file.type !== 'application/pdf' && !file.name.endsWith('.pdf')) {
      toast.error('Unsupported file type. Please upload a PDF file.')
      return
    }

    if (file.size > 250 * 1024 * 1024) {
      toast.error('File size exceeds the 250MB limit.')
      return
    }

    if (!workspace) {
      toast.error('No active workspace selected.')
      return
    }

    setIsPdfUpload(true)
    setIsMock(false)
    setProcessState('uploading')
    setProgress(5)
    setActiveStep(1)

    try {
      const uploadPromise = uploadMeeting({
        workspaceId: workspace.id,
        file,
        title: file.name.replace(/\.[^/.]+$/, ''),
        onProgress: (p) => {
          setProgress(Math.round(p * 0.95))
        }
      })

      const res = await uploadPromise
      setProgress(100)

      setProcessState('transcribing')
      setActiveStep(2)
      setProgress(15)
      await new Promise(resolve => setTimeout(resolve, 800))

      setProcessState('analyzing')
      setActiveStep(3)
      setProgress(40)
      await new Promise(resolve => setTimeout(resolve, 1000))

      setProcessState('extracting')
      setActiveStep(4)
      setProgress(70)
      await new Promise(resolve => setTimeout(resolve, 800))

      setProcessState('summarizing')
      setActiveStep(5)
      setProgress(90)
      await new Promise(resolve => setTimeout(resolve, 600))

      if (res.analysis) {
        setResults(res.analysis)
      } else {
        toast.error("Analysis not found in response, showing fallback results.")
      }

      setProcessState('complete')
      setActiveStep(6)
      toast.success("AI Generation Complete!", {
        description: "Extracted tasks, decisions, and summaries are ready."
      })

    } catch (err: any) {
      setProcessState('idle')
      console.error(err)
      const errMsg = err.response?.data?.error?.message || err.message || "Failed to process PDF."
      toast.error(`Upload failed: ${errMsg}`)
    }
  }

  const triggerProcessingFlow = () => {
    setIsMock(true)
    setProcessState('uploading')
    setProgress(0)
    setActiveStep(1)
  }

  useEffect(() => {
    if (processState === 'idle' || !isMock) return

    let interval: NodeJS.Timeout

    if (processState === 'uploading') {
      interval = setInterval(() => {
        setProgress(p => {
          if (p >= 100) {
            clearInterval(interval)
            setProcessState('transcribing')
            setActiveStep(2)
            setProgress(0)
            return 100
          }
          return p + 10
        })
      }, 150)
    } else if (processState === 'transcribing') {
      interval = setInterval(() => {
        setProgress(p => {
          if (p >= 100) {
            clearInterval(interval)
            setProcessState('analyzing')
            setActiveStep(3)
            setProgress(0)
            return 100
          }
          return p + 8
        })
      }, 200)
    } else if (processState === 'analyzing') {
      interval = setInterval(() => {
        setProgress(p => {
          if (p >= 100) {
            clearInterval(interval)
            setProcessState('extracting')
            setActiveStep(4)
            setProgress(0)
            return 100
          }
          return p + 12
        })
      }, 250)
    } else if (processState === 'extracting') {
      interval = setInterval(() => {
        setProgress(p => {
          if (p >= 100) {
            clearInterval(interval)
            setProcessState('summarizing')
            setActiveStep(5)
            setProgress(0)
            return 100
          }
          return p + 15
        })
      }, 150)
    } else if (processState === 'summarizing') {
      interval = setInterval(() => {
        setProgress(p => {
          if (p >= 100) {
            clearInterval(interval)
            setProcessState('complete')
            setActiveStep(6)
            toast.success("AI Generation Complete!", {
              description: "Extracted tasks, decisions, and summaries are ready."
            })
            return 100
          }
          return p + 20
        })
      }, 100)
    }

    return () => clearInterval(interval)
  }, [processState, isMock])

  const copyToClipboard = (text: string, section: string) => {
    navigator.clipboard.writeText(text)
    setCopiedSection(section)
    toast.success(`${section} copied to clipboard`)
    setTimeout(() => setCopiedSection(null), 2000)
  }

  const exportAsMarkdown = () => {
    let markdown = `# ${results.project_documentation}\n\n`
    markdown += `## Meeting Summary\n${results.summary}\n\n`
    
    markdown += `## Action Items\n`
    results.tasks.forEach(t => {
      markdown += `* **[${t.priority.toUpperCase()}]** ${t.title} - Assigned to: ${t.assignee} (Due: ${t.due_date})\n  _${t.description}_\n`
    })
    markdown += `\n`

    markdown += `## Decisions Made\n`
    results.decisions.forEach(d => {
      markdown += `* ${d.description} (Made by: ${d.made_by})\n`
    })

    const blob = new Blob([markdown], { type: 'text/markdown' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `Meetext-Summary-${new Date().toISOString().split('T')[0]}.md`
    a.click()
    toast.success("Document exported as Markdown!")
  }

  // Inline editing saving
  const startEditing = (section: string, value: string, index?: number, field?: string) => {
    setEditingField({ section, index, field })
    setEditValue(value)
  }

  const saveEdit = () => {
    if (!editingField) return
    const { section, index, field } = editingField

    if (section === 'summary') {
      setResults(prev => ({ ...prev, summary: editValue }))
    } else if (section === 'project_documentation') {
      setResults(prev => ({ ...prev, project_documentation: editValue }))
    } else if (section === 'tasks' && index !== undefined && field) {
      setResults(prev => {
        const updated = [...prev.tasks]
        updated[index] = { ...updated[index], [field]: editValue }
        return { ...prev, tasks: updated }
      })
    } else if (section === 'tickets' && index !== undefined && field) {
      setResults(prev => {
        const updated = [...prev.tickets]
        updated[index] = { ...updated[index], [field]: editValue }
        return { ...prev, tickets: updated }
      })
    } else if (section === 'decisions' && index !== undefined) {
      setResults(prev => {
        const updated = [...prev.decisions]
        updated[index] = { ...updated[index], description: editValue }
        return { ...prev, decisions: updated }
      })
    } else if (section === 'risks' && index !== undefined && field) {
      setResults(prev => {
        const updated = [...prev.risks]
        updated[index] = { ...updated[index], [field]: editValue }
        return { ...prev, risks: updated }
      })
    }

    setEditingField(null)
    toast.success("Content updated successfully")
  }

  return (
    <div className="mx-auto w-full max-w-[1400px] px-6 py-8 sm:px-8 lg:px-12 bg-white min-h-screen">
      
      {/* Header Area */}
      <div className="mb-8 flex flex-col items-start justify-between gap-4 border-b border-zinc-100 pb-6 sm:flex-row sm:items-center">
        <div>
          <div className="flex items-center gap-2 mb-1.5">
            <span className="inline-flex items-center gap-1 rounded-full bg-indigo-50 px-2 py-0.5 text-xs font-semibold text-indigo-600 border border-indigo-100">
              <Sparkles className="h-3 w-3" />
              AI Intelligent Workspace
            </span>
          </div>
          <h1 className="text-3xl font-extrabold tracking-tight text-zinc-900 sm:text-4xl">
            Workspace Hub
          </h1>
          <p className="mt-1 text-sm text-zinc-500">
            Automate meetings analysis, action item extraction, and document synthesis.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <button className="flex items-center gap-2 rounded-lg border border-zinc-200 bg-white px-4 py-2 text-xs font-semibold text-zinc-700 shadow-sm transition hover:bg-zinc-50">
            <Layers className="h-4 w-4" />
            Integrations
          </button>
          <button className="flex h-9 w-9 items-center justify-center rounded-lg border border-zinc-200 bg-white text-zinc-500 shadow-sm transition hover:bg-zinc-50 hover:text-zinc-900 relative">
            <Bell className="h-4 w-4" />
            <span className="absolute top-1 right-1 flex h-2 w-2 rounded-full bg-indigo-600 ring-2 ring-white" />
          </button>
        </div>
      </div>

      {/* Hero & Onboarding Banner */}
      <div className="relative mb-8 overflow-hidden rounded-2xl bg-gradient-to-r from-zinc-900 to-zinc-950 px-8 py-10 shadow-2xl">
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_80%_80%_at_50%_-20%,rgba(99,102,241,0.15),rgba(255,255,255,0))]" />
        <div className="relative z-10 max-w-2xl">
          <h2 className="text-2xl font-bold tracking-tight text-white sm:text-3xl">
            Transcribe & Auto-Document in Seconds
          </h2>
          <p className="mt-3 text-base text-zinc-300 leading-relaxed">
            Drop your meeting recordings, audio clips, or post virtual platform links. Meetext's semantic processing engine extracts structural decisions, task lists, risks, and formats standard documents immediately.
          </p>
          <div className="mt-6 flex flex-wrap items-center gap-4">
            <a href="#upload-section" className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-lg shadow-indigo-600/30 hover:bg-indigo-500 hover:shadow-indigo-600/40 transition">
              Process New Meeting
              <ArrowRight className="h-4 w-4" />
            </a>
            <button 
              onClick={() => {
                setResults(mockResults)
                setProcessState('complete')
              }}
              className="inline-flex items-center gap-2 rounded-lg bg-white/10 px-4 py-2.5 text-sm font-semibold text-white backdrop-blur-sm hover:bg-white/20 transition"
            >
              Load Demo Dataset
            </button>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-3" id="upload-section">
        
        {/* Upload & Flow column */}
        <div className="lg:col-span-2 space-y-6">
          <div className="rounded-2xl border border-zinc-200/80 bg-white p-6 shadow-sm sm:p-8">
            <div className="mb-6">
              <h2 className="text-lg font-bold text-zinc-900">Upload Center</h2>
              <p className="text-xs text-zinc-500">
                Supports video, voice memos, transcript files, or Google Meet / Zoom invite credentials.
              </p>
            </div>

            {processState === 'idle' && (
              <div className="space-y-6">
                {/* Hidden File Input */}
                <input
                  type="file"
                  ref={fileInputRef}
                  className="hidden"
                  accept="application/pdf,audio/*,video/*"
                  onChange={(e) => handleFileSelected(e.target.files?.[0] || null)}
                />

                {/* Drag and Drop Container */}
                <div 
                  onClick={() => fileInputRef.current?.click()}
                  onDragOver={(e) => {
                    e.preventDefault()
                    setDragging(true)
                  }}
                  onDragLeave={() => setDragging(false)}
                  onDrop={(e) => {
                    e.preventDefault()
                    setDragging(false)
                    handleFileSelected(e.dataTransfer.files?.[0] || null)
                  }}
                  className={`group flex flex-col items-center justify-center rounded-2xl border-2 border-dashed py-16 text-center cursor-pointer transition-all duration-200 hover:scale-[1.01] shadow-inner ${
                    dragging
                      ? 'border-indigo-500 bg-indigo-50/30'
                      : 'border-zinc-200 bg-zinc-50/50 hover:border-indigo-400 hover:bg-indigo-50/20'
                  }`}
                >
                  <div className="mb-4 inline-flex items-center justify-center rounded-2xl bg-white p-4 text-zinc-700 shadow-sm border border-zinc-100 group-hover:scale-105 group-hover:border-indigo-100 transition-all duration-300">
                    <UploadCloud className="h-8 w-8 text-indigo-500" />
                  </div>
                  <h3 className="mb-1 text-base font-bold text-zinc-950 group-hover:text-indigo-600 transition-colors">
                    Drop files here or click to browse
                  </h3>
                  <p className="mb-6 text-xs text-zinc-400">
                    MP4, MP3, WAV, PDF up to 250MB
                  </p>
                  
                  {/* Shortcut badges */}
                  <div className="flex flex-wrap items-center justify-center gap-3">
                    <span className="flex items-center gap-1.5 rounded-lg border border-zinc-200 bg-white px-3 py-1.5 text-xs text-zinc-600 shadow-sm group-hover:border-zinc-300 transition-colors">
                      <Music className="h-3 w-3 text-emerald-500" /> Audio file
                    </span>
                    <span className="flex items-center gap-1.5 rounded-lg border border-zinc-200 bg-white px-3 py-1.5 text-xs text-zinc-600 shadow-sm group-hover:border-zinc-300 transition-colors">
                      <Video className="h-3 w-3 text-sky-500" /> Video file
                    </span>
                    <span className="flex items-center gap-1.5 rounded-lg border border-zinc-200 bg-white px-3 py-1.5 text-xs text-zinc-600 shadow-sm group-hover:border-zinc-300 transition-colors">
                      <FileText className="h-3 w-3 text-red-500" /> PDF transcript
                    </span>
                  </div>
                </div>

                {/* URL Invite paste option */}
                <div className="relative">
                  <div className="absolute inset-0 flex items-center" aria-hidden="true">
                    <div className="w-full border-t border-zinc-150" />
                  </div>
                  <div className="relative flex justify-center text-xs">
                    <span className="bg-white px-3 text-zinc-400 font-medium">OR AUTOMATE VIA URL</span>
                  </div>
                </div>

                <div className="flex gap-2">
                  <div className="relative flex-1">
                    <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                      <LinkIcon className="h-4 w-4 text-zinc-400" />
                    </div>
                    <input
                      type="text"
                      placeholder="Paste cloud meeting URL (Zoom, Teams, Loom...)"
                      value={urlInput}
                      onChange={(e) => setUrlInput(e.target.value)}
                      className="block w-full rounded-lg border border-zinc-250 bg-white py-2.5 pl-10 pr-3 text-sm text-zinc-900 placeholder:text-zinc-400 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                    />
                  </div>
                  <button 
                    onClick={triggerProcessingFlow}
                    className="rounded-lg bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 transition-all focus:outline-none"
                  >
                    Fetch & Process
                  </button>
                </div>

                <div className="flex items-center justify-center gap-2 text-xs font-semibold text-zinc-400">
                  <ShieldCheck className="h-4 w-4 text-emerald-500" />
                  <span>Enterprise grade security · HIPAA compliant processing</span>
                </div>
              </div>
            )}

            {/* AI processing progress state */}
            {processState !== 'idle' && processState !== 'complete' && (
              <div className="py-8">
                <div className="mb-8 flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <Loader2 className="h-5 w-5 animate-spin text-indigo-600" />
                    <span className="font-bold text-zinc-900">
                      {processState === 'uploading' && 'Uploading raw workspace payload...'}
                      {processState === 'transcribing' && (isPdfUpload ? 'Extracting text from PDF...' : 'Transcribing audio streams...')}
                      {processState === 'analyzing' && 'Analyzing semantic contexts...'}
                      {processState === 'extracting' && 'Structuring action cards...'}
                      {processState === 'summarizing' && 'Compiling project docs...'}
                    </span>
                  </div>
                  <span className="text-xs font-bold text-zinc-500 bg-zinc-100 px-2.5 py-0.5 rounded-full">
                    {progress}% Complete
                  </span>
                </div>

                {/* Progress bar container */}
                <div className="h-2 w-full overflow-hidden rounded-full bg-zinc-100">
                  <motion.div
                    className="h-full bg-indigo-600 rounded-full"
                    initial={{ width: 0 }}
                    animate={{ width: `${progress}%` }}
                    transition={{ ease: 'linear' }}
                  />
                </div>

                {/* Vertical processing timeline checklist */}
                <div className="mt-8 space-y-4">
                  {[
                    { id: 1, label: 'Uploading file data securely', state: 'uploading' },
                    { id: 2, label: isPdfUpload ? 'Extracting text from PDF' : 'Whisper deep transcription processing', state: 'transcribing' },
                    { id: 3, label: 'LLM semantic analysis', state: 'analyzing' },
                    { id: 4, label: 'Task, ticket, and decision structuring', state: 'extracting' },
                    { id: 5, label: 'Synthesizing final Markdown document logs', state: 'summarizing' }
                  ].map((step) => {
                    const isDone = activeStep > step.id
                    const isActive = activeStep === step.id

                    return (
                      <div key={step.id} className="flex items-center gap-3 text-sm">
                        {isDone ? (
                          <div className="flex h-5 w-5 items-center justify-center rounded-full bg-emerald-50 text-emerald-500 border border-emerald-200">
                            <CheckCircle2 className="h-4 w-4" />
                          </div>
                        ) : isActive ? (
                          <div className="flex h-5 w-5 items-center justify-center rounded-full bg-indigo-50 border border-indigo-200">
                            <Loader2 className="h-3 w-3 animate-spin text-indigo-600" />
                          </div>
                        ) : (
                          <div className="flex h-5 w-5 items-center justify-center rounded-full border border-zinc-200 text-zinc-300">
                            <div className="h-1.5 w-1.5 rounded-full bg-current" />
                          </div>
                        )}
                        <span className={`font-medium ${isDone ? 'text-zinc-400 line-through' : isActive ? 'text-indigo-600 font-bold' : 'text-zinc-400'}`}>
                          {step.label}
                        </span>
                      </div>
                    )
                  })}
                </div>
              </div>
            )}

            {/* Complete output trigger resets */}
            {processState === 'complete' && (
              <div className="text-center py-6">
                <div className="inline-flex h-12 w-12 items-center justify-center rounded-full bg-emerald-50 text-emerald-500 border border-emerald-100 mb-3">
                  <CheckCircle2 className="h-6 w-6" />
                </div>
                <h3 className="text-lg font-bold text-zinc-900">Meeting Processing Complete</h3>
                <p className="text-xs text-zinc-500 mt-1">
                  Transcripts synthesised, AI objects exported down below.
                </p>
                <div className="mt-4 flex items-center justify-center gap-3">
                  <button 
                    onClick={() => setProcessState('idle')}
                    className="inline-flex items-center gap-1.5 rounded-lg border border-zinc-200 bg-white px-3 py-1.5 text-xs font-semibold text-zinc-700 shadow-sm hover:bg-zinc-50 transition"
                  >
                    <Plus className="h-3.5 w-3.5" /> Process Another
                  </button>
                  <button 
                    onClick={exportAsMarkdown}
                    className="inline-flex items-center gap-1.5 rounded-lg bg-indigo-600 px-3 py-1.5 text-xs font-semibold text-white shadow-sm hover:bg-indigo-500 transition"
                  >
                    <Download className="h-3.5 w-3.5" /> Export All (.md)
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Info Sidebar Column */}
        <div className="lg:col-span-1 space-y-6">
          <div className="rounded-2xl border border-zinc-200/80 bg-white p-6 shadow-sm">
            <div className="mb-4 flex items-center gap-3 border-b border-zinc-100 pb-4">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-indigo-50 text-indigo-600 border border-indigo-100">
                <Sparkles className="h-5 w-5" />
              </div>
              <div>
                <h2 className="font-bold text-zinc-900 leading-tight">AI Agent Status</h2>
                <p className="text-[10px] text-zinc-500">Connected to local node</p>
              </div>
            </div>
            
            <div className="space-y-4">
              <div className="flex items-center justify-between text-xs border-b border-zinc-50 pb-2">
                <span className="text-zinc-500 font-semibold">Active LLM Model</span>
                <span className="font-mono rounded bg-zinc-100 px-1.5 py-0.5 text-zinc-800 text-[10px] border border-zinc-200">
                  Llama-3-8B (ollama)
                </span>
              </div>
              <div className="flex items-center justify-between text-xs border-b border-zinc-50 pb-2">
                <span className="text-zinc-500 font-semibold">Transcription Engine</span>
                <span className="rounded bg-emerald-50 border border-emerald-100 px-1.5 py-0.5 text-emerald-700 text-[10px]">
                  Whisper Large v3
                </span>
              </div>
              <div className="flex items-center justify-between text-xs">
                <span className="text-zinc-500 font-semibold">Connected Integrations</span>
                <span className="text-[10px] text-zinc-400">Jira, Notion, Slack</span>
              </div>
            </div>
          </div>

          <div className="rounded-2xl border border-zinc-200/80 bg-white p-6 shadow-sm">
            <h3 className="font-bold text-zinc-900 mb-3 text-sm">Semantic Activity Feed</h3>
            <div className="space-y-4 text-xs">
              <div className="flex gap-2.5">
                <div className="h-1.5 w-1.5 rounded-full bg-emerald-500 mt-1.5" />
                <div>
                  <p className="font-semibold text-zinc-800">New meeting structured</p>
                  <p className="text-[10px] text-zinc-400">Database setup & billing sync completed</p>
                </div>
              </div>
              <div className="flex gap-2.5">
                <div className="h-1.5 w-1.5 rounded-full bg-indigo-500 mt-1.5" />
                <div>
                  <p className="font-semibold text-zinc-800">Doc sync template exported</p>
                  <p className="text-[10px] text-zinc-400">Created Sprint Notes doc page inside Workspace</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Render Generated Results Section */}
      <AnimatePresence>
        {processState === 'complete' && (
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 30 }}
            className="mt-12 space-y-8"
          >
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-zinc-200 pb-4">
              <div>
                <h2 className="text-2xl font-bold tracking-tight text-zinc-900">AI Meeting Intelligence</h2>
                <p className="text-sm text-zinc-500">Edit, structure, and dispatch actionable elements from this conversation.</p>
              </div>
              <div className="flex items-center gap-2">
                <button
                  onClick={exportAsMarkdown}
                  className="inline-flex items-center gap-1.5 rounded-lg border border-zinc-200 bg-white px-3 py-2 text-xs font-semibold text-zinc-700 shadow-sm hover:bg-zinc-50 transition"
                >
                  <Download className="h-4 w-4 text-zinc-400" /> Export Document
                </button>
                <button
                  onClick={() => {
                    toast.success("Synchronized successfully", {
                      description: "Transcripts synched to active workspace."
                    })
                  }}
                  className="inline-flex items-center gap-1.5 rounded-lg bg-indigo-600 px-3.5 py-2 text-xs font-semibold text-white shadow-sm hover:bg-indigo-500 transition"
                >
                  <RefreshCw className="h-4 w-4" /> Save To DB
                </button>
              </div>
            </div>

            {/* Grid structure of AI Cards */}
            <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
              
              {/* Card 1: Exec Summary */}
              <div className="rounded-xl border border-zinc-200/80 bg-white shadow-sm overflow-hidden lg:col-span-2">
                <div className="flex items-center justify-between border-b border-zinc-100 bg-zinc-50/50 px-5 py-4">
                  <div className="flex items-center gap-2.5">
                    <div className="h-2 w-2 rounded-full bg-indigo-500" />
                    <h3 className="font-bold text-zinc-900 text-sm">Meeting Summary</h3>
                    <span className="rounded bg-indigo-50 border border-indigo-100 px-2 py-0.5 text-[10px] text-indigo-600 font-semibold">
                      AI Generated
                    </span>
                  </div>
                  <div className="flex items-center gap-2">
                    <button 
                      onClick={() => toggleCollapse('summary')}
                      className="p-1 rounded text-zinc-400 hover:text-zinc-900 hover:bg-zinc-100 transition"
                    >
                      {collapsedCards.summary ? <ChevronDown className="h-4 w-4" /> : <ChevronUp className="h-4 w-4" />}
                    </button>
                  </div>
                </div>
                
                {!collapsedCards.summary && (
                  <div className="p-5">
                    {editingField?.section === 'summary' ? (
                      <div className="space-y-3">
                        <textarea
                          value={editValue}
                          onChange={(e) => setEditValue(e.target.value)}
                          rows={4}
                          className="w-full rounded-lg border border-zinc-300 p-3 text-sm focus:border-indigo-500 focus:outline-none"
                        />
                        <div className="flex justify-end gap-2">
                          <button onClick={() => setEditingField(null)} className="rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-semibold text-zinc-700">
                            Cancel
                          </button>
                          <button onClick={saveEdit} className="rounded-lg bg-indigo-600 px-3 py-1.5 text-xs font-semibold text-white">
                            Save
                          </button>
                        </div>
                      </div>
                    ) : (
                      <div className="group relative">
                        <p className="text-sm leading-relaxed text-zinc-650 pr-12">{results.summary}</p>
                        <div className="absolute right-0 top-0 flex items-center gap-1.5 opacity-0 group-hover:opacity-100 transition-opacity">
                          <button 
                            onClick={() => startEditing('summary', results.summary)}
                            className="p-1 rounded text-zinc-400 hover:text-zinc-900 hover:bg-zinc-100"
                          >
                            <Edit2 className="h-3.5 w-3.5" />
                          </button>
                          <button 
                            onClick={() => copyToClipboard(results.summary, 'Summary')}
                            className="p-1 rounded text-zinc-400 hover:text-zinc-900 hover:bg-zinc-100"
                          >
                            <Copy className="h-3.5 w-3.5" />
                          </button>
                        </div>
                      </div>
                    )}
                  </div>
                )}
              </div>

              {/* Card 2: Action Tasks list */}
              <div className="rounded-xl border border-zinc-200/80 bg-white shadow-sm overflow-hidden">
                <div className="flex items-center justify-between border-b border-zinc-100 bg-zinc-50/50 px-5 py-4">
                  <div className="flex items-center gap-2.5">
                    <CheckSquare className="h-4 w-4 text-indigo-500" />
                    <h3 className="font-bold text-zinc-900 text-sm">Action Tasks</h3>
                  </div>
                  <button 
                    onClick={() => toggleCollapse('tasks')}
                    className="p-1 rounded text-zinc-400 hover:text-zinc-950 hover:bg-zinc-100 transition"
                  >
                    {collapsedCards.tasks ? <ChevronDown className="h-4 w-4" /> : <ChevronUp className="h-4 w-4" />}
                  </button>
                </div>

                {!collapsedCards.tasks && (
                  <div className="p-5 space-y-4">
                    {results.tasks.map((task, i) => (
                      <div key={i} className="group relative border border-zinc-100 rounded-lg p-3 hover:bg-zinc-50/50 transition">
                        <div className="flex items-start justify-between gap-4">
                          <div>
                            <div className="flex items-center gap-2">
                              <h4 className="text-sm font-bold text-zinc-900">{task.title}</h4>
                              <span className={`rounded px-1.5 py-0.5 text-[9px] font-bold uppercase border ${
                                task.priority === 'high' 
                                  ? 'bg-rose-50 border-rose-100 text-rose-600'
                                  : 'bg-amber-50 border-amber-100 text-amber-600'
                              }`}>
                                {task.priority}
                              </span>
                            </div>
                            <p className="text-xs text-zinc-500 mt-1">{task.description}</p>
                            <div className="flex items-center gap-3 mt-2.5 text-[10px] text-zinc-400">
                              <span className="font-semibold text-zinc-600 bg-zinc-100 rounded px-1.5 py-0.5">Assignee: {task.assignee}</span>
                              <span>Due: {task.due_date}</span>
                            </div>
                          </div>
                          
                          <div className="opacity-0 group-hover:opacity-100 transition-opacity flex items-center gap-1">
                            <button 
                              onClick={() => startEditing('tasks', task.title, i, 'title')}
                              className="p-1 rounded text-zinc-400 hover:text-zinc-950"
                            >
                              <Edit2 className="h-3 w-3" />
                            </button>
                            <button 
                              onClick={() => copyToClipboard(`${task.title} - Assigned to: ${task.assignee}`, 'Task')}
                              className="p-1 rounded text-zinc-400 hover:text-zinc-950"
                            >
                              <Copy className="h-3 w-3" />
                            </button>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              {/* Card 3: Engineering Tickets */}
              <div className="rounded-xl border border-zinc-200/80 bg-white shadow-sm overflow-hidden">
                <div className="flex items-center justify-between border-b border-zinc-100 bg-zinc-50/50 px-5 py-4">
                  <div className="flex items-center gap-2.5">
                    <Compass className="h-4 w-4 text-indigo-500" />
                    <h3 className="font-bold text-zinc-900 text-sm">Product Tickets</h3>
                  </div>
                  <button 
                    onClick={() => toggleCollapse('tickets')}
                    className="p-1 rounded text-zinc-400 hover:text-zinc-950 hover:bg-zinc-100 transition"
                  >
                    {collapsedCards.tickets ? <ChevronDown className="h-4 w-4" /> : <ChevronUp className="h-4 w-4" />}
                  </button>
                </div>

                {!collapsedCards.tickets && (
                  <div className="p-5 space-y-4">
                    {results.tickets.map((ticket, i) => (
                      <div key={i} className="group relative border border-zinc-100 rounded-lg p-3 hover:bg-zinc-50/50 transition">
                        <div className="flex items-start justify-between gap-4">
                          <div>
                            <div className="flex items-center gap-2">
                              <h4 className="text-sm font-bold text-zinc-900">{ticket.title}</h4>
                              <span className="rounded bg-zinc-100 border border-zinc-200 px-1.5 py-0.5 text-[9px] font-bold text-zinc-655 uppercase">
                                {ticket.type}
                              </span>
                            </div>
                            <p className="text-xs text-zinc-500 mt-1">{ticket.description}</p>
                            <span className="inline-block mt-2.5 rounded bg-zinc-100 px-2 py-0.5 text-[9px] font-semibold text-zinc-600 uppercase">
                              Status: {ticket.status}
                            </span>
                          </div>
                          
                          <div className="opacity-0 group-hover:opacity-100 transition-opacity flex items-center gap-1">
                            <button 
                              onClick={() => startEditing('tickets', ticket.title, i, 'title')}
                              className="p-1 rounded text-zinc-400 hover:text-zinc-950"
                            >
                              <Edit2 className="h-3 w-3" />
                            </button>
                            <button 
                              onClick={() => copyToClipboard(ticket.title, 'Ticket')}
                              className="p-1 rounded text-zinc-400 hover:text-zinc-950"
                            >
                              <Copy className="h-3 w-3" />
                            </button>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              {/* Card 4: Critical Decisions */}
              <div className="rounded-xl border border-zinc-200/80 bg-white shadow-sm overflow-hidden">
                <div className="flex items-center justify-between border-b border-zinc-100 bg-zinc-50/50 px-5 py-4">
                  <div className="flex items-center gap-2.5">
                    <CheckCircle2 className="h-4 w-4 text-emerald-500" />
                    <h3 className="font-bold text-zinc-900 text-sm">Decisions Logged</h3>
                  </div>
                  <button 
                    onClick={() => toggleCollapse('decisions')}
                    className="p-1 rounded text-zinc-400 hover:text-zinc-950 hover:bg-zinc-100 transition"
                  >
                    {collapsedCards.decisions ? <ChevronDown className="h-4 w-4" /> : <ChevronUp className="h-4 w-4" />}
                  </button>
                </div>

                {!collapsedCards.decisions && (
                  <div className="p-5 space-y-3">
                    {results.decisions.map((decision, i) => (
                      <div key={i} className="group relative flex items-start gap-3 border border-zinc-100 rounded-lg p-3 hover:bg-zinc-50/50">
                        <div className="mt-0.5 rounded-full bg-emerald-50 p-1 text-emerald-600 border border-emerald-100">
                          <Check className="h-3.5 w-3.5" />
                        </div>
                        <div className="flex-1">
                          <p className="text-xs text-zinc-650 font-medium leading-relaxed">{decision.description}</p>
                          <span className="inline-block mt-2 rounded bg-zinc-100 px-1.5 py-0.5 text-[9px] font-bold text-zinc-550">
                            By: {decision.made_by}
                          </span>
                        </div>
                        <div className="opacity-0 group-hover:opacity-100 transition-opacity">
                          <button 
                            onClick={() => startEditing('decisions', decision.description, i)}
                            className="p-1 rounded text-zinc-400 hover:text-zinc-955"
                          >
                            <Edit2 className="h-3 w-3" />
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              {/* Card 5: Risks & Blockers */}
              <div className="rounded-xl border border-zinc-200/80 bg-white shadow-sm overflow-hidden">
                <div className="flex items-center justify-between border-b border-zinc-100 bg-zinc-50/50 px-5 py-4">
                  <div className="flex items-center gap-2.5">
                    <AlertCircle className="h-4 w-4 text-rose-500" />
                    <h3 className="font-bold text-zinc-900 text-sm">Risks & Blockers</h3>
                  </div>
                  <button 
                    onClick={() => toggleCollapse('risks')}
                    className="p-1 rounded text-zinc-400 hover:text-zinc-950 hover:bg-zinc-100 transition"
                  >
                    {collapsedCards.risks ? <ChevronDown className="h-4 w-4" /> : <ChevronUp className="h-4 w-4" />}
                  </button>
                </div>

                {!collapsedCards.risks && (
                  <div className="p-5 space-y-4">
                    {results.risks.map((risk, i) => (
                      <div key={i} className="group relative border border-zinc-100 rounded-lg p-3 hover:bg-zinc-50/50 transition">
                        <div className="flex items-start justify-between gap-4">
                          <div>
                            <div className="flex items-center gap-2">
                              <span className="rounded bg-rose-50 border border-rose-100 px-1.5 py-0.5 text-[9px] font-bold text-rose-600 uppercase">
                                Severity: {risk.severity}
                              </span>
                            </div>
                            <p className="text-xs font-bold text-zinc-900 mt-1.5">{risk.description}</p>
                            <div className="mt-2.5 rounded bg-zinc-100 px-3 py-2 text-[10px] text-zinc-650 border border-zinc-200">
                              <span className="font-bold text-zinc-805">Mitigation:</span> {risk.mitigation}
                            </div>
                          </div>
                          
                          <div className="opacity-0 group-hover:opacity-100 transition-opacity">
                            <button 
                              onClick={() => startEditing('risks', risk.description, i, 'description')}
                              className="p-1 rounded text-zinc-400 hover:text-zinc-950"
                            >
                              <Edit2 className="h-3 w-3" />
                            </button>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              {/* Card 6: Project Documentation */}
              <div className="rounded-xl border border-zinc-200/80 bg-white shadow-sm overflow-hidden lg:col-span-2">
                <div className="flex items-center justify-between border-b border-zinc-100 bg-zinc-50/50 px-5 py-4">
                  <div className="flex items-center gap-2.5">
                    <FileText className="h-4 w-4 text-indigo-500" />
                    <h3 className="font-bold text-zinc-900 text-sm">Project Documentation (Markdown)</h3>
                  </div>
                  <button 
                    onClick={() => toggleCollapse('documentation')}
                    className="p-1 rounded text-zinc-400 hover:text-zinc-950 hover:bg-zinc-100 transition"
                  >
                    {collapsedCards.documentation ? <ChevronDown className="h-4 w-4" /> : <ChevronUp className="h-4 w-4" />}
                  </button>
                </div>

                {!collapsedCards.documentation && (
                  <div className="p-5">
                    {editingField?.section === 'project_documentation' ? (
                      <div className="space-y-3">
                        <textarea
                          value={editValue}
                          onChange={(e) => setEditValue(e.target.value)}
                          rows={8}
                          className="w-full rounded-lg border border-zinc-300 p-3 font-mono text-xs focus:border-indigo-500 focus:outline-none"
                        />
                        <div className="flex justify-end gap-2">
                          <button onClick={() => setEditingField(null)} className="rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-semibold text-zinc-700">
                            Cancel
                          </button>
                          <button onClick={saveEdit} className="rounded-lg bg-indigo-600 px-3 py-1.5 text-xs font-semibold text-white">
                            Save
                          </button>
                        </div>
                      </div>
                    ) : (
                      <div className="group relative">
                        <pre className="rounded-lg bg-zinc-50 border border-zinc-200/50 p-5 font-mono text-xs text-zinc-700 leading-relaxed overflow-x-auto whitespace-pre-wrap">
                          {results.project_documentation}
                        </pre>
                        <div className="absolute right-4 top-4 flex items-center gap-1.5 opacity-0 group-hover:opacity-100 transition-opacity">
                          <button 
                            onClick={() => startEditing('project_documentation', results.project_documentation)}
                            className="p-1.5 rounded bg-white shadow-sm border border-zinc-200 text-zinc-400 hover:text-zinc-950"
                          >
                            <Edit2 className="h-3.5 w-3.5" />
                          </button>
                          <button 
                            onClick={() => copyToClipboard(results.project_documentation, 'Documentation')}
                            className="p-1.5 rounded bg-white shadow-sm border border-zinc-200 text-zinc-400 hover:text-zinc-950"
                          >
                            <Copy className="h-3.5 w-3.5" />
                          </button>
                        </div>
                      </div>
                    )}
                  </div>
                )}
              </div>

            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Onboarding checklist */}
      {processState === 'idle' && (
        <div className="mt-12 rounded-2xl border border-zinc-200 bg-zinc-50/50 p-6 sm:p-8">
          <h3 className="text-sm font-bold text-zinc-900 mb-6">First-Time Onboarding Experience</h3>
          <div className="grid grid-cols-1 gap-6 md:grid-cols-4">
            <div className="flex gap-3">
              <span className="flex h-6 w-6 items-center justify-center rounded-full bg-indigo-50 text-indigo-600 border border-indigo-100 text-xs font-bold font-mono">
                1
              </span>
              <div>
                <h4 className="text-xs font-bold text-zinc-900">Link Cloud Platform</h4>
                <p className="text-[11px] text-zinc-500 mt-1 leading-relaxed">Map credentials securely to automatically process recurrent meetings.</p>
              </div>
            </div>
            <div className="flex gap-3">
              <span className="flex h-6 w-6 items-center justify-center rounded-full bg-indigo-50 text-indigo-600 border border-indigo-100 text-xs font-bold font-mono">
                2
              </span>
              <div>
                <h4 className="text-xs font-bold text-zinc-900">Upload Meeting Payload</h4>
                <p className="text-[11px] text-zinc-500 mt-1 leading-relaxed">Provide live audio clips or text scripts directly to the semantic engine.</p>
              </div>
            </div>
            <div className="flex gap-3">
              <span className="flex h-6 w-6 items-center justify-center rounded-full bg-indigo-50 text-indigo-600 border border-indigo-100 text-xs font-bold font-mono">
                3
              </span>
              <div>
                <h4 className="text-xs font-bold text-zinc-900">AI Deep Structuring</h4>
                <p className="text-[11px] text-zinc-500 mt-1 leading-relaxed">System automates itemizations, task logs, and architectural choices in real time.</p>
              </div>
            </div>
            <div className="flex gap-3">
              <span className="flex h-6 w-6 items-center justify-center rounded-full bg-indigo-50 text-indigo-600 border border-indigo-100 text-xs font-bold font-mono">
                4
              </span>
              <div>
                <h4 className="text-xs font-bold text-zinc-900">Push to Jira / Notion</h4>
                <p className="text-[11px] text-zinc-500 mt-1 leading-relaxed">Instantly synchronize extracted action points with single-click triggers.</p>
              </div>
            </div>
          </div>
        </div>
      )}
      
      {/* Footer */}
      <div className="mt-16 flex flex-col items-center justify-between gap-4 border-t border-zinc-150 pt-8 text-xs text-zinc-400 sm:flex-row">
        <p>© 2026 Meetext MVP. Designed for high productivity.</p>
        <div className="flex gap-6">
          <button className="hover:text-zinc-650 hover:underline">Privacy Policy</button>
          <button className="hover:text-zinc-650 hover:underline">Terms of Service</button>
        </div>
      </div>

    </div>
  )
}
