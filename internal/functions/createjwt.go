package functions

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("ixczuVrYyDZSZa7uUbdRYopxhOf1sSUHfxxdz4")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(userID int) (string, error) {

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func UserIdFromToken(r *http.Request) (int, error) {
	// Obtenha o token JWT do cabeçalho Authorization
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return 0, fmt.Errorf("empty header")
	}

	// Verifique se o cabeçalho Authorization está no formato correto
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return 0, fmt.Errorf("error with token format")
	}

	// Parse o token JWT
	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		return jwtKey, nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to access token: %v", err)
	}

	// Verifique se o token é válido
	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Extrair o ID do usuário do token decodificado
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("token without map")
	}

	userIDFloat, _ := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("id not found")
	}
	userID := int(userIDFloat)
	return userID, nil
}
