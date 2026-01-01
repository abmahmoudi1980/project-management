package repositories

import (
	"context"
	"project-management/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetStatistics(ctx context.Context, userID uuid.UUID, userRole string) (models.DashboardStatistics, error) {
	var stats models.DashboardStatistics

	// 1. Active Projects
	// Current
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM projects 
		WHERE status IN ('active', 'In Progress', 'On Track')
		AND (user_id = $1 OR created_by = $1 OR is_public = true OR $2 = 'admin')
	`, userID, userRole).Scan(&stats.ActiveProjects.Current)
	if err != nil {
		return stats, err
	}

	// Previous (7 days ago)
	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM projects 
		WHERE status IN ('active', 'In Progress', 'On Track')
		AND (user_id = $1 OR created_by = $1 OR is_public = true OR $2 = 'admin')
		AND (updated_at <= NOW() - INTERVAL '7 days' OR created_at <= NOW() - INTERVAL '7 days')
	`, userID, userRole).Scan(&stats.ActiveProjects.Previous)
	if err != nil {
		return stats, err
	}
	stats.ActiveProjects.Change = stats.ActiveProjects.Current - stats.ActiveProjects.Previous

	// 2. Pending Tasks
	// Current
	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM tasks 
		WHERE completed = false
		AND (assignee_id = $1 OR created_by = $1 OR $2 = 'admin')
	`, userID, userRole).Scan(&stats.PendingTasks.Current)
	if err != nil {
		return stats, err
	}

	// Previous (7 days ago)
	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM tasks 
		WHERE completed = false
		AND (assignee_id = $1 OR created_by = $1 OR $2 = 'admin')
		AND (updated_at <= NOW() - INTERVAL '7 days' OR created_at <= NOW() - INTERVAL '7 days')
	`, userID, userRole).Scan(&stats.PendingTasks.Previous)
	if err != nil {
		return stats, err
	}
	stats.PendingTasks.Change = stats.PendingTasks.Current - stats.PendingTasks.Previous

	// 3. Team Members (Admin only)
	if userRole == "admin" {
		err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE is_active = true`).Scan(&stats.TeamMembers.Current)
		if err != nil {
			return stats, err
		}
		err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE is_active = true AND created_at <= NOW() - INTERVAL '7 days'`).Scan(&stats.TeamMembers.Previous)
		if err != nil {
			return stats, err
		}
		stats.TeamMembers.Change = stats.TeamMembers.Current - stats.TeamMembers.Previous
	}

	// 4. Upcoming Deadlines (next 7 days)
	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM (
			SELECT id FROM tasks 
			WHERE completed = false AND due_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '7 days'
			AND (assignee_id = $1 OR created_by = $1 OR $2 = 'admin')
			UNION ALL
			SELECT id FROM projects 
			WHERE status != 'Completed' AND due_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '7 days'
			AND (user_id = $1 OR created_by = $1 OR is_public = true OR $2 = 'admin')
		) AS deadlines
	`, userID, userRole).Scan(&stats.UpcomingDeadlines.Current)
	if err != nil {
		return stats, err
	}

	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM (
			SELECT id FROM tasks 
			WHERE completed = false AND due_date BETWEEN CURRENT_DATE - INTERVAL '7 days' AND CURRENT_DATE
			AND (assignee_id = $1 OR created_by = $1 OR $2 = 'admin')
			UNION ALL
			SELECT id FROM projects 
			WHERE status != 'Completed' AND due_date BETWEEN CURRENT_DATE - INTERVAL '7 days' AND CURRENT_DATE
			AND (user_id = $1 OR created_by = $1 OR is_public = true OR $2 = 'admin')
		) AS deadlines
	`, userID, userRole).Scan(&stats.UpcomingDeadlines.Previous)
	if err != nil {
		return stats, err
	}
	stats.UpcomingDeadlines.Change = stats.UpcomingDeadlines.Current - stats.UpcomingDeadlines.Previous

	return stats, nil
}

