import api from '@/lib/api'
import type { Meeting } from '@/types'

export interface UploadMeetingInput {
  workspaceId: string
  file: File
  title?: string
  projectId?: string
  clientId?: string
}

export async function uploadMeeting(input: UploadMeetingInput): Promise<Meeting> {
  const form = new FormData()
  form.append('file', input.file)
  if (input.title) form.append('title', input.title)
  if (input.projectId) form.append('project_id', input.projectId)
  if (input.clientId) form.append('client_id', input.clientId)

  const { data } = await api.post<{ success: boolean; data: Meeting }>(
    `/workspaces/${input.workspaceId}/meetings`,
    form,
    { headers: { 'Content-Type': 'multipart/form-data' } }
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
