package game

import (
	"math/rand"
	"time"

	"github.com/user/go-rogue/internal/data"
	"github.com/user/go-rogue/internal/domain/entities"
	"github.com/user/go-rogue/internal/domain/world"
)

const (
	MaxLevels = 21
)

// Engine manages the game logic
type Engine struct {
	session     *entities.Session
	dataManager *data.Manager
	worldGen    *world.Generator
	combat      *Combat
	ai          *AI
	visibility  *Visibility
	difficulty  *DifficultyManager

	levelSeeds  []int64
	currentSeed int64
}

// NewEngine creates a new game engine
func NewEngine(dataManager *data.Manager) *Engine {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	return &Engine{
		dataManager: dataManager,
		worldGen:    world.NewGenerator(),
		combat:      NewCombat(),
		ai:          NewAI(),
		visibility:  NewVisibility(),
		difficulty:  NewDifficultyManager(),
		levelSeeds:  make([]int64, MaxLevels),
		currentSeed: seed,
	}
}

// NewGame starts a new game
func (e *Engine) NewGame() {
	// Generate seeds for all levels
	for i := 0; i < MaxLevels; i++ {
		e.levelSeeds[i] = rand.Int63()
	}

	// Create new session
	e.session = entities.NewSession()

	// Preserve difficulty from previous game in this session
	// (DifficultyManager persists across games, only resets on fresh terminal start)
	e.session.DifficultyModifier = e.difficulty.GetModifier()

	// Generate first level
	e.generateLevel(1)

	// Place character in starting room
	e.placeCharacterInStartRoom()

	// Initial visibility update
	e.updateVisibility()

	// Save initial game state so player can continue from level 1
	e.saveGame()

	e.session.AddMessage("Welcome to the dungeon! Find the exit (%) to descend.")
}

// ContinueGame loads a saved game
func (e *Engine) ContinueGame() bool {
	saveData, err := e.dataManager.LoadGame()
	if err != nil {
		return false
	}

	e.session = saveData.Session
	e.levelSeeds = saveData.AllLevelSeeds
	e.currentSeed = saveData.LevelSeed

	// Restore difficulty modifier from saved session
	e.difficulty.SetModifier(e.session.DifficultyModifier)

	// Store character position before regenerating level
	charPos := e.session.Character.Position

	// Regenerate the current level (same seed = same layout)
	e.generateLevel(e.session.CurrentLevel)

	// Restore character position
	e.session.Character.Position = charPos

	e.session.AddMessage("Welcome back, adventurer!")
	e.updateVisibility()

	return true
}

// CanContinue checks if there's a saved game
func (e *Engine) CanContinue() bool {
	return e.dataManager.HasSavedGame()
}

// GetSession returns the current session
func (e *Engine) GetSession() *entities.Session {
	return e.session
}

// GetLeaderboard returns the leaderboard
func (e *Engine) GetLeaderboard() *entities.Leaderboard {
	leaderboard, _ := e.dataManager.LoadLeaderboard()
	return leaderboard
}

// generateLevel generates a new dungeon level
func (e *Engine) generateLevel(levelNum int) {
	seed := e.levelSeeds[levelNum-1]
	e.currentSeed = seed

	level := e.worldGen.Generate(levelNum, seed, e.difficulty.GetModifier())
	e.session.Level = level
	e.session.CurrentLevel = levelNum
}

// placeCharacterInStartRoom places the character in the starting room
func (e *Engine) placeCharacterInStartRoom() {
	level := e.session.Level
	startRoom := level.GetStartRoom()
	if startRoom == nil {
		return
	}

	// Place in center of start room
	pos := startRoom.GetCenter()
	e.session.Character.Position = pos
}

// MovePlayer moves the player in a direction
func (e *Engine) MovePlayer(dir entities.Direction) bool {
	if e.session == nil || e.session.Character == nil {
		return false
	}

	char := e.session.Character
	level := e.session.Level

	// Calculate new position
	dx, dy := dir.GetOffset()
	newPos := char.Position.Add(dx, dy)

	// Check for enemy at new position
	if enemy := level.GetEnemyAt(newPos); enemy != nil {
		// If it's an unrevealed mimic, reveal and let it take its turn
		if enemy.Type == entities.EnemyMimic && !enemy.IsRevealed {
			enemy.RevealMimic()
			enemy.IsAggro = true
			e.session.AddMessage("It's a Mimic!")
			e.processTurn()
			return true
		}

		// Attack the enemy
		e.combat.PlayerAttack(e.session, enemy)
		e.processTurn()
		return true
	}

	// Check if walkable
	if !level.IsWalkable(newPos) {
		// Check for locked door
		tile := level.GetTile(newPos)
		if tile != nil && tile.Type == entities.TileDoor && tile.DoorLocked {
			doorColor := tile.DoorColor
			if e.tryUnlockDoor(newPos, tile) {
				e.session.AddMessage("You unlock the " + doorColor + " door with the " + doorColor + " key!")
			} else {
				e.session.AddMessage("The door is locked. You need a " + doorColor + " key.")
			}
		}
		// Still process turn even if movement failed (enemies act)
		e.processTurn()
		return false
	}

	// Move character
	char.Move(newPos)

	// Check for item pickup
	e.checkItemPickup(newPos)

	// Check for exit
	if level.ExitPos.Equals(newPos) {
		e.descendLevel()
		return true
	}

	// Update visibility
	e.updateVisibility()

	// Process turn
	e.processTurn()

	return true
}