func (r *DashboardRepository) GetRecentProjects(ctx context.Context, userID uuid.UUID, userRole string, limit int) ([]models.ProjectCard, error) {
	rows, err := r.db.Query(ctx, `
		SELECT 
			p.id, p.title, p.status, p.updated_at,
			COALESCE(p.due_date, (p.created_at + INTERVAL '30 days')::date) as due_date,
			COUNT(t.id) as total_tasks,
			COUNT(CASE WHEN t.completed = true THEN 1 END) as completed_tasks
		FROM projects p
		LEFT JOIN tasks t ON t.project_id = p.id
		WHERE (p.user_id = $1 OR p.created_by = $1 OR p.is_public = true OR $2 = 'admin')
		AND p.status IN ('active', 'Planning', 'In Progress', 'On Track', 'Review')
		GROUP BY p.id, p.title, p.status, p.updated_at, p.due_date, p.created_at
		ORDER BY p.updated_at DESC
		LIMIT $3
	`, userID, userRole, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := []models.ProjectCard{}
	for rows.Next() {
		var p models.ProjectCard
		var totalTasks, completedTasks int
		var updatedAt time.Time
		var dueDate *time.Time
		if err := rows.Scan(&p.ID, &p.Name, &p.Status, &updatedAt, &dueDate, &totalTasks, &completedTasks); err != nil {
			return nil, err
		}
		p.Client = p.Name // Default to project name as per spec
		if dueDate != nil {
			p.DueDate = *dueDate
		}
		if totalTasks > 0 {
			p.Progress = (completedTasks * 100) / totalTasks
		} else {
			p.Progress = 0
		}

		// Get team members (max 3)
		p.TeamMembers, p.TotalMembers, err = r.getProjectTeamMembers(ctx, p.ID)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, nil
}

func (r *DashboardRepository) getProjectTeamMembers(ctx context.Context, projectID uuid.UUID) ([]models.User, int, error) {
	// Total count
	var total int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		JOIN tasks t ON t.assignee_id = u.id
		WHERE t.project_id = $1
	`, projectID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Top 3
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT u.id, u.username, u.email, u.role, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN tasks t ON t.assignee_id = u.id
		WHERE t.project_id = $1
		LIMIT 3
	`, projectID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var members []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		members = append(members, u)
	}

	return members, total, nil
}

func (r *DashboardRepository) GetUserTasks(ctx context.Context, userID uuid.UUID, limit int) ([]models.TaskSummary, error) {
	rows, err := r.db.Query(ctx, `
		SELECT t.id, t.title, p.title as project_name, t.project_id, t.priority, COALESCE(t.due_date, (t.created_at + INTERVAL '7 days')::date), t.completed
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE (t.assignee_id = $1 OR t.created_by = $1)
		AND t.completed = false
		ORDER BY 
			CASE 
				WHEN t.priority = 'Critical' THEN 4
				WHEN t.priority = 'High' THEN 3
				WHEN t.priority = 'Medium' THEN 2
				WHEN t.priority = 'Low' THEN 1
				ELSE 0
			END DESC,
			t.due_date ASC NULLS LAST
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.TaskSummary{}
	for rows.Next() {
		var t models.TaskSummary
		var priorityStr string
		var completed bool
		var dueDate *time.Time
		if err := rows.Scan(&t.ID, &t.Title, &t.ProjectName, &t.ProjectID, &priorityStr, &dueDate, &completed); err != nil {
			return nil, err
		}
		if dueDate != nil {
			t.DueDate = *dueDate
		}
		t.PriorityLabel = priorityStr
		switch priorityStr {
		case "Critical":
			t.Priority = 4
		case "High":
			t.Priority = 3
		case "Medium":
			t.Priority = 2
		case "Low":
			t.Priority = 1
		default:
			t.Priority = 0
		}
		if completed {
			t.Status = "done"
		} else {
			t.Status = "todo"
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
