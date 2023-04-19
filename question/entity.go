package question

type Participant struct {
	GameCode       *string `json:"gameCode" db:"game_code"`
	Name           *string `json:"name" db:"name"`
	PhoneNumber    *string `json:"phoneNumber" db:"phone_number"`
	QuestionId     *int    `json:"questionId" db:"question_id"`
	Answer         *int    `json:"answer" db:"answer"`
	RegisteredTime *string `json:"registeredTime" db:"registered_time"`
}

type Question struct {
	Id            *int    `json:"id" db:"id"`
	QuestionText  *string `json:"questionText" db:"question_text"`
	ChoiceA       *string `json:"choiceA" db:"choice_a"`
	ChoiceB       *string `json:"choiceB" db:"choice_b"`
	ChoiceC       *string `json:"choiceC" db:"choice_c"`
	ChoiceD       *string `json:"choiceD" db:"choice_d"`
	ChoiceE       *string `json:"choiceE" db:"choice_e"`
	ChoiceF       *string `json:"choiceF" db:"choice_f"`
	CorrectAnswer *int    `json:"correctAnswer" db:"correct_answer"`
	Status        *string `json:"status" db:"status"`
}

type ParticipantWAnswer struct {
	QuestionId     *int    `json:"questionId" db:"question_id"`
	Answer         *int    `json:"answer" db:"answer"`
	CorrectAnswer  *int    `json:"correctAnswer" db:"correct_answer"`
	RegisteredTime *string `json:"registeredTime" db:"registered_time"`
}
