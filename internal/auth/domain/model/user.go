package model

import "time"

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"` // "-" to omit in JSON responses
	Name      string `json:"name"`
	IsActive  bool   `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}