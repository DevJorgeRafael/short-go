package service

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"short-go/internal/auth/domain/model"
	"short-go/internal/auth/domain/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Errores del dominio
var (
	ErrInvalidEmail       = errors.New("email inválido")
	ErrEmailExists        = errors.New("el email ya está registrado")
	ErrInvalidPassword    = errors.New("la contraseña debe tener al menos 8 caracteres")
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrInvalidCredentials = errors.New("credenciales inválidas")

	ErrResetCodeNotFound = errors.New("código de reseteo inválido o email incorrecto")
	ErrResetCodeExpired  = errors.New("el código de reseteo ha expirado")
)

type AuthService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	jwtSecret   string

	emailService EmailService
}

func NewAuthService(
	userRepo repository.UserRepository, 
	sessionRepo repository.SessionRepository, 
	jwtSecret string,
	emailService EmailService,
	) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
		emailService: emailService,
	}
}

// Register- Registrar un usuario
func (s *AuthService) Register(email, password, name string) (*model.User, error) {
	// 1. Validar email (formato, no duplicado)
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	existing, _ := s.userRepo.FindByEmail(email)
	if existing != nil {
		return nil, ErrEmailExists
	}

	// 2. Validar y Hash password con bcrypt
	if len(password) < 8 {
		return nil, ErrInvalidPassword
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Crear user
	user := &model.User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(hashedPassword),
		Name:      name,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 4. Guardar en DB con userRepo.Create()
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 5. Retornar usuario (SIN password)
	user.Password = ""
	return user, nil
}

// Login - Iniciar sesión
func (s *AuthService) Login(email, password string) (*model.User, string, string, bool, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return nil, "", "", false, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, "", "", false, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", "", false, ErrInvalidCredentials
	}

	s.cleanExpiredSessions(user.ID)

	activeSessions, _ := s.sessionRepo.CountByUserID(user.ID)
	sessionRemoved := false
	if activeSessions >= 3 {
		s.sessionRepo.DeleteOldestByUserID(user.ID)
		sessionRemoved = true
	}

	accessToken, err := s.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, "", "", false, err
	}

	refreshToken := uuid.New().String()

	// 6. Guardar sesión
	session := &model.Session{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now(),
	}

	if err = s.sessionRepo.Create(session); err != nil {
		return nil, "", "", false, err
	}

	userResponse := &model.User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userResponse, accessToken, refreshToken, sessionRemoved, nil
}

func (s *AuthService) Logout(userID string) error {
	return s.sessionRepo.DeleteByUserID(userID)
}

func (s *AuthService) RefreshToken(refreshToken string) (newAccessToken string, err error) {
	session, err := s.sessionRepo.FindByRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("refresh token inválido")
	}

	if session.IsExpired() {
		s.sessionRepo.DeleteByUserID(session.UserID)
		return "", errors.New("refresh token expirado")
	}

	user, err := s.userRepo.FindByID(session.UserID)
	if err != nil || user == nil {
		return "", errors.New("usuario no encontrado")
	}

	newAccessToken, err = s.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func (s *AuthService) GetActiveSessions(userId string) ([]*model.Session, error) {
	return s.sessionRepo.FindActiveByUserID(userId)
}

func (s *AuthService) ForgotPassword(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return nil
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al hashear el código de reseteo: %w", err)
	}

	tokenString := string(hashedCode)
	expiresAt := time.Now().Add(10 * time.Minute)

	user.ResetPasswordToken = &tokenString
	user.ResetPasswordExpiresAt = &expiresAt

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("error al actualizar el usuario con el código de reseteo: %w", err)
	}

	if err := s.emailService.SendPasswordResetCode(user.Email, code); err != nil {
		return fmt.Errorf("error al enviar el email de reseteo: %w", err)
	}

	return nil
}

func (s *AuthService) ResetPassword(email, code, newPassword string) error {
	if len(newPassword) < 8 {
		return ErrInvalidPassword
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return ErrResetCodeNotFound
	}

	if user.ResetPasswordToken == nil || user.ResetPasswordExpiresAt == nil {
		return ErrResetCodeNotFound
	}

	if time.Now().After(*user.ResetPasswordExpiresAt) {
		user.ResetPasswordToken = nil
		user.ResetPasswordExpiresAt = nil
		s.userRepo.Update(user)
		return ErrResetCodeExpired
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.ResetPasswordToken), []byte(code)); err != nil {
		return ErrResetCodeNotFound
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al hashear la nueva contraseña: %w", err)
	}

	user.Password = string(hashedPassword)
	user.ResetPasswordToken = nil
	user.ResetPasswordExpiresAt = nil

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("error al actualizar la contraseña del usuario: %w", err)
	}

	s.sessionRepo.DeleteByUserID(user.ID)

	return nil
}

// --------------------- Helpers ---------------------
func (s *AuthService) generateAccessToken(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"email":  email,
		"exp":    time.Now().Add(1 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// Elimina sesiones expiradas del usuario
func (s *AuthService) cleanExpiredSessions(userID string) {
	s.sessionRepo.DeleteExpiredByUserID(userID)
}
