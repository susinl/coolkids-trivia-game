package port

import "github.com/susinl/coolkids-trivia-game/domain/model"

// GameCodeRepository is an interface that defines the methods for storing and retrieving GameCode instances
type GameCodeRepository interface {
	// FindByCode retrieves a GameCode instance by its code
	FindByCode(code string) (*model.GameCode, error)

	// Save saves a GameCode instance to the repository
	Save(gameCode *model.GameCode) error

	// CountWinners returns the number of winners stored in the repository
	CountWinners() (int, error)
}
