-- Create task_attachments table
CREATE TABLE IF NOT EXISTS task_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    task_id UUID NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
    original_filename VARCHAR(255) NOT NULL,
    stored_filename VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT NOT NULL CHECK (file_size > 0),
    mime_type VARCHAR(100) NOT NULL,
    uploaded_by UUID REFERENCES users (id) ON DELETE SET NULL,
    has_thumbnail BOOLEAN DEFAULT FALSE,
    thumbnail_path TEXT,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for task_attachments
CREATE INDEX IF NOT EXISTS idx_task_attachments_task_id ON task_attachments (task_id);

CREATE INDEX IF NOT EXISTS idx_task_attachments_uploaded_by ON task_attachments (uploaded_by)
WHERE
    uploaded_by IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_task_attachments_created_at ON task_attachments (created_at DESC);

-- Add trigger to update updated_at timestamp for task_attachments
DROP TRIGGER IF EXISTS update_task_attachments_updated_at ON task_attachments;

CREATE TRIGGER update_task_attachments_updated_at
    BEFORE UPDATE ON task_attachments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();