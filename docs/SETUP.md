# 🚀 Setup Guide - Task Management API

Chi tiết hướng dẫn setup project từ đầu để bắt đầu phát triển.

## 📋 Checklist Chuẩn Bị

### 1. Cài Đặt Phần Mềm Cần Thiết

#### ✅ Go (Golang)
- **Version:** 1.21 trở lên
- **Download:** https://golang.org/dl/
- **Kiểm tra:**
```bash
go version
# Kết quả: go version go1.21.x windows/amd64
```

#### ✅ PostgreSQL
- **Version:** 14 trở lên
- **Download:** https://www.postgresql.org/download/
- **Hoặc dùng Docker:**
```bash
docker run --name postgres-taskapi -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:14
```
- **Kiểm tra:**
```bash
psql --version
# Kết quả: psql (PostgreSQL) 14.x
```

#### ✅ Git
- **Download:** https://git-scm.com/downloads
- **Kiểm tra:**
```bash
git --version
```

#### ✅ IDE/Editor (Tùy chọn)
- **Visual Studio Code** + Go extension (khuyên dùng)
- **GoLand** by JetBrains
- **Vim/Neovim** + Go plugins

#### ✅ API Testing Tool
- **Postman** - https://www.postman.com/downloads/
- **Thunder Client** - VS Code extension
- **Insomnia** - https://insomnia.rest/download

---

## 🔧 Bước 1: Setup PostgreSQL Database

### Option A: Using psql (Command Line)

1. **Kết nối PostgreSQL:**
```bash
# Windows (PowerShell)
psql -U postgres

# Linux/Mac
sudo -u postgres psql
```

2. **Tạo Database:**
```sql
CREATE DATABASE taskmanagement;

-- Kiểm tra
\l

-- Kết nối vào database
\c taskmanagement

-- Thoát
\q
```

### Option B: Using pgAdmin

1. Mở pgAdmin
2. Kết nối đến PostgreSQL server
3. Right-click "Databases" → "Create" → "Database..."
4. Nhập tên: `taskmanagement`
5. Click "Save"

### Option C: Using Docker

```bash
# Start PostgreSQL container
docker run --name postgres-taskapi \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=taskmanagement \
  -p 5432:5432 \
  -d postgres:14

# Kiểm tra container đang chạy
docker ps

# Kết nối vào database
docker exec -it postgres-taskapi psql -U postgres -d taskmanagement
```

---

## 📦 Bước 2: Clone và Setup Project

### 1. Clone Repository (hoặc tạo mới)

```bash
# Nếu clone từ Git
git clone https://github.com/yourusername/TaskManagementAPI.git
cd TaskManagementAPI

# Nếu tạo mới (project này)
# Bạn đã có folder này rồi
cd TaskManagementAPI
```

### 2. Install Dependencies

```bash
# Download tất cả dependencies
go mod download

# Hoặc
go mod tidy
```

### 3. Configure Environment Variables

```bash
# Copy file example
cp .env.example .env

# Hoặc trên Windows PowerShell
Copy-Item .env.example .env
```

**Chỉnh sửa file `.env`:**

```env
# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres          # ⚠️ Đổi thành password của bạn
DB_NAME=taskmanagement

# JWT Configuration
JWT_SECRET=my-super-secret-jwt-key-change-this-min-32-characters-long
JWT_EXPIRY_HOURS=24

# CORS Configuration (Optional)
CORS_ALLOW_ORIGINS=http://localhost:3000,http://localhost:5173
```

⚠️ **LƯU Ý QUAN TRỌNG:**
- `DB_PASSWORD`: Đổi thành password PostgreSQL của bạn
- `JWT_SECRET`: Đổi thành chuỗi bí mật của riêng bạn (tối thiểu 32 ký tự)
- **KHÔNG** commit file `.env` lên Git (đã có trong `.gitignore`)

---

## 🔨 Bước 3: Cấu Trúc Project (Đã Tạo)

```
TaskManagementAPI/
├── cmd/
│   └── api/              # ⏳ Sẽ tạo ở Phase 1
│       └── main.go
├── internal/
│   ├── config/           # ⏳ Sẽ tạo ở Phase 1
│   │   └── config.go
│   ├── models/           # ⏳ Sẽ tạo ở Phase 1
│   │   ├── user.go
│   │   ├── task.go
│   │   └── category.go
│   ├── database/         # ⏳ Sẽ tạo ở Phase 1
│   │   └── database.go
│   ├── middleware/       # ⏳ Sẽ tạo ở Phase 2
│   │   └── auth.go
│   ├── handlers/         # ⏳ Sẽ tạo ở Phase 2-4
│   │   ├── auth_handler.go
│   │   ├── task_handler.go
│   │   └── category_handler.go
│   ├── repository/       # ⏳ Sẽ tạo ở Phase 2-4
│   │   ├── user_repository.go
│   │   ├── task_repository.go
│   │   └── category_repository.go
│   └── services/         # ⏳ Sẽ tạo ở Phase 2-4
│       ├── auth_service.go
│       ├── task_service.go
│       └── category_service.go
├── docs/                 # ⏳ Swagger (Phase 5)
├── .env                  # ✅ Đã tạo
├── .env.example          # ✅ Đã tạo
├── .gitignore            # ✅ Đã tạo
├── README.md             # ✅ Đã tạo
├── API_DESIGN.md         # ✅ Đã tạo
├── SETUP.md              # ✅ Đang đọc file này
├── go.mod                # ✅ Đã tạo
└── go.sum                # ✅ Đã tạo
```

