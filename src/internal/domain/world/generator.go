package world

import (
	"math/rand"

	"github.com/user/go-rogue/internal/domain/entities"
)

const (
	MinRoomWidth  = 6
	MaxRoomWidth  = 12
	MinRoomHeight = 4
	MaxRoomHeight = 8
)

// Generator handles procedural level generation
type Generator struct {
	rng          *rand.Rand
	doorSystem   *DoorKeySystem
	mimicPlacer  *MimicPlacer
}

// NewGenerator creates a new level generator
func NewGenerator() *Generator {
	return &Generator{
		doorSystem:  NewDoorKeySystem(),
		mimicPlacer: NewMimicPlacer(),
	}
}

// Generate creates a new level with the given seed
func (g *Generator) Generate(levelNum int, seed int64, difficultyMod float64) *entities.Level {
	g.rng = rand.New(rand.NewSource(seed))

	level := entities.NewLevel(levelNum)

	// Generate 9 rooms in a 3x3 grid
	rooms := g.generateRooms(level)

	// Connect rooms with corridors
	g.connectRooms(level, rooms)

	// Place rooms on the tile map
	g.placeRoomsOnMap(level)

	// Place corridors on the tile map
	g.placeCorridorsOnMap(level)

	// Select start and exit rooms
	g.selectSpecialRooms(level)

	// Place enemies (not in start room)
	g.placeEnemies(level, levelNum, difficultyMod)

	// Place items (not in start room)
	g.placeItems(level, levelNum, difficultyMod)

	// Bonus Task 6: Add doors and keys (only on some levels)
	if levelNum > 2 && g.rng.Float64() < 0.5 {
		g.doorSystem.AddDoorsAndKeys(level, seed+1)
	}

	// Bonus Task 8: Add mimics
	g.mimicPlacer.PlaceMimics(level, levelNum, seed+2)

	return level
}

// generateRooms creates rooms in a 3x3 grid
func (g *Generator) generateRooms(level *entities.Level) []*entities.Room {
	rooms := make([]*entities.Room, 0, 9)

	for gridY := 0; gridY < 3; gridY++ {
		for gridX := 0; gridX < 3; gridX++ {
			room := g.generateRoom(len(rooms), gridX, gridY)
			rooms = append(rooms, room)
			level.AddRoom(room)
		}
	}

	return rooms
}

// generateRoom creates a single room in the specified grid cell
func (g *Generator) generateRoom(id, gridX, gridY int) *entities.Room {
	// Calculate the section bounds
	sectionX := gridX * entities.SectionWidth
	sectionY := gridY * entities.SectionHeight

	// Random room size
	width := MinRoomWidth + g.rng.Intn(MaxRoomWidth-MinRoomWidth+1)
	height := MinRoomHeight + g.rng.Intn(MaxRoomHeight-MinRoomHeight+1)

	// Random position within section (with padding)
	maxX := entities.SectionWidth - width - 2
	maxY := entities.SectionHeight - height - 2
	if maxX < 1 {
		maxX = 1
	}
	if maxY < 1 {
		maxY = 1
	}

	x := sectionX + 1 + g.rng.Intn(maxX)
	y := sectionY + 1 + g.rng.Intn(maxY)

	return entities.NewRoom(id, x, y, width, height, gridX, gridY)
}

// connectRooms creates corridors between adjacent rooms
func (g *Generator) connectRooms(level *entities.Level, rooms []*entities.Room) {
	corridorID := 0

	// Connect horizontally adjacent rooms
	for gridY := 0; gridY < 3; gridY++ {
		for gridX := 0; gridX < 2; gridX++ {
			room1 := rooms[gridY*3+gridX]
			room2 := rooms[gridY*3+gridX+1]
			corridor := g.createCorridor(corridorID, room1, room2, true)
			level.AddCorridor(corridor)
			corridorID++
		}
	}

	// Connect vertically adjacent rooms
	for gridY := 0; gridY < 2; gridY++ {
		for gridX := 0; gridX < 3; gridX++ {
			room1 := rooms[gridY*3+gridX]
			room2 := rooms[(gridY+1)*3+gridX]
			corridor := g.createCorridor(corridorID, room1, room2, false)
			level.AddCorridor(corridor)
			corridorID++
		}
	}
}

