# Authentication Research for Go Fiber + PostgreSQL + Svelte 5

**Research Date**: December 30, 2025  
**Project**: Project Management Application  
**Feature**: User Authentication (003-user-auth)

---

## 1. Password Hashing in Go

### Decision
**Use bcrypt** for password hashing (via `golang.org/x/crypto/bcrypt`)

### Rationale
1. **Simplicity**: bcrypt is straightforward to implement with excellent Go support
2. **Battle-tested**: Widely used and vetted by the security community
3. **Auto-salting**: Automatically handles salt generation
4. **Work factor adjustability**: Can increase cost factor as hardware improves
5. **Sufficient for this project**: For a project management app with moderate security requirements, bcrypt provides excellent security

### Alternatives Considered

**Argon2id** (`golang.org/x/crypto/argon2`):
- **Pros**: Winner of 2015 Password Hashing Competition, more resistant to GPU/ASIC attacks, configurable memory/time/parallelism
- **Cons**: More complex configuration (3 parameters vs 1), newer (less battle-tested), requires more memory
- **OWASP Recommendation**: Argon2id with 19 MiB memory, t=2, p=1
- **When to use**: High-security applications, compliance requirements, or if protecting against sophisticated attackers with GPU clusters

**Scrypt**:
- **Pros**: Memory-hard algorithm
- **Cons**: Less convenient API than bcrypt in Go, not as widely adopted
- **When to use**: When Argon2id unavailable and need memory-hard algorithm

### Implementation Notes

#### Standard Configuration
```go
import "golang.org/x/crypto/bcrypt"

// Recommended cost factor: 10-12
const bcryptCost = 10

// Hash password
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

// Verify password
err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
```

#### Key Points
- **Cost factor**: Use minimum of 10 (OWASP recommendation). Consider 12 for better security if server performance allows
- **Password length limit**: bcrypt has 72-byte limit - enforce max password length of 72 bytes
- **Test performance**: Hash generation should take < 1 second. Adjust cost based on server performance
- **Auto-upgrading**: Store cost factor with hash; can increase for existing users on next login
- **Salt**: Automatically generated and stored with hash
- **Hash format**: bcrypt output includes algorithm version, cost, salt, and hash

#### Security Considerations
- Never log passwords or hashes
- Hash on server-side only (never trust client-side hashing)
- Use constant-time comparison (handled by bcrypt library)
- Consider peppering for defense-in-depth (optional secret not in DB)

---

## 2. Session Management in Go Fiber

### Decision
**Use JWT tokens** with httpOnly cookies for session management

### Rationale
1. **Stateless**: No server-side session storage needed
2. **Scalable**: Works well with multiple server instances
3. **Standard**: Industry-standard authentication mechanism
4. **Mobile-friendly**: Can use bearer tokens for mobile apps
5. **Fiber-native**: Good middleware support in Fiber ecosystem

### Alternatives Considered

**Server-side sessions** (with Redis/PostgreSQL):
- **Pros**: Can revoke immediately, simpler security model, smaller cookie size
- **Cons**: Requires session storage (Redis), harder to scale, adds latency
- **When to use**: Need immediate revocation, strict security requirements, single-server deployment

**Plain cookies** (session ID only):
- **Pros**: Simple, browser-native
- **Cons**: Requires server-side storage, scaling issues
- **When to use**: Simple apps, single-server deployments

### Implementation Notes

#### JWT Strategy
**Access Token (short-lived)**: 
- Stored in httpOnly cookie
- Expires in 15-60 minutes
- Contains user ID, role, issued time

**Refresh Token (long-lived)**:
- Stored in separate httpOnly cookie
- Expires in 7-30 days
- Used to get new access tokens
- Optional: Store refresh token hash in DB for revocation

#### Cookie Configuration
```go
app.Use(func(c *fiber.Ctx) error {
    c.Cookie(&fiber.Cookie{
        Name:     "access_token",
        Value:    tokenString,
        HTTPOnly: true,        // Prevent JavaScript access
        Secure:   true,        // HTTPS only (production)
        SameSite: "Lax",      // CSRF protection
        MaxAge:   3600,        // 1 hour
        Path:     "/",
    })
    return c.Next()
})
```

#### Security Considerations
- **httpOnly**: Prevents XSS attacks from stealing tokens
- **Secure flag**: HTTPS only in production
- **SameSite**: Protection against CSRF (use "Lax" or "Strict")
- **Short expiration**: Limit token lifetime (15-60 min for access token)
- **Token rotation**: Rotate refresh tokens on use
- **Logout**: Clear cookies and optionally blacklist tokens

---

## 3. JWT Implementation in Go

### Decision
**Use `github.com/golang-jwt/jwt/v5`** (formerly dgrijalva/jwt-go)

### Rationale
1. **Most popular**: 8.8k+ stars, 72k+ dependent projects
2. **Well-maintained**: Active development and security updates
3. **Complete**: Supports all major signing algorithms
4. **Go-native**: Idiomatic Go API
5. **Production-ready**: Used by major projects

### Alternatives Considered

**`github.com/lestrrat-go/jwx`**:
- **Pros**: More feature-complete (JWK, JWS, JWE), modern API
- **Cons**: Larger library, more complex
- **When to use**: Need advanced JWT features (JWE encryption, JWK key sets)

### Implementation Notes

#### Installation
```bash
go get -u github.com/golang-jwt/jwt/v5
```

#### Standard Claims Structure
```go
type JWTClaims struct {
    UserID   uuid.UUID `json:"user_id"`
    Email    string    `json:"email"`
    Role     string    `json:"role"`      // "admin" or "user"
    jwt.RegisteredClaims
}
```

