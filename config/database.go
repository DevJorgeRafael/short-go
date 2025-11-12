package config

import (
	authGormModels "short-go/internal/auth/infrastructure/persistence/gorm"
	shortLinksGormModels "short-go/internal/short-links/infrastructure/persistence/gorm"
	analyticsGormModels "short-go/internal/analytics/infrastructure/persistence/gorm"
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

		&shortLinksGormModels.ShortLinkModel{},

		&analyticsGormModels.ClickModel{},
	); err != nil {
		return nil, err
	}

	return db, nil
}

// func seedTaskCatalogs(db *gorm.DB) error {
// 	// Aquí se agrega los datos estáticos iniciales de claves foráneas si es necesario

// 	return nil
// }