// createCorridor creates a corridor between two rooms
func (g *Generator) createCorridor(id int, room1, room2 *entities.Room, horizontal bool) *entities.Corridor {
	corridor := entities.NewCorridor(id, room1.ID, room2.ID)

	var start, end entities.Position

	if horizontal {
		// Horizontal corridor (room1 is left of room2)
		// Start from right wall of room1
		startX := room1.X + room1.Width - 1
		startY := room1.Y + 1 + g.rng.Intn(room1.Height-2)

		// End at left wall of room2
		endX := room2.X
		endY := room2.Y + 1 + g.rng.Intn(room2.Height-2)

		start = entities.Position{X: startX, Y: startY}
		end = entities.Position{X: endX, Y: endY}

		// Record entrances
		room1.AddEntrance(start)
		room2.AddEntrance(end)

		// Generate L-shaped corridor
		midX := (startX + endX) / 2

		// First horizontal segment
		for x := startX; x <= midX; x++ {
			corridor.AddPoint(entities.Position{X: x, Y: startY})
		}
		// Vertical segment
		if startY < endY {
			for y := startY; y <= endY; y++ {
				corridor.AddPoint(entities.Position{X: midX, Y: y})
			}
		} else {
			for y := startY; y >= endY; y-- {
				corridor.AddPoint(entities.Position{X: midX, Y: y})
			}
		}
		// Second horizontal segment
		for x := midX; x <= endX; x++ {
			corridor.AddPoint(entities.Position{X: x, Y: endY})
		}
	} else {
		// Vertical corridor (room1 is above room2)
		// Start from bottom wall of room1
		startX := room1.X + 1 + g.rng.Intn(room1.Width-2)
		startY := room1.Y + room1.Height - 1

		// End at top wall of room2
		endX := room2.X + 1 + g.rng.Intn(room2.Width-2)
		endY := room2.Y

		start = entities.Position{X: startX, Y: startY}
		end = entities.Position{X: endX, Y: endY}

		// Record entrances
		room1.AddEntrance(start)
		room2.AddEntrance(end)

		// Generate L-shaped corridor
		midY := (startY + endY) / 2

		// First vertical segment
		for y := startY; y <= midY; y++ {
			corridor.AddPoint(entities.Position{X: startX, Y: y})
		}
		// Horizontal segment
		if startX < endX {
			for x := startX; x <= endX; x++ {
				corridor.AddPoint(entities.Position{X: x, Y: midY})
			}
		} else {
			for x := startX; x >= endX; x-- {
				corridor.AddPoint(entities.Position{X: x, Y: midY})
			}
		}
		// Second vertical segment
		for y := midY; y <= endY; y++ {
			corridor.AddPoint(entities.Position{X: endX, Y: y})
		}
	}

	return corridor
}

// placeRoomsOnMap renders rooms onto the tile map
func (g *Generator) placeRoomsOnMap(level *entities.Level) {
	for _, room := range level.Rooms {
		// Draw walls and floor
		for y := room.Y; y < room.Y+room.Height; y++ {
			for x := room.X; x < room.X+room.Width; x++ {
				pos := entities.Position{X: x, Y: y}

				// Determine if wall or floor
				if x == room.X || x == room.X+room.Width-1 ||
					y == room.Y || y == room.Y+room.Height-1 {
					// Wall
					symbol := g.getWallSymbol(x, y, room)
					level.SetTile(pos, entities.TileWall, symbol)
				} else {
					// Floor
					level.SetTile(pos, entities.TileFloor, '.')
				}
			}
		}
	}
}

// getWallSymbol returns the appropriate wall character
func (g *Generator) getWallSymbol(x, y int, room *entities.Room) rune {
	isTop := y == room.Y
	isBottom := y == room.Y+room.Height-1
	isLeft := x == room.X
	isRight := x == room.X+room.Width-1

	// Corners
	if isTop && isLeft {
		return '┌'
	}
	if isTop && isRight {
		return '┐'
	}
	if isBottom && isLeft {
		return '└'
	}
	if isBottom && isRight {
		return '┘'
	}

	// Edges
	if isTop || isBottom {
		return '─'
	}
	if isLeft || isRight {
		return '│'
	}

	return '#'
}

// placeCorridorsOnMap renders corridors onto the tile map
func (g *Generator) placeCorridorsOnMap(level *entities.Level) {
	for _, corridor := range level.Corridors {
		for _, pos := range corridor.Points {
			tile := level.GetTile(pos)
			if tile == nil {
				continue
			}

			// Don't overwrite room floors
			if tile.Type == entities.TileFloor {
				continue
			}

			// At room entrances, make it an entrance tile
			isEntrance := false
			for _, room := range level.Rooms {
				if room.IsEntrance(pos) {
					isEntrance = true
					break
				}
			}

			if isEntrance {
				level.SetTile(pos, entities.TileEntrance, '\'')
			} else if tile.Type == entities.TileWall {
				// Corridor through wall - make entrance
				level.SetTile(pos, entities.TileEntrance, '\'')
			} else {
				level.SetTile(pos, entities.TileCorridor, '#')
			}
		}
	}
}

