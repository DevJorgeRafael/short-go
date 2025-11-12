package gorm

import (
	"short-go/internal/analytics/domain/model"
	"short-go/internal/analytics/domain/repository"

	"gorm.io/gorm"
)

type ClickRepositoryGorm struct {
	db *gorm.DB
}

func NewClickRepository(db *gorm.DB) repository.ClickRepository {
	return &ClickRepositoryGorm{db: db}
}

func (r *ClickRepositoryGorm) Save(click *model.Click) error {
	clickModel := &ClickModel{
		LinkCode:    click.LinkCode,
		ClickedAt:   click.ClickedAt,
		CountryCode:     click.CountryCode,
		Referrer:    click.Referrer,
		IPAddress: click.IPAddress,
		UserAgent: click.UserAgent,
	}

	if err := r.db.Create(clickModel).Error; err != nil {
		return err
	}

	return nil
}

// Métodos de lectura para analytics (consultas pesadas con GROUP BY)
func (r *ClickRepositoryGorm) CountTotal(linkCode string) (int64, error) {
	var count int64
	err := r.db.Model(&ClickModel{}).
			Where("link_code = ?", linkCode).
			Count(&count).Error
	
	return count, err
}

func (r *ClickRepositoryGorm) GetClicksByDate(linkCode string) ([]model.DailyStat, error) {
	var stats []model.DailyStat

	err := r.db.Model(&ClickModel{}).
			Select("TO_CHAR(clicked_at, 'YYYY-MM-DD') as date, COUNT(*) as count").
			Where("link_code = ?", linkCode).
			Group("TO_CHAR(clicked_at, 'YYYY-MM-DD')").
			Order("date ASC").
			Limit(30).
			Scan(&stats).Error

	return stats, err
}

func (r *ClickRepositoryGorm) GetTopCountries(linkCode string, limit int) ([]model.CountryStat, error) {
	var stats []model.CountryStat

	err := r.db.Model(&ClickModel{}).
			Select("country_code, COUNT(*) as count").
			Where("link_code = ?", linkCode).
			Group("country_code").
			Order("count DESC").
			Limit(limit).
			Scan(&stats).Error

	return stats, err
}

func (r *ClickRepositoryGorm) GetTopReferrers(linkCode string, limit int) ([]model.ReferrerStat, error) {
	var stats []model.ReferrerStat

	err := r.db.Model(&ClickModel{}).
			Select("referrer, COUNT(*) as count").
			Where("link_code = ?", linkCode).
			Group("referrer").
			Order("count DESC").
			Limit(limit).
			Scan(&stats).Error

	return stats, err
}

// O un método maestro que traiga todas las estadísticas juntas
func (r *ClickRepositoryGorm) GetLinkStats(linkCode string) (*model.LinkStats, error) {
	stats := &model.LinkStats{}

	var err error

	stats.TotalClicks, err = r.CountTotal(linkCode)
	if err != nil {
		return nil, err
	}

	stats.ClicksByDate, err = r.GetClicksByDate(linkCode)
	if err != nil {
		return nil, err
	}

	stats.TopCountries, err = r.GetTopCountries(linkCode, 5)
	if err != nil {
		return nil, err
	}

	stats.TopReferrers, err = r.GetTopReferrers(linkCode, 5)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
