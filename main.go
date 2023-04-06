package main

import (
	"github.com/susinl/coolkids-trivia-game/handler"
	"github.com/susinl/coolkids-trivia-game/repository"
	"github.com/susinl/coolkids-trivia-game/service"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db := repository.NewConnection("your_connection_string")
	repo := repository.NewRepository(db)
	serv := service.NewService(repo)
	h := handler.NewHandler(serv)

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/check_game_code", h.CheckGameCodeHandler).Methods(http.MethodPost)

	http.ListenAndServe(":8080", r)
}
