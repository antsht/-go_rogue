package game

import (
	"github.com/user/go-rogue/internal/domain/entities"
)

// DifficultyManager handles dynamic difficulty adjustment (Bonus Task 7)
type DifficultyManager struct {
	modifier        float64
	checkInterval   int
	turnsSinceCheck int
}

// NewDifficultyManager creates a new difficulty manager
func NewDifficultyManager() *DifficultyManager {
	return &DifficultyManager{
		modifier:      1.0,
		checkInterval: 25, // Check every 25 turns
	}
}

// GetModifier returns the current difficulty modifier
func (d *DifficultyManager) GetModifier() float64 {
	return d.modifier
}

// SetModifier sets the difficulty modifier (used when loading saved games)
func (d *DifficultyManager) SetModifier(mod float64) {
	if mod < 0.5 {
		mod = 0.5
	} else if mod > 1.5 {
		mod = 1.5
	}
	d.modifier = mod
}

// Update updates the difficulty based on player performance
func (d *DifficultyManager) Update(session *entities.Session) {
	d.turnsSinceCheck++

	if d.turnsSinceCheck < d.checkInterval {
		return
	}

	d.turnsSinceCheck = 0

	// Analyze player performance
	char := session.Character

	// Calculate health ratio
	healthRatio := float64(char.Health) / float64(char.MaxHealth)

	// Check recent deaths (from session tracking)
	recentDeaths := session.RecentDeaths
	recentEasyKills := session.RecentEasyKills

	// Adjust difficulty based on performance
	if recentDeaths > 0 || healthRatio < 0.3 {
		// Player is struggling - decrease difficulty
		d.modifier -= 0.1
		if d.modifier < 0.5 {
			d.modifier = 0.5
		}
		session.AddMessage("The dungeon seems slightly less hostile...")
	} else if recentEasyKills > 10 && healthRatio > 0.8 {
		// Player is doing too well - increase difficulty
		d.modifier += 0.1
		if d.modifier > 1.5 {
			d.modifier = 1.5
		}
		session.AddMessage("The dungeon grows more treacherous...")
	}

	// Reset tracking
	session.RecentDeaths = 0
	session.RecentEasyKills = 0

	// Update session modifier
	session.DifficultyModifier = d.modifier
}

// AdjustEnemyStats adjusts enemy stats based on difficulty
func (d *DifficultyManager) AdjustEnemyStats(enemy *entities.Enemy) {
	// Apply modifier to enemy stats
	enemy.Health = int(float64(enemy.Health) * d.modifier)
	enemy.MaxHealth = enemy.Health
	enemy.Strength = int(float64(enemy.Strength) * d.modifier)

	// Ensure minimums
	if enemy.Health < 1 {
		enemy.Health = 1
		enemy.MaxHealth = 1
	}
	if enemy.Strength < 1 {
		enemy.Strength = 1
	}
}

// AdjustItemSpawn adjusts item spawn rates based on difficulty
func (d *DifficultyManager) AdjustItemSpawn() float64 {
	// Lower difficulty = more items
	// Higher difficulty = fewer items
	return 2.0 - d.modifier // Range: 0.5 to 1.5
}

// ShouldSpawnExtraHeal returns true if we should spawn extra healing items
func (d *DifficultyManager) ShouldSpawnExtraHeal(session *entities.Session) bool {
	if session == nil || session.Character == nil {
		return false
	}

	// Spawn extra healing if player health is low and difficulty is reduced
	healthRatio := float64(session.Character.Health) / float64(session.Character.MaxHealth)
	return d.modifier < 0.8 && healthRatio < 0.5
}
