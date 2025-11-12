package handler

import (
	"net/http"
	"short-go/internal/analytics/application/service"
	sharedContext "short-go/internal/shared/context"
	sharedhttp "short-go/internal/shared/http"
	"github.com/go-chi/chi/v5"
)

type AnalyticsHandler struct {
	service *service.AnalyticsService
}

func NewAnalyticsHandler(serivce *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: serivce}
} 

// getStats - GET /api/stats/{code}
func (h *AnalyticsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	token := r.URL.Query().Get("token")

	rawUserID := sharedContext.GetUserID(r.Context())
	var userID *string
	if rawUserID != "" {
		userID = &rawUserID
	}

	stats, err := h.service.GetStats(code, token, userID)

	if err != nil {
		if err == service.ErrUnauthorized {
			sharedhttp.ErrorResponse(w, http.StatusUnauthorized ,err.Error())
			return
		}
		if err == service.ErrLinkNotFound {
			sharedhttp.ErrorResponse(w, http.StatusNotFound ,err.Error())
			return
		}
		sharedhttp.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sharedhttp.SuccessResponse(w, http.StatusOK, stats)
}