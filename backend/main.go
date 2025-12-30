package main

import (
	"log"

	"project-management/config"
	"project-management/handlers"
	"project-management/repositories"
	"project-management/routes"
	"project-management/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDB()

	app := fiber.New(fiber.Config{
		AppName: "Project Management API",
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}

		return c.Next()
	})

	projectRepo := repositories.NewProjectRepository(config.DB)
	taskRepo := repositories.NewTaskRepository(config.DB)
	timeLogRepo := repositories.NewTimeLogRepository(config.DB)

	projectService := services.NewProjectService(projectRepo)
	taskService := services.NewTaskService(taskRepo, projectRepo)
	timeLogService := services.NewTimeLogService(timeLogRepo, taskRepo)

	projectHandler := handlers.NewProjectHandler(projectService)
	taskHandler := handlers.NewTaskHandler(taskService)
	timeLogHandler := handlers.NewTimeLogHandler(timeLogService)

	routes.SetupRoutes(app, projectHandler, taskHandler, timeLogHandler)

	log.Println("Server starting on port 3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
