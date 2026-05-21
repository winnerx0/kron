import { Link, useRouterState } from '@tanstack/react-router'
import { LayoutDashboard, Calendar, Activity, Settings, Sun, Moon } from 'lucide-react'
import { useTheme } from '@/lib/theme'
import { KronMark } from '@/components/kron-mark'
import {
  Sidebar as ShadcnSidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarTrigger,
  useSidebar,
} from '@/components/ui/sidebar'

const NAV = [
  { to: '/dashboard',   label: 'Dashboard',  icon: LayoutDashboard },
  { to: '/jobs',        label: 'Jobs',        icon: Calendar },
  { to: '/executions',  label: 'Executions',  icon: Activity },
]

export function Sidebar() {
  const { theme, toggle } = useTheme()
  const { open } = useSidebar()
  const router = useRouterState()
  const currentPath = router.location.pathname

  const isActive = (to: string) => currentPath === to || (to !== '/dashboard' && currentPath.startsWith(to))

  return (
    <ShadcnSidebar>
      <SidebarHeader className="flex items-center px-3">
        <div className={`flex min-w-0 flex-1 items-center ${open ? 'gap-2.5 px-2' : 'justify-center'}`}>
          <KronMark className="h-5 w-5 text-foreground shrink-0" />
          {open ? <span className="font-semibold text-[15px] tracking-tight">Kron</span> : null}
        </div>
        <SidebarTrigger className={open ? '' : 'absolute left-16 top-3 border border-border bg-card shadow-sm'} />
      </SidebarHeader>

      <SidebarContent className="space-y-0.5">
        {open ? (
          <p className="font-mono text-[10px] uppercase tracking-[0.22em] text-muted-foreground px-2 mb-3">
            Menu
          </p>
        ) : null}
        {NAV.map(({ to, label, icon: Icon }) => (
          <Link
            key={to}
            to={to}
            title={open ? undefined : label}
            aria-label={label}
            className={`flex items-center rounded-full text-sm transition-colors ${
              open ? 'gap-3 px-2.5 py-2' : 'h-10 justify-center px-0'
            } ${
              isActive(to)
                ? 'bg-primary text-primary-foreground font-medium'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted'
            }`}
          >
            <Icon className="w-4 h-4 shrink-0" />
            {open ? label : null}
          </Link>
        ))}

        <div className="pt-5">
          {open ? (
            <p className="font-mono text-[10px] uppercase tracking-[0.22em] text-muted-foreground px-2 mb-3">
              System
            </p>
          ) : null}
          <Link
            to="/settings"
            title={open ? undefined : 'Settings'}
            aria-label="Settings"
            className={`flex items-center rounded-full text-sm transition-colors ${
              open ? 'gap-3 px-2.5 py-2' : 'h-10 justify-center px-0'
            } ${
              isActive('/settings')
                ? 'bg-primary text-primary-foreground font-medium'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted'
            }`}
          >
            <Settings className="w-4 h-4 shrink-0" />
            {open ? 'Settings' : null}
          </Link>
        </div>
      </SidebarContent>

      <SidebarFooter className="space-y-1">
        <button
          onClick={toggle}
          title={open ? undefined : theme === 'dark' ? 'Light mode' : 'Dark mode'}
          aria-label={theme === 'dark' ? 'Light mode' : 'Dark mode'}
          className={`flex w-full items-center rounded-full text-sm text-muted-foreground transition-colors hover:bg-muted hover:text-foreground ${
            open ? 'gap-3 px-2.5 py-2' : 'h-10 justify-center px-0'
          }`}
        >
          {theme === 'dark'
            ? <Sun className="w-4 h-4 shrink-0" />
            : <Moon className="w-4 h-4 shrink-0" />
          }
          {open ? theme === 'dark' ? 'Light mode' : 'Dark mode' : null}
        </button>
      </SidebarFooter>
    </ShadcnSidebar>
  )
}
