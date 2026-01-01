package repositories

import (
	"context"
	"project-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// GetStatistics retrieves dashboard statistics
func (r *DashboardRepository) GetStatistics(ctx context.Context, userID uuid.UUID, userRole string) (*models.DashboardStatistics, error) {
	stats := &models.DashboardStatistics{}

	// Get active projects count (current and 7 days ago)
	var currentProjects, previousProjects int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM projects 
		 WHERE status IN ('active', 'in_progress') 
		 AND updated_at >= CURRENT_DATE`,
	).Scan(&currentProjects)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM projects 
		 WHERE status IN ('active', 'in_progress') 
		 AND updated_at >= CURRENT_DATE - INTERVAL '7 days'
		 AND updated_at < CURRENT_DATE`,
	).Scan(&previousProjects)
	if err != nil {
		return nil, err
	}

	stats.ActiveProjects = models.StatValue{
		Current:  currentProjects,
		Previous: previousProjects,
		Change:   currentProjects - previousProjects,
	}

	// Get pending tasks count (current and 7 days ago)
	var currentTasks, previousTasks int
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM tasks 
		 WHERE completed = false 
		 AND created_at >= CURRENT_DATE`,
	).Scan(&currentTasks)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM tasks 
		 WHERE completed = false 
		 AND created_at >= CURRENT_DATE - INTERVAL '7 days'
		 AND created_at < CURRENT_DATE`,
	).Scan(&previousTasks)
	if err != nil {
		return nil, err
	}

	stats.PendingTasks = models.StatValue{
		Current:  currentTasks,
		Previous: previousTasks,
		Change:   currentTasks - previousTasks,
	}

	// Get active team members count
	var currentMembers, previousMembers int
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM users 
		 WHERE is_active = true 
		 AND created_at >= CURRENT_DATE`,
	).Scan(&currentMembers)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM users 
		 WHERE is_active = true 
		 AND created_at >= CURRENT_DATE - INTERVAL '7 days'
		 AND created_at < CURRENT_DATE`,
	).Scan(&previousMembers)
	if err != nil {
		return nil, err
	}

	stats.TeamMembers = models.StatValue{
		Current:  currentMembers,
		Previous: previousMembers,
		Change:   currentMembers - previousMembers,
	}

	// Get upcoming deadlines (next 7 days)
	var currentDeadlines, previousDeadlines int
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM tasks 
		 WHERE due_date IS NOT NULL
		 AND due_date >= CURRENT_DATE 
		 AND due_date <= CURRENT_DATE + INTERVAL '7 days'`,
	).Scan(&currentDeadlines)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM tasks 
		 WHERE due_date IS NOT NULL
		 AND due_date >= CURRENT_DATE - INTERVAL '7 days'
		 AND due_date < CURRENT_DATE`,
	).Scan(&previousDeadlines)
	if err != nil {
		return nil, err
	}

	stats.UpcomingDeadlines = models.StatValue{
		Current:  currentDeadlines,
		Previous: previousDeadlines,
		Change:   currentDeadlines - previousDeadlines,
	}

	return stats, nil
}

// GetRecentProjects retrieves the most recently updated projects
func (r *DashboardRepository) GetRecentProjects(ctx context.Context, userID uuid.UUID, userRole string, limit int) ([]models.ProjectCard, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, title, description, updated_at
		 FROM projects
		 WHERE status IN ('active', 'in_progress')
		 ORDER BY updated_at DESC
		 LIMIT $1`,
		limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.ProjectCard
	for rows.Next() {
		var p models.ProjectCard
		var desc *string
		if err := rows.Scan(&p.ID, &p.Name, &desc, &p.UpdatedAt); err != nil {
			return nil, err
		}

		p.Client = p.Name // Default to project name
		if desc != nil && *desc != "" {
			p.Client = *desc
		}
		p.Status = "In Progress"
		p.Progress = 0
		p.TeamMembers = make([]models.UserInfo, 0)

		projects = append(projects, p)
	}

	return projects, nil
}

// GetUserTasks retrieves the user's top pending tasks
func (r *DashboardRepository) GetUserTasks(ctx context.Context, userID uuid.UUID, limit int) ([]models.TaskSummary, error) {
	rows, err := r.db.Query(ctx,
		`SELECT t.id, t.title, p.title, t.priority, t.due_date, t.completed, t.created_at
		 FROM tasks t
		 JOIN projects p ON t.project_id = p.id
		 WHERE t.assignee_id = $1 AND t.completed = false
		 ORDER BY t.priority DESC, t.due_date ASC, t.created_at ASC
		 LIMIT $2`,
		userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.TaskSummary
	for rows.Next() {
		var t models.TaskSummary
		var completed bool
		if err := rows.Scan(&t.ID, &t.Title, &t.ProjectName, &t.Priority, &t.DueDate, &completed, &t.CreatedAt); err != nil {
			return nil, err
		}

		// Map priority to label
		switch t.Priority {
		case 4:
			t.PriorityLabel = "Critical"
		case 3:
			t.PriorityLabel = "High"
		case 2:
			t.PriorityLabel = "Medium"
		case 1:
			t.PriorityLabel = "Low"
		default:
			t.PriorityLabel = "Medium"
		}

		t.Status = "to_do"
		if completed {
			t.Status = "done"
		}

		t.ProjectID = uuid.Nil // Will be populated properly in service

		tasks = append(tasks, t)
	}

	return tasks, nil
}
