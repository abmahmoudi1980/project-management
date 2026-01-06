package main

import (
	"log"

	"project-management/config"
	"project-management/handlers"
	"project-management/repositories"
	"project-management/routes"
	"project-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDB()

	if err := config.InitFileStorage(); err != nil {
		log.Fatalf("Failed to initialize file storage: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: "Project Management API",
	})

	// CORS middleware with credentials support
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization",
	}))

	// Security headers
	app.Use(helmet.New())

	// Request logging (helps verify frontend API calls)
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Initialize repositories
	projectRepo := repositories.NewProjectRepository(config.DB)
	taskRepo := repositories.NewTaskRepository(config.DB)
	timeLogRepo := repositories.NewTimeLogRepository(config.DB)
	userRepo := repositories.NewUserRepository(config.DB)
	sessionRepo := repositories.NewSessionRepository(config.DB)
	passwordResetRepo := repositories.NewPasswordResetRepository(config.DB)
	commentRepo := repositories.NewCommentRepository(config.DB)
	dashboardRepo := repositories.NewDashboardRepository(config.DB)
	meetingRepo := repositories.NewMeetingRepository(config.DB)

	// Initialize services
	emailService := services.NewEmailService()
	projectService := services.NewProjectService(projectRepo)
	taskService := services.NewTaskService(taskRepo, projectRepo)
	timeLogService := services.NewTimeLogService(timeLogRepo, taskRepo)
	authService := services.NewAuthService(userRepo, sessionRepo, passwordResetRepo, emailService)
	userService := services.NewUserService(userRepo)
	commentService := services.NewCommentService(commentRepo, taskRepo)
	dashboardService := services.NewDashboardService(dashboardRepo, meetingRepo)
	meetingService := services.NewMeetingService(meetingRepo, userRepo)

	// Initialize handlers
	projectHandler := handlers.NewProjectHandler(projectService)
	taskHandler := handlers.NewTaskHandler(taskService)
	timeLogHandler := handlers.NewTimeLogHandler(timeLogService)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	commentHandler := handlers.NewCommentHandler(commentService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	meetingHandler := handlers.NewMeetingHandler(meetingService)

	routes.SetupRoutes(app, projectHandler, taskHandler, timeLogHandler, authHandler, userHandler, commentHandler, dashboardHandler, meetingHandler)

	log.Println("Server starting on port 3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
