import { useEffect, useState } from 'react'
import { Link } from '@tanstack/react-router'
import { Calendar, Activity, CheckCircle, XCircle, ArrowRight, Clock } from 'lucide-react'
import { getJobs, getExecutions, type Job, type Execution } from '@/lib/api'

function StatCard({
  label, value, icon: Icon, sub, delay = 0
}: {
  label: string; value: string | number; icon: React.ElementType; sub?: string; delay?: number
}) {
  return (
    <div
      className="bg-card border border-border rounded-lg p-5 flex flex-col gap-3 animate-slide-up"
      style={{ animationDelay: `${delay}ms` }}
    >
      <div className="flex items-center justify-between">
        <span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">{label}</span>
        <div className="p-1.5 rounded-md bg-muted">
          <Icon className="w-3.5 h-3.5 text-muted-foreground" />
        </div>
      </div>
      <div>
        <p className="text-3xl font-semibold tracking-tight leading-none">{value}</p>
        {sub && <p className="text-xs text-muted-foreground mt-1.5">{sub}</p>}
      </div>
    </div>
  )
}

const STATUS_DOT: Record<string, string> = {
  success:   'bg-emerald-500',
  completed: 'bg-emerald-500',
  failed:    'bg-red-500',
  running:   'bg-blue-500',
  pending:   'bg-amber-400',
}

function formatAgo(dateStr?: string) {
  if (!dateStr) return '—'
  const diff = Date.now() - new Date(dateStr).getTime()
  const m = Math.floor(diff / 60000)
  if (m < 1) return 'just now'
  if (m < 60) return `${m}m ago`
  const h = Math.floor(m / 60)
  if (h < 24) return `${h}h ago`
  return `${Math.floor(h / 24)}d ago`
}

export function DashboardPage() {
  const [jobs, setJobs] = useState<Job[]>([])
  const [executions, setExecutions] = useState<Execution[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([getJobs(), getExecutions()])
      .then(([j, e]) => { setJobs(j || []); setExecutions(e || []) })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  const total = jobs.length
  const successful = executions.filter(e => e.status === 'success' || e.status === 'completed').length
  const failed = executions.filter(e => e.status === 'failed').length
  const rate = executions.length > 0 ? Math.round((successful / executions.length) * 100) : 0
  const recent = [...executions].sort((a, b) =>
    new Date(b.startedAt ?? 0).getTime() - new Date(a.startedAt ?? 0).getTime()
  ).slice(0, 6)

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Dashboard</h1>
        <p className="text-sm text-muted-foreground mt-1">Overview of your job scheduler</p>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        {loading ? (
          [...Array(4)].map((_, i) => (
            <div key={i} className="bg-card border border-border rounded-lg p-5 h-28 skeleton" />
          ))
        ) : (
          <>
            <StatCard label="Total Jobs"    value={total}      icon={Calendar}     sub={`${total} scheduled`}      delay={0} />
            <StatCard label="Executions"    value={executions.length} icon={Activity} sub="all time"              delay={50} />
            <StatCard label="Successful"    value={successful} icon={CheckCircle}  sub={`${rate}% success rate`}  delay={100} />
            <StatCard label="Failed"        value={failed}     icon={XCircle}      sub={failed ? 'needs attention' : 'all clear'} delay={150} />
          </>
        )}
      </div>

      {/* Two-column: recent executions + jobs */}
      <div className="grid grid-cols-1 lg:grid-cols-5 gap-6">
        {/* Recent executions */}
        <div className="lg:col-span-3 bg-card border border-border rounded-lg animate-slide-up" style={{ animationDelay: '200ms' }}>
          <div className="flex items-center justify-between px-5 py-4 border-b border-border">
            <h2 className="text-sm font-semibold">Recent Executions</h2>
            <Link to="/executions" className="text-xs text-muted-foreground hover:text-foreground flex items-center gap-1 transition-colors">
              View all <ArrowRight className="w-3 h-3" />
            </Link>
          </div>
          <div className="divide-y divide-border">
            {loading ? (
              [...Array(4)].map((_, i) => (
                <div key={i} className="px-5 py-3.5 flex items-center gap-3">
                  <div className="skeleton w-2 h-2 rounded-full" />
                  <div className="skeleton h-3 w-24" />
                  <div className="skeleton h-3 w-20 ml-auto" />
                </div>
              ))
            ) : recent.length === 0 ? (
              <div className="px-5 py-8 text-center text-sm text-muted-foreground">No executions yet</div>
            ) : recent.map((ex, i) => {
              const dotColor = STATUS_DOT[ex.status] ?? 'bg-muted-foreground'
              return (
                <div key={ex.id} className="flex items-center gap-3 px-5 py-3.5 text-sm hover:bg-muted/30 transition-colors animate-slide-up" style={{ animationDelay: `${220 + i * 30}ms` }}>
                  <span className={`w-2 h-2 rounded-full shrink-0 ${dotColor} ${ex.status === 'running' ? 'animate-pulse' : ''}`} />
                  <span className="font-mono text-xs text-muted-foreground">{ex.jobID.slice(0, 8)}</span>
                  <span className="capitalize text-xs text-muted-foreground">{ex.status}</span>
                  <span className="ml-auto text-xs text-muted-foreground">{formatAgo(ex.startedAt)}</span>
                </div>
              )
            })}
          </div>
        </div>

        {/* Jobs list */}
        <div className="lg:col-span-2 bg-card border border-border rounded-lg animate-slide-up" style={{ animationDelay: '250ms' }}>
          <div className="flex items-center justify-between px-5 py-4 border-b border-border">
            <h2 className="text-sm font-semibold">Jobs</h2>
            <Link to="/jobs" className="text-xs text-muted-foreground hover:text-foreground flex items-center gap-1 transition-colors">
              Manage <ArrowRight className="w-3 h-3" />
            </Link>
          </div>
          <div className="divide-y divide-border">
            {loading ? (
              [...Array(3)].map((_, i) => (
                <div key={i} className="px-5 py-3.5 space-y-1.5">
                  <div className="skeleton h-3 w-28" />
                  <div className="skeleton h-3 w-20" />
                </div>
              ))
            ) : jobs.length === 0 ? (
              <div className="px-5 py-8 text-center">
                <p className="text-sm text-muted-foreground mb-3">No jobs yet</p>
                <Link to="/jobs" className="text-xs underline text-muted-foreground hover:text-foreground">Create one →</Link>
              </div>
            ) : jobs.map((job, i) => (
              <div key={job.id} className="px-5 py-3.5 animate-slide-up" style={{ animationDelay: `${270 + i * 30}ms` }}>
                <p className="text-sm font-medium truncate">{job.name}</p>
                <p className="font-mono text-xs text-muted-foreground mt-0.5">{job.schedule}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
