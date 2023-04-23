package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/susinl/coolkids-trivia-game/util"
	"go.uber.org/zap"
)

type middleware struct {
	Logger *zap.Logger
}

const TokenCtxKey = "token"

func NewMiddleware(logger *zap.Logger) *middleware {
	return &middleware{
		Logger: logger,
	}
}

func (m *middleware) JsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(util.AccessControl, "*")
		w.Header().Set(util.ContentType, util.ApplicationJSON)
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid jwt token",
			})
			return
		}

		authHeaderValues := strings.Split(authHeader, " ")
		if len(authHeaderValues) != 2 || authHeaderValues[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid jwt token",
			})
			return
		}

		tokenString := authHeaderValues[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Use the secret key to validate the token signature
			return []byte(viper.GetString("jwt.secret")), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid jwt token",
			})
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid jwt token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid jwt token",
			})
			return
		}

		// Set the token in the request context
		ctx := context.WithValue(r.Context(), TokenCtxKey, claims["code"].(string))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
