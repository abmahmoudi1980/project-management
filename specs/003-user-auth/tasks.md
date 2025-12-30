# Tasks: User Authentication (احراز هویت کاربر)

**Input**: Design documents from `/specs/003-user-auth/`  
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**Tests**: Manual testing workflow (no automated tests)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `- [ ] [ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- **File paths**: All paths are relative to repository root

## Implementation Strategy

This feature implements 5 user stories in priority order:
- **P1**: User Story 1 (Registration) + User Story 2 (Login) = MVP
- **P2**: User Story 3 (Password Reset) + User Story 4 (Logout)
- **P3**: User Story 5 (Admin Management)

Each user story is independently testable and can be delivered as an increment.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Initialize authentication infrastructure and dependencies

**Duration**: ~1 hour

- [X] T001 Add JWT dependency: `cd backend && go get github.com/golang-jwt/jwt/v5`
- [X] T002 Create environment configuration file: `.env` with JWT_SECRET, SMTP credentials
- [X] T003 Generate JWT secret: Run `openssl rand -base64 32` and add to `.env`
- [X] T004 Create backend/middleware/ directory for authentication middleware
- [X] T005 Create backend/config/jwt.go for JWT configuration constants

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core authentication infrastructure that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

**Duration**: ~2 hours

### Database Setup

- [X] T006 Create database migration file: backend/migration/002_add_user_authentication.sql
- [X] T007 In migration: Create users table (id, username, email, password_hash, role, is_active, failed_login_attempts, locked_until, created_at, updated_at, last_login_at)
- [X] T008 In migration: Create sessions table (id, user_id, refresh_token_hash, user_agent, ip_address, created_at, expires_at, revoked)
- [X] T009 In migration: Create password_reset_tokens table (id, user_id, token_hash, created_at, expires_at, used)
- [X] T010 In migration: Add indexes (users.email, users.role, sessions.user_id, sessions.refresh_token_hash, password_reset_tokens.token_hash)
- [X] T011 In migration: ALTER TABLE projects ADD COLUMN user_id UUID REFERENCES users(id)
- [X] T012 In migration: ALTER TABLE projects ADD COLUMN created_by UUID REFERENCES users(id)
- [X] T013 In migration: ALTER TABLE tasks ADD COLUMN created_by UUID REFERENCES users(id)
- [X] T014 In migration: Insert seed admin user (email: admin@example.com, password: Admin123!, role: admin)
- [X] T015 Run migration: `cd backend && go run run_migration.go` and verify tables created

### Core Models

- [X] T016 [P] Create backend/models/user.go with User struct (ID, Username, Email, PasswordHash, Role, IsActive, FailedLoginAttempts, LockedUntil, CreatedAt, UpdatedAt, LastLoginAt)
- [X] T017 [P] In user.go: Add CreateUserRequest, LoginRequest, LoginResponse structs
- [X] T018 [P] Create backend/models/session.go with Session struct
- [X] T019 [P] Create backend/models/password_reset_token.go with PasswordResetToken struct

### Core Repositories

- [X] T020 [P] Create backend/repositories/user_repository.go with interface: Create, GetByID, GetByEmail, Update, List, UpdateFailedAttempts, LockAccount
- [X] T021 [P] Implement user repository methods with PostgreSQL queries
- [X] T022 [P] Create backend/repositories/session_repository.go with interface: Create, GetByRefreshToken, Revoke, DeleteExpired
- [X] T023 [P] Create backend/repositories/password_reset_repository.go with interface: Create, GetByToken, MarkAsUsed, DeleteExpired

### Core Services

- [X] T024 Create backend/services/auth_service.go with interface: Register, Login, VerifyPassword, GenerateTokens, ValidateToken, RefreshToken
- [X] T025 In auth_service.go: Implement Register (validate input, hash password with bcrypt cost 10, create user, generate tokens)
- [X] T026 In auth_service.go: Implement Login (get user by email, check password, check account not locked/inactive, update last_login_at, generate tokens)
- [X] T027 In auth_service.go: Implement GenerateTokens (create access JWT 15min, refresh JWT 7days, store refresh token hash in sessions table)
- [X] T028 In auth_service.go: Implement ValidateToken (parse JWT, verify signature, check expiration)
- [X] T029 In auth_service.go: Implement HandleFailedLogin (increment failed_attempts, lock account after 5 attempts for 30 minutes)

### Authentication Middleware

- [X] T030 Create backend/middleware/auth.go with RequireAuth middleware (extract JWT from cookie, validate, set user info in context)
- [X] T031 In auth.go: Add RequireRole(roles ...string) middleware (check user role from context against allowed roles)
- [X] T032 In auth.go: Add helper function getUserFromContext(c *fiber.Ctx) to extract user info from Fiber context

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - ثبت‌نام کاربر (User Registration) (Priority: P1) 🎯 MVP PART 1

**Goal**: Users can create new accounts with username, email, and password. Accounts are created with "user" role by default. Validation ensures email uniqueness and password strength. All UI labels in Persian.

**Independent Test**: 
1. Open registration form at http://localhost:5173
2. Fill in username, email, password (8+ chars, uppercase, lowercase, digit)
3. Submit form
4. Verify account created in database
5. Verify user automatically logged in with JWT cookies
6. Verify can access protected routes

**Duration**: ~3 hours

### Backend Implementation

- [X] T033 [US1] Create backend/handlers/auth_handler.go with Register handler
- [X] T034 [US1] In Register handler: Validate request body (username 3-50 chars, valid email, password 8+ chars with uppercase/lowercase/digit)
- [X] T035 [US1] In Register handler: Return Persian error for duplicate email: "این ایمیل قبلاً ثبت شده است"
- [X] T036 [US1] In Register handler: Return Persian error for weak password: "رمز عبور باید حداقل 8 کاراکتر و شامل حروف بزرگ، کوچک و اعداد باشد"
- [X] T037 [US1] In Register handler: Call auth_service.Register, set httpOnly cookies (access_token, refresh_token), return user info (exclude password_hash)
- [X] T038 [US1] In backend/routes/routes.go: Add POST /api/auth/register route (public, no auth required)

### Frontend Implementation

- [X] T039 [P] [US1] Create frontend/src/stores/authStore.js with writable store (user, isAuthenticated, isLoading)
- [X] T040 [P] [US1] In authStore.js: Add register(username, email, password) method that POSTs to /api/auth/register
- [X] T041 [P] [US1] In authStore.js: Add checkAuth() method that GETs /api/auth/me to check if user is logged in
- [X] T042 [P] [US1] In authStore.js: Add logout() method that POSTs to /api/auth/logout and clears user state
- [X] T043 [P] [US1] Create frontend/src/components/RegisterForm.svelte with Svelte 5 runes
- [X] T044 [US1] In RegisterForm.svelte: Add form fields with Persian labels (نام کاربری, ایمیل, رمز عبور, تکرار رمز عبور)
- [X] T045 [US1] In RegisterForm.svelte: Add client-side validation (check password confirmation matches)
- [X] T046 [US1] In RegisterForm.svelte: Display Persian error messages from API
- [X] T047 [US1] In RegisterForm.svelte: Use $state for form fields, $derived for validation, call authStore.register() on submit
- [X] T048 [US1] Update frontend/src/lib/api.js: Add `credentials: 'include'` to all fetch calls to send cookies
- [X] T049 [US1] Update frontend/src/App.svelte: Add route logic to show RegisterForm when not authenticated

**Manual Test Checklist for US1**:
- ✅ Register with valid data → Account created, automatically logged in
- ✅ Register with duplicate email → Persian error shown
- ✅ Register with weak password → Persian validation message shown
- ✅ Register with empty fields → Persian error messages shown
- ✅ Password confirmation mismatch → Persian error shown

---

## Phase 4: User Story 2 - ورود کاربر (User Login) (Priority: P1) 🎯 MVP PART 2

**Goal**: Registered users can log in with email and password. System authenticates user and creates session with JWT tokens in httpOnly cookies. User role (admin/user) determines access level. All error messages in Persian.

**Independent Test**:
1. Open login form at http://localhost:5173
2. Enter email and password from registration
3. Submit form
4. Verify JWT cookies set in browser DevTools
5. Verify redirected to main app
6. Admin: Verify can see all projects/tasks
7. Regular user: Verify sees only own projects/tasks

**Duration**: ~4 hours

### Backend Implementation

- [X] T050 [US2] In backend/handlers/auth_handler.go: Add Login handler
- [X] T051 [US2] In Login handler: Validate request (email and password required)
- [X] T052 [US2] In Login handler: Check account not locked (locked_until < NOW), return Persian error: "حساب کاربری شما قفل شده است. لطفاً 30 دقیقه صبر کنید"
- [X] T053 [US2] In Login handler: Check account is active (is_active = true), return Persian error: "حساب کاربری شما غیرفعال شده است"
- [X] T054 [US2] In Login handler: Call auth_service.Login, handle invalid credentials with Persian error: "ایمیل یا رمز عبور نادرست است"
- [X] T055 [US2] In Login handler: On failed login, call auth_service.HandleFailedLogin to increment counter and lock after 5 attempts
- [X] T056 [US2] In Login handler: On successful login, reset failed_login_attempts to 0, set httpOnly cookies, return user info
- [X] T057 [US2] In backend/routes/routes.go: Add POST /api/auth/login route (public)
- [X] T058 [US2] In routes.go: Add GET /api/auth/me route (protected with RequireAuth middleware) to get current user info
- [X] T059 [US2] Create backend/handlers/auth_handler.go: Add GetCurrentUser handler that returns user info from context (exclude password_hash)

### Role-Based Access Control

- [X] T060 [US2] In routes.go: Wrap all existing routes (projects, tasks, timelogs) with RequireAuth middleware
- [ ] T061 [US2] Update backend/services/project_service.go: Add GetProjectsByUser(userID, role) method (admins see all, users see own)
- [ ] T062 [US2] Update backend/handlers/project_handler.go: In GetProjects, call GetProjectsByUser with user from context
- [ ] T063 [US2] Update backend/services/task_service.go: Add GetTasksByUser(userID, role) method (filter by project ownership)
- [ ] T064 [US2] Update backend/handlers/task_handler.go: In GetTasks, call GetTasksByUser with user from context

### Frontend Implementation

- [X] T065 [P] [US2] Create frontend/src/components/LoginForm.svelte with Svelte 5 runes
- [X] T066 [US2] In LoginForm.svelte: Add form with Persian labels (ایمیل, رمز عبور, ورود)
- [X] T067 [US2] In LoginForm.svelte: Display Persian error messages (invalid credentials, account locked, account deactivated)
- [X] T068 [US2] In LoginForm.svelte: Use $state for form fields, call authStore.login() on submit
- [X] T069 [US2] In authStore.js: Implement login(email, password) that POSTs to /api/auth/login and updates store on success
- [X] T070 [US2] Update frontend/src/App.svelte: Call authStore.checkAuth() in onMount to check if user is already logged in
- [X] T071 [US2] In App.svelte: Use $effect to redirect to login page if not authenticated
- [X] T072 [US2] In App.svelte: Show loading state while checkAuth() is running
- [X] T073 [US2] Update frontend/src/lib/api.js: Add global 401 handler that calls authStore.logout() and redirects to login

### Access Control UI

- [X] T074 [US2] Update frontend/src/components/ProjectList.svelte: Show all projects if admin, filter if regular user
- [X] T075 [US2] Update frontend/src/components/TaskList.svelte: Show all tasks if admin, filter by project ownership if regular user
- [X] T076 [US2] In App.svelte: Show user role badge (ادمین / کاربر عادی) in header

**Manual Test Checklist for US2**:
- ✅ Login with valid credentials → Logged in, cookies set, redirected to main app
- ✅ Login with invalid credentials → Persian error shown
- ✅ Login with inactive account → Persian error shown
- ✅ Failed login 5 times → Account locked with Persian message, cannot login for 30 minutes
- ✅ After 30 minutes → Account unlocked, can login again
- ✅ Admin user → Sees all projects/tasks
- ✅ Regular user → Sees only own projects/tasks
- ✅ Access protected route without login → Redirected to login page

---

## Phase 5: User Story 3 - بازیابی رمز عبور (Password Recovery) (Priority: P2)

**Goal**: Users who forgot their password can request a reset link via email. System sends time-limited (1 hour) single-use token. User clicks link, enters new password, and can login with new credentials. All messages in Persian.

**Independent Test**:
1. Open forgot password form
2. Enter registered email
3. Check email inbox for reset link (or check console logs in development)
4. Click reset link with token
5. Enter new password
6. Verify password changed in database
7. Login with new password → Success

**Duration**: ~4 hours

### Backend - Email Service

- [ ] T077 [P] [US3] Create backend/services/email_service.go with SendEmail(to, subject, body) function
- [ ] T078 [P] [US3] In email_service.go: Implement SMTP client using net/smtp with Gmail configuration from .env
- [ ] T079 [P] [US3] In email_service.go: Add SendPasswordResetEmail(to, resetLink) function with Persian email template
- [ ] T080 [P] [US3] Persian email template with RTL direction: "برای بازیابی رمز عبور خود، روی لینک زیر کلیک کنید"

### Backend - Password Reset Logic

- [ ] T081 [US3] In backend/services/auth_service.go: Add RequestPasswordReset(email) method
- [ ] T082 [US3] In RequestPasswordReset: Generate cryptographically secure 32-byte token using crypto/rand
- [ ] T083 [US3] In RequestPasswordReset: Hash token with SHA-256, store hash in password_reset_tokens table with 1-hour expiry
- [ ] T084 [US3] In RequestPasswordReset: Send plain token via email (format: APP_URL/reset-password?token=xxx)
- [ ] T085 [US3] In RequestPasswordReset: Always return success message (prevent email enumeration): "اگر ایمیل شما در سیستم ثبت باشد، لینک بازیابی ارسال می‌شود"
- [ ] T086 [US3] In auth_service.go: Add ResetPassword(token, newPassword) method
- [ ] T087 [US3] In ResetPassword: Hash received token, lookup in database, verify not expired and not used
- [ ] T088 [US3] In ResetPassword: Return Persian error for invalid/expired token: "لینک بازیابی نامعتبر یا منقضی شده است"
- [ ] T089 [US3] In ResetPassword: Validate new password strength, hash with bcrypt, update user password, mark token as used

### Backend - API Endpoints

- [ ] T090 [US3] In backend/handlers/auth_handler.go: Add ForgotPassword handler
- [ ] T091 [US3] In ForgotPassword handler: Rate limit to 3 requests per hour per IP
- [ ] T092 [US3] In ForgotPassword handler: Call auth_service.RequestPasswordReset, return success message
- [ ] T093 [US3] In auth_handler.go: Add ResetPassword handler
- [ ] T094 [US3] In ResetPassword handler: Validate token and newPassword from request body
- [ ] T095 [US3] In ResetPassword handler: Return Persian validation errors for weak password
- [ ] T096 [US3] In ResetPassword handler: Call auth_service.ResetPassword, return success message
- [ ] T097 [US3] In backend/routes/routes.go: Add POST /api/auth/forgot-password (public)
- [ ] T098 [US3] In routes.go: Add POST /api/auth/reset-password (public)

### Frontend Implementation

- [ ] T099 [P] [US3] Create frontend/src/components/ForgotPasswordForm.svelte with Svelte 5 runes
- [ ] T100 [US3] In ForgotPasswordForm.svelte: Add email input field with Persian label (ایمیل)
- [ ] T101 [US3] In ForgotPasswordForm.svelte: POST to /api/auth/forgot-password, show success message
- [ ] T102 [US3] In ForgotPasswordForm.svelte: Display Persian message: "اگر ایمیل شما در سیستم ثبت باشد، لینک بازیابی ارسال می‌شود"
- [ ] T103 [P] [US3] Create frontend/src/components/ResetPasswordForm.svelte
- [ ] T104 [US3] In ResetPasswordForm.svelte: Extract token from URL query parameter
- [ ] T105 [US3] In ResetPasswordForm.svelte: Add fields for new password and confirmation with Persian labels
- [ ] T106 [US3] In ResetPasswordForm.svelte: POST to /api/auth/reset-password with token and new password
- [ ] T107 [US3] In ResetPasswordForm.svelte: Display Persian success message: "رمز عبور شما با موفقیت تغییر یافت"
- [ ] T108 [US3] In ResetPasswordForm.svelte: Display Persian error for invalid/expired token
- [ ] T109 [US3] In App.svelte: Add route for /reset-password to show ResetPasswordForm

**Manual Test Checklist for US3**:
- ✅ Request password reset with registered email → Email received (check inbox or logs)
- ✅ Click reset link → ResetPasswordForm shown
- ✅ Enter new password → Password changed, success message in Persian
- ✅ Login with new password → Success
- ✅ Request reset with unregistered email → Generic success message (no error)
- ✅ Use expired token (after 1 hour) → Persian error shown
- ✅ Use token twice → Persian error shown (token already used)
- ✅ Request multiple resets quickly → Rate limit error after 3 attempts

---

## Phase 6: User Story 4 - خروج از سیستم (Logout) (Priority: P2)

**Goal**: Authenticated users can logout, which revokes their session, clears cookies, and redirects to login page. User cannot access protected routes after logout.

**Independent Test**:
1. Login as any user
2. Click logout button
3. Verify cookies cleared in browser DevTools
4. Verify session revoked in database (revoked = true)
5. Try to access protected route → Redirected to login

**Duration**: ~1 hour

### Backend Implementation

- [ ] T110 [US4] In backend/handlers/auth_handler.go: Add Logout handler
- [ ] T111 [US4] In Logout handler: Extract refresh_token from cookie
- [ ] T112 [US4] In Logout handler: Call session_repository.Revoke to mark session as revoked
- [ ] T113 [US4] In Logout handler: Clear access_token and refresh_token cookies (set Max-Age = 0)
- [ ] T114 [US4] In Logout handler: Return Persian success message: "با موفقیت خارج شدید"
- [ ] T115 [US4] In backend/routes/routes.go: Add POST /api/auth/logout (protected with RequireAuth)

### Frontend Implementation

- [ ] T116 [US4] Update frontend/src/stores/authStore.js: Implement logout() method that POSTs to /api/auth/logout
- [ ] T117 [US4] In authStore logout(): Clear user state (user = null, isAuthenticated = false)
- [ ] T118 [US4] In App.svelte: Add logout button in header with Persian label (خروج)
- [ ] T119 [US4] In App.svelte: On logout button click, call authStore.logout() and redirect to login page

**Manual Test Checklist for US4**:
- ✅ Logout → Cookies cleared, redirected to login page
- ✅ After logout, try to access protected route → Redirected to login
- ✅ After logout, check database → Session marked as revoked

---

## Phase 7: User Story 5 - مدیریت کاربران (Admin User Management) (Priority: P3)

**Goal**: Admin users can view list of all users, change user roles (admin/user), and activate/deactivate user accounts. Only admins can access user management. All UI in Persian.

**Independent Test**:
1. Login as admin user
2. Navigate to user management page
3. View list of all users with their roles and statuses
4. Change a user's role from "user" to "admin"
5. Verify that user now has admin access
6. Deactivate a user
7. Try to login as deactivated user → Error message
8. Login as regular user → User management page not accessible (403)

**Duration**: ~3 hours

### Backend - User Management Endpoints

- [ ] T120 [P] [US5] Create backend/handlers/user_handler.go with GetUsers handler
- [ ] T121 [US5] In GetUsers handler: Support pagination (page, limit query params), default limit 20, max 100
- [ ] T122 [US5] In GetUsers handler: Support filtering by role and is_active
- [ ] T123 [US5] In GetUsers handler: Call user_repository.List, return users with pagination metadata (exclude password_hash)
- [ ] T124 [P] [US5] In user_handler.go: Add GetUserByID handler that returns user details (admin only)
- [ ] T125 [P] [US5] In user_handler.go: Add UpdateUserRole handler
- [ ] T126 [US5] In UpdateUserRole: Validate role is "admin" or "user", return Persian error: "نقش نامعتبر است"
- [ ] T127 [US5] In UpdateUserRole: Call user_repository.Update to change role, return updated user
- [ ] T128 [P] [US5] In user_handler.go: Add UpdateUserActivation handler
- [ ] T129 [US5] In UpdateUserActivation: Prevent admin from deactivating themselves, return Persian error: "نمی‌توانید حساب کاربری خود را غیرفعال کنید"
- [ ] T130 [US5] In UpdateUserActivation: Prevent deactivating last admin, return Persian error: "نمی‌توانید آخرین ادمین را غیرفعال کنید"
- [ ] T131 [US5] In UpdateUserActivation: Call user_repository.Update to change is_active, return updated user

### Backend - Admin Routes

- [ ] T132 [US5] In backend/routes/routes.go: Add GET /api/users (protected with RequireAuth + RequireRole("admin"))
- [ ] T133 [US5] In routes.go: Add GET /api/users/:id (admin only)
- [ ] T134 [US5] In routes.go: Add PUT /api/users/:id/role (admin only)
- [ ] T135 [US5] In routes.go: Add PUT /api/users/:id/activate (admin only)

### Frontend - User Management Component

- [ ] T136 [P] [US5] Create frontend/src/components/UserManagement.svelte with Svelte 5 runes
- [ ] T137 [US5] In UserManagement.svelte: Fetch users from GET /api/users on component mount
- [ ] T138 [US5] In UserManagement.svelte: Display users in table with columns: نام کاربری, ایمیل, نقش, وضعیت, عملیات
- [ ] T139 [US5] In UserManagement.svelte: Add pagination controls (قبلی / بعدی)
- [ ] T140 [US5] In UserManagement.svelte: Add role dropdown for each user (ادمین / کاربر عادی)
- [ ] T141 [US5] In UserManagement.svelte: On role change, PUT to /api/users/:id/role and refresh list
- [ ] T142 [US5] In UserManagement.svelte: Add activate/deactivate toggle button (فعال / غیرفعال)
- [ ] T143 [US5] In UserManagement.svelte: On activation change, PUT to /api/users/:id/activate and refresh list
- [ ] T144 [US5] In UserManagement.svelte: Display Persian error messages from API
- [ ] T145 [US5] In UserManagement.svelte: Show confirmation dialog before role change or deactivation
- [ ] T146 [US5] In App.svelte: Add "User Management" menu item in navigation (only show if user is admin)
- [ ] T147 [US5] In App.svelte: Add route for user management page that shows UserManagement component

**Manual Test Checklist for US5**:
- ✅ Login as admin → User management menu visible
- ✅ Navigate to user management → See list of all users
- ✅ Change user role → User role updated, user gains/loses admin access
- ✅ Deactivate user → User cannot login
- ✅ Reactivate user → User can login again
- ✅ Try to deactivate self as admin → Persian error shown
- ✅ Try to deactivate last admin → Persian error shown
- ✅ Login as regular user → User management menu not visible, 403 if accessing URL directly

---

## Phase 8: Security & Polish (Cross-Cutting Concerns)

**Purpose**: Security hardening, performance optimization, and UI polish

**Duration**: ~2 hours

### Security Headers & Rate Limiting

- [ ] T148 [P] Configure Fiber helmet middleware in backend/main.go with security headers (CSP, HSTS, X-Frame-Options, X-Content-Type-Options)
- [ ] T149 [P] Configure CORS middleware in main.go (allow origin: http://localhost:5173, allow credentials: true)
- [ ] T150 [P] Add IP-based rate limiter middleware (10 requests/minute per IP) to all auth endpoints
- [ ] T151 [P] Add stricter rate limiter to login endpoint (5 attempts per 5 minutes per IP)

### Password Change & Profile Update

- [ ] T152 [P] In backend/handlers/auth_handler.go: Add UpdateProfile handler (update username, email)
- [ ] T153 [P] In auth_handler.go: Add ChangePassword handler (verify current password, update to new password)
- [ ] T154 [P] In backend/routes/routes.go: Add PUT /api/auth/me (protected)
- [ ] T155 [P] In routes.go: Add PUT /api/auth/me/password (protected)
- [ ] T156 [P] Create frontend/src/components/ProfileSettings.svelte with forms for profile update and password change

### UI Polish & Persian Labels

- [ ] T157 [P] Review all components for 100% Persian labels and error messages
- [ ] T158 [P] Add loading spinners to all forms during API calls
- [ ] T159 [P] Add success toast notifications for all operations (ثبت‌نام موفق، ورود موفق، etc.)
- [ ] T160 [P] Add RTL (direction: rtl) styling to all forms and text in Tailwind config
- [ ] T161 [P] Style all forms consistently with Tailwind classes (shadows, borders, focus states)

### Cleanup & Maintenance

- [ ] T162 [P] Create database cleanup job script: backend/scripts/cleanup_expired_sessions.go (delete sessions where expires_at < NOW)
- [ ] T163 [P] Create cleanup job for password reset tokens: Delete tokens where (expires_at < NOW OR used = true)
- [ ] T164 [P] Add documentation comment to .env.example with all required environment variables
- [ ] T165 [P] Update README.md with authentication setup instructions

---

## Dependencies Graph

**User Story Dependencies**:
```
Setup (Phase 1)
  ↓
