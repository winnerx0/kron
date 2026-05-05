import { Link, useRouterState } from '@tanstack/react-router'
import { LayoutDashboard, Calendar, Activity, Settings, Sun, Moon } from 'lucide-react'
import { useTheme } from '@/lib/theme'
import { KronMark } from '@/components/kron-mark'

const NAV = [
  { to: '/dashboard',   label: 'Dashboard',  icon: LayoutDashboard },
  { to: '/jobs',        label: 'Jobs',        icon: Calendar },
  { to: '/executions',  label: 'Executions',  icon: Activity },
]

export function Sidebar() {
  const { theme, toggle } = useTheme()
  const router = useRouterState()
  const currentPath = router.location.pathname

  const isActive = (to: string) => currentPath === to || (to !== '/dashboard' && currentPath.startsWith(to))

  return (
    <aside className="fixed inset-y-0 left-0 w-56 flex flex-col border-r border-border bg-card z-30">
      {/* Logo */}
      <div className="h-14 flex items-center px-5 border-b border-border shrink-0">
        <div className="flex items-center gap-2.5">
          <KronMark className="h-5 w-5 text-foreground shrink-0" />
          <span className="font-semibold text-[15px] tracking-tight">Kron</span>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-5 space-y-0.5 overflow-y-auto">
        <p className="font-mono text-[10px] uppercase tracking-[0.22em] text-muted-foreground px-2 mb-3">
          Menu
        </p>
        {NAV.map(({ to, label, icon: Icon }) => (
          <Link
            key={to}
            to={to}
            className={`flex items-center gap-3 px-2.5 py-2 rounded-full text-sm transition-colors ${
              isActive(to)
                ? 'bg-primary text-primary-foreground font-medium'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted'
            }`}
          >
            <Icon className="w-4 h-4 shrink-0" />
            {label}
          </Link>
        ))}

        <div className="pt-5">
          <p className="font-mono text-[10px] uppercase tracking-[0.22em] text-muted-foreground px-2 mb-3">
            System
          </p>
          <Link
            to="/settings"
            className={`flex items-center gap-3 px-2.5 py-2 rounded-full text-sm transition-colors ${
              isActive('/settings')
                ? 'bg-primary text-primary-foreground font-medium'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted'
            }`}
          >
            <Settings className="w-4 h-4 shrink-0" />
            Settings
          </Link>
        </div>
      </nav>

      {/* Theme toggle */}
      <div className="px-3 py-4 border-t border-border shrink-0 space-y-1">
        <button
          onClick={toggle}
          className="w-full flex items-center gap-3 px-2.5 py-2 rounded-full text-sm text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
        >
          {theme === 'dark'
            ? <Sun className="w-4 h-4 shrink-0" />
            : <Moon className="w-4 h-4 shrink-0" />
          }
          {theme === 'dark' ? 'Light mode' : 'Dark mode'}
        </button>
        <p className="font-mono text-[10px] uppercase tracking-[0.22em] text-muted-foreground px-2.5 pt-1">
          v0.1.0
        </p>
      </div>
    </aside>
  )
}
