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
- Start room is random corner; exit is opposite corner with `%` tile.
- Enemies and items placed per level depth; difficulty modifier influences
  counts.
- Bonus features: colored doors/keys (softlock checks) and mimics.
- Mimics are spawned during standard enemy placement (with level-based caps).

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
