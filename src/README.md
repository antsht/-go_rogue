# Go Rogue - A Roguelike Game

A console-based roguelike game inspired by the classic 1980 Rogue game, implemented in Go using the tcell library.

## Features

### Core Gameplay (Tasks 0-5)
- **21 Dungeon Levels**: Progress through increasingly difficult dungeon levels
- **Procedural Generation**: Each level is randomly generated with 9 rooms in a 3x3 grid
- **Turn-Based Combat**: Strategic combat with hit chance based on dexterity
- **5 Enemy Types**:
  - **Zombie** (green `z`): High health, medium strength, slow
  - **Vampire** (red `v`): Drains max health, first hit always misses
  - **Ghost** (white `g`): Teleports, becomes invisible
  - **Ogre** (yellow `O`): Moves 2 tiles per turn, guaranteed counterattack after rest
  - **Snake-Mage** (white `s`): Diagonal movement, can put player to sleep
- **Item System**: Food, Elixirs (temporary buffs), Scrolls (permanent buffs), Weapons
- **Fog of War**: Ray casting visibility system
- **Save/Load**: JSON-based game persistence
- **Leaderboard**: Track your best runs sorted by gold collected

### Bonus Features (Tasks 6-8)
- **Colored Doors & Keys** (Task 6): DOOM-style colored key system with softlock prevention
- **Dynamic Difficulty** (Task 7): Game adjusts difficulty based on player performance
- **Mimic Enemy** (Task 8): Enemy that disguises itself as items

## Controls

### Movement
- `W` / `↑` - Move up
- `S` / `↓` - Move down
- `A` / `←` - Move left
- `D` / `→` - Move right

### Items
- `H` - Select weapon from backpack
- `J` - Use food from backpack
- `K` - Use elixir from backpack
- `E` - Use scroll from backpack
- `I` - Open inventory view

### Menu
- `N` - New game
- `C` - Continue saved game
- `L` - View leaderboard
- `Q` - Quit
- `ESC` - Return to menu / Cancel

## Building

```bash
cd src
go mod tidy
go build -o rogue.exe ./cmd/rogue
```

## Running

```bash
./rogue.exe
```

## Architecture

The game follows clean architecture principles:

```
src/
├── cmd/rogue/           # Application entry point
├── internal/
│   ├── domain/          # Business logic layer
│   │   ├── entities/    # Game entities (Character, Enemy, Item, etc.)
│   │   ├── game/        # Game mechanics (Combat, AI, Visibility)
│   │   └── world/       # Level generation
│   ├── presentation/    # UI layer
│   │   ├── renderer/    # tcell screen rendering
│   │   ├── input/       # Input handling
│   │   └── views/       # Different game views
│   └── data/            # Data persistence layer
└── go.mod
```

## Gameplay Tips

1. **Explore carefully**: Rooms are revealed when you enter them
2. **Manage resources**: Food heals, elixirs give temporary buffs, scrolls give permanent buffs
3. **Choose your battles**: Some enemies are better avoided at low levels
4. **Watch for Mimics**: At higher levels, that treasure might be a monster!
5. **Find the exit (%)**: Descend through all 21 levels to win

## Symbols

| Symbol | Description |
|--------|-------------|
| `@` | Player character |
| `%` | Level exit |
| `.` | Floor |
| `#` | Corridor |
| `+` | Locked door |
| `'` | Open door/entrance |
| `*` | Treasure |
| `:` | Food |
| `!` | Elixir |
| `?` | Scroll |
| `)` | Weapon |
| `k` | Key |

## License

Educational project - Team 01 Go Bootcamp
