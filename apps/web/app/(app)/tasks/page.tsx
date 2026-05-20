'use client'

import React, { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  CheckSquare, Plus, Search, Filter, Calendar, 
  User, Sparkles, ChevronRight, X, AlertCircle, 
  CheckCircle2, Clock, Trash2, ArrowRight
} from 'lucide-react'
import { toast } from 'sonner'

interface MockTask {
  id: string
  title: string
  description: string
  priority: 'low' | 'medium' | 'high' | 'urgent'
  assignee: string
  dueDate: string
  status: 'todo' | 'in_progress' | 'review' | 'done'
  project: string
  aiConfidence?: number
}

const initialTasks: MockTask[] = [
  {
    id: "task-1",
    title: "Build Authentication Endpoint",
    description: "Design JWT secure authentication handler routines with database validation scopes.",
    priority: "high",
    assignee: "Sarah",
    dueDate: "May 22, 2026",
    status: "todo",
    project: "Mobile Auth SDK",
    aiConfidence: 94
  },
  {
    id: "task-2",
    title: "Finalize SVG Logo Assets Pack",
    description: "Export full multi-scale SVG branding logos for dynamic light/dark web headers.",
    priority: "medium",
    assignee: "David",
    dueDate: "May 21, 2026",
    status: "in_progress",
    project: "Branding",
    aiConfidence: 88
  },
  {
    id: "task-3",
    title: "Address Outdated Stripe webhook callbacks",
    description: "Sandbox endpoints suffer mismatches under live testing. Connect Stripe debug tools.",
    priority: "urgent",
    assignee: "Sarah",
    dueDate: "May 20, 2026",
    status: "review",
    project: "Stripe Core Integrations",
    aiConfidence: 91
  },
  {
    id: "task-4",
    title: "Select Database Storage Architecture",
    description: "Formally align PostgreSQL storage parameters rather than MongoDB clusters.",
    priority: "medium",
    assignee: "Team",
    dueDate: "May 19, 2026",
    status: "done",
    project: "Database Selection",
    aiConfidence: 99
  }
]

const columns: { id: MockTask['status']; label: string; color: string }[] = [
  { id: 'todo', label: 'Todo', color: 'bg-zinc-100 border-zinc-200 text-zinc-800' },
  { id: 'in_progress', label: 'In Progress', color: 'bg-indigo-50 border-indigo-100 text-indigo-700' },
  { id: 'review', label: 'In Review', color: 'bg-amber-50 border-amber-100 text-amber-700' },
  { id: 'done', label: 'Done', color: 'bg-emerald-50 border-emerald-100 text-emerald-700' }
]

