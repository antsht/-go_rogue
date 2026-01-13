package input

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/user/go-rogue/internal/domain/entities"
	"github.com/user/go-rogue/internal/domain/game"
	"github.com/user/go-rogue/internal/presentation/renderer"
	"github.com/user/go-rogue/internal/presentation/views"
)

// Action represents a game action
type Action int

const (
	ActionNone Action = iota
	ActionQuit
	ActionMoveUp
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionUseWeapon // h
	ActionUseFood   // j
	ActionUseElixir // k
	ActionUseScroll // e
	ActionSelect1
	ActionSelect2
	ActionSelect3
	ActionSelect4
	ActionSelect5
	ActionSelect6
	ActionSelect7
	ActionSelect8
	ActionSelect9
	ActionSelect0 // Unequip weapon
	ActionCancel
	ActionConfirm
	ActionNewGame
	ActionContinue
	ActionLeaderboard
	ActionPause
)

// Handler handles user input
type Handler struct {
	screen      *renderer.Screen
	viewManager *views.Manager
	gameEngine  *game.Engine

	// Debounce for cancel keys to prevent key repeat issues
	lastCancelTime int64
}

// NewHandler creates a new input handler
func NewHandler(screen *renderer.Screen, viewManager *views.Manager, gameEngine *game.Engine) *Handler {
	return &Handler{
		screen:      screen,
		viewManager: viewManager,
		gameEngine:  gameEngine,
	}
}

// HandleInput processes input and returns the resulting action
func (h *Handler) HandleInput() Action {
	ev := h.screen.PollEvent()

	switch ev := ev.(type) {
	case *tcell.EventResize:
		h.screen.UpdateSize()
		h.screen.Clear()
		return ActionNone

	case *tcell.EventKey:
		return h.handleKeyEvent(ev)
	}

	return ActionNone
}

// handleKeyEvent processes keyboard events
func (h *Handler) handleKeyEvent(ev *tcell.EventKey) Action {
	currentView := h.viewManager.CurrentView()

	// Get session state for item selection check
	session := h.gameEngine.GetSession()
	selectingItem := false
	if session != nil {
		selectingItem = session.SelectingItem
	}

	// Check for ESC key OR Backspace as alternative cancel (Windows ESC workaround)
	isEscapeAction := ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2

	// Handle cancel/escape for item selection and inventory FIRST (before global handler)
	if isEscapeAction {
		// Debounce: ignore cancel keys within 100ms of last cancel (prevents key repeat issues)
		now := time.Now().UnixMilli()
		if now-h.lastCancelTime < 100 {
			return ActionNone
		}
		h.lastCancelTime = now

		// If in GameView with item selection active, cancel selection first
		if currentView == views.GameView && selectingItem {
			h.gameEngine.CancelItemSelection()
			return ActionCancel
		}

		// If in InventoryView, return to game
		if currentView == views.InventoryView {
			h.viewManager.SetView(views.GameView)
			return ActionCancel
		}

		// If in LeaderboardView, return to menu
		if currentView == views.LeaderboardView {
			h.viewManager.SetView(views.MainMenu)
			return ActionCancel
		}

		// If in GameOverView, return to menu
		if currentView == views.GameOverView || currentView == views.VictoryView {
			h.viewManager.SetView(views.MainMenu)
			return ActionConfirm
		}

		// If in GameView (not selecting), go to main menu
		if currentView == views.GameView {
			h.viewManager.SetView(views.MainMenu)
			return ActionNone
		}

		// In main menu, quit
		if currentView == views.MainMenu {
			return ActionQuit
		}
	}

	if ev.Key() == tcell.KeyCtrlC {
		return ActionQuit
	}

	switch currentView {
	case views.MainMenu:
		return h.handleMenuInput(ev)
	case views.GameView:
		return h.handleGameInput(ev)
	case views.InventoryView:
		return h.handleInventoryInput(ev)
	case views.LeaderboardView:
		return h.handleLeaderboardInput(ev)
	case views.GameOverView:
		return h.handleGameOverInput(ev)
	}

	return ActionNone
}

// handleMenuInput processes main menu input
func (h *Handler) handleMenuInput(ev *tcell.EventKey) Action {
	switch ev.Rune() {
	case 'n', 'N':
		h.gameEngine.NewGame()
		h.viewManager.SetView(views.GameView)
		return ActionNewGame
	case 'c', 'C':
		if h.gameEngine.CanContinue() {
			h.gameEngine.ContinueGame()
			h.viewManager.SetView(views.GameView)
			return ActionContinue
		}
	case 'l', 'L':
		h.viewManager.SetView(views.LeaderboardView)
		return ActionLeaderboard
	case 'q', 'Q':
		return ActionQuit
	}

	switch ev.Key() {
	case tcell.KeyEnter:
		h.gameEngine.NewGame()
		h.viewManager.SetView(views.GameView)
		return ActionNewGame
	}

	return ActionNone
}

