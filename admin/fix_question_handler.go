package admin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type fixQuestionHandler struct {
	Logger                 *zap.Logger
	UpdateQuestionStatusFn UpdateQuestionStatusFn
}

func NewFixQuestionHandler(logger *zap.Logger, updateQuestionStatusFn UpdateQuestionStatusFn) *fixQuestionHandler {
	return &fixQuestionHandler{
		Logger:                 logger,
		UpdateQuestionStatusFn: updateQuestionStatusFn,
	}
}

func (h *fixQuestionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	affects, err := h.UpdateQuestionStatusFn(r.Context())
	if err != nil {
		h.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	h.Logger.Debug(fmt.Sprintf("row affects: %d", affects))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rowAffected": affects,
	})
}
