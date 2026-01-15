package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/user/go-rogue/internal/domain/entities"
	"github.com/user/go-rogue/internal/domain/game"
	"github.com/user/go-rogue/internal/presentation/renderer"
)

// GameViewRender renders the main game view
type GameViewRender struct {
	screen     *renderer.Screen
	gameEngine *game.Engine
}

// NewGameViewRender creates a new game view renderer
func NewGameViewRender(screen *renderer.Screen, gameEngine *game.Engine) *GameViewRender {
	return &GameViewRender{
		screen:     screen,
		gameEngine: gameEngine,
	}
}

// Render draws the game view
func (v *GameViewRender) Render() {
	session := v.gameEngine.GetSession()
	if session == nil || session.Level == nil {
		return
	}

	level := session.Level
	char := session.Character

	// Get offset for centering the game area
	offsetX, offsetY := v.screen.GetGameAreaOffset()

	// Draw the level tiles
	v.screen.DrawLevel(level, char.Position, offsetX, offsetY)

	// Draw items in visible rooms
	for _, room := range level.Rooms {
		if room.Explored {
			for _, item := range room.Items {
				if level.Tiles[item.Position.Y][item.Position.X].Visible {
					v.screen.DrawItem(item, offsetX, offsetY)
				}
			}
		}
	}

	// Draw enemies in visible areas
	for _, room := range level.Rooms {
		for _, enemy := range room.Enemies {
			if enemy.IsAlive() && level.Tiles[enemy.Position.Y][enemy.Position.X].Visible {
				if enemy.IsVisible || enemy.IsAggro {
					v.screen.DrawEnemy(enemy, offsetX, offsetY)
				}
			}
		}
	}

	// Draw the player character
	v.screen.DrawCharacter(char.Position, offsetX, offsetY)

	// Draw status bar
	v.screen.DrawStatusBar(session, offsetX, offsetY)

	// Draw item selection UI if active
	if session.SelectingItem {
		v.renderItemSelection(session, offsetX, offsetY)
	}
}

// renderItemSelection draws the item selection overlay
func (v *GameViewRender) renderItemSelection(session *entities.Session, offsetX, offsetY int) {
	// Draw selection box (wider to fit stats) - positioned relative to game area
	boxWidth := 35
	boxX := offsetX + entities.MapWidth - boxWidth - 1
	boxY := offsetY + 1
	boxHeight := 15

	// Draw background
	for y := boxY; y < boxY+boxHeight; y++ {
		for x := boxX; x < boxX+boxWidth; x++ {
			v.screen.SetCell(x, y, ' ', tcell.ColorWhite, tcell.ColorDarkGray)
		}
	}

	// Draw border
	v.screen.DrawBox(boxX, boxY, boxWidth, boxHeight, tcell.ColorWhite, tcell.ColorDarkGray)

	// Title
	var title string
	var items []*entities.Item

	backpack := session.Character.Backpack

	switch session.SelectingItemType {
	case entities.ItemTypeWeapon:
		title = "Select Weapon"
		items = backpack.GetWeapons()
	case entities.ItemTypeFood:
		title = "Select Food"
		items = backpack.GetFood()
	case entities.ItemTypeElixir:
		title = "Select Elixir"
		items = backpack.GetElixirs()
	case entities.ItemTypeScroll:
		title = "Select Scroll"
		items = backpack.GetScrolls()
	}

	v.screen.DrawString(boxX+2, boxY+1, title, tcell.ColorYellow, tcell.ColorDarkGray)

	// Special option for weapons - unequip
	if session.SelectingItemType == entities.ItemTypeWeapon {
		v.screen.DrawString(boxX+2, boxY+3, "[0] Unequip", tcell.ColorWhite, tcell.ColorDarkGray)
	}

	// List items with stats
	startY := boxY + 4
	if session.SelectingItemType != entities.ItemTypeWeapon {
		startY = boxY + 3
	}

	for i, item := range items {
		if i >= 9 {
			break
		}
		line := "[" + string(rune('1'+i)) + "] " + item.Name + item.GetStatsString()
		v.screen.DrawString(boxX+2, startY+i, line, tcell.ColorWhite, tcell.ColorDarkGray)
	}

	if len(items) == 0 {
		v.screen.DrawString(boxX+2, startY, "No items", tcell.ColorGray, tcell.ColorDarkGray)
	}

	// Instructions
	v.screen.DrawString(boxX+2, boxY+boxHeight-2, "[X/Backspace] Cancel", tcell.ColorGray, tcell.ColorDarkGray)
}
