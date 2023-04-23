package code

type ParticipantAnswer struct {
	Code          *string `json:"code" db:"code"`
	QuestionId    *int    `json:"questionId" db:"question_id"`
	Name          *string `json:"name" db:"name"`
	Answer        *int    `json:"answer" db:"answer"`
	CorrectAnswer *int    `json:"correctAnswer" db:"correct_answer"`
	PhoneNumber   *string `json:"phoneNumber" db:"phone_number"`
	// RegisteredTime *string `json:"registeredTime" db:"registered_time"`
}
