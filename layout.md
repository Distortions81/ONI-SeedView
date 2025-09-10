# Project Structure

This document explains how the Oni-SeedView repository is organized. It expands upon the brief listing in `README.md` and documents where the core logic and assets live.

## Top-Level Directories

- `objects/` – Image files for geysers, points of interest and other world objects. They are embedded into the binary via `//go:embed` in `assets.go` and loaded at runtime.
- `biomes/` – Textures for each biome. Each PNG is 256×256 pixels and is mapped to a biome name in `colors.go`. The mapping is documented in `BIOME_TEXTURES.md`.
- `icons/` – Toolbar icons such as the camera, help and gear images.
- `html/` – WebAssembly loader pages (`index.html` and `view.html`).
- `data/` – Runtime fonts. `NotoSansMono.ttf` is embedded by `fonts.go`.
- `scripts/` – Helper scripts used for building, headless execution and font subsetting.
- `biomes`, `objects` and `icons` images are referenced by name and embedded using Go’s `embed` package.

## Key Go Files

- `main.go` – Application entry point. Parses command line flags, sets up the `Game` struct and starts Ebiten.
- `update.go` – Main update loop handling input and camera movement.
- `drawing.go`, `draw.go` and `draw_helpers.go` – Rendering helpers for biomes, icons, text and the UI.
- `game_helpers.go` – Definition of `Game` plus many helper methods for layout and state management.
- `layout.go` – Ebiten `Layout` function which resizes the view and clamps camera bounds.
- `options_menu.go`, `screenshot_menu.go`, `asteroid_menu.go` – Implement the various drop‑down menus.
- `assets.go` – Loads and caches embedded images. Also converts filenames to the camel case used by some assets.
- `colors.go` and `const.go` – Color definitions and user‑interface constants.
- `parse.go` – Converts biome path strings into coordinate lists.
- `types.go` – Data structures for geysers, POIs and asteroids.
- `net.go` – Performs HTTP requests to `https://oni-data.stefanoltmann.de/COORDINATE` and decodes protobuf data via Go's `google.golang.org/protobuf`.
- `fonts.go` – Handles font loading and size adjustments.
- `text_draw.go`, `textutil.go` – Text rendering utilities.
- `touch_input.go`, `mobile_detect.go` – Touch gesture handling and simple mobile detection.
- `url.go`, `url_wasm.go` – Helpers for parsing query parameters on desktop vs. WASM.

## Data Flow and State

Seed information is downloaded by `fetchSeedProto` in `net.go` and decoded to the
`SeedData` struct via `decodeSeedProto`. Each `Asteroid` entry stores geysers, POIs
and biome polygon paths. Functions in `parse.go` convert those paths into
coordinate lists. The resulting slices are stored on the `Game` struct defined
in `game_helpers.go` along with runtime assets, camera coordinates and menu
flags.

Icons, biome textures and the embedded font are loaded in `assets.go` and
`fonts.go` and cached inside the `Game` instance so `Update` and `Draw` only
reference in-memory data.

## Assets

The viewer embeds all PNGs located under `objects/`, `icons/` and `biomes/`. When running on the web, the files are served relative to the page URL using the `WebAssetBase` constant from `const.go`. The textures are 256×256 pixels and mapped to biome colors as shown in `BIOME_TEXTURES.md`. A subset of object icons is loaded on demand after the seed data is fetched to keep startup time low.

The font file `data/NotoSansMono.ttf` is minimized to ASCII only using `scripts/minimize_font.py` to keep binary size small.

Screenshot output (`screenshot.png`) is stored in the repository to demonstrate the interface and should be refreshed when the UI changes.

## Function Layout

`main.go` initializes a `Game` value and passes it to `ebiten.RunGame`. The
`Game` methods fulfilling Ebiten's `Game` interface are split across files:

- `Update` in `update.go` processes input and camera movement.
- `Draw` in `draw.go` delegates to helpers in `drawing.go`, `draw_helpers.go` and
  `display.go` for rendering.
- `Layout` in `layout.go` keeps the camera centered when the window resizes.

Supporting files like `update_helpers.go` provide smaller helpers for menu
logic and screenshot processing. UI components live in `options_menu.go`,
`screenshot_menu.go` and `asteroid_menu.go` while the WASM screenshot
implementation is in `screenshot_save_wasm.go`.

## Build and Tests

Go files are formatted with `gofmt` and tests are run with `go test ./...`. As noted in `README.md`, unit tests are located alongside the code.

