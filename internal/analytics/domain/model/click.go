package model

import "time"

type Click struct {
	ID          int       `json:"id"`
	LinkCode    string    `json:"linkCode"`
	IPAddress   string    `json:"ipAddress,omitempty"`
	UserAgent   string    `json:"userAgent,omitempty"`
	Referrer    string    `json:"referrer,omitempty"`
	CountryCode string    `json:"countryCode,omitempty"`
	ClickedAt   time.Time `json:"clickedAt"`
}

// Modelos adicionales para las estadisticas de un enlace
type LinkStats struct {
	TotalClicks  int64          `json:"totalClicks"`
	ClicksByDate []DailyStat    `json:"clicksByDate"`
	TopCountries []CountryStat  `json:"topCountries"`
	TopReferrers []ReferrerStat `json:"topReferrers"`
	LastClicks   []Click        `json:"lastClicks"` //ultimos 10 visitantes
}

// DailyStat agrupa clicks por fecha
type DailyStat struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type CountryStat struct {
	CountryCode string `json:"countryCode"`
	Count       int64  `json:"count"`
}

// ReferrerStat agrupa por fuetne de tráfico
type ReferrerStat struct {
	Referrer string `json:"referrer"`
	Count   int64  `json:"count"`
}
