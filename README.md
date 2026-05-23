# Kron

Kron is an authenticated HTTP job scheduler. It lets users create cron-based HTTP jobs, run them on demand, pause or disable them, stop active executions, and review execution history from a React dashboard.

## Features

- Google OAuth sign-in with JWT access tokens and rotating refresh tokens.
- Dashboard summary with total jobs, enabled and disabled jobs, execution counts, success rate, recent executions, and job status snapshots.
- Scheduled HTTP jobs using standard five-field cron expressions.
- Manual job controls for run now, stop, enable, disable, edit, and delete.
- HTTP methods for `GET`, `POST`, `PUT`, `PATCH`, and `DELETE`.
- Optional request headers and request body support for jobs that send payloads.
- Sensitive header encryption at rest with redacted secret handling.
- Execution tracking with status, start time, completion time, duration, filtering, pagination, and periodic refresh.
- User settings page with profile details and sign out.
- Swagger API documentation served by the backend.
- Docker Compose stack for Nginx, backend, frontend, and PostgreSQL.
- Optional PostHog analytics hooks in the frontend.

## Stack

- Backend: Go, Chi, GORM, PostgreSQL
- Frontend: React, Vite, TanStack Router, Tailwind CSS, lucide-react
- Auth: Google OAuth 2.0, JWT access tokens, rotating refresh tokens
- Runtime: Docker Compose with Nginx, backend, frontend, and Postgres services

## Requirements

- Docker and Docker Compose for the full local stack
- Go 1.26+ for local backend development
- Bun for local frontend development
- A Google OAuth client configured for the callback URL used by your environment
- Optional PostHog project credentials for frontend analytics

## Environment

Create a root `.env` file from the example:

```sh
cp .env.example .env
```

Then fill in the Google OAuth credentials and replace the placeholder secrets. Do not commit real `.env` values. The encryption key must be valid for AES-GCM initialization.

Environment values include:

- Database connection settings for PostgreSQL.
- `KRON_ENCRYPTION_KEY` for encrypting sensitive job header values.
- Google OAuth client credentials and redirect URL.
- JWT signing secrets for access and refresh tokens.
- `FRONTEND_URL` for OAuth callback redirects.
- Optional `VITE_POSTHOG_API_KEY` and `VITE_POSTHOG_API_HOST` for analytics.

## Run With Docker

```sh
docker compose up --build
```

Local services:

- App through Nginx: `http://localhost:80`
- Frontend dev server: `http://localhost:3000`
- Backend API: `http://localhost:5000`
- Postgres: `localhost:5432`
- Swagger UI: `http://localhost:80/swagger/index.html`

To stop services:

```sh
docker compose down
```

To reset the local database volume:

```sh
docker compose down -v
docker compose up --build
```

## Local Development

Backend:

```sh
go test ./...
go run ./cmd/kron
```

Frontend:

```sh
cd web
bun install
bun run dev
bun run type-check
```

The frontend expects `VITE_API_URL`; Docker sets it to `http://localhost:80/api`.

## Authentication Flow

1. The frontend sends users to `/oauth/google/login`.
2. The backend redirects to Google OAuth.
3. Google calls back to `/oauth/google/callback`.
4. The backend creates or finds the user, issues tokens, and redirects to the frontend callback page.
5. API requests use `Authorization: Bearer <access_token>`.
6. Expired access tokens are refreshed through `POST /api/auth/refresh`.

## App Workflows

- Dashboard: view job totals, execution totals, success rate, recent executions, and quick links into job and execution management.
- Jobs: create and edit HTTP jobs, choose a cron preset or custom cron expression, add headers, add request bodies for payload methods, run jobs immediately, stop active jobs, toggle enabled state, and delete jobs.
- Executions: review execution history with status filters, duration calculation, pagination, and periodic refresh for changing statuses.
- Settings: view the signed-in Google profile and sign out.

## API Overview

Public routes:

- `GET /health`
- `GET /swagger/*`
- `GET /oauth/google/login`
- `GET /oauth/google/callback`
- `POST /api/auth/refresh`

Authenticated routes:

- `GET /api/dashboard`
- `GET /api/user/me`
- `GET /api/job/all`
- `POST /api/job/create`
- `PUT /api/job/{jobID}`
- `DELETE /api/job/{jobID}`
- `POST /api/job/{jobID}/run`
- `POST /api/job/{jobID}/stop`
- `PATCH /api/job/{jobID}/status`
- `GET /api/execution/all`

Create job payload:

```json
{
  "name": "Daily health check",
  "description": "Ping the service every morning",
  "schedule": "0 9 * * *",
  "endpoint": "https://example.com/health",
  "method": "GET",
  "headers": {},
  "body": ""
}
```

Schedules use standard five-field cron expressions.

`GET /api/execution/all` supports pagination query parameters:

- `page`
- `pageSize`
- `jobID` for filtering by a specific job

## Project Layout

```text
cmd/kron/                 Backend entrypoint
internal/auth/            Token issuing and refresh logic
internal/http/            HTTP handlers and route wiring
internal/job/             Job domain, scheduler, execution logic
internal/execution/       Execution persistence and pagination
internal/dashboard/       Dashboard summary aggregation
internal/user/            User model, repository, DTOs
internal/oauth/           Google OAuth login and callback logic
internal/refresh_token/   Refresh token persistence
internal/secret/          Header secret encryption/redaction
web/                      React frontend
nginx/                    Development and production Nginx configs
```

## Verification

Run both checks before pushing changes:

```sh
go test ./...
cd web && bun run type-check
```
