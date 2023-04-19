package gameCode

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/susinl/coolkids-trivia-game/util"
)

type ValidateGameCodeRequest struct {
	Code string `json:"code" example:"code"`
}

func (req *ValidateGameCodeRequest) validate() error {
	if req.Code == "" {
		return errors.Wrapf(errors.New(fmt.Sprintf("'code' must be REQUIRED field but the input is '%v'.", req.Code)), util.ValidateFieldError)
	}
	return nil
}

type ValidateGameCodeResponse struct {
	JwtToken string `json:"jwtToken" example:"xxxxx.yyyyy.zzzzz"`
}

type ValidateGameCodeErrorResponse struct {
	Error string `json:"error" example:"code is invalid"`
}

type CheckStatusResponse struct {
	StatusCode int              `json:"statusCode" example:"0"`
	Data       *CheckStatusData `json:"data,omitempty"`
}

type CheckStatusData struct {
	FullName string `json:"fullname" example:"John Doe"`
	Email    string `json:"email" example:"example@gmail.com"`
	Code     string `json:"code" example:"123456"`
}
