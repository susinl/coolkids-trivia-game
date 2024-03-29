package code

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/susinl/coolkids-trivia-game/util"
)

type VerifyUserGoogleApiResponse struct {
	Success     bool     `json:"success" example:"true"`
	Score       float64  `json:"score" example:"1.0"`
	Action      string   `json:"action" example:"submit_form"`
	ChallengeTs string   `json:"challenge_ts" example:"2022-04-12T07:12:13Z"`
	Hostname    string   `json:"hostname" example:"localhost"`
	ErrorCodes  []string `json:"error-codes" example:"['timeout-or-duplicate']"`
}

type ValidateCodeRequest struct {
	Code string `json:"code" example:"code"`
}

func (req *ValidateCodeRequest) validate() error {
	if req.Code == "" {
		return errors.Wrapf(errors.New(fmt.Sprintf("'code' must be REQUIRED field but the input is '%v'.", req.Code)), util.ValidateFieldError)
	}
	return nil
}

type ValidateCodeResponse struct {
	JwtToken string `json:"jwtToken" example:"xxxxx.yyyyy.zzzzz"`
}

type ValidateCodeErrorResponse struct {
	Error string `json:"error" example:"code is invalid"`
}

type CheckStatusResponse struct {
	StatusCode int              `json:"statusCode" example:"0"`
	Data       *CheckStatusData `json:"data,omitempty"`
}

type CheckStatusData struct {
	Name        string `json:"fullname" example:"John Doe"`
	PhoneNumber string `json:"phoneNumber" example:"123456"`
	Code        string `json:"code" example:"123456"`
}

type CheckActiveResponse struct {
	ActiveCode int `json:"activeCode" example:"0"`
}