Foundational (Phase 2) ← MUST complete first
  ↓
  ├─→ US1: Registration (P1) ──┐
  │                             ├─→ MVP READY
  ├─→ US2: Login (P1) ─────────┘
  │
  ├─→ US3: Password Reset (P2)
  │
  ├─→ US4: Logout (P2)
  │
  └─→ US5: Admin Management (P3)
```

**Key Insights**:
- **Foundational phase BLOCKS all user stories** - must be completed first
- **US1 and US2 are independent** after Foundation - can be developed in parallel by 2 developers
- **US3, US4, US5 depend on US1+US2** (need auth to test)
- **US3, US4, US5 are independent of each other** - can be parallelized

**Parallel Execution Examples**:

**2-Developer Team**:
- Dev 1: Foundation → US1 → US3 → US5
- Dev 2: Foundation → US2 → US4 → Polish

**Single Developer**:
- Day 1: Foundation + US1 (6 hours)
- Day 2: US2 (4 hours)
- Day 3: US3 + US4 (5 hours)
- Day 4: US5 + Polish (5 hours)

---

## Task Summary

**Total Tasks**: 165  
**Parallelizable Tasks**: 48 (marked with [P])

**By Phase**:
- Phase 1 (Setup): 5 tasks (~1 hour)
- Phase 2 (Foundational): 27 tasks (~2 hours) - BLOCKING
- Phase 3 (US1 - Registration): 17 tasks (~3 hours)
- Phase 4 (US2 - Login): 27 tasks (~4 hours)
- Phase 5 (US3 - Password Reset): 33 tasks (~4 hours)
- Phase 6 (US4 - Logout): 10 tasks (~1 hour)
- Phase 7 (US5 - Admin Management): 28 tasks (~3 hours)
- Phase 8 (Security & Polish): 18 tasks (~2 hours)

**Total Estimated Time**: 20 hours (single developer, sequential)  
**With Parallelization**: ~14 hours (2 developers)

**MVP Scope** (US1 + US2): 44 tasks, ~8 hours → Core authentication working

---

## Format Validation

✅ All tasks follow the required format: `- [ ] [ID] [P?] [Story?] Description`  
✅ All task IDs are sequential (T001 - T165)  
✅ Parallelizable tasks marked with [P]  
✅ User story tasks labeled with [US1]-[US5]  
✅ All tasks include file paths  
✅ Tasks organized by user story for independent implementation  
✅ Each phase has clear checkpoints for validation

---

**Document Version**: 1.0  
**Last Updated**: 2025-12-30  
**Status**: Ready for Implementation ✅
