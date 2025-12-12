package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hoanghnt/TaskManagementAPI/internal/services"
	"github.com/hoanghnt/TaskManagementAPI/internal/utils"
)

type StatsHandler struct {
	statsService *services.StatsService
}

// NewStatsHandler creates a new stats handler
func NewStatsHandler(statsService *services.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

// GetDashboardStats godoc
// @Summary Get dashboard statistics
// @Description Get comprehensive task statistics for dashboard
// @Tags Statistics
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.TaskStats
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stats/dashboard [get]
func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get statistics
	stats, err := h.statsService.GetDashboardStats(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve statistics")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Statistics retrieved successfully", stats)
}

// GetUpcomingTasks godoc
// @Summary Get upcoming tasks
// @Description Get tasks due in the next N days
// @Tags Statistics
// @Produce json
// @Security BearerAuth
// @Param days query int false "Number of days (default: 7)"
// @Success 200 {array} models.Task
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stats/upcoming [get]
func (h *StatsHandler) GetUpcomingTasks(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse days parameter
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	// Get upcoming tasks
	tasks, err := h.statsService.GetUpcomingTasks(userID.(uint), days)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve upcoming tasks")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Upcoming tasks retrieved successfully", tasks)
}

// GetOverdueTasks godoc
// @Summary Get overdue tasks
// @Description Get tasks that are past their due date
// @Tags Statistics
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Task
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stats/overdue [get]
func (h *StatsHandler) GetOverdueTasks(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get overdue tasks
	tasks, err := h.statsService.GetOverdueTasks(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve overdue tasks")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Overdue tasks retrieved successfully", tasks)
}
