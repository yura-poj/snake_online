package snake_game

import (
	"math/rand"
	"snake_online/snake"
	"time"
)

type Treat struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type SnakeGame struct {
	Treats []Treat        `json:"treats"`
	Snakes []*snake.Snake `json:"snakes"`

	Width      int `json:"width"`
	Height     int `json:"height"`
	BestScore  int `json:"bestScore"`
	NumberFood int `json:"-"`
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

func (g *SnakeGame) NewSnake(ip string) int {
	s := snake.NewSnake(0, 0, ip)
	g.Snakes = append(g.Snakes, s)
	return s.Id
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

func (g *SnakeGame) Step() {
	for _, user_snake := range g.Snakes {
		user_snake.Move()
		checkCollision(user_snake)
		checkLunch(user_snake)
		user_snake.Untouchable -= 1

		if checkCollision(user_snake) {
			user_snake.GameOver = true
		}

		checkLunch(user_snake)
		checkBestScore(user_snake)
	}
}

func (g *SnakeGame) RemoveDeadSnakes() {
	activeSnakes := []*snake.Snake{}
	for _, user_snake := range g.Snakes {
		if user_snake.GameOver {
			continue
		}
		activeSnakes = append(activeSnakes, user_snake)
	}
	g.Snakes = activeSnakes
}

func checkCollision(user_snake *snake.Snake) bool {
	head := user_snake.Head()

	if head.X > g.Width || head.Y > g.Height || head.X < 0 || head.Y < 0 {
		return true
	}

	if user_snake.Untouchable >= 0 {
		return false
	}

	for _, new_snake := range g.Snakes {
		if new_snake == user_snake {
			for _, block := range new_snake.Body[1:] {
				if head.X == block.X && head.Y == block.Y {
					return true
				}
			}
			continue
		}
		for _, block := range new_snake.Body {
			if head.X == block.X && head.Y == block.Y {
				return true
			}
		}
	}

	return false
}

func checkLunch(user_snake *snake.Snake) {
	head := user_snake.Head()
	for i, treat := range g.Treats {
		if head.X == treat.X && head.Y == treat.Y {
			user_snake.Grow()
			appearTreat(i)
		}
	}
}

func checkBestScore(user_snake *snake.Snake) {
	if g.BestScore < user_snake.Score {
		g.BestScore = user_snake.Score
	}
}
