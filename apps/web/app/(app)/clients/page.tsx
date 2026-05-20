'use client'

import React, { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  Users, Mail, Globe, MapPin, Video, CheckSquare, 
  FileText, Plus, X, Search, Calendar, ChevronRight,
  TrendingUp, Award, ExternalLink
} from 'lucide-react'
import { toast } from 'sonner'

interface MockClient {
  id: string
  companyName: string
  contactName: string
  contactEmail: string
  domain: string
  location: string
  meetingsLinked: number
  outstandingTasks: number
  docsGenerated: number
  recentNotes: string
}

const mockClients: MockClient[] = [
  {
    id: "cli-1",
    companyName: "Acme Corp",
    contactName: "Sarah Jenkins",
    contactEmail: "sarah.j@acme.com",
    domain: "acme.corp",
    location: "San Francisco, CA",
    meetingsLinked: 6,
    outstandingTasks: 3,
    docsGenerated: 4,
    recentNotes: "Acme requested automated PDF exports and custom authentication API hooks. Critical delivery targets are scheduled for Friday."
  },
  {
    id: "cli-2",
    companyName: "Stripe Integrator",
    contactName: "Alex Rivera",
    contactEmail: "alex.r@stripe.com",
    domain: "stripe-integrator.com",
    location: "Austin, TX",
    meetingsLinked: 2,
    outstandingTasks: 1,
    docsGenerated: 1,
    recentNotes: "Resolving sandbox environment webhook callback signature checks. Scheduled debugging follow-up."
  },
  {
    id: "cli-3",
    companyName: "Vercel Partner",
    contactName: "Emily Chen",
    contactEmail: "emily.c@vercel.com",
    domain: "vercel-partner.io",
    location: "New York, NY",
    meetingsLinked: 1,
    outstandingTasks: 0,
    docsGenerated: 2,
    recentNotes: "Syncing deployment modules and setting up multi-zone isolated sandbox configurations."
  }
]

