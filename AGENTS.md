This repository contains a small ebiten viewer written in Go.

## Development
- Binary files can not be committed due to current limitations, don't do this or you will not be able to commit changes.
- Format Go code with `gofmt -w` before committing.
- Run tests normally:
  ```bash
  gofmt -w *.go
  go test ./...
  ```
- Image assets are committed to the repo and embedded at build time, so no
  download step is required.
- Run `./install_deps.sh` and related scripts **only** when you need the packages they install. These helpers speed up environment setup when you require tools like Go, Xvfb or the `webp` utilities.
- If you install additional packages while working, add them to
  `scripts/install_deps.sh` so they are included in the setup process.

Additional tasks to keep the project healthy:
 - Keep `README.md` up to date when features or build steps change.
 - Update the help menu text in `main.go` whenever controls change and keep it
   in sync with the **Controls** section of `README.md`.
 - Refresh `screenshot.png` if the UI changes significantly.
 - Run `go mod tidy` after adding dependencies.
