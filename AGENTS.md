This repository contains a small ebiten viewer written in Go.

## Development
- Format Go code with `gofmt -w` before committing.
- Run tests with the `test` build tag so the ebiten game code is excluded:
  ```bash
  gofmt -w *.go
  go test -tags test ./...
  ```
- The assets directory is normally empty. If you need the image assets
  run `./scripts/download_assets.sh`.

Additional tasks to keep the project healthy:
 - Keep `README.md` up to date when features or build steps change.
 - Refresh `screenshot.png` if the UI changes significantly.
 - Run `go mod tidy` after adding dependencies.

