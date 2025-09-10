package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"os"
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
			zap.L().Error("Error setting snake direction", zap.Error(err))
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
	file, err := os.OpenFile("snake.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	addLogger(file)

	g = snake_game.Init(20, 20, 3)
	r := mux.NewRouter()
	defer zap.L().Sync()
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

func addLogger(file *os.File) {
	fileWriteSyncer := zapcore.AddSync(file)
	consoleWriteSyncer := zapcore.AddSync(os.Stdout)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), fileWriteSyncer, zapcore.InfoLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), consoleWriteSyncer, zapcore.InfoLevel),
	)

	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
}
