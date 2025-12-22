# Amazon Q Development Rules - ReadAgain Go/Fiber Project

## General Rules

- Always commit with single line commit message
- All imports at top of file, grouped: stdlib, external, internal
- Implement proper separation of concerns
- Follow Go best practices and idioms
- Use meaningful variable and function names
- Keep functions small and focused (max 50 lines)
- No code duplication - extract to functions/packages

## Code Organization

### Package Structure
- `cmd/` - Application entry points
- `internal/` - Private application code
- `pkg/` - Public reusable packages
- `migrations/` - Database migrations
- `scripts/` - Utility scripts

### File Naming
- Use snake_case for file names: `user_service.go`
- Test files: `user_service_test.go`
- One main type per file
- Group related functionality in same package

## Go Best Practices

### Imports
```
import (
    // Standard library
    "context"
    "fmt"
    
    // External packages
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
    
    // Internal packages
    "readagain/internal/models"
    "readagain/internal/services"
)
```

### Error Handling
- Always check errors immediately
- Return errors, don't panic
- Wrap errors with context: `fmt.Errorf("failed to create user: %w", err)`
- Use custom error types for business logic errors
- Log errors before returning

### Naming Conventions
- Use camelCase for unexported: `userService`
- Use PascalCase for exported: `UserService`
- Interface names: `UserRepository`, `EmailSender`
- Avoid stuttering: `user.User` not `user.UserModel`
- Boolean variables: `isActive`, `hasPermission`

### Functions
- Keep functions small (max 50 lines)
- Single responsibility principle
- Return early on errors
- Use named return values for documentation only
- Avoid naked returns

### Structs
- Group related fields together
- Add JSON/GORM tags on same line
- Use pointers for optional fields
- Embed structs for composition
- Add validation tags

### Interfaces
- Keep interfaces small (1-3 methods)
- Define interfaces where used, not where implemented
- Accept interfaces, return structs
- Use `io.Reader`, `io.Writer` when applicable

### Concurrency
- Use channels for communication
- Use sync.WaitGroup for goroutine coordination
- Always handle context cancellation
- Avoid shared memory, use channels
- Close channels in sender, not receiver

### Context
- First parameter in functions: `ctx context.Context`
- Pass context through call chain
- Use context for cancellation and timeouts
- Don't store context in structs

## Project-Specific Rules

### Models
- Use GORM tags for database mapping
- Add JSON tags for API responses
- Embed `BaseModel` for common fields
- Use pointers for optional foreign keys
- Add validation tags from validator package

### Handlers
- Keep handlers thin (5-15 lines)
- Delegate to services
- Parse request, call service, return response
- Use fiber.Map for JSON responses
- Return proper HTTP status codes

### Services
- Business logic only
- No HTTP concerns
- Return domain errors
- Use repository for data access
- Keep transactions in service layer

### Repository
- Database operations only
- Return models and errors
- Use GORM best practices
- Implement pagination
- Use preloading for relationships

### Middleware
- Single responsibility
- Call `c.Next()` to continue chain
- Return error to stop chain
- Store data in `c.Locals()`
- Keep middleware stateless

## Database

### GORM Best Practices
- Use `db.WithContext(ctx)` for all queries
- Preload relationships: `db.Preload("Role").Find(&users)`
- Use transactions for multiple operations
- Add indexes for foreign keys and search fields
- Use soft deletes: `gorm.DeletedAt`

### Migrations
- One migration per change
- Use timestamp prefix: `20250109_create_users_table.go`
- Always provide up and down migrations
- Test migrations before committing

## API Design

### Endpoints
- Use RESTful conventions
- Plural nouns: `/api/books`, `/api/users`
- Use HTTP methods correctly: GET, POST, PUT, DELETE
- Version API: `/api/v1/books`
- Use query params for filtering: `?status=active&page=1`

### Request/Response
- Use JSON for all requests/responses
- Validate all inputs
- Return consistent error format
- Include pagination metadata
- Use proper HTTP status codes

