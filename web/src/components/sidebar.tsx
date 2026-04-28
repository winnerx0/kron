import { Link, useRouterState } from '@tanstack/react-router'
import { LayoutDashboard, Calendar, Activity, Settings, Sun, Moon } from 'lucide-react'
import { useTheme } from '@/lib/theme'

const NAV = [
  { to: '/',            label: 'Dashboard',  icon: LayoutDashboard },
  { to: '/jobs',        label: 'Jobs',        icon: Calendar },
  { to: '/executions',  label: 'Executions',  icon: Activity },
]

export function Sidebar() {
  const { theme, toggle } = useTheme()
  const router = useRouterState()
  const currentPath = router.location.pathname

  const isActive = (to: string) =>
    to === '/' ? currentPath === '/' : currentPath.startsWith(to)

  return (
    <aside className="fixed inset-y-0 left-0 w-56 flex flex-col border-r border-border bg-card z-30">
      {/* Logo */}
      <div className="h-14 flex items-center px-5 border-b border-border shrink-0">
        <div className="flex items-center gap-2.5">
          <svg width="18" height="18" viewBox="0 0 20 20" fill="none" className="shrink-0">
            <circle cx="10" cy="10" r="8" stroke="currentColor" strokeWidth="1.5"/>
            <line x1="10" y1="10" x2="10" y2="4.5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
            <line x1="10" y1="10" x2="13.5" y2="12" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
            <circle cx="10" cy="10" r="1" fill="currentColor"/>
          </svg>
          <span className="font-semibold text-[15px] tracking-tight">Kron</span>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-4 space-y-0.5 overflow-y-auto">
        <p className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground px-2 mb-3">
          Menu
        </p>
        {NAV.map(({ to, label, icon: Icon }) => (
          <Link
            key={to}
            to={to}
            className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm transition-colors ${
              isActive(to)
                ? 'bg-primary text-primary-foreground font-medium'
                : 'text-muted-foreground hover:text-foreground hover:bg-accent'
            }`}
          >
            <Icon className="w-4 h-4 shrink-0" />
            {label}
          </Link>
        ))}

        <div className="pt-4">
          <p className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground px-2 mb-3">
            System
          </p>
          <Link
            to="/settings"
            className={`flex items-center gap-3 px-2.5 py-2 rounded-md text-sm transition-colors ${
              isActive('/settings')
                ? 'bg-primary text-primary-foreground font-medium'
                : 'text-muted-foreground hover:text-foreground hover:bg-accent'
            }`}
          >
            <Settings className="w-4 h-4 shrink-0" />
            Settings
          </Link>
        </div>
      </nav>

      {/* Theme toggle + version */}
      <div className="px-3 py-4 border-t border-border shrink-0 space-y-3">
        <button
          onClick={toggle}
          className="w-full flex items-center gap-3 px-2.5 py-2 rounded-md text-sm text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
        >
          {theme === 'dark'
            ? <Sun className="w-4 h-4 shrink-0" />
            : <Moon className="w-4 h-4 shrink-0" />
          }
          {theme === 'dark' ? 'Light mode' : 'Dark mode'}
        </button>
        <p className="text-[10px] text-muted-foreground px-2 font-mono">v0.1.0</p>
      </div>
    </aside>
  )
}
