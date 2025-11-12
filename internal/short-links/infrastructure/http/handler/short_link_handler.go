package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"short-go/config"
	analyticsService "short-go/internal/analytics/application/service"
	sharedContext "short-go/internal/shared/context"
	sharedhttp "short-go/internal/shared/http"
	format "short-go/internal/shared/http/utils"
	sharedValidation "short-go/internal/shared/validation"
	"short-go/internal/short-links/application/service"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type ShortLinkHandler struct {
	shortLinkService *service.ShortLinkService
	analyticsService *analyticsService.AnalyticsService
	validator        *validator.Validate
	config           *config.Config
}

func NewShortLinkHandler(
	shortLinkService *service.ShortLinkService, 
	analyticsService *analyticsService.AnalyticsService, 
	cfg *config.Config,
) *ShortLinkHandler {
	return &ShortLinkHandler{
		shortLinkService: shortLinkService,
		analyticsService: analyticsService,
		validator:        sharedValidation.NewValidator(),
		config:           cfg,
	}
}

type ShortLinkRequest struct {
	OriginalURL string `json:"originalUrl" validate:"required"`
}

type ShortLinkResponse struct {
	ShortUrl    string  `json:"shortUrl"`
	OriginalUrl string  `json:"originalUrl"`
	StatsUrl    string  `json:"statsUrl"`
	QrUrl       string  `json:"qrUrl,omitempty"`
	ExpiresAt   string  `json:"expiresAt,omitempty"`
	UserID      *string `json:"userId,omitempty"`
}

// CreateShortLink - POST /api/short-links
func (h *ShortLinkHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	rawUserID := sharedContext.GetUserID(r.Context())

	var userID *string
	if rawUserID != "" {
		userID = &rawUserID
	}

	var req ShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, format.FormatValidationError(err))
		return
	}

	shortLink, err := h.shortLinkService.CreateShortLink(req.OriginalURL, userID)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Construcción de Enlaces
	baseUrl := h.config.Domain
	if h.config.Port != "" && h.config.Domain == "http://localhost" {
		baseUrl = fmt.Sprintf("%s:%s", h.config.Domain, h.config.Port)
	}

	fullShortUrl := fmt.Sprintf("%s/%s", baseUrl, shortLink.Code)
	fullQrUrl := fmt.Sprintf("%s/api/qr/%s", baseUrl, shortLink.Code)

	// Estructura: <Base>/api/stats/<Code>?token=<Token>
	fullStatsUrl := fmt.Sprintf("%s/api/stats/%s?token=%s", baseUrl, shortLink.Code, shortLink.ManagementToken)

	resp := ShortLinkResponse{
		ShortUrl:    fullShortUrl,
		OriginalUrl: shortLink.OriginalURL,
		StatsUrl:    fullStatsUrl,
		QrUrl:       fullQrUrl,
		ExpiresAt:   shortLink.ExpiresAt.Format(time.RFC3339),
		UserID:      shortLink.UserID,
	}

	sharedhttp.SuccessResponse(w, http.StatusCreated, resp)
}

// Redirect - GET /{code}
func (h *ShortLinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	shortLink, err := h.shortLinkService.GetShortLinkByCode(code)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusNotFound, "Enlace no encontrado")
		return
	}

	fmt.Println("Redirigiendo al enlace original:", shortLink.OriginalURL)

	// Extrae los metadatos básisocs
	ip := r.RemoteAddr  // Nota: En producción, usa r.Header.Get("X-Forwarded-For")
	userAgent := r.UserAgent()
	referer := r.Referer()

	h.analyticsService.TrackClick(code, ip, userAgent, referer)

	http.Redirect(w, r, shortLink.OriginalURL, http.StatusFound)
}