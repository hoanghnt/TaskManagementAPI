# 📋 API Design Documentation

## Base URL
```
http://localhost:8080/api/v1
```

---

## 🔐 Authentication Endpoints

### 1. Register User

**Endpoint:** `POST /auth/register`

**Description:** Register a new user account

**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "full_name": "John Doe"
}
```

**Validation Rules:**
- `username`: required, min 3 chars, max 50 chars, unique
- `email`: required, valid email format, unique
- `password`: required, min 6 chars
- `full_name`: required

**Success Response (201):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "created_at": "2024-01-10T10:00:00Z",
    "updated_at": "2024-01-10T10:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Validation errors
- `409 Conflict`: Username or email already exists

---

### 2. Login User

**Endpoint:** `POST /auth/login`

**Description:** Authenticate user and receive JWT token

**Request Body:**
```json
{
  "username": "johndoe",
  "password": "password123"
}
```

**Success Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "created_at": "2024-01-10T10:00:00Z",
    "updated_at": "2024-01-10T10:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Missing fields
- `401 Unauthorized`: Invalid credentials

---

### 3. Get Current User

**Endpoint:** `GET /auth/me`

**Description:** Get current authenticated user information

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Success Response (200):**
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "full_name": "John Doe",
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T10:00:00Z"
}
```

**Error Responses:**
- `401 Unauthorized`: Invalid or missing token

---

## 📂 Category Endpoints

> **All category endpoints require authentication**

### 1. Get All Categories

**Endpoint:** `GET /categories`

**Description:** Get all categories for the authenticated user with pagination

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10, max: 100)

**Success Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Work",
      "description": "Work-related tasks",
      "color": "#FF5733",
      "user_id": 1,
      "created_at": "2024-01-10T10:00:00Z",
      "updated_at": "2024-01-10T10:00:00Z"
    },
    {
      "id": 2,
      "name": "Personal",
      "description": "Personal tasks",
      "color": "#33FF57",
      "user_id": 1,
      "created_at": "2024-01-10T10:00:00Z",
      "updated_at": "2024-01-10T10:00:00Z"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total_items": 2,
  "total_pages": 1
}
```

---

### 2. Create Category

**Endpoint:** `POST /categories`

**Description:** Create a new category

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "name": "Work",
  "description": "Work-related tasks",
  "color": "#FF5733"
}
```

