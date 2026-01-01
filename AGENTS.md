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
│   ├── services/     # Business logic layer
│   ├── repositories/ # Data access (pgx)
│   ├── models/       # Domain entities
│   ├── routes/       # Route definitions
│   └── main.go       # Entry point
├── frontend/
│   └── src/
│       ├── components/ # Svelte 5 components (13)
│       ├── stores/     # State management
│       └── lib/        # API client
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
- Svelte components: 13+ total (includes dashboard components)
- Persian language support via `jalali-moment` library
- Default admin: `admin@example.com` / `Admin123!` (change after first login)
- **New in 005-dashboard**: Meetings entity, dashboard aggregations, 30-second auto-refresh pattern
- **Dashboard tech**: Single `/api/dashboard` endpoint for efficiency, client-side Jalali date formatting, generated avatars with initials
