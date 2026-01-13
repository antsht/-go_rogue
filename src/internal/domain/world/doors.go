package world

import (
	"math/rand"

	"github.com/user/go-rogue/internal/domain/entities"
)

// DoorKeySystem handles the colored door and key system (Bonus Task 6)
type DoorKeySystem struct {
	rng *rand.Rand
}

// NewDoorKeySystem creates a new door/key system
func NewDoorKeySystem() *DoorKeySystem {
	return &DoorKeySystem{}
}

// DoorColor represents different door colors
type DoorColor struct {
	Name    string
	KeyType entities.ItemSubtype
}

var doorColors = []DoorColor{
	{"red", entities.SubtypeRedKey},
	{"blue", entities.SubtypeBlueKey},
	{"green", entities.SubtypeGreenKey},
	{"yellow", entities.SubtypeYellowKey},
}

// AddDoorsAndKeys adds doors and keys to a level
func (d *DoorKeySystem) AddDoorsAndKeys(level *entities.Level, seed int64) {
	d.rng = rand.New(rand.NewSource(seed))

	// Select corridors to add doors to
	numDoors := 1 + d.rng.Intn(3) // 1-3 doors per level

	// Track which keys we need to place
	keysNeeded := make([]DoorColor, 0)

	// Add doors to random corridors
	doorsAdded := 0
	for _, corridor := range level.Corridors {
		if doorsAdded >= numDoors {
			break
		}

		// 40% chance to add a door to this corridor
		if d.rng.Float64() < 0.4 && len(corridor.Points) > 2 {
			// Place door in middle of corridor
			midIdx := len(corridor.Points) / 2
			doorPos := corridor.Points[midIdx]

			// Select random color
			color := doorColors[d.rng.Intn(len(doorColors))]

			// Add door to corridor
			corridor.AddDoor(doorPos, color.Name, color.KeyType)

			// Update tile
			tile := level.GetTile(doorPos)
			if tile != nil {
				tile.Type = entities.TileDoor
				tile.DoorColor = color.Name
				tile.DoorLocked = true
				tile.DoorKeyType = color.KeyType
				tile.Symbol = '+'
			}

			keysNeeded = append(keysNeeded, color)
			doorsAdded++
		}
	}

	// Place keys in accessible rooms
	d.placeKeys(level, keysNeeded)

	// Verify no softlocks
	d.verifySolvable(level)
}

// placeKeys places keys in rooms that are accessible without the key
func (d *DoorKeySystem) placeKeys(level *entities.Level, keysNeeded []DoorColor) {
	// Simple placement: put keys in rooms near the start
	startRoom := level.GetStartRoom()
	if startRoom == nil {
		return
	}

	// Get accessible rooms from start (BFS)
	accessibleRooms := d.getAccessibleRooms(level, startRoom.ID)

	for i, color := range keysNeeded {
		// Find a room to place the key
		// Prefer rooms that are accessible before the door
		if len(accessibleRooms) > 0 {
			roomIdx := i % len(accessibleRooms)
			room := accessibleRooms[roomIdx]

			// Create and place key
			key := entities.NewKey(color.KeyType)
			pos := room.GetRandomFloorPosition(entities.NewRNG(d.rng.Int63()))
			key.Position = pos
			room.AddItem(key)
		}
	}
}

// getAccessibleRooms returns rooms accessible from start without keys (BFS)
func (d *DoorKeySystem) getAccessibleRooms(level *entities.Level, startRoomID int) []*entities.Room {
	visited := make(map[int]bool)
	accessible := make([]*entities.Room, 0)
	queue := []int{startRoomID}

	for len(queue) > 0 {
		roomID := queue[0]
		queue = queue[1:]

		if visited[roomID] {
			continue
		}
		visited[roomID] = true

		room := level.GetRoomByID(roomID)
		if room != nil && !room.IsStart {
			accessible = append(accessible, room)
		}

		// Find connected rooms through unlocked corridors
		for _, corridor := range level.Corridors {
			if corridor.FromRoom == roomID || corridor.ToRoom == roomID {
				// Check if corridor is blocked by locked door
				hasLockedDoor := false
				for _, door := range corridor.Doors {
					if door.Locked {
						hasLockedDoor = true
						break
					}
				}

				if !hasLockedDoor {
					nextRoom := corridor.ToRoom
					if nextRoom == roomID {
						nextRoom = corridor.FromRoom
					}
					if !visited[nextRoom] {
						queue = append(queue, nextRoom)
					}
				}
			}
		}
	}

	return accessible
}

// verifySolvable verifies the level can be completed (no softlocks)
func (d *DoorKeySystem) verifySolvable(level *entities.Level) {
	// Collect all keys available from start
	startRoom := level.GetStartRoom()
	if startRoom == nil {
		return
	}

	// Simulate collecting keys and unlocking doors
	collectedKeys := make(map[entities.ItemSubtype]bool)
	visitedRooms := make(map[int]bool)

	d.collectKeysAndUnlock(level, startRoom.ID, collectedKeys, visitedRooms)

	// Check if exit is reachable
	exitRoom := level.GetExitRoom()
	if exitRoom == nil {
		return
	}

	if !visitedRooms[exitRoom.ID] {
		// Exit not reachable - unlock all doors as fallback
		for _, corridor := range level.Corridors {
			for i := range corridor.Doors {
				corridor.Doors[i].Locked = false
				// Update tile
				tile := level.GetTile(corridor.Doors[i].Position)
				if tile != nil {
					tile.DoorLocked = false
				}
			}
		}
	}
}

// collectKeysAndUnlock simulates key collection and door unlocking
func (d *DoorKeySystem) collectKeysAndUnlock(level *entities.Level, roomID int, keys map[entities.ItemSubtype]bool, visited map[int]bool) {
	if visited[roomID] {
		return
	}
	visited[roomID] = true

	room := level.GetRoomByID(roomID)
	if room == nil {
		return
	}

	// Collect keys in this room
	for _, item := range room.Items {
		if item.Type == entities.ItemTypeKey {
			keys[item.Subtype] = true
		}
	}

	// Try to traverse corridors
	for _, corridor := range level.Corridors {
		if corridor.FromRoom != roomID && corridor.ToRoom != roomID {
			continue
		}

		// Check if we can pass through
		canPass := true
		for i := range corridor.Doors {
			door := &corridor.Doors[i]
			if door.Locked {
				if keys[door.KeyType] {
					// Unlock the door
					door.Locked = false
					tile := level.GetTile(door.Position)
					if tile != nil {
						tile.DoorLocked = false
					}
				} else {
					canPass = false
				}
			}
		}

		if canPass {
			nextRoom := corridor.ToRoom
			if nextRoom == roomID {
				nextRoom = corridor.FromRoom
			}
			d.collectKeysAndUnlock(level, nextRoom, keys, visited)
		}
	}
}
