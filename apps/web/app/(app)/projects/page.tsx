'use client'

import React, { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  FolderKanban, Calendar, Clock, AlertCircle, 
  CheckCircle2, Plus, ArrowRight, X, Video, 
  CheckSquare, FileText, Activity, ShieldCheck,
  TrendingUp, Users, FileCode
} from 'lucide-react'
import { toast } from 'sonner'

interface MockProject {
  id: string
  name: string
  client: string
  status: 'planning' | 'active' | 'review' | 'completed'
  progress: number
  health: 'healthy' | 'at_risk' | 'critical'
  meetingsCount: number
  tasksCount: number
  docsCount: number
  description: string
  startDate: string
  endDate: string
  recentActivity: string[]
}

const mockProjects: MockProject[] = [
  {
    id: "proj-1",
    name: "Mobile Client Authentication SDK",
    client: "Acme Corp",
    status: "active",
    progress: 68,
    health: "healthy",
    meetingsCount: 4,
    tasksCount: 8,
    docsCount: 3,
    description: "Design and implement custom secure authentication workflows supporting multi-workspace isolation configurations.",
    startDate: "May 01, 2026",
    endDate: "June 15, 2026",
    recentActivity: [
      "David updated technical documentation pages.",
      "Sarah created authentication integration tasks.",
      "Ollama extracted storage selection decisions (PostgreSQL)."
    ]
  },
  {
    id: "proj-2",
    name: "Stripe Recurring Billing Engine",
    client: "Stripe Integrator",
    status: "active",
    progress: 45,
    health: "at_risk",
    meetingsCount: 2,
    tasksCount: 5,
    docsCount: 1,
    description: "Debugging sandboxed recurrent subscription callback endpoints and webhooks integration failures.",
    startDate: "May 10, 2026",
    endDate: "June 05, 2026",
    recentActivity: [
      "Stripe webhook debugging meeting failed signature.",
      "David resolved sandbox payload mismatch errors."
    ]
  },
  {
    id: "proj-3",
    name: "Onboarding Flow UX/UI Refactoring",
    client: "Internal Product",
    status: "planning",
    progress: 15,
    health: "healthy",
    meetingsCount: 1,
    tasksCount: 4,
    docsCount: 2,
    description: "Modernizing core workspace initialization pipelines and dropzone states with premium CSS/Framer motion.",
    startDate: "May 18, 2026",
    endDate: "July 01, 2026",
    recentActivity: [
      "Refactored mobile sidebar overlay layout modules."
    ]
  }
]

