package code

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	token "github.com/susinl/coolkids-trivia-game/jwt"
	"github.com/susinl/coolkids-trivia-game/util"
	"go.uber.org/zap"
)

type validateCode struct {
	Logger                 *zap.Logger
	CheckReCaptchaClientFn CheckReCaptchaClientFn
	QueryValidateCodeFn    QueryValidateCodeFn
}

func NewValidateCode(logger *zap.Logger, checkReCaptchaClientFn CheckReCaptchaClientFn, queryValidateCodeFn QueryValidateCodeFn) http.Handler {
	return &validateCode{
		Logger:                 logger,
		CheckReCaptchaClientFn: checkReCaptchaClientFn,
		QueryValidateCodeFn:    queryValidateCodeFn,
	}
}

func (s *validateCode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req ValidateCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.Logger.Error(err.Error(), zap.String("code", req.Code))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	defer r.Body.Close()

	if err := req.validate(); err != nil {
		s.Logger.Error(err.Error(), zap.String("code", req.Code))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	reCaptcha := r.Header.Get(util.ReCaptchaTokenHeader)
	if reCaptcha == "" {
		err := errors.New("no reCaptcha token")
		s.Logger.Error(err.Error(), zap.String("code", req.Code))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	clientResp, err := s.CheckReCaptchaClientFn(reCaptcha)
	if err != nil {
		s.Logger.Error(err.Error(), zap.String("code", req.Code))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	if clientResp.ErrorCodes != nil {
		err := errors.New(strings.Join(clientResp.ErrorCodes, ","))
		s.Logger.Error(err.Error(), zap.String("code", req.Code))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if !clientResp.Success || clientResp.Score < viper.GetFloat64("recaptcha.minimum-score") {
		err := errors.New("you're robot.")
		s.Logger.Error(err.Error(), zap.String("code", req.Code))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	count, err := s.QueryValidateCodeFn(r.Context(), req.Code)
	if err != nil {
		s.Logger.Error(err.Error(), zap.String("code", req.Code))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "code is invalid",
		})
		return
	}

	if count > 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "code has duplicates in db",
		})
		return
	}

	token, err := token.CreateToken(req.Code)
	if err != nil {
		s.Logger.Error(err.Error(), zap.String("code", req.Code))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	s.Logger.Debug("validateCode",
		zap.String("code", req.Code),
		zap.String("token", token),
	)

	resp := ValidateCodeResponse{
		JwtToken: token,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
