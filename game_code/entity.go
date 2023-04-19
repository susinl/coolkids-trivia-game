package gameCode

type ParticipantAnswer struct {
	GameCode      *string `json:"gameCode" db:"game_code"`
	QuestionId    *int    `json:"questionId" db:"question_id"`
	Name          *string `json:"name" db:"name"`
	Email         *string `json:"email" db:"email"`
	Answer        *int    `json:"answer" db:"answer"`
	CorrectAnswer *int    `json:"correctAnswer" db:"correct_answer"`
	// RegisteredTime *string `json:"registeredTime" db:"registered_time"`
}