// handleGameInput processes in-game input
func (h *Handler) handleGameInput(ev *tcell.EventKey) Action {
	session := h.gameEngine.GetSession()
	if session == nil {
		return ActionNone
	}

	// Check if selecting item
	if session.SelectingItem {
		return h.handleItemSelection(ev)
	}

	// Check if player is asleep
	if session.Character.Asleep {
		// Any key wakes up (processes sleep turn)
		h.gameEngine.ProcessPlayerSleep()
		return ActionNone
	}

	// Movement keys (WASD)
	switch ev.Rune() {
	case 'w', 'W':
		h.gameEngine.MovePlayer(entities.DirUp)
		return ActionMoveUp
	case 's', 'S':
		h.gameEngine.MovePlayer(entities.DirDown)
		return ActionMoveDown
	case 'a', 'A':
		h.gameEngine.MovePlayer(entities.DirLeft)
		return ActionMoveLeft
	case 'd', 'D':
		h.gameEngine.MovePlayer(entities.DirRight)
		return ActionMoveRight

	// Item usage keys
	case 'h', 'H':
		h.gameEngine.StartItemSelection(entities.ItemTypeWeapon)
		return ActionUseWeapon
	case 'j', 'J':
		h.gameEngine.StartItemSelection(entities.ItemTypeFood)
		return ActionUseFood
	case 'k', 'K':
		h.gameEngine.StartItemSelection(entities.ItemTypeElixir)
		return ActionUseElixir
	case 'e', 'E':
		h.gameEngine.StartItemSelection(entities.ItemTypeScroll)
		return ActionUseScroll

	// Inventory
	case 'i', 'I':
		h.viewManager.SetView(views.InventoryView)
		return ActionNone
	}

	// Arrow key movement
	switch ev.Key() {
	case tcell.KeyUp:
		h.gameEngine.MovePlayer(entities.DirUp)
		return ActionMoveUp
	case tcell.KeyDown:
		h.gameEngine.MovePlayer(entities.DirDown)
		return ActionMoveDown
	case tcell.KeyLeft:
		h.gameEngine.MovePlayer(entities.DirLeft)
		return ActionMoveLeft
	case tcell.KeyRight:
		h.gameEngine.MovePlayer(entities.DirRight)
		return ActionMoveRight
	}

	return ActionNone
}

// handleItemSelection processes item selection input
func (h *Handler) handleItemSelection(ev *tcell.EventKey) Action {
	session := h.gameEngine.GetSession()

	// Cancel selection with ESC, Backspace, or X key
	if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 ||
		ev.Rune() == 'x' || ev.Rune() == 'X' {
		h.gameEngine.CancelItemSelection()
		return ActionCancel
	}

	// Number key selection
	num := -1
	switch ev.Rune() {
	case '0':
		num = 0
	case '1':
		num = 1
	case '2':
		num = 2
	case '3':
		num = 3
	case '4':
		num = 4
	case '5':
		num = 5
	case '6':
		num = 6
	case '7':
		num = 7
	case '8':
		num = 8
	case '9':
		num = 9
	}

	if num >= 0 {
		// 0 is unequip for weapons
		if num == 0 && session.SelectingItemType == entities.ItemTypeWeapon {
			h.gameEngine.UnequipWeapon()
		} else if num > 0 {
			h.gameEngine.UseItem(num - 1) // Convert to 0-based index
		}
		h.gameEngine.CancelItemSelection()
		return Action(ActionSelect0 + Action(num))
	}

	return ActionNone
}

// handleInventoryInput processes inventory view input
func (h *Handler) handleInventoryInput(ev *tcell.EventKey) Action {
	switch ev.Key() {
	case tcell.KeyEscape:
		h.viewManager.SetView(views.GameView)
		return ActionCancel
	}

	switch ev.Rune() {
	// Use Q to close inventory (not 'I' - key repeat causes immediate close)
	case 'q', 'Q':
		h.viewManager.SetView(views.GameView)
		return ActionCancel
	}

	return ActionNone
}

// handleLeaderboardInput processes leaderboard view input
func (h *Handler) handleLeaderboardInput(ev *tcell.EventKey) Action {
	switch ev.Key() {
	case tcell.KeyEscape, tcell.KeyEnter:
		h.viewManager.SetView(views.MainMenu)
		return ActionCancel
	}

	switch ev.Rune() {
	case 'q', 'Q':
		h.viewManager.SetView(views.MainMenu)
		return ActionCancel
	}

	return ActionNone
}

// handleGameOverInput processes game over view input
func (h *Handler) handleGameOverInput(ev *tcell.EventKey) Action {
	switch ev.Key() {
	case tcell.KeyEscape, tcell.KeyEnter:
		h.viewManager.SetView(views.MainMenu)
		return ActionConfirm
	}

	switch ev.Rune() {
	case 'n', 'N':
		h.gameEngine.NewGame()
		h.viewManager.SetView(views.GameView)
		return ActionNewGame
	case 'q', 'Q':
		h.viewManager.SetView(views.MainMenu)
		return ActionConfirm
	}

	return ActionNone
}
