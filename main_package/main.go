package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"snake_online/snake"
	"snake_online/snake_game"
	"strconv"
	"time"
)

var g = &snake_game.SnakeGame{}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/home.html")
}

func play(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/game.html")
}

func init(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	id := g.NewSnake(host)
	return
}

func getGameState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(g)
}

func setDirection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	direction := vars["direction"]
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 || id >= len(g.Snakes) {
		http.Error(w, "invalid snake id", http.StatusBadRequest)
		return
	}
	var dir snake.Action
	switch direction {
	case "up":
		dir = snake.UP
	case "down":
		dir = snake.DOWN
	case "left":
		dir = snake.LEFT
	case "right":
		dir = snake.RIGHT
	default:
		http.Error(w, "invalid direction", http.StatusBadRequest)
		return
	}
	err = g.Snakes[id].SetDirection(dir)
	if err != nil {
		fmt.Fprintln(w, "Unsuccess, cant't move opposite way")
	} else {
		fmt.Fprintln(w, "Success")
	}
}

func main() {
	g = snake_game.Init(20, 20, 3)
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/play", play)
	r.HandleFunc("/state", getGameState)
	r.HandleFunc("/direction/{direction}/{id}", setDirection)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			g.Step()
		}
	}()

	http.ListenAndServe(":8080", r)
}
