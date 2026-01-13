package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/user/go-rogue/internal/domain/game"
	"github.com/user/go-rogue/internal/presentation/renderer"
)

// GameOverViewRender renders the game over / victory view
type GameOverViewRender struct {
	screen     *renderer.Screen
	gameEngine *game.Engine
}

// NewGameOverViewRender creates a new game over view renderer
func NewGameOverViewRender(screen *renderer.Screen, gameEngine *game.Engine) *GameOverViewRender {
	return &GameOverViewRender{
		screen:     screen,
		gameEngine: gameEngine,
	}
}

// Render draws the game over or victory screen
func (v *GameOverViewRender) Render(victory bool) {
	width, height := v.screen.Size()
	centerX := width / 2
	centerY := height / 2

	session := v.gameEngine.GetSession()
	if session == nil {
		return
	}

	char := session.Character

	if victory {
		// Victory screen
		title := "═══ VICTORY! ═══"
		v.screen.DrawString(centerX-len(title)/2, centerY-8, title, tcell.ColorGreen, tcell.ColorBlack)

		subtitle := "You have conquered the dungeon!"
		v.screen.DrawString(centerX-len(subtitle)/2, centerY-6, subtitle, tcell.ColorYellow, tcell.ColorBlack)
	} else {
		// Game over screen
		title := "═══ GAME OVER ═══"
		v.screen.DrawString(centerX-len(title)/2, centerY-8, title, tcell.ColorRed, tcell.ColorBlack)

		subtitle := "Your adventure has come to an end..."
		v.screen.DrawString(centerX-len(subtitle)/2, centerY-6, subtitle, tcell.ColorGray, tcell.ColorBlack)
	}

	// Stats box
	statsY := centerY - 3
	v.screen.DrawString(centerX-15, statsY, "╔═══════════════════════════╗", tcell.ColorOrange, tcell.ColorBlack)
	v.screen.DrawString(centerX-15, statsY+1, "║       FINAL STATS         ║", tcell.ColorOrange, tcell.ColorBlack)
	v.screen.DrawString(centerX-15, statsY+2, "╠═══════════════════════════╣", tcell.ColorOrange, tcell.ColorBlack)

	// Stats
	stats := []struct {
		label string
		value int
		color tcell.Color
	}{
		{"Level Reached", session.CurrentLevel, tcell.ColorTeal},
		{"Gold Collected", char.Gold, tcell.ColorYellow},
		{"Enemies Defeated", char.Stats.EnemiesDefeated, tcell.ColorRed},
		{"Tiles Traveled", char.Stats.TilesTraveled, tcell.ColorGreen},
		{"Hits Dealt", char.Stats.HitsDealt, tcell.ColorPurple},
		{"Hits Received", char.Stats.HitsReceived, tcell.ColorPurple},
	}

	for i, stat := range stats {
		y := statsY + 3 + i
		line := "║ " + stat.label + ": "
		// Pad to align values
		for len(line) < 22 {
			line += " "
		}
		v.screen.DrawString(centerX-15, y, line, tcell.ColorOrange, tcell.ColorBlack)
		v.screen.DrawString(centerX+7, y, itoa(stat.value), stat.color, tcell.ColorBlack)
		v.screen.DrawString(centerX+13, y, "║", tcell.ColorOrange, tcell.ColorBlack)
	}

	// Close box
	v.screen.DrawString(centerX-15, statsY+3+len(stats), "╚═══════════════════════════╝", tcell.ColorOrange, tcell.ColorBlack)

	// Options
	optionsY := centerY + 7
	v.screen.DrawString(centerX-10, optionsY, "[N] New Game", tcell.ColorWhite, tcell.ColorBlack)
	v.screen.DrawString(centerX-10, optionsY+1, "[Q] Main Menu", tcell.ColorWhite, tcell.ColorBlack)
}