// selectSpecialRooms selects start and exit rooms
func (g *Generator) selectSpecialRooms(level *entities.Level) {
	// Start room is random corner
	corners := []int{0, 2, 6, 8}
	startIdx := corners[g.rng.Intn(len(corners))]
	level.StartRoom = startIdx
	level.Rooms[startIdx].IsStart = true

	// Exit room is opposite corner
	exitIdx := 8 - startIdx
	level.ExitRoom = exitIdx
	level.Rooms[exitIdx].IsExit = true

	// Place exit tile in exit room
	exitRoom := level.Rooms[exitIdx]
	exitPos := exitRoom.GetCenter()
	level.ExitPos = exitPos
	level.SetTile(exitPos, entities.TileExit, '%')
}

// placeEnemies places enemies in rooms
func (g *Generator) placeEnemies(level *entities.Level, levelNum int, difficultyMod float64) {
	// More enemies at deeper levels
	baseEnemies := 2 + levelNum/3
	maxEnemies := baseEnemies + g.rng.Intn(3)

	// Apply difficulty modifier
	maxEnemies = int(float64(maxEnemies) * difficultyMod)
	if maxEnemies < 1 {
		maxEnemies = 1
	}

	enemiesPlaced := 0

	for _, room := range level.Rooms {
		// Skip start room
		if room.IsStart {
			continue
		}

		// Random number of enemies per room
		roomEnemies := g.rng.Intn(3)
		if room.IsExit {
			roomEnemies++ // More enemies guarding exit
		}

		for i := 0; i < roomEnemies && enemiesPlaced < maxEnemies; i++ {
			enemy := entities.CreateEnemyForLevel(levelNum)

			// Random position in room
			pos := room.GetRandomFloorPosition(entities.NewRNG(g.rng.Int63()))
			
			// Make sure not on exit
			if pos.Equals(level.ExitPos) {
				continue
			}

			enemy.Position = pos
			room.AddEnemy(enemy)
			enemiesPlaced++
		}
	}
}

// placeItems places items in rooms
func (g *Generator) placeItems(level *entities.Level, levelNum int, difficultyMod float64) {
	// Fewer items at deeper levels
	baseItems := 8 - levelNum/4
	if baseItems < 2 {
		baseItems = 2
	}

	// Difficulty modifier affects item count (lower difficulty = more items)
	itemMultiplier := 2.0 - difficultyMod
	maxItems := int(float64(baseItems) * itemMultiplier)
	if maxItems < 1 {
		maxItems = 1
	}

	itemsPlaced := 0

	for _, room := range level.Rooms {
		// Skip start room
		if room.IsStart {
			continue
		}

		// Random number of items per room
		roomItems := g.rng.Intn(2) + 1

		for i := 0; i < roomItems && itemsPlaced < maxItems; i++ {
			item := g.generateItem(levelNum)
			if item == nil {
				continue
			}

			// Random position in room
			pos := room.GetRandomFloorPosition(entities.NewRNG(g.rng.Int63()))

			// Make sure not on exit or occupied
			if pos.Equals(level.ExitPos) || room.GetItemAt(pos) != nil {
				continue
			}

			item.Position = pos
			room.AddItem(item)
			itemsPlaced++
		}
	}
}

// generateItem creates a random item appropriate for the level
func (g *Generator) generateItem(levelNum int) *entities.Item {
	roll := g.rng.Intn(100)

	if roll < 35 {
		// Food (35%)
		subtypes := []entities.ItemSubtype{
			entities.SubtypeRation,
			entities.SubtypeFruit,
			entities.SubtypeMeat,
		}
		return entities.NewFood(subtypes[g.rng.Intn(len(subtypes))])
	} else if roll < 55 {
		// Elixir (20%)
		subtypes := []entities.ItemSubtype{
			entities.SubtypeStrengthElixir,
			entities.SubtypeDexterityElixir,
			entities.SubtypeHealthElixir,
		}
		return entities.NewElixir(subtypes[g.rng.Intn(len(subtypes))])
	} else if roll < 75 {
		// Scroll (20%)
		subtypes := []entities.ItemSubtype{
			entities.SubtypeStrengthScroll,
			entities.SubtypeDexterityScroll,
			entities.SubtypeHealthScroll,
		}
		return entities.NewScroll(subtypes[g.rng.Intn(len(subtypes))])
	} else {
		// Weapon (25%)
		// Better weapons at deeper levels
		var subtype entities.ItemSubtype
		if levelNum < 5 {
			subtype = entities.SubtypeDagger
		} else if levelNum < 10 {
			if g.rng.Intn(2) == 0 {
				subtype = entities.SubtypeDagger
			} else {
				subtype = entities.SubtypeSword
			}
		} else if levelNum < 15 {
			options := []entities.ItemSubtype{
				entities.SubtypeSword,
				entities.SubtypeMace,
			}
			subtype = options[g.rng.Intn(len(options))]
		} else {
			options := []entities.ItemSubtype{
				entities.SubtypeSword,
				entities.SubtypeMace,
				entities.SubtypeAxe,
			}
			subtype = options[g.rng.Intn(len(options))]
		}
		return entities.NewWeapon(subtype)
	}
}
