# 📝 Task Management API

A RESTful API for managing tasks and categories built with Go, Gin, GORM, PostgreSQL, and JWT authentication.

## 🚀 Features

- ✅ User Authentication (Register/Login) with JWT
- ✅ CRUD operations for Tasks and Categories
- ✅ Advanced filtering, sorting, and pagination
- ✅ Category-based task organization
- ✅ Task priority and status management
- ✅ Swagger API documentation
- ✅ Soft delete support
- ✅ RESTful API design

## 🛠️ Tech Stack

- **Language**: Go 1.25.4
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt/jwt)
- **Documentation**: Swagger (swaggo/swag)
- **Password Hashing**: bcrypt

## 📋 Prerequisites

Before you begin, ensure you have the following installed:
- Go 1.21 or higher
- PostgreSQL 14 or higher
- Git

## 🔧 Installation & Setup

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/TaskManagementAPI.git
cd TaskManagementAPI
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Setup PostgreSQL Database

Create a new PostgreSQL database:

```sql
CREATE DATABASE taskmanagement;
```

Or using command line:

```bash
psql -U postgres -c "CREATE DATABASE taskmanagement;"
```

### 4. Configure Environment Variables

Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

Edit `.env` file with your configuration:

```env
SERVER_PORT=8080
GIN_MODE=debug

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=taskmanagement

JWT_SECRET=your-secret-key-at-least-32-characters-long
JWT_EXPIRY_HOURS=24
```

### 5. Run the Application

```bash
go run cmd/api/main.go
```

The server will start at `http://localhost:8080`

### 6. Generate Swagger Documentation (Optional)

Install swag CLI:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Generate docs:

```bash
swag init -g cmd/api/main.go -o docs
```

Access Swagger UI at: `http://localhost:8080/swagger/index.html`

## 📚 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Register new user | No |
| POST | `/api/v1/auth/login` | Login user | No |
| GET | `/api/v1/auth/me` | Get current user | Yes |

### Categories

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/categories` | Get all categories | Yes |
| POST | `/api/v1/categories` | Create category | Yes |
| GET | `/api/v1/categories/:id` | Get category by ID | Yes |
| PUT | `/api/v1/categories/:id` | Update category | Yes |
| DELETE | `/api/v1/categories/:id` | Delete category | Yes |

### Tasks

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/tasks` | Get all tasks (with filters) | Yes |
| POST | `/api/v1/tasks` | Create task | Yes |
| GET | `/api/v1/tasks/:id` | Get task by ID | Yes |
| PUT | `/api/v1/tasks/:id` | Update task | Yes |
| PATCH | `/api/v1/tasks/:id/status` | Update task status | Yes |
| DELETE | `/api/v1/tasks/:id` | Delete task | Yes |

### Query Parameters for Tasks

- `status`: Filter by status (pending, in_progress, completed)
- `priority`: Filter by priority (low, medium, high)
- `category_id`: Filter by category
- `search`: Search in title and description
- `sort_by`: Sort by field (created_at, updated_at, due_date, priority)
- `sort_order`: Sort order (asc, desc)
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 10, max: 100)

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication. After login, include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## 📖 Example Usage

### Register a new user

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123",
    "full_name": "John Doe"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123"
  }'
```

### Create a task

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "title": "Complete project documentation",
    "description": "Write comprehensive README",
    "priority": "high",
    "status": "in_progress"
  }'
```

### Get tasks with filters

```bash
curl -X GET "http://localhost:8080/api/v1/tasks?status=pending&priority=high&sort_by=due_date&sort_order=asc&page=1&page_size=10" \
  -H "Authorization: Bearer <your_token>"
```

## 🗂️ Project Structure

```
TaskManagementAPI/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── models/
│   │   ├── user.go              # User model
│   │   ├── task.go              # Task model
│   │   └── category.go          # Category model
│   ├── database/
│   │   └── database.go          # Database connection
│   ├── middleware/
│   │   └── auth.go              # JWT authentication middleware
│   ├── handlers/
│   │   ├── auth_handler.go      # Auth endpoints
│   │   ├── task_handler.go      # Task endpoints
│   │   └── category_handler.go  # Category endpoints
│   ├── repository/
│   │   ├── user_repository.go   # User data access
│   │   ├── task_repository.go   # Task data access
│   │   └── category_repository.go # Category data access
│   └── services/
│       ├── auth_service.go      # Auth business logic
│       ├── task_service.go      # Task business logic
│       └── category_service.go  # Category business logic
├── docs/                        # Swagger documentation (auto-generated)
├── .env                         # Environment variables (not in git)
├── .env.example                 # Example environment variables
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## 🧪 Testing

Run tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## 🚢 Deployment

### Build for production

```bash
go build -o bin/api cmd/api/main.go
```

### Run production build

```bash
./bin/api
```

### Docker (Optional)

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api cmd/api/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./api"]
```

Build and run:

```bash
docker build -t taskmanagement-api .
docker run -p 8080:8080 taskmanagement-api
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License.

## 👨‍💻 Author

Your Name - [GitHub](https://github.com/yourusername)

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [golang-jwt](https://github.com/golang-jwt/jwt)
- [Swaggo](https://github.com/swaggo/swag)

