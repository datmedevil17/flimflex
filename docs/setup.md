# Setup Guide

This guide provides instructions on how to set up and run the Micro-Flex microservices ecosystem.

## Prerequisites

- [Docker](https://www.docker.com/) and Docker Compose
- [Go](https://go.dev/) (1.21+)
- [Make](https://www.gnu.org/software/make/) (optional, but recommended)

## 1. Database Setup

The project uses PostgreSQL and Redis. You can spin them up easily using Docker Compose.

```bash
# Start Postgres and Redis
docker-compose up -d postgres redis

# Verify they are running
docker ps
```

### Database Credentials
| Service | Host | Port | User | Password | DB Name |
|---------|------|------|------|----------|---------|
| Postgres| localhost | 5432 | postgres | postgres | netflix |
| Redis   | localhost | 6379 | - | - | - |

## 2. Environment Configuration

Copy the example environment file if you haven't already (create one if needed based on `docker-compose.yml`):

```bash
# Example .env content (save as .env in root or specific service dirs if required)
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=netflix
```

Note: Most services in this project are configured to read from environment variables which are set in `docker-compose.yml`. For local development (`go run`), you may need to export these variables or use a `.env` file.

## 3. Running Services

You can run the services using Docker (recommended for full system) or locally using Go (for development).

### Option A: Using Docker (All Services)

To run the entire ecosystem including all microservices and the API Gateway:

```bash
# Build and start all services
make up

# Check logs
make logs
```

To stop everything:
```bash
make down
```

### Option B: Running Locally (Development)

You can run individual services locally using the provided Make commands. proper DB connection is required (Step 1).

**1. Auth Service**
```bash
make dev-auth
# Runs on port 50051
```

**2. User Service**
```bash
make dev-user
# Runs on port 50052
```

**3. Movie Service**
```bash
make dev-movie
# Runs on port 50053
```

**4. Upload Service**
```bash
make dev-upload
# Runs on port 50054
```

**5. Streaming Service**
```bash
make dev-streaming
# Runs on port 50055
```

**6. Recommendation Service**
```bash
make dev-recommendation
# Runs on port 50056
```

**7. API Gateway (Entry Point)**
```bash
make dev-gateway
# Runs on port 8080
```

## 4. Useful Commands

| Command | Description |
|---------|-------------|
| `make db-shell` | Open a psql shell inside the running Postgres container |
| `make db-reset` | Reset the database (down -v and up) |
| `make test` | Run all tests |
| `make proto` | Regenerate all Protobuf files |
