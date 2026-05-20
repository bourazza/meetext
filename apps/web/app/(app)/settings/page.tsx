'use client'

import React, { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  Settings as SettingsIcon, User, Layers, Sparkles, Bell, 
  Share2, CreditCard, Shield, Plus, Check, Key, Eye, 
  Trash2, Mail, ExternalLink, RefreshCw, Compass
} from 'lucide-react'
import { toast } from 'sonner'

type TabId = 'profile' | 'workspace' | 'ai' | 'notifications' | 'integrations' | 'billing' | 'security'

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState<TabId>('profile')

  // Profile States
  const [fullName, setFullName] = useState('Zaki Bourazza')
  const [email, setEmail] = useState('zaki@meetext.ai')
  const [userRole, setUserRole] = useState('Administrator')

  // Workspace States
  const [workspaceName, setWorkspaceName] = useState('Bourazza Hub')
  const [inviteEmail, setInviteEmail] = useState('')
  const [members, setMembers] = useState([
    { email: 'sarah.j@meetext.ai', role: 'admin' },
    { email: 'david.r@meetext.ai', role: 'member' }
  ])

  // AI settings States
  const [aiModel, setAiModel] = useState('llama3-8b')
  const [ollamaUrl, setOllamaUrl] = useState('http://localhost:11434')
  const [promptPreference, setPromptPreference] = useState('technical')

  // Integrations states
  const [linkedInteg, setLinkedInteg] = useState({
    notion: true,
    jira: false,
    slack: true,
    linear: false
  })

  // API Key States
  const [apiKeys, setApiKeys] = useState([
    { id: 'key-1', name: 'Production Webhook', key: 'mt_live_••••••••••••e92a' }
  ])
  const [newKeyName, setNewKeyName] = useState('')

  const handleInvite = (e: React.FormEvent) => {
    e.preventDefault()
    if (!inviteEmail) return
    setMembers([...members, { email: inviteEmail, role: 'member' }])
    setInviteEmail('')
    toast.success(`Invite sent successfully to ${inviteEmail}`)
  }

  const handleGenerateKey = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newKeyName) return
    const rand = Math.random().toString(36).substring(7)
    setApiKeys([...apiKeys, { id: `key-${Date.now()}`, name: newKeyName, key: `mt_live_••••••••••••${rand}` }])
    setNewKeyName('')
    toast.success("API key successfully generated!")
  }

  const handleProfileSave = (e: React.FormEvent) => {
    e.preventDefault()
    toast.success("Profile parameters updated!")
  }

  const handleWorkspaceSave = (e: React.FormEvent) => {
    e.preventDefault()
    toast.success("Workspace parameters synchronized!")
  }

  const handleAISave = (e: React.FormEvent) => {
    e.preventDefault()
    toast.success("AI Configuration registered!")
  }

  const toggleInteg = (key: keyof typeof linkedInteg) => {
    setLinkedInteg(prev => ({ ...prev, [key]: !prev[key] }))
    toast.success(`${key.toUpperCase()} integration state toggled!`)
  }

  return (
    <div className="mx-auto w-full max-w-[1400px] px-6 py-8 sm:px-8 lg:px-12 bg-white min-h-screen">
      
      {/* Page Title */}
      <div className="mb-8 border-b border-zinc-150 pb-6">
        <h1 className="text-3xl font-extrabold tracking-tight text-zinc-900">
          Workspace Settings
        </h1>
        <p className="mt-1 text-sm text-zinc-500">
          Configure profile details, manage localized LLMs, configure web hooks, and review invoices.
        </p>
      </div>

      <div className="flex flex-col lg:flex-row gap-8 items-start">
        
        {/* Left Tab navigation */}
        <aside className="w-full lg:w-64 shrink-0 space-y-1">
          {[
            { id: 'profile', label: 'Profile Settings', icon: User },
            { id: 'workspace', label: 'Workspace Hub', icon: Layers },
            { id: 'ai', label: 'AI Configurations', icon: Sparkles },
            { id: 'notifications', label: 'Notifications & Alerts', icon: Bell },
            { id: 'integrations', label: 'Integrations Sync', icon: Share2 },
            { id: 'billing', label: 'Billing & Plan', icon: CreditCard },
            { id: 'security', label: 'Security & Access', icon: Shield }
          ].map((tab) => {
            const isActive = activeTab === tab.id
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id as any)}
                className={`w-full flex items-center gap-3 rounded-lg px-3.5 py-2.5 text-xs font-semibold transition ${
                  isActive 
                    ? 'bg-zinc-950 text-white shadow-sm'
                    : 'text-zinc-650 hover:bg-zinc-50 hover:text-zinc-900'
                }`}
              >
                <tab.icon className={`h-4 w-4 ${isActive ? 'text-indigo-400' : 'text-zinc-400'}`} />
                <span>{tab.label}</span>
              </button>
            )
          })}
        </aside>

        {/* Main tabs container */}
        <div className="flex-1 rounded-2xl border border-zinc-200 bg-white p-6 sm:p-8 shadow-sm w-full min-h-[500px]">
          <AnimatePresence mode="wait">
            <motion.div
              key={activeTab}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.2 }}
            >
              {/* Profile Tab */}
              {activeTab === 'profile' && (
                <form onSubmit={handleProfileSave} className="space-y-6">
                  <div className="border-b border-zinc-100 pb-3 mb-6">
                    <h2 className="text-base font-bold text-zinc-950">Profile Settings</h2>
                    <p className="text-xs text-zinc-500 mt-0.5">Manage details regarding your individual active profile credentials.</p>
                  </div>

                  <div className="flex items-center gap-4 mb-6">
                    <div className="flex h-16 w-16 items-center justify-center rounded-full bg-indigo-900/40 text-indigo-300 font-bold border border-indigo-500/20 text-lg">
                      {fullName.charAt(0).toUpperCase()}
                    </div>
                    <div>
                      <button type="button" className="rounded-lg border border-zinc-200 px-3.5 py-1.5 text-xs font-semibold text-zinc-700 hover:bg-zinc-50 shadow-sm transition">
                        Change Avatar
                      </button>
                    </div>
                  </div>

                  <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                    <div>
                      <label className="block text-xs font-semibold text-zinc-500 mb-1">Full Name</label>
                      <input
                        type="text"
                        value={fullName}
                        onChange={(e) => setFullName(e.target.value)}
                        className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-xs focus:border-indigo-500 focus:outline-none"
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-semibold text-zinc-500 mb-1">Email Address</label>
                      <input
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-xs focus:border-indigo-500 focus:outline-none"
                      />
                    </div>
                  </div>

                  <div>
                    <label className="block text-xs font-semibold text-zinc-500 mb-1">Workspace Assignment Role</label>
                    <input
                      type="text"
                      disabled
                      value={userRole}
                      className="block w-full rounded-lg border border-zinc-200 bg-zinc-50 px-3 py-2 text-xs text-zinc-400 font-semibold cursor-not-allowed"
                    />
                  </div>

                  <button type="submit" className="rounded-lg bg-indigo-600 px-4 py-2 text-xs font-semibold text-white hover:bg-indigo-500 shadow-md shadow-indigo-600/20">
                    Save Profiles Changes
                  </button>
                </form>
              )}

              {/* Workspace Tab */}
              {activeTab === 'workspace' && (
                <div className="space-y-8">
                  <form onSubmit={handleWorkspaceSave} className="space-y-4">
                    <div className="border-b border-zinc-100 pb-3 mb-6">
                      <h2 className="text-base font-bold text-zinc-950">Workspace Preferences</h2>
                      <p className="text-xs text-zinc-500 mt-0.5">Customize properties regarding the joint corporate repository hubs.</p>
                    </div>

                    <div>
                      <label className="block text-xs font-semibold text-zinc-500 mb-1">Workspace Identifier</label>
                      <input
                        type="text"
                        value={workspaceName}
                        onChange={(e) => setWorkspaceName(e.target.value)}
                        className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-xs focus:border-indigo-500 focus:outline-none"
                      />
                    </div>

                    <button type="submit" className="rounded-lg bg-indigo-600 px-4 py-2 text-xs font-semibold text-white hover:bg-indigo-500">
                      Save Workspace Details
                    </button>
                  </form>

                  {/* Members Invites */}
                  <div className="border-t border-zinc-150 pt-6 space-y-4">
                    <h3 className="text-xs font-bold text-zinc-800 uppercase">Invite Workspace Members</h3>
                    
                    <form onSubmit={handleInvite} className="flex gap-2">
                      <input
                        type="email"
                        required
                        placeholder="collaborator@company.com"
                        value={inviteEmail}
                        onChange={(e) => setInviteEmail(e.target.value)}
                        className="flex-1 rounded-lg border border-zinc-250 bg-white px-3 py-2 text-xs focus:border-indigo-500 focus:outline-none"
                      />
                      <button type="submit" className="rounded-lg bg-indigo-600 px-4 py-2 text-xs font-semibold text-white hover:bg-indigo-500">
                        Send Invite
                      </button>
                    </form>

                    <div className="space-y-3">
                      {members.map((member, i) => (
                        <div key={i} className="flex items-center justify-between text-xs border border-zinc-100 rounded-lg p-3">
                          <span className="font-semibold text-zinc-800">{member.email}</span>
                          <span className="rounded bg-zinc-100 border border-zinc-200 px-1.5 py-0.5 text-[9px] uppercase font-bold text-zinc-550">
                            {member.role}
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              )}

              {/* AI Config Tab */}
              {activeTab === 'ai' && (
                <form onSubmit={handleAISave} className="space-y-6">
                  <div className="border-b border-zinc-100 pb-3 mb-6">
                    <h2 className="text-base font-bold text-zinc-950">AI Orchestration Settings</h2>
                    <p className="text-xs text-zinc-500 mt-0.5">Define connected LLM nodes, prompt rules, and Whisper templates.</p>
                  </div>

                  <div>
                    <label className="block text-xs font-semibold text-zinc-500 mb-1">Active LLM Model</label>
                    <select
                      value={aiModel}
                      onChange={(e) => setAiModel(e.target.value)}
                      className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-xs focus:border-indigo-500 focus:outline-none"
                    >
                      <option value="llama3-8b">Ollama: Llama 3 (8B)</option>
                      <option value="mistral-7b">Ollama: Mistral (7B)</option>
                      <option value="gpt-4o">OpenAI GPT-4o (Cloud)</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-xs font-semibold text-zinc-500 mb-1">Local Ollama API Endpoint</label>
                    <input
                      type="text"
                      value={ollamaUrl}
                      onChange={(e) => setOllamaUrl(e.target.value)}
                      placeholder="http://localhost:11434"
                      className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-xs focus:border-indigo-500 focus:outline-none"
                    />
                  </div>

                  <div>
                    <label className="block text-xs font-semibold text-zinc-500 mb-1">Prompt Spec Templates</label>
                    <select
                      value={promptPreference}
                      onChange={(e) => setPromptPreference(e.target.value)}
                      className="block w-full rounded-lg border border-zinc-250 bg-white px-3 py-2 text-xs focus:border-indigo-500 focus:outline-none"
                    >
                      <option value="technical">Engineering & Architecture Detail focused</option>
                      <option value="business">Executive Business brief focused</option>
                      <option value="sprint">Agile sprint backlog task focused</option>
                    </select>
                  </div>

                  <button type="submit" className="rounded-lg bg-indigo-600 px-4 py-2 text-xs font-semibold text-white hover:bg-indigo-500">
                    Save AI Settings
                  </button>
                </form>
              )}

              {/* Notifications Tab */}
              {activeTab === 'notifications' && (
                <div className="space-y-6">
                  <div className="border-b border-zinc-100 pb-3 mb-6">
                    <h2 className="text-base font-bold text-zinc-950">Notification Preferences</h2>
                    <p className="text-xs text-zinc-500 mt-0.5">Determine when and where transactional meeting summaries are sent.</p>
                  </div>

                  <div className="space-y-4">
                    {[
                      { title: 'Email Digest Notifications', desc: 'Recurrent daily summaries of action task item allocations.' },
                      { title: 'Slack Instant Alerts', desc: 'Pushes extracted Jira tickets directly into linked Slack channels.' },
                      { title: 'Webhook callback events', desc: 'Dispatches payload bodies immediately upon Whisper completion.' }
                    ].map((notif, i) => (
                      <div key={i} className="flex items-start justify-between border border-zinc-100 rounded-lg p-4">
                        <div>
                          <h4 className="text-xs font-bold text-zinc-800">{notif.title}</h4>
                          <p className="text-[10px] text-zinc-400 mt-0.5">{notif.desc}</p>
                        </div>
                        <input type="checkbox" defaultChecked className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-zinc-300 rounded cursor-pointer mt-1" />
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Integrations Tab */}
              {activeTab === 'integrations' && (
                <div className="space-y-6">
                  <div className="border-b border-zinc-100 pb-3 mb-6">
                    <h2 className="text-base font-bold text-zinc-950">Active Sync Platforms</h2>
                    <p className="text-xs text-zinc-500 mt-0.5">Directly map structured action cards to industry standard boards.</p>
                  </div>

                  <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                    {[
                      { id: 'notion', label: 'Notion Sync', desc: 'Sync spec docs directly inside company databases.', status: linkedInteg.notion },
                      { id: 'jira', label: 'Jira Software', desc: 'Auto-allocate sprint tasks inside agile boards.', status: linkedInteg.jira },
                      { id: 'slack', label: 'Slack channels', desc: 'Alert channels when decisions or risks occur.', status: linkedInteg.slack },
                      { id: 'linear', label: 'Linear App', desc: 'Create technical task backlogs directly.', status: linkedInteg.linear }
                    ].map((integ) => (
                      <div key={integ.id} className="border border-zinc-200 rounded-xl p-4 flex flex-col justify-between hover:border-zinc-300 transition">
                        <div>
                          <h4 className="text-xs font-bold text-zinc-950">{integ.label}</h4>
                          <p className="text-[10px] text-zinc-450 mt-1 leading-relaxed">{integ.desc}</p>
                        </div>
                        <div className="mt-4 flex items-center justify-between border-t border-zinc-50 pt-3">
                          <span className={`text-[10px] font-bold ${integ.status ? 'text-emerald-600' : 'text-zinc-400'}`}>
                            {integ.status ? '● Connected' : '○ Disabled'}
                          </span>
                          <button
                            onClick={() => toggleInteg(integ.id as any)}
                            className={`rounded px-2.5 py-1 text-[10px] font-bold border transition ${
                              integ.status 
                                ? 'bg-zinc-100 text-zinc-650 hover:bg-zinc-200' 
                                : 'bg-indigo-600 text-white hover:bg-indigo-500'
                            }`}
                          >
                            {integ.status ? 'Disconnect' : 'Connect'}
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Billing Tab */}
              {activeTab === 'billing' && (
                <div className="space-y-6">
                  <div className="border-b border-zinc-100 pb-3 mb-6">
                    <h2 className="text-base font-bold text-zinc-950">Billing Details</h2>
                    <p className="text-xs text-zinc-500 mt-0.5">Manage active plan tier billing, invoicing logs, and usage quotas.</p>
                  </div>

                  <div className="rounded-xl border border-indigo-150 bg-indigo-50/10 p-5 flex items-start justify-between">
                    <div>
                      <span className="rounded bg-indigo-50 border border-indigo-100 px-2 py-0.5 text-[9px] font-bold text-indigo-600 uppercase">
                        Active tier plan
                      </span>
                      <h3 className="font-extrabold text-zinc-950 text-base mt-2">Enterprise Trial (Pro)</h3>
                      <p className="text-[10px] text-zinc-500 mt-0.5">Enables deep Ollama models and structural Jira mapping APIs.</p>
                    </div>
                    <button className="rounded-lg bg-indigo-600 px-4 py-2 text-xs font-semibold text-white hover:bg-indigo-500 transition shadow-md shadow-indigo-600/20">
                      Manage Subscription
                    </button>
                  </div>

                  {/* Usage Quota */}
                  <div className="space-y-2">
                    <div className="flex justify-between text-xs text-zinc-450 font-semibold">
                      <span>Monthly Whisper transcription hours quota</span>
                      <span className="text-zinc-800">4.5 hrs / 25 hrs</span>
                    </div>
                    <div className="h-1.5 w-full bg-zinc-100 rounded-full overflow-hidden">
                      <div className="h-full bg-indigo-600 rounded-full" style={{ width: '18%' }} />
                    </div>
                  </div>
                </div>
              )}

              {/* Security Tab */}
              {activeTab === 'security' && (
                <div className="space-y-8">
                  {/* API Keys */}
                  <div className="space-y-4">
                    <h3 className="text-xs font-bold text-zinc-800 uppercase">API Tokens Generation</h3>
                    
                    <form onSubmit={handleGenerateKey} className="flex gap-2">
                      <input
                        type="text"
                        required
                        placeholder="e.g. n8n integration client"
                        value={newKeyName}
                        onChange={(e) => setNewKeyName(e.target.value)}
                        className="flex-1 rounded-lg border border-zinc-250 bg-white px-3 py-2 text-xs focus:border-indigo-500 focus:outline-none"
                      />
                      <button type="submit" className="rounded-lg bg-indigo-600 px-4 py-2 text-xs font-semibold text-white hover:bg-indigo-500">
                        Generate Key
                      </button>
                    </form>

                    <div className="space-y-3">
                      {apiKeys.map((keyObj) => (
                        <div key={keyObj.id} className="flex items-center justify-between text-xs border border-zinc-150 rounded-xl p-4 bg-zinc-50/50">
                          <div>
                            <span className="font-bold text-zinc-850 block">{keyObj.name}</span>
                            <span className="font-mono text-[10px] text-zinc-400 mt-1 block">{keyObj.key}</span>
                          </div>
                          
                          <button
                            onClick={() => {
                              setApiKeys(prev => prev.filter(k => k.id !== keyObj.id))
                              toast.success("API key successfully revoked.")
                            }}
                            className="rounded p-1.5 hover:bg-rose-50 text-rose-500 hover:border-rose-100 border border-transparent transition"
                          >
                            <Trash2 className="h-4 w-4" />
                          </button>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              )}

            </motion.div>
          </AnimatePresence>
        </div>

      </div>

    </div>
  )
}
