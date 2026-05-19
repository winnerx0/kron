import { useEffect, useState } from "react";
import {
  Plus,
  Pencil,
  Trash2,
  CalendarClock,
  Play,
  Square,
  Power,
  PowerOff,
} from "lucide-react";
import { usePostHog } from "posthog-js/react";
import {
  getJobs,
  deleteJob,
  runJob,
  stopJob,
  toggleJobStatus,
  type Job,
} from "@/lib/api";
import { CreateJobDialog } from "@/components/create-job-dialog";

const METHOD_COLORS: Record<string, string> = {
  GET: "text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/40 border-emerald-200 dark:border-emerald-900",
  POST: "text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-950/40 border-blue-200 dark:border-blue-900",
  PUT: "text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-950/40 border-amber-200 dark:border-amber-900",
  PATCH:
    "text-purple-600 dark:text-purple-400 bg-purple-50 dark:bg-purple-950/40 border-purple-200 dark:border-purple-900",
  DELETE:
    "text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950/40 border-red-200 dark:border-red-900",
};

function MethodBadge({ method }: { method: string }) {
  return (
    <span
      className={`inline-block font-mono text-[10px] font-semibold px-2 py-0.5 rounded border ${METHOD_COLORS[method] ?? "text-muted-foreground bg-muted border-border"}`}
    >
      {method}
    </span>
  );
}

function StatusBadge({ enabled }: { enabled: boolean }) {
  return (
    <span
      className={`inline-flex items-center gap-1.5 font-mono text-[10px] font-semibold px-2 py-0.5 rounded-full border ${
        enabled
          ? "text-emerald-700 dark:text-emerald-300 bg-emerald-50 dark:bg-emerald-950/40 border-emerald-200 dark:border-emerald-900"
          : "text-muted-foreground bg-muted border-border"
      }`}
    >
      <span
        className={`h-1.5 w-1.5 rounded-full ${
          enabled ? "bg-emerald-500" : "bg-muted-foreground/50"
        }`}
      />
      {enabled ? "Enabled" : "Disabled"}
    </span>
  );
}

function SkeletonRow() {
  return (
    <tr className="border-b border-border">
      <td className="px-5 py-3.5">
        <div className="skeleton h-3.5 w-28" />
      </td>
      <td className="px-5 py-3.5">
        <div className="skeleton h-5 w-12 rounded" />
      </td>
      <td className="px-5 py-3.5">
        <div className="skeleton h-3.5 w-52" />
      </td>
      <td className="px-5 py-3.5">
        <div className="skeleton h-6 w-24 rounded-full" />
      </td>
      <td className="px-5 py-3.5">
        <div className="skeleton h-5 w-20 rounded-full" />
      </td>
      <td className="px-5 py-3.5">
        <div className="skeleton h-3.5 w-10 ml-auto" />
      </td>
    </tr>
  );
}

