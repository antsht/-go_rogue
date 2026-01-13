package entities

// Character represents the player character
type Character struct {
	Position      Position  `json:"position"`
	MaxHealth     int       `json:"max_health"`
	Health        int       `json:"health"`
	Dexterity     int       `json:"dexterity"`
	Strength      int       `json:"strength"`
	Weapon        *Item     `json:"weapon"`
	Backpack      *Backpack `json:"backpack"`
	Experience    int       `json:"experience"`
	ExperienceMax int       `json:"experience_max"`
	Level         int       `json:"level"`
	Gold          int       `json:"gold"`
	Armor         int       `json:"armor"`

	// Status effects
	Asleep       bool `json:"asleep"`
	SleepTurns   int  `json:"sleep_turns"`
	ActiveEffects []Effect `json:"active_effects"`

	// Statistics
	Stats CharacterStats `json:"stats"`
}

// CharacterStats tracks gameplay statistics
type CharacterStats struct {
	EnemiesDefeated int `json:"enemies_defeated"`
	FoodConsumed    int `json:"food_consumed"`
	ElixirsDrunk    int `json:"elixirs_drunk"`
	ScrollsRead     int `json:"scrolls_read"`
	HitsDealt       int `json:"hits_dealt"`
	HitsReceived    int `json:"hits_received"`
	TilesTraveled   int `json:"tiles_traveled"`
}

// Effect represents a temporary stat modification
type Effect struct {
	Type          EffectType `json:"type"`
	Value         int        `json:"value"`
	TurnsRemaining int       `json:"turns_remaining"`
}

// EffectType represents different effect types
type EffectType int

const (
	EffectStrength EffectType = iota
	EffectDexterity
	EffectMaxHealth
)

// NewCharacter creates a new player character with default stats
func NewCharacter() *Character {
	return &Character{
		MaxHealth:     22,
		Health:        22,
		Dexterity:     10,
		Strength:      16,
		Backpack:      NewBackpack(),
		Experience:    0,
		ExperienceMax: 28,
		Level:         1,
		Gold:          0,
		Armor:         5,
		ActiveEffects: make([]Effect, 0),
		Stats:         CharacterStats{},
	}
}

// IsAlive returns true if the character has health remaining
func (c *Character) IsAlive() bool {
	return c.Health > 0
}

// TakeDamage reduces health by the specified amount
func (c *Character) TakeDamage(damage int) {
	c.Health -= damage
	c.Stats.HitsReceived++
	if c.Health < 0 {
		c.Health = 0
	}
}

// Heal restores health up to max health
func (c *Character) Heal(amount int) {
	c.Health += amount
	if c.Health > c.MaxHealth {
		c.Health = c.MaxHealth
	}
}

// GetEffectiveStrength returns strength including active effects
func (c *Character) GetEffectiveStrength() int {
	str := c.Strength
	for _, effect := range c.ActiveEffects {
		if effect.Type == EffectStrength {
			str += effect.Value
		}
	}
	return str
}

// GetEffectiveDexterity returns dexterity including active effects
func (c *Character) GetEffectiveDexterity() int {
	dex := c.Dexterity
	for _, effect := range c.ActiveEffects {
		if effect.Type == EffectDexterity {
			dex += effect.Value
		}
	}
	return dex
}

// GetEffectiveMaxHealth returns max health including active effects
func (c *Character) GetEffectiveMaxHealth() int {
	maxHP := c.MaxHealth
	for _, effect := range c.ActiveEffects {
		if effect.Type == EffectMaxHealth {
			maxHP += effect.Value
		}
	}
	return maxHP
}

// AddEffect adds a temporary effect to the character
func (c *Character) AddEffect(effectType EffectType, value, turns int) {
	c.ActiveEffects = append(c.ActiveEffects, Effect{
		Type:          effectType,
		Value:         value,
		TurnsRemaining: turns,
	})
}

// UpdateEffects decrements effect timers and removes expired effects
func (c *Character) UpdateEffects() {
	remaining := make([]Effect, 0)
	for _, effect := range c.ActiveEffects {
		effect.TurnsRemaining--
		if effect.TurnsRemaining > 0 {
			remaining = append(remaining, effect)
		} else {
			// Handle max health effect expiry
			if effect.Type == EffectMaxHealth && c.Health > c.MaxHealth {
				c.Health = c.MaxHealth
				if c.Health <= 0 {
					c.Health = 1 // Keep minimum health as per spec
				}
			}
		}
	}
	c.ActiveEffects = remaining
}

// IncreaseMaxHealth permanently increases max health and current health
func (c *Character) IncreaseMaxHealth(amount int) {
	c.MaxHealth += amount
	c.Health += amount
}

// AddGold adds gold to the character
func (c *Character) AddGold(amount int) {
	c.Gold += amount
}

// GetDamage calculates damage dealt by the character
func (c *Character) GetDamage() int {
	baseDamage := c.GetEffectiveStrength() / 3
	if c.Weapon != nil {
		baseDamage += c.Weapon.Strength
	}
	if baseDamage < 1 {
		baseDamage = 1
	}
	return baseDamage
}

// Move updates the character position and tracks statistics
func (c *Character) Move(newPos Position) {
	c.Position = newPos
	c.Stats.TilesTraveled++
}

// WakeUp resets sleep status
func (c *Character) WakeUp() {
	c.Asleep = false
	c.SleepTurns = 0
}

// PutToSleep puts character to sleep for specified turns
func (c *Character) PutToSleep(turns int) {
	c.Asleep = true
	c.SleepTurns = turns
}

// ProcessSleep handles sleep turn decrement
func (c *Character) ProcessSleep() bool {
	if c.Asleep {
		c.SleepTurns--
		if c.SleepTurns <= 0 {
			c.WakeUp()
		}
		return true
	}
	return false
}
