package repository

import "short-go/internal/short-links/domain/model"

type ShortLinkRepository interface {
	Create(shortLink *model.ShortLink) error
	FindByCode(code string) (*model.ShortLink, error)
	FindByManagementToken(token string) (*model.ShortLink, error)
	DeleteByCode(code string) error
}
