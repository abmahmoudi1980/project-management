-- Migration: 002_add_user_authentication.sql
-- Feature: User Authentication System
-- Date: 2025-12-30

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP
);

-- Create index on email for faster login queries
CREATE INDEX idx_users_email ON users(email);

-- Create index on role for admin queries
CREATE INDEX idx_users_role ON users(role);

-- Create sessions table (for refresh tokens)
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR(255) NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT false
);

-- Create index on user_id for faster session lookups
CREATE INDEX idx_sessions_user_id ON sessions(user_id);

-- Create index on refresh_token_hash for validation
CREATE INDEX idx_sessions_refresh_token_hash ON sessions(refresh_token_hash);

-- Create password_reset_tokens table
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN NOT NULL DEFAULT false
);

-- Create index on token_hash for fast token validation
CREATE INDEX idx_password_reset_tokens_token_hash ON password_reset_tokens(token_hash);

-- Create index on user_id for cleanup queries
CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);

-- Add user_id (owner) to projects table
ALTER TABLE projects ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id) ON DELETE SET NULL;

-- Add created_by to projects table
ALTER TABLE projects ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL;

-- Add created_by to tasks table
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL;

-- Insert seed admin user
-- Email: admin@example.com
-- Password: Admin123!
-- Password hash generated with bcrypt cost 10
INSERT INTO users (username, email, password_hash, role, is_active)
VALUES (
    'admin',
    'admin@example.com',
    '$2a$10$ecwWYilm18sAI94mxQp7Jupe9JWtAovLfY/BKA.baJpPQE1tIoTi2',
    'admin',
    true
)
ON CONFLICT (email) DO NOTHING;

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for users table
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