#### Token Generation
```go
import (
    "github.com/golang-jwt/jwt/v5"
    "time"
)

func GenerateToken(userID uuid.UUID, email, role string, secret []byte) (string, error) {
    claims := JWTClaims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "project-management-app",
            Subject:   userID.String(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secret)
}
```

#### Token Verification
```go
func VerifyToken(tokenString string, secret []byte) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secret, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}
```

#### Key Configuration
- **Algorithm**: Use HS256 (HMAC-SHA256) for simplicity
- **Secret key**: Minimum 32 bytes, cryptographically random
- **Store secret**: Environment variable, never commit to git
- **Key rotation**: Consider periodic key rotation strategy

#### Standard Claims to Include
1. **sub** (Subject): User ID
2. **iat** (Issued At): Token creation time
3. **exp** (Expiration): Token expiry time
4. **iss** (Issuer): Your application name
5. **nbf** (Not Before): Token not valid before this time (optional)

#### Custom Claims
- `user_id`: User UUID
- `email`: User email
- `role`: User role (admin/user)

#### Refresh Token Strategy
```go
// Longer expiration for refresh tokens
RefreshTokenExpiration = 7 * 24 * time.Hour  // 7 days

// Optional: Store refresh token hash in database for revocation
type RefreshToken struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    TokenHash string
    ExpiresAt time.Time
    CreatedAt time.Time
}
```

#### Security Best Practices
- **Algorithm verification**: Always verify `alg` header matches expected algorithm
- **Secret management**: Store JWT secret in environment variables
- **Token expiration**: Short-lived access tokens (15-60 min)
- **Refresh tokens**: Longer-lived, stored separately, rotated on use
- **Token blacklist**: For critical logout, maintain blacklist (optional)

---

## 4. Account Lockout Mechanisms

### Decision
**Use Fiber rate limiting middleware** + **custom account lockout logic**

### Rationale
1. **Two-layer defense**: IP-based + account-based protection
2. **Fiber native**: Built-in rate limiting support
3. **Simple implementation**: No external dependencies for basic protection
4. **Sufficient**: Adequate for project management app security needs

### Alternatives Considered

**Redis-based rate limiting**:
- **Pros**: Distributed, shared across servers, persistent
- **Cons**: Requires Redis infrastructure
- **When to use**: Multi-server deployment, need shared state

**fail2ban** (system-level):
- **Pros**: OS-level protection, works for all services
- **Cons**: Requires root access, configuration complexity
- **When to use**: Dedicated servers, sysadmin available

### Implementation Notes

#### Layer 1: IP-Based Rate Limiting (Fiber Middleware)
```go
import "github.com/gofiber/fiber/v2/middleware/limiter"

app.Use("/api/auth/login", limiter.New(limiter.Config{
    Max:        5,                           // 5 requests
    Expiration: 15 * time.Minute,            // per 15 minutes
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.IP()                        // Limit by IP
    },
    LimitReached: func(c *fiber.Ctx) error {
        return c.Status(429).JSON(fiber.Map{
            "error": "بیش از حد مجاز تلاش کرده‌اید. لطفاً 15 دقیقه صبر کنید",
        })
    },
}))
```

#### Layer 2: Account-Based Lockout (Database)
```sql
-- Add to users table
ALTER TABLE users ADD COLUMN failed_login_attempts INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN account_locked_until TIMESTAMP;
ALTER TABLE users ADD COLUMN last_failed_login TIMESTAMP;
```

```go
type User struct {
    // ... existing fields
    FailedLoginAttempts int       `json:"-"`
    AccountLockedUntil  *time.Time `json:"-"`
    LastFailedLogin     *time.Time `json:"-"`
}

const (
    MaxFailedAttempts = 5
    LockoutDuration   = 30 * time.Minute
)

func (s *AuthService) Login(email, password string) (*User, error) {
    user, err := s.repo.GetByEmail(email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    // Check if account is locked
    if user.AccountLockedUntil != nil && time.Now().Before(*user.AccountLockedUntil) {
        return nil, fmt.Errorf("حساب کاربری قفل شده است. لطفاً بعداً تلاش کنید")
    }

    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        // Increment failed attempts
        user.FailedLoginAttempts++
        user.LastFailedLogin = timePtr(time.Now())

        if user.FailedLoginAttempts >= MaxFailedAttempts {
            lockUntil := time.Now().Add(LockoutDuration)
            user.AccountLockedUntil = &lockUntil
        }

        s.repo.Update(user)
        return nil, ErrInvalidCredentials
    }

    // Successful login - reset counters
    user.FailedLoginAttempts = 0
    user.AccountLockedUntil = nil
    user.LastFailedLogin = nil
    s.repo.Update(user)

    return user, nil
}
```

#### Standard Practices
- **Failed attempts threshold**: 5 attempts
- **Lockout duration**: 15-30 minutes
- **Progressive delays**: Optional - increase lockout time for repeat offenses
- **Admin unlock**: Provide admin interface to unlock accounts
- **Notification**: Email user on account lockout (optional)
- **Logging**: Log all failed login attempts for monitoring

#### Additional Protection
- **CAPTCHA**: Add after 3 failed attempts (optional)
- **Email notification**: Alert user of suspicious activity
- **IP monitoring**: Track IPs with many failed attempts
- **Honeypot**: Add hidden form fields to catch bots

---

## 5. Password Reset Tokens

### Decision
**Use cryptographically secure random tokens** stored in database with expiration

### Rationale
1. **Secure**: Crypto-random tokens prevent guessing
2. **Time-limited**: Tokens expire after short period
3. **Single-use**: Tokens invalidated after use
4. **Simple**: No external dependencies

### Implementation Notes

