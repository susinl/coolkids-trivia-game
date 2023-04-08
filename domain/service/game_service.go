package service

import (
	"errors"

	"github.com/susinl/coolkids-trivia-game/domain/model"
)

var (
	ErrQuotaReached        = errors.New("quota has been reached")
	ErrInvalidGameCode     = errors.New("invalid game code")
	ErrGameCodeAlreadyUsed = errors.New("game code has already been used")
	ErrIncorrectAnswer     = errors.New("incorrect answer")
	ErrTimeout             = errors.New("time is up")
)

type GameCodeRepository interface {
	FindByCode(code string) (*model.GameCode, error)
	Save(gameCode *model.GameCode) error
	CountWinners() (int, error)
}

type QuestionRepository interface {
	FindAvailable() (*model.Question, error)
	Save(question *model.Question) error
}

type GameService struct {
	gameCodeRepo GameCodeRepository
	questionRepo QuestionRepository
	maxWinners   int
}

func NewGameService(gameCodeRepo GameCodeRepository, questionRepo QuestionRepository, maxWinners int) *GameService {
	return &GameService{
		gameCodeRepo: gameCodeRepo,
		questionRepo: questionRepo,
		maxWinners:   maxWinners,
	}
}

func (gs *GameService) StartGame(gameCode *model.GameCode) (*model.Question, error) {
	winners, err := gs.gameCodeRepo.CountWinners()
	if err != nil {
		return nil, err
	}

	if winners >= gs.maxWinners {
		return nil, ErrQuotaReached
	}

	question, err := gs.questionRepo.FindAvailable()
	if err != nil {
		return nil, err
	}

	gameCode.SetPlaying(question.ID)
	if err := gs.gameCodeRepo.Save(gameCode); err != nil {
		return nil, err
	}

	question.SetInUse()
	if err := gs.questionRepo.Save(question); err != nil {
		return nil, err
	}

	return question, nil
}

func (gs *GameService) SubmitAnswer(gameCode *model.GameCode, question *model.Question, answer int) error {
	if question.IsValidAnswer(answer) {
		winners, err := gs.gameCodeRepo.CountWinners()
		if err != nil {
			return err
		}

		if winners >= gs.maxWinners {
			return ErrQuotaReached
		}

		gameCode.SetWon()
		question.SetClaimed()
	} else {
		gameCode.SetLost()
	}

	if err := gs.gameCodeRepo.Save(gameCode); err != nil {
		return err
	}

	if err := gs.questionRepo.Save(question); err != nil {
		return err
	}

	return nil
}
