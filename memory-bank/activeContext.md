Current focus
- Core gameplay complete; project in maintenance/polishing phase.
- Only Task 9 (3D rendering mode) remains unimplemented.

Recent changes
- Updated fog of war raycasting system for better corridor visibility:
  - Rays now stop at empty tiles (corridor boundaries/void).
  - Rays now stop at locked doors (door visible but not beyond).
  - Reduced ray angle step from 5.0° to 1.5° for better coverage.
  - Improved room peek angle step from 2.0° to 1.0° for finer granularity.
- Added gold/treasure generation to item spawning (15% chance).
- Gold value scales with level depth: base `10 + levelNum*5` with variance.
- Increased maximum room dimensions: width 12→13, height 5→6.
- Increased MapHeight from 24→26 to accommodate larger rooms.
- Item distribution rebalanced: Treasure 15%, Food 30%, Elixir 15%, Scroll 15%, Weapon 25%.

Next steps
- Optional: Implement Task 9 (3D first-person rendering with ray casting).
- Consider adding automated tests for regression prevention.
- Monitor for balance issues or bugs during playtesting.
