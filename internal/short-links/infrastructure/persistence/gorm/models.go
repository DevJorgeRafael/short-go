package gorm

import (
	"time"

	authGormModels "short-go/internal/auth/infrastructure/persistence/gorm"
)

type ShortLinkModel struct {
	Code string `gorm:"primaryKey;size:32"`
	OriginalURL string `gorm:"not null"`

	// Clave de la lógica anónima/autenticada
	UserID *string `gorm:"index"`

	ManagementToken *string `gorm:"type:text;uniqueIndex"`

	// Fin de la lógica Clave
	ExpiresAt *time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Relación (GORM usará UserID como Foreign key)
	User authGormModels.UserModel `gorm:"foreignKey:UserID"`
}

func (ShortLinkModel) TableName() string {
	return "short_links"
}