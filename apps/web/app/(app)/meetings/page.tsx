'use client'

import React, { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  Video, Music, FileText, Search, Filter, 
  Calendar, Clock, User, ArrowRight, ExternalLink, 
  Download, Copy, Check, X, ChevronRight, Sparkles,
  Play, FileCode, CheckSquare, Compass
} from 'lucide-react'
import { toast } from 'sonner'

interface MockMeeting {
  id: string
  title: string
  date: string
  duration: string
  type: 'video' | 'audio' | 'pdf'
  status: 'completed' | 'processing' | 'failed'
  client: string
  project: string
  summary: string
  transcript: { time: string; speaker: string; text: string }[]
  tasks: { title: string; assignee: string; priority: 'low' | 'medium' | 'high' }[]
  tickets: string[]
}

const mockMeetings: MockMeeting[] = [
  {
    id: "meet-1",
    title: "Mobile App API Integration & Auth Review",
    date: "May 20, 2026",
    duration: "45 mins",
    type: "video",
    status: "completed",
    client: "Acme Corp",
    project: "Mobile Client SDK",
    summary: "Finalized standard security protocol utilizing robust JWT authorization strategies. Formally agreed on standardizing database selection with PostgreSQL to resolve relational consistency. Sarah leading endpoint construction targeted for Friday delivery.",
    transcript: [
      { time: "00:12", speaker: "Sarah", text: "Okay, let's get started. Thanks everyone for joining. First on the agenda is the new API integration for the mobile app." },
      { time: "01:05", speaker: "Sarah", text: "I'll take the lead on building the authentication endpoint. We should be able to deliver it by Friday." },
      { time: "02:30", speaker: "David", text: "Sounds excellent. What's the plan for database architecture selection? PostgreSQL or MongoDB?" },
      { time: "03:15", speaker: "Sarah", text: "We formally agreed on PostgreSQL standard to leverage transactional capabilities and JSONB fields for flexibility." }
    ],
    tasks: [
      { title: "Implement JWT Security Endpoint", assignee: "Sarah", priority: "high" },
      { title: "Deliver SVG Logo Assets Pack", assignee: "David", priority: "medium" }
    ],
    tickets: ["FEAT-104: PDF Export Mod", "BUG-219: Hook Refactor"]
  },
  {
    id: "meet-2",
    title: "Acme Corp PDF Export Architecture Alignment",
    date: "May 18, 2026",
    duration: "18 mins",
    type: "audio",
    status: "completed",
    client: "Acme Corp",
    project: "Statement Sheet Generator",
    summary: "Reviewed billing sheet generator requirements. Outlined technical specs for PDF synthesis engine using standard serverless routines. Scheduled follow-up mapping for Wednesday.",
    transcript: [
      { time: "00:05", speaker: "John (Client)", text: "We need automated export to PDF and CSV formats for all monthly billing sheets." },
      { time: "01:45", speaker: "David", text: "We committed to supporting that. I can wire the document generation service using a simple node microservice." }
    ],
    tasks: [
      { title: "Create PDF export templates", assignee: "David", priority: "high" }
    ],
    tickets: ["FEAT-109: PDF Invoice generator"]
  },
  {
    id: "meet-3",
    title: "Payment Webhook Troubleshooting Session",
    date: "May 17, 2026",
    duration: "30 mins",
    type: "video",
    status: "failed",
    client: "Stripe Integrator",
    project: "Billing Engine",
    summary: "Debug session for recurrent webhook authorization signature failure. Failed to pinpoint precise environment variable issue. Scheduled second debug slot.",
    transcript: [
      { time: "00:30", speaker: "David", text: "We are getting key mismatches under Sandbox callbacks. I need to review direct Stripe documentation." }
    ],
    tasks: [],
    tickets: ["BUG-222: Verify Stripe webhook callback"]
  },
  {
    id: "meet-4",
    title: "Weekly Sync & Next Steps",
    date: "May 15, 2026",
    duration: "60 mins",
    type: "pdf",
    status: "completed",
    client: "Internal Team",
    project: "Meetext MVP Core",
    summary: "General team sync discussing deployment, marketing assets, and product onboarding walkthrough. All modules are currently deployed under isolated sandbox clusters.",
    transcript: [],
    tasks: [],
    tickets: []
  }
]

