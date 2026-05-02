"use client";

import { useEffect, useState } from "react";
import { Activity } from "lucide-react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { getExecutions, type Execution } from "@/lib/api";

const STATUS_CONFIG: Record<string, { dot: string; label: string }> = {
  pending: { dot: "bg-amber-400", label: "text-amber-700" },
  running: { dot: "bg-blue-500", label: "text-blue-700" },
  success: { dot: "bg-emerald-500", label: "text-emerald-700" },
  failed: { dot: "bg-red-500", label: "text-red-700" },
};

function StatusDot({ status }: { status: string }) {
  const cfg = STATUS_CONFIG[status] ?? STATUS_CONFIG["pending"];
  const isRunning = status === "pending";
  return (
    <span className="inline-flex items-center gap-1.5">
      <span
        className={`relative inline-block w-1.5 h-1.5 rounded-full ${cfg.dot}`}
      >
        {isRunning && (
          <span
            className={`absolute inset-0 rounded-full ${cfg.dot} animate-ping opacity-75`}
          />
        )}
      </span>
      <span className={`text-xs font-medium capitalize ${cfg.label}`}>
        {status}
      </span>
    </span>
  );
}

function SkeletonRow() {
  return (
    <div className="flex items-center gap-6 px-6 py-3.5 border-b border-border last:border-0">
      <div className="skeleton h-3.5 w-20 rounded-full" />
      <div className="skeleton h-3.5 w-16" />
      <div className="skeleton h-3.5 w-32" />
      <div className="skeleton h-3.5 w-32" />
    </div>
  );
}

function EmptyState() {
  return (
    <div className="flex flex-col items-center justify-center py-16 px-6 text-center animate-fade-in">
      <div className="mb-5 text-muted-foreground/30">
        <Activity strokeWidth={1} className="w-14 h-14 mx-auto" />
      </div>
      <p className="text-sm font-medium text-foreground mb-1">
        No executions yet
      </p>
      <p className="text-sm text-muted-foreground max-w-xs">
        Jobs will appear here once they start executing.
      </p>
    </div>
  );
}

function formatDate(dateStr?: string) {
  if (!dateStr) return "—";
  const date = new Date(dateStr);
  if (date.getFullYear() <= 1) return "—";
  return date.toLocaleString("en-US", {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

function isUnfinished(status: string) {
  return status === "pending" || status === "running";
}

export function ExecutionsList() {
  const [executions, setExecutions] = useState<Execution[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastRefresh, setLastRefresh] = useState(Date.now());

  useEffect(() => {
    const load = async () => {
      try {
        setLoading(true);
        const data = await getExecutions();
        setExecutions(data || []);
        setError(null);
        setLastRefresh(Date.now());
      } catch (err) {
        setError(
          err instanceof Error ? err.message : "Failed to load executions",
        );
      } finally {
        setLoading(false);
      }
    };
    load();
    const interval = setInterval(load, 5000);
    return () => clearInterval(interval);
  }, []);

  const secondsAgo = Math.round((Date.now() - lastRefresh) / 1000);

  return (
    <Card>
      <CardHeader className="px-6 py-4 flex flex-row items-center justify-between space-y-0 border-b border-border">
        <div>
          <h2 className="text-base font-semibold tracking-tight">Executions</h2>
          <p className="text-xs text-muted-foreground mt-0.5">
            Live execution history
          </p>
        </div>
        <div className="flex items-center gap-1.5">
          <span className="relative flex h-1.5 w-1.5">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75" />
            <span className="relative inline-flex rounded-full h-1.5 w-1.5 bg-emerald-500" />
          </span>
          <span className="text-xs text-muted-foreground">
            {secondsAgo < 5 ? "Just updated" : `${secondsAgo}s ago`}
          </span>
        </div>
      </CardHeader>

      <CardContent className="p-0">
        {error && (
          <div className="mx-6 mt-4 px-4 py-3 bg-red-50 border border-red-200 text-red-700 rounded-md text-xs">
            {error}
          </div>
        )}

        {loading ? (
          <div className="divide-y divide-border">
            {[...Array(3)].map((_, i) => (
              <SkeletonRow key={i} />
            ))}
          </div>
        ) : executions.length === 0 ? (
          <EmptyState />
        ) : (
          <div className="divide-y divide-border">
            {/* Column headers */}
            <div className="flex items-center gap-6 px-6 py-2 bg-muted/40">
              <span className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground w-20 shrink-0">
                Status
              </span>
              <span className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground w-20 shrink-0">
                Job ID
              </span>
              <span className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground flex-1">
                Started
              </span>
              <span className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground flex-1">
                Completed
              </span>
            </div>

            {executions.map((execution, idx) => (
              <div
                key={execution.id}
                className="flex items-center gap-6 px-6 py-3.5 hover:bg-muted/20 transition-colors animate-slide-up"
                style={{ animationDelay: `${idx * 30}ms` }}
              >
                <div className="w-20 shrink-0">
                  <StatusDot status={execution.status} />
                </div>
                <span className="mono text-xs text-muted-foreground w-20 shrink-0 truncate">
                  {execution.jobID.slice(0, 8)}
                </span>
                <span className="text-xs text-muted-foreground flex-1">
                  {formatDate(execution.startedAt)}
                </span>
                <span className="text-xs text-muted-foreground flex-1">
                  {isUnfinished(execution.status)
                    ? ""
                    : formatDate(execution.finishedAt)}
                </span>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
