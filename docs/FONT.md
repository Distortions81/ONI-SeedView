## Minimizing the Font

If you change the `.ttf`, run the helper script to subset it to the ASCII range and shrink the binary size:

```bash
python3 scripts/minimize_font.py data/NotoSansMono.ttf data/NotoSansMono.ttf
```

The script requires the `fonttools` package (`pip install fonttools`) and overwrites the font in place. Use `git checkout -- data/NotoSansMono.ttf` to restore the original file.
