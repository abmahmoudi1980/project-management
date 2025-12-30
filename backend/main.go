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

	// Initialize repositories
	projectRepo := repositories.NewProjectRepository(config.DB)
	taskRepo := repositories.NewTaskRepository(config.DB)
	timeLogRepo := repositories.NewTimeLogRepository(config.DB)
	userRepo := repositories.NewUserRepository(config.DB)
	sessionRepo := repositories.NewSessionRepository(config.DB)

	// Initialize services
	projectService := services.NewProjectService(projectRepo)
	taskService := services.NewTaskService(taskRepo, projectRepo)
	timeLogService := services.NewTimeLogService(timeLogRepo, taskRepo)
	authService := services.NewAuthService(userRepo, sessionRepo)

	// Initialize handlers
	projectHandler := handlers.NewProjectHandler(projectService)
	taskHandler := handlers.NewTaskHandler(taskService)
	timeLogHandler := handlers.NewTimeLogHandler(timeLogService)
	authHandler := handlers.NewAuthHandler(authService)

	routes.SetupRoutes(app, projectHandler, taskHandler, timeLogHandler, authHandler)

	log.Println("Server starting on port 3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
