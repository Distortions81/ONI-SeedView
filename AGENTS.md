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
- Always run `scripts/install_deps.sh` after cloning the repository or when your
  environment might be missing build tools. It installs Go, Xvfb and the `webp`
  utilities so that you can compile and test the project without additional
  steps.
- If you install additional packages while working, add them to
  `scripts/install_deps.sh` so the setup process stays comprehensive.

Additional tasks to keep the project healthy:
 - Keep `README.md` up to date when features or build steps change.
 - Update the help menu text in `main.go` whenever controls change and keep it
   in sync with the **Controls** section of `README.md`.
 - Refresh `screenshot.png` if the UI changes significantly.
 - Run `go mod tidy` after adding dependencies.
