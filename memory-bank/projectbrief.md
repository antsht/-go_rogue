Project Team 01 — Go Rogue

Goal: build a console-based roguelike inspired by Rogue (1980) in Go using
tcell, with clean architecture separation and core gameplay mechanics,
procedural generation, and persistence.

Core requirements:
- 21 dungeon levels, each a 3x3 grid of rooms connected by corridors
- Turn-based movement/combat and keyboard controls
- Fog of war rendering
- Items, enemies, and combat per Rogue-inspired rules
- Save/load to JSON and a leaderboard of runs

Scope constraints:
- Work inside `src/`
- Clean architecture separation: domain, presentation, data
