package port

import "github.com/susinl/coolkids-trivia-game/domain/model"

// QuestionRepository is an interface that defines the methods for storing and retrieving Question instances
type QuestionRepository interface {
	// FindAvailable retrieves a Question instance that is not in use and has not been claimed
	FindAvailable() (*model.Question, error)

	// Save saves a Question instance to the repository
	Save(question *model.Question) error
}
