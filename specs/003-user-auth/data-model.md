# Data Model: User Authentication

**Feature**: 003-user-auth  
**Date**: 2025-12-30  
**Purpose**: Define data entities and relationships for authentication system

## Entity Overview

This feature introduces 3 new entities and modifies 3 existing entities:

**New Entities**:
1. User - Core user account information
2. Session - User authentication sessions (refresh tokens)
3. PasswordResetToken - Temporary tokens for password recovery

**Modified Entities**:
1. Project - Add user_id (owner) and created_by fields
2. Task - Add created_by field
3. TimeLog - Link to user (already has user_id)

---

## 1. User Entity

### Purpose
Represents a user account in the system with authentication credentials and role-based permissions.

### Attributes

| Attribute | Type | Constraints | Description |
|-----------|------|-------------|-------------|
| id | UUID | PRIMARY KEY, NOT NULL | Unique user identifier |
| username | VARCHAR(50) | UNIQUE, NOT NULL | User's display name |
| email | VARCHAR(255) | UNIQUE, NOT NULL | User's email address (used for login) |
| password_hash | VARCHAR(255) | NOT NULL | bcrypt-hashed password |
| role | VARCHAR(20) | NOT NULL, DEFAULT 'user' | User role: 'admin' or 'user' |
| is_active | BOOLEAN | NOT NULL, DEFAULT true | Account active status |
| failed_login_attempts | INTEGER | NOT NULL, DEFAULT 0 | Counter for account lockout |
| locked_until | TIMESTAMP | NULLABLE | Account unlock time (NULL if not locked) |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | Account creation timestamp |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | Last update timestamp |
| last_login_at | TIMESTAMP | NULLABLE | Last successful login timestamp |

### Validation Rules

- **username**: 
  - Length: 3-50 characters
  - Pattern: Alphanumeric, underscores, hyphens (no spaces)
  - Case-insensitive uniqueness check
  
- **email**:
  - Valid email format (RFC 5322 basic validation)
  - Case-insensitive uniqueness check
  - Normalized to lowercase before storage
  
- **password** (before hashing):
  - Minimum length: 8 characters
  - Must contain: uppercase letter, lowercase letter, digit
  - No maximum length (bcrypt handles truncation at 72 bytes)
  
- **role**:
  - Allowed values: 'admin', 'user'
  - Default: 'user' for new registrations
  
- **is_active**:
  - true: User can log in
  - false: User account is deactivated (by admin)

### Business Rules

1. **Account Lockout**:
   - After 5 failed login attempts, account is locked for 30 minutes
   - `failed_login_attempts` is incremented on each failed login
   - `locked_until` is set to NOW() + 30 minutes
   - On successful login, `failed_login_attempts` resets to 0 and `locked_until` to NULL

2. **Role Assignment**:
   - New users are created with role 'user' by default
   - Only admins can change user roles
   - There must be at least one admin user in the system

3. **Account Deactivation**:
   - Only admins can deactivate users (set `is_active` = false)
   - Deactivated users cannot log in
   - Admins cannot deactivate themselves

4. **Email Verification** (Future enhancement - not in P1):
   - Currently not implemented
   - Can add `email_verified` boolean field later

### Relationships

- **1:N** → Project (user owns many projects)
- **1:N** → Task (user creates many tasks, user is assigned to many tasks)
- **1:N** → TimeLog (user logs many time entries)
- **1:N** → Session (user has many active sessions)
- **1:N** → PasswordResetToken (user requests many reset tokens over time)

### Indexes

```sql
CREATE INDEX idx_users_email ON users(LOWER(email));
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);
```

---

## 2. Session Entity

### Purpose
Stores refresh tokens for JWT-based authentication. Access tokens are stateless (not stored), but refresh tokens need validation against the database.

### Attributes

| Attribute | Type | Constraints | Description |
|-----------|------|-------------|-------------|
| id | UUID | PRIMARY KEY, NOT NULL | Unique session identifier |
| user_id | UUID | FOREIGN KEY (users.id), NOT NULL | Reference to user |
| refresh_token_hash | VARCHAR(255) | UNIQUE, NOT NULL | SHA-256 hash of refresh token |
| user_agent | VARCHAR(500) | NULLABLE | Browser/client information |
| ip_address | VARCHAR(45) | NULLABLE | Client IP address (IPv6 compatible) |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | Session creation time |
| expires_at | TIMESTAMP | NOT NULL | Session expiration time |
| revoked | BOOLEAN | NOT NULL, DEFAULT false | Manual revocation flag |