export default function ClientsPage() {
  const [clients, setClients] = useState<MockClient[]>(mockClients)
  const [search, setSearch] = useState('')
  const [selectedClient, setSelectedClient] = useState<MockClient | null>(null)
  
  // Add client states
  const [showAddModal, setShowAddModal] = useState(false)
  const [newCompanyName, setNewCompanyName] = useState('')
  const [newContactName, setNewContactName] = useState('')
  const [newContactEmail, setNewContactEmail] = useState('')
  const [newDomain, setNewDomain] = useState('')
  const [newLocation, setNewLocation] = useState('')

  const handleCreateClient = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newCompanyName || !newContactEmail) return

    const newClient: MockClient = {
      id: `cli-${clients.length + 1}`,
      companyName: newCompanyName,
      contactName: newContactName || 'N/A',
      contactEmail: newContactEmail,
      domain: newDomain || 'N/A',
      location: newLocation || 'N/A',
      meetingsLinked: 0,
      outstandingTasks: 0,
      docsGenerated: 0,
      recentNotes: "New client registered. No active meeting logs available."
    }

    setClients([newClient, ...clients])
    setShowAddModal(false)
    setNewCompanyName('')
    setNewContactName('')
    setNewContactEmail('')
    setNewDomain('')
    setNewLocation('')
    toast.success("New Client Profile registered!")
  }

  const filteredClients = clients.filter(c => 
    c.companyName.toLowerCase().includes(search.toLowerCase()) ||
    c.contactName.toLowerCase().includes(search.toLowerCase()) ||
    c.contactEmail.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div className="mx-auto w-full max-w-[1400px] px-6 py-8 sm:px-8 lg:px-12 bg-white min-h-screen">
      
      {/* Page Title & Action */}
      <div className="mb-8 flex flex-col justify-between items-start gap-4 border-b border-zinc-150 pb-6 sm:flex-row sm:items-center">
        <div>
          <h1 className="text-3xl font-extrabold tracking-tight text-zinc-900">
            Client Hub
          </h1>
          <p className="mt-1 text-sm text-zinc-500">
            Track customer relations, analyze outstanding action items, and view active meetings.
          </p>
        </div>
        <button
          onClick={() => setShowAddModal(true)}
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2.5 text-xs font-semibold text-white shadow-lg shadow-indigo-600/30 hover:bg-indigo-500 hover:shadow-indigo-600/40 transition"
        >
          <Plus className="h-4 w-4" />
          Add Client Profile
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
            placeholder="Search clients..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="block w-full rounded-lg border border-zinc-250 bg-white py-2 pl-10 pr-3 text-sm placeholder:text-zinc-400 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
          />
        </div>
      </div>

      {/* Grid List of Clients */}
      <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        {filteredClients.map((client) => (
          <motion.div
            key={client.id}
            whileHover={{ y: -4, transition: { duration: 0.2 } }}
            onClick={() => setSelectedClient(client)}
            className="group cursor-pointer rounded-xl border border-zinc-200 bg-white p-6 shadow-sm hover:border-indigo-250 hover:shadow-md transition relative flex flex-col"
          >
            <div className="flex items-start justify-between gap-4 mb-4">
              <div>
                <h3 className="font-extrabold text-zinc-950 text-sm group-hover:text-indigo-600 transition-colors leading-snug">
                  {client.companyName}
                </h3>
                <span className="text-[10px] text-zinc-400 mt-1 block">Contact: {client.contactName}</span>
              </div>

              <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-zinc-50 border border-zinc-150 text-zinc-500 group-hover:bg-indigo-50 group-hover:text-indigo-650 transition-colors">
                <Users className="h-4.5 w-4.5" />
              </div>
            </div>

            {/* Metas info */}
            <div className="space-y-2 mb-6">
              <div className="flex items-center gap-2 text-xs text-zinc-500">
                <Mail className="h-3.5 w-3.5 text-zinc-400" />
                <span className="truncate">{client.contactEmail}</span>
              </div>
              <div className="flex items-center gap-2 text-xs text-zinc-500">
                <Globe className="h-3.5 w-3.5 text-zinc-400" />
                <span className="truncate">{client.domain}</span>
              </div>
            </div>

            {/* Metas count grid */}
            <div className="border-t border-zinc-50 pt-4 flex items-center justify-between text-xs text-zinc-500">
              <div className="flex items-center gap-3">
                <span className="flex items-center gap-1">
                  <Video className="h-3.5 w-3.5 text-zinc-400" /> {client.meetingsLinked}
                </span>
                <span className="flex items-center gap-1">
                  <CheckSquare className="h-3.5 w-3.5 text-zinc-400" /> {client.outstandingTasks}
                </span>
                <span className="flex items-center gap-1">
                  <FileText className="h-3.5 w-3.5 text-zinc-400" /> {client.docsGenerated}
                </span>
              </div>
              <span className="text-[10px] text-zinc-400">{client.location}</span>
            </div>
          </motion.div>
        ))}
      </div>

      {/* Add Client Modal Overlay */}
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
                <h3 className="font-extrabold text-zinc-950 text-base">New Client Profile</h3>
                <button onClick={() => setShowAddModal(false)} className="rounded text-zinc-400 hover:bg-zinc-100 p-1">
                  <X className="h-5 w-5" />
                </button>
              </div>

              <form onSubmit={handleCreateClient} className="space-y-4">
                <div>
                  <label className="block text-xs font-semibold text-zinc-500 mb-1">Company Name</label>
                  <input
                    type="text"
                    required
                    placeholder="e.g. Acme Corporation"
                    value={newCompanyName}
                    onChange={(e) => setNewCompanyName(e.target.value)}
                    className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-zinc-500 mb-1">Contact Name</label>
                  <input
                    type="text"
                    placeholder="e.g. Sarah Jenkins"
                    value={newContactName}
                    onChange={(e) => setNewContactName(e.target.value)}
                    className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-zinc-500 mb-1">Contact Email</label>
                  <input
                    type="email"
                    required
                    placeholder="e.g. sarah@acme.corp"
                    value={newContactEmail}
                    onChange={(e) => setNewContactEmail(e.target.value)}
                    className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-zinc-500 mb-1">Website Domain</label>
                    <input
                      type="text"
                      placeholder="e.g. acme.corp"
                      value={newDomain}
                      onChange={(e) => setNewDomain(e.target.value)}
                      className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-zinc-500 mb-1">Location</label>
                    <input
                      type="text"
                      placeholder="e.g. Austin, TX"
                      value={newLocation}
                      onChange={(e) => setNewLocation(e.target.value)}
                      className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                    />
                  </div>
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
                    Register Client
                  </button>
                </div>
              </form>
            </motion.div>
          </div>
        )}
      </AnimatePresence>

      {/* Slide-over detailed client workspace drawer */}
      <AnimatePresence>
        {selectedClient && (
          <div className="fixed inset-0 z-50 flex justify-end">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setSelectedClient(null)}
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
                      {selectedClient.domain}
                    </span>
                    <span className="text-[10px] text-zinc-405 font-medium">· Profile Details</span>
                  </div>
                  <h2 className="text-base font-extrabold text-zinc-950 leading-snug">{selectedClient.companyName}</h2>
                </div>
                <button
                  onClick={() => setSelectedClient(null)}
                  className="rounded-lg p-1.5 text-zinc-450 hover:bg-zinc-100 hover:text-zinc-900 transition"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>

              {/* Drawer Content */}
              <div className="flex-1 overflow-y-auto p-6 space-y-8">
                
                {/* Notes */}
                <div>
                  <h3 className="text-xs font-bold text-zinc-450 uppercase mb-2">Latest Action Directives</h3>
                  <p className="text-xs text-zinc-600 leading-relaxed bg-zinc-50 border border-zinc-150 rounded-xl p-4">
                    {selectedClient.recentNotes}
                  </p>
                </div>

                {/* Info parameters */}
                <div className="space-y-4 border border-zinc-100 rounded-xl p-4 text-xs">
                  <div className="flex justify-between border-b border-zinc-50 pb-2">
                    <span className="text-zinc-400 font-semibold">Primary Contact</span>
                    <span className="font-bold text-zinc-800">{selectedClient.contactName}</span>
                  </div>
                  <div className="flex justify-between border-b border-zinc-50 pb-2">
                    <span className="text-zinc-400 font-semibold">Contact Email</span>
                    <span className="font-bold text-zinc-800">{selectedClient.contactEmail}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-zinc-400 font-semibold">HQ Location</span>
                    <span className="font-bold text-zinc-800">{selectedClient.location}</span>
                  </div>
                </div>

                {/* Counts statistics */}
                <div className="grid grid-cols-3 gap-4">
                  <div className="border border-zinc-100 rounded-xl p-4 text-center hover:bg-zinc-50/50 transition">
                    <Video className="mx-auto h-5 w-5 text-zinc-400 mb-2" />
                    <span className="block text-lg font-extrabold text-zinc-900">{selectedClient.meetingsLinked}</span>
                    <span className="text-[10px] text-zinc-500 font-semibold">Meetings linked</span>
                  </div>
                  <div className="border border-zinc-100 rounded-xl p-4 text-center hover:bg-zinc-50/50 transition">
                    <CheckSquare className="mx-auto h-5 w-5 text-zinc-400 mb-2" />
                    <span className="block text-lg font-extrabold text-zinc-900">{selectedClient.outstandingTasks}</span>
                    <span className="text-[10px] text-zinc-500 font-semibold">Action tasks</span>
                  </div>
                  <div className="border border-zinc-100 rounded-xl p-4 text-center hover:bg-zinc-50/50 transition">
                    <FileText className="mx-auto h-5 w-5 text-zinc-400 mb-2" />
                    <span className="block text-lg font-extrabold text-zinc-900">{selectedClient.docsGenerated}</span>
                    <span className="text-[10px] text-zinc-500 font-semibold">Specs generated</span>
                  </div>
                </div>

              </div>

              {/* Footer */}
              <div className="border-t border-zinc-150 px-6 py-4 flex items-center justify-end bg-zinc-50/50">
                <button 
                  onClick={() => setSelectedClient(null)}
                  className="inline-flex items-center gap-1.5 rounded-lg bg-indigo-600 px-4 py-2 text-xs font-semibold text-white shadow-sm hover:bg-indigo-500 transition"
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
