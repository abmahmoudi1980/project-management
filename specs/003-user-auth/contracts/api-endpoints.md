# API Contracts: User Authentication

**Feature**: 003-user-auth  
**Date**: 2025-12-30  
**Base URL**: `http://localhost:3000/api`

## Overview

This document defines all HTTP API endpoints for the user authentication system. All endpoints return JSON responses with Persian error messages.

### Authentication Flow

1. **Public Endpoints** (no authentication required):
   - `POST /auth/register` - User registration
   - `POST /auth/login` - User login
   - `POST /auth/forgot-password` - Request password reset
   - `POST /auth/reset-password` - Confirm password reset

2. **Protected Endpoints** (requires valid access token):
   - `POST /auth/logout` - User logout
   - `GET /auth/me` - Get current user info
   - `PUT /auth/me` - Update current user info
   - `PUT /auth/me/password` - Change password

3. **Admin-Only Endpoints** (requires admin role):
   - `GET /users` - List all users
   - `GET /users/:id` - Get user by ID
   - `PUT /users/:id/role` - Change user role
   - `PUT /users/:id/activate` - Activate/deactivate user

---

## Common Patterns

### Authentication Header

Protected endpoints require JWT access token in httpOnly cookie:

```http
Cookie: access_token=eyJhbGciOiJIUzI1NiIs...
```

The frontend automatically sends cookies with `credentials: 'include'` in fetch requests.

### Success Response Format

```json
{
  "success": true,
  "data": { /* response data */ }
}
```

### Error Response Format

```json
{
  "success": false,
  "error": {
    "message": "پیام خطا به فارسی",
    "code": "ERROR_CODE",
    "details": { /* optional error details */ }
  }
}
```

### HTTP Status Codes

| Code | Meaning | Persian Message |
|------|---------|----------------|
| 200 | OK | موفقیت‌آمیز |
| 201 | Created | ایجاد شد |
| 400 | Bad Request | درخواست نامعتبر |
| 401 | Unauthorized | احراز هویت نشده است |
| 403 | Forbidden | دسترسی غیرمجاز |
| 404 | Not Found | یافت نشد |
| 409 | Conflict | تداخل داده |
| 429 | Too Many Requests | تعداد درخواست‌ها بیش از حد مجاز |
| 500 | Internal Server Error | خطای سرور |

---

## Public Endpoints

### 1. Register User

Create a new user account with "user" role.

**Endpoint**: `POST /api/auth/register`

**Request Body**:
```json
{
  "username": "علی_احمدی",
  "email": "ali@example.com",
  "password": "MyP@ssw0rd",
  "password_confirmation": "MyP@ssw0rd"
}
```

**Validation Rules**:
- `username`: 3-50 chars, alphanumeric + underscore/hyphen
- `email`: Valid email format, unique
- `password`: Min 8 chars, must include uppercase, lowercase, digit
- `password_confirmation`: Must match `password`

**Success Response** (201 Created):
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "علی_احمدی",
      "email": "ali@example.com",
      "role": "user",
      "is_active": true,
      "created_at": "2025-12-30T10:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

Sets httpOnly cookies:
- `access_token` (expires in 15 minutes)
- `refresh_token` (expires in 7 days)

**Error Responses**:

- **400 Bad Request** (Validation error):
```json
{
  "success": false,
  "error": {
    "message": "خطای اعتبارسنجی",
    "code": "VALIDATION_ERROR",
    "details": {
      "password": "رمز عبور باید حداقل 8 کاراکتر و شامل حروف بزرگ، کوچک و اعداد باشد"
    }
  }
}
```

- **409 Conflict** (Duplicate email):
```json
{
  "success": false,
  "error": {
    "message": "این ایمیل قبلاً ثبت شده است",
    "code": "EMAIL_EXISTS"
  }
}
```

---

### 2. Login User

Authenticate user and create session.

**Endpoint**: `POST /api/auth/login`

