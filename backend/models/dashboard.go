package models

import (
	"time"

	"github.com/google/uuid"
)

// StatValue represents a statistic with current, previous, and change values
type StatValue struct {
	Current  int `json:"current"`
	Previous int `json:"previous"`
	Change   int `json:"change"`
}

// DashboardStatistics contains all dashboard statistics
type DashboardStatistics struct {
	ActiveProjects    StatValue `json:"active_projects"`
	PendingTasks      StatValue `json:"pending_tasks"`
	TeamMembers       StatValue `json:"team_members"`
	UpcomingDeadlines StatValue `json:"upcoming_deadlines"`
}

// ProjectCard represents a project in the recent projects list
type ProjectCard struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Client       string     `json:"client"`
	Status       string     `json:"status"`
	Progress     int        `json:"progress"` // 0-100
	DueDate      *time.Time `json:"due_date,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	TeamMembers  []UserInfo `json:"team_members"`
	TotalMembers int        `json:"total_members"`
}

// TaskSummary represents a task in the user's task list
type TaskSummary struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	ProjectName   string     `json:"project_name"`
	ProjectID     uuid.UUID  `json:"project_id"`
	Priority      int        `json:"priority"` // 1=Low, 2=Medium, 3=High, 4=Critical
	PriorityLabel string     `json:"priority_label"`
	DueDate       *time.Time `json:"due_date,omitempty"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
}

// DashboardResponse is the complete dashboard data response
type DashboardResponse struct {
	Statistics     DashboardStatistics   `json:"statistics"`
	RecentProjects []ProjectCard         `json:"recent_projects"`
	UserTasks      []TaskSummary         `json:"user_tasks"`
	NextMeeting    *MeetingWithAttendees `json:"next_meeting"`
}
