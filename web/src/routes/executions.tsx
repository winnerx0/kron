import { useEffect, useState } from 'react'
import { Activity } from 'lucide-react'
import { getExecutions, type Execution } from '@/lib/api'

const STATUS_CFG: Record<string, { dot: string; text: string; bg: string }> = {
  pending:   { dot: 'bg-amber-400',   text: 'text-amber-600 dark:text-amber-400',  bg: 'bg-amber-50 dark:bg-amber-950/30' },
  running:   { dot: 'bg-blue-500',    text: 'text-blue-600 dark:text-blue-400',    bg: 'bg-blue-50 dark:bg-blue-950/30' },
  success:   { dot: 'bg-emerald-500',  text: 'text-emerald-700',                      bg: 'bg-emerald-50 dark:bg-emerald-950/30' },
  completed: { dot: 'bg-emerald-500', text: 'text-emerald-600 dark:text-emerald-400', bg: 'bg-emerald-50 dark:bg-emerald-950/30' },
  failed:    { dot: 'bg-red-500',     text: 'text-red-600 dark:text-red-400',      bg: 'bg-red-50 dark:bg-red-950/30' },
}

function StatusPill({ status }: { status: string }) {
  const cfg = STATUS_CFG[status.toLowerCase()] ?? STATUS_CFG['pending']
  return (
    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border ${cfg.bg} ${cfg.text}`}
      style={{ borderColor: 'transparent' }}>
      <span className={`w-1.5 h-1.5 rounded-full shrink-0 ${cfg.dot} ${status === 'running' ? 'animate-pulse' : ''}`} />
      <span className="capitalize">{status}</span>
    </span>
  )
}

function fmt(d?: string) {
  if (!d) return '—'
  return new Date(d).toLocaleString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function duration(start?: string, end?: string) {
  if (!start || !end) return '—'
  const ms = new Date(end).getTime() - new Date(start).getTime()
  if (ms < 1000) return `${ms}ms`
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`
  return `${Math.floor(ms / 60000)}m ${Math.floor((ms % 60000) / 1000)}s`
}

export function ExecutionsPage() {
  const [executions, setExecutions] = useState<Execution[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [filter, setFilter] = useState<string>('all')
  const [lastRefresh, setLastRefresh] = useState(Date.now())

  useEffect(() => {
    const load = async () => {
      try {
        setLoading(true)
        setExecutions((await getExecutions()) || [])
        setError(null)
        setLastRefresh(Date.now())
      } catch (e) {
        setError(e instanceof Error ? e.message : 'Failed to load')
      } finally {
        setLoading(false)
      }
    }
    load()
    const t = setInterval(load, 5000)
    return () => clearInterval(t)
  }, [])

  const statuses = ['all', 'running', 'completed', 'success', 'failed', 'pending']
  const filtered = filter === 'all' ? executions : executions.filter(e => e.status === filter)
  const ago = Math.round((Date.now() - lastRefresh) / 1000)

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-end justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Executions</h1>
          <p className="text-sm text-muted-foreground mt-1">Real-time job execution history</p>
        </div>
        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          <span className="relative flex h-1.5 w-1.5">
            <span className="animate-ping absolute h-full w-full rounded-full bg-emerald-400 opacity-75" />
            <span className="relative rounded-full h-1.5 w-1.5 bg-emerald-500" />
          </span>
          {ago < 5 ? 'Live' : `${ago}s ago`}
        </div>
      </div>

      {/* Filters */}
      <div className="flex items-center gap-1.5 flex-wrap">
        {statuses.map(s => (
          <button
            key={s}
            onClick={() => setFilter(s)}
            className={`px-3 py-1.5 rounded-full text-xs font-medium border transition-colors capitalize ${
              filter === s
                ? 'bg-primary text-primary-foreground border-primary'
                : 'bg-card text-muted-foreground border-border hover:text-foreground hover:border-foreground/40'
            }`}
          >
            {s}
          </button>
        ))}
        {!loading && (
          <span className="ml-auto text-xs text-muted-foreground">{filtered.length} record{filtered.length !== 1 ? 's' : ''}</span>
        )}
      </div>

      {error && (
        <div className="px-4 py-3 bg-destructive/10 border border-destructive/20 text-destructive rounded-lg text-sm">{error}</div>
      )}

      <div className="bg-card border border-border rounded-lg overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border bg-muted/40">
              <th className="px-5 py-3 text-left text-[10px] font-semibold uppercase tracking-widest text-muted-foreground">Status</th>
              <th className="px-5 py-3 text-left text-[10px] font-semibold uppercase tracking-widest text-muted-foreground">Job ID</th>
              <th className="px-5 py-3 text-left text-[10px] font-semibold uppercase tracking-widest text-muted-foreground">Started</th>
              <th className="px-5 py-3 text-left text-[10px] font-semibold uppercase tracking-widest text-muted-foreground">Completed</th>
              <th className="px-5 py-3 text-right text-[10px] font-semibold uppercase tracking-widest text-muted-foreground">Duration</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              [...Array(5)].map((_, i) => (
                <tr key={i} className="border-b border-border last:border-0">
                  <td className="px-5 py-3.5"><div className="skeleton h-6 w-24 rounded-full" /></td>
                  <td className="px-5 py-3.5"><div className="skeleton h-3.5 w-20" /></td>
                  <td className="px-5 py-3.5"><div className="skeleton h-3.5 w-36" /></td>
                  <td className="px-5 py-3.5"><div className="skeleton h-3.5 w-36" /></td>
                  <td className="px-5 py-3.5"><div className="skeleton h-3.5 w-12 ml-auto" /></td>
                </tr>
              ))
            ) : filtered.length === 0 ? (
              <tr>
                <td colSpan={5}>
                  <div className="flex flex-col items-center justify-center py-16 text-center">
                    <Activity className="w-10 h-10 text-muted-foreground/30 mb-4" strokeWidth={1} />
                    <p className="text-sm font-medium mb-1">No executions</p>
                    <p className="text-sm text-muted-foreground">
                      {filter !== 'all' ? `No ${filter} executions found` : 'Jobs will appear here once they run'}
                    </p>
                  </div>
                </td>
              </tr>
            ) : filtered.map((ex, idx) => (
              <tr
                key={ex.id}
                className="border-b border-border last:border-0 hover:bg-muted/30 transition-colors animate-slide-up"
                style={{ animationDelay: `${idx * 25}ms` }}
              >
                <td className="px-5 py-3"><StatusPill status={ex.status} /></td>
                <td className="px-5 py-3 font-mono text-xs text-muted-foreground">{ex.jobID.slice(0, 8)}</td>
                <td className="px-5 py-3 text-sm text-muted-foreground">{fmt(ex.startedAt)}</td>
                <td className="px-5 py-3 text-sm text-muted-foreground">{fmt(ex.finishedAt)}</td>
                <td className="px-5 py-3 text-right font-mono text-xs text-muted-foreground">{ex.finishedAt !== "0001-01-01T00:13:35+00:13" ? duration(ex.startedAt, ex.finishedAt) : undefined}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