### Validation Rules

- **refresh_token_hash**:
  - SHA-256 hash (64 hex characters)
  - Unique across all sessions
  
- **expires_at**:
  - Must be > created_at
  - Standard expiry: 7 days from creation

### Business Rules

1. **Session Creation**:
   - Created on successful login
   - `expires_at` = NOW() + 7 days
   - Refresh token is hashed before storage (like passwords)

2. **Session Validation**:
   - Check that `expires_at` > NOW()
   - Check that `revoked` = false
   - Check that user `is_active` = true

3. **Session Cleanup**:
   - Expired sessions should be periodically deleted (cleanup job)
   - Users can manually revoke sessions (logout)

4. **Multiple Sessions**:
   - Users can have multiple active sessions (different devices)
   - Admin UI can show all active sessions per user

### Relationships

- **N:1** → User (many sessions belong to one user)

### Indexes

```sql
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_sessions_refresh_token_hash ON sessions(refresh_token_hash);
```

---

## 3. PasswordResetToken Entity

### Purpose
Stores temporary tokens for password reset flow. Tokens are single-use and time-limited.

### Attributes

| Attribute | Type | Constraints | Description |
|-----------|------|-------------|-------------|
| id | UUID | PRIMARY KEY, NOT NULL | Unique token record identifier |
| user_id | UUID | FOREIGN KEY (users.id), NOT NULL | Reference to user |
| token_hash | VARCHAR(255) | UNIQUE, NOT NULL | SHA-256 hash of reset token |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | Token creation time |
| expires_at | TIMESTAMP | NOT NULL | Token expiration time |
| used | BOOLEAN | NOT NULL, DEFAULT false | Whether token has been used |

### Validation Rules

- **token_hash**:
  - SHA-256 hash (64 hex characters)
  - Unique across all tokens
  
- **expires_at**:
  - Must be > created_at
  - Standard expiry: 1 hour from creation

### Business Rules

1. **Token Generation**:
   - Cryptographically secure random token (32 bytes)
   - Hashed with SHA-256 before storage
   - Plain token sent via email (not stored)
   - `expires_at` = NOW() + 1 hour

2. **Token Validation**:
   - Check that `expires_at` > NOW()
   - Check that `used` = false
   - Check that user `is_active` = true

3. **Token Usage**:
   - After successful password reset, mark `used` = true
   - Tokens cannot be reused

4. **Rate Limiting**:
   - Users can request max 3 reset tokens per hour (enforced in application layer)
   - Old unused tokens for same user remain valid until expiry

5. **Cleanup**:
   - Expired and used tokens should be periodically deleted (cleanup job)

### Relationships

- **N:1** → User (many tokens belong to one user over time)

### Indexes

```sql
CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX idx_password_reset_tokens_token_hash ON password_reset_tokens(token_hash);
CREATE INDEX idx_password_reset_tokens_expires_at ON password_reset_tokens(expires_at);
```

---

## 4. Modified Entity: Project

### New Attributes

| Attribute | Type | Constraints | Description |
|-----------|------|-------------|-------------|
| user_id | UUID | FOREIGN KEY (users.id), NULLABLE | Project owner (creator) |
| created_by | UUID | FOREIGN KEY (users.id), NULLABLE | User who created the project |

### Business Rules

- **Ownership**:
  - `user_id` is set to the user who creates the project
  - Admins can view/edit all projects regardless of owner
  - Regular users can only view/edit their own projects
  
- **Created By**:
  - `created_by` tracks the original creator (for audit purposes)
  - Immutable after creation

- **Migration**:
  - Existing projects will have NULL user_id (created before auth system)
  - Option 1: Assign all existing projects to first admin user
  - Option 2: Keep as NULL and show as "System Projects" accessible to all

### Relationships

- **N:1** → User (many projects belong to one user)

---

## 5. Modified Entity: Task

### New Attributes

| Attribute | Type | Constraints | Description |
|-----------|------|-------------|-------------|
| created_by | UUID | FOREIGN KEY (users.id), NULLABLE | User who created the task |

### Business Rules

- **Creator Tracking**:
  - `created_by` is set to the user who creates the task
  - Immutable after creation
  - Used for audit trails

- **Assignment** (existing field):
  - `assignee_id` already exists in current schema (from feature 001)
  - Links task to assigned user
  - Admin can assign tasks to any user
  - Regular users can assign tasks to themselves

- **Access Control**:
  - Admins can view/edit all tasks
  - Regular users can only view/edit tasks in their own projects
  - Regular users can only view/edit tasks assigned to them

