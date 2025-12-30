# Implementation Plan: User Authentication (احراز هویت کاربر)

**Branch**: `003-user-auth` | **Date**: 2025-12-30 | **Spec**: [spec.md](spec.md)  
**Input**: Feature specification from `/specs/003-user-auth/spec.md`

## Summary

Implement a complete user authentication system for the project management application with:
- User registration and login with JWT tokens stored in httpOnly cookies
- Two user roles (Admin and Regular User) with role-based access control
- Password reset flow via email
- Account security features (password hashing with bcrypt, account lockout after failed attempts)
- Admin user management interface
- All UI text in Persian language

Technical approach: JWT-based stateless authentication with refresh tokens, bcrypt for password hashing, custom Fiber middleware for RBAC, Svelte 5 runes for frontend state management.

## Technical Context

**Language/Version**: Go 1.21+, JavaScript ES6+ (Svelte 5.0.0)  
**Primary Dependencies**: 
  - Backend: Fiber v2.52.0, pgx/v5 v5.5.0, golang-jwt/jwt/v5 (new), bcrypt (crypto/x)
  - Frontend: Svelte 5.0.0, Vite 6.0.0, Tailwind CSS 3.4.0
  
**Storage**: PostgreSQL 12+ (3 new tables: users, sessions, password_reset_tokens; modified: projects, tasks)  
**Testing**: Manual testing workflow (no automated test framework)  
**Target Platform**: Linux server (backend), Web browsers (frontend - Chrome, Firefox, Safari)  
**Project Type**: Web application (frontend + backend)  
**Performance Goals**: <2s authentication operations, 100 concurrent users, <200ms API response time (p95)  
**Constraints**: 
  - All UI text must be 100% Persian
  - httpOnly cookies only (no localStorage for tokens)
  - Email via Gmail SMTP (free tier: 500 emails/day)
  
**Scale/Scope**: 
  - Initial scale: Up to 1000 users
  - 12 new API endpoints
  - 5 new Svelte components
  - 1 new store
  - 3 new database tables
  - Estimated implementation: 14-22 hours

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Note**: No `.specify/memory/constitution.md` file found in project. Skipping constitution checks. This is acceptable as the project does not have documented architectural constraints beyond the existing layered backend structure and Svelte 5 frontend patterns.

**Architectural Compliance**:
- ✅ Follows existing layered backend architecture (models → repositories → services → handlers)
- ✅ Uses existing Svelte 5 patterns (runes, stores)
- ✅ Maintains backward compatibility with existing API structure
- ✅ No violations of existing patterns

## Project Structure

### Documentation (this feature)

```text
specs/003-user-auth/
├── plan.md                        # This file (/speckit.plan output)
├── research.md                    # Phase 0 output ✅
├── data-model.md                  # Phase 1 output ✅
├── quickstart.md                  # Phase 1 output ✅
├── contracts/                     # Phase 1 output ✅
│   └── api-endpoints.md
├── checklists/
│   └── requirements.md            # Spec validation checklist ✅
├── spec.md                        # Feature specification ✅
└── tasks.md                       # Phase 2 output (NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
backend/
├── models/
│   ├── user.go                    # New: User entity
│   ├── session.go                 # New: Session entity
│   ├── password_reset_token.go    # New: PasswordResetToken entity
│   ├── project.go                 # Modified: Add user_id, created_by
│   └── task.go                    # Modified: Add created_by
├── repositories/
│   ├── user_repository.go         # New: User CRUD operations
│   ├── session_repository.go      # New: Session management
│   ├── password_reset_repository.go # New: Reset token operations
│   ├── project_repository.go      # Modified: Filter by user_id
│   └── task_repository.go         # Modified: Filter by user_id
├── services/
│   ├── auth_service.go            # New: Authentication logic
│   ├── email_service.go           # New: Email sending (SMTP)
│   ├── project_service.go         # Modified: User ownership checks
│   └── task_service.go            # Modified: User ownership checks
├── handlers/
│   ├── auth_handler.go            # New: Auth endpoints
│   ├── user_handler.go            # New: User management (admin)
│   ├── project_handler.go         # Modified: Add ownership checks
│   └── task_handler.go            # Modified: Add ownership checks
├── middleware/
│   └── auth.go                    # New: Auth & RBAC middleware
├── routes/
│   └── routes.go                  # Modified: Add auth routes, protect existing
├── migration/
│   └── 002_add_user_authentication.sql # New: Database migration
└── config/
    └── jwt.go                     # New: JWT configuration

frontend/
├── src/
│   ├── components/
│   │   ├── LoginForm.svelte       # New: Login form
│   │   ├── RegisterForm.svelte    # New: Registration form
│   │   ├── ForgotPasswordForm.svelte # New: Forgot password
│   │   ├── ResetPasswordForm.svelte # New: Reset password
│   │   ├── UserManagement.svelte  # New: Admin user management
│   │   ├── App.svelte             # Modified: Auth check on mount
│   │   ├── ProjectList.svelte     # Modified: User filtering
│   │   └── TaskList.svelte        # Modified: User filtering
│   ├── stores/
│   │   └── authStore.js           # New: Auth state management
│   └── lib/
│       └── api.js                 # Modified: Add auth headers, handle 401
└── .env                           # New: Environment variables
```

