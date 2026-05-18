// ── Auth ──────────────────────────────────────────────────────────────────────

export interface User {
  id: string
  full_name: string
  email: string
  avatar_url: string | null
  plan: 'free' | 'pro' | 'business'
  created_at: string
  updated_at: string
}

export interface AuthResponse {
  user: User
  workspace?: Workspace
  access_token: string
  refresh_token: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
}

// ── Workspace ─────────────────────────────────────────────────────────────────

export interface Workspace {
  id: string
  owner_id: string
  name: string
  logo_url: string | null
  created_at: string
}

export interface WorkspaceMember {
  id: string
  workspace_id: string
  user_id: string
  role: 'owner' | 'admin' | 'member'
  created_at: string
}

// ── Client ────────────────────────────────────────────────────────────────────

export interface Client {
  id: string
  workspace_id: string
  company_name: string
  contact_name: string | null
  contact_email: string | null
  logo_url: string | null
  notes: string | null
  created_at: string
}

// ── Project ───────────────────────────────────────────────────────────────────

export type ProjectStatus = 'planning' | 'active' | 'review' | 'completed'

export interface Project {
  id: string
  workspace_id: string
  client_id: string | null
  name: string
  description: string | null
  status: ProjectStatus
  progress: number
  start_date: string | null
  end_date: string | null
  created_at: string
  updated_at: string
}

// ── Meeting ───────────────────────────────────────────────────────────────────

export type UploadType = 'audio' | 'video' | 'pdf' | 'docx'
export type MeetingStatus = 'uploaded' | 'processing' | 'completed' | 'failed' | 'needs_review'

export interface Meeting {
  id: string
  workspace_id: string
  project_id: string | null
  client_id: string | null
  title: string
  upload_type: UploadType
  original_file_url: string
  transcript: string | null
  ai_summary: string | null
  duration_seconds: number | null
  language: string | null
  status: MeetingStatus
  processing_started_at: string | null
  processing_completed_at: string | null
  uploaded_by: string | null
  created_at: string
}

// ── Task ──────────────────────────────────────────────────────────────────────

export type TaskStatus = 'todo' | 'in_progress' | 'review' | 'done'
export type TaskPriority = 'low' | 'medium' | 'high' | 'urgent'

export interface Task {
  id: string
  workspace_id: string
  project_id: string
  meeting_id: string | null
  title: string
  description: string | null
  status: TaskStatus
  priority: TaskPriority
  due_date: string | null
  ai_generated: boolean
  ai_confidence: number | null
  assigned_to: string | null
  created_at: string
  updated_at: string
}

// ── Document ──────────────────────────────────────────────────────────────────

export type DocumentType =
  | 'summary'
  | 'requirements'
  | 'technical_doc'
  | 'sprint_plan'
  | 'client_notes'
  | 'decision_log'

export interface Document {
  id: string
  workspace_id: string
  project_id: string
  meeting_id: string | null
  title: string
  type: DocumentType
  content: string | null
  generated_by_ai: boolean
  created_at: string
  updated_at: string
}

// ── API Response envelope ─────────────────────────────────────────────────────

export interface ApiResponse<T> {
  success: boolean
  data: T
  message?: string
}

export interface ApiError {
  success: false
  error: {
    code: string
    message: string
    fields?: Record<string, string>
  }
}
