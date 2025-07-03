//go:build screenshot

package main_test

import (
	"os"
	"os/exec"
	"testing"
)

func TestHeadlessScreenshot(t *testing.T) {
	if _, err := exec.LookPath("xvfb-run"); err != nil {
		t.Skip("xvfb-run not installed")
	}
	if err := os.RemoveAll("screenshot.png"); err != nil && !os.IsNotExist(err) {
		t.Fatalf("cleanup: %v", err)
	}
	cmd := exec.Command("bash", "scripts/run_headless.sh", "-screenshot", "screenshot.png")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run failed: %v\n%s", err, string(out))
	}
	if fi, err := os.Stat("screenshot.png"); err != nil || fi.Size() == 0 {
		t.Fatalf("screenshot not created")
	}
}
