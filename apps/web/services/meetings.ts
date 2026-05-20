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

export async function uploadMeeting(input: UploadMeetingInput): Promise<{ meeting: Meeting; analysis?: any }> {
  const form = new FormData()
  form.append('file', input.file)
  if (input.title) form.append('title', input.title)
  if (input.projectId) form.append('project_id', input.projectId)
  if (input.clientId) form.append('client_id', input.clientId)

  const { data } = await api.post<{ success: boolean; data: { meeting: Meeting; analysis?: any } | Meeting }>(
    `/workspaces/${input.workspaceId}/meetings`,
    form,
    {
      headers: { 'Content-Type': 'multipart/form-data' },
      timeout: 120000,
      onUploadProgress: (event) => {
        if (!input.onProgress || !event.total) return
        input.onProgress(Math.round((event.loaded / event.total) * 100))
      },
    }
  )

  if (data.data && 'meeting' in data.data) {
    return data.data
  }
  return { meeting: data.data as Meeting }
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
