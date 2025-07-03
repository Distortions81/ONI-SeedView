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

