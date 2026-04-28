import { Outlet } from '@tanstack/react-router'
import { ThemeProvider } from '@/lib/theme'
import { Sidebar } from '@/components/sidebar'

export function RootLayout() {
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
