# Docker Deployment Guide

This guide covers deploying the Project Management application on Ubuntu 24.04 using Docker and Docker Compose.

## Prerequisites

- Ubuntu 24.04 LTS
- Docker 24.0+
- Docker Compose 2.20+
- At least 2GB RAM and 5GB disk space
- A domain name (optional, for SSL)

## Installation on Ubuntu 24.04

### 1. Update System

```bash
sudo apt update && sudo apt upgrade -y
```

### 2. Install Docker

```bash
# Install dependencies
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common

# Add Docker GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

# Add Docker repository
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

# Install Docker
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Verify installation
docker --version
docker compose version
```

### 3. Add User to Docker Group (Optional, for non-root usage)

```bash
sudo usermod -aG docker $USER
newgrp docker
```

### 4. Clone or Copy Project

```bash
cd /home/ubuntu
git clone <your-repo-url> project-management
cd project-management
```

## Configuration

### 1. Create Environment File

```bash
cp .env.docker .env
nano .env
```

Update with your configuration:
- `DB_PASSWORD`: Strong database password
- `JWT_SECRET`: A long random string for JWT signing
- `CORS_ORIGIN`: Your domain name
- `TZ`: Timezone (default: Asia/Tehran)

### 2. Example .env for Production

```
DB_USER=postgres
DB_PASSWORD=$(openssl rand -base64 32)
DB_NAME=project_management
DB_PORT=5432

JWT_SECRET=$(openssl rand -base64 64)

CORS_ORIGIN=https://your-domain.com

SERVER_PORT=3000
TZ=Asia/Tehran
```

To generate secure values:

```bash
# Generate secure password
openssl rand -base64 32

# Generate secure JWT secret
openssl rand -base64 64
```

## Deployment

### 1. Build Images

```bash
docker compose build
```

This builds:
- `project-management-backend`: Go API service
- `project-management-frontend`: Svelte + Nginx service

### 2. Start Services

```bash
# Start in detached mode
docker compose up -d

# Monitor logs
docker compose logs -f

# Specific service logs
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f postgres
```

### 3. Verify Services

```bash
# Check running containers
docker compose ps

# Test API endpoint
curl http://localhost:3000/api/health

# Access frontend
curl http://localhost/
```

## Post-Deployment

### 1. Database Initialization

The database is automatically initialized from SQL migration files during first container startup:
- `schema.sql` - Base schema
- `001_enhance_entities.sql` - Entity enhancements
- `002_add_user_authentication.sql` - Auth tables
- `003_add_comments.sql` - Comments functionality

To verify:

```bash
docker compose exec postgres psql -U postgres -d project_management -c "\dt"
```

### 2. Reset Admin Password (if needed)

```bash
docker compose exec backend go run scripts/reset_admin_password.go
```

### 3. Seed Persian Sample Data (optional)

```bash
docker compose exec backend go run scripts/seed_persian_data.go
```

## SSL/TLS Setup (Recommended)

### Using Let's Encrypt with Certbot

```bash
# Install Certbot
sudo apt install -y certbot python3-certbot-nginx

# Obtain certificate
sudo certbot certonly --standalone -d your-domain.com

# Update nginx.conf with SSL
# Copy certificates to project:
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem ./frontend/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem ./frontend/key.pem
sudo chown $(whoami):$(whoami) ./frontend/cert.pem ./frontend/key.pem
```

### Update nginx.conf for HTTPS

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /etc/nginx/cert.pem;
    ssl_certificate_key /etc/nginx/key.pem;
    
    # ... rest of config
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

## Maintenance

### View Logs

```bash
# All services
docker compose logs --tail=100 -f

# Specific service
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f postgres
```

### Database Backup

```bash
# Backup database
docker compose exec postgres pg_dump -U postgres project_management > backup.sql

# Restore from backup
cat backup.sql | docker compose exec -T postgres psql -U postgres project_management
```

### Stop Services

```bash
# Graceful stop
docker compose stop

# Force stop
docker compose kill

# Remove containers (keeps volumes)
docker compose down

# Remove everything including volumes
docker compose down -v
```

### Restart Services

```bash
# Restart all
docker compose restart

# Restart specific service
docker compose restart backend
```

### Update Application

```bash
# Stop services
docker compose down

# Pull latest code
git pull origin main

# Rebuild images
docker compose build --no-cache

# Start services
docker compose up -d
```

## Resource Limits

To prevent excessive resource usage, add to docker-compose.yml:

```yaml
services:
  backend:
    # ... existing config
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
  
  frontend:
    # ... existing config
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 256M
        reservations:
          cpus: '0.5'
          memory: 128M
```

## Monitoring

### Check Container Health

```bash
docker compose ps

# Detailed health status
docker inspect project-management-api --format='{{json .State.Health}}' | jq
```

### Memory & CPU Usage

```bash
docker stats
docker stats project-management-api
```

## Troubleshooting

### Database Connection Failed

```bash
# Check database logs
docker compose logs postgres

# Test connection
docker compose exec backend wget -O- http://localhost:3000/api/health

# Check environment variables
docker compose exec backend env | grep DATABASE
```

### Frontend Can't Reach Backend

```bash
# Check network connectivity
docker compose exec frontend wget -v http://backend:3000/api/health

# Verify DNS resolution
docker compose exec frontend getent hosts backend
```

### Permission Denied Errors

```bash
# Fix permission issues
docker compose down
sudo chown -R $USER:$USER .

# Rebuild and restart
docker compose build --no-cache
docker compose up -d
```

### Out of Disk Space

```bash
# Clean up unused Docker resources
docker system prune -a --volumes

# Check disk usage
docker system df
```

## Performance Tuning

### PostgreSQL

Update `docker-compose.yml` environment:

```yaml
postgres:
  environment:
    POSTGRES_INITDB_ARGS: "-c max_connections=200 -c shared_buffers=256MB"
```

### Backend

Add resource limits and optimize:

```yaml
backend:
  environment:
    GOMAXPROCS: '2'
    GOMEMLIMIT: '480MiB'
```

### Frontend (Nginx)

Already configured with:
- Gzip compression
- Caching headers
- Worker process optimization

## Security Best Practices

1. **Change Default Passwords**: Update `.env` with strong credentials
2. **Use HTTPS**: Configure SSL certificates
3. **Restrict Database Access**: Use firewall rules
4. **Regular Backups**: Automate database backups
5. **Update Images**: Regularly rebuild with latest dependencies
6. **Environment Files**: Never commit `.env` to version control
7. **Network Isolation**: Use Docker networks (already configured)

## Backup Script

Create `backup.sh`:

```bash
#!/bin/bash
BACKUP_DIR="/backups/project-management"
mkdir -p $BACKUP_DIR

# Database backup
docker compose exec -T postgres pg_dump -U postgres project_management \
  > "$BACKUP_DIR/db_$(date +%Y%m%d_%H%M%S).sql"

# Keep only last 7 days
find $BACKUP_DIR -name "db_*.sql" -mtime +7 -delete

echo "Backup completed: $BACKUP_DIR"
```

Set up cron job:

```bash
# Edit crontab
crontab -e

# Add (daily at 2 AM):
0 2 * * * /home/ubuntu/backup.sh >> /var/log/backup.log 2>&1
```

## Support

For issues, check:
1. Application logs: `docker compose logs`
2. Configuration: Verify `.env` values
3. Network: Ensure ports are available
4. Resources: Check available disk/memory

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Docker Documentation](https://hub.docker.com/_/postgres)
- [Nginx Documentation](https://nginx.org/en/docs/)
