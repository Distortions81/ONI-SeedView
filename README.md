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

