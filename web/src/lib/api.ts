const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:5000/api'

export interface Job {
  id: string
  name: string
  description?: string
  schedule: string
  endpoint: string
  method: string
  headers: Record<string, string>
  body?: string
  createdAt?: string
  updatedAt?: string
}

export interface Execution {
  id: string
  jobID: string
  status: string
  startedAt?: string
  finishedAt?: string
}

export interface CreateJobRequest {
  name: string
  description?: string
  schedule: string
  endpoint: string
  method: string
  headers: Record<string, string>
  body?: string
}

export async function getJobs(): Promise<Job[]> {
  const res = await fetch(`${API_URL}/job/all`)
  if (!res.ok) throw new Error('Failed to fetch jobs')
  return res.json()
}

export async function createJob(job: CreateJobRequest): Promise<Job> {
  const res = await fetch(`${API_URL}/job/create`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(job),
  })
  if (!res.ok) throw new Error('Failed to create job')
  return res.json()
}

export async function updateJob(id: string, job: CreateJobRequest): Promise<Job> {
  const res = await fetch(`${API_URL}/job/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(job),
  })
  if (!res.ok) throw new Error('Failed to update job')
  return res.json()
}

export async function deleteJob(id: string): Promise<void> {
  const res = await fetch(`${API_URL}/job/${id}`, { method: 'DELETE' })
  if (!res.ok) throw new Error('Failed to delete job')
}

export async function getExecutions(): Promise<Execution[]> {
  const res = await fetch(`${API_URL}/execution/all`)
  if (!res.ok) throw new Error('Failed to fetch executions')
  return res.json()
}
