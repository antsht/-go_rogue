package entities

// ItemType represents the category of an item
type ItemType int

const (
	ItemTypeTreasure ItemType = iota
	ItemTypeFood
	ItemTypeElixir
	ItemTypeScroll
	ItemTypeWeapon
	ItemTypeKey
)

// ItemSubtype represents specific variations within item types
type ItemSubtype int

const (
	// Food subtypes
	SubtypeRation ItemSubtype = iota
	SubtypeFruit
	SubtypeMeat

	// Elixir subtypes
	SubtypeStrengthElixir
	SubtypeDexterityElixir
	SubtypeHealthElixir

	// Scroll subtypes
	SubtypeStrengthScroll
	SubtypeDexterityScroll
	SubtypeHealthScroll

	// Weapon subtypes
	SubtypeDagger
	SubtypeSword
	SubtypeMace
	SubtypeAxe

	// Key subtypes (for bonus task)
	SubtypeRedKey
	SubtypeBlueKey
	SubtypeGreenKey
	SubtypeYellowKey
)

// Item represents a collectible item in the game
type Item struct {
	Type      ItemType    `json:"type"`
	Subtype   ItemSubtype `json:"subtype"`
	Name      string      `json:"name"`
	Position  Position    `json:"position"`
	Health    int         `json:"health"`     // HP restored (food)
	MaxHealth int         `json:"max_health"` // Max HP increased (scrolls/elixirs)
	Dexterity int         `json:"dexterity"`  // DEX increased
	Strength  int         `json:"strength"`   // STR increased (or weapon damage)
	Value     int         `json:"value"`      // Gold value (treasure)
	Duration  int         `json:"duration"`   // Effect duration for elixirs (in turns)
	Symbol    rune        `json:"symbol"`
	Color     string      `json:"color"`
}

// NewTreasure creates a treasure item
func NewTreasure(value int) *Item {
	return &Item{
		Type:   ItemTypeTreasure,
		Name:   "Gold",
		Value:  value,
		Symbol: '*',
		Color:  "yellow",
	}
}

// NewFood creates a food item
func NewFood(subtype ItemSubtype) *Item {
	item := &Item{
		Type:    ItemTypeFood,
		Subtype: subtype,
		Symbol:  ':',
		Color:   "brown",
	}

	switch subtype {
	case SubtypeRation:
		item.Name = "Ration"
		item.Health = 10
	case SubtypeFruit:
		item.Name = "Fruit"
		item.Health = 5
	case SubtypeMeat:
		item.Name = "Meat"
		item.Health = 15
	}

	return item
}

// NewElixir creates an elixir item
func NewElixir(subtype ItemSubtype) *Item {
	item := &Item{
		Type:     ItemTypeElixir,
		Subtype:  subtype,
		Symbol:   '!',
		Color:    "magenta",
		Duration: 20, // Default duration
	}

	switch subtype {
	case SubtypeStrengthElixir:
		item.Name = "Strength Elixir"
		item.Strength = 5
	case SubtypeDexterityElixir:
		item.Name = "Dexterity Elixir"
		item.Dexterity = 5
	case SubtypeHealthElixir:
		item.Name = "Health Elixir"
		item.MaxHealth = 10
	}

	return item
}

// NewScroll creates a scroll item
func NewScroll(subtype ItemSubtype) *Item {
	item := &Item{
		Type:    ItemTypeScroll,
		Subtype: subtype,
		Symbol:  '?',
		Color:   "white",
	}

	switch subtype {
	case SubtypeStrengthScroll:
		item.Name = "Strength Scroll"
		item.Strength = 2
	case SubtypeDexterityScroll:
		item.Name = "Dexterity Scroll"
		item.Dexterity = 2
	case SubtypeHealthScroll:
		item.Name = "Health Scroll"
		item.MaxHealth = 5
	}

	return item
}

// NewWeapon creates a weapon item
func NewWeapon(subtype ItemSubtype) *Item {
	item := &Item{
		Type:    ItemTypeWeapon,
		Subtype: subtype,
		Symbol:  ')',
		Color:   "cyan",
	}

	switch subtype {
	case SubtypeDagger:
		item.Name = "Dagger"
		item.Strength = 3
	case SubtypeSword:
		item.Name = "Sword"
		item.Strength = 5
	case SubtypeMace:
		item.Name = "Mace"
		item.Strength = 7
	case SubtypeAxe:
		item.Name = "Axe"
		item.Strength = 10
	}

	return item
}

// NewKey creates a key item for doors (bonus feature)
func NewKey(subtype ItemSubtype) *Item {
	item := &Item{
		Type:    ItemTypeKey,
		Subtype: subtype,
		Symbol:  'k',
	}

	switch subtype {
	case SubtypeRedKey:
		item.Name = "Red Key"
		item.Color = "red"
	case SubtypeBlueKey:
		item.Name = "Blue Key"
		item.Color = "blue"
	case SubtypeGreenKey:
		item.Name = "Green Key"
		item.Color = "green"
	case SubtypeYellowKey:
		item.Name = "Yellow Key"
		item.Color = "yellow"
	}

	return item
}

// IsConsumable returns true if the item is consumed on use
func (i *Item) IsConsumable() bool {
	return i.Type == ItemTypeFood || i.Type == ItemTypeElixir || i.Type == ItemTypeScroll
}

// GetDisplaySymbol returns the symbol to display for this item
func (i *Item) GetDisplaySymbol() rune {
	return i.Symbol
}

// GetDisplayColor returns the color name for this item
func (i *Item) GetDisplayColor() string {
	return i.Color
}
