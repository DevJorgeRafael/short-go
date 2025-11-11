package gorm

import "time"

type UserModel struct {
	ID        string    `gorm:"primaryKey;type:text"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Name      string    `gorm:"not null"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (UserModel) TableName() string {
	return "users"
}

// SessionModel - Representa la tabla sessions
type SessionModel struct {
	ID           string    `gorm:"primaryKey;type:text"`
	UserID       string    `gorm:"not null;index"`
	RefreshToken string    `gorm:"not null;index"`
	ExpiresAt    time.Time `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (SessionModel) TableName() string {
	return "sessions"
}