# ReadAgain Backend - Go/Fiber

Multi-tenant SaaS e-book platform backend built with Go and Fiber.

## Tech Stack

- **Framework:** Fiber v2
- **Database:** PostgreSQL 17
- **ORM:** GORM
- **Cache:** Redis
- **Auth:** JWT

## Prerequisites

- Go 1.21+
- PostgreSQL 17
- Redis 7+

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Start PostgreSQL:
```bash
docker-compose up -d
```

4. Run migrations:
```bash
go run cmd/api/main.go migrate
```

5. Start server:
```bash
go run cmd/api/main.go
```

Server runs on http://localhost:8000

## Project Structure

```
backend/
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/          # Configuration
│   ├── database/        # Database connections
│   ├── middleware/      # HTTP middlewares
│   ├── models/          # Data models
│   ├── handlers/        # HTTP handlers
│   ├── services/        # Business logic
│   ├── repository/      # Data access
│   └── utils/           # Utilities
├── pkg/                 # Public packages
├── migrations/          # Database migrations
└── scripts/             # Utility scripts
```

## Development

Run with hot reload:
```bash
air
```

Run tests:
```bash
go test ./...
```

## API Documentation

API docs available at http://localhost:8000/swagger
