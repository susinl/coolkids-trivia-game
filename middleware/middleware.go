package middleware

import (
	"net/http"

	"github.com/susinl/coolkids-trivia-game/util"
	"go.uber.org/zap"
)

type middleware struct {
	Logger *zap.Logger
}

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

func (m *middleware) ValidateJWT()
