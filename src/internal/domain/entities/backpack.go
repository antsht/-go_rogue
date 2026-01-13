package entities

const (
	MaxItemsPerType = 9
)

// Backpack represents the player's inventory
type Backpack struct {
	Food    []*Item `json:"food"`
	Elixirs []*Item `json:"elixirs"`
	Scrolls []*Item `json:"scrolls"`
	Weapons []*Item `json:"weapons"`
	Keys    []*Item `json:"keys"`
}

// NewBackpack creates a new empty backpack
func NewBackpack() *Backpack {
	return &Backpack{
		Food:    make([]*Item, 0, MaxItemsPerType),
		Elixirs: make([]*Item, 0, MaxItemsPerType),
		Scrolls: make([]*Item, 0, MaxItemsPerType),
		Weapons: make([]*Item, 0, MaxItemsPerType),
		Keys:    make([]*Item, 0, MaxItemsPerType),
	}
}

// AddItem adds an item to the backpack if there is space
func (b *Backpack) AddItem(item *Item) bool {
	switch item.Type {
	case ItemTypeFood:
		if len(b.Food) < MaxItemsPerType {
			b.Food = append(b.Food, item)
			return true
		}
	case ItemTypeElixir:
		if len(b.Elixirs) < MaxItemsPerType {
			b.Elixirs = append(b.Elixirs, item)
			return true
		}
	case ItemTypeScroll:
		if len(b.Scrolls) < MaxItemsPerType {
			b.Scrolls = append(b.Scrolls, item)
			return true
		}
	case ItemTypeWeapon:
		if len(b.Weapons) < MaxItemsPerType {
			b.Weapons = append(b.Weapons, item)
			return true
		}
	case ItemTypeKey:
		if len(b.Keys) < MaxItemsPerType {
			b.Keys = append(b.Keys, item)
			return true
		}
	}
	return false
}

// RemoveFood removes and returns food at the given index
func (b *Backpack) RemoveFood(index int) *Item {
	if index < 0 || index >= len(b.Food) {
		return nil
	}
	item := b.Food[index]
	b.Food = append(b.Food[:index], b.Food[index+1:]...)
	return item
}

// RemoveElixir removes and returns elixir at the given index
func (b *Backpack) RemoveElixir(index int) *Item {
	if index < 0 || index >= len(b.Elixirs) {
		return nil
	}
	item := b.Elixirs[index]
	b.Elixirs = append(b.Elixirs[:index], b.Elixirs[index+1:]...)
	return item
}

// RemoveScroll removes and returns scroll at the given index
func (b *Backpack) RemoveScroll(index int) *Item {
	if index < 0 || index >= len(b.Scrolls) {
		return nil
	}
	item := b.Scrolls[index]
	b.Scrolls = append(b.Scrolls[:index], b.Scrolls[index+1:]...)
	return item
}

// RemoveWeapon removes and returns weapon at the given index
func (b *Backpack) RemoveWeapon(index int) *Item {
	if index < 0 || index >= len(b.Weapons) {
		return nil
	}
	item := b.Weapons[index]
	b.Weapons = append(b.Weapons[:index], b.Weapons[index+1:]...)
	return item
}

// RemoveKey removes and returns key at the given index
func (b *Backpack) RemoveKey(index int) *Item {
	if index < 0 || index >= len(b.Keys) {
		return nil
	}
	item := b.Keys[index]
	b.Keys = append(b.Keys[:index], b.Keys[index+1:]...)
	return item
}

// HasKey checks if backpack has a key of the specified subtype
func (b *Backpack) HasKey(subtype ItemSubtype) bool {
	for _, key := range b.Keys {
		if key.Subtype == subtype {
			return true
		}
	}
	return false
}

// RemoveKeyBySubtype removes a key of the specified subtype
func (b *Backpack) RemoveKeyBySubtype(subtype ItemSubtype) *Item {
	for i, key := range b.Keys {
		if key.Subtype == subtype {
			return b.RemoveKey(i)
		}
	}
	return nil
}

// GetFood returns all food items
func (b *Backpack) GetFood() []*Item {
	return b.Food
}

// GetElixirs returns all elixir items
func (b *Backpack) GetElixirs() []*Item {
	return b.Elixirs
}

// GetScrolls returns all scroll items
func (b *Backpack) GetScrolls() []*Item {
	return b.Scrolls
}

// GetWeapons returns all weapon items
func (b *Backpack) GetWeapons() []*Item {
	return b.Weapons
}

// GetKeys returns all key items
func (b *Backpack) GetKeys() []*Item {
	return b.Keys
}

// FoodCount returns the number of food items
func (b *Backpack) FoodCount() int {
	return len(b.Food)
}

// ElixirCount returns the number of elixir items
func (b *Backpack) ElixirCount() int {
	return len(b.Elixirs)
}

// ScrollCount returns the number of scroll items
func (b *Backpack) ScrollCount() int {
	return len(b.Scrolls)
}

// WeaponCount returns the number of weapon items
func (b *Backpack) WeaponCount() int {
	return len(b.Weapons)
}

// KeyCount returns the number of key items
func (b *Backpack) KeyCount() int {
	return len(b.Keys)
}
