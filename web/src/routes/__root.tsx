import { useEffect } from 'react'
import { Outlet, useLocation, useNavigate } from '@tanstack/react-router'
import { ThemeProvider } from '@/lib/theme'
import { Sidebar } from '@/components/sidebar'
import { getAccessToken } from '@/lib/auth'

const PUBLIC_PATHS = new Set(['/login', '/callback'])

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
      <ThemeProvider>
        <Outlet />
      </ThemeProvider>
    )
  }

  return (
    <ThemeProvider>
      <div className="flex h-screen overflow-hidden bg-background">
        <Sidebar />
        <div className="flex-1 flex flex-col min-w-0 ml-56">
          <main className="flex-1 overflow-y-auto p-8">
            <Outlet />
          </main>
        </div>
      </div>
    </ThemeProvider>
  )
}
