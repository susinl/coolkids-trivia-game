package usecase

import (
	"github.com/susinl/coolkids-trivia-game/application/port"
	"github.com/susinl/coolkids-trivia-game/domain/model"
	"github.com/susinl/coolkids-trivia-game/domain/service"
)

type GameUsecase struct {
	gameService service.GameService
}

func NewGameUsecase(gameCodeRepo port.GameCodeRepository, questionRepo port.QuestionRepository, maxWinners int) *GameUsecase {
	gameService := service.NewGameService(gameCodeRepo, questionRepo, maxWinners)
	return &GameUsecase{
		gameService: *gameService,
	}
}

func (gu *GameUsecase) StartGame(gameCode *model.GameCode) (*model.Question, error) {
	return gu.gameService.StartGame(gameCode)
}

func (gu *GameUsecase) SubmitAnswer(gameCode *model.GameCode, question *model.Question, answer int) error {
	return gu.gameService.SubmitAnswer(gameCode, question, answer)
}
