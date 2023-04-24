package question

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/susinl/coolkids-trivia-game/util"
)

type StartQuestionRequest struct {
	Name        string `json:"name" example:"name"`
	PhoneNumber string `json:"phoneNumber" example:"012-345-6789"`
}

func (req *StartQuestionRequest) validate() error {
	req.Name = strings.TrimSpace(req.Name)
	req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)

	if !util.IsLetter(req.Name) {
		return errors.Wrapf(errors.New(fmt.Sprintf("'name' must be letter only.")), util.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.Name) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'name' must be REQUIRED field but the input is '%v'.", req.Name)), util.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.PhoneNumber) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'phoneNumber' must be REQUIRED field but the input is '%v'.", req.PhoneNumber)), util.ValidateFieldError)
	}
	if !util.IsPhone(req.PhoneNumber) {
		return errors.Wrapf(errors.New(fmt.Sprintf("'phoneNumber' must be format.")), util.ValidateFieldError)
	}
	return nil
}

type StartQuestionResponse struct {
	IsAvailable bool          `json:"isAvailable" example:"true"`
	Data        *QuestionData `json:"data,omitempty"`
}

type QuestionData struct {
	QuestionText string `json:"questionText" example:"What"`
	ChoiceA      string `json:"choiceA" example:"1"`
	ChoiceB      string `json:"choiceB" example:"2"`
	ChoiceC      string `json:"choiceC" example:"3"`
	ChoiceD      string `json:"choiceD" example:"4"`
	ChoiceE      string `json:"choiceE" example:"5"`
	ChoiceF      string `json:"choiceF" example:"6"`
}

type SubmitAnswerRequest struct {
	Answer int `json:"answer" example:"1"`
}

func (req *SubmitAnswerRequest) validate() error {
	if req.Answer < 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'answer' must be integer but the input is '%v'.", req.Answer)), util.ValidateFieldError)
	}
	return nil
}