**Structure Decision**: Web application structure (Option 2) with separate `backend/` and `frontend/` directories. This matches the existing project structure established in features 001 and 002.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No constitution violations. No complexity tracking needed.

---

## Phase 0: Research & Technical Decisions ✅

**Status**: Completed  
**Output**: [research.md](research.md)

### Research Summary

All technical unknowns from spec have been resolved:

1. **Password Hashing**: bcrypt (cost factor 10) chosen for simplicity and battle-tested security
2. **Session Management**: JWT tokens with httpOnly cookies (stateless, XSS-resistant)
3. **JWT Library**: golang-jwt/jwt/v5 (most popular, 8.8k+ stars)
4. **Account Lockout**: Two-layer approach (IP rate limiting + database lockout)
5. **Password Reset**: Cryptographically secure tokens with SHA-256 hashing, 1-hour expiry
6. **Email Service**: net/smtp stdlib with Gmail SMTP (no external dependencies)
7. **RBAC**: Custom Fiber middleware (sufficient for 2 roles)
8. **Svelte 5 Auth**: Runes ($state, $derived) + writable stores
9. **Persian Support**: Native UTF-8 (no special handling needed)
10. **Security Headers**: Fiber helmet middleware (CSP, HSTS, X-Frame-Options)

**Key Decisions**:
- All decisions prioritize simplicity and practicality for project scale
- No over-engineering: Using stdlib where possible
- Security best practices from OWASP applied
- Complete code examples provided for each area

See [research.md](research.md) for detailed rationale and implementation notes.

---

## Phase 1: Design & Data Model ✅

**Status**: Completed  
**Outputs**: 
- [data-model.md](data-model.md)
- [contracts/api-endpoints.md](contracts/api-endpoints.md)
- [quickstart.md](quickstart.md)

### Data Model Overview

**3 New Entities**:

1. **User**: User accounts with authentication credentials
   - Fields: id, username, email, password_hash, role, is_active, failed_login_attempts, locked_until, timestamps
   - Roles: 'admin' (full access), 'user' (limited to own resources)
   - Security: Account lockout after 5 failed attempts (30 min)

2. **Session**: Refresh tokens for JWT renewal
   - Fields: id, user_id, refresh_token_hash, user_agent, ip_address, expires_at, revoked
   - Expiry: 7 days
   - Purpose: Validate refresh tokens (access tokens are stateless)

3. **PasswordResetToken**: Temporary password reset tokens
   - Fields: id, user_id, token_hash, created_at, expires_at, used
   - Expiry: 1 hour
   - Single-use: Marked as used after password reset

**2 Modified Entities**:

1. **Project**: Add user ownership
   - New fields: user_id (owner), created_by (creator)
   - Access control: Regular users see only own projects, admins see all

2. **Task**: Add creator tracking
   - New field: created_by (creator)
   - Access control: Regular users see only tasks in own projects

**Database Migration**: `002_add_user_authentication.sql`
- Creates 3 new tables with indexes
- Alters 2 existing tables
- Seeds first admin user

See [data-model.md](data-model.md) for complete entity definitions, relationships, validation rules, and state transitions.

### API Contracts

**12 New Endpoints**:

**Public** (no auth):
1. `POST /api/auth/register` - Create account
2. `POST /api/auth/login` - Authenticate user
3. `POST /api/auth/forgot-password` - Request password reset
4. `POST /api/auth/reset-password` - Confirm password reset

**Protected** (requires auth):
5. `GET /api/auth/me` - Get current user
6. `PUT /api/auth/me` - Update current user
7. `PUT /api/auth/me/password` - Change password
8. `POST /api/auth/logout` - End session

**Admin-only** (requires admin role):
9. `GET /api/users` - List all users (paginated)
10. `GET /api/users/:id` - Get user details
11. `PUT /api/users/:id/role` - Change user role
12. `PUT /api/users/:id/activate` - Activate/deactivate user

**Modified Endpoints** (add auth):
- All existing endpoints now require authentication
- Regular users see filtered results (own resources only)
- Admins see all resources

