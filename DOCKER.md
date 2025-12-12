# 🐳 Docker Deployment Guide

## Quick Start

### Start Full Stack (App + Database)

docker-compose up --buildAccess API at: http://localhost:8080

### Stop Services

docker-compose down---

## Prerequisites

- Docker installed
- Docker Compose installed

---

## Configuration

Environment variables are set in `docker-compose.yml`.

For production, update:
- `JWT_SECRET` - Strong random string
- `DB_PASSWORD` - Strong database password
- `GIN_MODE` - Set to `release`

---

## Development Setup

### Option 1: Docker PostgreSQL Only

# Start PostgreSQL
docker-compose -f docker-compose.dev.yml up -d

# Run Go app locally
go run cmd/api/main.go

# Stop PostgreSQL
docker-compose -f docker-compose.dev.yml down### Option 2: Full Docker Stack

docker-compose up---

## Available Commands

# Build image
docker-compose build

# Start services (detached)
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Restart services
docker-compose restart

# Execute command in container
docker exec -it taskapi-app sh

# Access PostgreSQL
docker exec -it taskapi-postgres psql -U postgres -d taskmanagement---

## Makefile Commands

make docker-build      # Build Docker image
make docker-up         # Start all services
make docker-down       # Stop all services
make docker-logs       # View logs
make docker-restart    # Restart services
make docker-clean      # Clean up everything
make docker-dev-up     # Start dev PostgreSQL only
make docker-shell      # Shell into API container
make docker-db-shell   # Shell into PostgreSQL---

## Production Deployment

### 1. Update Environment Variables

Edit `docker-compose.yml` or use `.env` file:

JWT_SECRET=production-secret-very-long-and-random
DB_PASSWORD=strong-production-password
GIN_MODE=release### 2. Build Production Image

docker build -t taskmanagement-api:v1.0.0 .### 3. Run with Docker Compose

docker-compose up -d### 4. Check Health

curl http://localhost:8080/health---

## Troubleshooting

### Port Already in Use

# Change port in docker-compose.yml
ports:
  - "8081:8080"  # Host:Container### Database Connection Failed

# Check PostgreSQL is healthy
docker-compose ps

# View logs
docker-compose logs postgres### Rebuild After Code Changes

docker-compose up --build---

## Performance Tips

1. **Multi-stage builds** - Smaller image size
2. **Alpine base** - Minimal footprint
3. **Health checks** - Auto-recovery
4. **Volumes** - Persistent data
5. **Networks** - Service isolation

---

## Security Best Practices

1. ✅ Don't hardcode secrets in Dockerfile
2. ✅ Use environment variables
3. ✅ Run as non-root user (can add)
4. ✅ Scan images for vulnerabilities
5. ✅ Keep base images updated

---

## Docker Hub Publishing

# Tag image
docker tag taskmanagement-api:latest yourusername/taskmanagement-api:v1.0.0

# Push to Docker Hub
docker push yourusername/taskmanagement-api:v1.0.0