#### Database Schema
```sql
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    used_at TIMESTAMP,
    INDEX idx_token_hash (token_hash),
    INDEX idx_expires_at (expires_at)
);
```

#### Token Generation
```go
import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
)

func GenerateResetToken() (token string, hash string, err error) {
    // Generate 32 random bytes (256 bits)
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", "", err
    }

    // Token sent to user (URL-safe base64)
    token = base64.URLEncoding.EncodeToString(bytes)

    // Hash stored in database (SHA-256)
    hashBytes := sha256.Sum256([]byte(token))
    hash = hex.EncodeToString(hashBytes[:])

    return token, hash, nil
}

func HashToken(token string) string {
    hashBytes := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hashBytes[:])
}
```

#### Token Storage and Verification
```go
type PasswordResetToken struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    TokenHash string
    ExpiresAt time.Time
    CreatedAt time.Time
    UsedAt    *time.Time
}

func (s *AuthService) RequestPasswordReset(email string) error {
    user, err := s.repo.GetByEmail(email)
    if err != nil {
        // Don't reveal if email exists
        return nil
    }

    // Generate token
    token, hash, err := GenerateResetToken()
    if err != nil {
        return err
    }

    // Store in database
    resetToken := &PasswordResetToken{
        ID:        uuid.New(),
        UserID:    user.ID,
        TokenHash: hash,
        ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiration
        CreatedAt: time.Now(),
    }
    
    if err := s.repo.SaveResetToken(resetToken); err != nil {
        return err
    }

    // Send email with reset link
    resetURL := fmt.Sprintf("https://yourapp.com/reset-password?token=%s", token)
    s.emailService.SendPasswordResetEmail(user.Email, resetURL)

    return nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
    // Hash token to match database
    tokenHash := HashToken(token)

    // Find and validate token
    resetToken, err := s.repo.GetResetToken(tokenHash)
    if err != nil || resetToken.UsedAt != nil {
        return fmt.Errorf("توکن نامعتبر یا منقضی شده است")
    }

    if time.Now().After(resetToken.ExpiresAt) {
        return fmt.Errorf("توکن منقضی شده است")
    }

    // Hash new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcryptCost)
    if err != nil {
        return err
    }

    // Update password
    if err := s.repo.UpdatePassword(resetToken.UserID, string(hashedPassword)); err != nil {
        return err
    }

    // Mark token as used
    now := time.Now()
    resetToken.UsedAt = &now
    s.repo.UpdateResetToken(resetToken)

    return nil
}
```

#### Security Best Practices
- **Token length**: Minimum 32 bytes (256 bits) of randomness
- **Token expiration**: 1 hour standard, max 24 hours
- **Single use**: Invalidate after successful password reset
- **Hashed storage**: Store SHA-256 hash, not plain token
- **Rate limiting**: Limit reset requests per email (e.g., 3 per hour)
- **No user enumeration**: Same response whether email exists or not
- **Clean up**: Periodically delete expired tokens
- **Notification**: Email user when password changed
- **Invalidate all tokens**: When password changed manually

#### Standard Expiration Times
- **Standard**: 1 hour
- **Strict security**: 15-30 minutes
- **Lenient**: Up to 24 hours
- **Recommended**: 1 hour

---

## 6. Email Sending in Go

### Decision
**Use `net/smtp` standard library** with external SMTP service (Gmail/SendGrid/Mailgun)

### Rationale
1. **Built-in**: No external dependencies for basic email
2. **Simple**: Straightforward API for simple use cases
3. **Flexible**: Works with any SMTP service
4. **Sufficient**: Adequate for authentication emails

### Alternatives Considered

**`github.com/go-mail/mail`**:
- **Pros**: More user-friendly API, better HTML support
- **Cons**: External dependency
- **When to use**: Complex email templates, attachments

**SendGrid SDK** (`github.com/sendgrid/sendgrid-go`):
- **Pros**: Native API, better tracking/analytics, templates
- **Cons**: Vendor lock-in, requires SendGrid account
- **When to use**: High volume, need analytics, using SendGrid already

**Mailgun SDK**:
- **Pros**: Reliable, good API, affordable
- **Cons**: Vendor lock-in
- **When to use**: Production app, need reliability

### Implementation Notes

