Tech stack
- Language: Go (module `github.com/user/go-rogue`, go 1.25.5)
- UI: `github.com/gdamore/tcell/v2`
- Data: JSON files for savegame and leaderboard

Project layout (under `src/`)
- `cmd/rogue/main.go` entry point
- `internal/domain/` entities, game logic, world generation
- `internal/presentation/` renderer, views, input handling
- `internal/data/manager.go` persistence manager

Build/run
- From `src/`: `go build -o rogue.exe ./cmd/rogue`
- Run: `./rogue.exe`

Runtime data
- Save and leaderboard files stored in executable directory via `data.Manager`.

Testing
- No automated tests found.
