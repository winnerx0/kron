import { useEffect, useState } from "react";
import {
  Plus,
  Pencil,
  Trash2,
  CalendarClock,
  Play,
  Square,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { getJobs, deleteJob, runJob, stopJob, type Job } from "@/lib/api";
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
        <div className="skeleton h-3.5 w-10 ml-auto" />
      </td>
    </tr>
  );
}

export function JobsPage() {
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
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to delete");
    }
  };

  const handleRun = async (id: string) => {
    try {
      setBusyJobID(id);
      await runJob(id);
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
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to stop");
    } finally {
      setBusyJobID(null);
    }
  };

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-end justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Jobs</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage your scheduled HTTP tasks
          </p>
        </div>
        <Button
          size="sm"
          className="gap-2 h-8 text-xs"
          onClick={() => {
            setSelected(null);
            setShowDialog(true);
          }}
        >
          <Plus className="w-3.5 h-3.5" /> New Job
        </Button>
      </div>

      {error && (
        <div className="px-4 py-3 bg-destructive/10 border border-destructive/20 text-destructive rounded-lg text-sm">
          {error}
        </div>
      )}

      <div className="bg-card border border-border rounded-lg overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border bg-muted/40">
              {["Name", "Method", "Endpoint", "Schedule", ""].map((h) => (
                <th
                  key={h}
                  className={`px-5 py-3 text-[10px] font-semibold uppercase tracking-widest text-muted-foreground ${h === "" ? "text-right" : "text-left"}`}
                >
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {loading ? (
              [...Array(4)].map((_, i) => <SkeletonRow key={i} />)
            ) : jobs.length === 0 ? (
              <tr>
                <td colSpan={5}>
                  <div className="flex flex-col items-center justify-center py-16 text-center">
                    <CalendarClock
                      className="w-10 h-10 text-muted-foreground/30 mb-4"
                      strokeWidth={1}
                    />
                    <p className="text-sm font-medium mb-1">No jobs yet</p>
                    <p className="text-sm text-muted-foreground mb-5">
                      Create your first scheduled job to get started
                    </p>
                    <Button
                      size="sm"
                      className="gap-2 h-8 text-xs"
                      onClick={() => {
                        setSelected(null);
                        setShowDialog(true);
                      }}
                    >
                      <Plus className="w-3.5 h-3.5" /> Create job
                    </Button>
                  </div>
                </td>
              </tr>
            ) : (
              jobs.map((job, idx) => (
                <tr
                  key={job.id}
                  className="border-b border-border last:border-0 hover:bg-muted/30 transition-colors group animate-slide-up"
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
                  <td className="px-5 py-3.5">
                    <div className="flex items-center justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button
                        onClick={() => handleRun(job.id)}
                        disabled={busyJobID === job.id}
                        title="Run job"
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
              ))
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