export function JobsPage() {
  const posthog = usePostHog();
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selected, setSelected] = useState<Job | null>(null);
  const [showDialog, setShowDialog] = useState(false);
  const [busyJobID, setBusyJobID] = useState<string | null>(null);

  const load = async () => {
    try {
      setLoading(true);
      setJobs((await getJobs()) || []);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this job?")) return;
    try {
      await deleteJob(id);
      setJobs((j) => j.filter((x) => x.id !== id));
      posthog.capture('job_deleted', { job_id: id });
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to delete");
    }
  };

  const handleRun = async (id: string) => {
    try {
      setBusyJobID(id);
      await runJob(id);
      posthog.capture('job_run', { job_id: id });
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to run");
    } finally {
      setBusyJobID(null);
    }
  };

  const handleStop = async (id: string) => {
    try {
      setBusyJobID(id);
      await stopJob(id);
      posthog.capture('job_stopped', { job_id: id });
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to stop");
    } finally {
      setBusyJobID(null);
    }
  };

  const handleToggleStatus = async (job: Job) => {
    const currentEnabled = job.enabled !== false;
    const nextEnabled = !currentEnabled;
    try {
      setBusyJobID(job.id);
      setJobs((current) =>
        current.map((item) =>
          item.id === job.id ? { ...item, enabled: nextEnabled } : item,
        ),
      );
      await toggleJobStatus(job.id);
      posthog.capture('job_status_toggled', { job_id: job.id, job_name: job.name, enabled: nextEnabled });
      setError(null);
    } catch (e) {
      setJobs((current) =>
        current.map((item) =>
          item.id === job.id ? { ...item, enabled: currentEnabled } : item,
        ),
      );
      setError(e instanceof Error ? e.message : "Failed to update status");
    } finally {
      setBusyJobID(null);
    }
  };

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-end justify-between">
        <div>
          <h1 className="font-display text-[36px] leading-tight">Jobs</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage your scheduled HTTP tasks
          </p>
        </div>
        <button
          className="inline-flex items-center gap-2 h-9 rounded-full bg-primary px-4 text-sm font-medium text-primary-foreground transition-all hover:opacity-90 active:scale-[0.99]"
          onClick={() => {
            setSelected(null);
            setShowDialog(true);
          }}
        >
          <Plus className="w-3.5 h-3.5" /> New Job
        </button>
      </div>

      {error && (
        <div className="px-4 py-3 bg-destructive/10 border border-destructive/20 text-destructive rounded-lg text-sm">
          {error}
        </div>
      )}

      <div className="bg-card border border-border rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border bg-muted/40">
              {["Name", "Method", "Endpoint", "Schedule", "Status", ""].map(
                (h) => (
                  <th
                    key={h}
                    className={`px-5 py-3 font-mono text-[10px] uppercase tracking-[0.22em] text-muted-foreground ${h === "" ? "text-right" : "text-left"}`}
                  >
                    {h}
                  </th>
                ),
              )}
            </tr>
          </thead>
          <tbody>
            {loading ? (
              [...Array(4)].map((_, i) => <SkeletonRow key={i} />)
            ) : jobs.length === 0 ? (
              <tr>
                <td colSpan={6}>
                  <div className="flex flex-col items-center justify-center py-16 text-center">
                    <CalendarClock
                      className="w-10 h-10 text-muted-foreground/30 mb-4"
                      strokeWidth={1}
                    />
                    <p className="text-sm font-medium mb-1">No jobs yet</p>
                    <p className="text-sm text-muted-foreground mb-5">
                      Create your first scheduled job to get started
                    </p>
                    <button
                      className="inline-flex items-center gap-2 h-9 rounded-full bg-primary px-4 text-sm font-medium text-primary-foreground transition-all hover:opacity-90 active:scale-[0.99]"
                      onClick={() => {
                        setSelected(null);
                        setShowDialog(true);
                      }}
                    >
                      <Plus className="w-3.5 h-3.5" /> Create job
                    </button>
                  </div>
                </td>
              </tr>
            ) : (
              jobs.map((job, idx) => {
                const enabled = job.enabled !== false;
                return (
                  <tr
                    key={job.id}
                    className={`border-b border-border last:border-0 hover:bg-muted/30 transition-colors group animate-slide-up ${
                      enabled ? "" : "bg-muted/10"
                    }`}
                    style={{ animationDelay: `${idx * 35}ms` }}
                  >
                    <td className="px-5 py-3.5 font-medium whitespace-nowrap">
                      {job.name}
                    </td>
                    <td className="px-5 py-3.5 whitespace-nowrap">
                      <MethodBadge method={job.method} />
                    </td>
                    <td className="px-5 py-3.5 font-mono text-xs text-muted-foreground max-w-xs truncate">
                      {job.endpoint}
                    </td>
                    <td className="px-5 py-3.5 whitespace-nowrap">
                      <span className="inline-block font-mono text-xs bg-muted text-muted-foreground px-2.5 py-1 rounded-full border border-border">
                        {job.schedule}
                      </span>
                    </td>
                    <td className="px-5 py-3.5 whitespace-nowrap">
                      <StatusBadge enabled={enabled} />
                    </td>
                    <td className="px-5 py-3.5">
                      <div className="flex items-center justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button
                          onClick={() => handleToggleStatus(job)}
                          disabled={busyJobID === job.id}
                          title={enabled ? "Disable job" : "Enable job"}
                          className={`p-1.5 rounded-md transition-colors disabled:opacity-50 ${
                            enabled
                              ? "text-muted-foreground hover:text-destructive hover:bg-destructive/10"
                              : "text-muted-foreground hover:text-foreground hover:bg-muted"
                          }`}
                        >
                          {enabled ? (
                            <PowerOff className="w-3.5 h-3.5" />
                          ) : (
                            <Power className="w-3.5 h-3.5" />
                          )}
                        </button>
                        <button
                          onClick={() => handleRun(job.id)}
                          disabled={busyJobID === job.id || !enabled}
                          title={enabled ? "Run job" : "Enable job to run"}
                          className="p-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted transition-colors disabled:opacity-50"
                        >
                          <Play className="w-3.5 h-3.5" />
                        </button>
                        <button
                          onClick={() => handleStop(job.id)}
                          disabled={busyJobID === job.id}
                          title="Stop job"
                          className="p-1.5 rounded-md text-muted-foreground hover:text-destructive hover:bg-destructive/10 transition-colors disabled:opacity-50"
                        >
                          <Square className="w-3.5 h-3.5" />
                        </button>
                        <button
                          onClick={() => {
                            setSelected(job);
                            setShowDialog(true);
                          }}
                          title="Edit job"
                          className="p-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
                        >
                          <Pencil className="w-3.5 h-3.5" />
                        </button>
                        <button
                          onClick={() => handleDelete(job.id)}
                          title="Delete job"
                          className="p-1.5 rounded-md text-muted-foreground hover:text-destructive hover:bg-destructive/10 transition-colors"
                        >
                          <Trash2 className="w-3.5 h-3.5" />
                        </button>
                      </div>
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>

        {!loading && jobs.length > 0 && (
          <div className="px-5 py-3 border-t border-border bg-muted/20">
            <span className="text-xs text-muted-foreground">
              {jobs.length} job{jobs.length !== 1 ? "s" : ""}
            </span>
          </div>
        )}
      </div>

      <CreateJobDialog
        open={showDialog}
        onOpenChange={setShowDialog}
        job={selected}
        onJobSaved={() => {
          setShowDialog(false);
          setSelected(null);
          load();
        }}
      />
    </div>
  );
}
