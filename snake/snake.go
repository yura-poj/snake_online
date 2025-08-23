package snake

import "errors"

type SnakePart struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Action int

const (
	UP Action = iota
	DOWN
	LEFT
	RIGHT
)

func (a Action) opposite() Action {
	if a%2 == 0 {
		return a + 1
	} else {
		return a - 1
	}
}

type Snake struct {
	Body             []SnakePart `json:"body"`
	MoveX, MoveY     int         `json:"-"`
	Previous         SnakePart   `json:"-"`
	LastDirection    Action      `json:"-"`
	CurrentDirection Action      `json:"-"`
	GameOver         bool        `json:"gameOver"`
	Untouchable      int         `json:"-"`
	Score            int         `json:"score"`
	Ip               string      `json:"ip"`
}

func NewSnake(startX, startY int, ip string) *Snake {
	s := &Snake{}
	for i := 0; i < 3; i++ {
		s.Body = append(s.Body, SnakePart{X: startX + i, Y: startY})
	}
	s.CurrentDirection = RIGHT
	s.LastDirection = RIGHT
	s.Previous = s.Body[len(s.Body)-1]
	s.Untouchable = 3
	s.Ip = ip
	return s
}

func (s *Snake) SetDirection(action Action) error {
	if action != s.LastDirection.opposite() {
		s.CurrentDirection = action
		return nil
	}

	return errors.New("cannot reverse direction")
}

func (s *Snake) releaseDirection() {
	switch s.CurrentDirection {
	case UP:
		s.MoveX = -1
		s.MoveY = 0
	case DOWN:
		s.MoveX = 1
		s.MoveY = 0
	case LEFT:
		s.MoveX = 0
		s.MoveY = -1
	case RIGHT:
		s.MoveX = 0
		s.MoveY = 1
	}
}

func (s *Snake) Move() {
	s.releaseDirection()
	s.LastDirection = s.CurrentDirection
	s.Previous = s.Body[len(s.Body)-1]

	for i := len(s.Body) - 1; i > 0; i-- {
		s.Body[i] = s.Body[i-1]
	}

	s.Body[0].X += s.MoveX
	s.Body[0].Y += s.MoveY
}

func (s *Snake) Grow() {
	s.Body = append(s.Body, s.Previous)
	s.Score += 1
}

func (s *Snake) Head() *SnakePart {
	return &s.Body[0]
}
