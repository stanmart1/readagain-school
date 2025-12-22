# Go Backend Development Rules

## Logging

### Always Use Centralized Logger
- **NEVER** use `log.Println()`, `log.Printf()`, or `fmt.Println()` directly
- **ALWAYS** use the centralized logger from `internal/utils/logger.go`

```go
// ❌ WRONG
log.Println("User created")
fmt.Println("Error occurred")

// ✅ CORRECT
utils.InfoLogger.Println("User created")
utils.ErrorLogger.Println("Error occurred")
utils.DebugLogger.Println("Debug info")
```

### Logger Types
- `utils.InfoLogger` - General information, successful operations
- `utils.ErrorLogger` - Errors, failures, exceptions
- `utils.DebugLogger` - Debug information, development only

### When to Log
- **Info:** Successful operations, important events
- **Error:** All errors with context
- **Debug:** Development debugging, verbose info

## Code Organization

### Separation of Concerns
- Handlers: HTTP request/response only
- Services: Business logic
- Repository: Database operations
- Utils: Shared utilities

### File Structure
- One main type per file
- Group related functions in same file
- Keep files under 300 lines

## Import Order
1. Standard library
2. External packages
3. Internal packages

```go
import (
    "context"
    "fmt"
    
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
    
    "readagain/internal/models"
    "readagain/internal/utils"
)
```

## Error Handling

### Always Check Errors
```go
// ❌ WRONG
result, _ := someFunction()

// ✅ CORRECT
result, err := someFunction()
if err != nil {
    utils.ErrorLogger.Printf("Failed to execute: %v", err)
    return err
}
```

### Wrap Errors with Context
```go
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}
```

## Naming Conventions

### Variables
- camelCase for unexported: `userService`
- PascalCase for exported: `UserService`

### Functions
- Descriptive names: `CreateUser`, `GetUserByID`
- Boolean functions: `IsActive`, `HasPermission`

### Constants
- UPPER_SNAKE_CASE: `MAX_RETRY_COUNT`

## Database

### Always Use Context
```go
db.WithContext(ctx).Find(&users)
```

### Use Transactions for Multiple Operations
```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

## API Responses

### Consistent Format
```go
// Success
return c.JSON(fiber.Map{
    "data": result,
})

// Error
return c.Status(400).JSON(fiber.Map{
    "error": "Invalid request",
})
```

## Security

### Never Log Sensitive Data
- Passwords
- Tokens
- API keys
- Credit card numbers
- Personal information

## Testing

### Test File Naming
- `user_service_test.go` for `user_service.go`

### Use Table-Driven Tests
```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {"valid", "test", "test", false},
    {"empty", "", "", true},
}
```

## Commit Messages

### Single Line Format
- Present tense: "Add user service"
- Be specific: "Add JWT authentication middleware"
- No periods at end

---

**Remember:** Follow these rules consistently across the entire codebase.
