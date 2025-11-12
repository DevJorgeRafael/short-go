package config

import (
	"short-go/internal/analytics/application/service"
	"short-go/internal/analytics/infrastructure/http/handler"
	gormAnalyticsRepo "short-go/internal/analytics/infrastructure/persistence/gorm"
	shortLinkRepo "short-go/internal/short-links/domain/repository"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type AnalyticsModule struct {
	Handler *handler.AnalyticsHandler
}

func NewAnalyticsModule(db *gorm.DB, linkRepo shortLinkRepo.ShortLinkRepository) *AnalyticsModule {
	// Repositories (ninguno por ahora)
	clickRepo := gormAnalyticsRepo.NewClickRepository(db)

	// Services
	analyticsService := service.NewAnalyticsService(clickRepo, linkRepo)

	// Handlers
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	
	return &AnalyticsModule{
		Handler: analyticsHandler,
	}
}

// RegisterRoutes registra las rutas del módulo analytics
func (m *AnalyticsModule) RegisterRoutes(r chi.Router) {
	r.Route("/api/stats", func(r chi.Router) {
		r.Get("/{code}", m.Handler.GetStats)
	})
}