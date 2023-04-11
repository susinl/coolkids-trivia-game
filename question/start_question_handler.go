package question

import (
	"encoding/json"
	"net/http"

	"github.com/susinl/coolkids-trivia-game/util"
	"go.uber.org/zap"
)

type startQuestion struct {
	Logger                                   *zap.Logger
	QueryParticipantByCodeFn                 QueryParticipantByCodeFn
	QueryQuestionByStatusFn                  QueryQuestionByStatusFn
	UpdateQuestionStatusAndParticipantInfoFn UpdateQuestionStatusAndParticipantInfoFn
}

func NewStartQuestion(logger *zap.Logger, queryParticipantByCodeFn QueryParticipantByCodeFn, queryQuestionByStatusFn QueryQuestionByStatusFn, updateQuestionStatusAndParticipantInfoFn UpdateQuestionStatusAndParticipantInfoFn) http.Handler {
	return &startQuestion{
		Logger:                                   logger,
		QueryParticipantByCodeFn:                 queryParticipantByCodeFn,
		QueryQuestionByStatusFn:                  queryQuestionByStatusFn,
		UpdateQuestionStatusAndParticipantInfoFn: updateQuestionStatusAndParticipantInfoFn,
	}
}

func (s *startQuestion) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := "I1o9Wp"

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

	if err := s.UpdateQuestionStatusAndParticipantInfoFn(r.Context(), code, req.Name, req.Email, req.PhoneNumber, *question.Id); err != nil {
		s.Logger.Error(err.Error(), zap.String("code", code))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	s.Logger.Debug("participant",
		zap.String("code", code),
		zap.String("name", *participant.Name),
		zap.String("phone", *participant.PhoneNumber),
		zap.Int("question id", *question.Id),
	)

	resp := StartQuestionResponse{
		QuestionText: *question.QuestionText,
		ChoiceA:      *question.ChoiceA,
		ChoiceB:      *question.ChoiceB,
		ChoiceC:      *question.ChoiceC,
		ChoiceD:      *question.ChoiceD,
		ChoiceE:      *question.ChoiceE,
		ChoiceF:      *question.ChoiceF,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": &resp,
	})
}
