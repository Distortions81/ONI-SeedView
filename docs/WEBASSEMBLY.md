## WebAssembly Build

Run the build script to compile the program and copy the runtime files. The WebAssembly module is optimized and compressed, with the runtime JavaScript and HTML files placed in `build/`.

```bash
./scripts/build_all.sh
```

Open `build/index.html` in a browser to enter a seed. Valid seeds redirect to `view.html` which loads `oni-view.wasm.br` and decompresses it with [brotli-dec-wasm](https://github.com/httptoolkit/brotli-wasm). You can also specify the seed in the viewer URL with `view.html?coord=<seed>` and an optional `asteroid=<id>`.

The page also supports `index.html?coord=<seed>` or `#<seed>` and will forward you automatically.
