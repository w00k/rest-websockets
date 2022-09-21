package middleware

import (
	"net/http"
	"strings"
	"w00k/go/rest-ws/models"
	"w00k/go/rest-ws/server"

	"github.com/golang-jwt/jwt"
)

// rutas a no validar
var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
	}
)

// shouldCheckToken: rutas que validan tokens
func shouldCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}

// CheckAuthMiddleware: validador de tokens
func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//ruta no protegida
			if !shouldCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			//obtenemos el token
			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
			//validamos con el error
			_, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})
			//si el error es distinto que nil, retornamos no autorizado
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			//en caso que todo este OK
			next.ServeHTTP(w, r)
		})
	}
}
