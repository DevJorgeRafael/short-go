package middleware

import (
	"context"
	"net/http"
	"short-go/internal/auth/domain/repository"
	sharedContext "short-go/internal/shared/context"
	sharedhttp "short-go/internal/shared/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtSecret   string
	sessionRepo repository.SessionRepository
}

func NewAuthMiddleware(jwtSecret string, sessionRepo repository.SessionRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:   jwtSecret,
		sessionRepo: sessionRepo,
	}
}

// extractToken intenta obtener el token desde cookie o Authorization header
func (m *AuthMiddleware) extractToken(r *http.Request) string {
	// 1. Intentar desde cookie primero
	cookie, err := r.Cookie("accessToken")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// 2. Si no hay cookie, buscar en Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Verificar formato "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// RequireAuth valida el JWT y extrae el userId
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := m.extractToken(r)
		
		if tokenString == "" {
			sharedhttp.ErrorResponse(w, http.StatusUnauthorized, "Token no proporcionado")
			return
		}

		ctx, err := m.validateAndSetContext(r.Context(), tokenString)
		if err != nil {
			sharedhttp.ErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth intenta validar el JWT si está presente y extrae el userId
// Si existe y el token es inválido, retorna un error 401
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := m.extractToken(r)

		// No envió el token, es anónimo -> continúa
		if tokenString == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Valida el token
		ctx, err := m.validateAndSetContext(r.Context(), tokenString)
		if err != nil {
			sharedhttp.ErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateAndSetContext valida el token y retorna un contexto con el userId
func (m *AuthMiddleware) validateAndSetContext(ctx context.Context, tokenString string) (context.Context, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	userID, ok := claims["userId"].(string)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	hasSession, err := m.sessionRepo.HasActiveSession(userID)
	if err != nil || !hasSession {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return context.WithValue(ctx, sharedContext.UserIdKey, userID), nil
}