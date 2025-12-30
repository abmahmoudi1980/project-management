# Quickstart Guide: User Authentication

**Feature**: 003-user-auth  
**For**: Developers implementing this feature  
**Time to Complete**: 30 seconds to understand the system

## 🎯 What This Feature Does

Adds complete user authentication to the project management app with:
- User registration and login (Persian UI)
- Two user roles: Admin (full access) and User (limited access)
- Password reset via email
- JWT-based session management with httpOnly cookies
- Account security (lockout after failed attempts)
- Admin user management

## 🏗️ Architecture at a Glance

```
Frontend (Svelte 5)          Backend (Go Fiber)          Database (PostgreSQL)
├─ Login Form                ├─ Auth Handlers            ├─ users table
├─ Register Form             ├─ Auth Middleware          ├─ sessions table
├─ Password Reset Form       ├─ JWT Generation           ├─ password_reset_tokens
├─ Auth Store (runes)        ├─ bcrypt Password Hash     └─ Modified: projects, tasks
└─ Protected Routes          ├─ Role-Based Access        
                             └─ Email Service (SMTP)
```

## 📊 Data Model Summary

**3 New Tables**:
1. **users**: User accounts (credentials, role, lockout status)
2. **sessions**: Refresh tokens (for JWT renewal)
3. **password_reset_tokens**: Temporary password reset tokens

**2 Modified Tables**:
- **projects**: Add `user_id` (owner) and `created_by`
- **tasks**: Add `created_by`

## 🔐 Authentication Flow

### Registration/Login
```
User → Register/Login Form (Svelte)
     → POST /api/auth/register or /api/auth/login
     → Backend validates & hashes password (bcrypt)
     → Generate JWT tokens (access + refresh)
     → Set httpOnly cookies
     → Return user info to frontend
     → Frontend updates auth store
     → User can access protected routes
```

### Password Reset
```
User → Forgot Password Form
     → POST /api/auth/forgot-password
     → Backend generates secure token
     → Email sent with reset link
     → User clicks link → Reset Password Form
     → POST /api/auth/reset-password with token
     → Backend validates token & updates password
```

### Session Management
```
Access Token (15 min) ──┐
                         ├─→ Every API request includes cookies
Refresh Token (7 days) ──┘   Backend validates token
                             If expired → 401 → Frontend redirects to login
```

## 🚀 Quick Implementation Checklist

### Phase 1: Backend Core (P1)
- [ ] Run database migration `002_add_user_authentication.sql`
- [ ] Add JWT dependency: `go get github.com/golang-jwt/jwt/v5`
- [ ] Create `models/user.go`, `models/session.go`, `models/password_reset_token.go`
- [ ] Create `repositories/user_repository.go`
- [ ] Create `services/auth_service.go` (bcrypt, JWT, email)
- [ ] Create `handlers/auth_handler.go` (register, login, logout)
- [ ] Create `middleware/auth.go` (RequireAuth, RequireRole)
- [ ] Update `routes/routes.go` (add auth routes, protect existing routes)

### Phase 2: Frontend Core (P1)
- [ ] Create `stores/authStore.js` (login, logout, checkAuth)
- [ ] Create `components/LoginForm.svelte` (Svelte 5 runes)
- [ ] Create `components/RegisterForm.svelte`
- [ ] Update `App.svelte` (check auth on mount, show login if not authenticated)
- [ ] Update `lib/api.js` (add `credentials: 'include'`, handle 401)
- [ ] Protect existing components (redirect to login if not authenticated)

### Phase 3: Password Reset (P2)
- [ ] Implement email service in backend (`services/email_service.go`)
- [ ] Create password reset handlers (forgot, reset)
- [ ] Create `components/ForgotPasswordForm.svelte`
- [ ] Create `components/ResetPasswordForm.svelte`

### Phase 4: Access Control (P2)
- [ ] Update project/task handlers to filter by user_id for regular users
- [ ] Ensure admins can access all resources
- [ ] Test role-based access in all endpoints

### Phase 5: Admin UI (P3)
- [ ] Create `components/UserManagement.svelte`
- [ ] Add admin routes for user listing, role change, activation
- [ ] Test admin-only features

## 🔑 Key Decisions (from Research)

| Decision | Choice | Why |
|----------|--------|-----|
| Password Hashing | bcrypt (cost 10) | Battle-tested, sufficient security |
| Session Management | JWT + httpOnly cookies | Stateless, secure (XSS-resistant) |
| JWT Library | golang-jwt/jwt/v5 | Most popular, well-maintained |
| Account Lockout | 5 attempts → 30 min | OWASP recommendation |
| Token Expiry | Access: 15 min, Refresh: 7 days | Balance security/UX |
| Email Service | net/smtp (Gmail) | No dependencies, simple |
| RBAC | Custom middleware | Only 2 roles, no need for complex library |
| Frontend Auth | Svelte 5 runes + stores | Native reactivity, simple |