**Request Body**:
```json
{
  "email": "ali@example.com",
  "password": "MyP@ssw0rd"
}
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "علی_احمدی",
      "email": "ali@example.com",
      "role": "user",
      "is_active": true,
      "last_login_at": "2025-12-30T10:30:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

Sets httpOnly cookies:
- `access_token` (expires in 15 minutes)
- `refresh_token` (expires in 7 days)

**Error Responses**:

- **401 Unauthorized** (Invalid credentials):
```json
{
  "success": false,
  "error": {
    "message": "ایمیل یا رمز عبور نادرست است",
    "code": "INVALID_CREDENTIALS"
  }
}
```

- **403 Forbidden** (Account locked):
```json
{
  "success": false,
  "error": {
    "message": "حساب کاربری شما به دلیل تلاش‌های ناموفق متوالی قفل شده است. لطفاً 30 دقیقه صبر کنید",
    "code": "ACCOUNT_LOCKED",
    "details": {
      "locked_until": "2025-12-30T11:00:00Z"
    }
  }
}
```

- **403 Forbidden** (Account deactivated):
```json
{
  "success": false,
  "error": {
    "message": "حساب کاربری شما غیرفعال شده است",
    "code": "ACCOUNT_DEACTIVATED"
  }
}
```

- **429 Too Many Requests** (Rate limit exceeded):
```json
{
  "success": false,
  "error": {
    "message": "تعداد تلاش‌های ورود بیش از حد مجاز است. لطفاً بعداً تلاش کنید",
    "code": "TOO_MANY_REQUESTS"
  }
}
```

---

### 3. Request Password Reset

Request a password reset email.

**Endpoint**: `POST /api/auth/forgot-password`

**Request Body**:
```json
{
  "email": "ali@example.com"
}
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "message": "اگر ایمیل شما در سیستم ثبت باشد، لینک بازیابی رمز عبور به آن ارسال می‌شود"
  }
}
```

**Notes**:
- Always returns success (even if email not found) to prevent email enumeration
- Rate limited: Max 3 requests per hour per IP
- Email contains reset link: `http://localhost:5173/reset-password?token=<token>`
- Token expires in 1 hour

**Error Responses**:

- **429 Too Many Requests** (Rate limit):
```json
{
  "success": false,
  "error": {
    "message": "شما درخواست‌های بازیابی رمز عبور زیادی ارسال کرده‌اید. لطفاً بعداً تلاش کنید",
    "code": "TOO_MANY_RESET_REQUESTS"
  }
}
```

---

### 4. Reset Password

Reset password using token from email.

**Endpoint**: `POST /api/auth/reset-password`

**Request Body**:
```json
{
  "token": "abc123def456...",
  "new_password": "NewP@ssw0rd",
  "new_password_confirmation": "NewP@ssw0rd"
}
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "message": "رمز عبور شما با موفقیت تغییر یافت. اکنون می‌توانید وارد شوید"
  }
}
```

**Error Responses**:

- **400 Bad Request** (Invalid token):
```json
{
  "success": false,
  "error": {
    "message": "لینک بازیابی نامعتبر یا منقضی شده است",
    "code": "INVALID_RESET_TOKEN"
  }
}
```

- **400 Bad Request** (Token already used):
```json
{
  "success": false,
  "error": {
    "message": "این لینک بازیابی قبلاً استفاده شده است",
    "code": "TOKEN_ALREADY_USED"
  }
}
```

- **400 Bad Request** (Password validation):
```json
{
  "success": false,
  "error": {
    "message": "خطای اعتبارسنجی",
    "code": "VALIDATION_ERROR",
    "details": {
      "new_password": "رمز عبور باید حداقل 8 کاراکتر و شامل حروف بزرگ، کوچک و اعداد باشد"
    }
  }
}
```

---

## Protected Endpoints (Requires Authentication)

### 5. Get Current User

Get authenticated user's information.

**Endpoint**: `GET /api/auth/me`

