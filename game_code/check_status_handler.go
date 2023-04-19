package gameCode

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const TokenCtxKey = "token"

type checkStatusHandler struct {
	Logger             *zap.Logger
	QueryCheckStatusFn QueryCheckStatusFn
}

func NewCheckStatusHandler(logger *zap.Logger, queryCheckStatusFn QueryCheckStatusFn) http.Handler {
	return &checkStatusHandler{
		Logger:             logger,
		QueryCheckStatusFn: queryCheckStatusFn,
	}
}

func (h *checkStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	tokenCtx := r.Context().Value(TokenCtxKey)
	if tokenCtx == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid jwt token",
		})
		return
	}

	participantAnswerCheck, err := h.QueryCheckStatusFn(r.Context(), code)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("code", code))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var resp CheckStatusResponse
	if participantAnswerCheck.QuestionId == nil {
		// Never have started the game
		resp = CheckStatusResponse{
			StatusCode: 0,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": &resp,
		})
		return
	} else if *participantAnswerCheck.Answer != *participantAnswerCheck.CorrectAnswer {
		// Lose
		resp = CheckStatusResponse{
			StatusCode: 2,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": &resp,
		})
		return
	} else {
		// Win the game
		resp = CheckStatusResponse{
			StatusCode: 1,
			Data: &CheckStatusData{
				FullName: *participantAnswerCheck.Name,
				Email:    *participantAnswerCheck.Email,
				Code:     *participantAnswerCheck.GameCode,
			},
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": &resp,
	})

}
