## Running Headless

If you need to run the viewer on a machine without a display, install `Xvfb` and use the provided helper script. The script starts a virtual framebuffer so the window can be created.

```bash
sudo apt-get install xvfb   # one-time setup
./scripts/run_headless.sh -coord SNDST-A-7-0-0-0
```
