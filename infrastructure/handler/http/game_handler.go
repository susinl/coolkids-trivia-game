package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/susinl/coolkids-trivia-game/application/usecase"
)

type GameHandler struct {
	gameUsecase usecase.GameUsecase
}

func NewGameHandler(gameUsecase usecase.GameUsecase) *GameHandler {
	return &GameHandler{
		gameUsecase: gameUsecase,
	}
}

func (h *GameHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/game/register", h.handleGameRegistration).Methods("POST")
	router.HandleFunc("/game/question", h.handleGetQuestion).Methods("GET")
	router.HandleFunc("/game/answer", h.handleAnswerQuestion).Methods("POST")
}

func (h *GameHandler) handleGameRegistration(w http.ResponseWriter, r *http.Request) {
	// Implement game registration logic
}

func (h *GameHandler) handleGetQuestion(w http.ResponseWriter, r *http.Request) {
	// Implement get question logic
}

func (h *GameHandler) handleAnswerQuestion(w http.ResponseWriter, r *http.Request) {
	// Implement answer question logic
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
