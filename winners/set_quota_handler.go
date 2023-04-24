package winners

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/susinl/coolkids-trivia-game/admin"
	"go.uber.org/zap"
)

type setQuotaHandler struct {
	Logger                 *zap.Logger
	UpdateQuotaFn          UpdateQuotaFn
	UpdateQuestionStatusFn admin.UpdateQuestionStatusFn
}

type EmptyResponse struct{}

func NewSetQuotaHandler(logger *zap.Logger, updateQuotaFn UpdateQuotaFn, updateQuestionStatusFn admin.UpdateQuestionStatusFn) http.Handler {
	return &setQuotaHandler{
		Logger:                 logger,
		UpdateQuotaFn:          updateQuotaFn,
		UpdateQuestionStatusFn: updateQuestionStatusFn,
	}
}

func (s *setQuotaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req UpdateQuotaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	defer r.Body.Close()

	// fmt.Println(req.NewQuota)

	if err := s.UpdateQuotaFn(r.Context(), req.NewQuota); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	affects, err := s.UpdateQuestionStatusFn(r.Context())
	if err != nil {
		s.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	s.Logger.Debug(fmt.Sprintf("row affects: %d", affects))

	resp := EmptyResponse{}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
