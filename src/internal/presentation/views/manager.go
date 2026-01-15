package views

import (
	"github.com/user/go-rogue/internal/domain/game"
	"github.com/user/go-rogue/internal/presentation/renderer"
)

// ViewType represents different game views
type ViewType int

const (
	MainMenu ViewType = iota
	GameView
	InventoryView
	LeaderboardView
	GameOverView
	VictoryView
)

// Manager manages game views
type Manager struct {
	screen      *renderer.Screen
	gameEngine  *game.Engine
	currentView ViewType

	// Individual view renderers
	menuView        *MenuView
	gameViewRender  *GameViewRender
	inventoryView   *InventoryViewRender
	leaderboardView *LeaderboardViewRender
	gameOverView    *GameOverViewRender
}

// NewManager creates a new view manager
func NewManager(screen *renderer.Screen, gameEngine *game.Engine) *Manager {
	m := &Manager{
		screen:      screen,
		gameEngine:  gameEngine,
		currentView: MainMenu,
	}

	// Initialize individual views
	m.menuView = NewMenuView(screen, gameEngine)
	m.gameViewRender = NewGameViewRender(screen, gameEngine)
	m.inventoryView = NewInventoryViewRender(screen, gameEngine)
	m.leaderboardView = NewLeaderboardViewRender(screen, gameEngine)
	m.gameOverView = NewGameOverViewRender(screen, gameEngine)

	return m
}

// SetView changes the current view
func (m *Manager) SetView(view ViewType) {
	m.currentView = view
}

// CurrentView returns the current view type
func (m *Manager) CurrentView() ViewType {
	return m.currentView
}

// Render renders the current view
func (m *Manager) Render() {
	m.screen.Clear()

	// Check if terminal is too small
	if m.screen.IsTooSmall() {
		m.screen.DrawTerminalTooSmall()
		return
	}

	switch m.currentView {
	case MainMenu:
		m.menuView.Render()
	case GameView:
		m.gameViewRender.Render()
	case InventoryView:
		m.inventoryView.Render()
	case LeaderboardView:
		m.leaderboardView.Render()
	case GameOverView:
		m.gameOverView.Render(false)
	case VictoryView:
		m.gameOverView.Render(true)
	}
}
