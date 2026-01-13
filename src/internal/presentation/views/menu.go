package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/user/go-rogue/internal/domain/game"
	"github.com/user/go-rogue/internal/presentation/renderer"
)

// MenuView renders the main menu
type MenuView struct {
	screen     *renderer.Screen
	gameEngine *game.Engine
}

// NewMenuView creates a new menu view
func NewMenuView(screen *renderer.Screen, gameEngine *game.Engine) *MenuView {
	return &MenuView{
		screen:     screen,
		gameEngine: gameEngine,
	}
}

// Render draws the main menu
func (v *MenuView) Render() {
	width, height := v.screen.Size()
	centerX := width / 2
	centerY := height / 2

	// Title
	title := "GO ROGUE"
	subtitle := "A Roguelike Adventure"

	v.screen.DrawString(centerX-len(title)/2, centerY-8, title, tcell.ColorOrange, tcell.ColorBlack)
	v.screen.DrawString(centerX-len(subtitle)/2, centerY-6, subtitle, tcell.ColorYellow, tcell.ColorBlack)

	// ASCII art dungeon
	art := []string{
		"    ┌─────────┐     ┌─────────┐",
		"    │.........│─────│.........│",
		"    │....@....│     │.........│",
		"    │.........│     │.........│",
		"    └────┬────┘     └────┬────┘",
		"         │               │",
		"    ┌────┴────┐     ┌────┴────┐",
		"    │.........│─────│....%....│",
		"    │.........│     │.........│",
		"    └─────────┘     └─────────┘",
	}

	for i, line := range art {
		v.screen.DrawString(centerX-len(line)/2, centerY-4+i, line, tcell.ColorOrange, tcell.ColorBlack)
	}

	// Menu options
	menuY := centerY + 8

	v.screen.DrawString(centerX-10, menuY, "[N] New Game", tcell.ColorWhite, tcell.ColorBlack)

	if v.gameEngine.CanContinue() {
		v.screen.DrawString(centerX-10, menuY+1, "[C] Continue", tcell.ColorGreen, tcell.ColorBlack)
	} else {
		v.screen.DrawString(centerX-10, menuY+1, "[C] Continue", tcell.ColorDarkGray, tcell.ColorBlack)
	}

	v.screen.DrawString(centerX-10, menuY+2, "[L] Leaderboard", tcell.ColorWhite, tcell.ColorBlack)
	v.screen.DrawString(centerX-10, menuY+3, "[Q] Quit", tcell.ColorWhite, tcell.ColorBlack)

	// Footer
	footer := "Press a key to select"
	v.screen.DrawString(centerX-len(footer)/2, height-2, footer, tcell.ColorGray, tcell.ColorBlack)
}