**Headers**:
```http
Cookie: access_token=eyJhbGciOiJIUzI1NiIs...
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "علی_احمدی",
      "email": "ali@example.com",
      "role": "user",
      "is_active": true,
      "created_at": "2025-12-30T10:00:00Z",
      "last_login_at": "2025-12-30T10:30:00Z"
    }
  }
}
```

**Error Responses**:

- **401 Unauthorized** (No token or invalid token):
```json
{
  "success": false,
  "error": {
    "message": "احراز هویت نشده است",
    "code": "UNAUTHORIZED"
  }
}
```

---

### 6. Update Current User

Update authenticated user's profile.

**Endpoint**: `PUT /api/auth/me`

**Headers**:
```http
Cookie: access_token=eyJhbGciOiJIUzI1NiIs...
```

**Request Body**:
```json
{
  "username": "علی_احمدزاده",
  "email": "ali.new@example.com"
}
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "علی_احمدزاده",
      "email": "ali.new@example.com",
      "role": "user",
      "is_active": true,
      "updated_at": "2025-12-30T11:00:00Z"
    }
  }
}
```

**Error Responses**:

- **409 Conflict** (Email already taken):
```json
{
  "success": false,
  "error": {
    "message": "این ایمیل توسط کاربر دیگری استفاده می‌شود",
    "code": "EMAIL_EXISTS"
  }
}
```

---

### 7. Change Password

Change authenticated user's password.

**Endpoint**: `PUT /api/auth/me/password`

**Headers**:
```http
Cookie: access_token=eyJhbGciOiJIUzI1NiIs...
```

**Request Body**:
```json
{
  "current_password": "MyP@ssw0rd",
  "new_password": "NewP@ssw0rd",
  "new_password_confirmation": "NewP@ssw0rd"
}
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "message": "رمز عبور شما با موفقیت تغییر یافت"
  }
}
```

**Error Responses**:

- **401 Unauthorized** (Wrong current password):
```json
{
  "success": false,
  "error": {
    "message": "رمز عبور فعلی نادرست است",
    "code": "INVALID_PASSWORD"
  }
}
```

---

### 8. Logout

Logout user and revoke refresh token.

**Endpoint**: `POST /api/auth/logout`

**Headers**:
```http
Cookie: access_token=eyJhbGciOiJIUzI1NiIs...
Cookie: refresh_token=eyJhbGciOiJIUzI1NiIs...
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "message": "با موفقیت خارج شدید"
  }
}
```

Clears httpOnly cookies:
- `access_token`
- `refresh_token`

---

## Admin-Only Endpoints (Requires Admin Role)

### 9. List All Users

Get paginated list of all users (admin only).

**Endpoint**: `GET /api/users`

**Query Parameters**:
- `page` (optional, default: 1)
- `limit` (optional, default: 20, max: 100)
- `role` (optional, filter by role: "admin" or "user")
- `is_active` (optional, filter by status: "true" or "false")

**Example**: `GET /api/users?page=1&limit=20&role=user&is_active=true`

**Headers**:
```http
Cookie: access_token=eyJhbGciOiJIUzI1NiIs...
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "username": "علی_احمدی",
        "email": "ali@example.com",
        "role": "user",
        "is_active": true,
        "created_at": "2025-12-30T10:00:00Z",
        "last_login_at": "2025-12-30T10:30:00Z"
      },
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "username": "مریم_رضایی",
        "email": "maryam@example.com",
        "role": "admin",
        "is_active": true,
        "created_at": "2025-12-29T09:00:00Z",
        "last_login_at": "2025-12-30T08:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 2,
      "total_pages": 1
    }
  }
}
```

**Error Responses**:

- **403 Forbidden** (Not admin):
```json
{
  "success": false,
  "error": {
    "message": "دسترسی غیرمجاز",
    "code": "FORBIDDEN"
  }
}
```

---

### 10. Get User By ID

Get specific user details (admin only).

**Endpoint**: `GET /api/users/:id`

**Example**: `GET /api/users/550e8400-e29b-41d4-a716-446655440000`

