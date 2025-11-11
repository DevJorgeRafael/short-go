package config

import (
	authGormModels "go-task-easy-list/internal/auth/infrastructure/persistence/gorm"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, err
	}


	// AutoMigrate: Crear tablas automáticamente
	if err := db.AutoMigrate(
		&authGormModels.UserModel{},
		&authGormModels.SessionModel{},

	); err != nil {
		return nil, err
	}

	return db, nil
}


func seedTaskCatalogs(db *gorm.DB) error {
	// Aquí se agrega los datos estáticos iniciales de claves foráneas si es necesario

	return nil
}