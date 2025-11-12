package infrastructure

import (
	"short-go/config"
	authConfig "short-go/internal/auth/infrastructure/config"
	gormRepo "short-go/internal/auth/infrastructure/persistence/gorm"
	"short-go/internal/shared/infrastructure/middleware"
	shortenerConfig "short-go/internal/short-links/infrastructure/config"

	qrConfig "short-go/internal/qr/infrastructure/config"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Container struct {
	AuthModule      *authConfig.AuthModule
	AuthMiddleware  *middleware.AuthMiddleware
	ShortenerModule *shortenerConfig.ShortenerModule
	QRModule        *qrConfig.QRModule
}

func NewContainer(db *gorm.DB, cfg *config.Config) *Container {
	sessionRepo := gormRepo.NewSessionRepository(db)

	return &Container{
		AuthModule:      authConfig.NewAuthModule(db, cfg.JWTSecret),
		AuthMiddleware:  middleware.NewAuthMiddleware(cfg.JWTSecret, sessionRepo),
		ShortenerModule: shortenerConfig.NewShortenerModule(db, cfg),
		QRModule:        qrConfig.NewQRModule(cfg),
	}
}

// RegisterRoutes registra las rutas de todos los módulos
func (c *Container) RegisterRoutes(r chi.Router) {
	c.AuthModule.RegisterRoutes(r, c.AuthMiddleware)
	c.ShortenerModule.RegisterRoutes(r, c.AuthMiddleware)
	c.QRModule.RegisterRoutes(r)
}