#### Using net/smtp (Standard Library)
```go
import (
    "fmt"
    "net/smtp"
    "os"
)

type EmailConfig struct {
    SMTPHost     string
    SMTPPort     string
    SMTPUsername string
    SMTPPassword string
    FromEmail    string
    FromName     string
}

func LoadEmailConfig() *EmailConfig {
    return &EmailConfig{
        SMTPHost:     os.Getenv("SMTP_HOST"),     // smtp.gmail.com
        SMTPPort:     os.Getenv("SMTP_PORT"),     // 587
        SMTPUsername: os.Getenv("SMTP_USERNAME"), // your@gmail.com
        SMTPPassword: os.Getenv("SMTP_PASSWORD"), // app password
        FromEmail:    os.Getenv("FROM_EMAIL"),
        FromName:     os.Getenv("FROM_NAME"),
    }
}

func (c *EmailConfig) SendEmail(to, subject, body string) error {
    // Setup authentication
    auth := smtp.PlainAuth("", c.SMTPUsername, c.SMTPPassword, c.SMTPHost)

    // Compose message
    msg := []byte(
        fmt.Sprintf("From: %s <%s>\r\n", c.FromName, c.FromEmail) +
        fmt.Sprintf("To: %s\r\n", to) +
        fmt.Sprintf("Subject: %s\r\n", subject) +
        "MIME-Version: 1.0\r\n" +
        "Content-Type: text/html; charset=UTF-8\r\n" +
        "\r\n" +
        body,
    )

    // Send email
    addr := fmt.Sprintf("%s:%s", c.SMTPHost, c.SMTPPort)
    return smtp.SendMail(addr, auth, c.FromEmail, []string{to}, msg)
}

func (c *EmailConfig) SendPasswordResetEmail(to, resetURL string) error {
    subject := "بازنشانی رمز عبور"
    body := fmt.Sprintf(`
        <html dir="rtl">
        <body style="font-family: Tahoma, Arial, sans-serif; direction: rtl; text-align: right;">
            <h2>بازنشانی رمز عبور</h2>
            <p>برای بازنشانی رمز عبور خود، روی لینک زیر کلیک کنید:</p>
            <p><a href="%s">بازنشانی رمز عبور</a></p>
            <p>این لینک تا 1 ساعت معتبر است.</p>
            <p>اگر شما درخواست بازنشانی نکرده‌اید، این ایمیل را نادیده بگیرید.</p>
        </body>
        </html>
    `, resetURL)

    return c.SendEmail(to, subject, body)
}

func (c *EmailConfig) SendWelcomeEmail(to, name string) error {
    subject := "خوش آمدید"
    body := fmt.Sprintf(`
        <html dir="rtl">
        <body style="font-family: Tahoma, Arial, sans-serif; direction: rtl; text-align: right;">
            <h2>خوش آمدید %s!</h2>
            <p>حساب کاربری شما با موفقیت ایجاد شد.</p>
        </body>
        </html>
    `, name)

    return c.SendEmail(to, subject, body)
}
```

#### SMTP Configuration (Gmail Example)
```bash
# .env file
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password  # Generate from Google Account settings
FROM_EMAIL=your-email@gmail.com
FROM_NAME="Project Management App"
```

#### Gmail Setup
1. Enable 2FA on Google account
2. Generate App Password: Google Account → Security → App passwords
3. Use app password (not regular password) in SMTP_PASSWORD

#### SMTP Providers Comparison

**Gmail** (Free tier):
- **Pros**: Free, easy setup, reliable
- **Cons**: 500 emails/day limit, requires app password
- **Best for**: Development, small apps

**SendGrid** (Free: 100 emails/day):
- **Pros**: Professional, good deliverability, analytics
- **Cons**: Requires signup, API key management
- **Best for**: Production, medium volume

**Mailgun** (Free: 5000 emails/month):
- **Pros**: Generous free tier, reliable
- **Cons**: Requires credit card, domain verification
- **Best for**: Production, growing apps

**Amazon SES**:
- **Pros**: Very cheap, scalable
- **Cons**: Complex setup, AWS account needed
- **Best for**: High volume, AWS infrastructure

#### Best Practices
- **Async sending**: Send emails in background goroutines
- **Error handling**: Log failures, implement retry logic
- **Rate limiting**: Respect SMTP provider limits
- **Templates**: Use HTML templates for maintainability
- **Testing**: Use Mailtrap or similar for development
- **Environment variables**: Never hardcode credentials
- **SPF/DKIM**: Configure for production (improves deliverability)

---

## 7. Role-Based Access Control (RBAC) in Go Fiber

### Decision
**Use custom Fiber middleware** for role-based authorization

### Rationale
1. **Simple**: Only 2 roles (admin, user), no complex hierarchy
2. **Fiber-native**: Clean integration with Fiber routing
3. **Flexible**: Easy to extend for more roles if needed
4. **No overhead**: No external dependencies

### Implementation Notes

#### Middleware Implementation
```go
// middleware/auth.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

// RequireAuth verifies JWT token
func RequireAuth(secret []byte) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Get token from cookie
        tokenString := c.Cookies("access_token")
        if tokenString == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "احراز هویت لازم است",
            })
        }

        // Parse and verify token
        token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method")
            }
            return secret, nil
        })

        if err != nil || !token.Valid {
            return c.Status(401).JSON(fiber.Map{
                "error": "توکن نامعتبر است",
            })
        }

        // Store claims in context
        claims := token.Claims.(*JWTClaims)
        c.Locals("user_id", claims.UserID)
        c.Locals("user_email", claims.Email)
        c.Locals("user_role", claims.Role)

        return c.Next()
    }
}

// RequireRole checks if user has required role
func RequireRole(roles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userRole := c.Locals("user_role")
        if userRole == nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "احراز هویت لازم است",
            })
        }

        // Check if user has one of required roles
        role := userRole.(string)
        for _, requiredRole := range roles {
            if role == requiredRole {
                return c.Next()
            }
        }

        return c.Status(403).JSON(fiber.Map{
            "error": "شما دسترسی لازم را ندارید",
        })
    }
}

// RequireAdmin is convenience function for admin-only routes
func RequireAdmin() fiber.Handler {
    return RequireRole("admin")
}

// GetCurrentUserID helper function
func GetCurrentUserID(c *fiber.Ctx) (uuid.UUID, error) {
    userID := c.Locals("user_id")
    if userID == nil {
        return uuid.Nil, fmt.Errorf("user not authenticated")
    }
    return userID.(uuid.UUID), nil
}

// GetCurrentUserRole helper function
func GetCurrentUserRole(c *fiber.Ctx) string {
    role := c.Locals("user_role")
    if role == nil {
        return ""
    }
    return role.(string)
}
```

