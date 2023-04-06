package service

import (
	"github.com/susinl/coolkids-trivia-game/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CheckGameCode(gameCode string) (bool, error) {
	return s.repo.CheckGameCode(gameCode)
}
