package repository

import "short-go/internal/analytics/domain/model"

type ClickRepository interface {
	Save(click *model.Click) error
	
	// Métodos de lectura para analytics (consultas pesadas con GROUP BY)
	CountTotal(linkCode string) (int64, error)
	GetClicksByDate(linkCode string) ([]model.DailyStat, error)
	GetTopCountries(linkCode string, limit int) ([]model.CountryStat, error)
	GetTopReferrers(linkCode string, limit int) ([]model.ReferrerStat, error)

	// O un método maestro que traiga todas las estadísticas juntas
	GetLinkStats(linkCode string) (*model.LinkStats, error)
}