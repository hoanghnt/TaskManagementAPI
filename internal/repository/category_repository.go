package repository

import (
	"errors"

	"github.com/hoanghnt/TaskManagementAPI/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// FindByID finds a category by ID for a specific user
func (r *CategoryRepository) FindByID(id uint, userID uint) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("User").Where("id = ? AND user_id = ?", id, userID).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// FindAllByUser finds all categories for a specific user with pagination
func (r *CategoryRepository) FindAllByUser(userID uint, page, pageSize int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	// Count total items
	if err := r.db.Model(&models.Category{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get paginated results
	err := r.db.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&categories).Error

	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// Update updates a category
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete soft deletes a category
func (r *CategoryRepository) Delete(id uint, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Category{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

// ExistsByID checks if a category exists for a specific user
func (r *CategoryRepository) ExistsByID(id uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).
		Where("id = ? AND user_id = ?", id, userID).
		Count(&count).Error
	return count > 0, err
}

// ExistsByNameAndUser checks if a category with the same name exists for a user
func (r *CategoryRepository) ExistsByNameAndUser(name string, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).
		Where("name = ? AND user_id = ?", name, userID).
		Count(&count).Error
	return count > 0, err
}

// ExistsByNameAndUserExcludingID checks if a category with the same name exists for a user, excluding a specific ID
func (r *CategoryRepository) ExistsByNameAndUserExcludingID(name string, userID uint, excludeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).
		Where("name = ? AND user_id = ? AND id != ?", name, userID, excludeID).
		Count(&count).Error
	return count > 0, err
}

// GetCategoriesWithTaskCount retrieves categories with task counts
func (r *CategoryRepository) GetCategoriesWithTaskCount(userID uint, page, pageSize int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	// Count total items
	if err := r.db.Model(&models.Category{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get categories
	err := r.db.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&categories).Error

	if err != nil {
		return nil, 0, err
	}

	// Load task counts for each category
	for i := range categories {
		// Total task count
		r.db.Model(&models.Task{}).
			Where("category_id = ?", categories[i].ID).
			Count(&categories[i].TaskCount)

		// Pending count
		r.db.Model(&models.Task{}).
			Where("category_id = ? AND status = ?", categories[i].ID, models.TaskStatusPending).
			Count(&categories[i].PendingCount)

		// Completed count
		r.db.Model(&models.Task{}).
			Where("category_id = ? AND status = ?", categories[i].ID, models.TaskStatusCompleted).
			Count(&categories[i].CompletedCount)
	}

	return categories, total, nil
}
