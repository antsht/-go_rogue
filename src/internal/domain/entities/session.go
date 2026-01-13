package entities

import "time"

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StateInventory
	StatePaused
	StateGameOver
	StateVictory
	StateLeaderboard
)

// Session represents a game session
type Session struct {
	ID           string     `json:"id"`
	Character    *Character `json:"character"`
	CurrentLevel int        `json:"current_level"`
	Level        *Level     `json:"level"`
	State        GameState  `json:"state"`
	TurnCount    int        `json:"turn_count"`
	StartTime    time.Time  `json:"start_time"`
	LastSaveTime time.Time  `json:"last_save_time"`
	Messages     []string   `json:"messages"`
	MaxMessages  int        `json:"-"`

	// For item selection UI
	SelectingItem     bool     `json:"selecting_item"`
	SelectingItemType ItemType `json:"selecting_item_type"`

	// Dynamic difficulty (bonus feature)
	DifficultyModifier float64 `json:"difficulty_modifier"`
	RecentDeaths       int     `json:"recent_deaths"`
	RecentEasyKills    int     `json:"recent_easy_kills"`
}

// NewSession creates a new game session
func NewSession() *Session {
	return &Session{
		ID:                 generateSessionID(),
		Character:          NewCharacter(),
		CurrentLevel:       1,
		State:              StatePlaying,
		TurnCount:          0,
		StartTime:          time.Now(),
		Messages:           make([]string, 0),
		MaxMessages:        5,
		DifficultyModifier: 1.0,
	}
}

// generateSessionID creates a unique session identifier
func generateSessionID() string {
	return time.Now().Format("20060102150405")
}

// AddMessage adds a game message
func (s *Session) AddMessage(msg string) {
	s.Messages = append(s.Messages, msg)
	if len(s.Messages) > s.MaxMessages {
		s.Messages = s.Messages[1:]
	}
}

// GetMessages returns recent messages
func (s *Session) GetMessages() []string {
	return s.Messages
}

// ClearMessages clears all messages
func (s *Session) ClearMessages() {
	s.Messages = make([]string, 0)
}

// IncrementTurn increments the turn counter
func (s *Session) IncrementTurn() {
	s.TurnCount++
}

// IsGameOver returns true if the game has ended
func (s *Session) IsGameOver() bool {
	return s.State == StateGameOver || s.State == StateVictory
}

// SetGameOver sets the game over state
func (s *Session) SetGameOver() {
	s.State = StateGameOver
}

// SetVictory sets the victory state
func (s *Session) SetVictory() {
	s.State = StateVictory
}

// GetResult returns the session result for leaderboard
func (s *Session) GetResult() SessionResult {
	return SessionResult{
		SessionID:       s.ID,
		LevelReached:    s.CurrentLevel,
		GoldCollected:   s.Character.Gold,
		EnemiesDefeated: s.Character.Stats.EnemiesDefeated,
		FoodConsumed:    s.Character.Stats.FoodConsumed,
		ElixirsDrunk:    s.Character.Stats.ElixirsDrunk,
		ScrollsRead:     s.Character.Stats.ScrollsRead,
		HitsDealt:       s.Character.Stats.HitsDealt,
		HitsReceived:    s.Character.Stats.HitsReceived,
		TilesTraveled:   s.Character.Stats.TilesTraveled,
		TurnCount:       s.TurnCount,
		Victory:         s.State == StateVictory,
		Timestamp:       time.Now(),
	}
}

// SessionResult represents the result of a completed session
type SessionResult struct {
	SessionID       string    `json:"session_id"`
	LevelReached    int       `json:"level_reached"`
	GoldCollected   int       `json:"gold_collected"`
	EnemiesDefeated int       `json:"enemies_defeated"`
	FoodConsumed    int       `json:"food_consumed"`
	ElixirsDrunk    int       `json:"elixirs_drunk"`
	ScrollsRead     int       `json:"scrolls_read"`
	HitsDealt       int       `json:"hits_dealt"`
	HitsReceived    int       `json:"hits_received"`
	TilesTraveled   int       `json:"tiles_traveled"`
	TurnCount       int       `json:"turn_count"`
	Victory         bool      `json:"victory"`
	Timestamp       time.Time `json:"timestamp"`
}

// SaveData represents all data needed to save/load a game
type SaveData struct {
	Session       *Session `json:"session"`
	LevelSeed     int64    `json:"level_seed"`
	AllLevelSeeds []int64  `json:"all_level_seeds"`
}

// Leaderboard represents the game leaderboard
type Leaderboard struct {
	Results []SessionResult `json:"results"`
}

// NewLeaderboard creates a new empty leaderboard
func NewLeaderboard() *Leaderboard {
	return &Leaderboard{
		Results: make([]SessionResult, 0),
	}
}

// AddResult adds a result to the leaderboard
func (l *Leaderboard) AddResult(result SessionResult) {
	l.Results = append(l.Results, result)
	// Sort by gold collected (descending)
	l.sortByGold()
}

// sortByGold sorts results by gold collected in descending order
func (l *Leaderboard) sortByGold() {
	for i := 0; i < len(l.Results)-1; i++ {
		for j := i + 1; j < len(l.Results); j++ {
			if l.Results[j].GoldCollected > l.Results[i].GoldCollected {
				l.Results[i], l.Results[j] = l.Results[j], l.Results[i]
			}
		}
	}
}

// GetTopResults returns the top N results
func (l *Leaderboard) GetTopResults(n int) []SessionResult {
	if n > len(l.Results) {
		n = len(l.Results)
	}
	return l.Results[:n]
}
