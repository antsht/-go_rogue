package main

import (
	"log"
	"os"

	"github.com/user/go-rogue/internal/data"
	"github.com/user/go-rogue/internal/domain/game"
	"github.com/user/go-rogue/internal/presentation/input"
	"github.com/user/go-rogue/internal/presentation/renderer"
	"github.com/user/go-rogue/internal/presentation/views"
)

func main() {
	// Initialize data layer
	dataManager := data.NewManager("savegame.json", "leaderboard.json")

	// Initialize domain layer
	gameEngine := game.NewEngine(dataManager)

	// Initialize presentation layer
	screen, err := renderer.NewScreen()
	if err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
	}
	defer screen.Close()

	// Initialize views
	viewManager := views.NewManager(screen, gameEngine)

	// Initialize input handler
	inputHandler := input.NewHandler(screen, viewManager, gameEngine)

	// Run the game loop
	if err := runGameLoop(inputHandler, viewManager, gameEngine, screen); err != nil {
		log.Printf("Game error: %v", err)
		os.Exit(1)
	}
}

func runGameLoop(inputHandler *input.Handler, viewManager *views.Manager, gameEngine *game.Engine, screen *renderer.Screen) error {
	// Show main menu first
	viewManager.SetView(views.MainMenu)

	for {
		// Render current view
		viewManager.Render()
		screen.Show()

		// Handle input
		action := inputHandler.HandleInput()

		if action == input.ActionQuit {
			return nil
		}

		// Process game logic if in game view
		if viewManager.CurrentView() == views.GameView {
			gameEngine.ProcessTurn()
		}
	}
}
