package e2e

import (
	"os"
	"testing"
)

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("file should exist: %s", path)
	}
}

func assertFileNotEmpty(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("cannot stat file: %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("file should not be empty: %s", path)
	}
}

func assertPDFHeader(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read file: %v", err)
	}
	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Fatalf("file should start with %%PDF-: %s", path)
	}
}