All responses include Persian error messages. All endpoints use httpOnly cookies for authentication.

See [contracts/api-endpoints.md](contracts/api-endpoints.md) for complete request/response schemas and error codes.

### Frontend Components

**5 New Components**:

1. **LoginForm.svelte**: Email/password login with Persian labels
2. **RegisterForm.svelte**: User registration with validation
3. **ForgotPasswordForm.svelte**: Email input for password reset
4. **ResetPasswordForm.svelte**: New password input with token validation
5. **UserManagement.svelte**: Admin UI for user management (list, role change, activate/deactivate)

**1 New Store**:

1. **authStore.js**: Authentication state management
   - Methods: login(), logout(), checkAuth()
   - State: user, isAuthenticated, isLoading
   - Used by all components to check auth status

**Modified Components**:
- **App.svelte**: Check auth on mount, show login if not authenticated
- **ProjectList.svelte**: Filter projects by user_id (unless admin)
- **TaskList.svelte**: Filter tasks by project ownership
- **api.js**: Add `credentials: 'include'`, handle 401 (token expiry)

All components use Svelte 5 runes ($state, $derived, $effect) for reactivity.

---

## Phase 2: Implementation Tasks

**Status**: Not started (requires `/speckit.tasks` command)  
**Output**: Will be created in [tasks.md](tasks.md)

Implementation is divided into 5 phases aligned with user story priorities:

### Phase 1: Core Authentication (P1)
- Database migration
- User model and repository
- Auth service (bcrypt, JWT)
- Register/login endpoints
- Auth middleware
- Login/register forms
- Auth store
- Protect existing routes

### Phase 2: Password Reset (P2)
- Email service
- Password reset token model and repository
- Forgot/reset password endpoints
- Forgot/reset password forms

### Phase 3: Access Control (P2)
- Role-based middleware
- Update project/task services for ownership
- Update project/task handlers for filtering
- Update frontend to show role-appropriate UI

### Phase 4: Admin Management (P3)
- User management endpoints
- User management component
- Admin-only routes

### Phase 5: Security Hardening
- Rate limiting
- Security headers
- Account lockout mechanism
- Session cleanup job
- Token cleanup job

Run `/speckit.tasks` to generate detailed task breakdown with acceptance criteria.

---

## Implementation Priority Matrix

| Priority | User Story | Backend Tasks | Frontend Tasks | Estimated Hours |
|----------|------------|---------------|----------------|-----------------|
| P1 | User Registration | User model, auth service, register endpoint, middleware | RegisterForm, authStore, App routing | 4-6 hours |
| P1 | User Login | Login endpoint, JWT generation, session creation | LoginForm, authStore integration | 3-4 hours |
| P2 | Password Recovery | Email service, reset token model, forgot/reset endpoints | ForgotPasswordForm, ResetPasswordForm | 4-6 hours |
| P2 | Logout | Logout endpoint, session revocation | Logout button, auth cleanup | 1-2 hours |
| P2 | Access Control | Update all handlers, RBAC middleware | Filter UI based on role | 3-5 hours |
| P3 | Admin User Management | User management endpoints | UserManagement component | 2-4 hours |

**Total Estimated Time**: 14-22 hours

---

## Testing Strategy

**Manual Testing Workflow** (no automated tests):

### Phase 1 Testing (Core Auth)
1. ✅ Run migration → Verify tables created
2. ✅ Register user → Check database, receive JWT cookies
3. ✅ Login → Verify cookies set, user info returned
4. ✅ Access protected route → Should work with valid token
5. ✅ Access protected route without login → Should redirect to login
6. ✅ Invalid login → Should return Persian error message

### Phase 2 Testing (Password Reset)
7. ✅ Request password reset → Verify email sent (check console/mailbox)
8. ✅ Click reset link → Verify token validated
9. ✅ Reset password → Verify can login with new password
10. ✅ Use expired token → Should show error
11. ✅ Use token twice → Should show error

### Phase 3 Testing (Access Control)
12. ✅ Regular user login → Should see only own projects
13. ✅ Admin login → Should see all projects
14. ✅ Regular user try admin route → Should return 403
15. ✅ Admin access admin route → Should work

### Phase 4 Testing (Admin Management)
16. ✅ Admin view user list → Should see all users
17. ✅ Admin change user role → User should gain/lose admin access
18. ✅ Admin deactivate user → User cannot login
19. ✅ Admin try to deactivate self → Should fail with error

### Phase 5 Testing (Security)
20. ✅ Failed login 5 times → Account locked for 30 minutes
21. ✅ Too many password reset requests → Rate limit error
22. ✅ Token expiry → Should redirect to login with message

