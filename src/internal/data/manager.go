package data

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/user/go-rogue/internal/domain/entities"
)

// Manager handles game data persistence
type Manager struct {
	saveFile        string
	leaderboardFile string
	dataDir         string
}

// NewManager creates a new data manager
func NewManager(saveFile, leaderboardFile string) *Manager {
	// Get executable directory for data storage (so data persists with the build)
	dataDir := getExecutableDir()

	// Create data directory if it doesn't exist
	os.MkdirAll(dataDir, 0755)

	return &Manager{
		saveFile:        filepath.Join(dataDir, saveFile),
		leaderboardFile: filepath.Join(dataDir, leaderboardFile),
		dataDir:         dataDir,
	}
}

// getExecutableDir returns the directory where the executable is located
func getExecutableDir() string {
	// Try to get executable path
	exePath, err := os.Executable()
	if err == nil {
		return filepath.Dir(exePath)
	}

	// Fallback to current working directory
	cwd, err := os.Getwd()
	if err == nil {
		return cwd
	}

	// Last resort fallback
	return "."
}

// SaveGame saves the current game state
func (m *Manager) SaveGame(data *entities.SaveData) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.saveFile, jsonData, 0644)
}

// LoadGame loads a saved game state
func (m *Manager) LoadGame() (*entities.SaveData, error) {
	data, err := os.ReadFile(m.saveFile)
	if err != nil {
		return nil, err
	}

	var saveData entities.SaveData
	if err := json.Unmarshal(data, &saveData); err != nil {
		return nil, err
	}

	return &saveData, nil
}

// HasSavedGame checks if a saved game exists
func (m *Manager) HasSavedGame() bool {
	_, err := os.Stat(m.saveFile)
	return err == nil
}

// DeleteSave removes the save file
func (m *Manager) DeleteSave() error {
	return os.Remove(m.saveFile)
}

// SaveLeaderboard saves the leaderboard
func (m *Manager) SaveLeaderboard(leaderboard *entities.Leaderboard) error {
	jsonData, err := json.MarshalIndent(leaderboard, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.leaderboardFile, jsonData, 0644)
}

// LoadLeaderboard loads the leaderboard
func (m *Manager) LoadLeaderboard() (*entities.Leaderboard, error) {
	data, err := os.ReadFile(m.leaderboardFile)
	if err != nil {
		if os.IsNotExist(err) {
			return entities.NewLeaderboard(), nil
		}
		return nil, err
	}

	var leaderboard entities.Leaderboard
	if err := json.Unmarshal(data, &leaderboard); err != nil {
		return nil, err
	}

	return &leaderboard, nil
}

// AddToLeaderboard adds a result to the leaderboard and saves
func (m *Manager) AddToLeaderboard(result entities.SessionResult) error {
	leaderboard, err := m.LoadLeaderboard()
	if err != nil {
		leaderboard = entities.NewLeaderboard()
	}

	leaderboard.AddResult(result)
	return m.SaveLeaderboard(leaderboard)
}
