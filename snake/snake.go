package snake

type SnakePart struct {
	X, Y int
}

type Action int

const (
	UP Action = iota
	DOWN
	LEFT
	RIGHT
)

func (a Action) Opposite() Action {
	if a%2 == 0 {
		return a + 1
	} else {
		return a - 1
	}
}

type Snake struct {
	Body             []SnakePart
	MoveX, MoveY     int
	Previous         SnakePart
	LastDirection    Action
	CurrentDirection Action
	GameOver         bool
	Untouchable      int
	Score            int
	Ip               string
}

func NewSnake(startX, startY int, ip string) *Snake {
	s := &Snake{}
	for i := 0; i < 3; i++ {
		s.Body = append(s.Body, SnakePart{X: startX + i, Y: startY})
	}
	s.CurrentDirection = LEFT
	s.LastDirection = LEFT
	s.Previous = s.Body[len(s.Body)-1]
	s.Untouchable = 3
	s.Ip = ip
	return s
}

func (s *Snake) SetDirection(action Action) {
	if action != s.LastDirection.Opposite() {
		s.CurrentDirection = action
	}
}

func (s *Snake) ReleaseDirection() {
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
	s.ReleaseDirection()
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
}

func (s *Snake) Head() *SnakePart {
	return &s.Body[0]
}
