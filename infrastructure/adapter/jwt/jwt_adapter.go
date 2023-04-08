package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/susinl/coolkids-trivia-game/application/port"
)

type JWTAdapter struct {
	secretKey string
}

func NewJWTAdapter(secretKey string) port.AuthenticationService {
	return &JWTAdapter{
		secretKey: secretKey,
	}
}

func (a *JWTAdapter) GenerateToken(id int64, email string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a *JWTAdapter) ParseToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := int64(claims["id"].(float64))
		return id, nil
	}

	return 0, port.ErrInvalidToken
}
