# Oni View

Oni View is a small utility for inspecting **Oxygen Not Included** seed data. It downloads seed information from [Maps Not Included](https://mapsnotincluded.org) and displays an interactive map of an asteroid using [Ebiten](https://ebiten.org/). The viewer renders geysers, points of interest and biome boundaries with scaling icons and a small legend. You can pan around and zoom in or out while the center of the map remains fixed.

![Screenshot](screenshot.png)

## Features

* Fetches seed data in [CBOR](https://en.wikipedia.org/wiki/CBOR) format from the MapsNotIncluded ingest API.
* Converts the data to JSON and optionally saves it to disk with the `-out` flag.
* Smooth panning via arrow keys, **WASD** or by dragging with the mouse.
* Touch controls for panning and pinch-to-zoom on supported devices.
* Mouse wheel and `+`/`-` keys control zoom. The window can be resized at any time.
* Legends for biome colors along with icons for geysers and points of interest.
* Supports running in headless environments using `Xvfb`.

## Getting Started

1. **Install Go** – Go 1.24 or newer is required to build the program.
2. **Download the image assets** – Icons used by the viewer are stored in the [`oni-seed-browser`](https://github.com/MapsNotIncluded/oni-seed-browser) repository. Run the helper script to clone the repo and copy the assets:

   ```bash
   ./scripts/download_assets.sh
   ```

   After the script completes the images will be placed in the `assets/` directory.

3. **Build and run** – Execute the program with a seed coordinate. The window defaults to 1280×720 but can be resized.

   ```bash
   go run . -coord SNDST-A-7-0-0-0
   ```

   To save the decoded seed as JSON use the optional `-out` flag:

   ```bash
   go run . -coord SNDST-A-7-0-0-0 -out seed.json
   ```

## Running Headless

If you need to run the viewer on a machine without a display, install `Xvfb` and use the provided script. The script starts a virtual framebuffer so the window can be created.

```bash
sudo apt-get install xvfb   # one-time setup
./scripts/run_headless.sh -coord SNDST-A-7-0-0-0
```

## Repository Layout

```
assets/     # Image files used by the viewer (downloaded separately)
main.go     # Program entry point and rendering logic
scripts/    # Helper scripts for asset download and headless execution
names.json  # Mapping tables for translating internal IDs to display names
```

Unit tests live alongside the code (`main_test.go`, `zoom_test.go`). Run them with:

```bash
go test ./...
```

## License

The project is released under the MIT License. See [LICENSE](LICENSE) for full details.

## Attribution

Seed information and image assets come from the [Maps Not Included](https://mapsnotincluded.org) project. Their source code repositories can be found on GitHub under [MapsNotIncluded](https://github.com/MapsNotIncluded).

