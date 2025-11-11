package config

import (
	"short-go/internal/shared/infrastructure/middleware"
	"short-go/internal/tasks/application/service"
	"short-go/internal/tasks/infrastructure/http/handler"
	gormRepo "short-go/internal/tasks/infrastructure/persistence/gorm"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type TaskModule struct {
	Handler *handler.TaskHandler
}

func NewTaskModule(db *gorm.DB) *TaskModule {
	// Repositories
	taskRepo := gormRepo.NewTaskRepository(db)

	// Services
	taskService := service.NewTaskService(taskRepo)

	// Handlers
	taskHandler := handler.NewTaskHandler(taskService)

	return &TaskModule{
		Handler: taskHandler,
	}
}

// RegisterRoutes registra las rutas del módulo tasks
func (m *TaskModule) RegisterRoutes(r chi.Router, authMiddleware *middleware.AuthMiddleware) {
	r.Route("/api/tasks", func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)
		r.Post("/", m.Handler.CreateTask)
		r.Get("/", m.Handler.GetTasks)
		r.Get("/{id}", m.Handler.GetTask)
		r.Put("/{id}", m.Handler.UpdateTask)
		r.Delete("/{id}", m.Handler.DeleteTask)
	})
}
