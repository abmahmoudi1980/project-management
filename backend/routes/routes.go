package routes

import (
	"project-management/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, projectHandler *handlers.ProjectHandler, taskHandler *handlers.TaskHandler, timeLogHandler *handlers.TimeLogHandler) {
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	})

	api := app.Group("/api")

	projects := api.Group("/projects")
	projects.Get("/", projectHandler.GetAllProjects)
	projects.Post("/", projectHandler.CreateProject)
	projects.Get("/:id", projectHandler.GetProject)
	projects.Put("/:id", projectHandler.UpdateProject)
	projects.Delete("/:id", projectHandler.DeleteProject)

	projects.Get("/:projectId/tasks", taskHandler.GetTasksByProject)
	projects.Post("/:projectId/tasks", taskHandler.CreateTask)

	tasks := api.Group("/tasks")
	tasks.Get("/:id", taskHandler.GetTask)
	tasks.Put("/:id", taskHandler.UpdateTask)
	tasks.Patch("/:id/complete", taskHandler.ToggleTaskCompletion)
	tasks.Delete("/:id", taskHandler.DeleteTask)

	tasks.Get("/:taskId/timelogs", timeLogHandler.GetTimeLogsByTask)
	tasks.Post("/:taskId/timelogs", timeLogHandler.CreateTimeLog)

	timelogs := api.Group("/timelogs")
	timelogs.Get("/:id", timeLogHandler.GetTimeLog)
	timelogs.Delete("/:id", timeLogHandler.DeleteTimeLog)
}
