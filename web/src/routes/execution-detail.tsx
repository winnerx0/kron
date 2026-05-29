import { useEffect, useState } from "react";
import { useNavigate, useParams } from "@tanstack/react-router";
import { ArrowLeft, FileText, AlertTriangle } from "lucide-react";
import { getExecution, type ExecutionDetail } from "@/lib/api";

const STATUS_CFG: Record<string, { dot: string; text: string; bg: string }> = {
  pending: {
    dot: "bg-amber-400",
    text: "text-amber-600 dark:text-amber-400",
    bg: "bg-amber-50 dark:bg-amber-950/30",
  },
  running: {
    dot: "bg-blue-500",
    text: "text-blue-600 dark:text-blue-400",
    bg: "bg-blue-50 dark:bg-blue-950/30",
  },
  success: {
    dot: "bg-emerald-500",
    text: "text-emerald-700 dark:text-emerald-400",
    bg: "bg-emerald-50 dark:bg-emerald-950/30",
  },
  failed: {
    dot: "bg-red-500",
    text: "text-red-600 dark:text-red-400",
    bg: "bg-red-50 dark:bg-red-950/30",
  },
  stopped: {
    dot: "bg-zinc-500",
    text: "text-zinc-600 dark:text-zinc-400",
    bg: "bg-zinc-50 dark:bg-zinc-900/40",
  },
};

const METHOD_COLORS: Record<string, string> = {
  GET: "text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/40 border-emerald-200 dark:border-emerald-900",
  POST: "text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-950/40 border-blue-200 dark:border-blue-900",
  PUT: "text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-950/40 border-amber-200 dark:border-amber-900",
  PATCH:
    "text-purple-600 dark:text-purple-400 bg-purple-50 dark:bg-purple-950/40 border-purple-200 dark:border-purple-900",
  DELETE:
    "text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950/40 border-red-200 dark:border-red-900",
};

