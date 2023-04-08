package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/susinl/coolkids-trivia-game/application/usecase"
)

type AuthenticationHandler struct {
	authUsecase usecase.AuthenticationUsecase
}

func NewAuthenticationHandler(authUsecase usecase.AuthenticationUsecase) *AuthenticationHandler {
	return &AuthenticationHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthenticationHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/auth/validate", h.handleValidate).Methods("POST")
}

func (h *AuthenticationHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Implement login logic
}

func (h *AuthenticationHandler) handleValidate(w http.ResponseWriter, r *http.Request) {
	// Implement token validation logic
}
