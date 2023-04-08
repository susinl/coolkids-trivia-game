package usecase

import (
	"github.com/susinl/coolkids-trivia-game/application/port"
)

type AuthenticationUsecase struct {
	authService port.AuthenticationService
}

func NewAuthenticationUsecase(authService port.AuthenticationService) *AuthenticationUsecase {
	return &AuthenticationUsecase{
		authService: authService,
	}
}

func (au *AuthenticationUsecase) GenerateToken(gameCode string, claims map[string]interface{}) (string, error) {
	return au.authService.GenerateToken(gameCode, claims)
}

func (au *AuthenticationUsecase) ValidateToken(token string) (string, map[string]interface{}, error) {
	return au.authService.ValidateToken(token)
}
