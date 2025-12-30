package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"project-management/config"
	"project-management/models"
	"project-management/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("ایمیل یا رمز عبور نادرست است")
	ErrAccountLocked      = errors.New("حساب کاربری شما قفل شده است. لطفاً 30 دقیقه صبر کنید")
	ErrAccountDeactivated = errors.New("حساب کاربری شما غیرفعال شده است")
	ErrEmailExists        = errors.New("این ایمیل قبلاً ثبت شده است")
	ErrWeakPassword       = errors.New("رمز عبور باید حداقل 8 کاراکتر و شامل حروف بزرگ، کوچک و اعداد باشد")
	ErrInvalidToken       = errors.New("توکن نامعتبر یا منقضی شده است")
)

type AuthService interface {
	Register(ctx context.Context, req models.CreateUserRequest) (*models.User, string, string, error)
	Login(ctx context.Context, req models.LoginRequest, userAgent, ipAddress string) (*models.User, string, string, error)
	VerifyPassword(hashedPassword, password string) error
	GenerateTokens(userID uuid.UUID, role string) (accessToken, refreshToken string, err error)
	ValidateAccessToken(tokenString string) (*jwt.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	HandleFailedLogin(ctx context.Context, userID uuid.UUID) error
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	RevokeSession(ctx context.Context, refreshToken string) error
	UpdateProfile(ctx context.Context, userID uuid.UUID, req models.UpdateUserRequest) (*models.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req models.ChangePasswordRequest) error
}

type authService struct {
	userRepo          repositories.UserRepository
	sessionRepo       repositories.SessionRepository
	passwordResetRepo repositories.PasswordResetRepository
	emailService      *EmailService
}

func NewAuthService(userRepo repositories.UserRepository, sessionRepo repositories.SessionRepository, passwordResetRepo repositories.PasswordResetRepository, emailService *EmailService) AuthService {
	return &authService{
		userRepo:          userRepo,
		sessionRepo:       sessionRepo,
		passwordResetRepo: passwordResetRepo,
		emailService:      emailService,
	}
}

func (s *authService) Register(ctx context.Context, req models.CreateUserRequest) (*models.User, string, string, error) {
	// Validate input
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return nil, "", "", errors.New("نام کاربری باید بین 3 تا 50 کاراکتر باشد")
	}

	if !isValidEmail(req.Email) {
		return nil, "", "", errors.New("فرمت ایمیل نامعتبر است")
	}

	if req.Password != req.PasswordConfirmation {
		return nil, "", "", errors.New("رمز عبور و تکرار آن مطابقت ندارند")
	}

	if !isStrongPassword(req.Password) {
		return nil, "", "", ErrWeakPassword
	}

	// Check if email already exists
	req.Email = strings.ToLower(req.Email)
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, "", "", ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, "", "", err
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
		IsActive:     true,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, "", "", err
	}

	// Generate tokens
	accessToken, refreshToken, err := s.GenerateTokens(user.ID, user.Role)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest, userAgent, ipAddress string) (*models.User, string, string, error) {
	// Get user by email
	req.Email = strings.ToLower(req.Email)
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", "", ErrInvalidCredentials
	}

	// Check if account is locked
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return nil, "", "", ErrAccountLocked
	}

	// Check if account is active
	if !user.IsActive {
		return nil, "", "", ErrAccountDeactivated
	}

	// Verify password
	err = s.VerifyPassword(user.PasswordHash, req.Password)
	if err != nil {
		// Handle failed login
		s.HandleFailedLogin(ctx, user.ID)
		return nil, "", "", ErrInvalidCredentials
	}

	// Reset failed login attempts on successful login
	if user.FailedLoginAttempts > 0 {
		s.userRepo.UpdateFailedAttempts(ctx, user.ID, 0)
		s.userRepo.LockAccount(ctx, user.ID, time.Time{}) // Clear lock
	}

	// Update last login timestamp
	now := time.Now()
	user.LastLoginAt = &now
	user.UpdatedAt = now
	s.userRepo.Update(ctx, user)

	// Generate tokens
	accessToken, refreshToken, err := s.GenerateTokens(user.ID, user.Role)
	if err != nil {
		return nil, "", "", err
	}

	// Store refresh token in sessions table
	refreshTokenHash := hashToken(refreshToken)
	session := &models.Session{
		UserID:           user.ID,
		RefreshTokenHash: refreshTokenHash,
		UserAgent:        userAgent,
		IPAddress:        ipAddress,
		ExpiresAt:        time.Now().Add(config.JWTRefreshExpiry),
		Revoked:          false,
	}
	s.sessionRepo.Create(ctx, session)

	return user, accessToken, refreshToken, nil
}

