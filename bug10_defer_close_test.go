package artifactsign_test

import (
	"path/filepath"
	"testing"

	"example.com/artifactsign"
)

func TestBug10_CloseFlushesBeforeClearKeys(t *testing.T) {
	dir := t.TempDir()
	s := artifactsign.New(artifactsign.WithPersistDir(dir))
	mustKey(t, s, "k1", []byte("secret-key-material!!"))
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	// 新实例加载，应能恢复密钥
	s2 := artifactsign.New(artifactsign.WithPersistDir(dir))
	defer s2.Close()
	if err := s2.LoadPersist(); err != nil {
		t.Fatal(err)
	}
	if s2.KeyCount() != 1 {
		t.Fatalf("want 1 key after reload, got %d (close flushed empty?)", s2.KeyCount())
	}
	id, err := s2.ActiveKeyID()
	if err != nil || id != "k1" {
		t.Fatalf("active=%s err=%v", id, err)
	}
	_ = filepath.Join(dir, "keys.json")
}
