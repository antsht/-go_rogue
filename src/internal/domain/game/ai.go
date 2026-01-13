package game

import (
	"math/rand"

	"github.com/user/go-rogue/internal/domain/entities"
)

// AI handles enemy artificial intelligence
type AI struct {
	combat *Combat
}

// NewAI creates a new AI handler
func NewAI() *AI {
	return &AI{
		combat: NewCombat(),
	}
}

// ProcessEnemies processes all enemy actions for a turn
func (ai *AI) ProcessEnemies(session *entities.Session) {
	level := session.Level
	playerPos := session.Character.Position

	for _, room := range level.Rooms {
		for _, enemy := range room.Enemies {
			if !enemy.IsAlive() {
				continue
			}

			ai.processEnemy(session, enemy, playerPos, room)
		}
	}
}

// processEnemy handles a single enemy's turn
func (ai *AI) processEnemy(session *entities.Session, enemy *entities.Enemy, playerPos entities.Position, room *entities.Room) {
	distance := enemy.Position.Distance(playerPos)

	// Check if player is in hostility range
	if distance <= enemy.Hostility {
		enemy.IsAggro = true
	}

	// Special enemy type handling
	switch enemy.Type {
	case entities.EnemyGhost:
		ai.processGhost(session, enemy, playerPos, room)
		return
	case entities.EnemyOgre:
		ai.processOgre(session, enemy, playerPos, room)
		return
	case entities.EnemySnakeMage:
		ai.processSnakeMage(session, enemy, playerPos, room)
		return
	}

	// Standard enemy behavior
	if enemy.IsAggro {
		ai.chasePlayer(session, enemy, playerPos, room)
	} else {
		ai.randomMove(session, enemy, room)
	}
}

// processGhost handles ghost-specific behavior
func (ai *AI) processGhost(session *entities.Session, enemy *entities.Enemy, playerPos entities.Position, room *entities.Room) {
	// Teleport randomly within room
	if rand.Float64() < 0.3 && !enemy.IsAggro {
		newPos := room.GetRandomFloorPosition(entities.NewRNG(rand.Int63()))
		if !newPos.Equals(playerPos) {
			enemy.Position = newPos
		}
	}

	// Toggle visibility when not in combat
	if !enemy.IsAggro && rand.Float64() < 0.2 {
		enemy.IsVisible = !enemy.IsVisible
	}

	// If aggro, become visible and chase
	if enemy.IsAggro {
		enemy.IsVisible = true
		ai.chasePlayer(session, enemy, playerPos, room)
	}
}

// processOgre handles ogre-specific behavior
func (ai *AI) processOgre(session *entities.Session, enemy *entities.Enemy, playerPos entities.Position, room *entities.Room) {
	// If resting after attack, skip turn and prepare counterattack
	if enemy.IsResting {
		enemy.IsResting = false
		// Guaranteed counterattack next to player
		if enemy.Position.Distance(playerPos) <= 1 {
			ai.combat.EnemyAttack(session, enemy)
		}
		return
	}

	if enemy.IsAggro {
		// Ogre moves two tiles per turn
		for i := 0; i < 2; i++ {
			if enemy.Position.Distance(playerPos) <= 1 {
				ai.combat.EnemyAttack(session, enemy)
				return
			}
			ai.moveToward(session, enemy, playerPos, room)
		}
	} else {
		ai.randomMove(session, enemy, room)
	}
}

// processSnakeMage handles snake-mage specific behavior
func (ai *AI) processSnakeMage(session *entities.Session, enemy *entities.Enemy, playerPos entities.Position, room *entities.Room) {
	if enemy.IsAggro {
		// If adjacent, attack
		if enemy.Position.Distance(playerPos) <= 1 {
			ai.combat.EnemyAttack(session, enemy)
			return
		}
		ai.chasePlayer(session, enemy, playerPos, room)
	} else {
		// Move diagonally
		dx, dy := enemy.MoveDirection.GetOffset()
		newPos := enemy.Position.Add(dx, dy)

		// If can't move in current direction, switch
		if !room.Contains(newPos) || newPos.Equals(playerPos) {
			enemy.SwitchDiagonalDirection()
			dx, dy = enemy.MoveDirection.GetOffset()
			newPos = enemy.Position.Add(dx, dy)
		}

		if room.Contains(newPos) && !newPos.Equals(playerPos) {
			enemy.Position = newPos
		}
	}
}

// chasePlayer moves enemy toward player
func (ai *AI) chasePlayer(session *entities.Session, enemy *entities.Enemy, playerPos entities.Position, room *entities.Room) {
	// If adjacent to player, attack
	if enemy.Position.Distance(playerPos) <= 1 {
		ai.combat.EnemyAttack(session, enemy)
		return
	}

	// Move toward player
	ai.moveToward(session, enemy, playerPos, room)
}

// moveToward moves enemy one step toward target
func (ai *AI) moveToward(session *entities.Session, enemy *entities.Enemy, target entities.Position, room *entities.Room) {
	// Find best direction using simple pathfinding
	bestDir := entities.DirNone
	bestDist := enemy.Position.Distance(target)

	directions := []entities.Direction{
		entities.DirUp, entities.DirDown, entities.DirLeft, entities.DirRight,
	}

	for _, dir := range directions {
		dx, dy := dir.GetOffset()
		newPos := enemy.Position.Add(dx, dy)

		// Check if valid move
		if !session.Level.IsWalkable(newPos) {
			continue
		}

		// Check if occupied by another enemy
		if ai.isOccupied(session, newPos) {
			continue
		}

		dist := newPos.Distance(target)
		if dist < bestDist {
			bestDist = dist
			bestDir = dir
		}
	}

	if bestDir != entities.DirNone {
		dx, dy := bestDir.GetOffset()
		enemy.Position = enemy.Position.Add(dx, dy)
	}
}

// randomMove moves enemy randomly within room
func (ai *AI) randomMove(session *entities.Session, enemy *entities.Enemy, room *entities.Room) {
	// 50% chance to not move
	if rand.Float64() < 0.5 {
		return
	}

	directions := []entities.Direction{
		entities.DirUp, entities.DirDown, entities.DirLeft, entities.DirRight,
	}

	// Shuffle directions
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	for _, dir := range directions {
		dx, dy := dir.GetOffset()
		newPos := enemy.Position.Add(dx, dy)

		// Must stay in room
		if !room.Contains(newPos) {
			continue
		}

		// Check if walkable and not occupied
		if session.Level.IsWalkable(newPos) && !ai.isOccupied(session, newPos) {
			enemy.Position = newPos
			return
		}
	}
}

// isOccupied checks if a position is occupied by an enemy or player
func (ai *AI) isOccupied(session *entities.Session, pos entities.Position) bool {
	// Check player
	if session.Character.Position.Equals(pos) {
		return true
	}

	// Check other enemies
	for _, room := range session.Level.Rooms {
		for _, enemy := range room.Enemies {
			if enemy.IsAlive() && enemy.Position.Equals(pos) {
				return true
			}
		}
	}

	return false
}
