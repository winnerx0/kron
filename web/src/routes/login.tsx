import { useEffect, useMemo, useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { FcGoogle } from 'react-icons/fc'
import { Button } from '@/components/ui/button'
import { getAccessToken, LOGIN_URL } from '@/lib/auth'

export function LoginPage() {
  const navigate = useNavigate()
  const [now, setNow] = useState<Date>(() => new Date())

  useEffect(() => {
    if (getAccessToken()) navigate({ to: '/', replace: true })
  }, [navigate])

  useEffect(() => {
    const id = window.setInterval(() => setNow(new Date()), 1000)
    return () => window.clearInterval(id)
  }, [])

  const utc = useMemo(() => formatUTC(now), [now])
  const session = useMemo(() => sessionId(), [])

  return (
    <div className="relative min-h-screen w-full overflow-hidden bg-background text-foreground">
      <GridPlate />

      {/* corner ticks */}
      <CornerTick className="top-6 left-6" />
      <CornerTick className="top-6 right-6 rotate-90" />
      <CornerTick className="bottom-6 left-6 -rotate-90" />
      <CornerTick className="bottom-6 right-6 rotate-180" />

      {/* top meta strip */}
      <header className="absolute inset-x-0 top-0 px-8 pt-8 flex items-center justify-between text-[10px] font-mono uppercase tracking-[0.2em] text-muted-foreground select-none">
        <div className="flex items-center gap-2 animate-fade-in">
          <span className="inline-block w-1.5 h-1.5 rounded-full bg-[hsl(14_88%_55%)] shadow-[0_0_8px_hsl(14_88%_55%/0.7)]" />
          <span>Kron / Control</span>
        </div>
        <div className="hidden sm:flex items-center gap-6 animate-fade-in" style={{ animationDelay: '120ms' }}>
          <span>Build 2026.05</span>
          <span>{utc}</span>
        </div>
      </header>

      <main className="relative z-10 min-h-screen grid lg:grid-cols-[1.05fr_0.95fr]">
        {/* LEFT — editorial column */}
        <section className="hidden lg:flex flex-col justify-between px-12 xl:px-20 pt-32 pb-16 border-r border-border/70">
          <div className="space-y-10">
            <div
              className="text-[10px] font-mono uppercase tracking-[0.3em] text-muted-foreground animate-slide-up"
              style={{ animationDelay: '60ms' }}
            >
              §01 — Authenticate
            </div>

            <h1
              className="font-light leading-[0.92] tracking-[-0.04em] animate-slide-up"
              style={{ animationDelay: '120ms', fontSize: 'clamp(3rem, 6vw, 5.25rem)' }}
            >
              Schedule
              <br />
              <span className="italic font-serif text-[hsl(14_88%_55%)]" style={{ fontFamily: '"Cormorant Garamond", "Times New Roman", serif' }}>
                with intent.
              </span>
            </h1>

            <p
              className="max-w-md text-[15px] leading-relaxed text-muted-foreground animate-slide-up"
              style={{ animationDelay: '200ms' }}
            >
              A precise control surface for cron jobs and background work — every
              run logged, every failure surfaced, every secret encrypted at rest.
            </p>
          </div>

          <ManifestBlock delay={320} />
        </section>

        {/* RIGHT — sign-in panel */}
        <section className="relative flex items-center justify-center px-6 sm:px-12 py-24">
          <div
            className="relative w-full max-w-[420px] animate-slide-up"
            style={{ animationDelay: '180ms' }}
          >
            {/* hairline frame */}
            <div className="relative bg-card border border-border">
              <div className="absolute -top-px left-6 right-6 h-px bg-foreground/30" />
              <div className="absolute -bottom-px left-6 right-6 h-px bg-foreground/10" />

              <div className="px-8 pt-10 pb-8 space-y-8">
                <div className="flex items-center justify-between">
                  <Monogram />
                  <span className="text-[10px] font-mono uppercase tracking-[0.25em] text-muted-foreground">
                    {session}
                  </span>
                </div>

                <div className="space-y-2">
                  <div className="text-[10px] font-mono uppercase tracking-[0.25em] text-muted-foreground">
                    Sign in
                  </div>
                  <h2 className="text-2xl font-medium tracking-[-0.02em]">
                    Open the console.
                  </h2>
                  <p className="text-sm text-muted-foreground">
                    Continue with your workspace identity to access scheduled jobs and execution history.
                  </p>
                </div>

                <Button
                  size="lg"
                  variant="outline"
                  className="group w-full h-12 justify-between gap-3 px-5 text-[13px] font-medium tracking-tight"
                  onClick={() => { window.location.href = LOGIN_URL }}
                >
                  <span className="flex items-center gap-3">
                    <span className="grid place-items-center w-7 h-7 rounded-full bg-background border border-border">
                      <FcGoogle className="w-4 h-4" />
                    </span>
                    Continue with Google
                  </span>
                  <span
                    aria-hidden
                    className="font-mono text-muted-foreground transition-transform duration-200 group-hover:translate-x-0.5"
                  >
                    →
                  </span>
                </Button>

                <div className="flex items-center gap-3 text-[10px] font-mono uppercase tracking-[0.25em] text-muted-foreground">
                  <span className="h-px flex-1 bg-border" />
                  <span>SSO · OAuth 2.0</span>
                  <span className="h-px flex-1 bg-border" />
                </div>

                <ul className="space-y-2.5 text-[12px] font-mono text-muted-foreground">
                  <FootRow k="ENC" v="AES-256-GCM at rest" />
                  <FootRow k="JWT" v="HS256 / 15-min access" />
                  <FootRow k="LOG" v="every run, retained 30d" />
                </ul>
              </div>
            </div>

            <p className="mt-6 text-[11px] font-mono text-muted-foreground/80 leading-relaxed">
              By continuing you agree to operate Kron in accordance with your
              workspace policy. Tokens are scoped to this device.
            </p>
          </div>
        </section>
      </main>

      {/* bottom meta strip */}
      <footer className="absolute inset-x-0 bottom-0 px-8 pb-6 flex items-center justify-between text-[10px] font-mono uppercase tracking-[0.2em] text-muted-foreground/80 select-none">
        <span>kron.scheduler / v0.0.1</span>
        <span className="hidden sm:inline">node-01 · ready</span>
        <span>{utc.split(' ')[1]}</span>
      </footer>
    </div>
  )
}

function FootRow({ k, v }: { k: string; v: string }) {
  return (
    <li className="flex items-center gap-3">
      <span className="text-foreground/70 tracking-[0.15em]">{k}</span>
      <span className="h-px flex-1 bg-border" />
      <span>{v}</span>
    </li>
  )
}

function Monogram() {
  return (
    <div className="flex items-center gap-2">
      <div className="relative w-7 h-7 grid place-items-center border border-foreground/80">
        <span className="font-serif italic text-[15px] leading-none" style={{ fontFamily: '"Cormorant Garamond", "Times New Roman", serif' }}>
          k
        </span>
        <span className="absolute -top-1 -right-1 w-1.5 h-1.5 bg-[hsl(14_88%_55%)]" />
      </div>
      <span className="text-[13px] font-medium tracking-tight">Kron</span>
    </div>
  )
}

function ManifestBlock({ delay = 0 }: { delay?: number }) {
  const lines: Array<[string, string]> = [
    ['NAMESPACE', 'production'],
    ['JOBS', '128 scheduled / 4 active'],
    ['LATENCY', '7ms p50 · 41ms p99'],
    ['UPTIME', '99.982% / 30d'],
  ]

  return (
    <div
      className="font-mono text-[11px] leading-relaxed text-muted-foreground animate-slide-up"
      style={{ animationDelay: `${delay}ms` }}
    >
      <div className="flex items-center gap-2 text-foreground/70 mb-3 tracking-[0.2em] uppercase">
        <span className="inline-block w-3 h-px bg-foreground/70" />
        Manifest
      </div>
      <ul className="space-y-1.5">
        {lines.map(([k, v], i) => (
          <li
            key={k}
            className="flex items-baseline gap-3 animate-fade-in"
            style={{ animationDelay: `${delay + 100 + i * 70}ms` }}
          >
            <span className="text-foreground/60 w-24 tracking-[0.12em]">{k}</span>
            <span className="flex-1 border-b border-dashed border-border/80 translate-y-[-3px]" />
            <span className="text-foreground/80">{v}</span>
          </li>
        ))}
      </ul>
    </div>
  )
}

function CornerTick({ className = '' }: { className?: string }) {
  return (
    <div
      aria-hidden
      className={`pointer-events-none absolute w-3 h-3 ${className}`}
    >
      <span className="absolute top-0 left-0 w-3 h-px bg-foreground/40" />
      <span className="absolute top-0 left-0 w-px h-3 bg-foreground/40" />
    </div>
  )
}

function GridPlate() {
  return (
    <>
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 opacity-[0.35] dark:opacity-[0.18]"
        style={{
          backgroundImage:
            'linear-gradient(hsl(var(--border)) 1px, transparent 1px), linear-gradient(90deg, hsl(var(--border)) 1px, transparent 1px)',
          backgroundSize: '64px 64px',
          maskImage:
            'radial-gradient(ellipse 70% 60% at 50% 40%, black 30%, transparent 80%)',
          WebkitMaskImage:
            'radial-gradient(ellipse 70% 60% at 50% 40%, black 30%, transparent 80%)',
        }}
      />
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0"
        style={{
          background:
            'radial-gradient(ellipse 60% 50% at 80% 100%, hsl(14 88% 55% / 0.06), transparent 60%)',
        }}
      />
    </>
  )
}

function pad(n: number) { return n.toString().padStart(2, '0') }

function formatUTC(d: Date): string {
  const yyyy = d.getUTCFullYear()
  const mm = pad(d.getUTCMonth() + 1)
  const dd = pad(d.getUTCDate())
  const hh = pad(d.getUTCHours())
  const mi = pad(d.getUTCMinutes())
  const ss = pad(d.getUTCSeconds())
  return `${yyyy}.${mm}.${dd} ${hh}:${mi}:${ss}Z`
}

function sessionId(): string {
  const bytes = new Uint8Array(3)
  if (typeof crypto !== 'undefined' && crypto.getRandomValues) {
    crypto.getRandomValues(bytes)
  } else {
    for (let i = 0; i < bytes.length; i++) bytes[i] = Math.floor(Math.random() * 256)
  }
  return 'SID·' + Array.from(bytes).map(b => b.toString(16).padStart(2, '0')).join('').toUpperCase()
}
