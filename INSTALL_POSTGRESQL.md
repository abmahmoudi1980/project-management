# PostgreSQL Installation Guide for Windows

## Option 1: Direct Download (Recommended)

1. Download PostgreSQL 16 from the official site:
   https://www.enterprisedb.com/downloads/postgres-postgresql/downloads/16/postgresql-16.1-2-windows-x64.exe

2. Run the installer and follow the prompts:
   - Set password for postgres user (remember this password)
   - Set port to 5432 (default)
   - Leave other settings as default

3. Complete the installation and restart if needed

## Option 2: Using Chocolatey (if available)

```bash
choco install postgresql
```

## Option 3: Using winget (retry)

```bash
winget install PostgreSQL.PostgreSQL.16
```

## After Installation

### 1. Create Database

Open "pgAdmin 4" (installed with PostgreSQL) or Command Prompt:

```cmd
psql -U postgres
CREATE DATABASE project_management;
\q
```

Or use pgAdmin:
1. Open pgAdmin 4
2. Connect to server (localhost:5432)
3. Right-click on Databases → Create → Database
4. Name: `project_management`
5. Click Save

### 2. Run Schema

From `D:\apps\project-management` directory:

```cmd
psql -U postgres -d project_management -f schema.sql
```

Or in pgAdmin:
1. Select `project_management` database
2. Click "Query Tool" (lightning icon)
3. Copy contents of schema.sql
4. Execute query

### 3. Update Backend Connection String

If you changed the default password, update `D:\apps\project-management\backend\config\database.go`:

Default connection string (if password is 'postgres'):
```
postgres://postgres:postgres@localhost:5432/project_management?sslmode=disable
```

If you set a custom password, replace `postgres:postgres` with `postgres:YOUR_PASSWORD`

### 4. Environment Variable (Optional)

Set DATABASE_URL environment variable:

**Windows (Command Prompt):**
```cmd
set DATABASE_URL=postgres://postgres:YOUR_PASSWORD@localhost:5432/project_management?sslmode=disable
```

**Windows (PowerShell):**
```powershell
$env:DATABASE_URL="postgres://postgres:YOUR_PASSWORD@localhost:5432/project_management?sslmode=disable"
```

## Verify Installation

```bash
psql --version
```

## Start Backend

```bash
cd D:\apps\project-management\backend
go mod tidy
go run main.go
```

## Start Frontend

```bash
cd D:\apps\project-management\frontend
npm install
npm run dev
```

## Access the Application

- Frontend: http://localhost:5173
- Backend API: http://localhost:3000
- API Docs: http://localhost:3000/api/projects

## Troubleshooting

### PostgreSQL service not starting
1. Open "Services" (services.msc)
2. Find "postgresql-x64-16"
3. Right-click → Start

### Connection refused
1. Check if PostgreSQL service is running
2. Verify port 5432 is not blocked
3. Check firewall settings

### psql command not found
1. Add PostgreSQL bin path to System PATH
2. Default location: `C:\Program Files\PostgreSQL\16\bin`
3. Or use full path: `"C:\Program Files\PostgreSQL\16\bin\psql.exe"`
