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

func new_snake(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	if findSnakeIp(host) != nil {
		http.Error(w, "Snake has already been created", http.StatusForbidden)
		return
	}
	playerID := g.NewSnake(host)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(map[string]int{"id": playerID})
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
	if err != nil {
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
	s := findSnake(id)
	if s == nil {
		http.Error(w, "invalid snake id", http.StatusBadRequest)
	} else {
		err = s.SetDirection(dir)
		if err != nil {
			fmt.Fprintln(w, "Unsuccess, cant't move opposite way")
		} else {
			fmt.Fprintln(w, "Success")
		}
	}
}

func findSnake(id int) *snake.Snake {
	for _, user_snake := range g.Snakes {
		if user_snake.Id == id {
			return user_snake
		}
	}
	return nil
}

func findSnakeIp(ip string) *snake.Snake {
	for _, user_snake := range g.Snakes {
		if user_snake.Ip == ip {
			return user_snake
		}
	}
	return nil
}

func main() {
	g = snake_game.Init(20, 20, 3)
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/play", play)
	r.HandleFunc("/state", getGameState)
	r.HandleFunc("/new_snake", new_snake).Methods(http.MethodPost)
	r.HandleFunc("/direction/{direction}/{id}", setDirection)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			g.Step()
			g.RemoveDeadSnakes()
		}
	}()

	http.ListenAndServe(":80", r)
}
