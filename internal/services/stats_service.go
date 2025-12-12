package services

import (
	"github.com/hoanghnt/TaskManagementAPI/internal/models"
	"github.com/hoanghnt/TaskManagementAPI/internal/repository"
)

type StatsService struct {
	statsRepo *repository.StatsRepository
}

// NewStatsService creates a new stats service
func NewStatsService(statsRepo *repository.StatsRepository) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
	}
}

// GetDashboardStats retrieves comprehensive dashboard statistics
func (s *StatsService) GetDashboardStats(userID uint) (*models.TaskStats, error) {
	return s.statsRepo.GetTaskStats(userID)
}

// GetUpcomingTasks retrieves tasks due in the next N days
func (s *StatsService) GetUpcomingTasks(userID uint, days int) ([]models.Task, error) {
	if days <= 0 {
		days = 7 // default to 7 days
	}
	return s.statsRepo.GetUpcomingTasks(userID, days)
}

// GetOverdueTasks retrieves overdue tasks
func (s *StatsService) GetOverdueTasks(userID uint) ([]models.Task, error) {
	return s.statsRepo.GetOverdueTasks(userID)
}
