package handler

import (
	"encoding/json"
	"net/http"

	"github.com/susinl/coolkids-trivia-game/service"
)

type Handler struct {
	serv *service.Service
}

func NewHandler(serv *service.Service) *Handler {
	return &Handler{
		serv: serv,
	}
}

func (h *Handler) CheckGameCodeHandler(w http.ResponseWriter, r *http.Request) {
	var gameCodeData struct {
		GameCode string `json:"game_code"`
	}

	err := json.NewDecoder(r.Body).Decode(&gameCodeData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exists, err := h.serv.CheckGameCode(gameCodeData.GameCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Exists bool `json:"exists"`
	}{
		Exists: exists,
	})
}
