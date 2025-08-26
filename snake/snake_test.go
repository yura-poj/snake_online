package snake

import (
	"testing"
)

func TestNewSnake(t *testing.T) {
	snake := NewSnake(5, 5, "192.168.1.1")

	if snake.Body[0].X != 5 || snake.Body[0].Y != 5 {
		t.Errorf("Expected head at (5,5), got (%d,%d)", snake.Body[0].X, snake.Body[0].Y)
	}

	if snake.Ip != "192.168.1.1" {
		t.Errorf("Expected IP '192.168.1.1', got '%s'", snake.Ip)
	}

	if snake.Score != 0 {
		t.Errorf("Expected initial score 0, got %d", snake.Score)
	}

	if snake.GameOver {
		t.Error("Expected GameOver to be false initially")
	}
}

func TestActionOpposite(t *testing.T) {
	tests := []struct {
		action   Action
		expected Action
	}{
		{UP, DOWN},
		{DOWN, UP},
		{LEFT, RIGHT},
		{RIGHT, LEFT},
	}

	for _, test := range tests {
		result := test.action.opposite()
		if result != test.expected {
			t.Errorf("Expected opposite of %d to be %d, got %d", test.action, test.expected, result)
		}
	}
}

func TestSetDirection(t *testing.T) {
	snake := NewSnake(5, 5, "192.168.1.1")

	err := snake.SetDirection(UP)
	if err != nil {
		t.Errorf("Expected no error when setting valid direction, got %v", err)
	}
	if snake.CurrentDirection != UP {
		t.Errorf("Expected direction UP, got %d", snake.CurrentDirection)
	}

	snake.Move()
	err = snake.SetDirection(DOWN)
	if err == nil {
		t.Error("Expected error when setting opposite direction")
	}
	if err != nil && err.Error() != "cannot reverse direction" {
		t.Errorf("Expected 'cannot reverse direction' error, got '%v'", err)
	}
}

func TestSnakeMove(t *testing.T) {
	snake := NewSnake(5, 5, "192.168.1.1")
	initialHead := snake.Body[0]

	snake.Move()

	if snake.Body[0].X != initialHead.X+1 {
		t.Errorf("Expected head X to be %d, got %d", initialHead.X+1, snake.Body[0].X)
	}
	if snake.Body[0].Y != initialHead.Y {
		t.Errorf("Expected head Y to remain %d, got %d", initialHead.Y, snake.Body[0].Y)
	}
}

func TestSnakeMoveDirections(t *testing.T) {
	tests := []struct {
		direction Action
		expectedX int
		expectedY int
	}{
		{UP, 0, -1},
		{DOWN, 0, 1},
		{LEFT, -1, 0},
		{RIGHT, 1, 0},
	}

	for _, test := range tests {
		snake := NewSnake(10, 10, "192.168.1.1")
		if test.direction == LEFT {
			snake.SetDirection(UP)
			snake.Move()
		}
		snake.SetDirection(test.direction)
		initialHead := snake.Body[0]

		snake.Move()

		expectedX := initialHead.X + test.expectedX
		expectedY := initialHead.Y + test.expectedY

		if snake.Body[0].X != expectedX || snake.Body[0].Y != expectedY {
			t.Errorf("Direction %d: expected head at (%d,%d), got (%d,%d)",
				test.direction, expectedX, expectedY, snake.Body[0].X, snake.Body[0].Y)
		}
	}
}

func TestSnakeGrow(t *testing.T) {
	snake := NewSnake(5, 5, "192.168.1.1")
	initialLength := len(snake.Body)
	initialScore := snake.Score

	snake.Move()
	snake.Grow()

	if len(snake.Body) != initialLength+1 {
		t.Errorf("Expected body length %d, got %d", initialLength+1, len(snake.Body))
	}

	if snake.Score != initialScore+1 {
		t.Errorf("Expected score %d, got %d", initialScore+1, snake.Score)
	}
}

func TestSnakeHead(t *testing.T) {
	snake := NewSnake(5, 5, "192.168.1.1")
	head := snake.Head()

	if head != &snake.Body[0] {
		t.Error("Head() should return pointer to first body part")
	}

	if head.X != 5 || head.Y != 5 {
		t.Errorf("Expected head at (5,5), got (%d,%d)", head.X, head.Y)
	}
}

func TestSnakeUniqueIds(t *testing.T) {
	snake1 := NewSnake(0, 0, "192.168.1.1")
	snake2 := NewSnake(0, 0, "192.168.1.2")

	if snake1.Id == snake2.Id {
		t.Error("Expected different IDs for different snakes")
	}
}
