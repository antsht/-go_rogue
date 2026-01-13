package world

import (
	"math/rand"

	"github.com/user/go-rogue/internal/domain/entities"
)

// MimicPlacer handles placing mimic enemies (Bonus Task 8)
type MimicPlacer struct {
	rng *rand.Rand
}

// NewMimicPlacer creates a new mimic placer
func NewMimicPlacer() *MimicPlacer {
	return &MimicPlacer{}
}

// PlaceMimics adds mimic enemies to a level
func (m *MimicPlacer) PlaceMimics(level *entities.Level, levelNum int, seed int64) {
	m.rng = rand.New(rand.NewSource(seed))

	// Mimics only appear after level 5
	if levelNum < 5 {
		return
	}

	// Number of mimics based on level
	numMimics := (levelNum - 4) / 4 // 1 at level 5-8, 2 at 9-12, etc.
	if numMimics > 3 {
		numMimics = 3
	}

	mimicsPlaced := 0

	for _, room := range level.Rooms {
		if mimicsPlaced >= numMimics {
			break
		}

		// Skip start room
		if room.IsStart {
			continue
		}

		// 20% chance per room
		if m.rng.Float64() < 0.2 {
			mimic := entities.NewMimic(levelNum)

			// Choose what item to mimic
			mimickedItem := m.chooseMimickedItem()
			mimic.SetMimickedItem(mimickedItem)

			// Position
			pos := room.GetRandomFloorPosition(entities.NewRNG(m.rng.Int63()))

			// Make sure not on exit
			if pos.Equals(level.ExitPos) {
				continue
			}

			mimic.Position = pos
			room.AddEnemy(mimic)
			mimicsPlaced++
		}
	}
}

// chooseMimickedItem selects what item type the mimic appears as
func (m *MimicPlacer) chooseMimickedItem() *entities.Item {
	roll := m.rng.Intn(100)

	if roll < 40 {
		// Treasure (most common - enticing!)
		return &entities.Item{
			Type:   entities.ItemTypeTreasure,
			Symbol: '*',
			Color:  "yellow",
			Name:   "Gold",
		}
	} else if roll < 60 {
		// Weapon
		return &entities.Item{
			Type:   entities.ItemTypeWeapon,
			Symbol: ')',
			Color:  "cyan",
			Name:   "Weapon",
		}
	} else if roll < 80 {
		// Scroll
		return &entities.Item{
			Type:   entities.ItemTypeScroll,
			Symbol: '?',
			Color:  "white",
			Name:   "Scroll",
		}
	} else {
		// Elixir
		return &entities.Item{
			Type:   entities.ItemTypeElixir,
			Symbol: '!',
			Color:  "magenta",
			Name:   "Elixir",
		}
	}
}

// IsMimicAt checks if there's a mimic at a position that looks like an item
func (m *MimicPlacer) IsMimicAt(level *entities.Level, pos entities.Position) *entities.Enemy {
	for _, room := range level.Rooms {
		for _, enemy := range room.Enemies {
			if enemy.Type == entities.EnemyMimic &&
				enemy.Position.Equals(pos) &&
				!enemy.IsRevealed &&
				enemy.IsAlive() {
				return enemy
			}
		}
	}
	return nil
}
