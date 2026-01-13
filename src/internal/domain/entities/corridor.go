package entities

// Corridor represents a passage between rooms
type Corridor struct {
	ID       int        `json:"id"`
	Points   []Position `json:"points"`   // All positions that make up the corridor
	FromRoom int        `json:"from_room"` // Room ID where corridor starts
	ToRoom   int        `json:"to_room"`   // Room ID where corridor ends
	Explored bool       `json:"explored"`

	// For doors (bonus feature)
	Doors []Door `json:"doors,omitempty"`
}

// Door represents a door in a corridor (bonus feature)
type Door struct {
	Position Position    `json:"position"`
	Color    string      `json:"color"`
	Locked   bool        `json:"locked"`
	KeyType  ItemSubtype `json:"key_type"`
}

// NewCorridor creates a new corridor between two rooms
func NewCorridor(id, fromRoom, toRoom int) *Corridor {
	return &Corridor{
		ID:       id,
		FromRoom: fromRoom,
		ToRoom:   toRoom,
		Points:   make([]Position, 0),
		Doors:    make([]Door, 0),
	}
}

// AddPoint adds a position to the corridor path
func (c *Corridor) AddPoint(pos Position) {
	c.Points = append(c.Points, pos)
}

// Contains checks if a position is part of this corridor
func (c *Corridor) Contains(pos Position) bool {
	for _, p := range c.Points {
		if p.Equals(pos) {
			return true
		}
	}
	return false
}

// AddDoor adds a door to the corridor (bonus feature)
func (c *Corridor) AddDoor(pos Position, color string, keyType ItemSubtype) {
	c.Doors = append(c.Doors, Door{
		Position: pos,
		Color:    color,
		Locked:   true,
		KeyType:  keyType,
	})
}

// GetDoorAt returns the door at a position, if any
func (c *Corridor) GetDoorAt(pos Position) *Door {
	for i := range c.Doors {
		if c.Doors[i].Position.Equals(pos) {
			return &c.Doors[i]
		}
	}
	return nil
}

// UnlockDoor unlocks a door at the specified position
func (c *Corridor) UnlockDoor(pos Position) bool {
	door := c.GetDoorAt(pos)
	if door != nil && door.Locked {
		door.Locked = false
		return true
	}
	return false
}
