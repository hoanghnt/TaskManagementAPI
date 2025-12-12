package services

import (
	"errors"
	"strings"

	"github.com/hoanghnt/TaskManagementAPI/internal/models"
	"github.com/hoanghnt/TaskManagementAPI/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// Create creates a new category
func (s *CategoryService) Create(req *models.CreateCategoryRequest, userID uint) (*models.Category, error) {
	// Validate and sanitize input
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.Color = strings.TrimSpace(req.Color)

	if req.Name == "" {
		return nil, errors.New("category name is required")
	}

	// Check if category with same name exists for this user
	exists, err := s.categoryRepo.ExistsByNameAndUser(req.Name, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("category with this name already exists")
	}

	// Create category
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		UserID:      userID,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, errors.New("failed to create category")
	}

	return category, nil
}

// GetByID retrieves a category by ID for a specific user
func (s *CategoryService) GetByID(id uint, userID uint) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// GetAllByUser retrieves all categories for a user with pagination
func (s *CategoryService) GetAllByUser(userID uint, page, pageSize int) ([]models.Category, int64, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Use the new method with task counts
	categories, total, err := s.categoryRepo.GetCategoriesWithTaskCount(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// Update updates a category
func (s *CategoryService) Update(id uint, req *models.UpdateCategoryRequest, userID uint) (*models.Category, error) {
	// Find existing category
	category, err := s.categoryRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		req.Name = strings.TrimSpace(req.Name)

		// Check if new name conflicts with existing category
		exists, err := s.categoryRepo.ExistsByNameAndUserExcludingID(req.Name, userID, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("category with this name already exists")
		}

		category.Name = req.Name
	}

	if req.Description != "" {
		category.Description = strings.TrimSpace(req.Description)
	}

	if req.Color != "" {
		category.Color = strings.TrimSpace(req.Color)
	}

	// Save updates
	if err := s.categoryRepo.Update(category); err != nil {
		return nil, errors.New("failed to update category")
	}

	return category, nil
}

// Delete deletes a category
func (s *CategoryService) Delete(id uint, userID uint) error {
	return s.categoryRepo.Delete(id, userID)
}