// tryUnlockDoor attempts to unlock a door
func (e *Engine) tryUnlockDoor(pos entities.Position, tile *entities.Tile) bool {
	backpack := e.session.Character.Backpack

	if backpack.HasKey(tile.DoorKeyType) {
		backpack.RemoveKeyBySubtype(tile.DoorKeyType)
		tile.DoorLocked = false
		return true
	}
	return false
}

// checkItemPickup checks for and handles item pickup
func (e *Engine) checkItemPickup(pos entities.Position) {
	level := e.session.Level
	item := level.GetItemAt(pos)

	if item == nil {
		return
	}

	// Handle treasure separately
	if item.Type == entities.ItemTypeTreasure {
		e.session.Character.AddGold(item.Value)
		e.session.AddMessage("You found " + itoa(item.Value) + " gold!")
		level.RemoveItem(item)
		return
	}

	// Try to add to backpack
	if e.session.Character.Backpack.AddItem(item) {
		e.session.AddMessage("You pick up " + item.Name + ".")
		level.RemoveItem(item)
	} else {
		e.session.AddMessage("Your backpack is full!")
	}
}

// descendLevel moves to the next level
func (e *Engine) descendLevel() {
	if e.session.CurrentLevel >= MaxLevels {
		// Victory!
		e.victory()
		return
	}

	// Clear all keys from inventory - keys are only valid for current level
	keyCount := e.session.Character.Backpack.KeyCount()
	e.session.Character.Backpack.ClearKeys()
	if keyCount > 0 {
		e.session.AddMessage("Your keys crumble to dust as you descend...")
	}

	// Generate next level
	e.session.CurrentLevel++
	e.generateLevel(e.session.CurrentLevel)
	e.placeCharacterInStartRoom()
	e.updateVisibility()

	// Save progress AFTER generating new level and placing character
	e.saveGame()

	e.session.AddMessage("You descend to level " + itoa(e.session.CurrentLevel) + "...")
}

// victory handles game victory
func (e *Engine) victory() {
	e.session.SetVictory()
	e.recordResult()
	e.dataManager.DeleteSave()
}

// gameOver handles player death
func (e *Engine) gameOver() {
	e.session.SetGameOver()
	e.recordResult()
	e.dataManager.DeleteSave()
}

// recordResult records the session result to leaderboard
func (e *Engine) recordResult() {
	result := e.session.GetResult()
	e.dataManager.AddToLeaderboard(result)
}

// saveGame saves the current game state
func (e *Engine) saveGame() {
	saveData := &entities.SaveData{
		Session:       e.session,
		LevelSeed:     e.currentSeed,
		AllLevelSeeds: e.levelSeeds,
	}
	e.dataManager.SaveGame(saveData)
}

// processTurn processes a game turn
func (e *Engine) processTurn() {
	e.session.IncrementTurn()

	// Update character effects
	e.session.Character.UpdateEffects()

	// Process enemy actions
	e.ai.ProcessEnemies(e.session)

	// Update difficulty based on performance
	e.difficulty.Update(e.session)

	// Check for player death
	if !e.session.Character.IsAlive() {
		e.gameOver()
	}
}

// ProcessTurn is called by the main loop after player action
func (e *Engine) ProcessTurn() {
	// This is called externally but actual processing happens after player moves
}

// ProcessPlayerSleep processes a turn while player is asleep
func (e *Engine) ProcessPlayerSleep() {
	if e.session.Character.ProcessSleep() {
		e.session.AddMessage("You are asleep...")
		e.processTurn()
	}
}

// updateVisibility updates the fog of war
func (e *Engine) updateVisibility() {
	if e.session == nil || e.session.Level == nil {
		return
	}

	e.visibility.Update(e.session.Level, e.session.Character.Position)
}

