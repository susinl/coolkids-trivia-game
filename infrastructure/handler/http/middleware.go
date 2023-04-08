package http

import (
	"net/http"

	"github.com/susinl/coolkids-trivia-game/application/usecase"
)

type Middleware struct {
	authUsecase usecase.AuthenticationUsecase
}

func NewMiddleware(authUsecase usecase.AuthenticationUsecase) *Middleware {
	return &Middleware{
		authUsecase: authUsecase,
	}
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement authentication middleware logic
	})
}
