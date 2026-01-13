package entities

// Position represents a 2D coordinate in the game world
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// NewPosition creates a new position
func NewPosition(x, y int) Position {
	return Position{X: x, Y: y}
}

// Add returns a new position offset by dx, dy
func (p Position) Add(dx, dy int) Position {
	return Position{X: p.X + dx, Y: p.Y + dy}
}

// Equals checks if two positions are the same
func (p Position) Equals(other Position) bool {
	return p.X == other.X && p.Y == other.Y
}

// Distance calculates Manhattan distance to another position
func (p Position) Distance(other Position) int {
	dx := p.X - other.X
	dy := p.Y - other.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

// Direction represents movement direction
type Direction int

const (
	DirNone Direction = iota
	DirUp
	DirDown
	DirLeft
	DirRight
	DirUpLeft
	DirUpRight
	DirDownLeft
	DirDownRight
)

// GetOffset returns the x,y offset for a direction
func (d Direction) GetOffset() (int, int) {
	switch d {
	case DirUp:
		return 0, -1
	case DirDown:
		return 0, 1
	case DirLeft:
		return -1, 0
	case DirRight:
		return 1, 0
	case DirUpLeft:
		return -1, -1
	case DirUpRight:
		return 1, -1
	case DirDownLeft:
		return -1, 1
	case DirDownRight:
		return 1, 1
	default:
		return 0, 0
	}
}
