package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

func CreateToken(code string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["code"] = code
	tokenString, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
