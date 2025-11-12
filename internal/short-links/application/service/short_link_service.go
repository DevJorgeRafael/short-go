package service

import (
	"crypto/rand"
	"errors"
	"math/big"
	"short-go/internal/short-links/domain/model"
	"short-go/internal/short-links/domain/repository"
	"time"
)

// Errores de dominio
var (
	ErrShortLinkNotFound      = errors.New("enlace corto no encontrado")
	ErrUnauthorizedAccess     = errors.New("acceso no autorizado al enlace corto")
	ErrInvalidOriginalURL     = errors.New("URL original inválida")
	ErrManagementTokenInvalid = errors.New("token de gestión inválido")
)

type ShortLinkService struct {
	shortLinkRepo repository.ShortLinkRepository
}

func NewShortLinkService(shortLinkRepo repository.ShortLinkRepository) *ShortLinkService {
	return &ShortLinkService{shortLinkRepo: shortLinkRepo}
}

func (s *ShortLinkService) CreateShortLink(originalURL string, userID *string) (*model.ShortLink, error) {
	if originalURL == "" {
		return nil, ErrInvalidOriginalURL
	}

	codeManagement := generateRandomString(6)
	managementToken := generateRandomString(16)

	newShortLink := &model.ShortLink{
		Code:            codeManagement,
		OriginalURL:     originalURL,
		ManagementToken: managementToken,
		ExpiresAt:       time.Now().AddDate(0, 2, 0),
		UserID:          userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.shortLinkRepo.Create(newShortLink); err != nil {
		return nil, err
	}

	return newShortLink, nil
}

// ------------------------------ HELPERS -----------------------------------
func generateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        // Usa crypto/rand para obtener un número seguro
        num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        if err != nil {
            return "" // Manejar error si es crítico, o retornar string vacío
        }
        b[i] = charset[num.Int64()]
    }
    return string(b)
}