#### Route Protection Examples
```go
// routes/routes.go
package routes

func SetupRoutes(app *fiber.App, authSecret []byte) {
    api := app.Group("/api")

    // Public routes (no authentication)
    auth := api.Group("/auth")
    auth.Post("/register", authHandler.Register)
    auth.Post("/login", authHandler.Login)
    auth.Post("/forgot-password", authHandler.ForgotPassword)
    auth.Post("/reset-password", authHandler.ResetPassword)

    // Protected routes (authentication required)
    api.Use(middleware.RequireAuth(authSecret))

    // User routes (any authenticated user)
    api.Get("/profile", userHandler.GetProfile)
    api.Put("/profile", userHandler.UpdateProfile)
    api.Post("/logout", authHandler.Logout)

    // Projects (authenticated users)
    projects := api.Group("/projects")
    projects.Get("/", projectHandler.GetAllProjects)       // List projects
    projects.Post("/", projectHandler.CreateProject)        // Create project
    projects.Get("/:id", projectHandler.GetProject)         // View project
    projects.Put("/:id", projectHandler.UpdateProject)      // Update project
    projects.Delete("/:id", projectHandler.DeleteProject)   // Delete project

    // Admin-only routes
    admin := api.Group("/admin")
    admin.Use(middleware.RequireAdmin())
    admin.Get("/users", adminHandler.ListUsers)
    admin.Put("/users/:id/role", adminHandler.ChangeUserRole)
    admin.Delete("/users/:id", adminHandler.DeleteUser)
    admin.Post("/users/:id/unlock", adminHandler.UnlockAccount)
}
```

#### Resource-Based Authorization
```go
// Example: Only project creator or admin can delete
func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
    projectID, _ := uuid.Parse(c.Params("id"))
    userID, _ := middleware.GetCurrentUserID(c)
    userRole := middleware.GetCurrentUserRole(c)

    project, err := h.service.GetProject(c.Context(), projectID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "پروژه یافت نشد"})
    }

    // Check if user is owner or admin
    if project.CreatedBy != userID && userRole != "admin" {
        return c.Status(403).JSON(fiber.Map{
            "error": "شما مجاز به حذف این پروژه نیستید",
        })
    }

    if err := h.service.DeleteProject(c.Context(), projectID); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "خطا در حذف پروژه"})
    }

    return c.Status(204).Send(nil)
}
```

#### Database Schema
```sql
-- Add to users table
ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'user' CHECK (role IN ('admin', 'user'));
CREATE INDEX idx_users_role ON users(role);

-- Add created_by to projects table
ALTER TABLE projects ADD COLUMN created_by UUID REFERENCES users(id);
```

#### Permission Matrix

| Resource | Action | User | Admin |
|----------|--------|------|-------|
| Register | POST | ✅ | ✅ |
| Login | POST | ✅ | ✅ |
| Profile | GET/PUT | ✅ Own | ✅ All |
| Projects | GET | ✅ Own | ✅ All |
| Projects | POST | ✅ | ✅ |
| Projects | PUT | ✅ Own | ✅ All |
| Projects | DELETE | ✅ Own | ✅ All |
| Tasks | All | ✅ Project member | ✅ All |
| Users | Manage | ❌ | ✅ |

---

## 8. Svelte 5 Authentication Patterns

### Decision
**Use Svelte 5 runes** ($state, $derived, $effect) with writable stores for auth state

### Rationale
1. **Modern**: Leverages Svelte 5's new reactivity system
2. **Simple**: Runes provide cleaner state management
3. **Flexible**: Works with existing store pattern
4. **Type-safe**: Better TypeScript integration

### Implementation Notes

#### Auth Store (Svelte 5 Pattern)
```javascript
// src/stores/authStore.js
import { writable } from 'svelte/store';
import { api } from '../lib/api.js';

function createAuthStore() {
    const { subscribe, set, update } = writable({
        user: null,
        isAuthenticated: false,
        isLoading: true,
    });

    return {
        subscribe,

        async init() {
            try {
                // Check if user is authenticated (validate existing cookie)
                const user = await api.auth.me();
                set({ user, isAuthenticated: true, isLoading: false });
            } catch (error) {
                set({ user: null, isAuthenticated: false, isLoading: false });
            }
        },

        async login(email, password) {
            const user = await api.auth.login(email, password);
            set({ user, isAuthenticated: true, isLoading: false });
            return user;
        },

        async register(userData) {
            const user = await api.auth.register(userData);
            set({ user, isAuthenticated: true, isLoading: false });
            return user;
        },

        async logout() {
            await api.auth.logout();
            set({ user: null, isAuthenticated: false, isLoading: false });
        },

        async updateProfile(profileData) {
            const user = await api.auth.updateProfile(profileData);
            update(state => ({ ...state, user }));
            return user;
        },
    };
}

export const auth = createAuthStore();
```

#### API Layer (with Cookies)
```javascript
// src/lib/api.js
const API_BASE = '/api';

async function apiCall(endpoint, options = {}) {
    const response = await fetch(`${API_BASE}${endpoint}`, {
        credentials: 'include',  // Important: Include cookies
        headers: {
            'Content-Type': 'application/json',
            ...options.headers,
        },
        ...options,
    });

    if (!response.ok) {
        if (response.status === 401) {
            // Unauthorized - clear auth state
            auth.logout();
        }
        const error = await response.json().catch(() => ({ error: 'خطایی رخ داده است' }));
        throw new Error(error.error || 'خطایی رخ داده است');
    }

    if (response.status === 204) {
        return null;
    }

    return response.json();
}

export const api = {
    auth: {
        login: (email, password) => 
            apiCall('/auth/login', { 
                method: 'POST', 
                body: JSON.stringify({ email, password }) 
            }),
        register: (userData) => 
            apiCall('/auth/register', { 
                method: 'POST', 
                body: JSON.stringify(userData) 
            }),
        logout: () => 
            apiCall('/auth/logout', { method: 'POST' }),
        me: () => 
            apiCall('/auth/me'),
        updateProfile: (profileData) => 
            apiCall('/profile', { 
                method: 'PUT', 
                body: JSON.stringify(profileData) 
            }),
        forgotPassword: (email) => 
            apiCall('/auth/forgot-password', { 
                method: 'POST', 
                body: JSON.stringify({ email }) 
            }),
        resetPassword: (token, password) => 
            apiCall('/auth/reset-password', { 
                method: 'POST', 
                body: JSON.stringify({ token, password }) 
            }),
    },
    // ... existing project/task/timelog methods
};
```

