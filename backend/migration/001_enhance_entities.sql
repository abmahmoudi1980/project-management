-- Phase 1: Setup (Database Migration)
-- Add new columns to existing tables without data loss

-- Projects table enhancements
ALTER TABLE projects ADD COLUMN
IF NOT EXISTS identifier VARCHAR
(50) UNIQUE;
ALTER TABLE projects ADD COLUMN
IF NOT EXISTS homepage VARCHAR
(255);
ALTER TABLE projects ADD COLUMN
IF NOT EXISTS is_public BOOLEAN DEFAULT FALSE;

-- Tasks table enhancements
ALTER TABLE tasks ADD COLUMN
IF NOT EXISTS description TEXT;
ALTER TABLE tasks ADD COLUMN
IF NOT EXISTS assignee_id UUID;
ALTER TABLE tasks ADD COLUMN
IF NOT EXISTS author_id UUID;
ALTER TABLE tasks ADD COLUMN
IF NOT EXISTS category VARCHAR
(100);
ALTER TABLE tasks ADD COLUMN
IF NOT EXISTS start_date DATE;
ALTER TABLE tasks ADD COLUMN
IF NOT EXISTS due_date DATE;
ALTER TABLE tasks ADD COLUMN
IF NOT EXISTS estimated_hours DECIMAL
(10,2);
ALTER TABLE tasks ADD COLUMN
IF NOT EXISTS done_ratio INTEGER DEFAULT 0;

-- Time logs table enhancements
ALTER TABLE time_logs ADD COLUMN
IF NOT EXISTS user_id UUID;
ALTER TABLE time_logs ADD COLUMN
IF NOT EXISTS activity_type VARCHAR
(50);
ALTER TABLE time_logs ADD COLUMN
IF NOT EXISTS project_id UUID;

-- Add indexes for performance
CREATE INDEX
IF NOT EXISTS idx_projects_identifier ON projects
(identifier);
CREATE INDEX
IF NOT EXISTS idx_tasks_assignee_id ON tasks
(assignee_id);
CREATE INDEX
IF NOT EXISTS idx_tasks_due_date ON tasks
(due_date);
CREATE INDEX
IF NOT EXISTS idx_time_logs_user_id ON time_logs
(user_id);
CREATE INDEX
IF NOT EXISTS idx_time_logs_activity_type ON time_logs
(activity_type);

-- Add constraints
-- done_ratio check 0-100
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS check_done_ratio;
ALTER TABLE tasks ADD CONSTRAINT check_done_ratio CHECK (done_ratio >= 0 AND done_ratio <= 100);

-- Data migration for existing projects (T093)
-- Auto-generate identifiers for existing projects if they are null
UPDATE projects 
SET identifier = LOWER(REGEXP_REPLACE(title, '[^a-zA-Z0-9]+', '-', 'g'))
WHERE identifier IS NULL;

-- Ensure identifiers are unique if there were duplicates after generation (unlikely but safe)
-- This is a simple approach, might need more complex logic if many projects have same title
UPDATE projects p1
SET identifier
= identifier || '-' || SUBSTRING
(id::text, 1, 4)
WHERE EXISTS
(
    SELECT 1
FROM projects p2
WHERE p1.identifier = p2.identifier AND p1.id <> p2.id
);
