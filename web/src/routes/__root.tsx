import { useEffect } from 'react'
import { Outlet, useLocation, useNavigate } from '@tanstack/react-router'
import posthog from 'posthog-js'
import { PostHogProvider } from 'posthog-js/react'
import { ThemeProvider } from '@/lib/theme'
import { Sidebar } from '@/components/sidebar'
import { SidebarInset, SidebarOverlay, SidebarProvider, SidebarTrigger, useSidebar } from '@/components/ui/sidebar'
import { getAccessToken } from '@/lib/auth'

const PUBLIC_PATHS = new Set(['/', '/login', '/callback'])

function MobileSidebarTrigger() {
  const { isMobile, open } = useSidebar()

  if (!isMobile || open) return null

  return (
    <SidebarTrigger className="fixed left-4 top-4 z-40 border border-border bg-card shadow-sm md:hidden" />
  )
}

export function RootLayout() {
  const location = useLocation()
  const navigate = useNavigate()
  const isPublic = PUBLIC_PATHS.has(location.pathname)

  useEffect(() => {
    if (isPublic) return
    if (!getAccessToken()) navigate({ to: '/login', replace: true })
  }, [isPublic, location.pathname, navigate])

  if (isPublic) {
    return (
      <PostHogProvider client={posthog}>
        <ThemeProvider>
          <Outlet />
        </ThemeProvider>
      </PostHogProvider>
    )
  }

  return (
    <ThemeProvider>
      <SidebarProvider>
        <MobileSidebarTrigger />
        <SidebarOverlay />
        <Sidebar />
        <SidebarInset className="h-screen overflow-hidden">
          <main className="flex-1 overflow-y-auto p-4 pt-16 md:p-8">
            <Outlet />
          </main>
        </SidebarInset>
      </SidebarProvider>
    </ThemeProvider>
  )
}
