package question

import (
	"encoding/json"
	"net/http"

	"github.com/spf13/viper"
	"github.com/susinl/coolkids-trivia-game/util"
	"go.uber.org/zap"
)

type startQuestion struct {
	Logger                                   *zap.Logger
	QueryParticipantByCodeFn                 QueryParticipantByCodeFn
	QueryQuestionByStatusFn                  QueryQuestionByStatusFn
	QueryCountTotalWinnerFn                  QueryCountTotalWinnerFn
	UpdateQuestionStatusAndParticipantInfoFn UpdateQuestionStatusAndParticipantInfoFn
}

func NewStartQuestion(logger *zap.Logger, queryParticipantByCodeFn QueryParticipantByCodeFn, queryQuestionByStatusFn QueryQuestionByStatusFn, queryCountTotalWinnerFn QueryCountTotalWinnerFn, updateQuestionStatusAndParticipantInfoFn UpdateQuestionStatusAndParticipantInfoFn) http.Handler {
	return &startQuestion{
		Logger:                                   logger,
		QueryParticipantByCodeFn:                 queryParticipantByCodeFn,
		QueryQuestionByStatusFn:                  queryQuestionByStatusFn,
		QueryCountTotalWinnerFn:                  queryCountTotalWinnerFn,
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

	if participant.Name != nil {
		err := "game code's already used"
		s.Logger.Error(err, zap.String("code", code))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err,
		})
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

	byPass := totalWinner >= viper.GetInt("question.quota")
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

	var resp StartQuestionResponse
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
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": &resp,
	})
}
