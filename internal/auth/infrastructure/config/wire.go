package config

import (
	"short-go/internal/auth/application/service"
	"short-go/internal/auth/infrastructure/email"
	"short-go/internal/auth/infrastructure/http/handler"
	gormRepo "short-go/internal/auth/infrastructure/persistence/gorm"
	"short-go/internal/shared/infrastructure/middleware"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type AuthModule struct {
	Handler *handler.AuthHandler
}

func NewAuthModule(db *gorm.DB, jwtSecret string, emailsApiKey string, senderEmail string) *AuthModule {
	// Repositories
	userRepo := gormRepo.NewUserRepository(db)
	sessionRepo := gormRepo.NewSessionRepository(db)

	// Services
	emailService := email.NewBrevoEmailService(emailsApiKey, senderEmail)
	authService := service.NewAuthService(userRepo, sessionRepo, jwtSecret, emailService)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)

	return &AuthModule{
		Handler: authHandler,
	}
}

// RegisterRoutes registra las rutas del módulo auth
func (m *AuthModule) RegisterRoutes(r chi.Router, authMiddleware *middleware.AuthMiddleware) {
	r.Route("/api/auth", func(r chi.Router) {
		// Rutas públicas sin autenticación
		r.Post("/register", m.Handler.Register)
		r.Post("/login", m.Handler.Login)
		r.Post("/refresh", m.Handler.RefreshToken)

		// Reset-password
		r.Post("/forgot-password", m.Handler.ForgotPassword)
		r.Post("/reset-password", m.Handler.ResetPassword)


		// Rutas protegidas (requiren JWT)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Post("/logout", m.Handler.Logout)
			r.Get("/sessions", m.Handler.GetSessions)
		})
	})
}
