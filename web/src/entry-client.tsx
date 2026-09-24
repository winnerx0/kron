import React from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider } from '@tanstack/react-router'
import posthog from 'posthog-js'
import { router } from './app'
import './index.css'

posthog.init(import.meta.env.VITE_POSTHOG_API_KEY || '', {
  api_host: import.meta.env.VITE_POSTHOG_API_HOST || 'https://us.i.posthog.com',
  person_profiles: 'identified_only',
  capture_exceptions: true,
  defaults: '2026-01-30',
})

const rootEl = document.getElementById('root')
if (!rootEl) throw new Error('Root element not found')

ReactDOM.createRoot(rootEl).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
)
