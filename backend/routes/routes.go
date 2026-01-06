package routes

import (
	"project-management/handlers"
	"project-management/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func SetupRoutes(
	app *fiber.App,
	projectHandler *handlers.ProjectHandler,
	taskHandler *handlers.TaskHandler,
	timeLogHandler *handlers.TimeLogHandler,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	commentHandler *handlers.CommentHandler,
	dashboardHandler *handlers.DashboardHandler,
	meetingHandler *handlers.MeetingHandler,
	attachmentHandler *handlers.AttachmentHandler,
) {
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	})

	api := app.Group("/api")

	// Public auth routes (no authentication required)
	auth := api.Group("/auth", limiter.New(limiter.Config{
		Max: 10, // 10 requests per minute per IP
	}))
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", limiter.New(limiter.Config{
		Max:        5,      // 5 attempts per 5 minutes per IP for login
		Expiration: 5 * 60, // 5 minutes
	}), authHandler.Login)
	auth.Post("/forgot-password", authHandler.ForgotPassword)
	auth.Post("/reset-password", authHandler.ResetPassword)

	// Protected auth routes (require authentication)
	auth.Get("/me", middleware.RequireAuth, authHandler.GetCurrentUser)
	auth.Post("/logout", middleware.RequireAuth, authHandler.Logout)
	auth.Put("/me", middleware.RequireAuth, authHandler.UpdateProfile)
	auth.Put("/me/password", middleware.RequireAuth, authHandler.ChangePassword)

	// Protected project routes
	projects := api.Group("/projects", middleware.RequireAuth)
	projects.Get("/", projectHandler.GetAllProjects)
	projects.Post("/", projectHandler.CreateProject)
	projects.Get("/:id", projectHandler.GetProject)
	projects.Put("/:id", projectHandler.UpdateProject)
	projects.Delete("/:id", projectHandler.DeleteProject)

	projects.Get("/:projectId/tasks", taskHandler.GetTasksByProject)
	projects.Post("/:projectId/tasks", taskHandler.CreateTask)

	// Protected task routes
	tasks := api.Group("/tasks", middleware.RequireAuth)
	tasks.Get("/:id", taskHandler.GetTask)
	tasks.Put("/:id", taskHandler.UpdateTask)
	tasks.Patch("/:id/complete", taskHandler.ToggleTaskCompletion)
	tasks.Delete("/:id", taskHandler.DeleteTask)

	tasks.Get("/:taskId/timelogs", timeLogHandler.GetTimeLogsByTask)
	tasks.Post("/:taskId/timelogs", timeLogHandler.CreateTimeLog)

	tasks.Get("/:taskId/comments", commentHandler.GetCommentsByTask)
	tasks.Post("/:taskId/comments", commentHandler.CreateComment)

	// Task attachment routes
	tasks.Get("/:taskId/attachments", attachmentHandler.ListAttachments)
	tasks.Post("/:taskId/attachments", attachmentHandler.UploadAttachments)

	// Protected timelog routes
	timelogs := api.Group("/timelogs", middleware.RequireAuth)
	timelogs.Get("/:id", timeLogHandler.GetTimeLog)
	timelogs.Delete("/:id", timeLogHandler.DeleteTimeLog)

	comments := api.Group("/comments", middleware.RequireAuth)
	comments.Put("/:id", commentHandler.UpdateComment)
	comments.Delete("/:id", commentHandler.DeleteComment)

	// Attachment routes
	attachments := api.Group("/attachments", middleware.RequireAuth)
	attachments.Get("/:id/download", attachmentHandler.DownloadAttachment)
	attachments.Get("/:id/thumbnail", attachmentHandler.GetThumbnail)
	attachments.Delete("/:id", attachmentHandler.DeleteAttachment)

	// Admin user management routes (admin only)
	users := api.Group("/users", middleware.RequireAuth, middleware.RequireRole("admin"))
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUserByID)
	users.Put("/:id/role", userHandler.UpdateUserRole)
	users.Put("/:id/activate", userHandler.UpdateUserActivation)

	// Dashboard route
	api.Get("/dashboard", middleware.RequireAuth, dashboardHandler.GetDashboard)

	// Meeting routes
	meetings := api.Group("/meetings", middleware.RequireAuth)
	meetings.Get("/", meetingHandler.ListMeetings)
	meetings.Post("/", meetingHandler.CreateMeeting)
	meetings.Get("/next", meetingHandler.GetNextMeeting)
	meetings.Get("/:id", meetingHandler.GetMeeting)
}
