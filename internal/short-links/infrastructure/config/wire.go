package config

import (
	"short-go/config"
	"short-go/internal/shared/infrastructure/middleware"
	"short-go/internal/short-links/application/service"
	"short-go/internal/short-links/infrastructure/http/handler"
	gormRepo "short-go/internal/short-links/infrastructure/persistence/gorm"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type ShortenerModule struct {
	Handler *handler.ShortLinkHandler
}

func NewShortenerModule(db *gorm.DB, cfg *config.Config) *ShortenerModule {
	// Repositories
	shortLinkRepo := gormRepo.NewShortLinkRepository(db)

	// Services
	shortLinkService := service.NewShortLinkService(shortLinkRepo)

	// Handlers
	shortLinkHandler := handler.NewShortLinkHandler(shortLinkService, cfg)

	return &ShortenerModule{
		Handler: shortLinkHandler,
	}
}

// RegisterRoutes registra las rutas del módulo shortener
func (m *ShortenerModule) RegisterRoutes(r chi.Router, authMiddleware *middleware.AuthMiddleware) {
	r.Route("/api/short-links", func(r chi.Router) {
		r.Use(authMiddleware.OptionalAuth)
		r.Post("/", m.Handler.CreateShortLink)
	})
}
