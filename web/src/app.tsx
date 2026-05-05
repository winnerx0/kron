import { createRootRoute, createRoute, Router } from '@tanstack/react-router'
import { RootLayout } from './routes/__root'
import { IndexRoute } from './routes/index'
import { DashboardPage } from './routes/dashboard'
import { JobsPage } from './routes/jobs'
import { ExecutionsPage } from './routes/executions'
import { SettingsPage } from './routes/settings'
import { CallbackPage } from './routes/callback'
import { LoginPage } from './routes/login'

const rootRoute = createRootRoute({ component: RootLayout })

const indexRoute     = createRoute({ getParentRoute: () => rootRoute, path: '/',           component: IndexRoute })
const dashboardRoute = createRoute({ getParentRoute: () => rootRoute, path: '/dashboard',  component: DashboardPage })
const jobsRoute      = createRoute({ getParentRoute: () => rootRoute, path: '/jobs',       component: JobsPage })
const executionsRoute= createRoute({ getParentRoute: () => rootRoute, path: '/executions', component: ExecutionsPage })
const settingsRoute  = createRoute({ getParentRoute: () => rootRoute, path: '/settings',   component: SettingsPage })
const callbackRoute  = createRoute({ getParentRoute: () => rootRoute, path: '/callback',   component: CallbackPage })
const loginRoute     = createRoute({ getParentRoute: () => rootRoute, path: '/login',      component: LoginPage })

const routeTree = rootRoute.addChildren([
  indexRoute, dashboardRoute, jobsRoute, executionsRoute, settingsRoute, callbackRoute, loginRoute,
])

export const router = new Router({ routeTree })

declare module '@tanstack/react-router' {
  interface Register { router: typeof router }
}
