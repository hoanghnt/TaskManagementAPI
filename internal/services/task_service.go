package services

import (
	"errors"
	"time"

	"github.com/hoanghnt/TaskManagementAPI/internal/models"
	"github.com/hoanghnt/TaskManagementAPI/internal/repository"
	"github.com/hoanghnt/TaskManagementAPI/internal/utils"
)

type TaskService struct {
	taskRepo     *repository.TaskRepository
	categoryRepo *repository.CategoryRepository
}

// NewTaskService creates a new task service
func NewTaskService(taskRepo *repository.TaskRepository, categoryRepo *repository.CategoryRepository) *TaskService {
	return &TaskService{
		taskRepo:     taskRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(userID uint, req models.CreateTaskRequest) (*models.Task, error) {
	// Validate category if provided
	if req.CategoryID != nil {
		exists, err := s.categoryRepo.ExistsByID(*req.CategoryID, userID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, utils.ErrCategoryNotFound
		}
	}

	// Validate due date (optional: cannot be in the past)
	if req.DueDate != nil && req.DueDate.Before(time.Now()) {
		return nil, errors.New("due date cannot be in the past")
	}

	// Set default values if not provided
	status := req.Status
	if status == "" {
		status = models.TaskStatusPending
	}

	priority := req.Priority
	if priority == "" {
		priority = models.TaskPriorityMedium
	}

	// Create task
	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    priority,
		DueDate:     req.DueDate,
		UserID:      userID,
		CategoryID:  req.CategoryID,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	// Reload with relationships
	return s.taskRepo.FindByID(task.ID, userID)
}

// GetTaskByID retrieves a task by ID
func (s *TaskService) GetTaskByID(id uint, userID uint) (*models.Task, error) {
	task, err := s.taskRepo.FindByID(id, userID)
	if err != nil {
		return nil, utils.ErrTaskNotFound
	}
	return task, nil
}

// GetAllTasks retrieves all tasks with filtering and pagination
func (s *TaskService) GetAllTasks(userID uint, filter models.TaskFilter) ([]models.Task, int64, error) {
	// Set default values for pagination
	filter.SetDefaults()

	// Validate category if provided in filter
	if filter.CategoryID > 0 {
		exists, err := s.categoryRepo.ExistsByID(filter.CategoryID, userID)
		if err != nil {
			return nil, 0, err
		}
		if !exists {
			return nil, 0, utils.ErrCategoryNotFound
		}
	}

	return s.taskRepo.FindAllByUser(userID, filter)
}

// UpdateTask updates an existing task
func (s *TaskService) UpdateTask(id uint, userID uint, req models.UpdateTaskRequest) (*models.Task, error) {
	// Check if task exists
	task, err := s.taskRepo.FindByID(id, userID)
	if err != nil {
		return nil, utils.ErrTaskNotFound
	}

	// Validate category if being updated
	if req.CategoryID != nil {
		// Allow null category (set to 0 to remove category)
		if *req.CategoryID > 0 {
			exists, err := s.categoryRepo.ExistsByID(*req.CategoryID, userID)
			if err != nil {
				return nil, err
			}
			if !exists {
				return nil, utils.ErrCategoryNotFound
			}
		}
		task.CategoryID = req.CategoryID
	}

	// Validate due date if being updated
	if req.DueDate != nil && req.DueDate.Before(time.Now()) {
		return nil, errors.New("due date cannot be in the past")
	}

	// Update fields only if provided (partial update)
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	// Save updates
	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	// Reload with relationships
	return s.taskRepo.FindByID(task.ID, userID)
}

// UpdateTaskStatus updates only the status of a task
func (s *TaskService) UpdateTaskStatus(id uint, userID uint, status models.TaskStatus) (*models.Task, error) {
	// Check if task exists
	exists, err := s.taskRepo.ExistsByID(id, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, utils.ErrTaskNotFound
	}

	// Update status
	if err := s.taskRepo.UpdateStatus(id, userID, status); err != nil {
		return nil, err
	}

	// Return updated task
	return s.taskRepo.FindByID(id, userID)
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(id uint, userID uint) error {
	// Check if task exists
	exists, err := s.taskRepo.ExistsByID(id, userID)
	if err != nil {
		return err
	}
	if !exists {
		return utils.ErrTaskNotFound
	}

	return s.taskRepo.Delete(id, userID)
}

// BulkUpdateStatus updates status for multiple tasks
func (s *TaskService) BulkUpdateStatus(userID uint, req models.BulkUpdateStatusRequest) (*models.BulkUpdateResponse, error) {
	// Validate task IDs not empty
	if len(req.TaskIDs) == 0 {
		return nil, errors.New("task IDs cannot be empty")
	}

	// Update status for all tasks
	rowsAffected, err := s.taskRepo.BulkUpdateStatus(req.TaskIDs, userID, req.Status)
	if err != nil {
		return nil, err
	}

	// Calculate success/failed counts
	totalCount := len(req.TaskIDs)
	successCount := int(rowsAffected)
	failedCount := totalCount - successCount

	response := &models.BulkUpdateResponse{
		SuccessCount: successCount,
		FailedCount:  failedCount,
		TotalCount:   totalCount,
		Message:      "Bulk status update completed",
	}

	return response, nil
}
