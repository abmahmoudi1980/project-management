# Backend Knowledge

**Generated:** 2025-12-31
**Commit:** 39fed23

## OVERVIEW

Go 1.25+ API server with Fiber v2, PostgreSQL (pgx), layered architecture.

## STRUCTURE

```
backend/
├── config/          # Database & JWT configuration
├── handlers/        # HTTP handlers (5)
├── services/        # Business logic (6)
├── repositories/    # Data access (6)
├── models/          # Domain entities (6)
├── routes/          # Route definitions
├── middleware/      # Auth middleware
├── migration/       # SQL migrations (2)
├── scripts/         # Utility scripts
├── main.go          # API entry point
└── go.mod          # Go module (separate from root)
```

## WHERE TO LOOK

| Task                   | Location           | Notes                                    |
| ---------------------- | ------------------ | ---------------------------------------- |
| Add new model          | `models/`          | Use UUID primary keys, JSON tags         |
| Add repository methods | `repositories/`    | Use pgx directly, no ORMs                |
| Add business logic     | `services/`        | Call repositories, return domain objects |
| Add HTTP endpoint      | `handlers/`        | Receive request → call service → respond |
| Register new route     | `routes/routes.go` | Add to existing router                   |
| Add middleware         | `middleware/`      | Fiber middleware functions               |
| Database changes       | `migration/`       | SQL files with version prefix (XXX\_)    |

## CONVENTIONS

- **Layered strictly**: handlers → services → repositories (never skip)
- **pgx not database/sql**: Direct pgx v5 usage
- **Repository pattern**: CRUD methods in repos, business rules in services
- **Error handling**: Fiber error responses with proper HTTP status codes
- **UUID primary keys**: Use github.com/google/uuid
- **JSON tags**: All struct fields must have `json:"field_name"` tag
- **PascalCase structs, snake_case DB columns**

## ANTI-PATTERNS

- **NEVER call repository directly from handler** (must go through service)
- **NEVER return database errors to client** (log and return generic message)
- **NEVER use database/sql** (use pgx)
- **NEVER expose DB connection** outside config/ package
- **NEVER hardcode values** (use env vars via godotenv)

## COMMANDS

```bash
cd backend
go run main.go              # Start API (port 3000)
go run ./cmd/migrate         # Run DB migrations
go build                   # Build binary
go mod tidy                # Clean dependencies
```

## NOTES

- Separate Go module (module: `project-management`)
- Imports reflect backend/ subdirectory path
- 37 total Go files, 6 each in services/repositories/models/handlers
- No test files currently
