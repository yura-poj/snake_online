package snake_game

import (
	"snake_online/snake"
	"testing"
)

func TestInit(t *testing.T) {
	game := Init(20, 20, 3)

	if game.Width != 20 {
		t.Errorf("Expected width 20, got %d", game.Width)
	}

	if game.Height != 20 {
		t.Errorf("Expected height 20, got %d", game.Height)
	}

	if len(game.Treats) != 3 {
		t.Errorf("Expected 3 treats, got %d", len(game.Treats))
	}

	if game.BestScore != 0 {
		t.Errorf("Expected initial best score 0, got %d", game.BestScore)
	}
}

func TestNewSnake(t *testing.T) {
	game := Init(20, 20, 3)

	id1 := game.NewSnake("192.168.1.1")
	id2 := game.NewSnake("192.168.1.2")

	if len(game.Snakes) != 2 {
		t.Errorf("Expected 2 snakes, got %d", len(game.Snakes))
	}

	if id1 == id2 {
		t.Error("Expected different IDs for different snakes")
	}

	if game.Snakes[0].Ip != "192.168.1.1" {
		t.Errorf("Expected first snake IP '192.168.1.1', got '%s'", game.Snakes[0].Ip)
	}
}

func TestBlockEmpty(t *testing.T) {
	game := Init(20, 20, 1)
	game.NewSnake("192.168.1.1")

	snakeHead := game.Snakes[0].Head()
	if blockEmpty(snakeHead.X, snakeHead.Y) {
		t.Error("Expected snake head position to not be empty")
	}

	treat := game.Treats[0]
	if blockEmpty(treat.X, treat.Y) {
		t.Error("Expected treat position to not be empty")
	}

	cont := false
	for x := 0; x < 20; x++ {
		for y := 0; y < 20; y++ {
			if !blockEmpty(x, y) {
				if x == game.Treats[0].X && y == game.Treats[0].Y {
					continue
				}
				cont = false

				for _, part := range game.Snakes[0].Body {
					if x == part.X && y == part.Y {
						cont = true
					}
				}
				if cont {
					continue
				}
				t.Error("Expected not to find at one not empty position on the field except one treat and snake")
			}
		}
	}
}

func TestCheckCollisionBoundaries(t *testing.T) {
	tests := []struct {
		name        string
		headX       int
		headY       int
		untouchable int
		want        bool
	}{
		{name: "collision with left boundary", headX: -1, headY: 5, untouchable: -1, want: true},
		{name: "collision with right boundary", headX: 21, headY: 5, untouchable: -1, want: true},
		{name: "collision with top boundary", headX: 5, headY: 21, untouchable: -1, want: true},
		{name: "collision with bottom boundary", headX: 5, headY: -1, untouchable: -1, want: true},
		{name: "no collision with boundaries", headX: 5, headY: 5, untouchable: -1, want: false},
	}
	game := Init(20, 20, 1)
	testSnake := snake.NewSnake(0, 0, "192.168.1.1")
	game.Snakes = []*snake.Snake{testSnake}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSnake.Head().X = tt.headX
			testSnake.Head().Y = tt.headY
			testSnake.Untouchable = tt.untouchable
			got := checkCollision(testSnake)

			if got != tt.want {
				t.Errorf("CheckCollisionBoundaries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckCollisionUntouchable(t *testing.T) {
	game := Init(20, 20, 1)
	testSnake := snake.NewSnake(5, 5, "192.168.1.1")
	game.Snakes = []*snake.Snake{testSnake}

	testSnake.Body[0].X = -1
	testSnake.Body[0].Y = 10
	testSnake.Untouchable = 10
	if checkCollision(testSnake) {
		t.Error("Expected no collision during untouchable period")
	}

	testSnake.Untouchable = -1
	testSnake.Body[0].X = 10
	testSnake.Body[0].Y = 10
	if checkCollision(testSnake) {
		t.Error("Expected no collision in empty space after untouchable period")
	}
}

func TestCheckCollisionSelfCollision(t *testing.T) {
	game := Init(20, 20, 1)
	testSnake := snake.NewSnake(5, 5, "192.168.1.1")

	testSnake.Body = append(testSnake.Body, snake.SnakePart{X: 4, Y: 5})
	testSnake.Body = append(testSnake.Body, snake.SnakePart{X: 3, Y: 5})
	testSnake.Untouchable = -1 // Make sure untouchable period is over

	game.Snakes = []*snake.Snake{testSnake}

	testSnake.Body[0].X = 4
	testSnake.Body[0].Y = 5

	if !checkCollision(testSnake) {
		t.Error("Expected self collision")
	}
}

func TestCheckLunch(t *testing.T) {
	game := Init(20, 20, 1)
	testSnake := snake.NewSnake(5, 5, "192.168.1.1")
	game.Snakes = []*snake.Snake{testSnake}

	initialScore := testSnake.Score
	initialLength := len(testSnake.Body)

	treat := game.Treats[0]
	testSnake.Body[0].X = treat.X
	testSnake.Body[0].Y = treat.Y

	checkLunch(testSnake)

	if testSnake.Score != initialScore+1 {
		t.Errorf("Expected score to increase by 1, got %d", testSnake.Score)
	}

	if len(testSnake.Body) != initialLength+1 {
		t.Errorf("Expected snake length to increase by 1, got %d", len(testSnake.Body))
	}
}

func TestCheckBestScore(t *testing.T) {
	game := Init(20, 20, 1)
	testSnake := snake.NewSnake(5, 5, "192.168.1.1")
	game.Snakes = []*snake.Snake{testSnake}

	testSnake.Score = 10
	game.BestScore = 5

	checkBestScore(testSnake)

	if game.BestScore != 10 {
		t.Errorf("Expected best score to be updated to 10, got %d", game.BestScore)
	}

	testSnake.Score = 3
	checkBestScore(testSnake)

	if game.BestScore != 10 {
		t.Errorf("Expected best score to remain 10, got %d", game.BestScore)
	}
}

func TestStep(t *testing.T) {
	game := Init(20, 20, 1)
	testSnake := snake.NewSnake(5, 5, "192.168.1.1")
	game.Snakes = []*snake.Snake{testSnake}

	initialUntouchable := testSnake.Untouchable
	initialX := testSnake.Body[0].X

	game.Step()

	if testSnake.Body[0].X == initialX {
		t.Error("Expected snake to move during step")
	}

	if testSnake.Untouchable != initialUntouchable-1 {
		t.Errorf("Expected untouchable to decrease by 1, got %d", testSnake.Untouchable)
	}
}

func TestRemoveDeadSnakes(t *testing.T) {
	game := Init(20, 20, 1)

	snake1 := snake.NewSnake(5, 5, "192.168.1.1")
	snake2 := snake.NewSnake(10, 10, "192.168.1.2")
	snake3 := snake.NewSnake(15, 15, "192.168.1.3")

	snake2.GameOver = true

	game.Snakes = []*snake.Snake{snake1, snake2, snake3}

	game.RemoveDeadSnakes()

	if len(game.Snakes) != 2 {
		t.Errorf("Expected 2 alive snakes, got %d", len(game.Snakes))
	}

	for _, s := range game.Snakes {
		if s.GameOver {
			t.Error("Found dead snake in active snakes list")
		}
	}
}
