package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	Color       string         `gorm:"size:7" json:"color"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Tasks       []Task         `gorm:"foreignKey:CategoryID" json:"-"`

	// Computed fields (not stored in DB)
	TaskCount      int64 `gorm:"-" json:"task_count,omitempty"`
	PendingCount   int64 `gorm:"-" json:"pending_count,omitempty"`
	CompletedCount int64 `gorm:"-" json:"completed_count,omitempty"`
}

// CreateCategoryRequest represents category creation input
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=255"`
	Color       string `json:"color" binding:"omitempty,len=7"` // #RRGGBB format
}

// UpdateCategoryRequest represents category update input
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Description string `json:"description" binding:"omitempty,max=255"`
	Color       string `json:"color" binding:"omitempty,len=7"`
}
