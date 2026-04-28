'use client'

import { useEffect, useState } from 'react'
import { Trash2, Pencil, Plus, CalendarClock } from 'lucide-react'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { getJobs, deleteJob, type Job } from '@/lib/api'
import { CreateJobDialog } from './create-job-dialog'

function SkeletonRow() {
  return (
    <div className="flex items-center gap-6 px-6 py-4 border-b border-border last:border-0">
      <div className="skeleton h-4 w-36" />
      <div className="skeleton h-4 w-56 flex-1" />
      <div className="skeleton h-6 w-28 rounded-full" />
      <div className="skeleton h-7 w-16 rounded-md" />
    </div>
  )
}

function EmptyState({ onNew }: { onNew: () => void }) {
  return (
    <div className="flex flex-col items-center justify-center py-16 px-6 text-center animate-fade-in">
      <div className="mb-5 text-muted-foreground/30">
        <CalendarClock strokeWidth={1} className="w-14 h-14 mx-auto" />
      </div>
      <p className="text-sm font-medium text-foreground mb-1">No jobs scheduled</p>
      <p className="text-sm text-muted-foreground mb-6 max-w-xs">
        Create your first job to start scheduling recurring work.
      </p>
      <Button size="sm" onClick={onNew} className="gap-2">
        <Plus className="w-3.5 h-3.5" />
        Create job
      </Button>
    </div>
  )
}

export function JobsList() {
  const [jobs, setJobs] = useState<Job[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedJob, setSelectedJob] = useState<Job | null>(null)
  const [showDialog, setShowDialog] = useState(false)

  const loadJobs = async () => {
    try {
      setLoading(true)
      const data = await getJobs()
      setJobs(data || [])
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load jobs')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { loadJobs() }, [])

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this job?')) return
    try {
      await deleteJob(id)
      setJobs(jobs.filter(j => j.id !== id))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete job')
    }
  }

  const handleJobSaved = () => {
    setShowDialog(false)
    setSelectedJob(null)
    loadJobs()
  }

  return (
    <>
      <Card>
        <CardHeader className="px-6 py-4 flex flex-row items-center justify-between space-y-0 border-b border-border">
          <div>
            <h2 className="text-base font-semibold tracking-tight">Jobs</h2>
            <p className="text-xs text-muted-foreground mt-0.5">Scheduled recurring tasks</p>
          </div>
          <Button
            size="sm"
            onClick={() => { setSelectedJob(null); setShowDialog(true) }}
            className="gap-1.5 h-8 text-xs px-3"
          >
            <Plus className="w-3.5 h-3.5" />
            New Job
          </Button>
        </CardHeader>

        <CardContent className="p-0">
          {error && (
            <div className="mx-6 mt-4 px-4 py-3 bg-red-50 border border-red-200 text-red-700 rounded-md text-xs">
              {error}
            </div>
          )}

          {loading ? (
            <div className="divide-y divide-border">
              {[...Array(3)].map((_, i) => <SkeletonRow key={i} />)}
            </div>
          ) : jobs.length === 0 ? (
            <EmptyState onNew={() => { setSelectedJob(null); setShowDialog(true) }} />
          ) : (
            <div className="divide-y divide-border">
              {/* Column headers */}
              <div className="flex items-center gap-6 px-6 py-2 bg-muted/40">
                <span className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground w-36 shrink-0">Name</span>
                <span className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground flex-1">Description</span>
                <span className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground w-28 shrink-0">Schedule</span>
                <span className="w-16 shrink-0" />
              </div>

              {jobs.map((job, idx) => (
                <div
                  key={job.id}
                  className="flex items-center gap-6 px-6 py-3.5 hover:bg-muted/30 transition-colors group animate-slide-up"
                  style={{ animationDelay: `${idx * 40}ms` }}
                >
                  <span className="text-sm font-medium w-36 shrink-0 truncate">{job.name}</span>
                  <span className="text-sm text-muted-foreground flex-1 truncate">{job.description || '—'}</span>
                  <span className="mono text-xs bg-muted text-muted-foreground px-2.5 py-1 rounded-full border border-border w-28 shrink-0 truncate block">
                    {job.schedule}
                  </span>
                  <div className="flex gap-1 w-16 shrink-0 justify-end opacity-0 group-hover:opacity-100 transition-opacity">
                    <button
                      onClick={() => { setSelectedJob(job); setShowDialog(true) }}
                      className="p-1.5 rounded-md hover:bg-muted text-muted-foreground hover:text-foreground transition-colors"
                    >
                      <Pencil className="w-3.5 h-3.5" />
                    </button>
                    <button
                      onClick={() => handleDelete(job.id)}
                      className="p-1.5 rounded-md hover:bg-red-50 text-muted-foreground hover:text-destructive transition-colors"
                    >
                      <Trash2 className="w-3.5 h-3.5" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      <CreateJobDialog
        open={showDialog}
        onOpenChange={setShowDialog}
        job={selectedJob}
        onJobSaved={handleJobSaved}
      />
    </>
  )
}
