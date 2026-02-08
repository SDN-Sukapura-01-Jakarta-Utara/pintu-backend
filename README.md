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
- [Testing API dengan Postman](#testing-api-dengan-postman)
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

### Option 2: Using pgAdmin

1. Open pgAdmin
2. Connect to PostgreSQL server
3. Create new database: `sdn_sukapura_01`

## 🚀 Running Migrations & Seeders

### Run All Migrations

Migrations akan membuat semua tabel sesuai skema database:

```bash
go run ./cmd migrate:up
```

Migrations yang akan dijalankan (dalam urutan):
1. `20260206094322_create_roles_table.sql` - Tabel roles
2. `20260206094707_create_permissions_table.sql` - Tabel permissions
3. `20260206094811_create_users_table.sql` - Tabel users
4. `20260206094719_create_role_permissions_table.sql` - Pivot table role-permissions

### Run Specific Migration

Jalankan hanya migration tertentu:

```bash
# Contoh: jalankan migration tabel users
go run ./cmd migrate:file 20260206094811_create_users_table.sql

# Atau dengan nama file yang lebih pendek
go run ./cmd migrate:file create_users_table.sql
```

### Run All Seeders

Seeders akan mengisi data initial ke database:

```bash
go run ./cmd seed:run
```

Seeders yang akan dijalankan (dalam urutan):
1. **Permissions** - Insert 32 permissions untuk berbagai modul
2. **Roles** - Insert 2 roles: Administrator dan Kepala Sekolah
3. **Role Permissions** - Associate permissions ke roles
4. **Users** - Insert 2 default users:
   - Admin: username `admin`, password `admin123`
   - Kepala Sekolah: username `kepala_sekolah`, password `kepala123`

### Run Specific Seeder

Jalankan hanya seeder tertentu:

```bash
# Run permission seeder
go run ./cmd seed:specific permission

# Run role seeder
go run ./cmd seed:specific role

# Run role-permission seeder
go run ./cmd seed:specific role_permission

# Run user seeder
go run ./cmd seed:specific user
```

Available seeders: `permission`, `role`, `role_permission`, `user`

⚠️ **Important**: Change default passwords di production!

### Complete Setup Flow

```bash
# 1. Create database
"C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -c "CREATE DATABASE sdn_sukapura_01;"

# 2. Install dependencies
go mod tidy

# 3. Configure .env file
copy .env.example .env
# Edit .env with your database credentials

# 4. Run migrations
go run ./cmd migrate:up

# 5. Run seeders
go run ./cmd seed:run

# 6. Start application
go run main.go
```

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

### Authentication (Public)
```
POST   /api/v1/auth/login         - User login (returns JWT token)
```

### Protected Routes (Require Authentication)

**Permissions:**
```
POST   /api/v1/permissions/create-permission        - Create permission
POST   /api/v1/permissions/get-permissions          - Get all permissions (with pagination)
POST   /api/v1/permissions/get-permission           - Get permission by ID
POST   /api/v1/permissions/update-permission        - Update permission
POST   /api/v1/permissions/delete-permission        - Delete permission
POST   /api/v1/permissions/get-permissions-by-group - Get by group name
POST   /api/v1/permissions/get-permissions-by-system- Get by system
```

**Roles:**
```
POST   /api/v1/roles/create-role  - Create role
POST   /api/v1/roles/get-roles    - Get all roles
POST   /api/v1/roles/get-role     - Get role by ID
POST   /api/v1/roles/update-role  - Update role
POST   /api/v1/roles/delete-role  - Delete role
```

**Users:**
```
POST   /api/v1/users/create-user           - Create user
POST   /api/v1/users/get-users             - Get all users
POST   /api/v1/users/get-user              - Get user by ID
POST   /api/v1/users/update-user           - Update user
POST   /api/v1/users/update-user-password  - Update password
POST   /api/v1/users/delete-user           - Delete user
```

**Auth (Protected):**
```
POST   /api/v1/auth/logout        - Logout (requires token)
```

---

## 🧪 Testing API dengan Postman

### Step 1: Setup Postman Environment

1. **Buat Environment Baru**
   - Buka Postman → Environments → Create New
   - Nama: `PINTU Backend`
   - Tambah variables:

   | Variable | Initial Value | Current Value |
   |----------|---------------|---------------|
   | `base_url` | `http://localhost:3000` | `http://localhost:3000` |
   | `token` | (kosong) | (akan terisi setelah login) |
   | `username` | `admin` | `admin` |
   | `password` | `admin123` | `admin123` |

2. **Pilih Environment** dari dropdown (kanan atas Postman)

### Step 2: Testing Login

1. **Buat Request Baru:**
   ```
   Method: POST
   URL: {{base_url}}/api/v1/auth/login
   ```

2. **Headers:**
   ```
   Content-Type: application/json
   ```

3. **Body (raw JSON):**
   ```json
   {
     "username": "{{username}}",
     "password": "{{password}}"
   }
   ```

4. **Tambah Test Script** (tab Tests) untuk auto-save token:
   ```javascript
   if (pm.response.code === 200) {
       var jsonData = pm.response.json();
       var token = jsonData.data.token;
       pm.environment.set("token", token);
       console.log("✓ Token saved to environment!");
   }
   ```

5. **Send Request** → Status 200 OK

**Response Contoh:**
```json
{
  "status": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "nama": "Administrator",
      "username": "admin",
      "status": "active",
      "role_id": 1,
      "role_name": "Administrator",
      "accessible_system": ["dashboard", "users"],
      "created_at": "2026-02-06T10:30:45Z"
    },
    "expires_at": "2026-02-09T10:30:45Z"
  }
}
```

### Step 3: Testing Protected Routes (Dengan Authorization)

Setelah login, token tersimpan otomatis di environment `{{token}}`.

**Untuk setiap request ke protected route, tambahkan Authorization:**

**Option A: Pakai Tab Authorization (Recommended)**
- Tab: **Authorization**
- Type: **Bearer Token**
- Token: `{{token}}`

**Option B: Pakai Headers**
- Tab: **Headers**
- Key: `Authorization`
- Value: `Bearer {{token}}`

### Step 4: Contoh Testing Endpoints

#### **1. Get All Permissions**
```
POST {{base_url}}/api/v1/permissions/get-permissions
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "limit": 10,
  "offset": 0
}

Response (200):
{
  "data": [
    {
      "id": 1,
      "name": "CREATE_INFORMASI_SEKOLAH",
      "description": "Create school information",
      "group_name": "Informasi Sekolah",
      "system": "dashboard",
      "created_at": "2026-02-06T10:00:00Z",
      "updated_at": "2026-02-06T10:00:00Z"
    },
    ...
  ],
  "limit": 10,
  "offset": 0,
  "total": 32
}
```

#### **2. Get Permission by ID**
```
POST {{base_url}}/api/v1/permissions/get-permission
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "id": 1
}

Response (200):
{
  "data": {
    "id": 1,
    "name": "CREATE_INFORMASI_SEKOLAH",
    ...
  }
}
```

#### **3. Create Permission**
```
POST {{base_url}}/api/v1/permissions/create-permission
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "name": "VIEW_DASHBOARD",
  "description": "View dashboard",
  "group_name": "Dashboard",
  "system": "dashboard"
}

Response (201):
{
  "data": {
    "id": 33,
    "name": "VIEW_DASHBOARD",
    ...
  }
}
```

#### **4. Get All Users**
```
POST {{base_url}}/api/v1/users/get-users
Authorization: Bearer {{token}}
Content-Type: application/json

Response (200):
{
  "data": [
    {
      "id": 1,
      "nama": "Administrator",
      "username": "admin",
      "status": "active",
      "role_id": 1,
      "role_name": "Administrator",
      "accessible_system": ["dashboard", "users"],
      "created_at": "2026-02-06T10:00:00Z",
      "updated_at": "2026-02-06T10:00:00Z"
    },
    ...
  ]
}
```

#### **5. Create User**
```
POST {{base_url}}/api/v1/users/create-user
Authorization: Bearer {{token}}
Content-Type: application/json

Body:
{
  "nama": "Kepala Sekolah",
  "username": "kepala_sekolah",
  "password": "kepala123",
  "role_id": 2,
  "accessible_system": ["dashboard", "reports"],
  "status": "active"
}

Response (201):
{
  "data": {
    "id": 2,
    "nama": "Kepala Sekolah",
    ...
  }
}
```

### Step 5: Testing Error Scenarios

#### **Akses tanpa Token**
```
POST {{base_url}}/api/v1/permissions/get-permissions
(tanpa Authorization header)

Response (401):
{
  "error": "missing authorization header"
}
```

#### **Token Invalid/Expired**
```
POST {{base_url}}/api/v1/permissions/get-permissions
Authorization: Bearer invalid_token_here

Response (401):
{
  "error": "invalid or expired token"
}
```

#### **Login dengan Password Salah**
```
POST {{base_url}}/api/v1/auth/login

Body:
{
  "username": "admin",
  "password": "password_salah"
}

Response (401):
{
  "error": "username atau password salah"
}
```

#### **Logout Setelah Login**
```
POST {{base_url}}/api/v1/auth/logout
Authorization: Bearer {{token}}

Response (200):
{
  "status": "success",
  "message": "Logout successful, please delete your token",
  "user_id": 1
}
```

Setelah logout, client harus delete token dari environment dan login ulang untuk akses routes lain.

### Step 6: Tips & Tricks

- **Lihat Request History**: Klik History di sidebar kiri
- **Save Request ke Collection**: Klik Save setelah membuat request
- **Use Pre-request Scripts**: Set variables sebelum request dikirim
- **View Token Payload**: Paste token ke https://jwt.io untuk decode
- **Enable Auto-Save**: Postman → Settings → Auto-save requests

---

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
