package entities

import "math/rand"

// EnemyType represents different enemy types
type EnemyType int

const (
	EnemyZombie EnemyType = iota
	EnemyVampire
	EnemyGhost
	EnemyOgre
	EnemySnakeMage
	EnemyMimic // Bonus Task 8
)

// Enemy represents a hostile creature in the game
type Enemy struct {
	Type      EnemyType `json:"type"`
	Name      string    `json:"name"`
	Position  Position  `json:"position"`
	Health    int       `json:"health"`
	MaxHealth int       `json:"max_health"`
	Dexterity int       `json:"dexterity"`
	Strength  int       `json:"strength"`
	Hostility int       `json:"hostility"` // Detection range
	Symbol    rune      `json:"symbol"`
	Color     string    `json:"color"`

	// State tracking
	IsAggro        bool      `json:"is_aggro"`
	IsVisible      bool      `json:"is_visible"`
	IsResting      bool      `json:"is_resting"`       // Ogre resting after attack
	FirstHitMissed bool      `json:"first_hit_missed"` // Vampire mechanic
	MoveDirection  Direction `json:"move_direction"`   // For Snake-Mage diagonal movement

	// For Mimic - what item it mimics
	MimickedItem *Item `json:"mimicked_item"`
	IsRevealed   bool  `json:"is_revealed"` // Mimic revealed when attacked
}

// NewZombie creates a zombie enemy
func NewZombie(level int) *Enemy {
	return &Enemy{
		Type:      EnemyZombie,
		Name:      "Zombie",
		Health:    20 + level*3,
		MaxHealth: 20 + level*3,
		Dexterity: 5,
		Strength:  8 + level,
		Hostility: 5,
		Symbol:    'z',
		Color:     "green",
		IsVisible: true,
	}
}

// NewVampire creates a vampire enemy
func NewVampire(level int) *Enemy {
	return &Enemy{
		Type:      EnemyVampire,
		Name:      "Vampire",
		Health:    15 + level*2,
		MaxHealth: 15 + level*2,
		Dexterity: 12 + level,
		Strength:  7 + level,
		Hostility: 8,
		Symbol:    'v',
		Color:     "red",
		IsVisible: true,
	}
}

// NewGhost creates a ghost enemy
func NewGhost(level int) *Enemy {
	return &Enemy{
		Type:      EnemyGhost,
		Name:      "Ghost",
		Health:    8 + level,
		MaxHealth: 8 + level,
		Dexterity: 14 + level,
		Strength:  4 + level/2,
		Hostility: 4,
		Symbol:    'g',
		Color:     "white",
		IsVisible: true,
	}
}

// NewOgre creates an ogre enemy
func NewOgre(level int) *Enemy {
	return &Enemy{
		Type:      EnemyOgre,
		Name:      "Ogre",
		Health:    30 + level*4,
		MaxHealth: 30 + level*4,
		Dexterity: 4,
		Strength:  15 + level*2,
		Hostility: 6,
		Symbol:    'O',
		Color:     "yellow",
		IsVisible: true,
	}
}

// NewSnakeMage creates a snake-mage enemy
func NewSnakeMage(level int) *Enemy {
	e := &Enemy{
		Type:      EnemySnakeMage,
		Name:      "Snake-Mage",
		Health:    12 + level*2,
		MaxHealth: 12 + level*2,
		Dexterity: 16 + level,
		Strength:  6 + level,
		Hostility: 7,
		Symbol:    's',
		Color:     "white",
		IsVisible: true,
	}
	// Start with random diagonal direction
	directions := []Direction{DirUpLeft, DirUpRight, DirDownLeft, DirDownRight}
	e.MoveDirection = directions[rand.Intn(4)]
	return e
}

// NewMimic creates a mimic enemy (bonus task)
func NewMimic(level int) *Enemy {
	return &Enemy{
		Type:       EnemyMimic,
		Name:       "Mimic",
		Health:     18 + level*2,
		MaxHealth:  18 + level*2,
		Dexterity:  12 + level,
		Strength:   5 + level/2,
		Hostility:  3,
		Symbol:     '*', // Mimics treasure by default
		Color:      "yellow",
		IsVisible:  true,
		IsRevealed: false,
	}
}

