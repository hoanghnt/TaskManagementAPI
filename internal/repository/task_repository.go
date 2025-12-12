package repository

import (
	"errors"
	"fmt"

	"github.com/hoanghnt/TaskManagementAPI/internal/models"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create creates a new task
func (r *TaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

// FindByID finds a task by ID for a specific user
func (r *TaskRepository) FindByID(id uint, userID uint) (*models.Task, error) {
	var task models.Task

	// Preload User and Category relationships
	err := r.db.Preload("User").
		Preload("Category").
		Where("id = ? AND user_id = ?", id, userID).
		First(&task).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &task, nil
}

// FindAllByUser finds all tasks for a specific user with advanced filtering
func (r *TaskRepository) FindAllByUser(userID uint, filter models.TaskFilter) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	// Start building query
	query := r.db.Model(&models.Task{}).Where("user_id = ?", userID)

	// Apply filters
	query = r.applyFilters(query, filter)

	// Count total items (before pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	query = r.applySorting(query, filter)

	// Apply pagination
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Limit(filter.PageSize).Offset(offset)

	// Preload relationships and execute query
	err := query.Preload("User").Preload("Category").Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// applyFilters applies dynamic filters to the query
func (r *TaskRepository) applyFilters(query *gorm.DB, filter models.TaskFilter) *gorm.DB {
	// Filter by status
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Filter by priority
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}

	// Filter by category
	if filter.CategoryID > 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	// Search in title and description
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	return query
}

// applySorting applies dynamic sorting to the query
func (r *TaskRepository) applySorting(query *gorm.DB, filter models.TaskFilter) *gorm.DB {
	// Build ORDER BY clause
	orderClause := fmt.Sprintf("%s %s", filter.SortBy, filter.SortOrder)
	return query.Order(orderClause)
}

// Update updates a task
func (r *TaskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

// UpdateStatus updates only the status of a task
func (r *TaskRepository) UpdateStatus(id uint, userID uint, status models.TaskStatus) error {
	result := r.db.Model(&models.Task{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}

// Delete soft deletes a task
func (r *TaskRepository) Delete(id uint, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}

// ExistsByID checks if a task exists for a specific user
func (r *TaskRepository) ExistsByID(id uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Task{}).
		Where("id = ? AND user_id = ?", id, userID).
		Count(&count).Error
	return count > 0, err
}

// CategoryExistsForUser checks if a category belongs to the user
func (r *TaskRepository) CategoryExistsForUser(categoryID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).
		Where("id = ? AND user_id = ?", categoryID, userID).
		Count(&count).Error
	return count > 0, err
}

// BulkUpdateStatus updates status for multiple tasks
func (r *TaskRepository) BulkUpdateStatus(taskIDs []uint, userID uint, status models.TaskStatus) (int64, error) {
	// Update only tasks that belong to the user
	result := r.db.Model(&models.Task{}).
		Where("id IN ? AND user_id = ?", taskIDs, userID).
		Update("status", status)

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// FindByIDs finds multiple tasks by IDs for a specific user
func (r *TaskRepository) FindByIDs(taskIDs []uint, userID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Preload("User").
		Preload("Category").
		Where("id IN ? AND user_id = ?", taskIDs, userID).
		Find(&tasks).Error
	return tasks, err
}
