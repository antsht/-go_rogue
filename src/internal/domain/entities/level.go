package entities

const (
	// Map dimensions - based on classic Rogue
	MapWidth  = 80
	MapHeight = 26

	// Grid dimensions for room placement
	GridWidth  = 3
	GridHeight = 3

	// Section dimensions
	SectionWidth  = MapWidth / GridWidth
	SectionHeight = (MapHeight - 2) / GridHeight // Reserve 2 lines for status bar
)

// Level represents a single dungeon level
type Level struct {
	Number    int         `json:"number"`
	Rooms     []*Room     `json:"rooms"`
	Corridors []*Corridor `json:"corridors"`
	Tiles     [][]Tile    `json:"tiles"`
	StartRoom int         `json:"start_room"`
	ExitRoom  int         `json:"exit_room"`
	ExitPos   Position    `json:"exit_pos"`

	// For fog of war
	PlayerRoom *Room `json:"-"` // Current room player is in
}

// NewLevel creates a new empty level
func NewLevel(number int) *Level {
	tiles := make([][]Tile, MapHeight)
	for y := range tiles {
		tiles[y] = make([]Tile, MapWidth)
		for x := range tiles[y] {
			tiles[y][x] = Tile{Type: TileEmpty, Symbol: ' '}
		}
	}

	return &Level{
		Number:    number,
		Rooms:     make([]*Room, 0),
		Corridors: make([]*Corridor, 0),
		Tiles:     tiles,
	}
}

// AddRoom adds a room to the level
func (l *Level) AddRoom(room *Room) {
	l.Rooms = append(l.Rooms, room)
}

// AddCorridor adds a corridor to the level
func (l *Level) AddCorridor(corridor *Corridor) {
	l.Corridors = append(l.Corridors, corridor)
}

// GetRoomAt returns the room containing the position, or nil
func (l *Level) GetRoomAt(pos Position) *Room {
	for _, room := range l.Rooms {
		if room.Contains(pos) {
			return room
		}
	}
	return nil
}

// GetCorridorAt returns the corridor containing the position, or nil
func (l *Level) GetCorridorAt(pos Position) *Corridor {
	for _, corridor := range l.Corridors {
		if corridor.Contains(pos) {
			return corridor
		}
	}
	return nil
}

// GetRoomByID returns a room by its ID
func (l *Level) GetRoomByID(id int) *Room {
	for _, room := range l.Rooms {
		if room.ID == id {
			return room
		}
	}
	return nil
}

// GetTile returns the tile at a position
func (l *Level) GetTile(pos Position) *Tile {
	if pos.X < 0 || pos.X >= MapWidth || pos.Y < 0 || pos.Y >= MapHeight {
		return nil
	}
	return &l.Tiles[pos.Y][pos.X]
}

// SetTile sets the tile at a position
func (l *Level) SetTile(pos Position, tileType TileType, symbol rune) {
	if pos.X < 0 || pos.X >= MapWidth || pos.Y < 0 || pos.Y >= MapHeight {
		return
	}
	l.Tiles[pos.Y][pos.X].Type = tileType
	l.Tiles[pos.Y][pos.X].Symbol = symbol
}

// IsWalkable checks if a position can be walked on
func (l *Level) IsWalkable(pos Position) bool {
	tile := l.GetTile(pos)
	if tile == nil {
		return false
	}

	switch tile.Type {
	case TileFloor, TileCorridor, TileExit, TileEntrance:
		return true
	case TileDoor:
		// Check if door is locked
		return !tile.DoorLocked
	default:
		return false
	}
}

// IsInBounds checks if a position is within map bounds
func (l *Level) IsInBounds(pos Position) bool {
	return pos.X >= 0 && pos.X < MapWidth && pos.Y >= 0 && pos.Y < MapHeight
}

// GetAllEnemies returns all enemies on the level
func (l *Level) GetAllEnemies() []*Enemy {
	var enemies []*Enemy
	for _, room := range l.Rooms {
		enemies = append(enemies, room.Enemies...)
	}
	return enemies
}

// GetEnemyAt returns the enemy at a position, if any
func (l *Level) GetEnemyAt(pos Position) *Enemy {
	for _, room := range l.Rooms {
		if enemy := room.GetEnemyAt(pos); enemy != nil {
			return enemy
		}
	}
	return nil
}

// GetItemAt returns the item at a position, if any
func (l *Level) GetItemAt(pos Position) *Item {
	for _, room := range l.Rooms {
		if item := room.GetItemAt(pos); item != nil {
			return item
		}
	}
	return nil
}

// RemoveItem removes an item from the level
func (l *Level) RemoveItem(item *Item) {
	for _, room := range l.Rooms {
		for i, it := range room.Items {
			if it == item {
				room.Items = append(room.Items[:i], room.Items[i+1:]...)
				return
			}
		}
	}
}

// RemoveEnemy removes an enemy from the level
func (l *Level) RemoveEnemy(enemy *Enemy) {
	for _, room := range l.Rooms {
		room.RemoveEnemy(enemy)
	}
}

// GetStartRoom returns the starting room
func (l *Level) GetStartRoom() *Room {
	return l.GetRoomByID(l.StartRoom)
}

// GetExitRoom returns the exit room
func (l *Level) GetExitRoom() *Room {
	return l.GetRoomByID(l.ExitRoom)
}

// MarkExplored marks a position as explored
func (l *Level) MarkExplored(pos Position) {
	if tile := l.GetTile(pos); tile != nil {
		tile.Explored = true
	}
}

// MarkVisible marks a position as visible
func (l *Level) MarkVisible(pos Position, visible bool) {
	if tile := l.GetTile(pos); tile != nil {
		tile.Visible = visible
		if visible {
			tile.Explored = true
		}
	}
}

// ClearVisibility clears all visibility flags
func (l *Level) ClearVisibility() {
	for y := range l.Tiles {
		for x := range l.Tiles[y] {
			l.Tiles[y][x].Visible = false
		}
	}
}

// ExploreRoom marks all tiles in a room as explored
func (l *Level) ExploreRoom(room *Room) {
	room.Explored = true
	for y := room.Y; y < room.Y+room.Height; y++ {
		for x := room.X; x < room.X+room.Width; x++ {
			l.MarkExplored(Position{X: x, Y: y})
		}
	}
}

// GetAdjacentRooms returns rooms connected to the given room
func (l *Level) GetAdjacentRooms(roomID int) []*Room {
	var adjacent []*Room
	for _, corridor := range l.Corridors {
		if corridor.FromRoom == roomID {
			if room := l.GetRoomByID(corridor.ToRoom); room != nil {
				adjacent = append(adjacent, room)
			}
		} else if corridor.ToRoom == roomID {
			if room := l.GetRoomByID(corridor.FromRoom); room != nil {
				adjacent = append(adjacent, room)
			}
		}
	}
	return adjacent
}
