package models

import (
	"time"

	"github.com/google/uuid"
)

// DashboardResponse is the main payload for the dashboard
type DashboardResponse struct {
	Statistics     DashboardStatistics   `json:"statistics"`
	RecentProjects []ProjectCard         `json:"recent_projects"`
	UserTasks      []TaskSummary         `json:"user_tasks"`
	NextMeeting    *MeetingWithAttendees `json:"next_meeting"`
}

// DashboardStatistics contains aggregate metrics
type DashboardStatistics struct {
	ActiveProjects    StatValue `json:"active_projects"`
	PendingTasks      StatValue `json:"pending_tasks"`
	TeamMembers       StatValue `json:"team_members"`
	UpcomingDeadlines StatValue `json:"upcoming_deadlines"`
}

// StatValue represents a metric with its change compared to previous period
type StatValue struct {
	Current  int `json:"current"`
	Previous int `json:"previous"`
	Change   int `json:"change"`
}

// ProjectCard represents a project summary for the dashboard
type ProjectCard struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Client       string    `json:"client"`
	Status       string    `json:"status"`
	Progress     int       `json:"progress"`
	DueDate      time.Time `json:"due_date"`
	TeamMembers  []User    `json:"team_members"`
	TotalMembers int       `json:"total_members"`
}

// TaskSummary represents a task summary for the dashboard
type TaskSummary struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	ProjectName   string    `json:"project_name"`
	ProjectID     uuid.UUID `json:"project_id"`
	Priority      int       `json:"priority"`
	PriorityLabel string    `json:"priority_label"`
	DueDate       time.Time `json:"due_date"`
	Status        string    `json:"status"`
}
