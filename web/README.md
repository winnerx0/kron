# Kron Web UI

A modern web interface for Kron - a job scheduling system. Built with TanStack Start, React, and shadcn UI components.

## Features

- **Job Management**: Create, update, and delete scheduled jobs
- **Execution Tracking**: View job execution history and status
- **Real-time Updates**: Auto-refreshing execution list (5s interval)
- **Minimal Design**: Clean, focused UI using shadcn components

## Prerequisites

- Node.js 18+
- npm or yarn
- Kron backend running on `http://localhost:5000`

## Setup

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm run dev
```

The app will be available at `http://localhost:3000`

## Configuration

Set the API URL via environment variable:
```bash
API_URL=http://your-kron-api:5000 npm run dev
```

## Building for Production

```bash
npm run build
npm run preview
```

## Project Structure

```
src/
  ├── components/
  │   ├── ui/              # shadcn UI components
  │   ├── jobs-list.tsx    # Jobs management UI
  │   ├── executions-list.tsx  # Execution history
  │   └── create-job-dialog.tsx # Job form
  ├── routes/              # TanStack Router routes
  ├── lib/
  │   ├── utils.ts         # Helper utilities
  │   └── api.ts           # Kron API client
  └── index.css            # Global styles & theme
```

## API Endpoints

The UI communicates with the following Kron API endpoints:

- `GET /job/all` - List all jobs
- `POST /job/create` - Create a new job
- `PUT /job/{jobID}` - Update a job
- `DELETE /job/{jobID}` - Delete a job
- `GET /execution/all` - List all executions

## Theme

The UI uses a minimal light theme with:
- Clean typography
- Subtle gray accents
- Red for destructive actions
- Accessible color contrast
