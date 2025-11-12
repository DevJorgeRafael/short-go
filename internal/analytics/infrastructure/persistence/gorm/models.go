package gorm

import "time"

type ClickModel struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	LinkCode string `gorm:"not null;index"`

	IPAddress   string `gorm:"type:text;size:45"`
	UserAgent   string `gorm:"type:text"`
	Referrer    string `gorm:"type:text"`
	CountryCode string `gorm:"size:2;index"`

	ClickedAt time.Time `gorm:"autoCreateTime;index"`
}

func (ClickModel) TableName() string {
	return "clicks"
}
