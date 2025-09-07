# FitByte

Track your progress, show your growth!

## Prerequisites

- Go 1.25+
- Docker and Docker Compose

## Quick Start

### 1. Clone Repository
```bash
git clone [repo-url]
cd fitbyte
```

### 2. Environment Configuration
```bash
cp .env.sample .env
# Edit .env with your configuration including port settings
```

### 3. Start All Services
```bash
docker compose up -d
```

This will start and configure:
- PostgreSQL database with automatic migrations
- Redis for caching
- MinIO for file storage
- Application server

### 4. Verify Installation
```bash
curl http://localhost:8080/api/v1/healthz
```

## API Endpoints

### Authentication
- `POST /api/v1/login` - User login
- `POST /api/v1/register` - User registration

### User Management
- `GET /api/v1/user` - Get user profile (requires auth)
- `PATCH /api/v1/user` - Update user profile (requires auth)

### Activity Management
- `GET /api/v1/activity` - Get user activities with filtering (requires auth)
- `POST /api/v1/activity` - Create new activity (requires auth)
- `PATCH /api/v1/activity/:activityId` - Update activity (requires auth)
- `DELETE /api/v1/activity/:activityId` - Delete activity (requires auth)

### File Upload
- `POST /api/v1/file` - Upload profile image (requires auth)

### System
- `GET /api/v1/healthz` - Health check

## Development

### Port Configuration
Default ports are configured in `.env` file:
- Application: `HTTP_PORT=8080`
- PostgreSQL: `PG_PORT=5432` 
- Redis: `REDIS_PORT=6379`
- MinIO API: `MINIO_PORT=9000`
- MinIO Console: `MINIO_CONSOLE_PORT=9001`

Modify these values in your `.env` file if ports conflict with existing services.

### Running Locally
All dependencies are managed through Docker Compose. No need to install PostgreSQL, Redis, or MinIO separately.

### Database
Migrations run automatically when starting with Docker Compose. Database connection details are configured in `.env` file.

### File Storage
Uploaded files are stored in MinIO and accessible via the `MINIO_PUBLIC_ENDPOINT` configured in `.env`.

MinIO admin console available at the port specified by `MINIO_CONSOLE_PORT`.
Credentials are set via `MINIO_ACCESS_KEY` and `MINIO_SECRET_KEY` in `.env`.

### Cache
Redis connection details including password are configured in `.env` file.