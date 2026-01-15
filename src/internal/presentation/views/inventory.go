package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/user/go-rogue/internal/domain/entities"
	"github.com/user/go-rogue/internal/domain/game"
	"github.com/user/go-rogue/internal/presentation/renderer"
)

// InventoryViewRender renders the inventory view
type InventoryViewRender struct {
	screen     *renderer.Screen
	gameEngine *game.Engine
}

// NewInventoryViewRender creates a new inventory view renderer
func NewInventoryViewRender(screen *renderer.Screen, gameEngine *game.Engine) *InventoryViewRender {
	return &InventoryViewRender{
		screen:     screen,
		gameEngine: gameEngine,
	}
}

// Render draws the inventory view
func (v *InventoryViewRender) Render() {
	session := v.gameEngine.GetSession()
	if session == nil {
		return
	}

	width, height := v.screen.Size()
	char := session.Character
	backpack := char.Backpack

	// Get offset for centering the game area
	offsetX, offsetY := v.screen.GetGameAreaOffset()

	// Draw title centered
	title := "═══ INVENTORY ═══"
	v.screen.DrawString(width/2-len(title)/2, offsetY+1, title, tcell.ColorYellow, tcell.ColorBlack)

	// Draw character stats
	statsY := offsetY + 3
	v.screen.DrawString(offsetX+2, statsY, "CHARACTER STATS", tcell.ColorOrange, tcell.ColorBlack)
	v.screen.DrawString(offsetX+2, statsY+1, "────────────────", tcell.ColorOrange, tcell.ColorBlack)
	v.screen.DrawString(offsetX+2, statsY+2, "Health:    "+itoa(char.Health)+"/"+itoa(char.MaxHealth), tcell.ColorGreen, tcell.ColorBlack)
	v.screen.DrawString(offsetX+2, statsY+3, "Strength:  "+itoa(char.GetEffectiveStrength())+" ("+itoa(char.Strength)+")", tcell.ColorRed, tcell.ColorBlack)
	v.screen.DrawString(offsetX+2, statsY+4, "Dexterity: "+itoa(char.GetEffectiveDexterity())+" ("+itoa(char.Dexterity)+")", tcell.ColorTeal, tcell.ColorBlack)
	v.screen.DrawString(offsetX+2, statsY+5, "Armor:     "+itoa(char.Armor), tcell.ColorTeal, tcell.ColorBlack)
	v.screen.DrawString(offsetX+2, statsY+6, "Gold:      "+itoa(char.Gold), tcell.ColorYellow, tcell.ColorBlack)

	// Current weapon
	weaponStr := "None (Fists)"
	if char.Weapon != nil {
		weaponStr = char.Weapon.Name + " (+" + itoa(char.Weapon.Strength) + " ATK)"
	}
	v.screen.DrawString(offsetX+2, statsY+8, "Weapon:    "+weaponStr, tcell.ColorWhite, tcell.ColorBlack)

	// Draw keys section (keys are level-specific)
	v.renderKeysSection(offsetX+2, statsY+10, backpack.GetKeys())

	// Draw backpack sections
	sectionX := offsetX + 30
	sectionWidth := 25

	// Weapons section
	v.renderItemSection(sectionX, offsetY+3, "WEAPONS [h]", backpack.GetWeapons(), sectionWidth)

	// Food section
	v.renderItemSection(sectionX+sectionWidth+2, offsetY+3, "FOOD [j]", backpack.GetFood(), sectionWidth)

	// Elixirs section
	v.renderItemSection(sectionX, offsetY+16, "ELIXIRS [k]", backpack.GetElixirs(), sectionWidth)

	// Scrolls section
	v.renderItemSection(sectionX+sectionWidth+2, offsetY+16, "SCROLLS [e]", backpack.GetScrolls(), sectionWidth)

	// Instructions at bottom of game area
	instructY := offsetY + 28
	if instructY > height-3 {
		instructY = height - 3
	}
	v.screen.DrawString(offsetX+2, instructY, "Press [H/J/K/E] to use items", tcell.ColorGray, tcell.ColorBlack)
	v.screen.DrawString(offsetX+2, instructY+1, "Press [I], [Q] or [Backspace] to close", tcell.ColorGray, tcell.ColorBlack)
}

// renderItemSection renders a section of items
func (v *InventoryViewRender) renderItemSection(x, y int, title string, items []*entities.Item, width int) {
	v.screen.DrawString(x, y, title, tcell.ColorOrange, tcell.ColorBlack)
	v.screen.DrawString(x, y+1, "────────────────", tcell.ColorOrange, tcell.ColorBlack)

	if len(items) == 0 {
		v.screen.DrawString(x, y+2, "(empty)", tcell.ColorDarkGray, tcell.ColorBlack)
		return
	}

	for i, item := range items {
		if i >= 9 {
			break
		}
		line := "[" + string(rune('1'+i)) + "] " + item.Name + item.GetStatsString()
		if len(line) > width {
			line = line[:width-3] + "..."
		}
		v.screen.DrawString(x, y+2+i, line, tcell.ColorWhite, tcell.ColorBlack)
	}
}

// renderKeysSection renders the keys section with colored key symbols
func (v *InventoryViewRender) renderKeysSection(x, y int, keys []*entities.Item) {
	v.screen.DrawString(x, y, "KEYS (this level)", tcell.ColorOrange, tcell.ColorBlack)
	v.screen.DrawString(x, y+1, "────────────────", tcell.ColorOrange, tcell.ColorBlack)

	if len(keys) == 0 {
		v.screen.DrawString(x, y+2, "(none)", tcell.ColorDarkGray, tcell.ColorBlack)
		return
	}

	// Display keys in a row with their colors
	keyX := x
	for i, key := range keys {
		if i > 0 {
			keyX += 2 // Space between keys
		}

		// Get the color for this key
		keyColor := v.getKeyColor(key.Color)

		// Draw the key symbol
		v.screen.SetCell(keyX, y+2, 'k', keyColor, tcell.ColorBlack)
		keyX++
	}

	// Also show key names below
	for i, key := range keys {
		if i >= 4 {
			break // Max 4 keys displayed
		}
		keyColor := v.getKeyColor(key.Color)
		v.screen.DrawString(x, y+3+i, "• "+key.Name, keyColor, tcell.ColorBlack)
	}
}

// getKeyColor returns the tcell color for a key color name
func (v *InventoryViewRender) getKeyColor(colorName string) tcell.Color {
	switch colorName {
	case "red":
		return tcell.ColorRed
	case "blue":
		return tcell.ColorBlue
	case "green":
		return tcell.ColorGreen
	case "yellow":
		return tcell.ColorYellow
	default:
		return tcell.ColorWhite
	}
}

// itoa converts int to string
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

	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}
