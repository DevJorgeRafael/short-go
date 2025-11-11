package infrastructure

import (
	authConfig "go-task-easy-list/internal/auth/infrastructure/config"
	"go-task-easy-list/internal/shared/infrastructure/middleware"
	gormRepo "go-task-easy-list/internal/auth/infrastructure/persistence/gorm"
	taskConfig "go-task-easy-list/internal/tasks/infrastructure/config"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Container struct {
	AuthModule     *authConfig.AuthModule
	AuthMiddleware *middleware.AuthMiddleware
	TaskModule *taskConfig.TaskModule
}

func NewContainer(db *gorm.DB, jwtSecret string) *Container {
	sessionRepo := gormRepo.NewSessionRepository(db)

	return &Container {
		AuthModule: authConfig.NewAuthModule(db, jwtSecret),
		AuthMiddleware: middleware.NewAuthMiddleware(jwtSecret, sessionRepo),
		TaskModule: taskConfig.NewTaskModule(db),
	}
}

// RegisterRoutes registra las rutas de todos los módulos
func (c *Container) RegisterRoutes(r chi.Router) {
	c.AuthModule.RegisterRoutes(r, c.AuthMiddleware)
	c.TaskModule.RegisterRoutes(r, c.AuthMiddleware)
}