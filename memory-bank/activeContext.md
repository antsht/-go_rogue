Current focus
- Core gameplay complete; project in maintenance/polishing phase.
- Only Task 9 (3D rendering mode) remains unimplemented.

Recent changes
- Mimic generation integrated into standard enemy placement in `generator.go`.
- Mimics spawn with 20% chance during enemy placement (capped per level).
- RNG-aware enemy factory (`CreateEnemyForLevelWithRNG`) and mimic helpers added.
- Level-based mimic caps: `(levelNum / 3) + 1` starting from level 2, max 4.

Next steps
- Optional: Implement Task 9 (3D first-person rendering with ray casting).
- Consider adding automated tests for regression prevention.
- Monitor for balance issues or bugs during playtesting.
