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

func (s *DashboardService) GetDashboardData(ctx context.Context, userID uuid.UUID, userRole string) (*models.DashboardResponse, error) {
	var resp models.DashboardResponse
	var err error

	// 1. Get Statistics
	resp.Statistics, err = s.dashboardRepo.GetStatistics(ctx, userID, userRole)
	if err != nil {
		return nil, err
	}

	// 2. Get Recent Projects (limit 4)
	resp.RecentProjects, err = s.dashboardRepo.GetRecentProjects(ctx, userID, userRole, 4)
	if err != nil {
		return nil, err
	}

	// 3. Get User Tasks (limit 5)
	resp.UserTasks, err = s.dashboardRepo.GetUserTasks(ctx, userID, 5)
	if err != nil {
		return nil, err
	}

	// 4. Get Next Meeting
	resp.NextMeeting, err = s.meetingRepo.GetNextMeetingForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