#### Login Component (Svelte 5)
```svelte
<!-- src/components/Login.svelte -->
<script>
    import { auth } from '../stores/authStore';
    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();

    let email = $state('');
    let password = $state('');
    let error = $state('');
    let isLoading = $state(false);

    async function handleLogin(e) {
        e.preventDefault();
        error = '';
        isLoading = true;

        try {
            await auth.login(email, password);
            dispatch('success');
        } catch (err) {
            error = err.message;
        } finally {
            isLoading = false;
        }
    }
</script>

<form onsubmit={handleLogin} class="space-y-4">
    <h2 class="text-2xl font-bold">ورود</h2>

    {#if error}
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
            {error}
        </div>
    {/if}

    <div>
        <label for="email" class="block mb-2">ایمیل</label>
        <input
            type="email"
            id="email"
            bind:value={email}
            required
            class="w-full px-3 py-2 border rounded"
            placeholder="email@example.com"
        />
    </div>

    <div>
        <label for="password" class="block mb-2">رمز عبور</label>
        <input
            type="password"
            id="password"
            bind:value={password}
            required
            class="w-full px-3 py-2 border rounded"
        />
    </div>

    <button
        type="submit"
        disabled={isLoading}
        class="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600 disabled:opacity-50"
    >
        {isLoading ? 'در حال ورود...' : 'ورود'}
    </button>

    <div class="text-center text-sm">
        <a href="/forgot-password" class="text-blue-500 hover:underline">
            رمز عبور را فراموش کرده‌اید؟
        </a>
    </div>
</form>
```

#### Protected Route Component
```svelte
<!-- src/components/ProtectedRoute.svelte -->
<script>
    import { auth } from '../stores/authStore';
    import { onMount } from 'svelte';

    let { children, requireAdmin = false } = $props();

    const authState = $derived($auth);

    // Show loading while checking auth
    if (authState.isLoading) {
        return;
    }

    // Redirect if not authenticated
    if (!authState.isAuthenticated) {
        // In a real app, use router to redirect
        window.location.href = '/login';
        return;
    }

    // Check admin requirement
    if (requireAdmin && authState.user?.role !== 'admin') {
        window.location.href = '/unauthorized';
        return;
    }
</script>

{#if authState.isAuthenticated}
    {@render children?.()}
{:else}
    <div class="flex items-center justify-center h-screen">
        <div class="text-center">
            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
            <p class="mt-4">در حال بارگذاری...</p>
        </div>
    </div>
{/if}
```

#### App Root with Auth Initialization
```svelte
<!-- src/App.svelte -->
<script>
    import { onMount } from 'svelte';
    import { auth } from './stores/authStore';
    import Login from './components/Login.svelte';
    import Dashboard from './components/Dashboard.svelte';

    const authState = $derived($auth);

    onMount(() => {
        // Initialize auth on app load
        auth.init();
    });
</script>

{#if authState.isLoading}
    <div class="flex items-center justify-center h-screen">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
    </div>
{:else if !authState.isAuthenticated}
    <Login />
{:else}
    <Dashboard />
{/if}
```

#### Protecting Individual Components
```svelte
<!-- src/components/AdminPanel.svelte -->
<script>
    import { auth } from '../stores/authStore';

    const user = $derived($auth.user);
    const isAdmin = $derived(user?.role === 'admin');
</script>

{#if isAdmin}
    <div class="admin-panel">
        <h2>پنل مدیریت</h2>
        <!-- Admin content -->
    </div>
{:else}
    <div class="text-center py-8">
        <p>شما دسترسی لازم را ندارید</p>
    </div>
{/if}
```

#### Token Storage Best Practices
- **Use httpOnly cookies**: Prevents XSS attacks
- **Don't store tokens in localStorage**: Vulnerable to XSS
- **Automatic cookie handling**: Browser handles cookies automatically with `credentials: 'include'`
- **Token refresh**: Handle transparently in API layer

---

## 9. Persian Language Support

### Decision
**UTF-8 encoding throughout stack** - Go, PostgreSQL, and Svelte all handle Persian correctly

### Rationale
1. **Native support**: All technologies support UTF-8
2. **Already working**: Current app uses Persian successfully
3. **No special handling**: Works out of the box

### Implementation Notes

#### Backend (Go Fiber)
```go
// Go handles UTF-8 natively - no special configuration needed
type User struct {
    Name     string `json:"name"`      // Supports Persian
    Email    string `json:"email"`
    Password string `json:"password"`
}

// Validation messages in Persian
if user.Name == "" {
    return fmt.Errorf("نام الزامی است")
}
```

#### Database (PostgreSQL)
```sql
-- PostgreSQL with UTF-8 encoding (already configured)
-- Check encoding
SHOW SERVER_ENCODING;  -- Should show UTF8

-- Text columns support Persian natively
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255),        -- Supports Persian
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255)
);

-- Indexes work correctly with Persian text
CREATE INDEX idx_users_name ON users(name);

-- Searching in Persian
SELECT * FROM users WHERE name LIKE '%علی%';
```

