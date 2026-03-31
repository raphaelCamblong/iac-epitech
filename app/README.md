# Task Manager API

REST API for managing tasks with authentication.

## API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | /auth/register | No | Register user |
| POST | /auth/login | No | Login, get JWT |
| GET | /health | No | Health check |
| POST | /tasks | JWT | Create task |
| GET | /tasks | JWT | List tasks |
| GET | /tasks/{id} | JWT | Get task |
| PUT | /tasks/{id} | JWT | Update task |
| DELETE | /tasks/{id} | JWT | Delete task |


## Quick Start

### 1. Configure environment

```bash
cd app
export DATABASE_URL="postgres://user:pass@localhost:5432/taskmanager?sslmode=disable"
export JWT_SECRET="your-secret-min-32-chars"
```

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/taskmanager?sslmode=disable"
export JWT_SECRET="your-secret-min-32-chars"
go run ./cmd/server
```

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

## Build

```bash
cd app
go build -o server ./cmd/server
```

## Docker

```bash
cd app
docker build -t task-manager:latest .
docker run -e DATABASE_URL=... -e JWT_SECRET=... -p 8080:8080 task-manager:latest
```

# Structure

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


## Helm Deployment

```bash
helm install task-manager ./charts/task-manager \
  --set databaseUrl="postgres://..." \
  --set jwtSecret="your-secret" \
  --set image.repository=your-registry/task-manager \
  --set image.tag=latest
```

## Terraform Deployment

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
terraform init
terraform plan
terraform apply
```