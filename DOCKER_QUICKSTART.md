# Quick Start Guide - Docker Deployment

## What's Included

This Docker setup provides a complete, production-ready deployment for Ubuntu 24.04:

- **Backend**: Go API server with Fiber framework
- **Frontend**: Svelte 5 application served by Nginx with reverse proxy
- **Database**: PostgreSQL 16 with automatic schema initialization
- **Networking**: Isolated Docker network for secure communication
- **Health Checks**: Automatic health monitoring for all services
- **Volumes**: Persistent PostgreSQL data storage

## File Structure

```
project-management/
├── docker-compose.yml              # Development configuration
├── docker-compose.prod.yml         # Production configuration
├── docker-start.sh                 # Interactive startup script
├── .env.docker                     # Environment template
├── DOCKER_DEPLOYMENT.md            # Comprehensive guide
├── backend/
│   ├── Dockerfile                  # Multi-stage Go build
│   └── .dockerignore               # Build context cleanup
├── frontend/
│   ├── Dockerfile                  # Nginx + Svelte build
│   ├── nginx.conf                  # Production Nginx config
│   └── .dockerignore               # Build context cleanup
└── schema.sql                       # Database initialization
```

## Quick Start (5 minutes)

### 1. Prerequisites Check

```bash
docker --version    # Should be 24.0+
docker compose version  # Should be 2.20+
```

### 2. One-Command Start

```bash
cd /path/to/project-management
./docker-start.sh
# Select option 1 to start services
```

### 3. Access Application

- **Frontend**: http://localhost
- **API**: http://localhost:3000
- **Database**: localhost:5432 (only from host)

## Configuration

### Development (.env)

```bash
cp .env.docker .env
# Edit as needed for development
```

### Production

```bash
# Generate secure values
openssl rand -base64 32  # For DB_PASSWORD
openssl rand -base64 64  # For JWT_SECRET

# Create production .env
cp .env.docker .env
# Update all values, especially passwords and CORS_ORIGIN
```

## Common Commands

### Start Services
```bash
# Using script (recommended)
./docker-start.sh

# Direct command
docker compose up -d
```

### View Status
```bash
docker compose ps
```

### Check Logs
```bash
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f postgres
```

### Stop Services
```bash
docker compose down
```

### Reset Everything
```bash
docker compose down -v  # Includes removing database
docker compose up -d
```

### Execute Commands

```bash
# Backend shell
docker compose exec backend sh

# Database connection
docker compose exec postgres psql -U postgres -d project_management

# Run backend script
docker compose exec backend go run scripts/reset_admin_password.go
```

## Database

### Automatic Initialization

The database is automatically set up on first run with:
- Base schema (schema.sql)
- Entity enhancements
- Authentication tables
- Comments functionality

### Backup Database

```bash
docker compose exec postgres pg_dump -U postgres project_management > backup.sql
```

### Restore Database

```bash
cat backup.sql | docker compose exec -T postgres psql -U postgres project_management
```

## Troubleshooting

### Can't access frontend
```bash
docker compose logs frontend
# Check that nginx container is running
docker compose ps
```

### API connection errors
```bash
# Test API from within frontend container
docker compose exec frontend wget -v http://backend:3000/api/health

# Check backend logs
docker compose logs backend
```

### Database issues
```bash
# Check database logs
docker compose logs postgres

# Verify database is running
docker compose ps postgres

# Check database connection
docker compose exec backend wget -O- http://localhost:3000/api/health
```

### Port already in use
```bash
# Check what's using the port
sudo lsof -i :80    # For frontend
sudo lsof -i :3000  # For API
sudo lsof -i :5432  # For database

# Change ports in docker-compose.yml if needed
```

## Production Deployment

For production on Ubuntu 24.04:

1. **Use `docker-compose.prod.yml`**
   ```bash
   docker compose -f docker-compose.prod.yml up -d
   ```

2. **Key differences from development**:
   - Resource limits configured
   - Database port restricted to localhost only
   - No port exposure for internal services
   - Production-grade PostgreSQL settings
   - Security optimizations enabled

3. **Set up SSL/TLS**:
   - Install Certbot
   - Generate certificates with Let's Encrypt
   - Update nginx.conf with SSL configuration
   - Uncomment SSL volume mounts in docker-compose

4. **See DOCKER_DEPLOYMENT.md for**:
   - Detailed installation steps
   - SSL/TLS configuration
   - Monitoring and maintenance
   - Backup strategies
   - Performance tuning
   - Security best practices

## Network Architecture

```
┌─────────────────────────────────────┐
│     Ubuntu Host 24.04                │
├─────────────────────────────────────┤
│  ┌──────────────────────────────┐  │
│  │   Docker Network: app-network │  │
│  │  ┌──────────────────────────┐ │  │
│  │  │  Frontend (Nginx)         │ │  │
│  │  │  Port: 80/443            │ │  │
│  │  │  ├─→ Backend (Proxy)     │ │  │
│  │  └──────────────────────────┘ │  │
│  │              ↓                 │  │
│  │  ┌──────────────────────────┐ │  │
│  │  │  Backend (Go/Fiber)      │ │  │
│  │  │  Port: 3000 (internal)   │ │  │
│  │  │  ├─→ PostgreSQL          │ │  │
│  │  └──────────────────────────┘ │  │
│  │              ↓                 │  │
│  │  ┌──────────────────────────┐ │  │
│  │  │  PostgreSQL              │ │  │
│  │  │  Port: 5432 (internal)   │ │  │
│  │  │  Volume: postgres_data   │ │  │
│  │  └──────────────────────────┘ │  │
│  └──────────────────────────────┘  │
└─────────────────────────────────────┘
```

## Support & Further Help

- **Full documentation**: See [DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md)
- **Docker Docs**: https://docs.docker.com/
- **Docker Compose Docs**: https://docs.docker.com/compose/
- **Project Docs**: Check [AGENTS.md](AGENTS.md) for architecture details

## Default Credentials

After first startup:
- **Admin Email**: admin@example.com
- **Admin Password**: Admin123!

⚠️ **IMPORTANT**: Change these immediately after first login!

---

**Last Updated**: December 31, 2025
**Tested On**: Ubuntu 24.04 LTS with Docker 24.0+