function StatusPill({ status }: { status: string }) {
  const cfg = STATUS_CFG[status.toLowerCase()] ?? STATUS_CFG["pending"];
  return (
    <span
      className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium ${cfg.bg} ${cfg.text}`}
    >
      <span
        className={`w-1.5 h-1.5 rounded-full shrink-0 ${cfg.dot} ${status === "running" ? "animate-pulse" : ""}`}
      />
      <span className="capitalize">{status}</span>
    </span>
  );
}

function fmt(d?: string) {
  if (!d) return "—";
  const date = new Date(d);
  if (date.getFullYear() <= 1) return "—";
  return date.toLocaleString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

function duration(start?: string, end?: string) {
  if (!start || !end) return "—";
  const started = new Date(start);
  const finished = new Date(end);
  if (finished.getFullYear() <= 1) return "—";
  const ms = finished.getTime() - started.getTime();
  if (ms < 1000) return `${ms}ms`;
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`;
  return `${Math.floor(ms / 60000)}m ${Math.floor((ms % 60000) / 1000)}s`;
}

function isUnfinished(status: string) {
  return status === "pending" || status === "running";
}

function MetaRow({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex items-start justify-between gap-6 px-5 py-3.5 border-b border-border last:border-0">
      <span className="font-mono text-[10px] uppercase tracking-[0.22em] text-muted-foreground pt-0.5 shrink-0">
        {label}
      </span>
      <span className="min-w-0 text-right text-sm text-foreground">
        {children}
      </span>
    </div>
  );
}

export function ExecutionDetailPage() {
  const navigate = useNavigate();
  const { id } = useParams({ strict: false }) as { id: string };
  const [execution, setExecution] = useState<ExecutionDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    const load = async () => {
      try {
        setLoading(true);
        const data = await getExecution(id);
        if (cancelled) return;
        setExecution(data);
        setError(null);
      } catch (e) {
        if (cancelled) return;
        setError(e instanceof Error ? e.message : "Failed to load execution");
      } finally {
        if (!cancelled) setLoading(false);
      }
    };
    load();
    return () => {
      cancelled = true;
    };
  }, [id]);

  const failed = execution?.status === "failed";

  return (
    <div className="space-y-6 animate-fade-in max-w-3xl">
      <button
        onClick={() => navigate({ to: "/executions" })}
        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
      >
        <ArrowLeft className="w-3.5 h-3.5" /> Executions
      </button>

      {loading ? (
        <div className="space-y-6">
          <div className="skeleton h-10 w-64" />
          <div className="bg-card border border-border rounded-xl overflow-hidden">
            {[...Array(5)].map((_, i) => (
              <div
                key={i}
                className="flex items-center justify-between px-5 py-3.5 border-b border-border last:border-0"
              >
                <div className="skeleton h-3 w-24" />
                <div className="skeleton h-3.5 w-40" />
              </div>
            ))}
          </div>
        </div>
      ) : error ? (
        <div className="px-4 py-3 bg-destructive/10 border border-destructive/20 text-destructive rounded-lg text-sm">
          {error}
        </div>
      ) : execution ? (
        <>
          <div className="flex items-end justify-between gap-4">
            <div className="min-w-0">
              <h1 className="font-display text-[36px] leading-tight truncate">
                {execution.jobName}
              </h1>
              <p className="text-sm text-muted-foreground mt-1">
                Execution detail
              </p>
            </div>
            <StatusPill status={execution.status} />
          </div>

          <div className="bg-card border border-border rounded-xl overflow-hidden">
            <MetaRow label="Endpoint">
              <span className="inline-flex items-center gap-2 flex-wrap justify-end">
                <span
                  className={`inline-block font-mono text-[10px] font-semibold px-2 py-0.5 rounded border ${METHOD_COLORS[execution.method] ?? "text-muted-foreground bg-muted border-border"}`}
                >
                  {execution.method}
                </span>
                <span className="font-mono text-xs text-muted-foreground break-all">
                  {execution.endpoint}
                </span>
              </span>
            </MetaRow>
            <MetaRow label="Started">
              <span className="text-muted-foreground">
                {fmt(execution.startedAt)}
              </span>
            </MetaRow>
            <MetaRow label="Completed">
              <span className="text-muted-foreground">
                {isUnfinished(execution.status)
                  ? "—"
                  : fmt(execution.finishedAt)}
              </span>
            </MetaRow>
            <MetaRow label="Duration">
              <span className="font-mono text-xs text-muted-foreground">
                {isUnfinished(execution.status)
                  ? "—"
                  : duration(execution.startedAt, execution.finishedAt)}
              </span>
            </MetaRow>
            <MetaRow label="Execution ID">
              <span className="font-mono text-xs text-muted-foreground break-all">
                {execution.id}
              </span>
            </MetaRow>
          </div>

          <div className="bg-card border border-border rounded-xl overflow-hidden">
            <div className="flex items-center gap-2 px-5 py-3 border-b border-border bg-muted/40">
              {failed ? (
                <AlertTriangle className="w-3.5 h-3.5 text-red-500" />
              ) : (
                <FileText className="w-3.5 h-3.5 text-muted-foreground" />
              )}
              <span className="font-mono text-[10px] uppercase tracking-[0.22em] text-muted-foreground">
                {failed ? "Error response" : "Response body"}
              </span>
            </div>
            {execution.responseBody ? (
              <pre
                className={`px-5 py-4 text-xs font-mono leading-relaxed whitespace-pre-wrap break-all overflow-x-auto max-h-[480px] ${
                  failed ? "text-red-600 dark:text-red-400" : "text-foreground"
                }`}
              >
                {execution.responseBody}
              </pre>
            ) : (
              <div className="px-5 py-10 text-center text-sm text-muted-foreground">
                No response body was recorded for this execution.
              </div>
            )}
          </div>
        </>
      ) : null}
    </div>
  );
}
