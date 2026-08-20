package artifactsign_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"example.com/artifactsign"
)

func mustKey(t *testing.T, s *artifactsign.Service, id string, mat []byte) {
	t.Helper()
	if err := s.RegisterKey(id, mat); err != nil {
		t.Fatalf("RegisterKey: %v", err)
	}
}

func tempPersist(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func writePayloadFile(t *testing.T, dir string, name string, n int) string {
	t.Helper()
	p := filepath.Join(dir, name)
	b := bytes.Repeat([]byte("x"), n)
	if err := os.WriteFile(p, b, 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}
