import api from '@/lib/api'
import type { Meeting } from '@/types'

export interface UploadMeetingInput {
  workspaceId: string
  file: File
  title?: string
  projectId?: string
  clientId?: string
  onProgress?: (progress: number) => void
}

export interface MeetingStatus {
  id: string
  status: 'uploaded' | 'processing' | 'completed' | 'failed'
  ai_summary?: string
}

export async function uploadMeeting(input: UploadMeetingInput): Promise<{ meeting: Meeting }> {
  const form = new FormData()
  form.append('file', input.file)
  if (input.title) form.append('title', input.title)
  if (input.projectId) form.append('project_id', input.projectId)
  if (input.clientId) form.append('client_id', input.clientId)

  const { data } = await api.post<{ success: boolean; data: { meeting: Meeting } }>(
    `/workspaces/${input.workspaceId}/meetings`,
    form,
    {
      timeout: 60000,
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (event) => {
        if (!input.onProgress || !event.total) return
        input.onProgress(Math.round((event.loaded / event.total) * 100))
      },
    }
  )
  return data.data
}

export async function getMeetingStatus(workspaceId: string, meetingId: string): Promise<MeetingStatus> {
  const { data } = await api.get<{ success: boolean; data: MeetingStatus }>(
    `/workspaces/${workspaceId}/meetings/${meetingId}/status`
  )
  return data.data
}

export async function getMeetings(
  workspaceId: string,
  limit = 20,
  offset = 0
): Promise<Meeting[]> {
  const { data } = await api.get<{ success: boolean; data: Meeting[] }>(
    `/workspaces/${workspaceId}/meetings`,
    { params: { limit, offset } }
  )
  return data.data ?? []
}

export async function getMeeting(workspaceId: string, meetingId: string): Promise<Meeting> {
  const { data } = await api.get<{ success: boolean; data: Meeting }>(
    `/workspaces/${workspaceId}/meetings/${meetingId}`
  )
  return data.data
}

export async function deleteMeeting(workspaceId: string, meetingId: string): Promise<void> {
  await api.delete(`/workspaces/${workspaceId}/meetings/${meetingId}`)
}
