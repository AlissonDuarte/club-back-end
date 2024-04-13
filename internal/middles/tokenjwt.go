package middles

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("ixczuVrYyDZSZa7uUbdRYopxhOf1sSUHfxxdz4")

func extractToken(authorizationHeader string) (string, error) {
	const bearer = "Bearer "
	if !strings.HasPrefix(authorizationHeader, bearer) {
		return "", errors.New("invalid authorization format")
	}

	token := strings.TrimPrefix(authorizationHeader, bearer)
	return token, nil
}
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenClean, err := extractToken(tokenString)
		if err != nil {

			w.WriteHeader(http.StatusForbidden)
		}

		if tokenClean == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenClean, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metodo de assinatura inesperado")
			}
			return []byte(jwtKey), nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
