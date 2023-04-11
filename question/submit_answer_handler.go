package question

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type submitAnswer struct {
	Logger                             *zap.Logger
	QueryParticipantAndAnswerFn        QueryParticipantAndAnswerFn
	UpdateParticipantAnswerAndStatusFn UpdateParticipantAnswerAndStatusFn
}

func NewSubmitAnswer(logger *zap.Logger, queryParticipantAndAnswerFn QueryParticipantAndAnswerFn, updateParticipantAnswerAndStatusFn UpdateParticipantAnswerAndStatusFn) *submitAnswer {
	return &submitAnswer{
		Logger:                             logger,
		QueryParticipantAndAnswerFn:        queryParticipantAndAnswerFn,
		UpdateParticipantAnswerAndStatusFn: updateParticipantAnswerAndStatusFn,
	}
}

func (s *submitAnswer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := "I1o9Wp"
	// now := time.Now()

	var req SubmitAnswerRequest
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

	participantWAnswer, err := s.QueryParticipantAndAnswerFn(r.Context(), code)
	if err != nil {
		s.Logger.Error(err.Error(), zap.String("code", code))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if participantWAnswer.Answer != nil {
		err := "user's already answered"
		s.Logger.Error(err, zap.String("code", code))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err,
		})
		return
	}

	// t, _ := util.ParseDateTime(*participantWAnswer.RegisteredTime)
	// fmt.Println(now.Sub(t) > viper.GetDuration("question.timeout"), now.Sub(t))
	// if now.Sub(t) > viper.GetDuration("question.timeout") {
	// 	err := "game is timeout"
	// 	s.Logger.Error(err, zap.String("code", code))
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(map[string]string{
	// 		"error": err,
	// 	})
	// 	return
	// }

	status := "ready"
	if req.Answer == *participantWAnswer.CorrectAnswer {
		status = "used"
	}

	if err := s.UpdateParticipantAnswerAndStatusFn(r.Context(), code, req.Answer, *participantWAnswer.QuestionId, status); err != nil {
		s.Logger.Error(err.Error(), zap.String("code", code))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	s.Logger.Debug("submit answer",
		zap.String("code", code),
		zap.Int("answer", req.Answer),
		zap.Int("correct answer", *participantWAnswer.CorrectAnswer),
		zap.String("question status", status),
	)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": "success",
	})
}
