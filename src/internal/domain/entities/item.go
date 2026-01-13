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
	SubtypeHammer
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

// WeaponAttackRange defines the min and max attack bonus for each weapon type
type WeaponAttackRange struct {
	Min int
	Max int
}

// GetWeaponAttackRange returns the attack bonus range for a weapon subtype
func GetWeaponAttackRange(subtype ItemSubtype) WeaponAttackRange {
	switch subtype {
	case SubtypeDagger:
		return WeaponAttackRange{Min: 1, Max: 3}
	case SubtypeSword:
		return WeaponAttackRange{Min: 2, Max: 5}
	case SubtypeHammer:
		return WeaponAttackRange{Min: 3, Max: 8}
	case SubtypeMace:
		return WeaponAttackRange{Min: 4, Max: 7}
	case SubtypeAxe:
		return WeaponAttackRange{Min: 5, Max: 10}
	default:
		return WeaponAttackRange{Min: 1, Max: 1}
	}
}

// NewWeapon creates a weapon item with default (max) attack bonus
func NewWeapon(subtype ItemSubtype) *Item {
	attackRange := GetWeaponAttackRange(subtype)
	return NewWeaponWithBonus(subtype, attackRange.Max)
}

// NewWeaponWithBonus creates a weapon item with a specific attack bonus
func NewWeaponWithBonus(subtype ItemSubtype, attackBonus int) *Item {
	item := &Item{
		Type:     ItemTypeWeapon,
		Subtype:  subtype,
		Symbol:   ')',
		Color:    "cyan",
		Strength: attackBonus,
	}

	switch subtype {
	case SubtypeDagger:
		item.Name = "Dagger"
	case SubtypeSword:
		item.Name = "Sword"
	case SubtypeHammer:
		item.Name = "Hammer"
	case SubtypeMace:
		item.Name = "Mace"
	case SubtypeAxe:
		item.Name = "Axe"
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

// GetStatsString returns a formatted string with the item's stats
func (i *Item) GetStatsString() string {
	switch i.Type {
	case ItemTypeWeapon:
		return " (+" + intToStr(i.Strength) + " ATK)"
	case ItemTypeFood:
		return " (+" + intToStr(i.Health) + " HP)"
	case ItemTypeElixir:
		if i.Strength > 0 {
			return " (+" + intToStr(i.Strength) + " STR)"
		} else if i.Dexterity > 0 {
			return " (+" + intToStr(i.Dexterity) + " DEX)"
		} else if i.MaxHealth > 0 {
			return " (+" + intToStr(i.MaxHealth) + " MaxHP)"
		}
	case ItemTypeScroll:
		if i.Strength > 0 {
			return " (+" + intToStr(i.Strength) + " STR)"
		} else if i.Dexterity > 0 {
			return " (+" + intToStr(i.Dexterity) + " DEX)"
		} else if i.MaxHealth > 0 {
			return " (+" + intToStr(i.MaxHealth) + " MaxHP)"
		}
	}
	return ""
}

// intToStr converts int to string (internal helper)
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	negative := false
	if n < 0 {
		negative = true
		n = -n
	}
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}