### Relationships

- **N:1** → User (many tasks created by one user)
- **N:1** → User (many tasks assigned to one user) - via `assignee_id`

---

## 6. Modified Entity: TimeLog

### Note
TimeLog entity already has `user_id` field (from feature 001). No schema changes needed, but access control rules apply:

### Business Rules

- **Access Control**:
  - Admins can view/edit all time logs
  - Regular users can only view/edit their own time logs
  - Time logs are linked to tasks, which are linked to projects (ownership cascade)

---

## Entity Relationship Diagram

```
┌──────────────┐
│     User     │
│──────────────│
│ id (PK)      │
│ username     │
│ email        │
│ password_hash│
│ role         │
│ is_active    │
│ failed_...   │
│ locked_until │
│ created_at   │
│ updated_at   │
│ last_login_at│
└──────┬───────┘
       │
       │ 1:N
       │
       ├─────────────┬─────────────┬─────────────┐
       │             │             │             │
       ▼             ▼             ▼             ▼
┌─────────────┐ ┌─────────┐ ┌──────────────┐ ┌──────────────┐
│   Session   │ │ Project │ │     Task     │ │   TimeLog    │
│─────────────│ │─────────│ │──────────────│ │──────────────│
│ id (PK)     │ │ id (PK) │ │ id (PK)      │ │ id (PK)      │
│ user_id (FK)│ │ user_id │ │ created_by   │ │ user_id (FK) │
│ refresh_... │ │ created │ │ assignee_id  │ │ ...          │
│ expires_at  │ │ ...     │ │ ...          │ └──────────────┘
│ revoked     │ └─────────┘ └──────────────┘
└─────────────┘
       │
       │ 1:N
       │
       ▼
┌──────────────────────┐
│ PasswordResetToken   │
│──────────────────────│
│ id (PK)              │
│ user_id (FK)         │
│ token_hash           │
│ expires_at           │
│ used                 │
└──────────────────────┘
```

---

## Database Migration Strategy

### Migration File: `002_add_user_authentication.sql`

**Order of Operations**:
1. Create `users` table
2. Create `sessions` table
3. Create `password_reset_tokens` table
4. Alter `projects` table (add user_id, created_by)
5. Alter `tasks` table (add created_by)
6. Create indexes
7. Create first admin user (seed data)

### Rollback Strategy
- Down migration drops tables in reverse order
- Foreign keys must be dropped before parent tables

### Data Migration for Existing Records
- Existing projects/tasks will have NULL user references
- Option 1: Create a "system" admin user and assign all existing records to it
- Option 2: Keep as NULL and handle in application logic

---

## State Transitions

### User Account States

```
           ┌──────────────────┐
           │  Registration    │
           │  (role: user)    │
           └────────┬─────────┘
                    │
                    ▼
           ┌──────────────────┐
           │   Active         │◄──────┐
           │ (is_active=true) │       │
           └──┬────────────┬──┘       │
              │            │          │
    Failed    │            │ Admin    │ Admin
    Login     │            │ Action   │ Action
              │            │          │
              ▼            ▼          │
     ┌────────────┐  ┌──────────┐    │
     │  Locked    │  │Deactivated│────┘
     │(30 minutes)│  │(is_active │
     │            │  │  =false)  │
     └─────┬──────┘  └───────────┘
           │
           │ Time Expires
           │ or Successful Login
           │
           └──────────────────────────┘
```

### Session States

```
     ┌──────────────┐
     │   Created    │
     │  (on login)  │
     └──────┬───────┘
            │
            ▼
     ┌──────────────┐
     │    Active    │
     │(not expired) │
     └──┬───────┬───┘
        │       │
        │       │ Manual Logout
        │       │ or Admin Revoke
        │       │
        │       ▼
        │  ┌──────────┐
        │  │ Revoked  │
        │  │(revoked= │
        │  │  true)   │
        │  └──────────┘
        │
        │ Time > expires_at
        │
        ▼
     ┌──────────────┐
     │   Expired    │
     │ (cleanup job)│
     └──────────────┘
```

### Password Reset Token States

```
     ┌──────────────┐
     │   Created    │
     │ (on request) │
     └──────┬───────┘
            │
            ▼
     ┌──────────────┐
     │    Valid     │
     │ (not expired,│
     │  not used)   │
     └──┬───────┬───┘
        │       │
        │       │ User Resets Password
        │       │
        │       ▼
        │  ┌──────────┐
        │  │   Used   │
        │  │ (used=   │
        │  │  true)   │
        │  └──────────┘
        │
        │ Time > expires_at
        │
        ▼
     ┌──────────────┐
     │   Expired    │
     │ (cleanup job)│
     └──────────────┘
```

