package winners

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type getQuotaHandler struct {
	Logger          *zap.Logger
	QueryGetQuotaFn QueryGetQuotaFn
}

func NewGetQuotaHandler(logger *zap.Logger, queryGetQuotaFn QueryGetQuotaFn) http.Handler {
	return &getQuotaHandler{
		Logger:          logger,
		QueryGetQuotaFn: queryGetQuotaFn,
	}
}

func (h *getQuotaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	quota, err := h.QueryGetQuotaFn(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	resp := GetQuotaResponse{
		Quota: quota,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
