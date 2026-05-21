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
  const { isMobile, open, setOpen } = useSidebar()
  const router = useRouterState()
  const currentPath = router.location.pathname

  const isActive = (to: string) => currentPath === to || (to !== '/dashboard' && currentPath.startsWith(to))
  const closeOnMobile = () => {
    if (isMobile) setOpen(false)
  }
  const labelClass = `overflow-hidden whitespace-nowrap transition-[max-width,opacity,transform] duration-300 ease-out ${
    open ? 'max-w-32 translate-x-0 opacity-100' : 'max-w-0 -translate-x-1 opacity-0'
  }`
  const sectionLabelClass = `font-mono overflow-hidden whitespace-nowrap px-2 text-[10px] uppercase tracking-[0.22em] text-muted-foreground transition-[max-height,margin,opacity] duration-300 ease-out ${
    open ? 'mb-3 max-h-5 opacity-100' : 'mb-0 max-h-0 opacity-0'
  }`

  return (
    <ShadcnSidebar>
      <SidebarHeader className="flex items-center px-3">
        <div className={`flex min-w-0 items-center gap-2.5 px-2 overflow-hidden transition-all duration-300 ease-out ${open ? 'flex-1 max-w-full opacity-100' : 'flex-none max-w-0 opacity-0 px-0'}`}>
          <KronMark className="h-5 w-5 text-foreground shrink-0" />
          <span className="font-semibold text-[15px] tracking-tight whitespace-nowrap">Kron</span>
        </div>
        <SidebarTrigger className={`transition-[margin] duration-300 ease-out ${open ? '' : 'mx-auto'}`} />
      </SidebarHeader>

      <SidebarContent className="space-y-0.5">
        <p className={sectionLabelClass}>Menu</p>
        {NAV.map(({ to, label, icon: Icon }) => (
          <Link
            key={to}
            to={to}
            title={open ? undefined : label}
            aria-label={label}
            onClick={closeOnMobile}
            className={`flex items-center rounded-full text-sm transition-colors ${
              open ? 'gap-3 px-2.5 py-2' : 'h-10 justify-center px-0'
            } ${
              isActive(to)
                ? 'bg-primary text-primary-foreground font-medium'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted'
            }`}
          >
            <Icon className="w-4 h-4 shrink-0" />
            <span className={labelClass}>{label}</span>
          </Link>
        ))}

        <div className="pt-5">
          <p className={sectionLabelClass}>System</p>
          <Link
            to="/settings"
            title={open ? undefined : 'Settings'}
            aria-label="Settings"
            onClick={closeOnMobile}
            className={`flex items-center rounded-full text-sm transition-colors ${
              open ? 'gap-3 px-2.5 py-2' : 'h-10 justify-center px-0'
            } ${
              isActive('/settings')
                ? 'bg-primary text-primary-foreground font-medium'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted'
            }`}
          >
            <Settings className="w-4 h-4 shrink-0" />
            <span className={labelClass}>Settings</span>
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
          <span className={labelClass}>{theme === 'dark' ? 'Light mode' : 'Dark mode'}</span>
        </button>
      </SidebarFooter>
    </ShadcnSidebar>
  )
}
