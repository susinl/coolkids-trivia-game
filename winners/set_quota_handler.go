package winners

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type setQuotaHandler struct {
	Logger        *zap.Logger
	UpdateQuotaFn UpdateQuotaFn
}

type EmptyResponse struct{}

func NewSetQuotaHandler(logger *zap.Logger, updateQuotaFn UpdateQuotaFn) http.Handler {
	return &setQuotaHandler{
		Logger:        logger,
		UpdateQuotaFn: updateQuotaFn,
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

	err := s.UpdateQuotaFn(r.Context(), req.NewQuota)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	resp := EmptyResponse{}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
