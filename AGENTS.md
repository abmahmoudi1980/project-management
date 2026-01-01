# Project Knowledge Base

**Generated:** 2026-01-01
**Commit:** dashboard-feature
**Branch:** 005-dashboard

## OVERVIEW

Go + Svelte 5 project management system with PostgreSQL backend. Layered Go architecture (handlers→services→repositories), Svelte 5 with runes, Tailwind CSS. Dashboard provides real-time project statistics, task management, and meeting tracking.

## STRUCTURE

```
project-management/
├── backend/            # Go API server (separate module)
│   ├── config/       # DB & JWT config
│   ├── handlers/     # HTTP request handlers
│   │   ├── dashboard_handler.go  # NEW - Dashboard endpoint
│   │   └── meeting_handler.go    # NEW - Meeting CRUD
│   ├── services/     # Business logic layer
│   │   ├── dashboard_service.go  # NEW - Dashboard aggregation
│   │   └── meeting_service.go    # NEW - Meeting logic
│   ├── repositories/ # Data access (pgx)
│   │   ├── dashboard_repository.go # NEW - Dashboard queries
│   │   └── meeting_repository.go   # NEW - Meeting data access
│   ├── models/       # Domain entities
│   │   ├── dashboard.go  # NEW - Dashboard DTOs
│   │   └── meeting.go    # NEW - Meeting models
│   ├── routes/       # Route definitions
│   ├── migration/    # SQL migrations
│   │   └── 005_add_dashboard_meetings.sql  # NEW - Meetings tables
│   └── main.go       # Entry point (updated for dashboard)
├── frontend/
│   └── src/
│       ├── components/ # Svelte 5 components (19)
│       │   ├── Avatar.svelte       # NEW - User avatar with initials
│       │   ├── Dashboard.svelte    # NEW - Main dashboard page
│       │   ├── MeetingCard.svelte  # NEW - Meeting widget
│       │   ├── ProjectCard.svelte  # NEW - Project card
│       │   ├── StatCard.svelte     # NEW - Statistics card
│       │   └── TaskListItem.svelte # NEW - Task item with checkbox
│       ├── stores/     # State management
│       └── lib/        # API client
│           └── api.js  # Updated with dashboard/meeting endpoints
├── specs/             # Feature specifications
└── .github/           # Agents & prompts
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add API endpoint | `backend/handlers/` → `backend/services/` → `backend/repositories/` | Layered pattern |
| Add Svelte component | `frontend/src/components/` | Use Svelte 5 runes |
| Add DB table | `backend/models/` → `backend/migration/` | Update schema.sql |
| Update auth logic | `backend/services/auth_service.go` | JWT + bcrypt |
| Add route | `backend/routes/routes.go` | Register handler |

## CONVENTIONS

### Go (Backend)
- **Architecture**: handlers → services → repositories (strict layering)
- **DB driver**: pgx (not database/sql)
- **Framework**: Fiber v2
- **DB**: PostgreSQL, UUID primary keys
- **Auth**: JWT (access 15min, refresh 7d), httpOnly cookies
- **Passwords**: bcrypt (cost factor 10)
- **Deviation**: `backend/go.mod` (module per directory, not root)

### Svelte 5 (Frontend)
- **Props**: `let { prop } = $props()` (not `export let`)
- **State**: `let x = $state(0)` (reactive)
- **Derived**: `let computed = $derived(x * 2)` (not `$:`)
- **Effects**: `$effect(() => { ... })` (not `$:` blocks)
- **Stores**: `$store` syntax unchanged
- **Language**: JavaScript (not TypeScript)

### Database
- **Migration**: SQL files in `backend/migration/`
- **Runner**: `go run ./cmd/migrate`
- **Naming**: snake_case columns, PascalCase structs

## ANTI-PATTERNS

- **NEVER skip layers** in backend (handlers must call services, services must call repositories)
- **NEVER use `export let`** in Svelte 5 (use `$props()`)
- **NEVER use `database/sql`** (use pgx directly)
- **NEVER expose DB connection** outside `config/` package
- **NEVER store passwords** (bcrypt hash only)

## COMMANDS

```bash
# Backend
cd backend
go run main.go              # Start server (port 3000)
go run ./cmd/migrate         # Run migrations
go build                   # Build binary

# Frontend
cd frontend
npm install                 # Install deps
npm run dev                 # Dev server (port 5173)
npm run build               # Production build
```

## NOTES

- Backend is separate Go module (`backend/go.mod`), import paths reflect this
- No test files currently (manual testing only)
- Svelte components: 19 total (includes 6 new dashboard components)
- Persian language support via `jalali-moment` library
- Default admin: `admin@example.com` / `Admin123!` (change after first login)
- **New in 005-dashboard**: 
  - Meetings entity with meeting_attendees junction table
  - Dashboard aggregations (statistics, projects, tasks, next meeting)
  - 30-second auto-refresh pattern on Dashboard component
  - Single `/api/dashboard` endpoint for efficiency
  - Client-side Jalali date formatting using `dateToJalaliString()`
  - Generated avatars with user initials and color-coded backgrounds
  - 6 new Svelte 5 components: Avatar, Dashboard, MeetingCard, ProjectCard, StatCard, TaskListItem
  - Dashboard API methods added to `api.js`
  - Dashboard route added to App.svelte navigation