func (s *authService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *authService) GenerateTokens(userID uuid.UUID, role string) (accessToken, refreshToken string, err error) {
	// Generate access token (15 minutes)
	accessClaims := jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"type":    "access",
		"exp":     time.Now().Add(config.JWTAccessExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token (7 days)
	refreshClaims := jwt.MapClaims{
		"user_id": userID.String(),
		"type":    "refresh",
		"exp":     time.Now().Add(config.JWTRefreshExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Check token type
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return nil, ErrInvalidToken
	}

	return token, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Validate refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return "", "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", ErrInvalidToken
	}

	// Check token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", "", ErrInvalidToken
	}

	// Check if session is revoked
	refreshTokenHash := hashToken(refreshToken)
	session, err := s.sessionRepo.GetByRefreshToken(ctx, refreshTokenHash)
	if err != nil || session.Revoked || session.ExpiresAt.Before(time.Now()) {
		return "", "", ErrInvalidToken
	}

	// Get user
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return "", "", ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := s.GenerateTokens(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	// Revoke old refresh token
	s.sessionRepo.Revoke(ctx, refreshTokenHash)

	return newAccessToken, newRefreshToken, nil
}

func (s *authService) HandleFailedLogin(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Increment failed login attempts
	attempts := user.FailedLoginAttempts + 1
	s.userRepo.UpdateFailedAttempts(ctx, userID, attempts)

	// Lock account after 5 failed attempts
	if attempts >= 5 {
		lockUntil := time.Now().Add(30 * time.Minute)
		s.userRepo.LockAccount(ctx, userID, lockUntil)
	}

	return nil
}

// RequestPasswordReset generates a password reset token and sends it via email
func (s *authService) RequestPasswordReset(ctx context.Context, email string) error {
	// Always return success to prevent email enumeration
	// This means we don't reveal whether the email exists or not

	email = strings.ToLower(email)
	user, err := s.userRepo.GetByEmail(ctx, email)

	// If user doesn't exist, silently succeed (security best practice)
	if err != nil || user == nil {
		return nil
	}

	// Generate cryptographically secure 32-byte token
	token, err := generateSecureToken()
	if err != nil {
		return err
	}

	// Hash token for database storage
	tokenHash := hashToken(token)

	// Create password reset token with 1-hour expiry
	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	err = s.passwordResetRepo.Create(ctx, resetToken)
	if err != nil {
		return err
	}

	// Send email with plain token (not the hash)
	err = s.emailService.SendPasswordResetEmail(user.Email, token)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to send password reset email: %v\n", err)
	}

	return nil
}

// ResetPassword validates token and updates user password
func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate new password strength
	if !isStrongPassword(newPassword) {
		return ErrWeakPassword
	}

	// Hash the received token to match database
	tokenHash := hashToken(token)

	// Lookup token
	resetToken, err := s.passwordResetRepo.GetByToken(ctx, tokenHash)
	if err != nil || resetToken == nil {
		return ErrInvalidToken
	}

	// Check if token is expired
	if time.Now().After(resetToken.ExpiresAt) {
		return ErrInvalidToken
	}

	// Check if token was already used
	if resetToken.Used {
		return ErrInvalidToken
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return err
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
	if err != nil || user == nil {
		return ErrInvalidToken
	}

	// Update user password
	user.PasswordHash = string(hashedPassword)
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	// Mark token as used
	err = s.passwordResetRepo.MarkAsUsed(ctx, tokenHash)
	if err != nil {
		return err
	}

	// Reset failed login attempts if any
	err = s.userRepo.UpdateFailedAttempts(ctx, user.ID, 0)
	if err != nil {
		return err
	}

	return nil
}

// RevokeSession revokes a session by marking it as revoked in the database
func (s *authService) RevokeSession(ctx context.Context, refreshToken string) error {
	// Hash the refresh token to match database
	tokenHash := hashToken(refreshToken)

	// Revoke the session
	return s.sessionRepo.Revoke(ctx, tokenHash)
}

// UpdateProfile updates user profile information
func (s *authService) UpdateProfile(ctx context.Context, userID uuid.UUID, req models.UpdateUserRequest) (*models.User, error) {
	// Validate input
	if req.Username == "" || req.Email == "" {
		return nil, errors.New("نام کاربری و ایمیل نمی‌توانند خالی باشند")
	}

	if !isValidEmail(req.Email) {
		return nil, errors.New("ایمیل نامعتبر است")
	}

	// Check if email is already taken by another user
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser.ID != userID {
		return nil, ErrEmailExists
	}

	// Update user
	user := &models.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	// Return updated user
	return s.userRepo.GetByID(ctx, userID)
}

// ChangePassword changes user password after verifying current password
func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, req models.ChangePasswordRequest) error {
	// Validate input
	if req.CurrentPassword == "" || req.NewPassword == "" {
		return errors.New("رمز عبور فعلی و جدید نمی‌توانند خالی باشند")
	}

	if !isStrongPassword(req.NewPassword) {
		return ErrWeakPassword
	}

	// Get current user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("کاربر یافت نشد")
	}

	// Verify current password
	if err := s.VerifyPassword(user.PasswordHash, req.CurrentPassword); err != nil {
		return errors.New("رمز عبور فعلی نادرست است")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("خطا در هش کردن رمز عبور")
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	err = s.userRepo.Update(ctx, user)
	return err
}

// Helper functions

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
