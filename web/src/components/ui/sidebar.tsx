import * as React from 'react'
import { PanelLeftClose, PanelLeftOpen } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'

type SidebarContextValue = {
  open: boolean
  isMobile: boolean
  setOpen: React.Dispatch<React.SetStateAction<boolean>>
  toggleSidebar: () => void
}

const SidebarContext = React.createContext<SidebarContextValue | null>(null)
const MOBILE_QUERY = '(max-width: 767px)'

function useSidebar() {
  const context = React.useContext(SidebarContext)

  if (!context) {
    throw new Error('useSidebar must be used within a SidebarProvider.')
  }

  return context
}

function SidebarProvider({
  defaultOpen = true,
  children,
}: React.PropsWithChildren<{ defaultOpen?: boolean }>) {
  const [isMobile, setIsMobile] = React.useState(() =>
    typeof window === 'undefined' ? false : window.matchMedia(MOBILE_QUERY).matches,
  )
  const [desktopOpen, setDesktopOpen] = React.useState(defaultOpen)
  const [mobileOpen, setMobileOpen] = React.useState(false)
  const open = isMobile ? mobileOpen : desktopOpen

  React.useEffect(() => {
    const media = window.matchMedia(MOBILE_QUERY)
    const update = () => setIsMobile(media.matches)

    update()
    media.addEventListener('change', update)

    return () => media.removeEventListener('change', update)
  }, [])

  const setOpen = React.useCallback<React.Dispatch<React.SetStateAction<boolean>>>(
    (value) => {
      const apply = (current: boolean) =>
        typeof value === 'function' ? value(current) : value

      if (isMobile) {
        setMobileOpen(apply)
      } else {
        setDesktopOpen(apply)
      }
    },
    [isMobile],
  )

  const value = React.useMemo(
    () => ({
      open,
      isMobile,
      setOpen,
      toggleSidebar: () => setOpen((current) => !current),
    }),
    [isMobile, open, setOpen],
  )

  return (
    <SidebarContext.Provider value={value}>
      <div
        className="group/sidebar-wrapper flex min-h-screen w-full bg-background"
        data-state={open ? 'expanded' : 'collapsed'}
        style={
          {
            '--sidebar-width': '14rem',
            '--sidebar-width-icon': '4rem',
            '--sidebar-offset': open ? 'var(--sidebar-width)' : 'var(--sidebar-width-icon)',
          } as React.CSSProperties
        }
      >
        {children}
      </div>
    </SidebarContext.Provider>
  )
}

function Sidebar({ className, ...props }: React.HTMLAttributes<HTMLElement>) {
  const { open } = useSidebar()

  return (
    <aside
      className={cn(
        'fixed inset-y-0 left-0 z-30 flex w-[var(--sidebar-width)] flex-col border-r border-border bg-card shadow-xl transition-[width,transform,box-shadow] duration-300 ease-out md:w-[var(--sidebar-offset)] md:translate-x-0 md:shadow-none',
        open ? 'translate-x-0' : '-translate-x-full',
        className,
      )}
      data-state={open ? 'expanded' : 'collapsed'}
      {...props}
    />
  )
}

function SidebarHeader({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return <div className={cn('h-14 shrink-0 border-b border-border', className)} {...props} />
}

function SidebarContent({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return <div className={cn('flex-1 overflow-y-auto px-3 py-5', className)} {...props} />
}

function SidebarFooter({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return <div className={cn('shrink-0 border-t border-border px-3 py-4', className)} {...props} />
}

function SidebarInset({ className, style, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  const { open, isMobile } = useSidebar()
  const offset = open ? 'var(--sidebar-width)' : 'var(--sidebar-width-icon)'

  return (
    <div
      className={cn(
        'flex min-w-0 flex-1 flex-col transition-[margin-left,width] duration-300 ease-out',
        className,
      )}
      style={
        isMobile
          ? style
          : {
              marginLeft: offset,
              width: `calc(100% - ${offset})`,
              ...style,
            }
      }
      {...props}
    />
  )
}

function SidebarOverlay({ className, ...props }: React.HTMLAttributes<HTMLButtonElement>) {
  const { isMobile, open, setOpen } = useSidebar()

  if (!isMobile || !open) return null

  return (
    <button
      type="button"
      aria-label="Close sidebar"
      className={cn('fixed inset-0 z-20 bg-background/70 backdrop-blur-sm transition-opacity md:hidden', className)}
      onClick={() => setOpen(false)}
      {...props}
    />
  )
}

function SidebarTrigger({ className, ...props }: React.ButtonHTMLAttributes<HTMLButtonElement>) {
  const { open, toggleSidebar } = useSidebar()
  const Icon = open ? PanelLeftClose : PanelLeftOpen

  return (
    <Button
      type="button"
      variant="ghost"
      size="icon"
      className={cn('h-8 w-8 shrink-0', className)}
      aria-label={open ? 'Collapse sidebar' : 'Expand sidebar'}
      title={open ? 'Collapse sidebar' : 'Expand sidebar'}
      onClick={toggleSidebar}
      {...props}
    >
      <Icon className="h-4 w-4" />
    </Button>
  )
}

export {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarInset,
  SidebarOverlay,
  SidebarProvider,
  SidebarTrigger,
  useSidebar,
}
