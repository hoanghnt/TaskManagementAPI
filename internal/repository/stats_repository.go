package repository

import (
	"time"

	"github.com/hoanghnt/TaskManagementAPI/internal/models"
	"gorm.io/gorm"
)

type StatsRepository struct {
	db *gorm.DB
}

// NewStatsRepository creates a new stats repository
func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// GetTaskStats retrieves comprehensive task statistics for a user
func (r *StatsRepository) GetTaskStats(userID uint) (*models.TaskStats, error) {
	stats := &models.TaskStats{
		ByStatus:   make(map[string]int64),
		ByPriority: make(map[string]int64),
	}

	// 1. Get total tasks count
	if err := r.db.Model(&models.Task{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalTasks).Error; err != nil {
		return nil, err
	}

	// 2. Get count by status
	type StatusCount struct {
		Status string
		Count  int64
	}
	var statusCounts []StatusCount
	if err := r.db.Model(&models.Task{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, err
	}
	for _, sc := range statusCounts {
		stats.ByStatus[sc.Status] = sc.Count
	}

	// 3. Get count by priority
	type PriorityCount struct {
		Priority string
		Count    int64
	}
	var priorityCounts []PriorityCount
	if err := r.db.Model(&models.Task{}).
		Select("priority, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("priority").
		Scan(&priorityCounts).Error; err != nil {
		return nil, err
	}
	for _, pc := range priorityCounts {
		stats.ByPriority[pc.Priority] = pc.Count
	}

	// 4. Get count by category
	var categoryCounts []models.CategoryTaskCount
	if err := r.db.Model(&models.Task{}).
		Select("categories.id as category_id, categories.name as category_name, COUNT(tasks.id) as task_count").
		Joins("LEFT JOIN categories ON tasks.category_id = categories.id").
		Where("tasks.user_id = ?", userID).
		Where("categories.id IS NOT NULL").
		Group("categories.id, categories.name").
		Order("task_count DESC").
		Scan(&categoryCounts).Error; err != nil {
		return nil, err
	}
	stats.ByCategory = categoryCounts

	// 5. Calculate completion rate
	completedCount := stats.ByStatus["completed"]
	if stats.TotalTasks > 0 {
		stats.CompletionRate = float64(completedCount) / float64(stats.TotalTasks) * 100
	}

	// 6. Get overdue tasks count
	now := time.Now()
	if err := r.db.Model(&models.Task{}).
		Where("user_id = ? AND due_date < ? AND status != ?", userID, now, models.TaskStatusCompleted).
		Count(&stats.OverdueTasks).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// GetUpcomingTasks retrieves tasks due in the next N days
func (r *StatsRepository) GetUpcomingTasks(userID uint, days int) ([]models.Task, error) {
	var tasks []models.Task
	now := time.Now()
	futureDate := now.AddDate(0, 0, days)

	err := r.db.Preload("Category").
		Where("user_id = ? AND due_date BETWEEN ? AND ? AND status != ?",
			userID, now, futureDate, models.TaskStatusCompleted).
		Order("due_date ASC").
		Find(&tasks).Error

	return tasks, err
}

// GetOverdueTasks retrieves overdue tasks
func (r *StatsRepository) GetOverdueTasks(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	now := time.Now()

	err := r.db.Preload("Category").
		Where("user_id = ? AND due_date < ? AND status != ?",
			userID, now, models.TaskStatusCompleted).
		Order("due_date ASC").
		Find(&tasks).Error

	return tasks, err
}