**Headers**:
```http
Cookie: access_token=eyJhbGciOiJIUzI1NiIs...
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "علی_احمدی",
      "email": "ali@example.com",
      "role": "user",
      "is_active": true,
      "failed_login_attempts": 0,
      "locked_until": null,
      "created_at": "2025-12-30T10:00:00Z",
      "updated_at": "2025-12-30T10:00:00Z",
      "last_login_at": "2025-12-30T10:30:00Z"
    }
  }
}
```

**Error Responses**:

- **404 Not Found**:
```json
{
  "success": false,
  "error": {
    "message": "کاربر یافت نشد",
    "code": "USER_NOT_FOUND"
  }
}
```

---

### 11. Update User Role

Change user's role (admin only).

**Endpoint**: `PUT /api/users/:id/role`

**Example**: `PUT /api/users/550e8400-e29b-41d4-a716-446655440000/role`

**Headers**:
```http
Cookie: access_token=eyJhbGciOiJIUzI1NiIs...
```

**Request Body**:
```json
{
  "role": "admin"
}
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "علی_احمدی",
      "email": "ali@example.com",
      "role": "admin",
      "is_active": true,
      "updated_at": "2025-12-30T11:00:00Z"
    }
  }
}
```

**Error Responses**:

- **400 Bad Request** (Invalid role):
```json
{
  "success": false,
  "error": {
    "message": "نقش نامعتبر است. مقادیر مجاز: admin, user",
    "code": "INVALID_ROLE"
  }
}
```

---

### 12. Activate/Deactivate User

Toggle user active status (admin only).

**Endpoint**: `PUT /api/users/:id/activate`

**Example**: `PUT /api/users/550e8400-e29b-41d4-a716-446655440000/activate`

**Headers**:
```http
Cookie: access_token=eyJhbGciOiJIUzI1NiIs...
```

**Request Body**:
```json
{
  "is_active": false
}
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "علی_احمدی",
      "email": "ali@example.com",
      "role": "user",
      "is_active": false,
      "updated_at": "2025-12-30T11:00:00Z"
    }
  }
}
```

**Error Responses**:

- **400 Bad Request** (Cannot deactivate self):
```json
{
  "success": false,
  "error": {
    "message": "نمی‌توانید حساب کاربری خود را غیرفعال کنید",
    "code": "CANNOT_DEACTIVATE_SELF"
  }
}
```

- **400 Bad Request** (Cannot deactivate last admin):
```json
{
  "success": false,
  "error": {
    "message": "نمی‌توانید آخرین ادمین سیستم را غیرفعال کنید",
    "code": "CANNOT_DEACTIVATE_LAST_ADMIN"
  }
}
```

---

## Rate Limiting

All endpoints are rate-limited to prevent abuse:

| Endpoint Type | Rate Limit |
|---------------|------------|
| Public endpoints | 10 requests/minute per IP |
| Login endpoint | 5 attempts/5 minutes per IP |
| Password reset | 3 requests/hour per IP |
| Protected endpoints | 100 requests/minute per user |
| Admin endpoints | 100 requests/minute per user |

When rate limit is exceeded, returns **429 Too Many Requests**.

---

## CORS Configuration

**Development**:
- Allowed Origin: `http://localhost:5173` (Vite dev server)
- Allowed Methods: `GET, POST, PUT, DELETE`
- Allow Credentials: `true` (for cookies)

**Production** (future):
- Allowed Origin: Production frontend URL
- Same methods and credentials settings

---

## Security Headers

All responses include security headers:

```http
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
Content-Security-Policy: default-src 'self'
```

---

## Cookie Settings

All authentication cookies use these settings:

```
HttpOnly: true
Secure: true (production only, false in development)
SameSite: Strict
Path: /
```

- `access_token`: Max-Age 900 (15 minutes)
- `refresh_token`: Max-Age 604800 (7 days)

---

**Document Version**: 1.0  
**Last Updated**: 2025-12-30
