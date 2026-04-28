'use client'

import { useState, useEffect } from 'react'
import { X, Plus, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { createJob, updateJob, type Job, type CreateJobRequest } from '@/lib/api'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  job?: Job | null
  onJobSaved: () => void
}

const METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE']
const CRON_PRESETS = [
  { label: 'Every minute', value: '* * * * *' },
  { label: 'Hourly',       value: '0 * * * *' },
  { label: 'Daily',        value: '0 0 * * *' },
  { label: 'Weekly',       value: '0 0 * * 0' },
]

const EMPTY: CreateJobRequest = {
  name: '', description: '', schedule: '',
  endpoint: '', method: 'GET',
  headers: {}, body: '',
}

type HeaderRow = { key: string; value: string }

function headersToRows(h: Record<string, string>): HeaderRow[] {
  const rows = Object.entries(h).map(([key, value]) => ({ key, value }))
  return rows.length ? rows : []
}

function rowsToHeaders(rows: HeaderRow[]): Record<string, string> {
  return Object.fromEntries(rows.filter(r => r.key.trim()).map(r => [r.key.trim(), r.value]))
}

function Field({ label, hint, children }: { label: string; hint?: string; children: React.ReactNode }) {
  return (
    <div className="space-y-1.5">
      <div className="flex items-baseline gap-2">
        <label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">{label}</label>
        {hint && <span className="text-xs text-muted-foreground normal-case font-normal">{hint}</span>}
      </div>
      {children}
    </div>
  )
}

