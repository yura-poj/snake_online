package snake_game

import (
	"math/rand"
	"snake_online/snake"
	"time"
)

type Treat struct {
	X, Y int
}

type SnakeGame struct {
	Treats []Treat
	Snakes []*snake.Snake

	Width, Height int
	Best          int
	NumberFood    int
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var g = SnakeGame{}

func Init(width, height, numberFood int) *SnakeGame {
	g.Width, g.Height = width, height
	g.NumberFood = numberFood
	for i := 0; i < numberFood; i++ {
		g.Treats = append(g.Treats, Treat{X: 0, Y: 0})
		appearTreat(i)
	}

	return &g
}

func appearTreat(index int) {
	x, y := 0, 0
	for {
		x, y = rand.Intn(g.Width), rand.Intn(g.Height)
		if blockEmpty(x, y) {
			g.Treats[index] = Treat{x, y}
			break
		}
	}
}

func blockEmpty(x, y int) bool {
	for _, userSnake := range g.Snakes {
		for _, block := range userSnake.Body {
			if block.X == x && block.Y == y {
				return false
			}
		}
	}

	for _, treat := range g.Treats {
		if treat.X == x && treat.Y == y {
			return false
		}
	}

	return true
}

func NewSnake(ip string) {
	s := snake.NewSnake(0, 0, ip)
	g.Snakes = append(g.Snakes, s)
}

func move() {
	
}