// StartItemSelection begins item selection mode
func (e *Engine) StartItemSelection(itemType entities.ItemType) {
	e.session.SelectingItem = true
	e.session.SelectingItemType = itemType
}

// CancelItemSelection cancels item selection mode
func (e *Engine) CancelItemSelection() {
	e.session.SelectingItem = false
}

// UseItem uses an item from backpack
func (e *Engine) UseItem(index int) {
	if e.session == nil {
		return
	}

	char := e.session.Character
	backpack := char.Backpack

	switch e.session.SelectingItemType {
	case entities.ItemTypeWeapon:
		if weapon := backpack.RemoveWeapon(index); weapon != nil {
			// If already has weapon, drop it
			if char.Weapon != nil {
				e.dropWeapon()
			}
			char.Weapon = weapon
			e.session.AddMessage("You equip the " + weapon.Name + ".")
		}

	case entities.ItemTypeFood:
		if food := backpack.RemoveFood(index); food != nil {
			char.Heal(food.Health)
			char.Stats.FoodConsumed++
			e.session.AddMessage("You eat the " + food.Name + ". Healed " + itoa(food.Health) + " HP.")
		}

	case entities.ItemTypeElixir:
		if elixir := backpack.RemoveElixir(index); elixir != nil {
			e.applyElixir(elixir)
			char.Stats.ElixirsDrunk++
			e.session.AddMessage("You drink the " + elixir.Name + ".")
		}

	case entities.ItemTypeScroll:
		if scroll := backpack.RemoveScroll(index); scroll != nil {
			e.applyScroll(scroll)
			char.Stats.ScrollsRead++
			e.session.AddMessage("You read the " + scroll.Name + ".")
		}
	}
}

// UnequipWeapon unequips the current weapon
func (e *Engine) UnequipWeapon() {
	char := e.session.Character
	if char.Weapon != nil {
		// Add back to backpack if space
		if char.Backpack.AddItem(char.Weapon) {
			e.session.AddMessage("You unequip the " + char.Weapon.Name + ".")
			char.Weapon = nil
		} else {
			e.session.AddMessage("No room in backpack to store weapon.")
		}
	}
}

// dropWeapon drops current weapon on adjacent tile
func (e *Engine) dropWeapon() {
	char := e.session.Character
	level := e.session.Level

	if char.Weapon == nil {
		return
	}

	// Find adjacent empty tile
	directions := []entities.Direction{
		entities.DirUp, entities.DirDown, entities.DirLeft, entities.DirRight,
	}

	for _, dir := range directions {
		dx, dy := dir.GetOffset()
		dropPos := char.Position.Add(dx, dy)

		if level.IsWalkable(dropPos) && level.GetItemAt(dropPos) == nil {
			char.Weapon.Position = dropPos
			// Add to room
			if room := level.GetRoomAt(dropPos); room != nil {
				room.AddItem(char.Weapon)
			}
			e.session.AddMessage("You drop the " + char.Weapon.Name + ".")
			char.Weapon = nil
			return
		}
	}

	// No space to drop
	e.session.AddMessage("No space to drop weapon!")
}

// applyElixir applies an elixir's temporary effect
func (e *Engine) applyElixir(elixir *entities.Item) {
	char := e.session.Character

	if elixir.Strength > 0 {
		char.AddEffect(entities.EffectStrength, elixir.Strength, elixir.Duration)
	}
	if elixir.Dexterity > 0 {
		char.AddEffect(entities.EffectDexterity, elixir.Dexterity, elixir.Duration)
	}
	if elixir.MaxHealth > 0 {
		char.AddEffect(entities.EffectMaxHealth, elixir.MaxHealth, elixir.Duration)
		char.Health += elixir.MaxHealth // Also heal
	}
}

// applyScroll applies a scroll's permanent effect
func (e *Engine) applyScroll(scroll *entities.Item) {
	char := e.session.Character

	if scroll.Strength > 0 {
		char.Strength += scroll.Strength
		e.session.AddMessage("Your strength increases by " + itoa(scroll.Strength) + "!")
	}
	if scroll.Dexterity > 0 {
		char.Dexterity += scroll.Dexterity
		e.session.AddMessage("Your dexterity increases by " + itoa(scroll.Dexterity) + "!")
	}
	if scroll.MaxHealth > 0 {
		char.IncreaseMaxHealth(scroll.MaxHealth)
		e.session.AddMessage("Your max health increases by " + itoa(scroll.MaxHealth) + "!")
	}
}

// itoa converts int to string
func itoa(n int) string {
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
