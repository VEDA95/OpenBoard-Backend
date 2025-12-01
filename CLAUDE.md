# CLAUDE.md - OpenBoard Backend API

> **Comprehensive guide for AI assistants working on the OpenBoard Backend codebase**
> Last Updated: 2025-12-01

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture & Design Patterns](#architecture--design-patterns)
3. [Directory Structure](#directory-structure)
4. [Technology Stack](#technology-stack)
5. [Code Conventions](#code-conventions)
6. [Development Workflows](#development-workflows)
7. [Key Components](#key-components)
8. [Database & Migrations](#database--migrations)
9. [Authentication & Authorization](#authentication--authorization)
10. [API Design Patterns](#api-design-patterns)
11. [Error Handling](#error-handling)
12. [Testing Guidelines](#testing-guidelines)
13. [Common Tasks](#common-tasks)

---

## Project Overview

**OpenBoard Backend** is a REST API backend for a Kanban board application built with Go. The project serves as an open-source alternative to proprietary kanban board solutions.

### Key Features
- User authentication with session management
- Role-based access control (RBAC)
- Kanban boards with workspaces, lists, and cards
- Real-time updates via WebSocket
- File uploads for user thumbnails and card attachments
- Card comments, checklists, and labels
- Activity logging
- Email notifications

### Module Path
```go
module VEDA95/open_board/api
```

---

## Architecture & Design Patterns

### Layered Architecture

The codebase follows a **layered architecture** inspired by Clean Architecture principles:

```
┌─────────────────────────────────────┐
│   Presentation Layer (HTTP/WS)      │  ← Routes/Handlers
├─────────────────────────────────────┤
│   Business Logic Layer              │  ← Services
├─────────────────────────────────────┤
│   Data Access Layer                 │  ← Repositories
├─────────────────────────────────────┤
│   Domain Models                     │  ← Models/Entities
└─────────────────────────────────────┘
```

### Core Principles

1. **Separation of Concerns**: Each layer has distinct responsibilities
2. **Dependency Injection**: Manual DI via constructor functions
3. **Unidirectional Dependencies**: Layers only depend on layers below
4. **Single Responsibility**: One primary type per file

### Design Patterns Used

- **Repository Pattern**: Data access abstraction
- **Service Layer Pattern**: Business logic orchestration
- **Constructor Pattern**: `New*` functions for all components
- **Options Pattern**: `QueryOptions` for flexible queries
- **Middleware Pattern**: Request/response interception

---

## Directory Structure

```
/OpenBoard-Backend/
├── cmd/                          # Application entry points
│   ├── server/                   # Main API server
│   │   └── main.go              # Bootstrap & dependency injection
│   ├── loader/                   # Atlas schema loader
│   ├── admin/                    # Admin utilities
│   │   ├── seed/                # Database seeding
│   │   ├── createuser/          # User creation CLI
│   │   └── resetpassword/       # Password reset CLI
│   └── migrate/                  # Migration runner (deprecated)
│
├── internal/                     # Private application code
│   ├── http/                    # HTTP layer
│   │   ├── routes/              # HTTP handlers (controllers)
│   │   ├── middleware/          # HTTP middleware
│   │   ├── validators/          # Request validation
│   │   └── responses/           # Response formatting
│   ├── service/                 # Business logic layer
│   ├── db/                      # Database layer
│   │   ├── model/               # Domain models (entities)
│   │   └── repository/          # Data access layer
│   ├── auth/                    # Authentication utilities
│   ├── config/                  # Configuration management
│   ├── errors/                  # Error handling
│   ├── log/                     # Logging utilities
│   ├── email/                   # Email service
│   │   └── templates/           # Email templates
│   └── websocket/               # WebSocket management
│
├── migrations/                   # Database migrations (SQL)
├── deployment/                   # Deployment configurations
│   └── dev/                     # Development configs
│       ├── .air.linux.conf      # Hot reload config
│       └── docker-compose.postgres.yaml
├── env/                         # Environment files
│   ├── .env                     # Base configuration
│   └── .env.development         # Development overrides
├── build/                       # Build artifacts (gitignored)
├── docs/                        # Swagger documentation
│
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── Makefile                     # Build commands
├── atlas.hcl                    # Atlas migration config
└── README.md                    # Project documentation
```

### Key Locations Reference

| Purpose | Location |
|---------|----------|
| Add new endpoint | `/internal/http/routes/` |
| Add business logic | `/internal/service/` |
| Add database operations | `/internal/db/repository/` |
| Add/modify models | `/internal/db/model/` |
| Add middleware | `/internal/http/middleware/` |
| Add validation rules | `/internal/http/validators/` |
| Database migrations | `/migrations/` |
| Environment config | `/env/` |

---

## Technology Stack

### Core Framework
- **Go**: 1.24.6
- **Fiber**: v2.52.9 (Fast HTTP web framework)
- **GORM**: v1.31.0 (ORM)

### Database
- **PostgreSQL**: Primary database
- **pgx**: v5.7.6 (PostgreSQL driver)
- **Atlas**: Schema migration tool

### Authentication
- **Argon2id**: Password hashing (`alexedwards/argon2id`)
- **Session-based**: Custom session management

### Validation
- **go-playground/validator**: v10.25.0

### Logging
- **zerolog**: v1.33.0 (structured logging)
- **lumberjack**: v2.2.1 (log rotation)

### Development Tools
- **Air**: Hot reload for development
- **Swaggo**: API documentation generation

### Other Key Libraries
- **goccy/go-json**: Fast JSON encoding
- **gofrs/uuid**: UUID generation
- **wneessen/go-mail**: Email sending
- **fasthttp/websocket**: WebSocket support

---

## Code Conventions

### Naming Conventions

#### Files
- **Lowercase**, single word where possible
- Match primary struct/type: `user.go`, `board.go`
- One primary type per file

#### Packages
- **Singular nouns**: `service`, `repository`, `model`
- Descriptive: `middleware`, `validators`, `responses`

#### Functions/Methods
```go
// Exported: PascalCase
func NewUserService(...) *UserService
func FindByID(id string) (*User, error)

// Unexported: camelCase
func createSessionToken() (string, error)
func hashPassword(password string) (string, error)

// REST handlers: Uppercase HTTP method names
func (h *UserHandler) GET(c *fiber.Ctx) error
func (h *UserHandler) POST(c *fiber.Ctx) error
func (h *UserHandler) GETByID(c *fiber.Ctx) error
```

#### Variables
```go
// Descriptive names for important variables
userRepository := repository.NewUserRepository(db)
sessionToken := generateToken()

// Short receiver names (consistent per type)
func (u *User) IsAuthorized() bool
func (s *Session) IsValid() bool
```

### Constructor Pattern

**Always use `New*` constructors:**

```go
// Handler
func NewUserHandler(
    service *service.UserService,
    validator *validators.Validator,
) *UserHandler {
    return &UserHandler{
        service:   service,
        validator: validator,
    }
}

// Service
func NewUserService(
    userRepo *repository.UserRepository,
    roleRepo *repository.RoleRepository,
) *UserService {
    return &UserService{
        userRepo: userRepo,
        roleRepo: roleRepo,
    }
}

// Repository
func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}
```

### Error Handling

```go
// Services and Repositories: Return errors
if err != nil {
    return nil, err
}

// Check for specific GORM errors
if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
}

// Handlers: Return errors (global handler processes them)
user, err := h.service.GetUser(id)
if err != nil {
    return err  // ErrorHandler converts to JSON
}
```

### Logging

```go
import applogger "VEDA95/open_board/api/internal/log"

// Error logging
applogger.Global.Error().Err(err).Msg("failed to create user")

// Info logging
applogger.Global.Info().Msgf("user %s logged in", username)

// Fatal logging (exits program)
applogger.Global.Fatal().Err(err).Msg("database connection failed")

// Debug logging (only in development)
applogger.Global.Debug().Interface("user", user).Msg("user data")
```

### Import Organization

```go
import (
    // Standard library
    "fmt"
    "time"

    // Internal packages (alphabetically)
    "VEDA95/open_board/api/internal/config"
    "VEDA95/open_board/api/internal/db/model"
    "VEDA95/open_board/api/internal/service"

    // External packages (alphabetically)
    "github.com/gofiber/fiber/v2"
    "github.com/gofrs/uuid/v5"
    "gorm.io/gorm"
)
```

---

## Development Workflows

### Initial Setup

1. **Clone and setup:**
```bash
git clone <repository-url>
cd OpenBoard-Backend
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Setup PostgreSQL:**
```bash
cd deployment/dev
docker-compose -f docker-compose.postgres.yaml up -d
```

4. **Configure environment:**
```bash
# Edit env/.env.development with your settings
# Key variables:
# - DATABASE_URL
# - PORT, HOST
# - AUTH_SESSION_EXPIRES_IN
```

5. **Run migrations:**
```bash
# Install Atlas
curl -sSf https://atlasgo.sh | sh

# Apply migrations
atlas migrate apply --env gorm
```

6. **Seed database (optional):**
```bash
make run_dev_seed_program
```

### Development Server

**Using Air (hot reload):**
```bash
# Install Air if not present
go install github.com/air-verse/air@latest

# Run with hot reload
make run_dev_api
```

**Or build and run manually:**
```bash
make build_api
./build/backend_api
```

### Building

```bash
# Build all components
make install

# Build specific components
make build_api
make build_seed_program
make build_swagger_documentation
```

### Database Migrations

**Generate migration from model changes:**
```bash
# 1. Modify models in internal/db/model/
# 2. Generate migration
atlas migrate diff <migration_name> --env gorm

# 3. Apply migration
atlas migrate apply --env gorm
```

**Important:** This project uses Atlas for schema-as-code migrations. Do not manually create SQL migrations unless absolutely necessary.

### Swagger Documentation

**Generate API docs:**
```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
make build_swagger_documentation

# Access docs at: http://localhost:8080/swagger/
```

### Admin CLI Tools

**Create user:**
```bash
make run_dev_create_user_program
# Follow prompts
```

**Reset password:**
```bash
make run_dev_reset_password_program
# Follow prompts
```

---

## Key Components

### 1. HTTP Handlers (Routes)

**Location:** `/internal/http/routes/`

**Responsibilities:**
- Parse HTTP requests (params, body, files)
- Validate input using validators
- Call service layer methods
- Format and return responses
- Extract authentication context

**Standard Structure:**
```go
type UserHandler struct {
    service   *service.UserService
    validator *validators.Validator
}

func NewUserHandler(
    service *service.UserService,
    validator *validators.Validator,
) *UserHandler {
    return &UserHandler{
        service:   service,
        validator: validator,
    }
}

// REST methods
func (h *UserHandler) GET(c *fiber.Ctx) error           // List all
func (h *UserHandler) GETByID(c *fiber.Ctx) error       // Get one
func (h *UserHandler) POST(c *fiber.Ctx) error          // Create
func (h *UserHandler) PATCH(c *fiber.Ctx) error         // Update
func (h *UserHandler) DELETE(c *fiber.Ctx) error        // Delete
```

**Example Handler:**
```go
func (h *UserHandler) POST(c *fiber.Ctx) error {
    // 1. Parse and validate request
    validator := new(validators.CreateUserValidator)
    if err := c.BodyParser(validator); err != nil {
        return err
    }

    if errs := h.validator.Validate(validator); len(errs) > 0 {
        return errors.CreateValidationError(errs)
    }

    // 2. Call service
    user, err := h.service.CreateUser(validator)
    if err != nil {
        return err
    }

    // 3. Return response
    return responses.JSONResponse(
        c,
        fiber.StatusCreated,
        responses.OKResponse(fiber.StatusCreated, user),
    )
}
```

### 2. Services (Business Logic)

**Location:** `/internal/service/`

**Responsibilities:**
- Implement business rules
- Orchestrate repository calls
- Handle authorization logic
- Manage transactions
- Validate business constraints

**Standard Structure:**
```go
type UserService struct {
    userRepo *repository.UserRepository
    roleRepo *repository.RoleRepository
}

func NewUserService(
    userRepo *repository.UserRepository,
    roleRepo *repository.RoleRepository,
) *UserService {
    return &UserService{
        userRepo: userRepo,
        roleRepo: roleRepo,
    }
}
```

**Example Service Method:**
```go
func (s *UserService) CreateUser(
    data *validators.CreateUserValidator,
) (*models.User, error) {
    // 1. Business validation
    if s.userRepo.ExistsByUsernameOrEmail(data.Username, data.Email) {
        return nil, errors.New("user already exists")
    }

    // 2. Create entity
    user := &models.User{
        Username: data.Username,
        Email:    data.Email,
    }

    if err := user.HashPassword(data.Password); err != nil {
        return nil, err
    }

    // 3. Transaction with multiple operations
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    if data.Roles != nil {
        roles, err := s.roleRepo.FindByIDs(*data.Roles)
        if err != nil {
            return nil, err
        }
        user.Roles = roles
    }

    return user, nil
}
```

### 3. Repositories (Data Access)

**Location:** `/internal/db/repository/`

**Responsibilities:**
- CRUD operations
- Query building with GORM
- Transaction management
- Existence checks

**Standard Structure:**
```go
type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}
```

**Standard Methods:**
```go
func (r *UserRepository) FindAll(
    options repository.QueryOptions,
) ([]*models.User, error)

func (r *UserRepository) FindByID(
    id string,
    options repository.QueryOptions,
) (*models.User, error)

func (r *UserRepository) Create(
    user *models.User,
    options repository.QueryOptions,
) error

func (r *UserRepository) Update(
    user *models.User,
    options repository.QueryOptions,
) error

func (r *UserRepository) Delete(id string) error

func (r *UserRepository) Exists(id string) bool
```

**QueryOptions Pattern:**
```go
type QueryOptions struct {
    Select  []string  // Columns to select
    Omit    []string  // Columns to omit
    Preload []string  // Relations to eager load
}

// Usage
users, err := repo.FindAll(repository.QueryOptions{
    Select:  []string{"id", "username", "email"},
    Preload: []string{"Roles", "Roles.Permissions"},
})
```

### 4. Models (Domain Entities)

**Location:** `/internal/db/model/`

**Base Model Pattern:**
```go
type BaseID struct {
    ID string `gorm:"type:uuid;primaryKey" json:"id"`
}

type Base struct {
    BaseID
    CreatedAt time.Time  `gorm:"not null" json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`
}

// UUID auto-generation hook
func (b *BaseID) BeforeCreate(tx *gorm.DB) error {
    if b.ID == "" {
        id, err := uuid.NewV7()
        if err != nil {
            return err
        }
        b.ID = id.String()
    }
    return nil
}
```

**Example Model:**
```go
type User struct {
    Base
    Username string  `gorm:"unique;not null" json:"username"`
    Email    string  `gorm:"unique;not null" json:"email"`
    Password string  `gorm:"not null" json:"-"` // Never serialize
    Roles    []*Role `gorm:"many2many:user_roles" json:"roles,omitempty"`
}

// Business logic methods
func (u *User) HashPassword(password string) error {
    hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
    if err != nil {
        return err
    }
    u.Password = hash
    return nil
}

func (u *User) VerifyPassword(password string) (bool, error) {
    return argon2id.ComparePasswordAndHash(password, u.Password)
}
```

### 5. Validators

**Location:** `/internal/http/validators/`

**Validation Struct:**
```go
type CreateUserValidator struct {
    Username string    `json:"username" validate:"required,min=1"`
    Email    string    `json:"email" validate:"required,email"`
    Password string    `json:"password" validate:"required,min=8,max=32"`
    Roles    *[]string `json:"roles,omitempty" validate:"omitempty,min=1"`
}
```

**Common Validation Tags:**
- `required` - Field must be present
- `min`, `max` - Length constraints
- `email` - Email format validation
- `uuid` - UUID format validation
- `oneof` - Enum validation
- `omitempty` - Optional field

**Usage in Handler:**
```go
validator := new(validators.CreateUserValidator)
if err := c.BodyParser(validator); err != nil {
    return err
}

if errs := h.validator.Validate(validator); len(errs) > 0 {
    return errors.CreateValidationError(errs)
}
```

### 6. Middleware

**Location:** `/internal/http/middleware/`

**Authentication Middleware:**
```go
type AuthMiddleware struct {
    authService *service.AuthService
}

func (m *AuthMiddleware) RequireAuthentication() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Extract token from header or cookie
        token := extractToken(c)

        // Validate session
        session, err := m.authService.ValidateSession(token)
        if err != nil {
            return fiber.ErrUnauthorized
        }

        // Store in context
        c.Locals("auth_session", session)
        return c.Next()
    }
}

func (m *AuthMiddleware) RequireAuthorization(
    permissions ...string,
) fiber.Handler {
    return func(c *fiber.Ctx) error {
        session := c.Locals("auth_session").(*models.Session)

        if !session.User.IsAuthorized(permissions...) {
            return fiber.ErrForbidden
        }

        return c.Next()
    }
}
```

**Usage:**
```go
// In main.go
app.Get("/api/users",
    authMiddleware.RequireAuthentication(),
    userHandler.GET,
)

app.Delete("/api/users/:id",
    authMiddleware.RequireAuthentication(),
    authMiddleware.RequireAuthorization("users:delete"),
    userHandler.DELETE,
)
```

---

## Database & Migrations

### GORM Configuration

**Connection Setup** (`internal/db/db.go`):
```go
func NewDB() (*gorm.DB, error) {
    dsn := os.Getenv("DATABASE_URL")

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: customLogger,
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
        SkipDefaultTransaction: true,  // Performance optimization
    })

    if err != nil {
        return nil, err
    }

    // Connection pooling
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return db, nil
}
```

### Atlas Migrations

**Philosophy:** Schema-as-code. Models are the source of truth.

**Workflow:**

1. **Modify model** in `/internal/db/model/`:
```go
type User struct {
    Base
    Username  string  `gorm:"unique;not null"`
    Email     string  `gorm:"unique;not null"`
    // Add new field:
    FirstName string  `gorm:"type:varchar(100)"`
}
```

2. **Generate migration**:
```bash
atlas migrate diff add_user_first_name --env gorm
```

3. **Review generated SQL** in `/migrations/`:
```sql
-- migrations/20250825003923_add_user_first_name.sql
ALTER TABLE "users" ADD COLUMN "first_name" varchar(100) NULL;
```

4. **Apply migration**:
```bash
atlas migrate apply --env gorm
```

**Atlas Config** (`atlas.hcl`):
```hcl
data "external_schema" "gorm" {
  program = ["go", "run", "-mod=mod", "./cmd/loader/main.go"]
}

env "gorm" {
  src = data.external_schema.gorm.url
  url = getenv("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
}
```

### Transactions

**GORM Transactions:**
```go
// Automatic rollback on error
err := repo.db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(user).Error; err != nil {
        return err  // Rolls back
    }

    if err := tx.Model(user).Association("Roles").Append(roles); err != nil {
        return err  // Rolls back
    }

    return nil  // Commits
})
```

**Manual Transaction:**
```go
tx := repo.db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(session).Error; err != nil {
    tx.Rollback()
    return err
}

return tx.Commit().Error
```

---

## Authentication & Authorization

### Session-Based Authentication

**Session Model:**
```go
type Session struct {
    Base
    ExpiresOn        time.Time
    RefreshExpiresOn *time.Time
    AccessToken      string  `gorm:"unique;not null"`
    RefreshToken     *string `gorm:"unique"`
    RememberMe       bool
    UserID           string
    User             *User  `gorm:"foreignKey:UserID"`
}

func (s *Session) IsValid() bool {
    return time.Now().UTC().Before(s.ExpiresOn)
}

func (s *Session) IsRefreshValid() bool {
    if s.RefreshExpiresOn == nil {
        return false
    }
    return time.Now().UTC().Before(*s.RefreshExpiresOn)
}
```

**Login Flow:**
```go
// 1. Validate credentials
user, err := authService.userRepo.FindByUsername(username)
if err != nil {
    return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
}

valid, err := user.VerifyPassword(password)
if !valid {
    return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
}

// 2. Create session
session := &models.Session{
    UserID:     user.ID,
    RememberMe: rememberMe,
}

accessToken, _ := generateToken(32)
session.AccessToken = accessToken
session.ExpiresOn = time.Now().UTC().Add(sessionExpiry)

if rememberMe {
    refreshToken, _ := generateToken(32)
    session.RefreshToken = &refreshToken
    refreshExpiry := time.Now().UTC().Add(refreshSessionExpiry)
    session.RefreshExpiresOn = &refreshExpiry
}

// 3. Save session
if err := authService.sessionRepo.Create(session); err != nil {
    return nil, err
}

return session, nil
```

**Authentication Middleware:**
```go
func (m *AuthMiddleware) RequireAuthentication() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Extract token from header or cookie
        authHeader := c.Get("Authorization")
        token := strings.TrimPrefix(authHeader, "Bearer ")

        if token == "" {
            token = c.Cookies("access_token")
        }

        if token == "" {
            return fiber.ErrUnauthorized
        }

        // Validate session
        session, err := m.authService.ValidateSession(token)
        if err != nil {
            c.ClearCookie("access_token", "refresh_token")
            return fiber.ErrUnauthorized
        }

        // Store in context
        c.Locals("auth_session", session)
        return c.Next()
    }
}
```

### Role-Based Access Control (RBAC)

**RBAC Models:**
```go
type User struct {
    Base
    Roles []*Role `gorm:"many2many:user_roles"`
}

type Role struct {
    Base
    Name        string       `gorm:"unique;not null"`
    Permissions []Permission `gorm:"many2many:role_permissions"`
}

type Permission struct {
    Base
    Path string `gorm:"unique;not null"`  // e.g., "auth:superuser", "boards:manage"
}
```

**Permission Checking:**
```go
// On User model
func (u *User) IsAuthorized(paths ...string) bool {
    // Superuser bypass
    if u.IsSuperuser() {
        return true
    }

    // Check all required permissions
    for _, path := range paths {
        found := false
        for _, role := range u.Roles {
            for _, perm := range role.Permissions {
                if perm.Path == path {
                    found = true
                    break
                }
            }
        }
        if !found {
            return false
        }
    }
    return true
}

func (u *User) IsAuthorizedPartial(paths ...string) bool {
    // Requires ANY permission (not all)
    if u.IsSuperuser() {
        return true
    }

    for _, path := range paths {
        for _, role := range u.Roles {
            for _, perm := range role.Permissions {
                if perm.Path == path {
                    return true
                }
            }
        }
    }
    return false
}

func (u *User) IsSuperuser() bool {
    return u.IsAuthorized("auth:superuser")
}
```

**Authorization Middleware:**
```go
authMiddleware.RequireAuthorization("boards:delete", "boards:manage_all")
```

**Service-Level Authorization:**
```go
func (s *BoardService) DeleteBoard(
    boardID string,
    user *models.User,
) error {
    board, err := s.boardRepo.FindByID(boardID)
    if err != nil {
        return err
    }

    // Check ownership or permission
    if board.OwnerID != user.ID &&
       !user.IsAuthorized("boards:manage_all") {
        return fiber.ErrForbidden
    }

    return s.boardRepo.Delete(boardID)
}
```

### Password Management

**Hashing (Argon2id):**
```go
func (u *User) HashPassword(password string) error {
    hash, err := argon2id.CreateHash(
        password,
        argon2id.DefaultParams,
    )
    if err != nil {
        return err
    }
    u.Password = hash
    return nil
}

func (u *User) VerifyPassword(password string) (bool, error) {
    return argon2id.ComparePasswordAndHash(password, u.Password)
}
```

**Password Reset:**
```go
type PasswordReset struct {
    Base
    Token     string     `gorm:"unique;not null;size:6"`
    ExpiresOn time.Time  `gorm:"not null"`
    UserID    string
    User      *User      `gorm:"foreignKey:UserID"`
}

func (p *PasswordReset) IsValid() bool {
    return time.Now().UTC().Before(p.ExpiresOn)
}

// Two flows:
// 1. Authenticated reset (in-app)
// 2. Unauthenticated reset (forgot password with email)
```

---

## API Design Patterns

### REST Conventions

**Endpoint Structure:**
```
/auth/*        - Authentication endpoints (mixed auth)
/api/*         - Main API endpoints (requires auth)
```

**HTTP Methods:**
- `GET` - Retrieve resources (collection or single)
- `POST` - Create new resource
- `PATCH` - Partial update
- `DELETE` - Remove resource

**Route Patterns:**
```go
// Collection endpoints
app.Get("/api/users", handler.GET)
app.Post("/api/users", handler.POST)

// Resource endpoints
app.Get("/api/users/:id", handler.GETByID)
app.Patch("/api/users/:id", handler.PATCH)
app.Delete("/api/users/:id", handler.DELETE)

// Nested resources
app.Get("/api/boards/:id/lists", listHandler.GETByBoardID)
app.Post("/api/cards/:id/comments", commentHandler.POST)

// Special endpoints
app.Get("/auth/@me", authHandler.UserInfoGET)
app.Post("/auth/login", authHandler.LocalLogin)
```

### Request Patterns

**Path Parameters:**
```go
params := new(validators.ParamValidator)
if err := c.ParamsParser(params); err != nil {
    return err
}
// Access: params.Id
```

**Query Parameters:**
```go
query := new(validators.QueryValidator)
if err := c.QueryParser(query); err != nil {
    return err
}
// Access: query.Page, query.Limit
```

**JSON Body:**
```go
data := new(validators.CreateUserValidator)
if err := c.BodyParser(data); err != nil {
    return err
}
// Access: data.Username, data.Email
```

**File Uploads:**
```go
file, err := c.FormFile("file_upload")
if err != nil {
    return err
}

// Save file
if err := c.SaveFile(file, dest); err != nil {
    return err
}
```

### Response Patterns

**Success Response (Single Resource):**
```go
return responses.JSONResponse(
    c,
    fiber.StatusOK,
    responses.OKResponse(fiber.StatusOK, user),
)
```
```json
{
  "code": 200,
  "data": {
    "id": "uuid",
    "username": "john"
  }
}
```

**Success Response (Collection):**
```go
return responses.JSONResponse(
    c,
    fiber.StatusOK,
    responses.OKCollectionResponse(fiber.StatusOK, users),
)
```
```json
{
  "code": 200,
  "count": 5,
  "data": [...]
}
```

**Success with Message:**
```go
return responses.JSONResponse(
    c,
    fiber.StatusCreated,
    responses.OKResponse(fiber.StatusCreated, fiber.Map{
        "message": "User created successfully",
        "user": user,
    }),
)
```
```json
{
  "code": 201,
  "data": {
    "message": "User created successfully",
    "user": {...}
  }
}
```

**Error Response:**
```json
{
  "code": 422,
  "errors": {
    "email": {
      "failed_field": "email",
      "tag": "email",
      "value": "invalid",
      "err_value": "A valid email address must be provided"
    }
  }
}
```

### Status Codes

| Code | Usage |
|------|-------|
| 200 | Successful GET, PATCH, DELETE |
| 201 | Successful POST (creation) |
| 401 | Authentication required/failed |
| 403 | Forbidden (authorization failed) |
| 404 | Resource not found |
| 422 | Validation errors |
| 500 | Internal server error |

### Swagger Documentation

**Annotate handlers:**
```go
// UsersGET godoc
//
//  @Description   List all users
//  @Summary       List users
//  @Tags          users
//  @Success       200 {object} responses.OkCollectionResponse[models.User]
//  @Failure       401 {object} responses.ErrorResponse[responses.GenericMessage]
//  @Failure       500 {object} responses.ErrorResponse[responses.GenericMessage]
//  @Router        /api/users [get]
//  @Produce       json
//  @Security      BearerAuth
func (h *UserHandler) GET(c *fiber.Ctx) error {
    // Implementation
}
```

**Generate docs:**
```bash
make build_swagger_documentation
```

**Access:** `http://localhost:8080/swagger/`

---

## Error Handling

### Global Error Handler

**Location:** `/internal/errors/error.go`

```go
func ErrorHandler(ctx *fiber.Ctx, err error) error {
    code := fiber.StatusInternalServerError

    // Handle Fiber errors
    if e, ok := err.(*fiber.Error); ok {
        code = e.Code
        return ctx.Status(code).JSON(responses.ErrorResp(code, e.Message))
    }

    // Handle validation errors
    if e, ok := err.(*ValidationError); ok {
        return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
            responses.ErrorResp(fiber.StatusUnprocessableEntity, e.Errors),
        )
    }

    // Log 500 errors
    if code == fiber.StatusInternalServerError {
        applogger.Global.Error().Err(err).Msg("")
    }

    return ctx.Status(code).JSON(
        responses.ErrorResp(code, err.Error()),
    )
}
```

### Creating Errors

**Simple errors:**
```go
return nil, errors.New("user not found")
```

**Fiber errors (with status code):**
```go
return fiber.NewError(fiber.StatusNotFound, "user not found")
return fiber.NewError(fiber.StatusForbidden, "insufficient permissions")
return fiber.ErrUnauthorized
```

**Validation errors:**
```go
if errs := validator.Validate(data); len(errs) > 0 {
    return errors.CreateValidationError(errs)
}
```

### Error Checking

**GORM errors:**
```go
if errors.Is(err, gorm.ErrRecordNotFound) {
    return fiber.NewError(fiber.StatusNotFound, "resource not found")
}
```

**Custom error types:**
```go
type NotFoundError struct {
    Resource string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s not found", e.Resource)
}
```

---

## Testing Guidelines

### Test Structure

**Test files:** `*_test.go` in same package

```go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)

    // Act
    user, err := service.CreateUser(validData)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "john", user.Username)
}
```

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/service

# With coverage
go test -cover ./...

# Verbose
go test -v ./...
```

### Mocking

**Create mock repositories:**
```go
type MockUserRepository struct {
    FindByIDFunc func(id string) (*models.User, error)
    CreateFunc   func(user *models.User) error
}

func (m *MockUserRepository) FindByID(id string) (*models.User, error) {
    if m.FindByIDFunc != nil {
        return m.FindByIDFunc(id)
    }
    return nil, nil
}
```

---

## Common Tasks

### Adding a New Endpoint

1. **Create validator** (if needed):
```go
// internal/http/validators/board.go
type CreateBoardValidator struct {
    Name        string `json:"name" validate:"required,min=1,max=100"`
    Description string `json:"description,omitempty" validate:"max=500"`
}
```

2. **Add repository method** (if needed):
```go
// internal/db/repository/board.go
func (r *BoardRepository) Create(
    board *models.Board,
    options QueryOptions,
) error {
    return r.db.Create(board).Error
}
```

3. **Add service method**:
```go
// internal/service/board.go
func (s *BoardService) CreateBoard(
    data *validators.CreateBoardValidator,
    userID string,
) (*models.Board, error) {
    board := &models.Board{
        Name:        data.Name,
        Description: data.Description,
        OwnerID:     userID,
    }

    if err := s.boardRepo.Create(board); err != nil {
        return nil, err
    }

    return board, nil
}
```

4. **Add handler method**:
```go
// internal/http/routes/board.go
func (h *BoardHandler) POST(c *fiber.Ctx) error {
    validator := new(validators.CreateBoardValidator)
    if err := c.BodyParser(validator); err != nil {
        return err
    }

    if errs := h.validator.Validate(validator); len(errs) > 0 {
        return errors.CreateValidationError(errs)
    }

    session := c.Locals("auth_session").(*models.Session)

    board, err := h.service.CreateBoard(validator, session.UserID)
    if err != nil {
        return err
    }

    return responses.JSONResponse(
        c,
        fiber.StatusCreated,
        responses.OKResponse(fiber.StatusCreated, board),
    )
}
```

5. **Register route**:
```go
// cmd/server/main.go
apiGroup.Post("/boards",
    authMiddleware.RequireAuthentication(),
    boardHandler.POST,
)
```

6. **Generate Swagger docs**:
```bash
make build_swagger_documentation
```

### Adding a New Model

1. **Create model**:
```go
// internal/db/model/tag.go
type Tag struct {
    Base
    Name    string  `gorm:"unique;not null" json:"name"`
    Color   string  `gorm:"type:varchar(7)" json:"color"`
    BoardID string  `json:"board_id"`
    Board   *Board  `gorm:"foreignKey:BoardID" json:"board,omitempty"`
}
```

2. **Register in loader**:
```go
// cmd/loader/main.go
stmts, err := gormschema.New("postgres").Load(
    &models.User{},
    &models.Tag{},  // Add here
    // ... other models
)
```

3. **Generate migration**:
```bash
atlas migrate diff add_tags_table --env gorm
```

4. **Review and apply**:
```bash
atlas migrate apply --env gorm
```

### Adding Middleware

1. **Create middleware**:
```go
// internal/http/middleware/rate_limit.go
func RateLimit() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Rate limiting logic
        return c.Next()
    }
}
```

2. **Apply globally**:
```go
// cmd/server/main.go
app.Use(middleware.RateLimit())
```

3. **Or apply to specific routes**:
```go
apiGroup.Get("/sensitive",
    middleware.RateLimit(),
    handler.GET,
)
```

### Environment Configuration

**Add new environment variable:**

1. **Update env files**:
```bash
# env/.env.development
NEW_FEATURE_ENABLED=true
```

2. **Access in code**:
```go
enabled := os.Getenv("NEW_FEATURE_ENABLED") == "true"
```

3. **Document in README or this file**

---

## AI Assistant Best Practices

### When Making Changes

1. **Read before modifying**: Always read existing files before suggesting changes
2. **Follow patterns**: Use existing patterns in the codebase
3. **Maintain layer separation**: Don't skip layers (handler → service → repository)
4. **Validate at boundaries**: Validators in handlers, business logic in services
5. **Use constructors**: Always create `New*` functions
6. **Error handling**: Return errors up the stack, don't swallow them
7. **Log appropriately**: Error/Fatal for errors, Info for important events, Debug for details

### Code Quality Checklist

- [ ] Does it follow the layered architecture?
- [ ] Are there appropriate validators for input?
- [ ] Are errors handled and logged?
- [ ] Is authentication/authorization checked?
- [ ] Are database operations in repositories only?
- [ ] Are business rules in services only?
- [ ] Is the code documented with comments?
- [ ] Are Swagger annotations added?
- [ ] Does it follow existing naming conventions?

### Common Mistakes to Avoid

1. **Don't put DB operations in handlers**
   ```go
   // ❌ Bad
   func (h *UserHandler) GET(c *fiber.Ctx) error {
       var users []models.User
       h.db.Find(&users)  // NO!
   }

   // ✅ Good
   func (h *UserHandler) GET(c *fiber.Ctx) error {
       users, err := h.service.GetUsers()  // YES!
   }
   ```

2. **Don't skip validation**
   ```go
   // ❌ Bad
   func (h *UserHandler) POST(c *fiber.Ctx) error {
       var data map[string]interface{}
       c.BodyParser(&data)  // NO!
   }

   // ✅ Good
   func (h *UserHandler) POST(c *fiber.Ctx) error {
       validator := new(validators.CreateUserValidator)
       c.BodyParser(validator)
       if errs := h.validator.Validate(validator); len(errs) > 0 {
           return errors.CreateValidationError(errs)
       }
   }
   ```

3. **Don't put business logic in repositories**
   ```go
   // ❌ Bad
   func (r *UserRepository) CreateUser(...) error {
       if r.ExistsByEmail(email) {  // Business logic!
           return errors.New("user exists")
       }
   }

   // ✅ Good
   func (s *UserService) CreateUser(...) error {
       if s.userRepo.ExistsByEmail(email) {  // YES!
           return errors.New("user exists")
       }
   }
   ```

4. **Don't hardcode values**
   ```go
   // ❌ Bad
   session.ExpiresOn = time.Now().Add(24 * time.Hour)

   // ✅ Good
   expiry, _ := time.ParseDuration(os.Getenv("AUTH_SESSION_EXPIRES_IN"))
   session.ExpiresOn = time.Now().Add(expiry)
   ```

---

## Reference Files

### Essential Files to Reference

When working on specific tasks, refer to these files:

| Task | Reference Files |
|------|----------------|
| Adding endpoint | `internal/http/routes/user.go` |
| Adding service | `internal/service/user.go` |
| Adding repository | `internal/db/repository/user.go` |
| Adding model | `internal/db/model/user.go` |
| Adding validator | `internal/http/validators/user.go` |
| Authentication | `internal/service/auth.go`, `internal/http/middleware/auth.go` |
| Error handling | `internal/errors/error.go` |
| Response formatting | `internal/http/responses/` |
| Dependency injection | `cmd/server/main.go` |
| Database config | `internal/db/db.go` |
| Logging | `internal/log/log.go` |

---

## Conclusion

This guide provides a comprehensive overview of the OpenBoard Backend codebase. When in doubt:

1. **Follow existing patterns** - The codebase is consistent
2. **Check reference files** - See how similar features are implemented
3. **Maintain separation of concerns** - Keep layers distinct
4. **Ask for clarification** - If requirements are unclear

For questions or updates to this guide, consult the development team or submit a pull request.

**Last Updated:** 2025-12-01
**Version:** 1.0.0
