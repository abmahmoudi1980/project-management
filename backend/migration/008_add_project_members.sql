-- Migration: 008_add_project_members
-- Feature: Project Member Role Assignment
-- Date: 2026-02-26
-- Description: Add project_roles and project_members tables for role-based project membership

-- ============================================
-- TABLE: project_roles
-- ============================================
CREATE TABLE IF NOT EXISTS project_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Comment on table
COMMENT ON TABLE project_roles IS 'Predefined roles for project membership (owner, manager, contributor, viewer)';

-- ============================================
-- TABLE: project_members
-- ============================================
CREATE TABLE IF NOT EXISTS project_members (
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_role_id UUID NOT NULL REFERENCES project_roles(id),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    added_by UUID REFERENCES users(id),
    PRIMARY KEY (project_id, user_id)
);

-- Comments on table
COMMENT ON TABLE project_members IS 'Many-to-many relationship between projects and users with role assignment';

-- ============================================
-- INDEXES
-- ============================================
CREATE INDEX IF NOT EXISTS idx_project_members_project_id ON project_members(project_id);
CREATE INDEX IF NOT EXISTS idx_project_members_user_id ON project_members(user_id);
CREATE INDEX IF NOT EXISTS idx_project_members_project_role_id ON project_members(project_role_id);
CREATE INDEX IF NOT EXISTS idx_project_roles_is_active ON project_roles(is_active);

-- ============================================
-- SEED DATA: Predefined project roles
-- ============================================
INSERT INTO project_roles (name, display_name, is_active) VALUES
    ('owner', 'Project Owner', true),
    ('manager', 'Project Manager', true),
    ('contributor', 'Contributor', true),
    ('viewer', 'Viewer', true)
ON CONFLICT (name) DO NOTHING;

-- ============================================
-- UPDATED_AT TRIGGER for project_roles
-- ============================================
CREATE OR REPLACE FUNCTION update_project_roles_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_project_roles_updated_at ON project_roles;
CREATE TRIGGER trigger_project_roles_updated_at
    BEFORE UPDATE ON project_roles
    FOR EACH ROW
    EXECUTE FUNCTION update_project_roles_updated_at();
