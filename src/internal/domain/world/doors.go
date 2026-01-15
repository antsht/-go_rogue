package world

import (
	"math/rand"

	"github.com/user/go-rogue/internal/domain/entities"
)

// DoorGenerator handles door placement and optional colored keys (Bonus Task 6)
type DoorGenerator struct{}

// NewDoorGenerator creates a new door generator
func NewDoorGenerator() *DoorGenerator {
	return &DoorGenerator{}
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

// doorKeyPair tracks a placed door and its corresponding key
type doorKeyPair struct {
	corridor *entities.Corridor
	doorPos  entities.Position
	color    DoorColor
	keyRoom  *entities.Room
	keyPos   entities.Position
}

// AddDoors adds doors to a level; if withKeys is true, doors are locked with keys
func (d *DoorGenerator) AddDoors(level *entities.Level, rng *rand.Rand, withKeys bool) {
	if rng == nil {
		return
	}

	if withKeys {
		d.addLockedDoorsWithKeys(level, rng)
		return
	}

	d.addUnlockedDoors(level, rng)
}

func (d *DoorGenerator) addUnlockedDoors(level *entities.Level, rng *rand.Rand) {
	// Target number of doors: 2-5 per level
	targetDoors := 2 + rng.Intn(4)

	corridors := d.shuffledCorridors(level, rng)
	placedDoors := 0

	for _, corridor := range corridors {
		if placedDoors >= targetDoors {
			break
		}

		// Skip corridors that are too short
		if len(corridor.Points) <= 2 {
			continue
		}

		// 70% chance to try adding a door to this corridor
		if rng.Float64() > 0.7 {
			continue
		}

		// Determine door position (middle of corridor)
		midIdx := len(corridor.Points) / 2
		doorPos := corridor.Points[midIdx]

		corridor.AddDoor(doorPos, "", 0)
		tile := level.GetTile(doorPos)
		if tile != nil {
			tile.Type = entities.TileDoor
			tile.DoorColor = ""
			tile.DoorLocked = false
			tile.DoorKeyType = 0
			tile.Symbol = '+'
		}

		placedDoors++
	}
}

func (d *DoorGenerator) addLockedDoorsWithKeys(level *entities.Level, rng *rand.Rand) {

	startRoom := level.GetStartRoom()
	if startRoom == nil {
		return
	}

	// Target number of doors: 2-5 per level
	targetDoors := 2 + rng.Intn(4)

	// Shuffle corridors for random selection
	corridors := d.shuffledCorridors(level, rng)

	// Track placed doors and their keys
	placedPairs := make([]doorKeyPair, 0)
	// Track which key types are already used (to avoid duplicates)
	usedColors := make(map[entities.ItemSubtype]bool)

	// Try to place doors one by one, ensuring each has an accessible key
	colorIdx := 0
	for _, corridor := range corridors {
		if len(placedPairs) >= targetDoors {
			break
		}

		// Skip corridors that are too short
		if len(corridor.Points) <= 2 {
			continue
		}

		// 70% chance to try adding a door to this corridor
		if rng.Float64() > 0.7 {
			continue
		}

		// Find a color we haven't used yet (to ensure variety)
		var selectedColor DoorColor
		found := false
		for attempts := 0; attempts < len(doorColors); attempts++ {
			candidate := doorColors[colorIdx%len(doorColors)]
			colorIdx++
			if !usedColors[candidate.KeyType] {
				selectedColor = candidate
				found = true
				break
			}
		}
		if !found {
			// All colors used, pick any
			selectedColor = doorColors[rng.Intn(len(doorColors))]
		}

		// Determine door position (middle of corridor)
		midIdx := len(corridor.Points) / 2
		doorPos := corridor.Points[midIdx]

		// CRITICAL: Find where we can place the key BEFORE committing to the door
		// Simulate the level state with existing doors
		simulatedKeys := make(map[entities.ItemSubtype]bool)
		for _, pair := range placedPairs {
			simulatedKeys[pair.color.KeyType] = true
		}

		// Get rooms accessible with current keys (before this new door)
		accessibleRooms := d.getAccessibleRoomsWithKeys(level, startRoom.ID, simulatedKeys)
		// Always include start room
		accessibleRooms = append([]*entities.Room{startRoom}, accessibleRooms...)

		// Remove duplicates
		seen := make(map[int]bool)
		uniqueRooms := make([]*entities.Room, 0)
		for _, room := range accessibleRooms {
			if !seen[room.ID] {
				seen[room.ID] = true
				uniqueRooms = append(uniqueRooms, room)
			}
		}
		accessibleRooms = uniqueRooms

		if len(accessibleRooms) == 0 {
			continue // Can't place key, skip this door
		}

		// Choose a room for the key
		keyRoom := accessibleRooms[rng.Intn(len(accessibleRooms))]
		keyPos := keyRoom.GetRandomFloorPosition(entities.NewRNG(rng.Int63()))

		// Ensure key doesn't overlap with exit
		attempts := 0
		for keyPos.Equals(level.ExitPos) && attempts < 10 {
			keyPos = keyRoom.GetRandomFloorPosition(entities.NewRNG(rng.Int63()))
			attempts++
		}

		// NOW commit: place the door
		corridor.AddDoor(doorPos, selectedColor.Name, selectedColor.KeyType)
		tile := level.GetTile(doorPos)
		if tile != nil {
			tile.Type = entities.TileDoor
			tile.DoorColor = selectedColor.Name
			tile.DoorLocked = true
			tile.DoorKeyType = selectedColor.KeyType
			tile.Symbol = '+'
		}

		// Place the key
		key := entities.NewKey(selectedColor.KeyType)
		key.Position = keyPos
		keyRoom.AddItem(key)

		// Track this pair
		placedPairs = append(placedPairs, doorKeyPair{
			corridor: corridor,
			doorPos:  doorPos,
			color:    selectedColor,
			keyRoom:  keyRoom,
			keyPos:   keyPos,
		})
		usedColors[selectedColor.KeyType] = true
	}

	// Final verification - ensure exit is reachable
	d.verifySolvable(level)
}

func (d *DoorGenerator) shuffledCorridors(level *entities.Level, rng *rand.Rand) []*entities.Corridor {
	corridors := make([]*entities.Corridor, len(level.Corridors))
	copy(corridors, level.Corridors)
	rng.Shuffle(len(corridors), func(i, j int) {
		corridors[i], corridors[j] = corridors[j], corridors[i]
	})
	return corridors
}

// getAccessibleRoomsWithKeys returns rooms accessible from start with given keys (BFS)
func (d *DoorGenerator) getAccessibleRoomsWithKeys(level *entities.Level, startRoomID int, keys map[entities.ItemSubtype]bool) []*entities.Room {
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

		// Find connected rooms through corridors (checking door access)
		for _, corridor := range level.Corridors {
			if corridor.FromRoom == roomID || corridor.ToRoom == roomID {
				// Check if corridor is blocked by locked door we can't open
				canPass := true
				for _, door := range corridor.Doors {
					if door.Locked {
						// Check if we have the key
						if keys == nil || !keys[door.KeyType] {
							canPass = false
							break
						}
					}
				}

				if canPass {
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
func (d *DoorGenerator) verifySolvable(level *entities.Level) {
	startRoom := level.GetStartRoom()
	if startRoom == nil {
		return
	}

	exitRoom := level.GetExitRoom()
	if exitRoom == nil {
		return
	}

	// Simulate without modifying actual door state
	collectedKeys := make(map[entities.ItemSubtype]bool)
	visitedRooms := make(map[int]bool)
	unlockedInSim := make(map[entities.ItemSubtype]bool)

	d.simulateKeyCollection(level, startRoom.ID, collectedKeys, visitedRooms, unlockedInSim)

	if !visitedRooms[exitRoom.ID] {
		// Exit not reachable - unlock all doors and remove all keys as fallback
		for _, corridor := range level.Corridors {
			for i := range corridor.Doors {
				corridor.Doors[i].Locked = false
				tile := level.GetTile(corridor.Doors[i].Position)
				if tile != nil {
					tile.DoorLocked = false
				}
			}
		}
		// Remove all keys since all doors are now unlocked
		d.removeAllKeys(level)
	}
}

// simulateKeyCollection simulates key collection without modifying door states
func (d *DoorGenerator) simulateKeyCollection(level *entities.Level, roomID int, keys map[entities.ItemSubtype]bool, visited map[int]bool, unlocked map[entities.ItemSubtype]bool) {
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

		// Check if we can pass through (simulate having keys)
		canPass := true
		for _, door := range corridor.Doors {
			if door.Locked && !unlocked[door.KeyType] {
				if keys[door.KeyType] {
					// Mark as unlocked in simulation (don't modify actual door)
					unlocked[door.KeyType] = true
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
			d.simulateKeyCollection(level, nextRoom, keys, visited, unlocked)
		}
	}
}

// removeAllKeys removes all keys from all rooms
func (d *DoorGenerator) removeAllKeys(level *entities.Level) {
	for _, room := range level.Rooms {
		newItems := make([]*entities.Item, 0)
		for _, item := range room.Items {
			if item.Type != entities.ItemTypeKey {
				newItems = append(newItems, item)
			}
		}
		room.Items = newItems
	}
}
