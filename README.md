# FitByte

Track your progress, show your growth!

## Prerequisites

- Go 1.21+
- PostgreSQL (or Docker)

## Quick Start

```bash
# 1. Clone and setup
git clone [repo-url]
cd fitbyte
go mod tidy

# 2. Database setup
# Option A: Use Docker
docker compose up -d

# Option B: Use local PostgreSQL
# Create database: createdb fitbyte

# 3. Environment setup
cp .env.sample .env
# Edit .env with your database details

# 4. Run migrations
# Option A: Docker (automatic via compose)

# Option B: Local PostgreSQL
psql -d [db-name] -f migrations/001_create_users_table.sql
# Add more files as needed

# 5. Start server
go run cmd/main.go

# 6. Test it's working
curl http://localhost:[port]/api/v1/health
```

## API Endpoints

- `GET /api/v1/health` - Health check
- `GET /api/v1/users` - Get all users