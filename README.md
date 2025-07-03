# Oni View

This viewer uses image assets from the [oni-seed-browser](https://github.com/MapsNotIncluded/oni-seed-browser) project.
To keep this repository lightweight, the icons are not stored here.
Run the script below to download them into the `assets/` directory:

```bash
./scripts/download_assets.sh
```

Once downloaded, you can build and run the viewer normally.

The viewer now opens at **1280x720** and supports window resizing. Use the
arrow or **WASD** keys or drag with the mouse to pan, and the mouse wheel or
`+`/`-` keys to zoom. Zooming keeps the center of the screen fixed and the
icons scale with the current zoom level.

## Running headless

If you need to run the viewer without a physical display, install `Xvfb` and
use the `run_headless.sh` script:

```bash
sudo apt-get install xvfb   # one-time setup
./scripts/run_headless.sh -coord SNDST-A-7-0-0-0
```

The script uses `xvfb-run` to start a virtual framebuffer so the window can be
created even on headless machines.

