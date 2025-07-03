//go:build screenshot

package main_test

import (
	main "oni-view"
	"os"
	"os/exec"
	"testing"
)

func TestHeadlessScreenshot(t *testing.T) {
	if _, err := exec.LookPath("xvfb-run"); err != nil {
		t.Skip("xvfb-run not installed")
	}
	if err := os.RemoveAll(main.ScreenshotFile); err != nil && !os.IsNotExist(err) {
		t.Fatalf("cleanup: %v", err)
	}
	cmd := exec.Command("bash", "scripts/run_headless.sh", "-screenshot", main.ScreenshotFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run failed: %v\n%s", err, string(out))
	}
	if fi, err := os.Stat(main.ScreenshotFile); err != nil || fi.Size() == 0 {
		t.Fatalf("screenshot not created")
	}
}
