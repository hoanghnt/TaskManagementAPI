package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/hoanghnt/TaskManagementAPI/internal/config"
	"github.com/hoanghnt/TaskManagementAPI/internal/database"
	"github.com/hoanghnt/TaskManagementAPI/internal/handlers"
	"github.com/hoanghnt/TaskManagementAPI/internal/middleware"
	"github.com/hoanghnt/TaskManagementAPI/internal/repository"
	"github.com/hoanghnt/TaskManagementAPI/internal/services"
	"github.com/hoanghnt/TaskManagementAPI/internal/utils"

	_ "github.com/hoanghnt/TaskManagementAPI/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Task Management API
// @version 1.0
// @description RESTful API for managing tasks and categories
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@taskmanagement.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize database
	if err := database.InitDB(&cfg.Database); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	// Initialize JWT
	utils.InitJWT(cfg.JWT.Secret)
	log.Println("✅ JWT initialized")

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.GetDB())

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWT.ExpiryHours)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Category initialization
	categoryRepo := repository.NewCategoryRepository(database.GetDB())

	categoryService := services.NewCategoryService(categoryRepo)

	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Task initialization
	taskRepo := repository.NewTaskRepository(database.GetDB())

	taskService := services.NewTaskService(taskRepo, categoryRepo)

	taskHandler := handlers.NewTaskHandler(taskService)

	// Stats initialization
	statsRepo := repository.NewStatsRepository(database.GetDB())

	statsService := services.NewStatsService(statsRepo)

	statsHandler := handlers.NewStatsHandler(statsService)

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(corsMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Task Management API is running",
		})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Auth routes
			protected.GET("/auth/me", authHandler.GetMe)

			categories := protected.Group("/categories")
			{
				categories.GET("", categoryHandler.GetAll)
				categories.POST("", categoryHandler.Create)
				categories.GET("/:id", categoryHandler.GetByID)
				categories.PUT("/:id", categoryHandler.Update)
				categories.DELETE("/:id", categoryHandler.Delete)
			}

			tasks := protected.Group("/tasks")
			{
				tasks.GET("", taskHandler.GetAllTasks)
				tasks.POST("", taskHandler.CreateTask)
				tasks.GET("/:id", taskHandler.GetTaskByID)
				tasks.PUT("/:id", taskHandler.UpdateTask)
				tasks.PATCH("/:id/status", taskHandler.UpdateTaskStatus)
				tasks.PATCH("/bulk/status", taskHandler.BulkUpdateStatus)
				tasks.DELETE("/:id", taskHandler.DeleteTask)
			}

			stats := protected.Group("/stats")
			{
				stats.GET("/dashboard", statsHandler.GetDashboardStats)
				stats.GET("/upcoming", statsHandler.GetUpcomingTasks)
				stats.GET("/overdue", statsHandler.GetOverdueTasks)
			}
		}

		// API info endpoint
		v1.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Task Management API v1",
				"version": "1.0.0",
				"status":  "active",
				"endpoints": gin.H{
					"auth": gin.H{
						"register": "POST /api/v1/auth/register",
						"login":    "POST /api/v1/auth/login",
						"me":       "GET /api/v1/auth/me (protected)",
					},
					"categories": gin.H{
						"list":   "GET /api/v1/categories (protected)",
						"create": "POST /api/v1/categories (protected)",
						"get":    "GET /api/v1/categories/:id (protected)",
						"update": "PUT /api/v1/categories/:id (protected)",
						"delete": "DELETE /api/v1/categories/:id (protected)",
					},
					"tasks": gin.H{
						"list":          "GET /api/v1/tasks (protected)",
						"create":        "POST /api/v1/tasks (protected)",
						"get":           "GET /api/v1/tasks/:id (protected)",
						"update":        "PUT /api/v1/tasks/:id (protected)",
						"update_status": "PATCH /api/v1/tasks/:id/status (protected)",
						"delete":        "DELETE /api/v1/tasks/:id (protected)",
					},
					"stats": gin.H{
						"dashboard": "GET /api/v1/stats/dashboard (protected)",
						"upcoming":  "GET /api/v1/stats/upcoming (protected)",
						"overdue":   "GET /api/v1/stats/overdue (protected)",
					},
				},
			})
		})
	}

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("🛑 Shutting down server...")

		if err := database.CloseDB(); err != nil {
			log.Printf("❌ Error closing database: %v", err)
		}

		log.Println("✅ Server shutdown complete")
		os.Exit(0)
	}()

	// Start server
	serverAddr := ":" + cfg.Server.Port
	log.Println("🎉 Task Management API - Phase 4 Complete!")
	log.Printf("🚀 Server running on http://localhost%s", serverAddr)
	log.Println("📋 Available endpoints:")
	log.Println("   --- Authentication ---")
	log.Println("   POST   http://localhost" + serverAddr + "/api/v1/auth/register")
	log.Println("   POST   http://localhost" + serverAddr + "/api/v1/auth/login")
	log.Println("   GET    http://localhost" + serverAddr + "/api/v1/auth/me (protected)")
	log.Println("   --- Categories (protected) ---")
	log.Println("   GET    http://localhost" + serverAddr + "/api/v1/categories")
	log.Println("   POST   http://localhost" + serverAddr + "/api/v1/categories")
	log.Println("   GET    http://localhost" + serverAddr + "/api/v1/categories/:id")
	log.Println("   PUT    http://localhost" + serverAddr + "/api/v1/categories/:id")
	log.Println("   DELETE http://localhost" + serverAddr + "/api/v1/categories/:id")
	log.Println("   --- Tasks (protected) ---")
	log.Println("   GET    http://localhost" + serverAddr + "/api/v1/tasks")
	log.Println("   POST   http://localhost" + serverAddr + "/api/v1/tasks")
	log.Println("   GET    http://localhost" + serverAddr + "/api/v1/tasks/:id")
	log.Println("   PUT    http://localhost" + serverAddr + "/api/v1/tasks/:id")
	log.Println("   PATCH  http://localhost" + serverAddr + "/api/v1/tasks/:id/status")
	log.Println("   DELETE http://localhost" + serverAddr + "/api/v1/tasks/:id")
	log.Println("   --- Health ---")
	log.Println("   GET    http://localhost" + serverAddr + "/health")
	log.Printf("📚 Swagger: http://localhost%s/swagger/index.html (Phase 5)", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// corsMiddleware handles CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
