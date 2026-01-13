package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/user/go-rogue/internal/domain/game"
	"github.com/user/go-rogue/internal/presentation/renderer"
)

// LeaderboardViewRender renders the leaderboard view
type LeaderboardViewRender struct {
	screen     *renderer.Screen
	gameEngine *game.Engine
}

// NewLeaderboardViewRender creates a new leaderboard view renderer
func NewLeaderboardViewRender(screen *renderer.Screen, gameEngine *game.Engine) *LeaderboardViewRender {
	return &LeaderboardViewRender{
		screen:     screen,
		gameEngine: gameEngine,
	}
}

// Render draws the leaderboard view
func (v *LeaderboardViewRender) Render() {
	width, height := v.screen.Size()

	// Title
	title := "═══ LEADERBOARD ═══"
	v.screen.DrawString(width/2-len(title)/2, 1, title, tcell.ColorYellow, tcell.ColorBlack)

	// Get leaderboard data
	leaderboard := v.gameEngine.GetLeaderboard()
	if leaderboard == nil || len(leaderboard.Results) == 0 {
		msg := "No records yet. Go explore some dungeons!"
		v.screen.DrawString(width/2-len(msg)/2, height/2, msg, tcell.ColorGray, tcell.ColorBlack)
		v.screen.DrawString(width/2-10, height-2, "Press any key to return", tcell.ColorDarkGray, tcell.ColorBlack)
		return
	}

	// Header
	headerY := 4
	v.screen.DrawString(3, headerY, "RANK", tcell.ColorOrange, tcell.ColorBlack)
	v.screen.DrawString(10, headerY, "GOLD", tcell.ColorYellow, tcell.ColorBlack)
	v.screen.DrawString(20, headerY, "LEVEL", tcell.ColorTeal, tcell.ColorBlack)
	v.screen.DrawString(28, headerY, "ENEMIES", tcell.ColorRed, tcell.ColorBlack)
	v.screen.DrawString(38, headerY, "TILES", tcell.ColorGreen, tcell.ColorBlack)
	v.screen.DrawString(48, headerY, "HITS D/R", tcell.ColorPurple, tcell.ColorBlack)
	v.screen.DrawString(60, headerY, "STATUS", tcell.ColorWhite, tcell.ColorBlack)

	// Separator
	v.screen.DrawString(3, headerY+1, "────────────────────────────────────────────────────────────", tcell.ColorOrange, tcell.ColorBlack)

	// Results
	results := leaderboard.GetTopResults(15)
	for i, result := range results {
		y := headerY + 2 + i

		// Rank
		rank := itoa(i+1) + "."
		v.screen.DrawString(3, y, rank, tcell.ColorWhite, tcell.ColorBlack)

		// Gold
		v.screen.DrawString(10, y, itoa(result.GoldCollected), tcell.ColorYellow, tcell.ColorBlack)

		// Level reached
		v.screen.DrawString(20, y, itoa(result.LevelReached), tcell.ColorTeal, tcell.ColorBlack)

		// Enemies defeated
		v.screen.DrawString(28, y, itoa(result.EnemiesDefeated), tcell.ColorRed, tcell.ColorBlack)

		// Tiles traveled
		v.screen.DrawString(38, y, itoa(result.TilesTraveled), tcell.ColorGreen, tcell.ColorBlack)

		// Hits dealt/received
		hitsStr := itoa(result.HitsDealt) + "/" + itoa(result.HitsReceived)
		v.screen.DrawString(48, y, hitsStr, tcell.ColorPurple, tcell.ColorBlack)

		// Status
		status := "Dead"
		statusColor := tcell.ColorRed
		if result.Victory {
			status = "Victory!"
			statusColor = tcell.ColorGreen
		}
		v.screen.DrawString(60, y, status, statusColor, tcell.ColorBlack)
	}

	// Footer
	v.screen.DrawString(width/2-10, height-2, "Press any key to return", tcell.ColorGray, tcell.ColorBlack)
}
