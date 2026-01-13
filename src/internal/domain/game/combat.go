package game

import (
	"math/rand"

	"github.com/user/go-rogue/internal/domain/entities"
)

// Combat handles combat mechanics
type Combat struct{}

// NewCombat creates a new combat handler
func NewCombat() *Combat {
	return &Combat{}
}

// PlayerAttack handles player attacking an enemy
func (c *Combat) PlayerAttack(session *entities.Session, enemy *entities.Enemy) {
	char := session.Character

	// Reveal mimic on attack
	if enemy.Type == entities.EnemyMimic && !enemy.IsRevealed {
		enemy.RevealMimic()
		session.AddMessage("It's a Mimic!")
	}

	// Vampire first hit always misses
	if enemy.Type == entities.EnemyVampire && !enemy.FirstHitMissed {
		enemy.FirstHitMissed = true
		session.AddMessage("Your attack passes through the Vampire!")
		return
	}

	// Hit check
	hitChance := c.calculateHitChance(char.GetEffectiveDexterity(), enemy.Dexterity)
	if rand.Float64() > hitChance {
		session.AddMessage("You miss the " + enemy.Name + "!")
		return
	}

	// Calculate and apply damage
	damage := c.calculateDamage(char.GetDamage())
	enemy.TakeDamage(damage)
	char.Stats.HitsDealt++

	if enemy.IsAlive() {
		session.AddMessage("You hit the " + enemy.Name + " for " + itoa(damage) + " damage!")
	} else {
		// Enemy defeated
		treasure := enemy.GetTreasureValue()
		char.AddGold(treasure)
		char.Stats.EnemiesDefeated++
		session.Level.RemoveEnemy(enemy)
		session.AddMessage("You defeat the " + enemy.Name + "! +" + itoa(treasure) + " gold!")

		// Update difficulty tracking
		if session.DifficultyModifier > 0 {
			session.RecentEasyKills++
		}
	}

	// Ogre resting mechanic - after being attacked, ogre rests then counterattacks
	if enemy.Type == entities.EnemyOgre && enemy.IsAlive() {
		enemy.IsResting = true
	}
}

// EnemyAttack handles enemy attacking the player
func (c *Combat) EnemyAttack(session *entities.Session, enemy *entities.Enemy) {
	char := session.Character

	// Hit check
	hitChance := c.calculateHitChance(enemy.Dexterity, char.GetEffectiveDexterity())

	// Apply armor reduction to hit chance
	hitChance -= float64(char.Armor) * 0.03

	if rand.Float64() > hitChance {
		session.AddMessage("The " + enemy.Name + " misses you!")
		return
	}

	// Calculate damage
	damage := c.calculateDamage(enemy.GetDamage())

	// Special effects based on enemy type
	switch enemy.Type {
	case entities.EnemyVampire:
		// Vampire reduces max health
		if char.MaxHealth > 5 {
			char.MaxHealth--
			session.AddMessage("The Vampire drains your life force!")
		}

	case entities.EnemySnakeMage:
		// Chance to put player to sleep
		if rand.Float64() < 0.3 {
			char.PutToSleep(2)
			session.AddMessage("The Snake-Mage's magic puts you to sleep!")
		}
	}

	char.TakeDamage(damage)
	session.AddMessage("The " + enemy.Name + " hits you for " + itoa(damage) + " damage!")

	// Update difficulty tracking
	if !char.IsAlive() {
		session.RecentDeaths++
	}
}

// calculateHitChance calculates the chance to hit based on dexterity
func (c *Combat) calculateHitChance(attackerDex, defenderDex int) float64 {
	// Base 70% hit chance, modified by dex difference
	baseChance := 0.70
	dexDiff := float64(attackerDex - defenderDex)
	
	// Each point of dex difference adjusts by 3%
	modifier := dexDiff * 0.03
	
	chance := baseChance + modifier
	
	// Clamp between 10% and 95%
	if chance < 0.10 {
		chance = 0.10
	}
	if chance > 0.95 {
		chance = 0.95
	}
	
	return chance
}

// calculateDamage calculates actual damage with some variance
func (c *Combat) calculateDamage(baseDamage int) int {
	// Add ±20% variance
	variance := float64(baseDamage) * 0.2
	damage := float64(baseDamage) + (rand.Float64()*2-1)*variance
	
	if damage < 1 {
		damage = 1
	}
	
	return int(damage)
}
