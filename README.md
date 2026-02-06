# PINTU SDN Sukapura 01

Portal Informasi Terpadu (PINTU) adalah sistem informasi terintegrasi untuk SDN Sukapura 01. Backend dibangun menggunakan Go dengan framework Gin dan database PostgreSQL.

## 📋 Daftar Isi

- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Setup & Installation](#setup--installation)
- [Running the Application](#running-the-application)
- [Database Setup](#database-setup)
- [Code Generation](#code-generation)
- [API Endpoints](#api-endpoints)
- [Troubleshooting](#troubleshooting)

## 🛠 Tech Stack

- **Language**: Go 1.25.6
- **Framework**: Gin Gonic
- **Database**: PostgreSQL 18
- **ORM**: GORM
- **Containerization**: Docker & Docker Compose
- **Architecture**: Clean Architecture

## 📁 Project Structure

```
pintu-backend/
├── cmd/                           # Command line tools & generators
│   ├── main.go                   # Main CLI entry point
│   └── generator.go              # File generators logic
├── pkg/                          # Packages (database connection, etc)
├── src/
│   ├── config/                   # Configuration files
│   ├── middleware/               # Middleware handlers
│   ├── database/
│   │   ├── migrations/           # SQL migration files
│   │   └── seeders/              # Data seeders
│   ├── modules/
│   │   ├── controllers/          # HTTP request handlers
│   │   ├── models/               # Database models
│   │   ├── repositories/         # Data access layer
│   │   └── services/             # Business logic layer
│   ├── dtos/                     # Data Transfer Objects
│   └── routes/                   # API routes definition
├── main.go                       # Application entry point
├── Dockerfile                    # Docker image configuration
├── docker-compose.yml            # Docker Compose configuration
├── Makefile                      # Build & run shortcuts
├── go.mod                        # Go module dependencies
├── .env                          # Environment variables (local)
├── .env.example                  # Environment variables template
├── .gitignore                    # Git ignore rules
└── README.md                     # This file
```

## 📋 Prerequisites

- Go 1.25.6 or higher
- PostgreSQL 18 or higher
- Docker & Docker Compose (for containerization)
- Git (optional)

## 🚀 Setup & Installation

### 1. Clone/Download Project

```bash
git clone https://github.com/SDN-Sukapura-01-Jakarta-Utara/pintu-backend.git
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Configure Environment Variables

Copy `.env.example` to `.env` and update the values:

```bash
# Windows
copy .env.example .env
```

Edit `.env`:

```
APP_NAME=PINTU SDN Sukapura 01
GIN_MODE=debug
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=sdn_sukapura_01
DB_SSLMODE=disable
```

## 🗄️ Database Setup

### Option 1: Using Command Prompt (Without pgAdmin)

#### Step 1: Create Database

```bash
"C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -c "CREATE DATABASE sdn_sukapura_01;"
```

#### Step 2: Verify Database Created

```bash
"C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -c "\l"
```

You should see `sdn_sukapura_01` in the list.

#### Step 3: Run Migrations (Manual)

Edit migration files in `src/database/migrations/` and run:

```bash
"C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d sdn_sukapura_01 -f src/database/migrations/20260206074547_create_users_table.sql
```

### Option 2: Using pgAdmin

1. Open pgAdmin
2. Connect to PostgreSQL server
3. Create new database: `sdn_sukapura_01`
4. Run migration files through pgAdmin interface

## 🐳 Running the Application

### Option 1: Run Locally

```bash
# Download dependencies
go mod tidy

# Run application
go run main.go
```

Application will run on `http://localhost:8080`

### Option 2: Run with Docker

```bash
# Build Docker image
docker build -t pintu-backend:latest .

# Run with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop container
docker-compose down
```

### Using Makefile

```bash
# Build locally
make build

# Run locally
make run

# Build Docker image
make docker-build

# Start Docker container
make docker-up

# Stop Docker container
make docker-down

# View Docker logs
make docker-logs

# Rebuild and restart Docker
make docker-rebuild
```

### Test Application

```bash
# Test root endpoint
curl http://localhost:8080/

# Test health check
curl http://localhost:8080/health
```

Expected response:

```json
{
  "app": "PINTU SDN Sukapura 01",
  "message": "PINTU Backend is running"
}
```

## 📝 Code Generation

Use built-in generators to quickly create boilerplate code.

### Generate Migration File

```bash
go run ./cmd generate:migration create_users_table
```

Creates: `src/database/migrations/[timestamp]_create_users_table.sql`

### Generate Model

```bash
go run ./cmd generate:model User
```

Creates: `src/modules/models/user.go`

**Update the model with your fields:**

```go
type User struct {
    ID        uint            `gorm:"primaryKey"`
    Name      string          `gorm:"not null"`
    Email     string          `gorm:"uniqueIndex"`
    Password  string          `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt  `gorm:"index"`
}
```

### Generate Repository

```bash
go run ./cmd generate:repository User
```

Creates: `src/modules/repositories/user_repository.go`

Provides interface and implementation with methods:
- `Create(data *models.User) error`
- `GetByID(id uint) (*models.User, error)`
- `GetAll() ([]models.User, error)`
- `Update(data *models.User) error`
- `Delete(id uint) error`

### Generate Service

```bash
go run ./cmd generate:service User
```

Creates: `src/modules/services/user_service.go`

Business logic layer that uses repository.

### Generate Controller

```bash
go run ./cmd generate:controller User
```

Creates: `src/modules/controllers/user_controller.go`

HTTP handlers with methods:
- `Create()` - POST
- `GetByID()` - GET by ID
- `GetAll()` - GET all
- `Update()` - PUT
- `Delete()` - DELETE

### Generate DTO (Data Transfer Object)

```bash
go run ./cmd generate:dto User
```

Creates: `src/dtos/user_dto.go`

Includes:
- `UserCreateRequest`
- `UserUpdateRequest`
- `UserResponse`
- `UserListResponse`

### Generate Seeder

```bash
go run ./cmd generate:seeder User
```

Creates: `src/database/seeders/user_seeder.go`

For populating initial data.

### Generate All at Once

```bash
go run ./cmd generate:model User && go run ./cmd generate:repository User && go run ./cmd generate:service User && go run ./cmd generate:controller User && go run ./cmd generate:dto User
```

## 📚 Complete Example: Creating User Module

### Step 1: Create Migration

```bash
go run ./cmd generate:migration create_users_table
```

Edit `src/database/migrations/[timestamp]_create_users_table.sql`:

```sql
-- Migration: create_users_table

BEGIN;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);

COMMIT;
```

Run migration:

```bash
"C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -d sdn_sukapura_01 -f src/database/migrations/[timestamp]_create_users_table.sql
```

### Step 2: Generate All Files

```bash
go run ./cmd generate:model User
go run ./cmd generate:repository User
go run ./cmd generate:service User
go run ./cmd generate:controller User
go run ./cmd generate:dto User
go run ./cmd generate:seeder User
```

### Step 3: Update Model

Edit `src/modules/models/user.go`:

```go
type User struct {
    ID        uint            `gorm:"primaryKey"`
    Name      string          `gorm:"not null"`
    Email     string          `gorm:"uniqueIndex"`
    Password  string          `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt  `gorm:"index"`
}
```

### Step 4: Update DTO

Edit `src/dtos/user_dto.go`:

```go
type UserCreateRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type UserUpdateRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}
```

### Step 5: Register Routes

Update `src/routes/routes.go`:

```go
package routes

import (
    "pintu-backend/src/modules/controllers"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
    // User routes
    // userController := controllers.NewUserController(service)
    // api := router.Group("/api/v1")
    // api.POST("/users", userController.Create)
    // api.GET("/users", userController.GetAll)
    // api.GET("/users/:id", userController.GetByID)
    // api.PUT("/users/:id", userController.Update)
    // api.DELETE("/users/:id", userController.Delete)
}
```

## 📡 API Endpoints

Once User module is fully set up:

```
POST   /api/v1/users              - Create new user
GET    /api/v1/users              - Get all users
GET    /api/v1/users/:id          - Get user by ID
PUT    /api/v1/users/:id          - Update user
DELETE /api/v1/users/:id          - Delete user
```

## 🔧 Troubleshooting

### Docker Build Error

```bash
# Clean Docker cache
docker builder prune -a

# Rebuild
docker-compose up -d --build
```

### Database Connection Failed

1. Check PostgreSQL is running
2. Verify credentials in `.env`
3. Check port 5432 is accessible

```bash
# Test connection
"C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -h localhost
```

### Port Already in Use

```bash
# Change port in .env
PORT=3000

# Or kill process using port 8080
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

### .env File Not Found in Docker

Ensure `.env` file is in the project root and `docker-compose.yml` has `env_file: .env`

## 📖 Best Practices

1. **Separation of Concerns**: Keep business logic in services, data access in repositories
2. **DTOs**: Always use DTOs for API requests/responses
3. **Error Handling**: Handle errors properly in all layers
4. **Migrations**: Version control all migrations
5. **Environment Variables**: Use `.env` for local development only
6. **Testing**: Write tests for services and repositories
7. **Logging**: Add structured logging for debugging

## 📜 License

Copyright 2026 SDN Sukapura 01. All rights reserved.

## 👥 Contributors

- Development Team

## 📞 Support

For issues, questions, or support, please contact:

- **WhatsApp**: 08889125991
- **Developer Email**: sdnsukapura01.dev@gmail.com
- **School Email**: sdnsukapuraa01@gmail.com
- **Personal Email**: syahiraisnaeni15@gmail.com
