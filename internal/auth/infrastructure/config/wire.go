package config

import (
	"go-task-easy-list/internal/auth/application/service"
	"go-task-easy-list/internal/auth/infrastructure/http/handler"
	gormRepo "go-task-easy-list/internal/auth/infrastructure/persistence/gorm"
	"go-task-easy-list/internal/shared/infrastructure/middleware"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type AuthModule struct {
	Handler *handler.AuthHandler
}

func NewAuthModule(db *gorm.DB, jwtSecre string) *AuthModule {
	// Repositories
	userRepo := gormRepo.NewUserRepository(db)
	sessionRepo := gormRepo.NewSessionRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, sessionRepo, jwtSecre)

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

		// Rutas protegidas (requiren JWT)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Post("/logout", m.Handler.Logout)
			r.Get("/sessions", m.Handler.GetSessions)
		})
	})
}