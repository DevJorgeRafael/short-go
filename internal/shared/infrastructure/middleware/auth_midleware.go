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

// RequireAuth valida el JWT y extrae el userId
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return validateToken(m, next)
}

// OptionalAuth intenta validar el JWT si está presente y extrae el userId
// Si existe y el token es inválido, retorna un error 401
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// No envió el token, es anónimo -> continua
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Valida el token 
		validateToken(m, next).ServeHTTP(w, r)
	})
}

// Función compartida para manejar la validación del token
func validateToken(m *AuthMiddleware, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sharedhttp.ErrorResponse(w, http.StatusUnauthorized, "Token no proporcionado")
			return
		}

		// Verificar formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sharedhttp.ErrorResponse(w, http.StatusUnauthorized, "Formato de token inválido")
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			sharedhttp.ErrorResponse(w, http.StatusUnauthorized, "Token inválido o expirado")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			sharedhttp.ErrorResponse(w, http.StatusUnauthorized, "Claims inválidos")
			return
		}

		userID, ok := claims["userId"].(string)
		if !ok {
			sharedhttp.ErrorResponse(w, http.StatusUnauthorized, "userId no encontrado en token")
			return
		}

		hasSession, err := m.sessionRepo.HasActiveSession(userID)
		if err != nil || !hasSession {
			sharedhttp.ErrorResponse(w, http.StatusUnauthorized, "Sesión inválida o expirada")
			return
		}

		ctx := context.WithValue(r.Context(), sharedContext.UserIdKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}