package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hoanghnt/TaskManagementAPI/internal/models"
	"github.com/hoanghnt/TaskManagementAPI/internal/services"
	"github.com/hoanghnt/TaskManagementAPI/internal/utils"
)

type TaskHandler struct {
	taskService *services.TaskService
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body models.CreateTaskRequest true "Task creation data"
// @Success 201 {object} map[string]interface{} "Task created successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse request body
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Create task
	task, err := h.taskService.CreateTask(userID.(uint), req)
	if err != nil {
		if err.Error() == "category not found" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if err.Error() == "due date cannot be in the past" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create task")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Task created successfully", task)
}

// GetAllTasks godoc
// @Summary Get all tasks
// @Description Get all tasks for the authenticated user with filtering, sorting, and pagination
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status (pending, in_progress, completed)"
// @Param priority query string false "Filter by priority (low, medium, high)"
// @Param category_id query int false "Filter by category ID"
// @Param search query string false "Search in title and description"
// @Param sort_by query string false "Sort by field (created_at, updated_at, due_date, priority)" default(created_at)
// @Param sort_order query string false "Sort order (asc, desc)" default(desc)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size (max 100)" default(10)
// @Success 200 {object} map[string]interface{} "Tasks retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /tasks [get]
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse query parameters
	var filter models.TaskFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Get tasks with filtering
	tasks, total, err := h.taskService.GetAllTasks(userID.(uint), filter)
	if err != nil {
		if err.Error() == "category not found" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}

	// Set defaults if not already set
	filter.SetDefaults()

	// Return paginated response
	utils.PaginatedResponse(c, http.StatusOK, "Tasks retrieved successfully", tasks, total, filter.Page, filter.PageSize)
}

// GetTaskByID godoc
// @Summary Get task by ID
// @Description Get a specific task by ID for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]interface{} "Task retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid task ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Task not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse task ID from URL
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Get task
	task, err := h.taskService.GetTaskByID(uint(taskID), userID.(uint))
	if err != nil {
		if err.Error() == "task not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task retrieved successfully", task)
}

// UpdateTask godoc
// @Summary Update task
// @Description Update an existing task (partial update)
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param task body models.UpdateTaskRequest true "Task update data"
// @Success 200 {object} map[string]interface{} "Task updated successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Task not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse task ID from URL
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Parse request body
	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Update task
	task, err := h.taskService.UpdateTask(uint(taskID), userID.(uint), req)
	if err != nil {
		if err.Error() == "task not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "category not found" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if err.Error() == "due date cannot be in the past" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task updated successfully", task)
}

// UpdateTaskStatus godoc
// @Summary Update task status
// @Description Update only the status of a task (quick status change)
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param status body models.UpdateTaskStatusRequest true "Status update data"
// @Success 200 {object} map[string]interface{} "Task status updated successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Task not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /tasks/{id}/status [patch]
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse task ID from URL
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Parse request body
	var req models.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Update status
	task, err := h.taskService.UpdateTaskStatus(uint(taskID), userID.(uint), req.Status)
	if err != nil {
		if err.Error() == "task not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update task status")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task status updated successfully", task)
}

// DeleteTask godoc
// @Summary Delete task
// @Description Soft delete a task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]interface{} "Task deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid task ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Task not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse task ID from URL
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Delete task
	if err := h.taskService.DeleteTask(uint(taskID), userID.(uint)); err != nil {
		if err.Error() == "task not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task deleted successfully", nil)
}

// BulkUpdateStatus godoc
// @Summary Bulk update task status
// @Description Update status for multiple tasks at once
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.BulkUpdateStatusRequest true "Bulk status update data"
// @Success 200 {object} map[string]interface{} "Bulk update completed"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /tasks/bulk/status [patch]
func (h *TaskHandler) BulkUpdateStatus(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse request body
	var req models.BulkUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Bulk update
	response, err := h.taskService.BulkUpdateStatus(userID.(uint), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Bulk status update completed", response)
}
