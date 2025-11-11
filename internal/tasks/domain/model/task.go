package model

import "time"

type Task struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	StatusID    int       `json:"statusId"`
	PriorityID  int       `json:"priorityId"`
	StartsAt    time.Time `json:"startsAt,omitempty"`
	DueDate     time.Time `json:"dueDate,omitempty"`
	CompletedAt time.Time `json:"completedAt,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
