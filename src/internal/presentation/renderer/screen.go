package renderer

import (
	"github.com/gdamore/tcell/v2"
	"github.com/user/go-rogue/internal/domain/entities"
)

// Screen wraps tcell screen functionality
type Screen struct {
	screen tcell.Screen
	width  int
	height int

	// Color map for game colors
	colors map[string]tcell.Color
}

// NewScreen creates and initializes a new screen
func NewScreen() (*Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	// Set default style
	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

	width, height := screen.Size()

	s := &Screen{
		screen: screen,
		width:  width,
		height: height,
		colors: make(map[string]tcell.Color),
	}

	// Initialize color map
	s.initColors()

	return s, nil
}

// initColors sets up the color mapping
func (s *Screen) initColors() {
	s.colors["black"] = tcell.ColorBlack
	s.colors["white"] = tcell.ColorWhite
	s.colors["red"] = tcell.ColorRed
	s.colors["green"] = tcell.ColorGreen
	s.colors["blue"] = tcell.ColorBlue
	s.colors["yellow"] = tcell.ColorYellow
	s.colors["cyan"] = tcell.ColorTeal
	s.colors["magenta"] = tcell.ColorPurple
	s.colors["brown"] = tcell.ColorMaroon
	s.colors["orange"] = tcell.ColorOrange
	s.colors["gray"] = tcell.ColorGray
	s.colors["darkgray"] = tcell.ColorDarkGray
}

// GetColor returns a tcell color for a color name
func (s *Screen) GetColor(name string) tcell.Color {
	if color, ok := s.colors[name]; ok {
		return color
	}
	return tcell.ColorWhite
}

// Close cleans up the screen
func (s *Screen) Close() {
	s.screen.Fini()
}

// Clear clears the screen
func (s *Screen) Clear() {
	s.screen.Clear()
}

// Show updates the screen
func (s *Screen) Show() {
	s.screen.Show()
}

// Size returns the screen dimensions
func (s *Screen) Size() (int, int) {
	return s.width, s.height
}

// UpdateSize updates stored dimensions
func (s *Screen) UpdateSize() {
	s.width, s.height = s.screen.Size()
}

// PollEvent returns the next event
func (s *Screen) PollEvent() tcell.Event {
	return s.screen.PollEvent()
}

// SetCell sets a cell at position with given style
func (s *Screen) SetCell(x, y int, ch rune, fg, bg tcell.Color) {
	style := tcell.StyleDefault.Foreground(fg).Background(bg)
	s.screen.SetContent(x, y, ch, nil, style)
}

// SetCellStyle sets a cell with a named color
func (s *Screen) SetCellStyle(x, y int, ch rune, fgName, bgName string) {
	fg := s.GetColor(fgName)
	bg := s.GetColor(bgName)
	s.SetCell(x, y, ch, fg, bg)
}

// DrawString draws a string at position
func (s *Screen) DrawString(x, y int, str string, fg, bg tcell.Color) {
	col := 0
	for _, ch := range str {
		s.SetCell(x+col, y, ch, fg, bg)
		col++
	}
}

// DrawStringStyle draws a string with named colors
func (s *Screen) DrawStringStyle(x, y int, str string, fgName, bgName string) {
	fg := s.GetColor(fgName)
	bg := s.GetColor(bgName)
	s.DrawString(x, y, str, fg, bg)
}

// DrawBox draws a box at the specified position
func (s *Screen) DrawBox(x, y, width, height int, fg, bg tcell.Color) {
	// Corners
	s.SetCell(x, y, '┌', fg, bg)
	s.SetCell(x+width-1, y, '┐', fg, bg)
	s.SetCell(x, y+height-1, '└', fg, bg)
	s.SetCell(x+width-1, y+height-1, '┘', fg, bg)

	// Horizontal lines
	for i := x + 1; i < x+width-1; i++ {
		s.SetCell(i, y, '─', fg, bg)
		s.SetCell(i, y+height-1, '─', fg, bg)
	}

	// Vertical lines
	for i := y + 1; i < y+height-1; i++ {
		s.SetCell(x, i, '│', fg, bg)
		s.SetCell(x+width-1, i, '│', fg, bg)
	}
}

// DrawFilledBox draws a filled box
func (s *Screen) DrawFilledBox(x, y, width, height int, fg, bg tcell.Color, fillChar rune) {
	s.DrawBox(x, y, width, height, fg, bg)

	// Fill interior
	for dy := y + 1; dy < y+height-1; dy++ {
		for dx := x + 1; dx < x+width-1; dx++ {
			s.SetCell(dx, dy, fillChar, fg, bg)
		}
	}
}