export function CreateJobDialog({ open, onOpenChange, job, onJobSaved }: Props) {
  const [form, setForm] = useState<CreateJobRequest>(EMPTY)
  const [headerRows, setHeaderRows] = useState<HeaderRow[]>([])
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!open) return
    if (job) {
      setForm({ name: job.name, description: job.description ?? '', schedule: job.schedule,
                endpoint: job.endpoint, method: job.method, headers: job.headers ?? {}, body: job.body ?? '' })
      setHeaderRows(headersToRows(job.headers ?? {}))
    } else {
      setForm(EMPTY)
      setHeaderRows([])
    }
    setError(null)
  }, [job, open])

  const set = (k: keyof CreateJobRequest, v: string) => setForm(f => ({ ...f, [k]: v }))

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const payload = { ...form, headers: rowsToHeaders(headerRows) }
      job?.id ? await updateJob(job.id, payload) : await createJob(payload)
      onJobSaved()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong')
    } finally {
      setLoading(false)
    }
  }

  if (!open) return null

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4 animate-fade-in"
      style={{ backgroundColor: 'rgba(0,0,0,0.5)', backdropFilter: 'blur(4px)' }}
      onClick={e => { if (e.target === e.currentTarget) onOpenChange(false) }}
    >
      <div className="w-full max-w-xl bg-card border border-border rounded-lg shadow-2xl animate-slide-up flex flex-col max-h-[90vh]">
        {/* Header */}
        <div className="flex items-start justify-between px-6 pt-5 pb-4 border-b border-border shrink-0">
          <div>
            <h2 className="text-base font-semibold tracking-tight text-card-foreground">
              {job ? 'Edit job' : 'New job'}
            </h2>
            <p className="text-xs text-muted-foreground mt-0.5">
              {job ? 'Update job configuration' : 'Configure a new scheduled HTTP job'}
            </p>
          </div>
          <button
            onClick={() => onOpenChange(false)}
            className="p-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Scrollable body */}
        <form onSubmit={handleSubmit} className="overflow-y-auto flex-1">
          <div className="px-6 py-5 space-y-5">
            {error && (
              <div className="px-4 py-3 bg-destructive/10 border border-destructive/20 text-destructive rounded-md text-xs">
                {error}
              </div>
            )}

            {/* Name + Description */}
            <div className="grid grid-cols-2 gap-4">
              <Field label="Name">
                <Input required placeholder="Daily backup" value={form.name}
                  onChange={e => set('name', e.target.value)} disabled={loading} />
              </Field>
              <Field label="Description" hint="(optional)">
                <Input placeholder="What does this do?" value={form.description}
                  onChange={e => set('description', e.target.value)} disabled={loading} />
              </Field>
            </div>

            {/* Schedule */}
            <Field label="Cron Schedule">
              <div className="flex flex-wrap gap-1.5 mb-2">
                {CRON_PRESETS.map(p => (
                  <button key={p.value} type="button"
                    onClick={() => set('schedule', p.value)}
                    className={`text-[11px] px-2.5 py-1 rounded-full border transition-all ${
                      form.schedule === p.value
                        ? 'bg-primary text-primary-foreground border-primary'
                        : 'bg-card text-muted-foreground border-border hover:border-foreground/40 hover:text-foreground'
                    }`}
                  >
                    {p.label}
                  </button>
                ))}
              </div>
              <Input required placeholder="0 0 * * *" value={form.schedule}
                onChange={e => set('schedule', e.target.value)}
                disabled={loading} className="font-mono text-xs" />
              <p className="text-[11px] text-muted-foreground mt-1">minute · hour · day · month · weekday</p>
            </Field>

            <div className="border-t border-border pt-5 space-y-5">
              {/* Method + Endpoint */}
              <div className="flex gap-3">
                <Field label="Method">
                  <div className="flex gap-1.5 flex-wrap">
                    {METHODS.map(m => (
                      <button key={m} type="button"
                        onClick={() => set('method', m)}
                        className={`text-xs px-3 py-1.5 rounded-md border font-mono font-medium transition-all ${
                          form.method === m
                            ? 'bg-primary text-primary-foreground border-primary'
                            : 'bg-card text-muted-foreground border-border hover:border-foreground/40 hover:text-foreground'
                        }`}
                      >
                        {m}
                      </button>
                    ))}
                  </div>
                </Field>
              </div>

              <Field label="Endpoint URL">
                <Input required type="url" placeholder="https://api.example.com/endpoint"
                  value={form.endpoint} onChange={e => set('endpoint', e.target.value)}
                  disabled={loading} className="font-mono text-xs" />
              </Field>

              {/* Headers */}
              <Field label="Headers" hint="(optional)">
                <div className="space-y-2">
                  {headerRows.map((row, i) => (
                    <div key={i} className="flex gap-2 items-center">
                      <Input
                        placeholder="Key"
                        value={row.key}
                        onChange={e => setHeaderRows(rows => rows.map((r, j) => j === i ? { ...r, key: e.target.value } : r))}
                        className="font-mono text-xs"
                        disabled={loading}
                      />
                      <Input
                        placeholder="Value"
                        value={row.value}
                        onChange={e => setHeaderRows(rows => rows.map((r, j) => j === i ? { ...r, value: e.target.value } : r))}
                        className="font-mono text-xs"
                        disabled={loading}
                      />
                      <button type="button"
                        onClick={() => setHeaderRows(rows => rows.filter((_, j) => j !== i))}
                        className="p-1.5 text-muted-foreground hover:text-destructive hover:bg-destructive/10 rounded-md transition-colors shrink-0"
                      >
                        <Trash2 className="w-3.5 h-3.5" />
                      </button>
                    </div>
                  ))}
                  <button type="button"
                    onClick={() => setHeaderRows(rows => [...rows, { key: '', value: '' }])}
                    className="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
                  >
                    <Plus className="w-3.5 h-3.5" /> Add header
                  </button>
                </div>
              </Field>

              {/* Body — only for non-GET */}
              {form.method !== 'GET' && form.method !== 'DELETE' && (
                <Field label="Request Body" hint="(optional, JSON)">
                  <textarea
                    placeholder='{"key": "value"}'
                    value={form.body}
                    onChange={e => set('body', e.target.value)}
                    disabled={loading}
                    rows={4}
                    className="w-full rounded-md border border-input bg-card px-3 py-2 text-xs font-mono text-card-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring focus-visible:border-foreground/40 transition-all resize-none disabled:opacity-50"
                  />
                </Field>
              )}
            </div>
          </div>

          {/* Footer */}
          <div className="flex gap-2 justify-end px-6 py-4 border-t border-border bg-muted/30 shrink-0">
            <Button type="button" variant="outline" size="sm"
              onClick={() => onOpenChange(false)} disabled={loading}>
              Cancel
            </Button>
            <Button type="submit" size="sm" disabled={loading}>
              {loading ? 'Saving…' : job ? 'Save changes' : 'Create job'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
