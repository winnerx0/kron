import { createRootRoute, createRoute, Router } from '@tanstack/react-router'
import { RootLayout } from './routes/__root'
import { DashboardPage } from './routes/index'
import { JobsPage } from './routes/jobs'
import { ExecutionsPage } from './routes/executions'
import { SettingsPage } from './routes/settings'

const rootRoute = createRootRoute({ component: RootLayout })

const indexRoute = createRoute({ getParentRoute: () => rootRoute, path: '/', component: DashboardPage })
const jobsRoute = createRoute({ getParentRoute: () => rootRoute, path: '/jobs', component: JobsPage })
const executionsRoute = createRoute({ getParentRoute: () => rootRoute, path: '/executions', component: ExecutionsPage })
const settingsRoute = createRoute({ getParentRoute: () => rootRoute, path: '/settings', component: SettingsPage })

const routeTree = rootRoute.addChildren([indexRoute, jobsRoute, executionsRoute, settingsRoute])

export const router = new Router({ routeTree })

declare module '@tanstack/react-router' {
  interface Register { router: typeof router }
}