export default function ProjectsPage() {
  const [projects, setProjects] = useState<MockProject[]>(mockProjects)
  const [selectedProject, setSelectedProject] = useState<MockProject | null>(null)
  const [showAddModal, setShowAddModal] = useState(false)
  const [newProjectName, setNewProjectName] = useState('')
  const [newProjectClient, setNewProjectClient] = useState('')
  const [newProjectDesc, setNewProjectDesc] = useState('')

  const handleCreateProject = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newProjectName || !newProjectClient) return

    const newProject: MockProject = {
      id: `proj-${projects.length + 1}`,
      name: newProjectName,
      client: newProjectClient,
      status: "planning",
      progress: 0,
      health: "healthy",
      meetingsCount: 0,
      tasksCount: 0,
      docsCount: 0,
      description: newProjectDesc,
      startDate: new Date().toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' }),
      endDate: "TBD",
      recentActivity: ["Project initialized in workspace."]
    }

    setProjects([newProject, ...projects])
    setShowAddModal(false)
    setNewProjectName('')
    setNewProjectClient('')
    setNewProjectDesc('')
    toast.success("New project workspace registered!")
  }

  return (
    <div className="mx-auto w-full max-w-[1400px] px-6 py-8 sm:px-8 lg:px-12 bg-white min-h-screen">
      
      {/* Page Title & Add action */}
      <div className="mb-8 flex flex-col justify-between items-start gap-4 border-b border-zinc-150 pb-6 sm:flex-row sm:items-center">
        <div>
          <h1 className="text-3xl font-extrabold tracking-tight text-zinc-900">
            Project Workspaces
          </h1>
          <p className="mt-1 text-sm text-zinc-500">
            Organize meetings, document files, and actionable tickets grouped by product workflows.
          </p>
        </div>
        <button
          onClick={() => setShowAddModal(true)}
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2.5 text-xs font-semibold text-white shadow-lg shadow-indigo-600/30 hover:bg-indigo-500 hover:shadow-indigo-600/40 transition"
        >
          <Plus className="h-4 w-4" />
          Create Project
        </button>
      </div>

      {/* Grid List of Projects */}
      <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        {projects.map((project) => (
          <motion.div
            key={project.id}
            whileHover={{ y: -4, transition: { duration: 0.2 } }}
            onClick={() => setSelectedProject(project)}
            className="group cursor-pointer rounded-xl border border-zinc-200 bg-white p-6 shadow-sm hover:border-indigo-250 hover:shadow-md transition relative flex flex-col"
          >
            <div className="flex items-start justify-between gap-4 mb-3">
              <div>
                <span className="rounded bg-zinc-100 border border-zinc-200 px-2 py-0.5 text-[9px] font-bold text-zinc-550 uppercase">
                  {project.client}
                </span>
                <h3 className="font-extrabold text-zinc-950 text-sm group-hover:text-indigo-600 transition-colors mt-1.5 leading-snug">
                  {project.name}
                </h3>
              </div>

              {/* Health Indicator badge */}
              <span className={`rounded-full px-2.5 py-0.5 text-[9px] font-bold uppercase border flex items-center gap-1 ${
                project.health === 'healthy'
                  ? 'bg-emerald-50 border-emerald-100 text-emerald-600'
                  : project.health === 'at_risk'
                  ? 'bg-amber-50 border-amber-100 text-amber-600'
                  : 'bg-rose-50 border-rose-100 text-rose-600'
              }`}>
                <div className={`h-1.5 w-1.5 rounded-full ${
                  project.health === 'healthy' ? 'bg-emerald-500' : project.health === 'at_risk' ? 'bg-amber-500' : 'bg-rose-500'
                }`} />
                {project.health.replace('_', ' ')}
              </span>
            </div>

            <p className="text-xs text-zinc-500 leading-relaxed mb-6 flex-1">
              {project.description}
            </p>

            {/* Progress Bar */}
            <div className="mb-6 space-y-1.5">
              <div className="flex items-center justify-between text-[10px] font-semibold text-zinc-400">
                <span>AI Documentation Progress</span>
                <span className="text-zinc-650">{project.progress}%</span>
              </div>
              <div className="h-1.5 w-full bg-zinc-100 rounded-full overflow-hidden">
                <div 
                  className="h-full bg-indigo-600 rounded-full" 
                  style={{ width: `${project.progress}%` }}
                />
              </div>
            </div>

            {/* Counts & Metas */}
            <div className="border-t border-zinc-50 pt-4 flex items-center justify-between text-xs text-zinc-500">
              <div className="flex items-center gap-3">
                <span className="flex items-center gap-1">
                  <Video className="h-3.5 w-3.5 text-zinc-400" /> {project.meetingsCount}
                </span>
                <span className="flex items-center gap-1">
                  <CheckSquare className="h-3.5 w-3.5 text-zinc-400" /> {project.tasksCount}
                </span>
                <span className="flex items-center gap-1">
                  <FileText className="h-3.5 w-3.5 text-zinc-400" /> {project.docsCount}
                </span>
              </div>
              <span className="text-[10px] text-zinc-400">{project.startDate}</span>
            </div>
          </motion.div>
        ))}
      </div>

      {/* Add Project Modal overlay */}
      <AnimatePresence>
        {showAddModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setShowAddModal(false)}
              className="fixed inset-0 bg-black/60 backdrop-blur-sm"
            />
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="relative w-full max-w-md rounded-2xl border border-zinc-200 bg-white p-6 shadow-2xl z-10"
            >
              <div className="flex items-center justify-between border-b border-zinc-100 pb-3 mb-5">
                <h3 className="font-extrabold text-zinc-950 text-base">New Project Workspace</h3>
                <button onClick={() => setShowAddModal(false)} className="rounded text-zinc-400 hover:bg-zinc-100 p-1">
                  <X className="h-5 w-5" />
                </button>
              </div>

              <form onSubmit={handleCreateProject} className="space-y-4">
                <div>
                  <label className="block text-xs font-semibold text-zinc-500 mb-1">Project Name</label>
                  <input
                    type="text"
                    required
                    placeholder="e.g. Notion Auto-Sync Integration"
                    value={newProjectName}
                    onChange={(e) => setNewProjectName(e.target.value)}
                    className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-zinc-500 mb-1">Associated Client</label>
                  <input
                    type="text"
                    required
                    placeholder="e.g. Acme Corp or Internal"
                    value={newProjectClient}
                    onChange={(e) => setNewProjectClient(e.target.value)}
                    className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-zinc-500 mb-1">Description</label>
                  <textarea
                    rows={3}
                    placeholder="Short summary of project scope..."
                    value={newProjectDesc}
                    onChange={(e) => setNewProjectDesc(e.target.value)}
                    className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                  />
                </div>
                <div className="flex justify-end gap-2 pt-3 border-t border-zinc-100">
                  <button 
                    type="button" 
                    onClick={() => setShowAddModal(false)} 
                    className="rounded-lg border border-zinc-200 px-4 py-2 text-xs font-semibold text-zinc-700 hover:bg-zinc-50"
                  >
                    Cancel
                  </button>
                  <button 
                    type="submit" 
                    className="rounded-lg bg-indigo-600 px-4 py-2 text-xs font-semibold text-white hover:bg-indigo-500 shadow-md shadow-indigo-600/20"
                  >
                    Create Workspace
                  </button>
                </div>
              </form>
            </motion.div>
          </div>
        )}
      </AnimatePresence>

      {/* Slide-over detailed Project workspace drawer */}
      <AnimatePresence>
        {selectedProject && (
          <div className="fixed inset-0 z-50 flex justify-end">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setSelectedProject(null)}
              className="fixed inset-0 bg-black/55 backdrop-blur-sm"
            />
            <motion.div
              initial={{ x: '100%' }}
              animate={{ x: 0 }}
              exit={{ x: '100%' }}
              transition={{ type: 'spring', damping: 25, stiffness: 200 }}
              className="relative w-full max-w-2xl bg-white h-full shadow-2xl flex flex-col"
            >
              {/* Header */}
              <div className="flex items-start justify-between border-b border-zinc-150 px-6 py-5 bg-zinc-50/50">
                <div>
                  <div className="flex items-center gap-2 mb-1">
                    <span className="rounded bg-zinc-200 px-2 py-0.5 text-[9px] font-bold text-zinc-650 uppercase border border-zinc-300">
                      {selectedProject.client}
                    </span>
                    <span className="text-[10px] text-zinc-400">· Active Workspace</span>
                  </div>
                  <h2 className="text-base font-extrabold text-zinc-950 leading-snug">{selectedProject.name}</h2>
                </div>
                <button
                  onClick={() => setSelectedProject(null)}
                  className="rounded-lg p-1.5 text-zinc-450 hover:bg-zinc-100 hover:text-zinc-900 transition"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>

              {/* Drawer Content */}
              <div className="flex-1 overflow-y-auto p-6 space-y-8">
                
                {/* Scope */}
                <div>
                  <h3 className="text-xs font-bold text-zinc-450 uppercase mb-2">Workspace Objective</h3>
                  <p className="text-xs text-zinc-600 leading-relaxed bg-zinc-50 border border-zinc-150 rounded-xl p-4">
                    {selectedProject.description}
                  </p>
                </div>

                {/* Recent activity timeline */}
                <div>
                  <h3 className="text-sm font-bold text-zinc-900 mb-4 flex items-center gap-2">
                    <Activity className="h-4 w-4 text-indigo-500" />
                    AI Activity Log
                  </h3>
                  <div className="space-y-4">
                    {selectedProject.recentActivity.map((activity, i) => (
                      <div key={i} className="flex gap-3 text-xs leading-relaxed">
                        <div className="h-1.5 w-1.5 rounded-full bg-indigo-500 mt-1.5 shrink-0" />
                        <span className="text-zinc-650 font-semibold">{activity}</span>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Dynamic linked components */}
                <div className="grid grid-cols-3 gap-4">
                  <div className="border border-zinc-100 rounded-xl p-4 text-center hover:bg-zinc-50/50 transition">
                    <Video className="mx-auto h-5 w-5 text-zinc-400 mb-2" />
                    <span className="block text-lg font-extrabold text-zinc-900">{selectedProject.meetingsCount}</span>
                    <span className="text-[10px] text-zinc-500 font-semibold">Meetings linked</span>
                  </div>
                  <div className="border border-zinc-100 rounded-xl p-4 text-center hover:bg-zinc-50/50 transition">
                    <CheckSquare className="mx-auto h-5 w-5 text-zinc-400 mb-2" />
                    <span className="block text-lg font-extrabold text-zinc-900">{selectedProject.tasksCount}</span>
                    <span className="text-[10px] text-zinc-500 font-semibold">Action tasks</span>
                  </div>
                  <div className="border border-zinc-100 rounded-xl p-4 text-center hover:bg-zinc-50/50 transition">
                    <FileText className="mx-auto h-5 w-5 text-zinc-400 mb-2" />
                    <span className="block text-lg font-extrabold text-zinc-900">{selectedProject.docsCount}</span>
                    <span className="text-[10px] text-zinc-500 font-semibold">AI Documents</span>
                  </div>
                </div>

              </div>

              {/* Footer */}
              <div className="border-t border-zinc-150 px-6 py-4 flex items-center justify-between bg-zinc-50/50 text-xs">
                <span className="text-[10px] text-zinc-405 font-medium flex items-center gap-1.5">
                  <Calendar className="h-3.5 w-3.5 animate-pulse" /> Range: {selectedProject.startDate} - {selectedProject.endDate}
                </span>
                <button 
                  onClick={() => {
                    toast.success("Synchronized successfully")
                    setSelectedProject(null)
                  }}
                  className="inline-flex items-center gap-1.5 rounded-lg bg-indigo-600 px-3.5 py-2 text-xs font-semibold text-white shadow-sm hover:bg-indigo-500 transition"
                >
                  Open Board <ArrowRight className="h-4 w-4" />
                </button>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>

    </div>
  )
}
