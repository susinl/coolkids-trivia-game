package gameCode

import (
	"encoding/json"
	"net/http"

	token "github.com/susinl/coolkids-trivia-game/jwt"
	"go.uber.org/zap"
)

type validateGameCode struct {
	Logger                  *zap.Logger
	QueryValidateGameCodeFn QueryValidateGameCodeFn
}

func NewValidateGameCode(logger *zap.Logger, queryValidateGameCodeFn QueryValidateGameCodeFn) http.Handler {
	return &validateGameCode{
		Logger:                  logger,
		QueryValidateGameCodeFn: queryValidateGameCodeFn,
	}
}

func (s *validateGameCode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// code := "I1o9Wp"

	var req ValidateGameCodeRequest
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

	count, err := s.QueryValidateGameCodeFn(r.Context(), req.Code)
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

	s.Logger.Debug("validateGameCode",
		zap.String("code", req.Code),
		zap.String("token", token),
	)

	resp := ValidateGameCodeResponse{
		JwtToken: token,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
