package config

import (
	authGormModels "go-task-easy-list/internal/auth/infrastructure/persistence/gorm"
	tasksGormModels "go-task-easy-list/internal/tasks/infrastructure/persistence/gorm"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, err
	}

	db.Exec("PRAGMA foreign_keys = ON")

	// AutoMigrate: Crear tablas automáticamente
	if err := db.AutoMigrate(
		&authGormModels.UserModel{},
		&authGormModels.SessionModel{},

		&tasksGormModels.TaskStatusModel{},
		&tasksGormModels.TaskPriorityModel{},
		&tasksGormModels.TaskModel{},
	); err != nil {
		return nil, err
	}

	if err := seedTaskCatalogs(db); err != nil {
		return nil, err
	}

	return db, nil
}


func seedTaskCatalogs(db *gorm.DB) error {
	// Verificar si ya existen datos
	var count int64
	db.Model(&tasksGormModels.TaskStatusModel{}).Count(&count)
	if count > 0 {
		return nil // Ya hay datos, no hacer nada
	}

	now := time.Now()

	// Insertar estados
	statuses := []tasksGormModels.TaskStatusModel{
		{ID: 1, Code: "PENDING", Name: "Pendiente", Description: "Tarea pendiente por iniciar", CreatedAt: now},
		{ID: 2, Code: "IN_PROGRESS", Name: "En Progreso", Description: "Tarea en proceso", CreatedAt: now},
		{ID: 3, Code: "COMPLETED", Name: "Completada", Description: "Tarea finalizada", CreatedAt: now},
	}

	for _, status := range statuses {
		if err := db.Create(&status).Error; err != nil {
			return err
		}
	}

	// Insertar prioridades
	priorities := []tasksGormModels.TaskPriorityModel{
		{ID: 1, Code: "LOW", Name: "Baja", Level: 1, CreatedAt: now},
		{ID: 2, Code: "MEDIUM", Name: "Media", Level: 2, CreatedAt: now},
		{ID: 3, Code: "HIGH", Name: "Alta", Level: 3, CreatedAt: now},
	}

	for _, priority := range priorities {
		if err := db.Create(&priority).Error; err != nil {
			return err
		}
	}

	return nil
}