// CreateEnemyForLevel creates a random enemy appropriate for the level
func CreateEnemyForLevel(level int) *Enemy {
	// Harder enemies appear more frequently at deeper levels
	roll := rand.Intn(100)

	if level < 5 {
		// Early levels: mostly zombies and ghosts
		if roll < 50 {
			return NewZombie(level)
		} else if roll < 80 {
			return NewGhost(level)
		} else {
			return NewVampire(level)
		}
	} else if level < 10 {
		// Mid levels: add vampires and snake-mages
		if roll < 30 {
			return NewZombie(level)
		} else if roll < 50 {
			return NewGhost(level)
		} else if roll < 75 {
			return NewVampire(level)
		} else {
			return NewSnakeMage(level)
		}
	} else if level < 15 {
		// Later levels: add ogres
		if roll < 20 {
			return NewZombie(level)
		} else if roll < 35 {
			return NewGhost(level)
		} else if roll < 55 {
			return NewVampire(level)
		} else if roll < 80 {
			return NewSnakeMage(level)
		} else {
			return NewOgre(level)
		}
	} else {
		// Deep levels: all enemy types, more dangerous ones
		if roll < 15 {
			return NewZombie(level)
		} else if roll < 25 {
			return NewGhost(level)
		} else if roll < 45 {
			return NewVampire(level)
		} else if roll < 70 {
			return NewSnakeMage(level)
		} else {
			return NewOgre(level)
		}
	}
}

// IsAlive returns true if the enemy has health remaining
func (e *Enemy) IsAlive() bool {
	return e.Health > 0
}

// TakeDamage reduces enemy health
func (e *Enemy) TakeDamage(damage int) {
	e.Health -= damage
	if e.Health < 0 {
		e.Health = 0
	}
}

// GetDamage calculates damage dealt by this enemy
func (e *Enemy) GetDamage() int {
	return e.Strength / 2
}

// GetTreasureValue calculates treasure dropped on death
func (e *Enemy) GetTreasureValue() int {
	base := e.Hostility * 5
	base += e.Strength * 2
	base += e.MaxHealth / 2
	// Add some randomness
	return base + rand.Intn(base/2+1)
}

// GetDisplaySymbol returns the symbol to display
func (e *Enemy) GetDisplaySymbol() rune {
	if e.Type == EnemyMimic && !e.IsRevealed {
		if e.MimickedItem != nil {
			return e.MimickedItem.Symbol
		}
		return '*' // Default to treasure
	}
	return e.Symbol
}

// GetDisplayColor returns the color to display
func (e *Enemy) GetDisplayColor() string {
	if e.Type == EnemyMimic && !e.IsRevealed {
		if e.MimickedItem != nil {
			return e.MimickedItem.Color
		}
		return "yellow"
	}
	return e.Color
}

// RevealMimic reveals a mimic's true form
func (e *Enemy) RevealMimic() {
	if e.Type == EnemyMimic {
		e.IsRevealed = true
		e.Symbol = 'm'
		e.Color = "white"
	}
}

// SetMimickedItem sets what item the mimic appears as
func (e *Enemy) SetMimickedItem(item *Item) {
	if e.Type == EnemyMimic {
		e.MimickedItem = item
		e.Symbol = item.Symbol
		e.Color = item.Color
	}
}

// SwitchDiagonalDirection changes snake-mage movement direction
func (e *Enemy) SwitchDiagonalDirection() {
	if e.Type != EnemySnakeMage {
		return
	}

	switch e.MoveDirection {
	case DirUpLeft:
		if rand.Intn(2) == 0 {
			e.MoveDirection = DirUpRight
		} else {
			e.MoveDirection = DirDownLeft
		}
	case DirUpRight:
		if rand.Intn(2) == 0 {
			e.MoveDirection = DirUpLeft
		} else {
			e.MoveDirection = DirDownRight
		}
	case DirDownLeft:
		if rand.Intn(2) == 0 {
			e.MoveDirection = DirDownRight
		} else {
			e.MoveDirection = DirUpLeft
		}
	case DirDownRight:
		if rand.Intn(2) == 0 {
			e.MoveDirection = DirDownLeft
		} else {
			e.MoveDirection = DirUpRight
		}
	}
}
