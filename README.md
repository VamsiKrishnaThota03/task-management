# AI-Powered Task Management System

A full-stack task manager with a Go (Gin) API, a Next.js frontend, JWT authentication, WebSocket updates and Google Gemini task suggestions.

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.9-008ECF)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-GORM-4169E1?logo=postgresql&logoColor=white)
![Next.js](https://img.shields.io/badge/Next.js-15-000000?logo=nextdotjs&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?logo=typescript&logoColor=white)
![Gemini](https://img.shields.io/badge/Google-Gemini-8E75B2?logo=googlegemini&logoColor=white)

## Overview

Users register and log in, create tasks and assign them to other users, and move tasks through `PENDING`, `IN_PROGRESS` and `COMPLETED`. While writing a task, a user can ask Google Gemini to improve the description and break it into subtasks. Task status changes are broadcast over a WebSocket, so connected task lists refresh without reloading the page. Admins can see every task, delete tasks and list users.

## Features

- **Authentication**: register/login with bcrypt-hashed passwords and 24-hour HS256 JWTs. Roles are `user` (default) and `admin`.
- **Tasks**: create tasks with a title, description, assignee (user ID), due date and priority. Users see the tasks assigned to them; admins see all tasks.
- **Status updates**: `PENDING` → `IN_PROGRESS` → `COMPLETED`
- **AI suggestions**: `POST /api/tasks/suggestions` sends the description to Gemini (`gemini-pro`) and returns 2–3 suggestions for improving it plus 3–4 subtasks
- **Real-time refresh**: a gorilla/websocket hub at `/ws` broadcasts "Task updated" messages, and the frontend `useWebSocket` hook (with auto-reconnect) refetches tasks
- **Admin**: delete tasks and list all users
- **Auto-migrations**: GORM creates the `users`, `tasks`, `task_tags` and `task_task_tags` tables on startup
- **Frontend**: Next.js App Router pages for login, register/signup, the task list, a new-task form with an AI-suggestion button, task editing and the profile, plus a light/dark mode switch in the navbar. Styled with Bootstrap 5 (loaded from a CDN).

## Tech stack

| Layer | Technologies |
|-------|--------------|
| Backend | Go 1.23, Gin, GORM + PostgreSQL driver, golang-jwt/jwt v5, gorilla/websocket, gin-contrib/cors, godotenv, `golang.org/x/crypto/bcrypt` |
| AI | Google Gemini via `github.com/google/generative-ai-go` |
| Frontend | Next.js 15 (App Router), React 19, TypeScript, Axios, Bootstrap 5 (Tailwind/PostCSS are configured but not used yet) |

## Project structure

```
task-management/
├── backend/
│   ├── cmd/server/main.go        # Entry point (listens on :8080)
│   ├── internal/
│   │   ├── auth/                 # JWT generation and verification
│   │   ├── config/               # .env / environment loading
│   │   ├── database/             # GORM connection + auto-migrations
│   │   ├── handlers/             # task.go, user.go
│   │   ├── middleware/           # JWT auth, admin role check
│   │   ├── models/               # User, Task, TaskTag
│   │   ├── routes/               # Gin router + CORS
│   │   └── websocket/            # Broadcast hub
│   ├── pkg/ai/ai.go              # Gemini client
│   ├── .env.example
│   └── Dockerfile
└── frontend/
    └── app/
        ├── api/                  # auth.ts, task.ts (Axios helpers)
        ├── components/           # Navbar, TaskForm, TaskList, UserProfile
        ├── hooks/useWebSocket.ts
        ├── login/ register/ signup/ profile/
        └── tasks/                # list, new, edit/[id]
```

## Getting started

### Prerequisites

- Go 1.23+
- Node.js 18+ and npm
- PostgreSQL
- A Google Gemini API key ([Google AI Studio](https://aistudio.google.com/))

### Backend

```bash
git clone https://github.com/VamsiKrishnaThota03/task-management.git
cd task-management/backend
cp .env.example .env      # then fill in the values
go mod tidy
go run ./cmd/server       # http://localhost:8080
```

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `DB_URI` | yes | PostgreSQL connection string | `postgres://<user>:<password>@localhost:5432/<db>` |
| `GEMINI_KEY` | yes | Google Gemini API key (the server exits if it is missing) | `<your-gemini-api-key>` |
| `JWT_SECRET` | – | Loaded into config, but token signing currently uses a constant in `internal/auth/auth.go` | `<a-long-random-string>` |

The server listens on port `8080`. CORS allows a single origin, set in `internal/routes/routes.go`. Change it to your frontend URL (for example `http://localhost:3000`) for local development.

### Frontend

```bash
cd task-management/frontend
npm install
npm run dev               # http://localhost:3000
```

| Variable | Description | Example |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | API base URL used by `app/api/task.ts` | `http://localhost:8080` |
| `NEXT_PUBLIC_WS_URL` | WebSocket URL default in `useWebSocket.ts` | `ws://localhost:8080/ws` |

Most components currently call `http://localhost:8080` directly.

| Script | Description |
|--------|-------------|
| `npm run dev` | Start the Next.js dev server |
| `npm run build` | Production build |
| `npm start` | Serve the production build |

## API

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| `POST` | `/register` | none | Create a user (`username`, `password`, optional `role`) |
| `POST` | `/login` | none | Returns `{ token, role }` |
| `GET` | `/ws` | none | WebSocket; receives task-update broadcasts |
| `GET` | `/api/profile` | JWT | Current user's profile |
| `GET` | `/api/tasks` | JWT | Tasks assigned to the user (all tasks for admins) |
| `POST` | `/api/tasks` | JWT | Create a task (`title`, `description`, `assignedTo`, `dueDate`, `priority`) |
| `PUT` | `/api/tasks/:id/status` | JWT | Update status (`PENDING`, `IN_PROGRESS`, `COMPLETED`) and broadcast it |
| `POST` | `/api/tasks/suggestions` | JWT | Gemini suggestions for a `description` |
| `DELETE` | `/api/tasks/:id` | admin | Delete a task |
| `GET` | `/api/users` | admin | List users |

Send protected requests with `Authorization: Bearer <token>`.

## Known limitations

- Several frontend calls use `/tasks...` without the `/api` prefix, and the edit page calls `GET`/`PUT /tasks/:id`, which the backend does not define yet
- The admin middleware expects the raw token (without `Bearer `) in the `Authorization` header
- `main.go` logs the loaded configuration values on startup
- The Dockerfile builds the module root rather than `./cmd/server`

## Future enhancements

- Docker / Kubernetes deployment
- Slack or Discord notifications for task updates
- AI-assisted task prioritisation

## Contact

Thota Vamsi Krishna — [GitHub](https://github.com/VamsiKrishnaThota03)
