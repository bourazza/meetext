import api from '@/lib/api'
import type { Workspace, WorkspaceMember } from '@/types'

export async function getWorkspaces(): Promise<Workspace[]> {
  const { data } = await api.get<{ success: boolean; data: Workspace[] }>('/workspaces')
  return data.data
}

export async function getWorkspace(id: string): Promise<Workspace> {
  const { data } = await api.get<{ success: boolean; data: Workspace }>(`/workspaces/${id}`)
  return data.data
}

export async function updateWorkspace(id: string, name: string): Promise<Workspace> {
  const { data } = await api.patch<{ success: boolean; data: Workspace }>(`/workspaces/${id}`, { name })
  return data.data
}

export async function getMembers(workspaceId: string): Promise<WorkspaceMember[]> {
  const { data } = await api.get<{ success: boolean; data: WorkspaceMember[] }>(
    `/workspaces/${workspaceId}/members`
  )
  return data.data
}

export async function removeMember(workspaceId: string, userId: string): Promise<void> {
  await api.delete(`/workspaces/${workspaceId}/members/${userId}`)
}