**Validation Rules:**
- `name`: required, max 100 chars
- `description`: optional, max 255 chars
- `color`: optional, must be 7 chars (#RRGGBB format)

**Success Response (201):**
```json
{
  "id": 1,
  "name": "Work",
  "description": "Work-related tasks",
  "color": "#FF5733",
  "user_id": 1,
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T10:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Validation errors
- `401 Unauthorized`: Invalid token

---

### 3. Get Category by ID

**Endpoint:** `GET /categories/:id`

**Description:** Get a specific category by ID

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Success Response (200):**
```json
{
  "id": 1,
  "name": "Work",
  "description": "Work-related tasks",
  "color": "#FF5733",
  "user_id": 1,
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T10:00:00Z"
}
```

**Error Responses:**
- `404 Not Found`: Category not found or not owned by user

---

### 4. Update Category

**Endpoint:** `PUT /categories/:id`

**Description:** Update an existing category

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "name": "Updated Work",
  "description": "Updated description",
  "color": "#FF6633"
}
```

**Success Response (200):**
```json
{
  "id": 1,
  "name": "Updated Work",
  "description": "Updated description",
  "color": "#FF6633",
  "user_id": 1,
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T11:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Validation errors
- `404 Not Found`: Category not found or not owned by user

---

### 5. Delete Category

**Endpoint:** `DELETE /categories/:id`

**Description:** Delete a category (soft delete)

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Success Response (200):**
```json
{
  "message": "Category deleted successfully"
}
```

**Error Responses:**
- `404 Not Found`: Category not found or not owned by user

---

## ✅ Task Endpoints

> **All task endpoints require authentication**

### 1. Get All Tasks

**Endpoint:** `GET /tasks`

**Description:** Get all tasks with filtering, sorting, and pagination

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `status` (optional): Filter by status (pending, in_progress, completed)
- `priority` (optional): Filter by priority (low, medium, high)
- `category_id` (optional): Filter by category ID
- `search` (optional): Search in title and description
- `sort_by` (optional): Sort field (created_at, updated_at, due_date, priority)
- `sort_order` (optional): Sort order (asc, desc) - default: desc
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10, max: 100)

**Example Request:**
```
GET /tasks?status=pending&priority=high&sort_by=due_date&sort_order=asc&page=1&page_size=10
```

**Success Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Complete project documentation",
      "description": "Write comprehensive README and API docs",
      "status": "in_progress",
      "priority": "high",
      "due_date": "2024-01-15T23:59:59Z",
      "user_id": 1,
      "category_id": 1,
      "category": {
        "id": 1,
        "name": "Work",
        "color": "#FF5733"
      },
      "created_at": "2024-01-10T10:00:00Z",
      "updated_at": "2024-01-10T10:00:00Z"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total_items": 1,
  "total_pages": 1
}
```

---

### 2. Create Task

**Endpoint:** `POST /tasks`

**Description:** Create a new task

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "title": "Complete project documentation",
  "description": "Write comprehensive README and API docs",
  "status": "pending",
  "priority": "high",
  "due_date": "2024-01-15T23:59:59Z",
  "category_id": 1
}
```

**Validation Rules:**
- `title`: required, max 200 chars
- `description`: optional
- `status`: optional, must be one of: pending, in_progress, completed (default: pending)
- `priority`: optional, must be one of: low, medium, high (default: medium)
- `due_date`: optional, must be valid ISO 8601 datetime
- `category_id`: optional, must exist and belong to user

**Success Response (201):**
```json
{
  "id": 1,
  "title": "Complete project documentation",
  "description": "Write comprehensive README and API docs",
  "status": "pending",
  "priority": "high",
  "due_date": "2024-01-15T23:59:59Z",
  "user_id": 1,
  "category_id": 1,
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T10:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Validation errors
- `404 Not Found`: Category not found

---

### 3. Get Task by ID

**Endpoint:** `GET /tasks/:id`

**Description:** Get a specific task by ID

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Success Response (200):**
```json
{
  "id": 1,
  "title": "Complete project documentation",
  "description": "Write comprehensive README and API docs",
  "status": "in_progress",
  "priority": "high",
  "due_date": "2024-01-15T23:59:59Z",
  "user_id": 1,
  "category_id": 1,
  "category": {
    "id": 1,
    "name": "Work",
    "color": "#FF5733"
  },
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T10:00:00Z"
}
```

**Error Responses:**
- `404 Not Found`: Task not found or not owned by user

---

### 4. Update Task

**Endpoint:** `PUT /tasks/:id`

**Description:** Update an existing task

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "title": "Updated title",
  "description": "Updated description",
  "status": "completed",
  "priority": "medium",
  "due_date": "2024-01-20T23:59:59Z",
  "category_id": 2
}
```

**Success Response (200):**
```json
{
  "id": 1,
  "title": "Updated title",
  "description": "Updated description",
  "status": "completed",
  "priority": "medium",
  "due_date": "2024-01-20T23:59:59Z",
  "user_id": 1,
  "category_id": 2,
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T12:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Validation errors
- `404 Not Found`: Task or category not found

---

### 5. Update Task Status (Quick Action)

**Endpoint:** `PATCH /tasks/:id/status`

**Description:** Quickly update only the task status

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "status": "completed"
}
```

**Success Response (200):**
```json
{
  "id": 1,
  "title": "Complete project documentation",
  "description": "Write comprehensive README and API docs",
  "status": "completed",
  "priority": "high",
  "due_date": "2024-01-15T23:59:59Z",
  "user_id": 1,
  "category_id": 1,
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T13:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid status value
- `404 Not Found`: Task not found or not owned by user

---

### 6. Delete Task

**Endpoint:** `DELETE /tasks/:id`

**Description:** Delete a task (soft delete)

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Success Response (200):**
```json
{
  "message": "Task deleted successfully"
}
```

**Error Responses:**
- `404 Not Found`: Task not found or not owned by user

---

## 📊 Common Response Formats

### Error Response Format

All error responses follow this format:

```json
{
  "error": "Error message description"
}
```

### Validation Error Response

```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "email",
      "message": "must be a valid email address"
    },
    {
      "field": "password",
      "message": "must be at least 6 characters"
    }
  ]
}
```

---

## 🔑 HTTP Status Codes

- `200 OK`: Successful GET, PUT, PATCH, DELETE
- `201 Created`: Successful POST (resource created)
- `400 Bad Request`: Invalid input, validation errors
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: Valid token but insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists (duplicate)
- `500 Internal Server Error`: Server error

---

## 🛡️ Authentication Flow

1. **Register or Login**: Get JWT token
2. **Store Token**: Save token in client (localStorage, cookie, etc.)
3. **Make Requests**: Include token in Authorization header
4. **Token Expiry**: Token expires after 24 hours (configurable)
5. **Refresh**: Login again to get new token

**Token Format:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 📝 Database Schema

### Users Table
```
id: integer (PK, auto-increment)
username: varchar(50) (unique, not null)
email: varchar(100) (unique, not null)
password: text (hashed, not null)
full_name: varchar(100)
created_at: timestamp
updated_at: timestamp
deleted_at: timestamp (nullable, for soft delete)
```

### Categories Table
```
id: integer (PK, auto-increment)
name: varchar(100) (not null)
description: varchar(255)
color: varchar(7)
user_id: integer (FK -> users.id, not null)
created_at: timestamp
updated_at: timestamp
deleted_at: timestamp (nullable)
```

### Tasks Table
```
id: integer (PK, auto-increment)
title: varchar(200) (not null)
description: text
status: enum('pending', 'in_progress', 'completed')
priority: enum('low', 'medium', 'high')
due_date: timestamp (nullable)
user_id: integer (FK -> users.id, not null)
category_id: integer (FK -> categories.id, nullable)
created_at: timestamp
updated_at: timestamp
deleted_at: timestamp (nullable)
```

### Relationships
- User has many Tasks (1:N)
- User has many Categories (1:N)
- Category has many Tasks (1:N)
- Task belongs to User (N:1)
- Task belongs to Category (N:1, optional)

---

## 🎯 Business Rules

1. **User Isolation**: Users can only see/modify their own tasks and categories
2. **Soft Delete**: Deleted resources are marked with `deleted_at` timestamp, not physically removed
3. **Category Assignment**: Tasks can exist without a category
4. **Default Values**: 
   - Task status defaults to "pending"
   - Task priority defaults to "medium"
5. **Pagination**: Maximum page size is 100 items
6. **Token Expiry**: JWT tokens expire after 24 hours (configurable)
7. **Password Security**: Passwords are hashed using bcrypt before storage

