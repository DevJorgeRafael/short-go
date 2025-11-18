package config

import (
	"short-go/config"
	analyticsService "short-go/internal/analytics/application/service"
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

func NewShortenerModule(db *gorm.DB, cfg *config.Config, analyticsService *analyticsService.AnalyticsService) *ShortenerModule {
	// Repositories
	shortLinkRepo := gormRepo.NewShortLinkRepository(db)

	// Services
	shortLinkService := service.NewShortLinkService(shortLinkRepo)

	// Handlers
	shortLinkHandler := handler.NewShortLinkHandler(shortLinkService, analyticsService, cfg)

	return &ShortenerModule{
		Handler: shortLinkHandler,
	}
}

// RegisterRoutes registra las rutas del módulo shortener
func (m *ShortenerModule) RegisterRoutes(r chi.Router, authMiddleware *middleware.AuthMiddleware) {
    r.Route("/api/short-links", func(r chi.Router) {
		r.With(authMiddleware.OptionalAuth).Post("/", m.Handler.CreateShortLink)
	})

    r.Get("/{code}", m.Handler.Redirect)
}
