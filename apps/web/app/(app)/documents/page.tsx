'use client'

import React, { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  FileText, Search, Folder, Download, Copy, Plus, 
  ExternalLink, FileCode, Check, X, Calendar, User,
  FileSpreadsheet, Edit3
} from 'lucide-react'
import { toast } from 'sonner'

interface MockDoc {
  id: string
  title: string
  type: 'Technical Spec' | 'Client Report' | 'Sprint Notes' | 'General'
  project: string
  date: string
  author: string
  content: string
}

const mockDocs: MockDoc[] = [
  {
    id: "doc-1",
    title: "Mobile App JWT Authentication Spec",
    type: "Technical Spec",
    project: "Mobile Auth SDK",
    date: "May 20, 2026",
    author: "AI Writer",
    content: `## 1. Context & Architecture Selection
Following the alignment meeting on database selections, standard relational schemas will leverage **PostgreSQL** multi-workspace partition columns. For secure client validation, secure stateless JWT sessions will be deployed.

## 2. JWT Signature Scopes
* **Header**: standard SHA256 signature configurations.
* **Payload values**: 
  - \`workspace_id\`: strict tenant mapping key.
  - \`user_role\`: authorization permissions boundary.

## 3. Storage Security Parameters
Token expiration is configured with strict grace windows to protect local states. Web hooks rely on secure payload headers to prevent sandbox forgery.`
  },
  {
    id: "doc-2",
    title: "Acme Corp PDF Sheet Requirements",
    type: "Client Report",
    project: "Statement Generator",
    date: "May 18, 2026",
    author: "AI Writer",
    content: `## Acme Corp PDF Exporter Requirements

### 1. Delivery Criteria
Provide highly styled downloadable reports detailing monthly payment breakdowns and account states. PDF schemas must follow strict company asset templates.

### 2. Export Actions
* **Format support**: Automated PDF layout rendering & raw CSV spreadsheets.
* **Integrations**: Linked to Workspace billing sheets.`
  },
  {
    id: "doc-3",
    title: "Billing Setup Debugging Session Log",
    type: "Sprint Notes",
    project: "Billing Engine",
    date: "May 17, 2026",
    author: "AI Writer",
    content: `## Stripe Integration Debug Session

### 1. Webhook Signature Errors
Troubleshot callback integrity failures within sandbox environments.Mismatches seem localized inside webhook verification functions.

### 2. Resolution Checklist
* Re-verify webhook secret configurations.
* Setup Stripe CLI redirect tunnels.`
  }
]

