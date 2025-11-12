package gorm

import (
	"short-go/internal/short-links/domain/model"
	derefUtils "short-go/internal/shared/http/utils"
	"gorm.io/gorm"
)

type ShortLinkRepositoryGorm struct {
	db *gorm.DB
}

func NewShortLinkRepository(db *gorm.DB) *ShortLinkRepositoryGorm {
	return &ShortLinkRepositoryGorm{db: db}
}

func (r *ShortLinkRepositoryGorm) Create(shortLink *model.ShortLink) error {
	shortLinkModel := &ShortLinkModel{
		Code: shortLink.Code,
		OriginalURL: shortLink.OriginalURL,
		ManagementToken: &shortLink.ManagementToken,
		ExpiresAt: &shortLink.ExpiresAt,
		CreatedAt: shortLink.CreatedAt,
		UpdatedAt: shortLink.UpdatedAt,
		UserID: shortLink.UserID,
	}

	if err := r.db.Create(shortLinkModel).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShortLinkRepositoryGorm) FindByCode(code string) (*model.ShortLink, error) {
	var shortLinkModel ShortLinkModel
	if err := r.db.Where("code = ?", code).Find(&shortLinkModel).Error; err != nil {
		return nil, err
	}
	
	shortLink := &model.ShortLink{
		Code: shortLinkModel.Code,
		OriginalURL: shortLinkModel.OriginalURL,
		UserID: shortLinkModel.UserID,
		ManagementToken: derefUtils.DerefString(shortLinkModel.ManagementToken),
		ExpiresAt: derefUtils.DerefTime(shortLinkModel.ExpiresAt),
		CreatedAt: shortLinkModel.CreatedAt,
		UpdatedAt: shortLinkModel.UpdatedAt,
	}

	return shortLink, nil
}

func (r *ShortLinkRepositoryGorm) FindByManagementToken(token string) (*model.ShortLink, error) {
	var shortLinkModel ShortLinkModel
	if err := r.db.Where("management_token = ?", token).Find(&shortLinkModel).Error; err != nil {
		return nil, err
	}

	shortLink := &model.ShortLink{
		Code: shortLinkModel.Code,
		OriginalURL: shortLinkModel.OriginalURL,
		UserID: shortLinkModel.UserID,
		ManagementToken: derefUtils.DerefString(shortLinkModel.ManagementToken),
		ExpiresAt: derefUtils.DerefTime(shortLinkModel.ExpiresAt),
		CreatedAt: shortLinkModel.CreatedAt,
		UpdatedAt: shortLinkModel.UpdatedAt,
	}

	return shortLink, nil
}

func (r *ShortLinkRepositoryGorm) DeleteByCode(code string) error {
	if err := r.db.Where("code = ?", code).Delete(&ShortLinkModel{}).Error; err != nil {
		return err
	}
	return nil
}