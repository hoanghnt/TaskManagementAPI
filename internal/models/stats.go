package models

// TaskStats represents task statistics for dashboard
type TaskStats struct {
	TotalTasks     int64               `json:"total_tasks"`
	ByStatus       map[string]int64    `json:"by_status"`
	ByPriority     map[string]int64    `json:"by_priority"`
	ByCategory     []CategoryTaskCount `json:"by_category"`
	CompletionRate float64             `json:"completion_rate"`
	OverdueTasks   int64               `json:"overdue_tasks"`
}

// CategoryTaskCount represents task count per category
type CategoryTaskCount struct {
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	TaskCount    int64  `json:"task_count"`
}
