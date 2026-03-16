# Task Manager API

Production-ready Task Manager REST API built with Go, Clean Architecture, PostgreSQL, and deployed via Helm/Terraform.

## Requirements

- Go 1.23+
- PostgreSQL (CloudSQL, RDS, or local)
- Docker (for building image)
- Kubernetes cluster with Helm
- Terraform (optional, for deployment)

## Quick Start

### 1. Configure environment

```bash
cd app
export DATABASE_URL="postgres://user:pass@localhost:5432/taskmanager?sslmode=disable"
export JWT_SECRET="your-secret-min-32-chars"
```

### 2. Run the application

```bash
cd app
go run ./cmd/server
```

### 3. Register and use the API

```bash
# Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Login (use token from register)
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Create task
curl -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -H "correlation_id: abc123" \
  -d '{"title":"Write","content":"Prepare lesson","due_date":"2025-09-30","request_timestamp":"2025-09-25T20:00:00Z"}'
```

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

## Project Structure

```
app/

charts/task-manager/      # Helm chart
terraform/                # Terraform deployment
```
