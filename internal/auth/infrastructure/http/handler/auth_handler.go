package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"short-go/internal/auth/application/service"
	sharedContext "short-go/internal/shared/context"
	sharedhttp "short-go/internal/shared/http"
	format "short-go/internal/shared/http/utils"
	sharedValidation "short-go/internal/shared/validation"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService   *service.AuthService
	validator     *validator.Validate
	isDevelopment bool
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	// Detectar si estamos en desarrollo
	isDev := os.Getenv("ENVIRONMENT") != "production"
	
	return &AuthHandler{
		authService:   authService,
		validator:     sharedValidation.NewValidator(),
		isDevelopment: isDev,
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
	User    interface{} `json:"user"`
	Message string      `json:"message,omitempty"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Code        string `json:"code" validate:"required,len=6"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

// setCookie es un helper para establecer cookies con la configuración correcta
func (h *AuthHandler) setCookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   !h.isDevelopment, // Solo HTTPS en producción
		SameSite: http.SameSiteStrictMode,
	})
}

// deleteCookie elimina una cookie
func (h *AuthHandler) deleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   !h.isDevelopment,
		SameSite: http.SameSiteStrictMode,
	})
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

	// Establecer cookies
	h.setCookie(w, "accessToken", accessToken, 3600)        // 1 hora
	h.setCookie(w, "refreshToken", refreshToken, 2592000)   // 30 días

	message := ""
	if sessionRemoved {
		message = "Se cerró tu sesión más antigua porque alcanzaste el límite de 3 sesiones activas."
	}

	response := AuthResponse{
		User:    user,
		Message: message,
	}

	sharedhttp.SuccessResponse(w, http.StatusOK, response)
}

// Logout - POST /api/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID := sharedContext.GetUserID(r.Context())

	if err := h.authService.Logout(userID); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusInternalServerError, "Error al cerrar sesión")
		return
	}

	// Eliminar cookies
	h.deleteCookie(w, "accessToken")
	h.deleteCookie(w, "refreshToken")

	sharedhttp.SuccessResponse(w, http.StatusOK, map[string]string{"message": "Sesión cerrada exitosamente"})
}

// RefreshToken - POST /api/auth/refresh
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken,omitempty"`
	Message     string `json:"message,omitempty"`
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Intentar obtener refresh token desde cookie primero
	refreshToken := ""
	
	cookie, err := r.Cookie("refreshToken")
	if err == nil {
		refreshToken = cookie.Value
	} else {
		// Si no hay cookie, intentar desde el body
		var req RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sharedhttp.ErrorResponse(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		refreshToken = req.RefreshToken
	}

	if refreshToken == "" {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, "Refresh token requerido")
		return
	}

	accessToken, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Actualizar cookie de access token
	h.setCookie(w, "accessToken", accessToken, 3600)

	response := RefreshResponse{
		Message: "Token actualizado exitosamente",
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

// ForgotPassword - POST /api/auth/forgot-password
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, format.FormatValidationError(err))
		return
	}

	if err := h.authService.ForgotPassword(req.Email); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusInternalServerError, "Error al procesar la solicitud")
		return
	}

	msg := "Si el email está registrado, recibirás un código de recuperación."
	sharedhttp.SuccessResponse(w, http.StatusOK, map[string]string{"message": msg})
}

// ResetPassword - POST /api/auth/reset-password
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, format.FormatValidationError(err))
		return
	}

	if err := h.authService.ResetPassword(req.Email, req.Code, req.NewPassword); err != nil {
		status := http.StatusInternalServerError

		if err == service.ErrInvalidPassword {
			status = http.StatusBadRequest
		} else if err == service.ErrResetCodeNotFound || err == service.ErrResetCodeExpired {
			status = http.StatusUnauthorized
		}

		sharedhttp.ErrorResponse(w, status, err.Error())
		return
	}

	sharedhttp.SuccessResponse(w, http.StatusOK, map[string]string{"message": "Contraseña restablecida exitosamente"})
}