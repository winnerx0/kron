import { useEffect, useState } from "react";
import { Link } from "@tanstack/react-router";
import {
  Calendar,
  Activity,
  CheckCircle,
  XCircle,
  ArrowRight,
} from "lucide-react";
import { getDashboard, type DashboardData } from "@/lib/api";

function StatCard({
  label,
  value,
  icon: Icon,
  sub,
  delay = 0,
}: {
  label: string;
  value: string | number;
  icon: React.ElementType;
  sub?: string;
  delay?: number;
}) {
  return (
    <div
      className="bg-card border border-border rounded-xl p-5 flex flex-col gap-3 animate-slide-up"
      style={{ animationDelay: `${delay}ms` }}
    >
      <div className="flex items-center justify-between">
        <span className="font-mono text-[10px] uppercase tracking-[0.22em] text-muted-foreground">
          {label}
        </span>
        <div className="p-1.5 rounded-md bg-muted">
          <Icon className="w-3.5 h-3.5 text-muted-foreground" />
        </div>
      </div>
      <div>
        <p className="text-3xl font-semibold tracking-tight leading-none">
          {value}
        </p>
        {sub && <p className="text-xs text-muted-foreground mt-1.5">{sub}</p>}
      </div>
    </div>
  );
}

const STATUS_DOT: Record<string, string> = {
  success: "bg-foreground",
  completed: "bg-foreground",
  failed: "bg-destructive",
  running: "bg-foreground",
  pending: "bg-muted-foreground",
};

function formatAgo(dateStr?: string) {
  if (!dateStr) return "—";
  const diff = Date.now() - new Date(dateStr).getTime();
  const m = Math.floor(diff / 60000);
  if (m < 1) return "just now";
  if (m < 60) return `${m}m ago`;
  const h = Math.floor(m / 60);
  if (h < 24) return `${h}h ago`;
  return `${Math.floor(h / 24)}d ago`;
}

export function DashboardPage() {
  const [data, setData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getDashboard()
      .then((dashboard) => {
        setData(dashboard);
        setError(null);
      })
      .catch((err) => {
        setError(
          err instanceof Error ? err.message : "Failed to load dashboard",
        );
      })
      .finally(() => setLoading(false));
  }, []);

  const total = data?.totalJobs ?? 0;
  const enabledJobs = data?.enabledJobs ?? 0;
  const disabledJobs = data?.disabledJobs ?? 0;
  const jobsSummary =
    disabledJobs > 0
      ? `${enabledJobs} enabled, ${disabledJobs} disabled`
      : `${enabledJobs} enabled`;
  const executions = data?.totalExecutions ?? 0;
  const successful = data?.successfulExecutions ?? 0;
  const failed = data?.failedExecutions ?? 0;
  const rate = data?.successRate ?? 0;
  const recent = data?.recentExecutions ?? [];
  const jobs = data?.jobs ?? [];

  return (
    <div className="space-y-8 animate-fade-in">
      <div>
        <h1 className="font-display text-[36px] leading-tight">Dashboard</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Overview of your job scheduler
        </p>
      </div>

      {error && (
        <div className="px-4 py-3 bg-destructive/10 border border-destructive/20 text-destructive rounded-lg text-sm">
          {error}
        </div>
      )}

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        {loading ? (
          [...Array(4)].map((_, i) => (
            <div
              key={i}
              className="bg-card border border-border rounded-xl p-5 h-28 skeleton"
            />
          ))
        ) : (
          <>
            <StatCard
              label="Total Jobs"
              value={total}
              icon={Calendar}
              sub={jobsSummary}
              delay={0}
            />
            <StatCard
              label="Executions"
              value={executions}
              icon={Activity}
              sub="all time"
              delay={50}
            />
            <StatCard
              label="Successful"
              value={successful}
              icon={CheckCircle}
              sub={`${rate}% success rate`}
              delay={100}
            />
            <StatCard
              label="Failed"
              value={failed}
              icon={XCircle}
              sub={failed ? "needs attention" : "all clear"}
              delay={150}
            />
          </>
        )}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-5 gap-6">
        {/* Recent executions */}
        <div
          className="lg:col-span-3 bg-card border border-border rounded-xl animate-slide-up"
          style={{ animationDelay: "200ms" }}
        >
          <div className="flex items-center justify-between px-5 py-4 border-b border-border">
            <h2 className="font-mono text-[11px] uppercase tracking-[0.22em] text-muted-foreground">
              Recent Executions
            </h2>
            <Link
              to="/executions"
              className="text-xs text-muted-foreground hover:text-foreground flex items-center gap-1 transition-colors"
            >
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
              <div className="px-5 py-8 text-center text-sm text-muted-foreground">
                No executions yet
              </div>
            ) : (
              recent.map((ex, i) => {
                const dotColor = STATUS_DOT[ex.status] ?? "bg-muted-foreground";
                return (
                  <div
                    key={ex.id}
                    className="flex items-center gap-3 px-5 py-3.5 text-sm hover:bg-muted/30 transition-colors animate-slide-up"
                    style={{ animationDelay: `${220 + i * 30}ms` }}
                  >
                    <span
                      className={`w-1.5 h-1.5 rounded-full shrink-0 ${dotColor} ${ex.status === "running" ? "animate-pulse" : ""}`}
                    />
                    <span className="min-w-0 truncate text-xs text-muted-foreground">
                      {ex.jobName || ex.jobID.slice(0, 8)}
                    </span>
                    <span className="capitalize text-xs text-muted-foreground">
                      {ex.status}
                    </span>
                    <span className="ml-auto text-xs text-muted-foreground">
                      {formatAgo(ex.startedAt)}
                    </span>
                  </div>
                );
              })
            )}
          </div>
        </div>

        {/* Jobs list */}
        <div
          className="lg:col-span-2 bg-card border border-border rounded-xl animate-slide-up"
          style={{ animationDelay: "250ms" }}
        >
          <div className="flex items-center justify-between px-5 py-4 border-b border-border">
            <h2 className="font-mono text-[11px] uppercase tracking-[0.22em] text-muted-foreground">
              Jobs
            </h2>
            <Link
              to="/jobs"
              className="text-xs text-muted-foreground hover:text-foreground flex items-center gap-1 transition-colors"
            >
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
                <p className="text-sm text-muted-foreground mb-3">
                  No jobs yet
                </p>
                <Link
                  to="/jobs"
                  className="text-xs underline text-muted-foreground hover:text-foreground"
                >
                  Create one →
                </Link>
              </div>
            ) : (
              jobs.map((job, i) => (
                <div
                  key={job.id}
                  className="px-5 py-3.5 animate-slide-up"
                  style={{ animationDelay: `${270 + i * 30}ms` }}
                >
                  <div className="flex items-center gap-2">
                    <p className="text-sm font-medium truncate">{job.name}</p>
                    {job.enabled === false && (
                      <span className="shrink-0 rounded-full border border-border bg-muted px-2 py-0.5 font-mono text-[10px] text-muted-foreground">
                        Disabled
                      </span>
                    )}
                  </div>
                  <p className="font-mono text-xs text-muted-foreground mt-0.5">
                    {job.schedule}
                  </p>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
