Architecture
- Clean layering: `cmd/` entry, `internal/domain` (entities/game/world),
  `internal/presentation` (renderer/input/views), `internal/data`.
- Domain is framework-agnostic; presentation uses tcell; data uses JSON files.

Game engine flow
- `Engine` owns session, world generator, combat, AI, visibility, difficulty.
- `NewGame` seeds levels, generates level 1, places player, saves initial state.
- `MovePlayer` handles combat, doors, pickups, exit progression, and turn
  processing.
- Turn processing updates effects, runs enemy AI, updates difficulty, and
  checks game over.

World generation
- 3x3 room grid with L-shaped corridors; walls/floors placed on a tile map.
- Map dimensions: 80×26 (SectionWidth=26, SectionHeight=8).
- Room sizes: width 6-13, height 4-6 (with 1-tile margins per section).
- Start room is random corner; exit is opposite corner with `%` tile.
- Item generation: Treasure 15%, Food 30%, Elixir 15%, Scroll 15%, Weapon 25%.
- Gold value scales with level: `10 + levelNum*5` with ±25% variance.
- Enemies and items placed per level depth; difficulty modifier influences counts.
- Bonus features: colored doors/keys (softlock checks via BFS/DFS) and mimics.
- Mimics spawn during `placeEnemies` with 20% chance, capped at `(levelNum/3)+1` (max 4).
- Mimic item appearance randomized via `RandomMimickedItem` (treasure 40%, weapon/scroll/elixir).

Combat & AI
- Hit chance derived from dexterity; damage from strength with variance.
- Enemy-specific behaviors: ghost teleport/invisibility, ogre rest/counter,
  snake-mage diagonal move/sleep, vampire first-hit miss & max HP drain.
- Mimics are enemies that render as items until attacked.

Visibility
- Rooms are fully revealed when entered; corridors use ray casting with
  Bresenham line-of-sight to lift fog of war near entrances.

Persistence
- Save data stores session plus per-level RNG seeds to regenerate levels
  deterministically on continue.
- Leaderboard is stored in JSON and sorted by gold collected.
