package infrastructure

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"short-go/config"
	analyticsService "short-go/internal/analytics/application/service"
	analyticsConfig "short-go/internal/analytics/infrastructure/config"
	analyticsGorm "short-go/internal/analytics/infrastructure/persistence/gorm"
	authConfig "short-go/internal/auth/infrastructure/config"
	gormRepo "short-go/internal/auth/infrastructure/persistence/gorm"
	qrConfig "short-go/internal/qr/infrastructure/config"
	"short-go/internal/shared/infrastructure/middleware"
	shortenerConfig "short-go/internal/short-links/infrastructure/config"
	shortLinkGormRepo "short-go/internal/short-links/infrastructure/persistence/gorm"
)

type Container struct {
	AuthModule      *authConfig.AuthModule
	AuthMiddleware  *middleware.AuthMiddleware
	ShortenerModule *shortenerConfig.ShortenerModule
	QRModule        *qrConfig.QRModule
	AnalyticsModule *analyticsConfig.AnalyticsModule
}

func NewContainer(db *gorm.DB, cfg *config.Config) *Container {
	sessionRepo := gormRepo.NewSessionRepository(db)

	// Repos
	linkRepo := shortLinkGormRepo.NewShortLinkRepository(db)
	clickRepo := analyticsGorm.NewClickRepository(db)

	// Services
	analyticsService := analyticsService.NewAnalyticsService(clickRepo, linkRepo)

	return &Container{
		AuthModule:      authConfig.NewAuthModule(db, cfg.JWTSecret, cfg.EmailsAPIKey, cfg.SenderEmail),
		AuthMiddleware:  middleware.NewAuthMiddleware(cfg.JWTSecret, sessionRepo),
		ShortenerModule: shortenerConfig.NewShortenerModule(db, cfg, analyticsService),
		QRModule:        qrConfig.NewQRModule(cfg),
		AnalyticsModule: analyticsConfig.NewAnalyticsModule(db, linkRepo),
	}
}

// RegisterRoutes registra las rutas de todos los módulos
func (c *Container) RegisterRoutes(r chi.Router) {
	c.AuthModule.RegisterRoutes(r, c.AuthMiddleware)
	c.ShortenerModule.RegisterRoutes(r, c.AuthMiddleware)
	c.QRModule.RegisterRoutes(r)
	c.AnalyticsModule.RegisterRoutes(r, c.AuthMiddleware)
}