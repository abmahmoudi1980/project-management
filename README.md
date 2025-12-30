# Project Management Application

A high-performance Project Management application built with Go (Fiber + pgx) and Svelte + Tailwind CSS.

## Tech Stack

- **Backend**: Go with Fiber framework and pgx PostgreSQL driver
- **Frontend**: Svelte with Vite and Tailwind CSS
- **Database**: PostgreSQL

## Features

- Project Management (CRUD with status tracking)
- Task Management within projects
- Priority levels (Low, Medium, High)
- Task completion tracking
- Time logging for tasks
- User Authentication (Registration, Login, Password Reset)
- Role-based Access Control (Admin/User)
- Responsive UI with Persian language support

## Project Structure

```
project-management/
├── backend/
│   ├── config/          # Database configuration
│   ├── models/          # Data models
│   ├── repositories/     # Database operations
│   ├── services/        # Business logic
│   ├── handlers/        # HTTP handlers
│   ├── routes/          # API routes
│   ├── main.go          # Application entry point
│   └── go.mod          # Go dependencies
├── frontend/
│   ├── src/
│   │   ├── components/  # Svelte components
│   │   ├── stores/      # State management
│   │   ├── lib/         # API client
│   │   ├── App.svelte    # Main component
│   │   └── main.js      # Entry point
│   ├── static/          # Static assets
│   ├── package.json      # Node dependencies
│   └── vite.config.js   # Vite configuration
└── schema.sql          # Database schema
```

## Setup Instructions

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Node.js 18 or higher
- npm or yarn

### Database Setup

1. Create a PostgreSQL database:
```sql
CREATE DATABASE project_management;
```

2. Run the schema:
```bash
psql -U postgres -d project_management -f schema.sql
```

3. Run the authentication migration:
```bash
cd backend && go run run_migration.go
```

### Authentication Setup

1. Copy the environment template:
```bash
cp .env.example .env
```

2. Generate a secure JWT secret:
```bash
openssl rand -base64 32
```
Add this value to `JWT_SECRET` in `.env`

3. Configure Gmail SMTP (for password reset):
- Enable 2FA on your Gmail account
- Generate an App Password: https://support.google.com/accounts/answer/185833
- Set `SMTP_USER` and `SMTP_PASSWORD` in `.env`

4. The default admin user will be created during migration:
- Email: `admin@example.com`
- Password: `Admin123!`
- Change this password after first login!

### Backend Setup

1. Navigate to the backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set environment variable (optional):
```bash
# Default connection: postgres://postgres:postgres@localhost:5432/project_management?sslmode=disable
export DATABASE_URL="postgres://user:password@localhost:5432/project_management?sslmode=disable"
```

4. Run the server:
```bash
go run main.go
```

The API will be available at `http://localhost:3000`

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Run the development server:
```bash
npm run dev
```

The application will be available at `http://localhost:5173`

## API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `GET /api/auth/me` - Get current user info
- `PUT /api/auth/me` - Update user profile
- `PUT /api/auth/me/password` - Change password
- `POST /api/auth/forgot-password` - Request password reset
- `POST /api/auth/reset-password` - Reset password with token

### User Management (Admin only)
- `GET /api/users` - List all users
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id/role` - Change user role
- `PUT /api/users/:id/activate` - Activate/deactivate user

### Projects
- `GET /api/projects` - List all projects
- `POST /api/projects` - Create new project
- `GET /api/projects/:id` - Get project by ID
- `PUT /api/projects/:id` - Update project
- `DELETE /api/projects/:id` - Delete project

### Tasks
- `GET /api/projects/:projectId/tasks` - List tasks for project
- `POST /api/projects/:projectId/tasks` - Create task in project
- `GET /api/tasks/:id` - Get task by ID
- `PUT /api/tasks/:id` - Update task
- `PATCH /api/tasks/:id/complete` - Toggle task completion
- `DELETE /api/tasks/:id` - Delete task

### Time Logs
- `GET /api/tasks/:taskId/timelogs` - List time logs for task
- `POST /api/tasks/:taskId/timelogs` - Create time log
- `GET /api/timelogs/:id` - Get time log by ID
- `DELETE /api/timelogs/:id` - Delete time log

## Database Schema

### Projects
- id (UUID, primary key)
- title (VARCHAR 255, NOT NULL)
- description (TEXT)
- status (VARCHAR 50, NOT NULL, default: 'active')
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)

### Tasks
- id (UUID, primary key)
- project_id (UUID, foreign key → projects)
- title (VARCHAR 255, NOT NULL)
- priority (VARCHAR 10, NOT NULL, default: 'Medium')
- completed (BOOLEAN, NOT NULL, default: false)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)

### Time Logs
- id (UUID, primary key)
- task_id (UUID, foreign key → tasks)
- date (DATE, NOT NULL)
- duration_minutes (INTEGER, NOT NULL)
- note (TEXT)
- created_at (TIMESTAMP)

## Development

### Backend Development

The backend follows a layered architecture:
1. **Models**: Data structures
2. **Repositories**: Database operations (CRUD)
3. **Services**: Business logic and validation
4. **Handlers**: HTTP request/response handling
5. **Routes**: API endpoint definitions

### Frontend Development

The frontend uses:
- **Svelte stores** for state management
- **Tailwind CSS** for styling
- **Vite** for development and building

### Adding New Features

**Backend:**
1. Add model to `models/`
2. Add repository functions to `repositories/`
3. Add service logic to `services/`
4. Add handlers to `handlers/`
5. Register routes in `routes/routes.go`

**Frontend:**
1. Add API functions to `src/lib/api.js`
2. Add store in `src/stores/`
3. Create Svelte component in `src/components/`
4. Import and use in `App.svelte`

## License

MIT
