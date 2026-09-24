import { Link } from "@tanstack/react-router";
import { ArrowUpRight, Clock, History, RotateCw, ShieldCheck } from "lucide-react";
import { FaGithub } from "react-icons/fa";
import { KronMark } from "@/components/kron-mark";

const GITHUB_URL = "https://github.com/winnerx0/kron";

const FEATURES = [
  {
    icon: Clock,
    title: "Schedule with cron",
    body: "Standard 5-field expressions, validated on save. Set it once and Kron runs it on time, every time.",
  },
  {
    icon: History,
    title: "Full run history",
    body: "Every execution is logged with status, duration, and output — searchable and scoped per user.",
  },
  {
    icon: RotateCw,
    title: "Smart retries",
    body: "Backoff windows and ceilings out of the box. Failures surface; they don't silently retry forever.",
  },
  {
    icon: ShieldCheck,
    title: "Secure by default",
    body: "Google OAuth, scoped sessions, and refresh tokens. No DIY auth, no shared credentials.",
  },
];

export function LandingPage() {
  return (
    <div className="relative min-h-screen bg-background text-foreground">
      {/* Header */}
      <header className="mx-auto flex w-full max-w-5xl items-center justify-between px-6 py-6">
        <div className="flex items-center gap-2.5">
          <KronMark className="h-5 w-5 text-foreground" />
          <span className="text-[15px] font-semibold tracking-tight">Kron</span>
        </div>
        <nav className="hidden items-center gap-8 text-sm text-muted-foreground sm:flex">
          <a href="#features" className="transition-colors hover:text-foreground">Features</a>
          <a href="#how" className="transition-colors hover:text-foreground">How it works</a>
        </nav>
        <div className="flex items-center gap-5">
          <a
            href={GITHUB_URL}
            target="_blank"
            rel="noreferrer"
            aria-label="Kron on GitHub"
            className="text-muted-foreground transition-colors hover:text-foreground"
          >
            <FaGithub className="h-[18px] w-[18px]" />
          </a>
          <Link
            to="/login"
            className="text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            Sign in
          </Link>
        </div>
      </header>

      {/* Hero */}
      <section className="mx-auto w-full max-w-5xl px-6 pt-20 pb-24 text-center">
        <p className="mb-6 inline-flex items-center gap-2 rounded-full border border-border bg-card px-3 py-1 font-mono text-[11px] uppercase tracking-[0.18em] text-muted-foreground">
          <span className="h-1 w-1 rounded-full bg-foreground" />
          a job scheduler for humans
        </p>
        <h1 className="font-display text-[64px] leading-[0.98] sm:text-[88px]">
          Time, <span className="italic">on</span> schedule.
        </h1>
        <p className="mx-auto mt-6 max-w-lg text-base leading-relaxed text-muted-foreground">
          Kron schedules, runs, and tracks the recurring jobs that keep your
          systems alive &mdash; with cron you already know and a UI you'll
          actually want to open.
        </p>
        <div className="mt-10 flex items-center justify-center gap-3">
          <Link
            to="/login"
            className="group inline-flex h-11 items-center gap-2 rounded-full bg-primary px-5 text-sm font-medium text-primary-foreground transition-all hover:opacity-90 active:scale-[0.99]"
          >
            Get started
            <ArrowUpRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5 group-hover:-translate-y-0.5" />
          </Link>
          <a
            href={GITHUB_URL}
            target="_blank"
            rel="noreferrer"
            className="inline-flex h-11 items-center gap-2 rounded-full border border-border bg-card px-5 text-sm font-medium text-foreground transition-colors hover:bg-muted"
          >
            <FaGithub className="h-4 w-4" />
            View on GitHub
          </a>
        </div>
      </section>

      {/* Features */}
      <section id="features" className="border-t border-border">
        <div className="mx-auto w-full max-w-5xl px-6 py-20">
          <div className="mb-12 max-w-xl">
            <span className="font-mono text-[11px] uppercase tracking-[0.22em] text-muted-foreground">
              Features
            </span>
            <h2 className="font-display mt-3 text-[40px] leading-[1.05]">
              Everything you need to run jobs <span className="italic">reliably</span>.
            </h2>
          </div>

          <div className="grid grid-cols-1 gap-px overflow-hidden rounded-2xl border border-border bg-border sm:grid-cols-2">
            {FEATURES.map(({ icon: Icon, title, body }) => (
              <div key={title} className="bg-background p-8 transition-colors hover:bg-card">
                <Icon className="h-5 w-5 text-foreground" />
                <h3 className="mt-6 text-base font-semibold tracking-tight">{title}</h3>
                <p className="mt-2 text-sm leading-relaxed text-muted-foreground">{body}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* How it works */}
      <section id="how" className="border-t border-border">
        <div className="mx-auto w-full max-w-5xl px-6 py-20">
          <div className="mb-12 max-w-xl">
            <span className="font-mono text-[11px] uppercase tracking-[0.22em] text-muted-foreground">
              How it works
            </span>
            <h2 className="font-display mt-3 text-[40px] leading-[1.05]">
              From cron line to <span className="italic">running job</span>, in three steps.
            </h2>
          </div>

          <ol className="grid grid-cols-1 gap-6 sm:grid-cols-3">
            {[
              { n: "01", t: "Sign in", b: "Continue with Google. Your jobs and history stay scoped to your account." },
              { n: "02", t: "Define a job", b: "Give it a name, a command, and a cron schedule. Kron validates and previews the next runs." },
              { n: "03", t: "Watch it run", b: "Track every execution in the dashboard with status, output, and duration." },
            ].map((step) => (
              <li key={step.n} className="rounded-xl border border-border bg-card p-6">
                <span className="font-mono text-[11px] uppercase tracking-[0.22em] text-muted-foreground">
                  {step.n}
                </span>
                <h3 className="mt-4 text-base font-semibold tracking-tight">{step.t}</h3>
                <p className="mt-2 text-sm leading-relaxed text-muted-foreground">{step.b}</p>
              </li>
            ))}
          </ol>
        </div>
      </section>

      {/* CTA */}
      <section className="border-t border-border">
        <div className="mx-auto w-full max-w-5xl px-6 py-24 text-center">
          <h2 className="font-display mx-auto max-w-2xl text-[44px] leading-[1.05] sm:text-[56px]">
            Stop babysitting <span className="italic">cron tables.</span>
          </h2>
          <p className="mx-auto mt-5 max-w-md text-sm text-muted-foreground">
            Sign in and schedule your first job in under a minute.
          </p>
          <div className="mt-8">
            <Link
              to="/login"
              className="group inline-flex h-12 items-center gap-2 rounded-full bg-primary px-6 text-sm font-medium text-primary-foreground transition-all hover:opacity-90 active:scale-[0.99]"
            >
              Sign in to Kron
              <ArrowUpRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5 group-hover:-translate-y-0.5" />
            </Link>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-border">
        <div className="mx-auto flex w-full max-w-5xl items-center justify-between px-6 py-8 font-mono text-[11px] uppercase tracking-[0.22em] text-muted-foreground">
          <div className="flex items-center gap-2">
            <KronMark className="h-3.5 w-3.5 text-foreground" />
            <span>kron · scheduler</span>
          </div>
          <div className="flex items-center gap-5">
            <a
              href={GITHUB_URL}
              target="_blank"
              rel="noreferrer"
              className="inline-flex items-center gap-1.5 transition-colors hover:text-foreground"
            >
              <FaGithub className="h-3.5 w-3.5" />
              github
            </a>
            <span>© {new Date().getFullYear()}</span>
          </div>
        </div>
      </footer>
    </div>
  );
}
