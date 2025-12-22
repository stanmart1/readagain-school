# Upload API

Standalone file upload service for ReadnWin School platform.
Files are stored directly in Coolify persistent storage.

## Structure

```
upload-api/
├── cmd/api/          # Application entry point
└── internal/
    ├── handlers/     # HTTP handlers
    └── utils/        # Utilities (validation, file handling)
```

## Setup

```bash
cp .env.example .env
go mod download
go run cmd/api/main.go
```

## Endpoints

- `POST /upload/cover` - Upload book cover image
- `POST /upload/book` - Upload book file
- `GET /files/:filename` - Serve uploaded files
- `DELETE /files/:filename` - Delete file

## Deployment

Files are stored in Coolify persistent storage at `/data/storage`.
No authentication required - public upload endpoint.
