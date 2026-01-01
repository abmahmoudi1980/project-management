package services

import (
	"context"
	"project-management/models"
	"project-management/repositories"

	"github.com/google/uuid"
)

type DashboardService struct {
	dashboardRepo *repositories.DashboardRepository
	meetingRepo   *repositories.MeetingRepository
}

func NewDashboardService(dashboardRepo *repositories.DashboardRepository, meetingRepo *repositories.MeetingRepository) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
		meetingRepo:   meetingRepo,
	}
}

// GetDashboardData aggregates all dashboard information
func (s *DashboardService) GetDashboardData(ctx context.Context, userID uuid.UUID, userRole string) (*models.DashboardResponse, error) {
	response := &models.DashboardResponse{
		RecentProjects: make([]models.ProjectCard, 0),
		UserTasks:      make([]models.TaskSummary, 0),
	}

	// Get statistics
	stats, err := s.dashboardRepo.GetStatistics(ctx, userID, userRole)
	if err != nil {
		return nil, err
	}
	response.Statistics = *stats

	// Get recent projects (up to 4)
	projects, err := s.dashboardRepo.GetRecentProjects(ctx, userID, userRole, 4)
	if err != nil {
		return nil, err
	}
	response.RecentProjects = projects

	// Get user's tasks (up to 5)
	tasks, err := s.dashboardRepo.GetUserTasks(ctx, userID, 5)
	if err != nil {
		return nil, err
	}
	response.UserTasks = tasks

	// Get next meeting
	meeting, err := s.meetingRepo.GetNextMeetingForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	response.NextMeeting = meeting

	return response, nil
}