#### Frontend (Svelte)
```svelte
<!-- HTML already set to RTL Persian -->
<!-- index.html -->
<!DOCTYPE html>
<html lang="fa" dir="rtl">
  <head>
    <meta charset="UTF-8" />
    <!-- Persian font already configured -->
    <link href="https://fonts.googleapis.com/css2?family=Vazirmatn:wght@300;400;500;600;700&display=swap" rel="stylesheet">
  </head>
</html>

<!-- CSS already set for RTL -->
<!-- app.css -->
@layer base {
  body {
    font-family: 'Vazirmatn', sans-serif;
    direction: rtl;
  }
}

<!-- Components with Persian text -->
<script>
    let errorMessage = $state('رمز عبور نادرست است');
</script>

<div>
    <label>ایمیل</label>
    <input type="email" placeholder="email@example.com" />
</div>
```

#### Email Templates
```html
<!-- Persian email template -->
<html dir="rtl">
<body style="font-family: Tahoma, Arial, sans-serif; direction: rtl; text-align: right;">
    <h2>بازنشانی رمز عبور</h2>
    <p>برای بازنشانی رمز عبور خود، روی لینک زیر کلیک کنید:</p>
</body>
</html>
```

### Special Considerations

#### None Required
✅ Go: Native UTF-8 support  
✅ PostgreSQL: UTF-8 encoding already configured  
✅ Svelte: Native UTF-8, Vazirmatn font already in use  
✅ Email: HTML emails with dir="rtl" attribute  
✅ JSON: UTF-8 by specification  

#### Validation
```go
// Password validation allowing Persian characters
func ValidatePassword(password string) error {
    if utf8.RuneCountInString(password) < 8 {
        return fmt.Errorf("رمز عبور باید حداقل 8 کاراکتر باشد")
    }
    // bcrypt handles Persian characters correctly
    return nil
}

// Name validation
func ValidateName(name string) error {
    if utf8.RuneCountInString(name) < 2 {
        return fmt.Errorf("نام باید حداقل 2 کاراکتر باشد")
    }
    return nil
}
```

#### Testing
- ✅ Already tested with Persian project/task names
- ✅ Database stores Persian correctly
- ✅ Frontend displays Persian correctly
- ✅ Form validation works with Persian

---

## 10. Security Headers

### Decision
**Use Fiber middleware** for security headers + **Helmet-style configuration**

### Rationale
1. **Built-in**: Fiber has excellent middleware support
2. **Simple**: Easy to configure essential headers
3. **Standard**: Industry-standard security headers
4. **Sufficient**: Adequate protection for web app

### Implementation Notes

#### Security Headers Middleware
```go
// middleware/security.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
)

func SecurityHeaders() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Prevent MIME type sniffing
        c.Set("X-Content-Type-Options", "nosniff")

        // Prevent clickjacking
        c.Set("X-Frame-Options", "DENY")

        // Enable XSS protection (legacy browsers)
        c.Set("X-XSS-Protection", "1; mode=block")

        // Referrer policy
        c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

        // Permissions policy
        c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

        // Content Security Policy
        c.Set("Content-Security-Policy", 
            "default-src 'self'; " +
            "script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +  // Svelte needs unsafe-eval in dev
            "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
            "font-src 'self' https://fonts.gstatic.com; " +
            "img-src 'self' data: https:; " +
            "connect-src 'self'")

        // HSTS (HTTPS only - enable in production)
        if c.Protocol() == "https" {
            c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        }

        return c.Next()
    }
}
```

#### CORS Configuration
```go
import (
    "github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupCORS(isDevelopment bool) fiber.Handler {
    if isDevelopment {
        // Development: Allow localhost
        return cors.New(cors.Config{
            AllowOrigins:     "http://localhost:5173, http://localhost:3000",
            AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
            AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
            AllowCredentials: true,
            MaxAge:           3600,
        })
    }

    // Production: Restrict to your domain
    return cors.New(cors.Config{
        AllowOrigins:     "https://yourdomain.com",
        AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowCredentials: true,
        MaxAge:           86400,
    })
}
```

#### Rate Limiting (Anti-DOS)
```go
import (
    "github.com/gofiber/fiber/v2/middleware/limiter"
)

func SetupRateLimiting() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        100,                    // 100 requests
        Expiration: 1 * time.Minute,        // per minute
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(429).JSON(fiber.Map{
                "error": "درخواست‌های بیش از حد. لطفاً بعداً تلاش کنید",
            })
        },
    })
}
```

#### Complete Setup in main.go
```go
func main() {
    app := fiber.New(fiber.Config{
        AppName: "Project Management API",
        // Disable default security headers (we'll set custom ones)
        DisableStartupMessage: false,
    })

    // Security headers
    app.Use(middleware.SecurityHeaders())

    // CORS
    isDevelopment := os.Getenv("ENV") != "production"
    app.Use(middleware.SetupCORS(isDevelopment))

    // Rate limiting
    app.Use(middleware.SetupRateLimiting())

    // Request logging
    app.Use(logger.New(logger.Config{
        Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
    }))

    // Recover from panics
    app.Use(recover.New())

    // Routes...
    routes.SetupRoutes(app, authSecret)

    app.Listen(":3000")
}
```

#### Essential Security Headers Explained

**X-Content-Type-Options: nosniff**
- Prevents MIME type sniffing
- Protects against malicious file uploads

**X-Frame-Options: DENY**
- Prevents clickjacking attacks
- Use "SAMEORIGIN" if you need iframe embedding

**X-XSS-Protection: 1; mode=block**
- Legacy XSS protection for older browsers
- Modern browsers use CSP instead

**Content-Security-Policy (CSP)**
- Prevents XSS attacks
- Controls resource loading
- **Adjust for Svelte**: May need `'unsafe-eval'` in development