### Status Codes
- 200: Success
- 201: Created
- 400: Bad Request (validation error)
- 401: Unauthorized (not authenticated)
- 403: Forbidden (not authorized)
- 404: Not Found
- 500: Internal Server Error

## Security

### Authentication
- Use JWT for authentication
- Store tokens in Authorization header: `Bearer <token>`
- Validate tokens in middleware
- Use bcrypt for password hashing (cost 14)
- Implement token refresh mechanism

### Authorization
- Check permissions in middleware
- Use RBAC (Role-Based Access Control)
- Validate user owns resource
- Log security events
- Rate limit sensitive endpoints

### Input Validation
- Validate all user inputs
- Sanitize HTML inputs
- Use validator package
- Check file types and sizes
- Validate business rules

## Testing

### Unit Tests
- Test file next to source: `user_service_test.go`
- Use table-driven tests
- Mock external dependencies
- Test error cases
- Aim for 80%+ coverage

### Test Structure
- Arrange, Act, Assert pattern
- Use `t.Run()` for subtests
- Clean up resources in `defer`
- Use test fixtures
- Parallel tests when possible

## Performance

### Optimization
- Use connection pooling
- Implement caching (Redis)
- Use indexes on database
- Paginate large result sets
- Use goroutines for concurrent operations
- Profile before optimizing

### Caching
- Cache frequently accessed data
- Use Redis for distributed cache
- Set appropriate TTL
- Invalidate cache on updates
- Cache at service layer

## Logging

### Log Levels
- Debug: Development info
- Info: General information
- Warn: Warning messages
- Error: Error messages
- Fatal: Critical errors (exits app)

### What to Log
- All errors with context
- Authentication attempts
- Important business events
- Performance metrics
- Security events

### What NOT to Log
- Passwords or secrets
- Full credit card numbers
- Personal identifiable information (PII)
- Large payloads

## Documentation

### Code Comments
- Comment exported functions and types
- Explain "why", not "what"
- Use godoc format
- Keep comments up to date
- Add TODO/FIXME with ticket number

### API Documentation
- Use Swagger/OpenAPI
- Document all endpoints
- Include request/response examples
- Document error responses
- Keep docs in sync with code

## Git Workflow

### Commits
- Single line commit message
- Present tense: "Add user service" not "Added"
- Be specific: "Add JWT authentication" not "Update auth"
- Reference ticket if applicable: "Add user service (#123)"

### Branches
- `main` - Production code
- `develop` - Development branch
- `feature/feature-name` - New features
- `fix/bug-name` - Bug fixes
- `hotfix/critical-fix` - Production hotfixes

## Code Review Checklist

- [ ] All imports at top, properly grouped
- [ ] No code duplication
- [ ] Error handling implemented
- [ ] Input validation added
- [ ] Tests written
- [ ] Documentation updated
- [ ] No hardcoded values
- [ ] Proper separation of concerns
- [ ] Security considerations addressed
- [ ] Performance considerations addressed

## Common Patterns

### Service Pattern
```
type UserService struct {
    repo UserRepository
    jwt  *JWTService
}

func NewUserService(repo UserRepository, jwt *JWTService) *UserService {
    return &UserService{repo: repo, jwt: jwt}
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Validate
    // Business logic
    // Call repository
    // Return result
}
```

### Repository Pattern
```
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uint) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uint) error
}
```

### Handler Pattern
```
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    user, err := h.service.CreateUser(c.Context(), req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.Status(201).JSON(user)
}
```

## Environment Variables

- Never commit `.env` file
- Use `.env.example` as template
- Validate required env vars on startup
- Use meaningful names: `DB_HOST` not `HOST`
- Group related vars: `DB_*`, `REDIS_*`, `JWT_*`

## Dependencies

- Use `go mod tidy` regularly
- Pin versions in production
- Review dependencies before adding
- Keep dependencies updated
- Minimize external dependencies

## Deployment

- Use Docker for containerization
- Multi-stage builds for smaller images
- Health check endpoints
- Graceful shutdown
- Environment-based configuration

---

**Remember:** Write clean, maintainable, idiomatic Go code. When in doubt, check the official Go documentation and follow the Go proverbs.
