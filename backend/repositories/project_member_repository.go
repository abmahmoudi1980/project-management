package repositories

import (
	"context"
	"project-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectMemberRepository struct {
	db *pgxpool.Pool
}

func NewProjectMemberRepository(db *pgxpool.Pool) *ProjectMemberRepository {
	return &ProjectMemberRepository{db: db}
}

// AddMember inserts a new project membership
func (r *ProjectMemberRepository) AddMember(ctx context.Context, projectID, userID, roleID, addedBy uuid.UUID) error {
	query := `
		INSERT INTO project_members (project_id, user_id, project_role_id, added_by)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, query, projectID, userID, roleID, addedBy)
	return err
}

// ListMembers retrieves all members of a project with their role details
func (r *ProjectMemberRepository) ListMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMemberResponse, error) {
	query := `
		SELECT 
			pm.project_id,
			pm.user_id,
			pm.project_role_id,
			pm.joined_at,
			pm.added_by,
			u.username,
			u.email,
			pr.name AS role_name,
			pr.display_name AS role_display
		FROM project_members pm
		JOIN users u ON pm.user_id = u.id
		JOIN project_roles pr ON pm.project_role_id = pr.id
		WHERE pm.project_id = $1
		ORDER BY pm.joined_at DESC
	`

	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.ProjectMemberResponse
	for rows.Next() {
		var m models.ProjectMemberResponse
		var addedBy *uuid.UUID
		err := rows.Scan(
			&m.ProjectID,
			&m.UserID,
			&m.ProjectRoleID,
			&m.JoinedAt,
			&addedBy,
			&m.Username,
			&m.Email,
			&m.RoleName,
			&m.RoleDisplay,
		)
		if err != nil {
			return nil, err
		}
		m.AddedBy = addedBy
		members = append(members, m)
	}

	return members, nil
}

// GetMemberByProjectAndUser retrieves a specific membership
func (r *ProjectMemberRepository) GetMemberByProjectAndUser(ctx context.Context, projectID, userID uuid.UUID) (*models.ProjectMember, error) {
	query := `
		SELECT project_id, user_id, project_role_id, joined_at, added_by
		FROM project_members
		WHERE project_id = $1 AND user_id = $2
	`

	var m models.ProjectMember
	err := r.db.QueryRow(ctx, query, projectID, userID).Scan(
		&m.ProjectID,
		&m.UserID,
		&m.ProjectRoleID,
		&m.JoinedAt,
		&m.AddedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// ListEligibleUsers retrieves active users who are not already members of the project
func (r *ProjectMemberRepository) ListEligibleUsers(ctx context.Context, projectID uuid.UUID) ([]models.EligibleUser, error) {
	query := `
		SELECT u.id, u.username, u.email
		FROM users u
		WHERE u.is_active = true
			AND u.id NOT IN (
				SELECT user_id FROM project_members WHERE project_id = $1
			)
		ORDER BY u.username
	`

	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.EligibleUser
	for rows.Next() {
		var u models.EligibleUser
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// GetActiveRoles retrieves all active predefined roles
func (r *ProjectMemberRepository) GetActiveRoles(ctx context.Context) ([]models.ProjectRole, error) {
	query := `
		SELECT id, name, display_name, is_active, created_at, updated_at
		FROM project_roles
		WHERE is_active = true
		ORDER BY 
			CASE name
				WHEN 'owner' THEN 1
				WHEN 'manager' THEN 2
				WHEN 'contributor' THEN 3
				WHEN 'viewer' THEN 4
				ELSE 5
			END
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.ProjectRole
	for rows.Next() {
		var r models.ProjectRole
		if err := rows.Scan(&r.ID, &r.Name, &r.DisplayName, &r.IsActive, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}

	return roles, nil
}

// GetRoleByID retrieves a role by its ID
func (r *ProjectMemberRepository) GetRoleByID(ctx context.Context, roleID uuid.UUID) (*models.ProjectRole, error) {
	query := `
		SELECT id, name, display_name, is_active, created_at, updated_at
		FROM project_roles
		WHERE id = $1 AND is_active = true
	`

	var role models.ProjectRole
	err := r.db.QueryRow(ctx, query, roleID).Scan(
		&role.ID, &role.Name, &role.DisplayName, &role.IsActive, &role.CreatedAt, &role.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// IsUserProjectMember checks if a user is already a member of a project
func (r *ProjectMemberRepository) IsUserProjectMember(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM project_members 
			WHERE project_id = $1 AND user_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, projectID, userID).Scan(&exists)
	return exists, err
}
