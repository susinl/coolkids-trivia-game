package winners

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

const TokenCtxKey = "token"

type getWinnersHandler struct {
	Logger            *zap.Logger
	QueryWinnerListFn QueryWinnerListFn
}

func NewgetWinnersHandler(logger *zap.Logger, queryWinnerListFn QueryWinnerListFn) http.Handler {
	return &getWinnersHandler{
		Logger:            logger,
		QueryWinnerListFn: queryWinnerListFn,
	}
}

func (h *getWinnersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	winners, err := h.QueryWinnerListFn(r.Context())
	if err != nil {
		h.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := GetWinnersResponse{
		Data: winners.Winners,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

// type GetWinnersData struct {
// 	FullName    string `json:"fullname" example:"John Doe"`
// 	Email       string `json:"email" example:"johndoe@example.com"`
// 	PhoneNumber string `json:"phoneNumber" example:"123-456-7890"`
// 	Code        string `json:"code" example:"ABCD"`
// 	Timestamp   string `json:"timestamp" example:"2023-04-15T13:00:00Z"`
// }
