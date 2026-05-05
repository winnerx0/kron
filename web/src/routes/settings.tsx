import { useEffect, useMemo, useState } from 'react'
import { Check, LogOut, Moon, Sun, UserRound } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { clearTokens } from '@/lib/auth'
import { CurrentUser, getCurrentUser } from '@/lib/api'
import { useTheme } from '@/lib/theme'

function Section({
  icon: Icon,
  title,
  description,
  children,
}: {
  icon: React.ComponentType<{ className?: string }>
  title: string
  description: string
  children: React.ReactNode
}) {
  return (
    <section className="grid gap-5 border-t border-border py-6 md:grid-cols-[220px_1fr]">
      <div className="flex gap-3">
        <span className="grid h-9 w-9 shrink-0 place-items-center rounded-md border border-border bg-card">
          <Icon className="h-4 w-4 text-muted-foreground" />
        </span>
        <div>
          <h2 className="font-mono text-[11px] uppercase tracking-[0.22em] text-foreground">{title}</h2>
          <p className="mt-1 text-sm leading-5 text-muted-foreground">{description}</p>
        </div>
      </div>
      <div className="min-w-0">{children}</div>
    </section>
  )
}

function FieldRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex min-h-10 items-center justify-between gap-4 border-b border-border py-2 last:border-0">
      <span className="text-sm text-muted-foreground">{label}</span>
      <span className="min-w-0 truncate text-right text-sm font-medium">{value}</span>
    </div>
  )
}

export function SettingsPage() {
  
  const [user, setUser] = useState<CurrentUser | null>(null)
  const [loadingUser, setLoadingUser] = useState(true)
  const [userError, setUserError] = useState<string | null>(null)

  useEffect(() => {
    let active = true
    setLoadingUser(true)
    getCurrentUser()
      .then(profile => {
        if (!active) return
        setUser(profile)
        setUserError(null)
      })
      .catch(err => {
        if (!active) return
        setUserError(err instanceof Error ? err.message : 'Failed to load profile')
      })
      .finally(() => {
        if (active) setLoadingUser(false)
      })

    return () => {
      active = false
    }
  }, [])

  const joinedAt = useMemo(() => {
    if (!user?.joined_at) return 'Not available'
    return new Intl.DateTimeFormat(undefined, {
      dateStyle: 'medium',
      timeStyle: 'short',
    }).format(new Date(user.joined_at))
  }, [user?.joined_at])

  const signOut = () => {
    clearTokens()
    window.location.href = '/login'
  }

  return (
    <div className="mx-auto max-w-4xl animate-fade-in">
      <header className="mb-7 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="font-display text-[36px] leading-tight">Settings</h1>
          <p className="mt-1 text-sm text-muted-foreground">Manage account, display, and runtime preferences.</p>
        </div>
        <button
          className="inline-flex items-center gap-2 h-9 rounded-full border border-border bg-card px-4 text-sm font-medium text-foreground transition-colors hover:bg-muted self-start sm:self-auto"
          onClick={signOut}
        >
          <LogOut className="h-4 w-4" />
          Sign out
        </button>
      </header>

      <div className="rounded-xl border border-border bg-card px-5 sm:px-6">
        <Section icon={UserRound} title="Account" description="Profile details from your workspace identity.">
          {loadingUser ? (
            <div className="h-28 animate-pulse rounded-md border border-border bg-muted/30" />
          ) : userError ? (
            <div className="rounded-md border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive">
              {userError}
            </div>
          ) : user ? (
            <div className="space-y-4">
              <div className="flex items-center gap-4">
                <img
                  src={user.profile_picture}
                  alt=""
                  className="h-14 w-14 rounded-full border border-border object-cover"
                  referrerPolicy="no-referrer"
                />
                <div className="min-w-0">
                  <p className="truncate text-sm font-semibold">{user.name}</p>
                  <p className="truncate text-sm text-muted-foreground">{user.email}</p>
                </div>
              </div>
              <div className="rounded-md border border-border px-4">
                <FieldRow label="Joined" value={joinedAt} />
              </div>
            </div>
          ) : null}
        </Section>

      </div>
    </div>
  )
}
