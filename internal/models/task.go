package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
)

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null;size:200" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Status      TaskStatus     `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Priority    TaskPriority   `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	DueDate     *time.Time     `json:"due_date,omitempty"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
	CategoryID  *uint          `json:"category_id,omitempty"`
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateTaskRequest represents task creation input
type CreateTaskRequest struct {
	Title       string       `json:"title" binding:"required,max=200"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
	Priority    TaskPriority `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time   `json:"due_date"`
	CategoryID  *uint        `json:"category_id"`
}

// UpdateTaskRequest represents task update input
type UpdateTaskRequest struct {
	Title       string       `json:"title" binding:"omitempty,max=200"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
	Priority    TaskPriority `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time   `json:"due_date"`
	CategoryID  *uint        `json:"category_id"`
}

// UpdateTaskStatusRequest represents task status update input
type UpdateTaskStatusRequest struct {
	Status TaskStatus `json:"status" binding:"required,oneof=pending in_progress completed"`
}

// TaskFilter represents query parameters for filtering tasks
type TaskFilter struct {
	Status     string `form:"status" binding:"omitempty,oneof=pending in_progress completed"`
	Priority   string `form:"priority" binding:"omitempty,oneof=low medium high"`
	CategoryID uint   `form:"category_id"`
	Search     string `form:"search"` // search in title and description
	SortBy     string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at due_date priority"`
	SortOrder  string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
	Page       int    `form:"page" binding:"omitempty,min=1"`
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// BulkUpdateStatusRequest represents bulk status update input
type BulkUpdateStatusRequest struct {
	TaskIDs []uint     `json:"task_ids" binding:"required,min=1"`
	Status  TaskStatus `json:"status" binding:"required,oneof=pending in_progress completed"`
}

// BulkUpdateResponse represents bulk operation response
type BulkUpdateResponse struct {
	SuccessCount int    `json:"success_count"`
	FailedCount  int    `json:"failed_count"`
	TotalCount   int    `json:"total_count"`
	Message      string `json:"message"`
}

// SetDefaults sets default values for pagination
func (f *TaskFilter) SetDefaults() {
	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}
	if f.SortBy == "" {
		f.SortBy = "created_at"
	}
	if f.SortOrder == "" {
		f.SortOrder = "desc"
	}
}
