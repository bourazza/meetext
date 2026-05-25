import api, { pollingApi } from '@/lib/api'
import type { Meeting } from '@/types'

const DEBUG_POLLING = process.env.NODE_ENV === 'development'

function debugLog(message: string, data?: any) {
  if (DEBUG_POLLING) {
    console.log(`[Meetings Service] ${message}`, data || '')
  }
}

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
  // Full structured AI output — populated when status is 'completed'.
  // Parsed JSON string containing tasks, decisions, risks, tickets, etc.
  ai_result?: string
}

export interface PollStatusOptions {
  maxAttempts?: number
  pollInterval?: number
  onProgress?: (status: MeetingStatus, attempt: number) => void
  onError?: (error: any, attempt: number) => void
}

export async function uploadMeeting(input: UploadMeetingInput): Promise<{ meeting: Meeting }> {
  const form = new FormData()
  form.append('file', input.file)
  if (input.title) form.append('title', input.title)
  if (input.projectId) form.append('project_id', input.projectId)
  if (input.clientId) form.append('client_id', input.clientId)

  debugLog('Uploading meeting', {
    workspaceId: input.workspaceId,
    fileName: input.file.name,
    fileSize: input.file.size,
  })

  // IMPORTANT: Explicitly delete the instance-level 'Content-Type: application/json' header
  // for this multipart/form-data request. If the default header is not cleared, Axios can
  // fail to override it with the correct 'multipart/form-data; boundary=...' value,
  // causing Go's ParseMultipartForm to receive an application/json content-type and return 400.
  const { data } = await api.post<{ success: boolean; data: { meeting: Meeting } }>(
    `/workspaces/${input.workspaceId}/meetings`,
    form,
    {
      timeout: 60000,
      headers: {
        'Content-Type': undefined, // Let the browser/Axios set multipart/form-data + boundary automatically
      },
      onUploadProgress: (event) => {
        if (!input.onProgress || !event.total) return
        input.onProgress(Math.round((event.loaded / event.total) * 100))
      },
    }
  )

  debugLog('Upload successful', { meetingId: data.data.meeting.id })
  return data.data
}

export async function getMeetingStatus(workspaceId: string, meetingId: string): Promise<MeetingStatus> {
  debugLog('Fetching meeting status', { workspaceId, meetingId })

  const { data } = await pollingApi.get<{ success: boolean; data: MeetingStatus }>(
    `/workspaces/${workspaceId}/meetings/${meetingId}/status`
  )

  debugLog('Status fetched', data.data)
  return data.data
}

/**
 * Poll meeting status until completed or failed with exponential backoff
 */
export async function pollMeetingStatus(
  workspaceId: string,
  meetingId: string,
  options: PollStatusOptions = {}
): Promise<MeetingStatus> {
  const {
    maxAttempts = 180, // 180 attempts * 5s = 15 minutes max
    pollInterval = 5000, // 5 seconds base interval
    onProgress,
    onError,
  } = options

  debugLog('Starting status polling', {
    workspaceId,
    meetingId,
    maxAttempts,
    pollInterval,
  })

  let attempt = 0
  let consecutiveErrors = 0
  const maxConsecutiveErrors = 5

  while (attempt < maxAttempts) {
    attempt++

    try {
      const status = await getMeetingStatus(workspaceId, meetingId)

      // Reset error counter on success
      consecutiveErrors = 0

      debugLog(`Poll attempt ${attempt}/${maxAttempts}`, {
        status: status.status,
        hasSummary: !!status.ai_summary,
      })

      if (onProgress) {
        onProgress(status, attempt)
      }

      // Terminal states
      if (status.status === 'completed' || status.status === 'failed') {
        debugLog('Polling complete', { finalStatus: status.status, attempts: attempt })
        return status
      }

      // Calculate next poll interval with exponential backoff
      const backoffMultiplier = Math.min(Math.floor(attempt / 10), 3) // Increase every 10 attempts, max 3x
      const nextInterval = pollInterval * (1 + backoffMultiplier * 0.5)

      debugLog(`Waiting ${nextInterval}ms before next poll`)
      await new Promise(resolve => setTimeout(resolve, nextInterval))

    } catch (error: any) {
      consecutiveErrors++

      debugLog(`Poll error (${consecutiveErrors}/${maxConsecutiveErrors})`, {
        attempt,
        error: error.message,
        status: error.response?.status,
      })

      if (onError) {
        onError(error, attempt)
      }

      // If too many consecutive errors, give up
      if (consecutiveErrors >= maxConsecutiveErrors) {
        debugLog('Too many consecutive errors, stopping poll')
        throw new Error(`Failed to poll status after ${consecutiveErrors} consecutive errors`)
      }

      // If 401, let the interceptor handle refresh, then retry immediately
      if (error.response?.status === 401) {
        debugLog('401 error, auth interceptor will handle refresh')
        await new Promise(resolve => setTimeout(resolve, 2000))
        continue
      }

      // For other errors, wait with exponential backoff
      const errorBackoff = Math.min(2000 * Math.pow(2, consecutiveErrors - 1), 10000)
      debugLog(`Waiting ${errorBackoff}ms after error`)
      await new Promise(resolve => setTimeout(resolve, errorBackoff))
    }
  }

  debugLog('Max attempts reached', { attempts: attempt })
  throw new Error(`Polling timeout: max attempts (${maxAttempts}) reached`)
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