**Testing Tools**:
- Browser DevTools (check cookies, network requests, console)
- PostgreSQL client (verify database state)
- Email client or console logs (verify email sending)

---

## Environment Setup

### Required Environment Variables

Create `.env` file in project root:

```bash
# JWT Configuration
JWT_SECRET=<run: openssl rand -base64 32>
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h  # 7 days

# Email Configuration (Gmail)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-app@gmail.com
SMTP_PASSWORD=<generate-at-https://myaccount.google.com/apppasswords>

# Application URLs
APP_URL=http://localhost:5173
API_URL=http://localhost:3000

# Database (already configured in feature 001)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=project_management
```

### Dependency Installation

**Backend**:
```bash
cd backend
go get github.com/golang-jwt/jwt/v5
# bcrypt already available in golang.org/x/crypto
```

**Frontend**:
No new dependencies needed (Svelte 5, Tailwind CSS already installed)

---

## Security Considerations

### Password Security
- ✅ bcrypt hashing (cost factor 10)
- ✅ Minimum 8 characters with complexity requirements
- ✅ Never expose password_hash in API responses
- ✅ Secure password reset flow (hashed tokens, 1-hour expiry)

### Session Security
- ✅ JWT tokens stored in httpOnly cookies (prevents XSS)
- ✅ Refresh tokens hashed in database (like passwords)
- ✅ Short-lived access tokens (15 minutes)
- ✅ Long-lived refresh tokens (7 days, revocable)

### Account Security
- ✅ Account lockout after 5 failed attempts (30 minutes)
- ✅ Rate limiting on auth endpoints (10 req/min per IP)
- ✅ Password reset rate limiting (3 req/hour per IP)
- ✅ Single-use reset tokens

### API Security
- ✅ Security headers (CSP, HSTS, X-Frame-Options, etc.)
- ✅ CORS configuration (allow only frontend origin)
- ✅ Input validation (email format, password strength, etc.)
- ✅ Role-based access control on all endpoints

### Data Security
- ✅ All sensitive data hashed (passwords, tokens)
- ✅ No sensitive data in logs
- ✅ Persian error messages (no information leakage)
- ✅ Email enumeration prevention (always return success for password reset)

---

## Rollback Plan

If issues arise during implementation:

1. **Database Rollback**: Run down migration to drop new tables and columns
2. **Code Rollback**: Remove auth middleware from routes (routes become public again)
3. **Frontend Rollback**: Remove auth checks from App.svelte (remove login gate)

**Rollback Commands**:
```bash
# Database
psql -U postgres -d project_management -f backend/migration/002_add_user_authentication_down.sql

# Git (if needed)
git checkout main  # Return to stable branch
```

---

## Success Criteria Verification

After implementation, verify these measurable outcomes from the spec:

- [ ] **SC-001**: Users can complete registration in under 2 minutes
- [ ] **SC-002**: Users can login in under 30 seconds
- [ ] **SC-003**: 90% of users successfully complete registration/login on first attempt
- [ ] **SC-004**: Password recovery email received within 5 minutes
- [ ] **SC-005**: System handles 100 concurrent users without performance degradation
- [ ] **SC-006**: 100% of UI text is in Persian
- [ ] **SC-007**: Login error rate less than 5% (due to user errors, not system errors)
- [ ] **SC-008**: Authentication operations respond in under 2 seconds
- [ ] **SC-009**: 95% of brute force attacks detected and blocked

---

## Known Limitations & Future Enhancements

**Current Scope** (this feature):
- Basic email/password authentication
- Two static roles (admin, user)
- Gmail SMTP for email (500/day limit)

**Future Enhancements** (out of scope):
- Email verification (add email_verified boolean field)
- Two-factor authentication (TOTP)
- OAuth/SSO integration (Google, GitHub)
- Dynamic role/permission system (more than 2 roles)
- Professional email service (SendGrid, Mailgun)
- Remember me functionality (longer token expiry)
- Session management UI (view/revoke active sessions)
- Audit logging (track all auth events)

---

## References

- **Specification**: [spec.md](spec.md) - User requirements and acceptance criteria
- **Research**: [research.md](research.md) - Technical decisions and rationale
- **Data Model**: [data-model.md](data-model.md) - Entity definitions and relationships
- **API Contracts**: [contracts/api-endpoints.md](contracts/api-endpoints.md) - Endpoint specifications
- **Quickstart**: [quickstart.md](quickstart.md) - 30-second overview for developers
- **Agent Context**: [AGENTS.md](../../AGENTS.md) - Project-wide context for AI agents

---

**Document Version**: 1.0  
**Last Updated**: 2025-12-30  
**Status**: Planning Complete ✅ - Ready for `/speckit.tasks`