---

## Performance Considerations

### Query Patterns

**Most Frequent Queries**:
1. User lookup by email (every login) → Index on LOWER(email)
2. Session validation by refresh token → Index on refresh_token_hash
3. User role check (every protected request) → Index on role
4. Password reset token lookup → Index on token_hash

**Batch Operations**:
- Cleanup expired sessions (hourly cron job)
- Cleanup expired/used reset tokens (daily cron job)

### Scalability

**Current Scale**: Up to 1000 users (per spec assumptions)

**Optimizations for Larger Scale** (future):
- Add `last_activity_at` to sessions for inactive session cleanup
- Partition sessions table by creation date
- Add Redis cache for user role/status checks
- Add email verification queue (async)

---

## Security Considerations

### Data Protection

1. **Password Storage**:
   - Never store plain text passwords
   - Use bcrypt with cost factor 10
   - Password hashes are 60 characters (bcrypt format)

2. **Token Storage**:
   - Refresh tokens and reset tokens are hashed (SHA-256)
   - Plain tokens only exist in memory/transit, never in database

3. **Sensitive Fields**:
   - `password_hash`: Never expose in API responses
   - `refresh_token_hash`: Never expose in API responses
   - `locked_until`: Can expose to user (for feedback)
   - `failed_login_attempts`: Internal only, don't expose

### Access Patterns

1. **Read Operations**:
   - Users can read their own user record (excluding password_hash)
   - Admins can read all user records (excluding password_hash)

2. **Write Operations**:
   - Users can update their own username, email, password
   - Users cannot change their own role or is_active status
   - Admins can update any user's role or is_active status
   - Admins cannot delete the last admin user

---

## Testing Considerations

### Unit Tests (Model Layer)

- Password hashing and verification
- Token generation and hashing
- Validation rules (email format, password strength)
- State transitions (account lockout, session expiry)

### Integration Tests (Database Layer)

- User creation with duplicate email (should fail)
- Session creation and validation
- Password reset token lifecycle
- Cascading deletes (user → sessions, tokens)

### End-to-End Tests

- Full registration → login → access protected resource
- Failed login attempts → account lockout → unlock
- Password reset request → email → reset → login
- Admin role assignment → access admin-only resources

---

## Appendix: Go Model Structs

### User Model

```go
package models

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID                   uuid.UUID  `json:"id"`
    Username             string     `json:"username"`
    Email                string     `json:"email"`
    PasswordHash         string     `json:"-"` // Never expose in JSON
    Role                 string     `json:"role"`
    IsActive             bool       `json:"is_active"`
    FailedLoginAttempts  int        `json:"-"` // Internal only
    LockedUntil          *time.Time `json:"locked_until,omitempty"`
    CreatedAt            time.Time  `json:"created_at"`
    UpdatedAt            time.Time  `json:"updated_at"`
    LastLoginAt          *time.Time `json:"last_login_at,omitempty"`
}

type CreateUserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"` // Plain text, will be hashed
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    User         User   `json:"user"`
    AccessToken  string `json:"access_token"`  // Short-lived JWT
    RefreshToken string `json:"refresh_token"` // Long-lived JWT
}
```

### Session Model

```go
type Session struct {
    ID                uuid.UUID  `json:"id"`
    UserID            uuid.UUID  `json:"user_id"`
    RefreshTokenHash  string     `json:"-"` // Never expose
    UserAgent         *string    `json:"user_agent,omitempty"`
    IPAddress         *string    `json:"ip_address,omitempty"`
    CreatedAt         time.Time  `json:"created_at"`
    ExpiresAt         time.Time  `json:"expires_at"`
    Revoked           bool       `json:"revoked"`
}
```

### PasswordResetToken Model

```go
type PasswordResetToken struct {
    ID        uuid.UUID  `json:"id"`
    UserID    uuid.UUID  `json:"user_id"`
    TokenHash string     `json:"-"` // Never expose
    CreatedAt time.Time  `json:"created_at"`
    ExpiresAt time.Time  `json:"expires_at"`
    Used      bool       `json:"used"`
}

type PasswordResetRequest struct {
    Email string `json:"email"`
}

type PasswordResetConfirm struct {
    Token       string `json:"token"`
    NewPassword string `json:"new_password"`
}
```

---

**Document Version**: 1.0  
**Last Updated**: 2025-12-30
