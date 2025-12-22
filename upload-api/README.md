# Upload API

Standalone file upload service for ReadnWin School platform.

## Structure

```
upload-api/
├── cmd/api/          # Application entry point
├── internal/
│   ├── handlers/     # HTTP handlers
│   ├── middleware/   # Middleware (auth, cors, etc)
│   └── utils/        # Utilities (validation, file handling)
└── uploads/          # File storage
    ├── covers/       # Book cover images
    ├── books/        # Book files (PDF, EPUB)
    └── thumbnails/   # Generated thumbnails
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
- `GET /files/:type/:filename` - Serve uploaded files
- `DELETE /files/:type/:filename` - Delete file