export default function TasksPage() {
  const [tasks, setTasks] = useState<MockTask[]>(initialTasks)
  const [search, setSearch] = useState('')
  const [filterPriority, setFilterPriority] = useState<'all' | 'low' | 'medium' | 'high' | 'urgent'>('all')
  const [selectedTask, setSelectedTask] = useState<MockTask | null>(null)
  
  // Add task states
  const [showAddModal, setShowAddModal] = useState(false)
  const [newTitle, setNewTitle] = useState('')
  const [newDesc, setNewDesc] = useState('')
  const [newAssignee, setNewAssignee] = useState('')
  const [newPriority, setNewPriority] = useState<MockTask['priority']>('medium')
  const [newProject, setNewProject] = useState('')

  const handleCreateTask = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newTitle || !newAssignee) return

    const newTask: MockTask = {
      id: `task-${tasks.length + 1}`,
      title: newTitle,
      description: newDesc,
      priority: newPriority,
      assignee: newAssignee,
      dueDate: new Date(Date.now() + 86400000 * 3).toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' }),
      status: "todo",
      project: newProject || "General",
      aiConfidence: undefined
    }

    setTasks([...tasks, newTask])
    setShowAddModal(false)
    setNewTitle('')
    setNewDesc('')
    setNewAssignee('')
    setNewProject('')
    toast.success("Task created and placed in Todo column!")
  }

  const moveTaskStatus = (taskId: string, currentStatus: MockTask['status'], direction: 'forward' | 'backward') => {
    const statuses: MockTask['status'][] = ['todo', 'in_progress', 'review', 'done']
    const currentIndex = statuses.indexOf(currentStatus)
    let nextIndex = direction === 'forward' ? currentIndex + 1 : currentIndex - 1
    
    if (nextIndex < 0 || nextIndex >= statuses.length) return

    setTasks(prev => prev.map(t => {
      if (t.id === taskId) {
        return { ...t, status: statuses[nextIndex] }
      }
      return t
    }))
    toast.success(`Task moved to ${statuses[nextIndex].replace('_', ' ')}`)
  }

  const deleteTask = (taskId: string) => {
    setTasks(prev => prev.filter(t => t.id !== taskId))
    setSelectedTask(null)
    toast.success("Task deleted successfully")
  }

  const filteredTasks = tasks.filter(t => {
    const matchesSearch = t.title.toLowerCase().includes(search.toLowerCase()) || 
                          t.project.toLowerCase().includes(search.toLowerCase())
    const matchesPriority = filterPriority === 'all' || t.priority === filterPriority
    return matchesSearch && matchesPriority
  })

  return (
    <div className="mx-auto w-full max-w-[1400px] px-6 py-8 sm:px-8 lg:px-12 bg-white min-h-screen">
      
      {/* Page Title & Action */}
      <div className="mb-8 flex flex-col justify-between items-start gap-4 border-b border-zinc-150 pb-6 sm:flex-row sm:items-center">
        <div>
          <h1 className="text-3xl font-extrabold tracking-tight text-zinc-900">
            Action Tasks Board
          </h1>
          <p className="mt-1 text-sm text-zinc-500">
            Inspect, allocate, and synchronize tasks extracted dynamically from meetings.
          </p>
        </div>
        <button
          onClick={() => setShowAddModal(true)}
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2.5 text-xs font-semibold text-white shadow-lg shadow-indigo-600/30 hover:bg-indigo-500 hover:shadow-indigo-600/40 transition"
        >
          <Plus className="h-4 w-4" />
          Add Action Task
        </button>
      </div>

      {/* Filter panel */}
      <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center">
        <div className="relative flex-1 max-w-xs">
          <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
            <Search className="h-4 w-4 text-zinc-400" />
          </div>
          <input
            type="text"
            placeholder="Filter tasks..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="block w-full rounded-lg border border-zinc-250 bg-white py-2 pl-10 pr-3 text-sm placeholder:text-zinc-400 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
          />
        </div>

        <div className="flex items-center gap-1.5 rounded-lg border border-zinc-200 bg-white px-2.5 py-1.5 text-xs text-zinc-700 shadow-sm">
          <Filter className="h-3.5 w-3.5 text-zinc-400" />
          <select
            value={filterPriority}
            onChange={(e) => setFilterPriority(e.target.value as any)}
            className="bg-transparent focus:outline-none font-semibold cursor-pointer"
          >
            <option value="all">All Priorities</option>
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
            <option value="urgent">Urgent</option>
          </select>
        </div>
      </div>

      {/* Kanban Grid */}
      <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4 items-start">
        {columns.map((col) => {
          const colTasks = filteredTasks.filter(t => t.status === col.id)
          return (
            <div key={col.id} className="rounded-2xl border border-zinc-150 bg-zinc-50/50 p-4 min-h-[500px]">
              {/* Column Header */}
              <div className="mb-4 flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span className={`rounded px-2 py-0.5 text-xs font-semibold border ${col.color}`}>
                    {col.label}
                  </span>
                  <span className="text-xs font-bold text-zinc-400 bg-zinc-200/50 rounded-full px-2 py-0.5">
                    {colTasks.length}
                  </span>
                </div>
              </div>

              {/* Column Cards */}
              <div className="space-y-3">
                {colTasks.map((task) => (
                  <motion.div
                    key={task.id}
                    layoutId={`task-card-${task.id}`}
                    onClick={() => setSelectedTask(task)}
                    className="group cursor-pointer rounded-xl border border-zinc-200 bg-white p-4 shadow-sm hover:border-indigo-250 transition-all hover:shadow-md relative overflow-hidden"
                  >
                    {/* Top Row */}
                    <div className="flex items-center justify-between gap-2 mb-2">
                      <span className="rounded bg-zinc-100 border border-zinc-200 px-1.5 py-0.5 text-[9px] font-bold text-zinc-550 uppercase">
                        {task.project}
                      </span>
                      <span className={`rounded px-1.5 py-0.5 text-[9px] font-bold uppercase border ${
                        task.priority === 'urgent' || task.priority === 'high'
                          ? 'bg-rose-50 border-rose-100 text-rose-600'
                          : 'bg-zinc-100 border-zinc-200 text-zinc-600'
                      }`}>
                        {task.priority}
                      </span>
                    </div>

                    <h4 className="font-bold text-zinc-950 text-xs leading-snug group-hover:text-indigo-600 transition-colors">
                      {task.title}
                    </h4>

                    {/* AI Confidence badge */}
                    {task.aiConfidence && (
                      <div className="flex items-center gap-1 mt-2 text-[9px] text-indigo-500 font-semibold">
                        <Sparkles className="h-3 w-3" />
                        <span>AI Extraction · {task.aiConfidence}% confidence</span>
                      </div>
                    )}

                    {/* Card bottom details */}
                    <div className="mt-4 flex items-center justify-between border-t border-zinc-50 pt-3 text-[10px] text-zinc-400">
                      <span className="font-semibold text-zinc-600 bg-zinc-100 rounded px-1.5 py-0.5">
                        {task.assignee}
                      </span>
                      <span className="flex items-center gap-1">
                        <Calendar className="h-3 w-3" /> {task.dueDate}
                      </span>
                    </div>

                    {/* Move controls inside card hover */}
                    <div className="absolute inset-x-0 bottom-0 bg-zinc-950/90 text-white py-2 px-3 flex justify-between items-center opacity-0 group-hover:opacity-100 transition-all duration-200 text-[10px]">
                      <button 
                        onClick={(e) => { e.stopPropagation(); moveTaskStatus(task.id, task.status, 'backward') }}
                        disabled={task.status === 'todo'}
                        className="disabled:opacity-30 hover:underline px-1 py-0.5 font-bold"
                      >
                        ◀ Back
                      </button>
                      <button 
                        onClick={(e) => { e.stopPropagation(); moveTaskStatus(task.id, task.status, 'forward') }}
                        disabled={task.status === 'done'}
                        className="disabled:opacity-30 hover:underline px-1 py-0.5 font-bold"
                      >
                        Forward ▶
                      </button>
                    </div>
                  </motion.div>
                ))}

                {colTasks.length === 0 && (
                  <div className="py-12 text-center text-zinc-450 border border-dashed border-zinc-200 rounded-xl text-xs bg-white/50">
                    No active tasks
                  </div>
                )}
              </div>
            </div>
          )
        })}
      </div>

      {/* Task Creation Modal */}
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
                <h3 className="font-extrabold text-zinc-950 text-base">Allocate Action Task</h3>
                <button onClick={() => setShowAddModal(false)} className="rounded text-zinc-400 hover:bg-zinc-100 p-1">
                  <X className="h-5 w-5" />
                </button>
              </div>

              <form onSubmit={handleCreateTask} className="space-y-4">
                <div>
                  <label className="block text-xs font-semibold text-zinc-500 mb-1">Task Title</label>
                  <input
                    type="text"
                    required
                    placeholder="e.g. Map recurring subscription webhooks"
                    value={newTitle}
                    onChange={(e) => setNewTitle(e.target.value)}
                    className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-zinc-500 mb-1">Assignee</label>
                  <input
                    type="text"
                    required
                    placeholder="e.g. Sarah or David"
                    value={newAssignee}
                    onChange={(e) => setNewAssignee(e.target.value)}
                    className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-zinc-500 mb-1">Priority</label>
                    <select
                      value={newPriority}
                      onChange={(e) => setNewPriority(e.target.value as any)}
                      className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                    >
                      <option value="low">Low</option>
                      <option value="medium">Medium</option>
                      <option value="high">High</option>
                      <option value="urgent">Urgent</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-zinc-500 mb-1">Project Code</label>
                    <input
                      type="text"
                      placeholder="e.g. API Gateway"
                      value={newProject}
                      onChange={(e) => setNewProject(e.target.value)}
                      className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-zinc-500 mb-1">Task Details</label>
                  <textarea
                    rows={3}
                    placeholder="Short description of task criteria..."
                    value={newDesc}
                    onChange={(e) => setNewDesc(e.target.value)}
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
                    Allocate Task
                  </button>
                </div>
              </form>
            </motion.div>
          </div>
        )}
      </AnimatePresence>

      {/* Task detailed modal overlay */}
      <AnimatePresence>
        {selectedTask && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setSelectedTask(null)}
              className="fixed inset-0 bg-black/60 backdrop-blur-sm"
            />
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="relative w-full max-w-md rounded-2xl border border-zinc-200 bg-white p-6 shadow-2xl z-10"
            >
              <div className="flex items-center justify-between border-b border-zinc-100 pb-3 mb-4">
                <div className="flex items-center gap-2">
                  <span className="rounded bg-zinc-100 border border-zinc-200 px-2 py-0.5 text-[9px] font-bold text-zinc-650 uppercase">
                    {selectedTask.project}
                  </span>
                  <span className="text-[10px] text-zinc-400">· Action Point</span>
                </div>
                <button onClick={() => setSelectedTask(null)} className="rounded text-zinc-400 hover:bg-zinc-100 p-1">
                  <X className="h-5 w-5" />
                </button>
              </div>

              <h3 className="font-extrabold text-zinc-950 text-base leading-snug mb-2">
                {selectedTask.title}
              </h3>
              
              <p className="text-xs text-zinc-600 leading-relaxed mb-6 bg-zinc-50 border border-zinc-150 rounded-lg p-3">
                {selectedTask.description || "No deep parameters declared."}
              </p>

              <div className="space-y-2.5 text-xs border-t border-zinc-100 pt-4 mb-6">
                <div className="flex justify-between">
                  <span className="text-zinc-400 font-semibold">Assignee</span>
                  <span className="font-bold text-zinc-800">{selectedTask.assignee}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-zinc-400 font-semibold">Priority Status</span>
                  <span className="font-bold uppercase text-rose-600">{selectedTask.priority}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-zinc-400 font-semibold">Delivery Target</span>
                  <span className="text-zinc-600">{selectedTask.dueDate}</span>
                </div>
              </div>

              <div className="flex justify-between items-center pt-4 border-t border-zinc-100">
                <button 
                  onClick={() => deleteTask(selectedTask.id)} 
                  className="rounded-lg hover:bg-rose-50 text-rose-600 p-2 border border-transparent hover:border-rose-100"
                >
                  <Trash2 className="h-4.5 w-4.5" />
                </button>
                <button 
                  onClick={() => setSelectedTask(null)} 
                  className="rounded-lg bg-indigo-600 px-4 py-2 text-xs font-semibold text-white hover:bg-indigo-500"
                >
                  Close Details
                </button>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>

    </div>
  )
}
