-- Ensure tasks.priority check constraint accepts English priority values.
-- Keep Persian values temporarily for backward compatibility with existing rows.

ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_priority_check;

ALTER TABLE tasks
ADD CONSTRAINT tasks_priority_check
CHECK (
    priority IN (
        -- English
        'Low', 'Medium', 'High',
        -- Persian (legacy / existing data)
        'پایین', 'متوسط', 'بالا'
    )
);
