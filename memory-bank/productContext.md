Why this project exists
- Bootcamp team project to implement a classic Rogue-like game in Go.
- Demonstrates clean architecture layering, procedural generation, and
  terminal UI rendering.

Problems it solves
- Provides a complete, playable roguelike with save/continue and stats.
- Exercises game state management, AI, and rendering in a constrained
  terminal environment.

How it should work
- Player explores 21 procedurally generated levels, finds exits, and fights
  enemies in a turn-based loop.
- Items are picked up automatically and used via inventory/selection UI.
- Fog of war reveals rooms/corridors based on visibility rules.
- Runs are recorded in a leaderboard sorted by gold collected.

User experience goals
- Responsive keyboard controls with clear feedback messages.
- Readable, stable terminal UI with status bar and inventory views.
- Ability to continue from the last saved session.