---

## ✅ Bước 4: Verify Setup

### 1. Test Go Installation

```bash
go version
go env
```

### 2. Test Database Connection

```bash
# Kết nối PostgreSQL
psql -U postgres -d taskmanagement

# Trong psql, chạy:
SELECT version();
\dt
\q
```

### 3. Check Dependencies

```bash
go mod verify
```

---

## 🎯 Bước 5: Sẵn Sàng Phát Triển

Bạn đã hoàn thành Phase 0! ✅

### Tiếp Theo - Development Roadmap:

#### **Phase 1: Foundation** (3-4 giờ)
- [ ] Config management
- [ ] Database connection
- [ ] Models definition
- [ ] Auto-migration
- [ ] Basic error handling

#### **Phase 2: Authentication** (4-5 giờ)
- [ ] User repository
- [ ] Auth service
- [ ] JWT middleware
- [ ] Register/Login endpoints

#### **Phase 3: Categories CRUD** (3-4 giờ)
- [ ] Category repository
- [ ] Category service
- [ ] Category handlers
- [ ] CRUD operations

#### **Phase 4: Tasks CRUD** (5-6 giờ)
- [ ] Task repository
- [ ] Task service
- [ ] Task handlers
- [ ] Filtering & sorting
- [ ] Pagination

#### **Phase 5: Documentation** (2-3 giờ)
- [ ] Swagger annotations
- [ ] Generate docs
- [ ] Testing
- [ ] Polish

---

## 🐛 Troubleshooting

### Problem: PostgreSQL connection refused

**Solution:**
1. Kiểm tra PostgreSQL đang chạy:
```bash
# Windows
Get-Service postgresql*

# Linux
sudo systemctl status postgresql
```

2. Kiểm tra port 5432:
```bash
netstat -an | findstr 5432
```

### Problem: Go modules errors

**Solution:**
```bash
# Clean cache
go clean -modcache

# Re-download
go mod download
```

### Problem: Permission denied on PostgreSQL

**Solution:**
```bash
# Reset password
ALTER USER postgres WITH PASSWORD 'new_password';
```

### Problem: Port 8080 already in use

**Solution:**
Đổi `SERVER_PORT` trong file `.env`:
```env
SERVER_PORT=8081
```

---

## 📚 Tài Liệu Tham Khảo

### Go Documentation
- **Official Docs:** https://go.dev/doc/
- **Go by Example:** https://gobyexample.com/
- **Effective Go:** https://go.dev/doc/effective_go

### Framework & Libraries
- **Gin:** https://gin-gonic.com/docs/
- **GORM:** https://gorm.io/docs/
- **JWT Go:** https://github.com/golang-jwt/jwt

### Database
- **PostgreSQL Docs:** https://www.postgresql.org/docs/

### Best Practices
- **Go Code Review Comments:** https://github.com/golang/go/wiki/CodeReviewComments
- **Project Layout:** https://github.com/golang-standards/project-layout

---

## 🎓 Learning Resources

### Video Tutorials
- **Learn Go with Tests:** https://quii.gitbook.io/learn-go-with-tests/
- **Go Tutorial (FreeCodeCamp):** YouTube

### Articles
- **Building RESTful APIs in Go**
- **JWT Authentication in Go**
- **GORM Best Practices**

---

## ✅ Checklist Cuối Cùng

Trước khi bắt đầu code Phase 1:

```
□ Go đã cài đặt và test thành công
□ PostgreSQL đã cài đặt và chạy
□ Database 'taskmanagement' đã được tạo
□ File .env đã được config đúng
□ Dependencies đã được download
□ Đã đọc API_DESIGN.md
□ Đã hiểu cấu trúc project
□ Đã test connection database
□ IDE/Editor đã setup xong
□ Postman hoặc testing tool đã cài
```

---

## 🚀 Ready to Code!

Bây giờ bạn đã sẵn sàng bắt đầu **Phase 1: Foundation**!

Chạy lệnh để xác nhận mọi thứ OK:
```bash
go version
psql -U postgres -d taskmanagement -c "SELECT 1;"
go mod verify
```

Nếu tất cả đều thành công ✅, bạn có thể bắt đầu code!

**Next Step:** Phase 1 - Tạo Config, Database Connection, và Models

Good luck! 🎉

