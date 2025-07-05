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
- Run `scripts/bump_version.sh` before committing to update the timestamp in
  the `ClientVersion` constant found in `const.go`. The script preserves the
  existing `MAJOR.MINOR.PATCH` numbers. Increase these version numbers
  manually when releasing new features or fixes and **never** decrease them.

Additional tasks to keep the project healthy:
 - Keep `README.md` up to date when features or build steps change.
 - Refresh `screenshot.png` if the UI changes significantly.
 - Run `go mod tidy` after adding dependencies.
