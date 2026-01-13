package entities

// TileType represents different types of map tiles
type TileType int

const (
	TileEmpty TileType = iota
	TileFloor
	TileWall
	TileCorridor
	TileDoor
	TileExit
	TileEntrance
)

// Tile represents a single map tile
type Tile struct {
	Type     TileType `json:"type"`
	Explored bool     `json:"explored"`
	Visible  bool     `json:"visible"`
	Symbol   rune     `json:"symbol"`

	// For doors (bonus feature)
	DoorColor    string `json:"door_color,omitempty"`
	DoorLocked   bool   `json:"door_locked,omitempty"`
	DoorKeyType  ItemSubtype `json:"door_key_type,omitempty"`
}

// Room represents a room in the dungeon
type Room struct {
	ID         int       `json:"id"`
	X          int       `json:"x"`      // Top-left X position
	Y          int       `json:"y"`      // Top-left Y position
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	GridX      int       `json:"grid_x"` // Position in 3x3 grid (0-2)
	GridY      int       `json:"grid_y"` // Position in 3x3 grid (0-2)
	IsStart    bool      `json:"is_start"`
	IsExit     bool      `json:"is_exit"`
	Explored   bool      `json:"explored"`
	Enemies    []*Enemy  `json:"enemies"`
	Items      []*Item   `json:"items"`
	Entrances  []Position `json:"entrances"` // Door/opening positions
}

// NewRoom creates a new room at the specified position
func NewRoom(id, x, y, width, height, gridX, gridY int) *Room {
	return &Room{
		ID:        id,
		X:         x,
		Y:         y,
		Width:     width,
		Height:    height,
		GridX:     gridX,
		GridY:     gridY,
		Enemies:   make([]*Enemy, 0),
		Items:     make([]*Item, 0),
		Entrances: make([]Position, 0),
	}
}

// Contains checks if a position is inside the room (floor area)
func (r *Room) Contains(pos Position) bool {
	return pos.X > r.X && pos.X < r.X+r.Width-1 &&
		pos.Y > r.Y && pos.Y < r.Y+r.Height-1
}

// ContainsIncludingWalls checks if a position is inside the room including walls
func (r *Room) ContainsIncludingWalls(pos Position) bool {
	return pos.X >= r.X && pos.X < r.X+r.Width &&
		pos.Y >= r.Y && pos.Y < r.Y+r.Height
}

// GetCenter returns the center position of the room
func (r *Room) GetCenter() Position {
	return Position{
		X: r.X + r.Width/2,
		Y: r.Y + r.Height/2,
	}
}

// GetRandomFloorPosition returns a random position on the floor
func (r *Room) GetRandomFloorPosition(rng *RNG) Position {
	x := r.X + 1 + rng.Intn(r.Width-2)
	y := r.Y + 1 + rng.Intn(r.Height-2)
	return Position{X: x, Y: y}
}

// AddEnemy adds an enemy to the room
func (r *Room) AddEnemy(enemy *Enemy) {
	r.Enemies = append(r.Enemies, enemy)
}

// RemoveEnemy removes an enemy from the room
func (r *Room) RemoveEnemy(enemy *Enemy) {
	for i, e := range r.Enemies {
		if e == enemy {
			r.Enemies = append(r.Enemies[:i], r.Enemies[i+1:]...)
			return
		}
	}
}

// AddItem adds an item to the room
func (r *Room) AddItem(item *Item) {
	r.Items = append(r.Items, item)
}

// RemoveItem removes an item from the room
func (r *Room) RemoveItem(item *Item) {
	for i, it := range r.Items {
		if it == item {
			r.Items = append(r.Items[:i], r.Items[i+1:]...)
			return
		}
	}
}

// GetItemAt returns an item at the specified position, if any
func (r *Room) GetItemAt(pos Position) *Item {
	for _, item := range r.Items {
		if item.Position.Equals(pos) {
			return item
		}
	}
	return nil
}

// GetEnemyAt returns an enemy at the specified position, if any
func (r *Room) GetEnemyAt(pos Position) *Enemy {
	for _, enemy := range r.Enemies {
		if enemy.Position.Equals(pos) && enemy.IsAlive() {
			return enemy
		}
	}
	return nil
}

// AddEntrance adds a door/opening position to the room
func (r *Room) AddEntrance(pos Position) {
	r.Entrances = append(r.Entrances, pos)
}

// IsEntrance checks if a position is an entrance to the room
func (r *Room) IsEntrance(pos Position) bool {
	for _, entrance := range r.Entrances {
		if entrance.Equals(pos) {
			return true
		}
	}
	return false
}

// RNG is a simple wrapper for random number generation
type RNG struct {
	seed int64
}

// NewRNG creates a new RNG
func NewRNG(seed int64) *RNG {
	return &RNG{seed: seed}
}

// Intn returns a random integer in [0, n)
func (r *RNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	// Simple linear congruential generator
	r.seed = (r.seed*1103515245 + 12345) & 0x7fffffff
	return int(r.seed % int64(n))
}
