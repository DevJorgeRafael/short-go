package handler

import (
	"encoding/json"
	format "go-task-easy-list/internal/shared/http/utils"
	"go-task-easy-list/internal/auth/application/service"
	sharedhttp "go-task-easy-list/internal/shared/http"
	sharedContext "go-task-easy-list/internal/shared/context"
	sharedValidation "go-task-easy-list/internal/shared/validation"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService *service.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   sharedValidation.NewValidator(),
	}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User         interface{} `json:"user"`
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	Message      string      `json:"message"`
}

// Register - POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, format.FormatValidationError(err))
		return
	}

	user, err := h.authService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		status := http.StatusBadRequest
		if err == service.ErrEmailExists {
			status = http.StatusConflict
		}
		sharedhttp.ErrorResponse(w, status, err.Error())
		return
	}

	sharedhttp.SuccessResponse(w, http.StatusCreated, user)
}

// Login - POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, accessToken, refreshToken, sessionRemoved, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusUnauthorized, "Credenciales inválidas")
		return
	}

	message := ""
	if sessionRemoved {
		message = "Se cerró tu sesión más antigua porque alcanzaste el límite de 3 sesiones activas."
	}

	response := AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      message,
	}

	sharedhttp.SuccessResponse(w, http.StatusOK, response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Extraer userId del contexto (establecido por el middleware)
	userID := sharedContext.GetUserID(r.Context())

	if err := h.authService.Logout(userID); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusInternalServerError, "Error al cerrar sesión")
		return
	}

	sharedhttp.SuccessResponse(w, http.StatusOK, map[string]string{"message": "Sesión cerrada exitosamente"})
}


// ---------------------------- Refresh Token ---------------------------- //
type RefreshRequest struct {
	RefreshRequest string `json:"refreshToken" validate:"required"`
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

// RefreshToken - POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, format.FormatValidationError(err))
		return
	}

	accessToken, err := h.authService.RefreshToken(req.RefreshRequest)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	response := RefreshResponse{
		AccessToken: accessToken,
	}

	sharedhttp.SuccessResponse(w, http.StatusOK, response)
}

// GetSessions - GET /api/auth/sessions
func (h *AuthHandler) GetSessions(w http.ResponseWriter, r *http.Request) {
	userID := sharedContext.GetUserID(r.Context())

	sessions, err := h.authService.GetActiveSessions(userID)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusInternalServerError, "Error al obtener sesiones")
		return
	}

	sharedhttp.SuccessResponse(w, http.StatusOK, sessions)
}