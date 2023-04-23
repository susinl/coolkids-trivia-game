package question

import (
	"encoding/json"
	"net/http"

	"github.com/susinl/coolkids-trivia-game/util"
	"github.com/susinl/coolkids-trivia-game/winners"
	"go.uber.org/zap"
)

type startQuestion struct {
	Logger                                   *zap.Logger
	QueryParticipantByCodeFn                 QueryParticipantByCodeFn
	QueryQuestionByStatusFn                  QueryQuestionByStatusFn
	QueryCountTotalWinnerFn                  QueryCountTotalWinnerFn
	QueryGetQuotaFn                          winners.QueryGetQuotaFn
	UpdateQuestionStatusAndParticipantInfoFn UpdateQuestionStatusAndParticipantInfoFn
}

func NewStartQuestion(logger *zap.Logger, queryParticipantByCodeFn QueryParticipantByCodeFn, queryQuestionByStatusFn QueryQuestionByStatusFn, queryCountTotalWinnerFn QueryCountTotalWinnerFn, queryGetQuotaFn winners.QueryGetQuotaFn, updateQuestionStatusAndParticipantInfoFn UpdateQuestionStatusAndParticipantInfoFn) http.Handler {
	return &startQuestion{
		Logger:                                   logger,
		QueryParticipantByCodeFn:                 queryParticipantByCodeFn,
		QueryQuestionByStatusFn:                  queryQuestionByStatusFn,
		QueryCountTotalWinnerFn:                  queryCountTotalWinnerFn,
		QueryGetQuotaFn:                          queryGetQuotaFn,
		UpdateQuestionStatusAndParticipantInfoFn: updateQuestionStatusAndParticipantInfoFn,
	}
}

func (s *startQuestion) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := r.Context().Value(util.TokenCtxKey).(string)

	var req StartQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.Logger.Error(err.Error(), zap.String("code", code))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	defer r.Body.Close()

	if err := req.validate(); err != nil {
		s.Logger.Error(err.Error(), zap.String("code", code))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	participant, err := s.QueryParticipantByCodeFn(r.Context(), code)
	if err != nil {
		s.Logger.Error(err.Error(), zap.String("code", code))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var resp StartQuestionResponse
	if participant.Name != nil {
		err := "game code's already used"
		s.Logger.Error(err, zap.String("code", code))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&resp)
		return
	}

	totalWinner, err := s.QueryCountTotalWinnerFn(r.Context())
	if err != nil {
		s.Logger.Error(err.Error(), zap.String("code", code))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	question, err := s.QueryQuestionByStatusFn(r.Context(), util.ReadyStatus)
	if err != nil {
		s.Logger.Error(err.Error(), zap.String("code", code))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if question == nil {
		err := "running out of question"
		s.Logger.Error(err, zap.String("code", code))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err,
		})
		return
	}

	quota, err := s.QueryGetQuotaFn(r.Context())
	if err != nil {
		s.Logger.Error(err.Error(), zap.String("code", code))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	byPass := totalWinner >= quota
	if err := s.UpdateQuestionStatusAndParticipantInfoFn(r.Context(), code, req.Name, req.PhoneNumber, *question.Id, byPass); err != nil {
		s.Logger.Error(err.Error(), zap.String("code", code))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	s.Logger.Debug("participant",
		zap.String("code", code),
		zap.Int("question id", *question.Id),
		zap.Bool("by pass", byPass),
	)

	if !byPass {
		resp = StartQuestionResponse{
			IsAvailable: true,
			Data: &QuestionData{
				QuestionText: *question.QuestionText,
				ChoiceA:      *question.ChoiceA,
				ChoiceB:      *question.ChoiceB,
				ChoiceC:      *question.ChoiceC,
				ChoiceD:      *question.ChoiceD,
				ChoiceE:      *question.ChoiceE,
				ChoiceF:      *question.ChoiceF,
			},
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
