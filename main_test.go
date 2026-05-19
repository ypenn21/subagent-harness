package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary builds the appteam binary to a temp directory and returns the path.
func buildBinary(t *testing.T) string {
	t.Helper()
	binPath := filepath.Join(t.TempDir(), "appteam")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
	}
	return binPath
}

func TestHelpFlag(t *testing.T) {
	bin := buildBinary(t)

	for _, flag := range []string{"--help", "-h"} {
		cmd := exec.Command(bin, flag)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("%s: unexpected error: %v", flag, err)
		}
		output := string(out)
		if !strings.Contains(output, "Usage: appteam") {
			t.Errorf("%s: output should contain 'Usage: appteam', got:\n%s", flag, output)
		}
		if !strings.Contains(output, "-d, --dir") {
			t.Errorf("%s: output should contain '-d, --dir', got:\n%s", flag, output)
		}
	}
}

func TestVersionFlag(t *testing.T) {
	bin := buildBinary(t)

	for _, flag := range []string{"--version", "-v"} {
		cmd := exec.Command(bin, flag)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("%s: unexpected error: %v", flag, err)
		}
		output := string(out)
		if !strings.Contains(output, "appteam v0.18.0") {
			t.Errorf("%s: output should contain 'appteam v0.18.0', got:\n%s", flag, output)
		}
	}
}

func TestUnknownFlag(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "--bogus")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("--bogus should exit non-zero")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() == 0 {
		t.Error("--bogus should exit with non-zero code")
	}
	output := string(out)
	if !strings.Contains(output, "Unknown option") {
		t.Errorf("--bogus output should contain 'Unknown option', got:\n%s", output)
	}
}

func TestRegenerateWithoutSettings(t *testing.T) {
	bin := buildBinary(t)
	tmpDir := t.TempDir()

	cmd := exec.Command(bin, "-r")
	cmd.Dir = tmpDir
	cmd.Env = append(os.Environ(), "HOME="+tmpDir)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("-r without settings.json should exit non-zero")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() == 0 {
		t.Error("-r should exit with non-zero code when no settings exist")
	}
	output := string(out)
	if !strings.Contains(output, "No settings.json found") {
		t.Errorf("-r output should contain 'No settings.json found', got:\n%s", output)
	}
}

func TestDirFlag(t *testing.T) {
	bin := buildBinary(t)
	tmpDir := filepath.Join(t.TempDir(), "new-project")

	// Run with -d pointing to a non-existent directory; provide empty stdin so wizard exits
	cmd := exec.Command(bin, "-d", tmpDir)
	cmd.Stdin = strings.NewReader("")
	cmd.Run() // ignore exit code — wizard will fail due to empty stdin, but dir should be created

	info, err := os.Stat(tmpDir)
	if err != nil {
		t.Fatalf("-d should create directory, got error: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("-d target should be a directory")
	}
}

func TestDirFlagLongForm(t *testing.T) {
	bin := buildBinary(t)
	tmpDir := filepath.Join(t.TempDir(), "long-form-dir")

	cmd := exec.Command(bin, "--dir", tmpDir)
	cmd.Stdin = strings.NewReader("")
	cmd.Run()

	info, err := os.Stat(tmpDir)
	if err != nil {
		t.Fatalf("--dir should create directory, got error: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("--dir target should be a directory")
	}
}

func TestDirFlagWithRegenerate(t *testing.T) {
	bin := buildBinary(t)
	tmpDir := filepath.Join(t.TempDir(), "regen-dir")

	// -r -d with no settings.json should fail
	cmd := exec.Command(bin, "-r", "-d", tmpDir)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("-r -d without settings.json should exit non-zero")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() == 0 {
		t.Error("-r -d should exit non-zero when no settings exist")
	}
	output := string(out)
	if !strings.Contains(output, "No settings.json found") {
		t.Errorf("output should contain 'No settings.json found', got:\n%s", output)
	}

	// But the directory should have been created
	info, statErr := os.Stat(tmpDir)
	if statErr != nil {
		t.Fatalf("-d should create directory even when -r fails, got error: %v", statErr)
	}
	if !info.IsDir() {
		t.Fatalf("-d target should be a directory")
	}
}

func TestDirFlagReversedOrder(t *testing.T) {
	bin := buildBinary(t)
	tmpDir := filepath.Join(t.TempDir(), "reversed-dir")

	// -d before -r should work the same
	cmd := exec.Command(bin, "-d", tmpDir, "-r")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("-d -r without settings.json should exit non-zero")
	}
	output := string(out)
	if !strings.Contains(output, "No settings.json found") {
		t.Errorf("output should contain 'No settings.json found', got:\n%s", output)
	}

	// Directory should still be created
	info, statErr := os.Stat(tmpDir)
	if statErr != nil {
		t.Fatalf("-d should create directory, got error: %v", statErr)
	}
	if !info.IsDir() {
		t.Fatalf("-d target should be a directory")
	}
}

func TestDirFlagMissingValue(t *testing.T) {
	bin := buildBinary(t)

	// -d with no following argument
	cmd := exec.Command(bin, "-d")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("-d without value should exit non-zero")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() == 0 {
		t.Error("-d without value should exit non-zero")
	}
	output := string(out)
	if !strings.Contains(output, "requires a directory argument") {
		t.Errorf("-d output should mention missing argument, got:\n%s", output)
	}
}

func TestDirFlagNextArgIsFlag(t *testing.T) {
	bin := buildBinary(t)

	// -d followed by another flag (no value)
	cmd := exec.Command(bin, "-d", "-r")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("-d -r (no dir value) should exit non-zero")
	}
	output := string(out)
	if !strings.Contains(output, "requires a directory argument") {
		t.Errorf("output should mention missing argument, got:\n%s", output)
	}
}
