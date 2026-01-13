package game

import (
	"math"

	"github.com/user/go-rogue/internal/domain/entities"
)

// Visibility handles fog of war and line of sight calculations
type Visibility struct{}

// NewVisibility creates a new visibility handler
func NewVisibility() *Visibility {
	return &Visibility{}
}

// Update updates visibility based on player position
func (v *Visibility) Update(level *entities.Level, playerPos entities.Position) {
	// Clear all visibility
	level.ClearVisibility()

	// Find which room or corridor player is in
	room := level.GetRoomAt(playerPos)
	corridor := level.GetCorridorAt(playerPos)

	if room != nil {
		// Player is in a room - reveal entire room
		v.revealRoom(level, room)
		level.PlayerRoom = room

		// Also reveal adjacent corridor entrances
		for _, entrance := range room.Entrances {
			v.castRaysFromPoint(level, entrance, 3)
		}
	} else if corridor != nil {
		// Player is in a corridor - use ray casting
		corridor.Explored = true
		v.castRaysFromPoint(level, playerPos, 8)

		// Check if near a room entrance
		for _, roomPtr := range level.Rooms {
			for _, entrance := range roomPtr.Entrances {
				if playerPos.Distance(entrance) <= 1 {
					// Near entrance - cast rays into room
					v.castRaysIntoRoom(level, playerPos, roomPtr)
				}
			}
		}
	} else {
		// Fallback - just reveal immediate area
		v.revealRadius(level, playerPos, 2)
	}
}

// revealRoom reveals all tiles in a room
func (v *Visibility) revealRoom(level *entities.Level, room *entities.Room) {
	room.Explored = true

	for y := room.Y; y < room.Y+room.Height; y++ {
		for x := room.X; x < room.X+room.Width; x++ {
			pos := entities.Position{X: x, Y: y}
			level.MarkVisible(pos, true)
		}
	}
}

// revealRadius reveals tiles within a radius
func (v *Visibility) revealRadius(level *entities.Level, center entities.Position, radius int) {
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			pos := center.Add(dx, dy)
			if level.IsInBounds(pos) {
				dist := int(math.Sqrt(float64(dx*dx + dy*dy)))
				if dist <= radius {
					level.MarkVisible(pos, true)
				}
			}
		}
	}
}

// castRaysFromPoint casts rays in all directions from a point
func (v *Visibility) castRaysFromPoint(level *entities.Level, origin entities.Position, distance int) {
	// Cast rays in 360 degrees
	for angle := 0.0; angle < 360.0; angle += 5.0 {
		v.castRay(level, origin, angle, distance)
	}
}

// castRay casts a single ray using Bresenham's algorithm
func (v *Visibility) castRay(level *entities.Level, origin entities.Position, angle float64, maxDist int) {
	// Convert angle to radians
	rad := angle * math.Pi / 180.0

	// Calculate direction
	dx := math.Cos(rad)
	dy := math.Sin(rad)

	// Cast ray
	for i := 0; i <= maxDist; i++ {
		x := origin.X + int(math.Round(float64(i)*dx))
		y := origin.Y + int(math.Round(float64(i)*dy))
		pos := entities.Position{X: x, Y: y}

		if !level.IsInBounds(pos) {
			break
		}

		level.MarkVisible(pos, true)

		// Stop at walls
		tile := level.GetTile(pos)
		if tile != nil && tile.Type == entities.TileWall {
			break
		}
	}
}

// castRaysIntoRoom casts rays from entrance into a room using Bresenham
func (v *Visibility) castRaysIntoRoom(level *entities.Level, playerPos entities.Position, room *entities.Room) {
	// Get room center for direction
	center := room.GetCenter()

	// Calculate angle to room center
	dx := float64(center.X - playerPos.X)
	dy := float64(center.Y - playerPos.Y)
	baseAngle := math.Atan2(dy, dx) * 180.0 / math.Pi

	// Cast rays in a cone toward the room
	for angle := baseAngle - 45; angle <= baseAngle+45; angle += 2 {
		v.castRayIntoRoom(level, playerPos, angle, room)
	}
}

// castRayIntoRoom casts a ray into a specific room
func (v *Visibility) castRayIntoRoom(level *entities.Level, origin entities.Position, angle float64, room *entities.Room) {
	rad := angle * math.Pi / 180.0
	dx := math.Cos(rad)
	dy := math.Sin(rad)

	// Cast ray until we hit a wall or leave the room
	maxDist := room.Width + room.Height
	for i := 0; i <= maxDist; i++ {
		x := origin.X + int(math.Round(float64(i)*dx))
		y := origin.Y + int(math.Round(float64(i)*dy))
		pos := entities.Position{X: x, Y: y}

		if !level.IsInBounds(pos) {
			break
		}

		level.MarkVisible(pos, true)

		tile := level.GetTile(pos)
		if tile == nil {
			break
		}

		// Stop at walls (but reveal them)
		if tile.Type == entities.TileWall {
			// Check if this is the far wall of the room (not entrance wall)
			if !room.ContainsIncludingWalls(pos) {
				break
			}
		}
	}
}

// BresenhamLine returns points on a line between two positions
func BresenhamLine(start, end entities.Position) []entities.Position {
	var points []entities.Position

	x0, y0 := start.X, start.Y
	x1, y1 := end.X, end.Y

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	sx := 1
	if x0 > x1 {
		sx = -1
	}

	sy := 1
	if y0 > y1 {
		sy = -1
	}

	err := dx - dy

	for {
		points = append(points, entities.Position{X: x0, Y: y0})

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err

		if e2 > -dy {
			err -= dy
			x0 += sx
		}

		if e2 < dx {
			err += dx
			y0 += sy
		}
	}

	return points
}

// abs returns absolute value
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
