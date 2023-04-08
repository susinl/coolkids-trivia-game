package model

import (
	"time"
)

// GameCodeStatus is an enum-like type to represent the different states of a game code
type GameCodeStatus string

const (
	GameCodeStatusReady   GameCodeStatus = "ready"
	GameCodeStatusPlaying GameCodeStatus = "playing"
	GameCodeStatusLost    GameCodeStatus = "lost"
	GameCodeStatusWon     GameCodeStatus = "won"
)

// GameCode represents a unique game code associated with a guest
type GameCode struct {
	Code           string
	Name           string
	Email          string
	PhoneNumber    string
	QuestionID     int64
	Status         GameCodeStatus
	RegisteredTime time.Time
}

// NewGameCode creates a new GameCode instance with the initial status of "ready"
func NewGameCode(code, name, email, phoneNumber string) *GameCode {
	return &GameCode{
		Code:           code,
		Name:           name,
		Email:          email,
		PhoneNumber:    phoneNumber,
		Status:         GameCodeStatusReady,
		RegisteredTime: time.Now(),
	}
}

// SetPlaying sets the GameCode status to "playing" and assigns the QuestionID
func (gc *GameCode) SetPlaying(questionID int64) {
	gc.Status = GameCodeStatusPlaying
	gc.QuestionID = questionID
}

// SetLost sets the GameCode status to "lost"
func (gc *GameCode) SetLost() {
	gc.Status = GameCodeStatusLost
}

// SetWon sets the GameCode status to "won"
func (gc *GameCode) SetWon() {
	gc.Status = GameCodeStatusWon
}
