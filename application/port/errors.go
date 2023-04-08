package port

import "errors"

// ErrGameCodeNotFound is an error returned when a game code is not found in the repository
var ErrGameCodeNotFound = errors.New("game code not found")

// ErrNoAvailableQuestions is an error returned when no available questions are found in the repository
var ErrNoAvailableQuestions = errors.New("no available questions")

// ErrInvalidToken is an error returned when the jwt token is invalid
var ErrInvalidToken = errors.New("invalid token")
