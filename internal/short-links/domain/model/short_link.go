package model

import "time"

type ShortLink struct {
	Code            string     `json:"code"`
	OriginalURL     string     `json:"originalUrl"`
	ManagementToken string    `json:"managementToken,omitempty"`
	ExpiresAt       time.Time `json:"expiresAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	UserID          *string     `json:"userId,omitempty"`
}