// DrawLevel renders a dungeon level
func (s *Screen) DrawLevel(level *entities.Level, playerPos entities.Position) {
	for y := 0; y < entities.MapHeight; y++ {
		for x := 0; x < entities.MapWidth; x++ {
			tile := &level.Tiles[y][x]

			if !tile.Explored {
				s.SetCell(x, y, ' ', tcell.ColorBlack, tcell.ColorBlack)
				continue
			}

			fg := tcell.ColorWhite
			bg := tcell.ColorBlack
			ch := tile.Symbol

			switch tile.Type {
			case entities.TileWall:
				fg = tcell.ColorOrange
				bg = tcell.ColorBlack
			case entities.TileFloor:
				if tile.Visible {
					fg = tcell.ColorGreen
					ch = '.'
				} else {
					fg = tcell.ColorDarkGray
					ch = ' '
				}
			case entities.TileCorridor:
				fg = tcell.ColorWhite
				ch = '#'
			case entities.TileDoor:
				if tile.DoorLocked {
					fg = s.GetColor(tile.DoorColor)
					ch = '+'
				} else {
					fg = tcell.ColorWhite
					ch = '\''
				}
			case entities.TileExit:
				fg = tcell.ColorYellow
				ch = '%'
			case entities.TileEntrance:
				fg = tcell.ColorWhite
				ch = '\''
			}

			s.SetCell(x, y, ch, fg, bg)
		}
	}
}

// DrawCharacter draws the player character
func (s *Screen) DrawCharacter(pos entities.Position) {
	s.SetCell(pos.X, pos.Y, '@', tcell.ColorGreen, tcell.ColorBlack)
}

// DrawEnemy draws an enemy
func (s *Screen) DrawEnemy(enemy *entities.Enemy) {
	if !enemy.IsVisible {
		return
	}
	fg := s.GetColor(enemy.GetDisplayColor())
	s.SetCell(enemy.Position.X, enemy.Position.Y, enemy.GetDisplaySymbol(), fg, tcell.ColorBlack)
}

// DrawItem draws an item
func (s *Screen) DrawItem(item *entities.Item) {
	fg := s.GetColor(item.GetDisplayColor())
	s.SetCell(item.Position.X, item.Position.Y, item.GetDisplaySymbol(), fg, tcell.ColorBlack)
}

// DrawStatusBar draws the status bar at the bottom
func (s *Screen) DrawStatusBar(session *entities.Session) {
	y := entities.MapHeight
	char := session.Character

	// Clear status area (3 lines: stats + 2 message lines)
	for x := 0; x < s.width; x++ {
		s.SetCell(x, y, ' ', tcell.ColorWhite, tcell.ColorBlack)
		s.SetCell(x, y+1, ' ', tcell.ColorWhite, tcell.ColorBlack)
		s.SetCell(x, y+2, ' ', tcell.ColorWhite, tcell.ColorBlack)
	}

	// Format: Level:X  Hits:XX(XX)  Str:XX(XX)  Gold:XXX  Armor:X  Exp:X/XX
	status := []struct {
		label string
		value string
		color tcell.Color
	}{
		{"Level", string(rune('0' + session.CurrentLevel)), tcell.ColorWhite},
	}

	x := 0
	// Level
	s.DrawString(x, y, "Level:", tcell.ColorWhite, tcell.ColorBlack)
	x += 6
	s.DrawString(x, y, itoa(session.CurrentLevel), tcell.ColorYellow, tcell.ColorBlack)
	x += len(itoa(session.CurrentLevel)) + 4

	// Hits (current/max)
	s.DrawString(x, y, "Hits:", tcell.ColorWhite, tcell.ColorBlack)
	x += 5
	hitsStr := itoa(char.Health) + "(" + itoa(char.MaxHealth) + ")"
	hitsColor := tcell.ColorGreen
	if char.Health < char.MaxHealth/3 {
		hitsColor = tcell.ColorRed
	} else if char.Health < char.MaxHealth*2/3 {
		hitsColor = tcell.ColorYellow
	}
	s.DrawString(x, y, hitsStr, hitsColor, tcell.ColorBlack)
	x += len(hitsStr) + 4

	// Str
	s.DrawString(x, y, "Str:", tcell.ColorWhite, tcell.ColorBlack)
	x += 4
	strStr := itoa(char.GetEffectiveStrength()) + "(" + itoa(char.Strength) + ")"
	s.DrawString(x, y, strStr, tcell.ColorWhite, tcell.ColorBlack)
	x += len(strStr) + 4

	// Gold
	s.DrawString(x, y, "Gold:", tcell.ColorWhite, tcell.ColorBlack)
	x += 5
	s.DrawString(x, y, itoa(char.Gold), tcell.ColorYellow, tcell.ColorBlack)
	x += len(itoa(char.Gold)) + 4

	// Armor
	s.DrawString(x, y, "Armor:", tcell.ColorWhite, tcell.ColorBlack)
	x += 6
	s.DrawString(x, y, itoa(char.Armor), tcell.ColorTeal, tcell.ColorBlack)
	x += len(itoa(char.Armor)) + 4

	// Exp
	s.DrawString(x, y, "Exp:", tcell.ColorWhite, tcell.ColorBlack)
	x += 4
	expStr := itoa(char.Level) + "/" + itoa(char.Experience)
	s.DrawString(x, y, expStr, tcell.ColorPurple, tcell.ColorBlack)

	// Draw last two messages on status lines
	msgCount := len(session.Messages)
	if msgCount > 0 {
		// Show second-to-last message on first line (if exists)
		if msgCount > 1 {
			msg := session.Messages[msgCount-2]
			s.DrawString(0, y+1, msg, tcell.ColorGray, tcell.ColorBlack)
		}
		// Show last message on second line
		msg := session.Messages[msgCount-1]
		s.DrawString(0, y+2, msg, tcell.ColorWhite, tcell.ColorBlack)
	}

	// Ignore unused variable
	_ = status
}

// itoa converts int to string without importing strconv
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

	// Reverse
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}
