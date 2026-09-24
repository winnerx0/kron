import { getAccessToken, redirectToLogin, refreshTokens } from "./auth";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:5000/api";

export interface ErrorResponse {
  error: string;
}

export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

export interface Job {
  id: string;
  name: string;
  description?: string;
  schedule: string;
  endpoint: string;
  method: string;
  headers: Record<string, string>;
  body?: string;
  enabled: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface Execution {
  id: string;
  jobID: string;
  jobName: string;
  status: string;
  startedAt?: string;
  finishedAt?: string;
}

export interface ExecutionDetail {
  id: string;
  jobID: string;
  jobName: string;
  endpoint: string;
  method: string;
  status: string;
  startedAt?: string;
  finishedAt?: string;
  responseBody: string;
}

export interface DashboardJob {
  id: string;
  name: string;
  schedule: string;
  enabled: boolean;
}

export interface DashboardExecution {
  id: string;
  jobID: string;
  jobName: string;
  status: string;
  startedAt?: string;
  finishedAt?: string;
}

export interface DashboardData {
  totalJobs: number;
  enabledJobs: number;
  disabledJobs: number;
  totalExecutions: number;
  successfulExecutions: number;
  failedExecutions: number;
  successRate: number;
  recentExecutions: DashboardExecution[];
  jobs: DashboardJob[];
}

export interface PaginatedExecutions {
  items: Execution[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface CreateJobRequest {
  name: string;
  description?: string;
  schedule: string;
  endpoint: string;
  method: string;
  headers: Record<string, string>;
  body?: string;
}

export interface CurrentUser {
  id: string;
  name: string;
  email: string;
  profile_picture: string;
  joined_at: string;
}

async function readError(res: Response, fallback: string): Promise<string> {
  try {
    const data = (await res.json()) as Partial<ErrorResponse>;
    if (data && typeof data.error === "string") return data.error;
  } catch {
    // body wasn't JSON
  }
  return fallback;
}

async function authedFetch(
  path: string,
  init: RequestInit = {},
  retry = true,
): Promise<Response> {
  const headers = new Headers(init.headers);
  const token = getAccessToken();
  if (!token) {
    redirectToLogin();
    throw new ApiError(401, "Unauthorized");
  }
  headers.set("Authorization", `Bearer ${token}`);

  const res = await fetch(`${API_URL}${path}`, { ...init, headers });

  if (res.status === 401 && retry) {
    const refreshed = await refreshTokens();
    if (!refreshed) {
      redirectToLogin();
      throw new ApiError(401, "Unauthorized");
    }
    return authedFetch(path, init, false);
  }

  return res;
}

async function expectOk(res: Response, fallback: string): Promise<void> {
  if (!res.ok) {
    throw new ApiError(res.status, await readError(res, fallback));
  }
}

export async function getJobs(): Promise<Job[]> {
  const res = await authedFetch("/job/all");
  await expectOk(res, "Failed to fetch jobs");
  return res.json();
}

export async function getCurrentUser(): Promise<CurrentUser> {
  const res = await authedFetch("/user/me");
  await expectOk(res, "Failed to fetch user profile");
  return res.json();
}

export async function getDashboard(): Promise<DashboardData> {
  const res = await authedFetch("/dashboard");
  await expectOk(res, "Failed to fetch dashboard");
  return res.json();
}

export async function createJob(job: CreateJobRequest): Promise<Job> {
  const res = await authedFetch("/job/create", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(job),
  });
  await expectOk(res, "Failed to create job");
  return res.json();
}

export async function updateJob(
  id: string,
  job: CreateJobRequest,
): Promise<Job> {
  const res = await authedFetch(`/job/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(job),
  });
  await expectOk(res, "Failed to update job");
  return res.json();
}

export async function deleteJob(id: string): Promise<void> {
  const res = await authedFetch(`/job/${id}`, { method: "DELETE" });
  await expectOk(res, "Failed to delete job");
}

export async function runJob(id: string): Promise<void> {
  const res = await authedFetch(`/job/${id}/run`, { method: "POST" });
  await expectOk(res, "Failed to run job");
}

export async function stopJob(id: string): Promise<void> {
  const res = await authedFetch(`/job/${id}/stop`, { method: "POST" });
  await expectOk(res, "Failed to stop job");
}

export async function toggleJobStatus(id: string): Promise<void> {
  const res = await authedFetch(`/job/${id}/status`, { method: "PATCH" });
  await expectOk(res, "Failed to update job status");
}

export async function getExecutionsPage(
  page = 1,
  pageSize = 10,
  jobID?: string,
): Promise<PaginatedExecutions> {
  const params = new URLSearchParams({
    page: String(page),
    pageSize: String(pageSize),
  });
  if (jobID) {
    params.append("jobID", jobID);
  }
  const res = await authedFetch(`/execution/all?${params.toString()}`);
  await expectOk(res, "Failed to fetch executions");
  return res.json();
}

export async function getExecutions(jobID?: string): Promise<Execution[]> {
  const data = await getExecutionsPage(1, 10, jobID);
  return data.items;
}

export async function getExecution(id: string): Promise<ExecutionDetail> {
  const res = await authedFetch(`/execution/${id}`);
  await expectOk(res, "Failed to fetch execution");
  return res.json();
}
