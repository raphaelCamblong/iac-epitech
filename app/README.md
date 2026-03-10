# Task Manager API

REST API for managing tasks with authentication.

## Run tests

From the `app/` directory:

```bash
# All tests
go test ./...

# Verbose output
go test ./... -v

# Specific package
go test ./internal/usecase/...
go test ./internal/adapter/http/...

# With coverage
go test ./... -cover
```

## Run the app

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/taskmanager?sslmode=disable"
export JWT_SECRET="your-secret-min-32-chars"
go run ./cmd/server
```

## Structure

```
app/
├── cmd/server/           # Entry point + bootstrap
├── internal/
│   ├── entity/           # Domain models
│   ├── port/             # Repository interfaces
│   ├── usecase/
│   │   ├── task/         # Task logic + tests
│   │   └── auth/         # Auth logic + tests
│   └── adapter/
│       ├── http/         # HTTP handlers + middleware + tests
│       └── repository/   # GORM persistence
├── pkg/                  # Shared utilities
└── config/               # Configuration
```