**Strict-Transport-Security (HSTS)**
- Forces HTTPS
- Only enable in production with valid SSL
- `max-age=31536000` = 1 year

**Referrer-Policy**
- Controls referrer information
- `strict-origin-when-cross-origin` balances privacy and functionality

**Permissions-Policy**
- Controls browser features
- Disable unused features (geolocation, camera, etc.)

#### CSP for Svelte in Production
```go
// Stricter CSP for production
if !isDevelopment {
    c.Set("Content-Security-Policy", 
        "default-src 'self'; " +
        "script-src 'self'; " +  // Remove unsafe-eval in production
        "style-src 'self' https://fonts.googleapis.com; " +
        "font-src 'self' https://fonts.gstatic.com; " +
        "img-src 'self' data: https:; " +
        "connect-src 'self'; " +
        "frame-ancestors 'none'")  // Equivalent to X-Frame-Options: DENY
}
```

#### Testing Security Headers
```bash
# Test headers with curl
curl -I http://localhost:3000/api/projects

# Use online tools
# - securityheaders.com
# - observatory.mozilla.org
```

#### Production Checklist
- ✅ Enable HSTS with valid SSL certificate
- ✅ Tighten CSP (remove unsafe-inline/unsafe-eval)
- ✅ Configure CORS for production domain only
- ✅ Enable rate limiting
- ✅ Set up logging and monitoring
- ✅ Regular security updates
- ✅ Penetration testing

---

## Summary Table

| Topic | Decision | Key Library/Tool | Complexity |
|-------|----------|------------------|------------|
| Password Hashing | bcrypt | `golang.org/x/crypto/bcrypt` | Low |
| Session Management | JWT + httpOnly cookies | `github.com/golang-jwt/jwt/v5` | Medium |
| JWT Implementation | golang-jwt v5 | `github.com/golang-jwt/jwt/v5` | Medium |
| Account Lockout | IP rate limit + DB lockout | Fiber limiter + custom logic | Medium |
| Password Reset | Crypto tokens in DB | `crypto/rand` + PostgreSQL | Medium |
| Email Sending | net/smtp (Gmail SMTP) | `net/smtp` (stdlib) | Low |
| RBAC | Custom Fiber middleware | Custom middleware | Low |
| Svelte Auth | Runes + stores | Svelte 5 runes | Low |
| Persian Support | UTF-8 (native) | Built-in | None |
| Security Headers | Fiber middleware | Custom middleware | Low |

---

## Implementation Priorities

### Phase 1: Core Authentication (P1)
1. User model and database migrations
2. Password hashing with bcrypt
3. JWT token generation and verification
4. Login/Register endpoints
5. Auth middleware for protected routes

### Phase 2: Security Features (P1-P2)
1. Account lockout mechanism
2. Rate limiting middleware
3. Security headers
4. Email configuration (Gmail SMTP)
5. Password reset flow

### Phase 3: Frontend Integration (P2)
1. Auth store with Svelte 5 runes
2. Login/Register components
3. Protected route component
4. Profile management
5. Persian error messages

### Phase 4: RBAC & Admin (P2-P3)
1. Role-based middleware
2. Admin user creation
3. Admin panel components
4. User management endpoints

### Phase 5: Production Readiness (P3)
1. Comprehensive logging
2. Security header hardening
3. Rate limit tuning
4. Email templates refinement
5. Documentation

---

## Dependencies to Add

```bash
# Backend dependencies
go get -u github.com/golang-jwt/jwt/v5
go get -u golang.org/x/crypto/bcrypt
# Note: net/smtp is in standard library

# Fiber already installed, no additional packages needed for:
# - Rate limiting (fiber/v2/middleware/limiter)
# - CORS (fiber/v2/middleware/cors)
# - Logger (fiber/v2/middleware/logger)
```

---

## Environment Variables Needed

```bash
# .env file for development
# Database (already configured)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=project_management

# JWT Configuration
JWT_SECRET=your-super-secret-key-min-32-bytes-long-change-in-production
JWT_ACCESS_EXPIRATION=15m
JWT_REFRESH_EXPIRATION=7d

# Email Configuration (Gmail)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-gmail-app-password
FROM_EMAIL=your-email@gmail.com
FROM_NAME=Project Management App

# Application
ENV=development  # or production
APP_URL=http://localhost:5173
API_URL=http://localhost:3000

# Rate Limiting
RATE_LIMIT_MAX=100
RATE_LIMIT_WINDOW=1m

# Account Lockout
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION=30m
```

---

## Security Best Practices Summary

1. **Never log sensitive data**: Passwords, tokens, personal info
2. **Use environment variables**: All secrets in .env, never commit
3. **HTTPS in production**: Always use SSL certificates
4. **Short token expiration**: 15-60 minutes for access tokens
5. **httpOnly cookies**: Prevent XSS token theft
6. **Rate limiting**: Protect against brute force
7. **Input validation**: Validate all user input
8. **Prepared statements**: Already using (pgx parameterized queries)
9. **Security headers**: All essential headers configured
10. **Regular updates**: Keep dependencies updated
11. **Audit logs**: Log authentication events
12. **Monitor failed logins**: Alert on suspicious activity

---

## References

- [OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [OWASP Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
- [JWT Best Practices](https://datatracker.ietf.org/doc/html/rfc8725)
- [Go Fiber Documentation](https://docs.gofiber.io/)
- [golang-jwt Documentation](https://pkg.go.dev/github.com/golang-jwt/jwt/v5)
- [Svelte 5 Runes Documentation](https://svelte-5-preview.vercel.app/docs/runes)
- [bcrypt Go Package](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

---

**Next Steps**: Use this research to create implementation plan (plan.md) and detailed tasks (tasks.md)
