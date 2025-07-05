This repository contains a small ebiten viewer written in Go.

## Development
- Binary files can not be committed due to current limitations, don't do this or you will not be able to commit changes.
- Format Go code with `gofmt -w` before committing.
- Run tests with the `test` build tag so the ebiten game code is excluded:
  ```bash
  gofmt -w *.go
  go test -tags test ./...
  ```
- Image assets are committed to the repo and embedded at build time, so no
  download step is required.
- Run `scripts/bump_version.sh` before committing to update `ClientVersion`
  in `const.go`. The script writes `v0.0.0-YYMMDDHHMM` using the current UTC
  time.

Additional tasks to keep the project healthy:
 - Keep `README.md` up to date when features or build steps change.
 - Refresh `screenshot.png` if the UI changes significantly.
 - Run `go mod tidy` after adding dependencies.
