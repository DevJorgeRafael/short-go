package gorm

import "time"

type TaskStatusModel struct {
	ID          int    `gorm:"primaryKey;autoIncrement"`
	Code        string `gorm:"unique;not null"`
	Name        string `gorm:"not null"`
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (TaskStatusModel) TableName() string {
	return "task_statuses"
}

type TaskPriorityModel struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Code      string    `gorm:"unique;not null"`
	Name      string    `gorm:"not null"`
	Level     int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (TaskPriorityModel) TableName() string {
	return "task_priorities"
}

type TaskModel struct {
	ID          string `gorm:"primaryKey;type:text"`
	UserID      string `gorm:"not null;index"`
	Title       string `gorm:"not null"`
	Description string
	StatusID    int `gorm:"not null;index"`
	PriorityID  int `gorm:"not null;index"`
	StartsAt    *time.Time
	DueDate     *time.Time `gorm:"index"`
	CompletedAt *time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Relaciones (GORM cargará estos automáticamente con Preload)
	Status   TaskStatusModel   `gorm:"foreignKey:StatusID"`
	Priority TaskPriorityModel `gorm:"foreignKey:PriorityID"`
}

func (TaskModel) TableName() string {
	return "tasks"
}
