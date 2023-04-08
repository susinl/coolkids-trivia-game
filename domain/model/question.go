package model

// QuestionStatus is an enum-like type to represent the different states of a question
type QuestionStatus string

const (
	QuestionStatusReady   QuestionStatus = "ready"
	QuestionStatusInUse   QuestionStatus = "in_use"
	QuestionStatusClaimed QuestionStatus = "claimed"
)

// Question represents a trivia question with its choices and correct answer
type Question struct {
	ID            int64
	QuestionText  string
	Choices       [6]string
	CorrectAnswer int
	Status        QuestionStatus
}

// NewQuestion creates a new Question instance with the initial status of "ready"
func NewQuestion(id int64, questionText string, choices [6]string, correctAnswer int) *Question {
	return &Question{
		ID:            id,
		QuestionText:  questionText,
		Choices:       choices,
		CorrectAnswer: correctAnswer,
		Status:        QuestionStatusReady,
	}
}

// SetInUse sets the Question status to "in_use"
func (q *Question) SetInUse() {
	q.Status = QuestionStatusInUse
}

// SetClaimed sets the Question status to "claimed"
func (q *Question) SetClaimed() {
	q.Status = QuestionStatusClaimed
}

// IsValidAnswer checks if the given answer is correct for the question
func (q *Question) IsValidAnswer(answer int) bool {
	return q.CorrectAnswer == answer
}
