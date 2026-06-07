-- Add parent_id to projects for the project hierarchy feature (001-project-hierarchy).
-- See specs/001-project-hierarchy/spec.md (FR-001..FR-010) and research.md (R1, R6).

ALTER TABLE projects
    ADD COLUMN IF NOT EXISTS parent_id UUID;

ALTER TABLE projects
    ADD CONSTRAINT fk_projects_parent
    FOREIGN KEY (parent_id) REFERENCES projects(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_projects_parent_id ON projects(parent_id);

COMMENT ON COLUMN projects.parent_id IS
    'Optional parent project. NULL means top-level. ON DELETE SET NULL keeps orphans accessible.';
