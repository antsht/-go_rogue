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

	// Draw title
	title := "═══ INVENTORY ═══"
	v.screen.DrawString(width/2-len(title)/2, 1, title, tcell.ColorYellow, tcell.ColorBlack)

	// Draw character stats
	statsY := 3
	v.screen.DrawString(2, statsY, "CHARACTER STATS", tcell.ColorOrange, tcell.ColorBlack)
	v.screen.DrawString(2, statsY+1, "────────────────", tcell.ColorOrange, tcell.ColorBlack)
	v.screen.DrawString(2, statsY+2, "Health:    "+itoa(char.Health)+"/"+itoa(char.MaxHealth), tcell.ColorGreen, tcell.ColorBlack)
	v.screen.DrawString(2, statsY+3, "Strength:  "+itoa(char.GetEffectiveStrength())+" ("+itoa(char.Strength)+")", tcell.ColorRed, tcell.ColorBlack)
	v.screen.DrawString(2, statsY+4, "Dexterity: "+itoa(char.GetEffectiveDexterity())+" ("+itoa(char.Dexterity)+")", tcell.ColorTeal, tcell.ColorBlack)
	v.screen.DrawString(2, statsY+5, "Armor:     "+itoa(char.Armor), tcell.ColorTeal, tcell.ColorBlack)
	v.screen.DrawString(2, statsY+6, "Gold:      "+itoa(char.Gold), tcell.ColorYellow, tcell.ColorBlack)

	// Current weapon
	weaponStr := "None (Fists)"
	if char.Weapon != nil {
		weaponStr = char.Weapon.Name + " (+" + itoa(char.Weapon.Strength) + ")"
	}
	v.screen.DrawString(2, statsY+8, "Weapon:    "+weaponStr, tcell.ColorWhite, tcell.ColorBlack)

	// Draw backpack sections
	sectionX := 30
	sectionWidth := 20

	// Weapons section
	v.renderItemSection(sectionX, 3, "WEAPONS [h]", backpack.GetWeapons(), sectionWidth)

	// Food section
	v.renderItemSection(sectionX+sectionWidth+2, 3, "FOOD [j]", backpack.GetFood(), sectionWidth)

	// Elixirs section
	v.renderItemSection(sectionX, 14, "ELIXIRS [k]", backpack.GetElixirs(), sectionWidth)

	// Scrolls section
	v.renderItemSection(sectionX+sectionWidth+2, 14, "SCROLLS [e]", backpack.GetScrolls(), sectionWidth)

	// Instructions
	v.screen.DrawString(2, height-3, "Press [H/J/K/E] to use items", tcell.ColorGray, tcell.ColorBlack)
	v.screen.DrawString(2, height-2, "Press [I], [Q] or [Backspace] to close", tcell.ColorGray, tcell.ColorBlack)
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
		line := "[" + string(rune('1'+i)) + "] " + item.Name
		if len(line) > width {
			line = line[:width-3] + "..."
		}
		v.screen.DrawString(x, y+2+i, line, tcell.ColorWhite, tcell.ColorBlack)
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