## 📦 Dependencies to Add

**Backend**:
```bash
go get github.com/golang-jwt/jwt/v5
# bcrypt already available in golang.org/x/crypto
```

**Frontend**:
No new dependencies needed (Svelte 5, Tailwind already installed)

## 🌍 Environment Variables

Add to `.env` file:

```bash
# JWT Configuration
JWT_SECRET=<run: openssl rand -base64 32>
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Email Configuration (Gmail)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-app@gmail.com
SMTP_PASSWORD=<gmail-app-password>

# Application URLs
APP_URL=http://localhost:5173
API_URL=http://localhost:3000
```

**Gmail App Password**: Generate at https://myaccount.google.com/apppasswords

## 🧪 Testing Strategy

**Manual Testing Priority** (no test framework yet):
1. ✅ Register new user → Check database, try login
2. ✅ Login → Check cookies in browser DevTools
3. ✅ Access protected route → Should work after login
4. ✅ Logout → Cookies cleared, redirect to login
5. ✅ Failed login 5 times → Account locked for 30 min
6. ✅ Request password reset → Check email received
7. ✅ Reset password → Try login with new password
8. ✅ Admin: Change user role → User sees admin features
9. ✅ Admin: Deactivate user → User cannot login
10. ✅ Regular user: Try admin route → 403 Forbidden

## 📝 Persian Language Labels

All form labels and error messages must be in Persian. Examples:

**Forms**:
- Username: نام کاربری
- Email: ایمیل
- Password: رمز عبور
- Confirm Password: تکرار رمز عبور
- Login: ورود
- Register: ثبت‌نام
- Forgot Password: فراموشی رمز عبور
- Reset Password: بازیابی رمز عبور
- Logout: خروج

**Error Messages**:
- Invalid credentials: ایمیل یا رمز عبور نادرست است
- Email exists: این ایمیل قبلاً ثبت شده است
- Account locked: حساب کاربری قفل شده است
- Weak password: رمز عبور باید حداقل 8 کاراکتر و شامل حروف بزرگ، کوچک و اعداد باشد
- Unauthorized: احراز هویت نشده است
- Forbidden: دسترسی غیرمجاز

## 🚨 Common Pitfalls

1. **❌ Storing JWT in localStorage**: Use httpOnly cookies instead (prevents XSS)
2. **❌ Not setting `credentials: 'include'` in fetch**: Cookies won't be sent
3. **❌ Forgetting to hash tokens**: Always hash refresh tokens and reset tokens before DB storage
4. **❌ Not checking user role in handlers**: Middleware checks auth, handlers must check ownership
5. **❌ Hardcoding secret key**: Always use environment variables
6. **❌ Not handling 401 globally**: Add global fetch wrapper to redirect on token expiry

## 🔗 Reference Documents

- **Full Specification**: [spec.md](../spec.md)
- **Data Model Details**: [data-model.md](../data-model.md)
- **API Contracts**: [contracts/api-endpoints.md](../contracts/api-endpoints.md)
- **Research & Decisions**: [research.md](../research.md)
- **Implementation Plan**: [plan.md](../plan.md) ← Start here for step-by-step

## 💡 Pro Tips

1. **Start with P1 (Core Auth)**: Get login/register working first, then add extras
2. **Test on every commit**: Manual testing workflow is fast for auth features
3. **Use browser DevTools**: Check cookies, network requests, console errors
4. **Test both roles**: Always test as both admin and regular user
5. **Keep it simple**: Don't add features not in spec (email verification, 2FA, etc.)
6. **Persian first**: All UI text must be Persian from the start, not English placeholders

## 🎬 Getting Started

1. Read this quickstart (you're here! ✅)
2. Read [plan.md](../plan.md) for detailed implementation steps
3. Start with Phase 1: Backend database migration
4. Follow checklist above sequentially
5. Test as you go (don't wait until the end)
6. Refer to [research.md](../research.md) for technical details

**Estimated Implementation Time**: 
- P1 (Core): 8-12 hours
- P2 (Security): 4-6 hours
- P3 (Admin): 2-4 hours
- **Total**: ~14-22 hours

---

**Document Version**: 1.0  
**Last Updated**: 2025-12-30  
**Branch**: 003-user-auth
