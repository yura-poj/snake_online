package main

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"snake_online/snake_game"
	"time"
)

var g = &snake_game.SnakeGame{}

func main() {
	file, logger := addLogger()
	defer finishLogging(file, logger)

	g = snake_game.Init(20, 20, 3)
	r := mux.NewRouter()

	defer func() {
		if r := recover(); r != nil {
			zap.L().Panic("Panic recovered", zap.Any("panic", r))
		}
	}()

	addHandlers(r)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	go gameplay(ticker)

	zap.L().Info("Server start listening")
	http.ListenAndServe(":80", r)
}

func addHandlers(r *mux.Router) {
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/play", play)
	r.HandleFunc("/state", getGameState)
	r.HandleFunc("/new_snake", new_snake).Methods(http.MethodPost)
	r.HandleFunc("/direction/{direction}/{id}", setDirection)
	fs := http.FileServer(http.Dir("./"))
	r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", fs))

}

func gameplay(ticker *time.Ticker) {
	for range ticker.C {
		g.Step()
		g.RemoveDeadSnakes()
	}
}