export default function DocumentsPage() {
  const [docs, setDocs] = useState<MockDoc[]>(mockDocs)
  const [search, setSearch] = useState('')
  const [selectedDoc, setSelectedDoc] = useState<MockDoc | null>(null)
  const [copiedText, setCopiedText] = useState(false)
  const [isEditing, setIsEditing] = useState(false)
  const [editingContent, setEditingContent] = useState('')

  const handleCopy = (text: string) => {
    navigator.clipboard.writeText(text)
    setCopiedText(true)
    toast.success("Document copied to clipboard")
    setTimeout(() => setCopiedText(false), 2000)
  }

  const handleDownload = (doc: MockDoc) => {
    const blob = new Blob([doc.content], { type: 'text/markdown' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${doc.title.replace(/\s+/g, '_')}.md`
    a.click()
    toast.success("Markdown file downloaded successfully")
  }

  const saveEditedDoc = () => {
    if (!selectedDoc) return
    setDocs(prev => prev.map(d => {
      if (d.id === selectedDoc.id) {
        return { ...d, content: editingContent }
      }
      return d
    }))
    setSelectedDoc(prev => prev ? { ...prev, content: editingContent } : null)
    setIsEditing(false)
    toast.success("Document updated successfully")
  }

  const filteredDocs = docs.filter(d => 
    d.title.toLowerCase().includes(search.toLowerCase()) ||
    d.project.toLowerCase().includes(search.toLowerCase()) ||
    d.type.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div className="mx-auto w-full max-w-[1400px] px-6 py-8 sm:px-8 lg:px-12 bg-white min-h-screen">
      
      {/* Page Header */}
      <div className="mb-8 border-b border-zinc-150 pb-6">
        <h1 className="text-3xl font-extrabold tracking-tight text-zinc-900">
          AI Documents Workspace
        </h1>
        <p className="mt-1 text-sm text-zinc-500">
          Manage, customize, and download structured documents compiled from meetings.
        </p>
      </div>

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
        
        {/* List of Documents column */}
        <div className="lg:col-span-1 space-y-6 border-r border-zinc-100 pr-0 lg:pr-8">
          <div className="relative">
            <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
              <Search className="h-4 w-4 text-zinc-400" />
            </div>
            <input
              type="text"
              placeholder="Search documentation..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="block w-full rounded-lg border border-zinc-250 bg-white py-2.5 pl-10 pr-3 text-sm placeholder:text-zinc-400 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
            />
          </div>

          <div className="space-y-3">
            {filteredDocs.map((doc) => (
              <div
                key={doc.id}
                onClick={() => {
                  setSelectedDoc(doc)
                  setEditingContent(doc.content)
                  setIsEditing(false)
                }}
                className={`group cursor-pointer rounded-xl border p-4 transition-all duration-200 ${
                  selectedDoc?.id === doc.id
                    ? 'border-indigo-500 bg-indigo-50/10 shadow-sm'
                    : 'border-zinc-200 bg-white hover:border-zinc-300 hover:shadow-sm'
                }`}
              >
                <div className="flex items-center gap-2 mb-2">
                  <span className={`rounded-full px-2 py-0.5 text-[9px] font-bold uppercase border ${
                    doc.type === 'Technical Spec' 
                      ? 'bg-rose-50 border-rose-100 text-rose-600'
                      : doc.type === 'Client Report'
                      ? 'bg-emerald-50 border-emerald-100 text-emerald-600'
                      : 'bg-zinc-100 border-zinc-200 text-zinc-650'
                  }`}>
                    {doc.type}
                  </span>
                </div>
                <h3 className={`font-bold text-xs leading-snug transition-colors ${
                  selectedDoc?.id === doc.id ? 'text-indigo-650' : 'text-zinc-900 group-hover:text-indigo-600'
                }`}>
                  {doc.title}
                </h3>
                <span className="text-[10px] text-zinc-400 mt-2 block font-semibold">
                  Project: {doc.project}
                </span>
              </div>
            ))}

            {filteredDocs.length === 0 && (
              <div className="text-center py-12 border border-dashed border-zinc-200 rounded-xl text-xs text-zinc-400">
                No matching documents
              </div>
            )}
          </div>
        </div>

        {/* Preview Document details column */}
        <div className="lg:col-span-2 space-y-6">
          {selectedDoc ? (
            <div className="rounded-2xl border border-zinc-200 bg-white shadow-sm overflow-hidden flex flex-col min-h-[600px]">
              {/* Header */}
              <div className="border-b border-zinc-150 px-6 py-4 flex flex-col justify-between gap-4 sm:flex-row sm:items-center bg-zinc-50/50">
                <div>
                  <h2 className="font-extrabold text-zinc-950 text-base leading-snug">{selectedDoc.title}</h2>
                  <div className="flex items-center gap-3 mt-1.5 text-[10px] text-zinc-400">
                    <span className="font-bold text-zinc-600">By: {selectedDoc.author}</span>
                    <span>Date: {selectedDoc.date}</span>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  <button
                    onClick={() => handleCopy(selectedDoc.content)}
                    className="p-2 rounded-lg border border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-600 transition"
                  >
                    {copiedText ? <Check className="h-4 w-4 text-emerald-500" /> : <Copy className="h-4 w-4" />}
                  </button>
                  <button
                    onClick={() => handleDownload(selectedDoc)}
                    className="p-2 rounded-lg border border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-600 transition"
                  >
                    <Download className="h-4 w-4" />
                  </button>
                  <button
                    onClick={() => setIsEditing(!isEditing)}
                    className={`inline-flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-semibold border transition ${
                      isEditing 
                        ? 'bg-indigo-600 border-indigo-600 text-white hover:bg-indigo-500' 
                        : 'bg-white border-zinc-200 hover:bg-zinc-50 text-zinc-700'
                    }`}
                  >
                    <Edit3 className="h-3.5 w-3.5" />
                    {isEditing ? 'Save Changes' : 'Edit Spec'}
                  </button>
                </div>
              </div>

              {/* Editor / Markdown preview panel */}
              <div className="flex-1 p-6 overflow-y-auto">
                {isEditing ? (
                  <div className="space-y-4 h-full flex flex-col">
                    <textarea
                      value={editingContent}
                      onChange={(e) => setEditingContent(e.target.value)}
                      className="w-full flex-1 min-h-[400px] rounded-lg border border-zinc-300 p-4 font-mono text-xs focus:border-indigo-500 focus:outline-none"
                    />
                    <div className="flex justify-end gap-2">
                      <button 
                        onClick={() => setIsEditing(false)} 
                        className="rounded-lg border border-zinc-200 px-4 py-2 text-xs font-semibold text-zinc-700 hover:bg-zinc-50"
                      >
                        Cancel
                      </button>
                      <button 
                        onClick={saveEditedDoc} 
                        className="rounded-lg bg-indigo-600 px-4 py-2 text-xs font-semibold text-white hover:bg-indigo-500"
                      >
                        Save Spec Document
                      </button>
                    </div>
                  </div>
                ) : (
                  <div className="prose prose-sm max-w-none text-zinc-650 leading-relaxed space-y-4">
                    {/* Simplified render formatting logic */}
                    {selectedDoc.content.split('\n').map((line, index) => {
                      if (line.startsWith('## ')) {
                        return <h2 key={index} className="text-base font-bold text-zinc-950 border-b border-zinc-100 pb-1 pt-4">{line.replace('## ', '')}</h2>
                      }
                      if (line.startsWith('### ')) {
                        return <h3 key={index} className="text-sm font-bold text-zinc-900 pt-2">{line.replace('### ', '')}</h3>
                      }
                      if (line.startsWith('* ')) {
                        return <li key={index} className="list-disc ml-5 text-xs">{line.replace('* ', '')}</li>
                      }
                      if (line.trim() === '') {
                        return <div key={index} className="h-2" />
                      }
                      return <p key={index} className="text-xs">{line}</p>
                    })}
                  </div>
                )}
              </div>
            </div>
          ) : (
            <div className="rounded-2xl border border-dashed border-zinc-200 bg-zinc-50/50 p-12 text-center flex flex-col items-center justify-center min-h-[600px]">
              <Folder className="h-10 w-10 text-zinc-300 mb-3" />
              <h3 className="font-extrabold text-zinc-800 text-sm">Select a Document</h3>
              <p className="text-xs text-zinc-400 mt-1">Review structural specs, sprint summaries, and transcripts compiled by Ollama.</p>
            </div>
          )}
        </div>

      </div>

    </div>
  )
}