export default function MeetingsPage() {
  const [meetings, setMeetings] = useState<MockMeeting[]>(mockMeetings)
  const [search, setSearch] = useState('')
  const [filterType, setFilterType] = useState<'all' | 'video' | 'audio' | 'pdf'>('all')
  const [filterStatus, setFilterStatus] = useState<'all' | 'completed' | 'processing' | 'failed'>('all')
  const [selectedMeeting, setSelectedMeeting] = useState<MockMeeting | null>(null)
  const [copiedTranscript, setCopiedTranscript] = useState<number | null>(null)

  const filteredMeetings = meetings.filter(m => {
    const matchesSearch = m.title.toLowerCase().includes(search.toLowerCase()) || 
                          m.client.toLowerCase().includes(search.toLowerCase()) || 
                          m.project.toLowerCase().includes(search.toLowerCase())
    const matchesType = filterType === 'all' || m.type === filterType
    const matchesStatus = filterStatus === 'all' || m.status === filterStatus
    return matchesSearch && matchesType && matchesStatus
  })

  const copyTranscriptText = (text: string, index: number) => {
    navigator.clipboard.writeText(text)
    setCopiedTranscript(index)
    toast.success("Transcript segment copied")
    setTimeout(() => setCopiedTranscript(null), 2000)
  }

  return (
    <div className="mx-auto w-full max-w-[1400px] px-6 py-8 sm:px-8 lg:px-12 bg-white min-h-screen">
      
      {/* Page Title */}
      <div className="mb-8 border-b border-zinc-100 pb-6">
        <h1 className="text-3xl font-extrabold tracking-tight text-zinc-900">
          Meetings Library
        </h1>
        <p className="mt-1 text-sm text-zinc-500">
          Browse, query, and review transcripts or structural AI assets from past conversations.
        </p>
      </div>

      {/* Filter and Search Bar */}
      <div className="mb-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="relative flex-1 max-w-md">
          <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
            <Search className="h-4 w-4 text-zinc-400" />
          </div>
          <input
            type="text"
            placeholder="Search meetings by title, client, or project..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="block w-full rounded-lg border border-zinc-250 bg-white py-2 pl-10 pr-3 text-sm placeholder:text-zinc-400 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
          />
        </div>

        <div className="flex flex-wrap items-center gap-3">
          {/* Type Filter */}
          <div className="flex items-center gap-1.5 rounded-lg border border-zinc-200 bg-white px-2.5 py-1.5 text-xs text-zinc-700 shadow-sm">
            <Filter className="h-3.5 w-3.5 text-zinc-400" />
            <select
              value={filterType}
              onChange={(e) => setFilterType(e.target.value as any)}
              className="bg-transparent focus:outline-none font-semibold cursor-pointer"
            >
              <option value="all">All Formats</option>
              <option value="video">Video</option>
              <option value="audio">Audio</option>
              <option value="pdf">PDF Upload</option>
            </select>
          </div>

          {/* Status Filter */}
          <div className="flex items-center gap-1.5 rounded-lg border border-zinc-200 bg-white px-2.5 py-1.5 text-xs text-zinc-700 shadow-sm">
            <select
              value={filterStatus}
              onChange={(e) => setFilterStatus(e.target.value as any)}
              className="bg-transparent focus:outline-none font-semibold cursor-pointer"
            >
              <option value="all">All Statuses</option>
              <option value="completed">Completed</option>
              <option value="processing">Processing</option>
              <option value="failed">Failed</option>
            </select>
          </div>
        </div>
      </div>

      {/* Grid List of Meetings */}
      <div className="grid grid-cols-1 gap-5 md:grid-cols-2 lg:grid-cols-3">
        {filteredMeetings.length > 0 ? (
          filteredMeetings.map((meeting) => (
            <motion.div
              key={meeting.id}
              whileHover={{ y: -4, transition: { duration: 0.2 } }}
              onClick={() => setSelectedMeeting(meeting)}
              className="group cursor-pointer rounded-xl border border-zinc-200 bg-white p-5 shadow-sm transition hover:border-indigo-250 hover:shadow-md relative overflow-hidden"
            >
              <div className="flex items-start justify-between gap-4 mb-4">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-zinc-50 border border-zinc-150 text-zinc-650 group-hover:bg-indigo-50 group-hover:text-indigo-650 transition-colors">
                  {meeting.type === 'video' && <Video className="h-5 w-5" />}
                  {meeting.type === 'audio' && <Music className="h-5 w-5" />}
                  {meeting.type === 'pdf' && <FileText className="h-5 w-5" />}
                </div>

                <span className={`rounded-full px-2.5 py-0.5 text-[10px] font-bold uppercase border ${
                  meeting.status === 'completed' 
                    ? 'bg-emerald-50 border-emerald-100 text-emerald-600'
                    : meeting.status === 'failed'
                    ? 'bg-rose-50 border-rose-100 text-rose-600'
                    : 'bg-indigo-50 border-indigo-100 text-indigo-600'
                }`}>
                  {meeting.status}
                </span>
              </div>

              <h3 className="font-bold text-zinc-950 text-sm group-hover:text-indigo-600 transition-colors leading-snug">
                {meeting.title}
              </h3>
              
              <div className="mt-3 flex items-center gap-3 text-[10px] text-zinc-400">
                <span className="font-semibold text-zinc-600 bg-zinc-100 rounded px-1.5 py-0.5">{meeting.client}</span>
                <span>{meeting.project}</span>
              </div>

              <div className="mt-5 flex items-center justify-between border-t border-zinc-50 pt-4 text-xs text-zinc-500">
                <div className="flex items-center gap-1">
                  <Calendar className="h-3.5 w-3.5" />
                  <span>{meeting.date}</span>
                </div>
                <div className="flex items-center gap-1">
                  <Clock className="h-3.5 w-3.5" />
                  <span>{meeting.duration}</span>
                </div>
              </div>
            </motion.div>
          ))
        ) : (
          <div className="col-span-full py-16 text-center border border-dashed border-zinc-200 rounded-xl bg-zinc-50/50">
            <Video className="mx-auto h-8 w-8 text-zinc-300 mb-3" />
            <h3 className="font-bold text-zinc-800 text-sm">No matches found</h3>
            <p className="text-xs text-zinc-400 mt-1">Refine your query or filters above.</p>
          </div>
        )}
      </div>

      {/* Slide-over details drawer panel */}
      <AnimatePresence>
        {selectedMeeting && (
          <div className="fixed inset-0 z-50 flex justify-end">
            {/* Backdrop */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setSelectedMeeting(null)}
              className="fixed inset-0 bg-black/55 backdrop-blur-sm"
            />

            {/* Panel */}
            <motion.div
              initial={{ x: '100%' }}
              animate={{ x: 0 }}
              exit={{ x: '100%' }}
              transition={{ type: 'spring', damping: 25, stiffness: 200 }}
              className="relative w-full max-w-2xl bg-white h-full shadow-2xl flex flex-col"
            >
              {/* Drawer Header */}
              <div className="flex items-start justify-between border-b border-zinc-150 px-6 py-5 bg-zinc-50/50">
                <div>
                  <div className="flex items-center gap-2 mb-1">
                    <span className="rounded bg-zinc-200 px-2 py-0.5 text-[9px] font-bold text-zinc-650 uppercase border border-zinc-300">
                      {selectedMeeting.client}
                    </span>
                    <span className="text-[10px] text-zinc-400">· {selectedMeeting.project}</span>
                  </div>
                  <h2 className="text-base font-extrabold text-zinc-950 leading-snug">{selectedMeeting.title}</h2>
                </div>
                <button
                  onClick={() => setSelectedMeeting(null)}
                  className="rounded-lg p-1.5 text-zinc-450 hover:bg-zinc-100 hover:text-zinc-900 transition"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>

              {/* Drawer Content */}
              <div className="flex-1 overflow-y-auto p-6 space-y-8">
                
                {/* Section 1: Executive Summary */}
                <div>
                  <div className="flex items-center gap-2 mb-3">
                    <span className="inline-flex items-center gap-1 rounded-full bg-indigo-50 px-2 py-0.5 text-[10px] font-semibold text-indigo-600 border border-indigo-100">
                      <Sparkles className="h-3 w-3" />
                      AI Executive Summary
                    </span>
                  </div>
                  <p className="text-sm leading-relaxed text-zinc-600 bg-zinc-50 border border-zinc-150 rounded-xl p-4">
                    {selectedMeeting.summary}
                  </p>
                </div>

                {/* Section 2: Action Tasks list */}
                {selectedMeeting.tasks.length > 0 && (
                  <div>
                    <h3 className="text-sm font-bold text-zinc-900 mb-3 flex items-center gap-2">
                      <CheckSquare className="h-4 w-4 text-indigo-500" />
                      Extracted Tasks
                    </h3>
                    <div className="space-y-3">
                      {selectedMeeting.tasks.map((task, i) => (
                        <div key={i} className="flex items-center justify-between border border-zinc-100 rounded-lg p-3 hover:bg-zinc-50/50 transition">
                          <div>
                            <h4 className="text-xs font-bold text-zinc-800">{task.title}</h4>
                            <span className="text-[10px] text-zinc-400 mt-1 block">Assignee: {task.assignee}</span>
                          </div>
                          <span className={`rounded-full px-2 py-0.5 text-[9px] font-bold uppercase border ${
                            task.priority === 'high' 
                              ? 'bg-rose-50 border-rose-100 text-rose-600'
                              : 'bg-amber-50 border-amber-100 text-amber-600'
                          }`}>
                            {task.priority}
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Section 3: Tickets sync */}
                {selectedMeeting.tickets.length > 0 && (
                  <div>
                    <h3 className="text-sm font-bold text-zinc-900 mb-3 flex items-center gap-2">
                      <Compass className="h-4 w-4 text-indigo-500" />
                      Linked Product Tickets
                    </h3>
                    <div className="flex flex-wrap gap-2">
                      {selectedMeeting.tickets.map((ticket, i) => (
                        <span key={i} className="inline-flex items-center gap-1.5 rounded-lg border border-zinc-200 bg-white px-3 py-1.5 text-xs text-zinc-700 shadow-sm">
                          <FileCode className="h-3.5 w-3.5 text-indigo-500" />
                          {ticket}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Section 4: Deep Transcript timestamps */}
                {selectedMeeting.transcript.length > 0 ? (
                  <div>
                    <h3 className="text-sm font-bold text-zinc-900 mb-3">Transcript Stream</h3>
                    <div className="space-y-4 max-h-[300px] overflow-y-auto pr-2 border border-zinc-100 rounded-xl p-4 bg-zinc-50/30">
                      {selectedMeeting.transcript.map((t, i) => (
                        <div key={i} className="group relative text-xs">
                          <div className="flex items-center justify-between mb-1">
                            <div className="flex items-center gap-2">
                              <span className="font-bold text-zinc-900">{t.speaker}</span>
                              <span className="font-mono text-[10px] text-zinc-400 bg-zinc-100 rounded px-1">{t.time}</span>
                            </div>
                            
                            <button
                              onClick={() => copyTranscriptText(t.text, i)}
                              className="opacity-0 group-hover:opacity-100 transition-opacity p-0.5 rounded text-zinc-400 hover:text-zinc-900"
                            >
                              {copiedTranscript === i ? <Check className="h-3 w-3 text-emerald-500" /> : <Copy className="h-3 w-3" />}
                            </button>
                          </div>
                          <p className="text-zinc-650 leading-relaxed pr-8">{t.text}</p>
                        </div>
                      ))}
                    </div>
                  </div>
                ) : (
                  <div className="text-center py-6 border border-dashed border-zinc-200 rounded-xl">
                    <h3 className="text-xs font-bold text-zinc-500">No text transcript stream</h3>
                    <p className="text-[10px] text-zinc-400">PDF files or manual summaries lack step-by-step transcript streams.</p>
                  </div>
                )}

              </div>

              {/* Drawer Footer */}
              <div className="border-t border-zinc-150 px-6 py-4 flex items-center justify-between bg-zinc-50/50">
                <span className="text-[10px] text-zinc-405 font-medium flex items-center gap-1.5">
                  <Clock className="h-3.5 w-3.5" /> Checked by Semantic Engine
                </span>
                <div className="flex gap-2">
                  <button 
                    onClick={() => {
                      toast.success("Markdown exported successfully")
                      setSelectedMeeting(null)
                    }}
                    className="inline-flex items-center gap-1.5 rounded-lg border border-zinc-200 bg-white px-3.5 py-2 text-xs font-semibold text-zinc-700 shadow-sm hover:bg-zinc-50 transition"
                  >
                    <Download className="h-3.5 w-3.5" /> Export Specs
                  </button>
                </div>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>

    </div>
  )
}
