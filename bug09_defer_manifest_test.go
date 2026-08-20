package artifactsign_test

import (
	"os"
	"path/filepath"
	"testing"

	"example.com/artifactsign"
)

func TestBug09_ManifestTempFileClosed(t *testing.T) {
	s := artifactsign.New()
	defer s.Close()
	view, err := s.BuildManifest(map[string][]byte{
		"a.bin": []byte("aaa"),
		"b.bin": []byte("bbb"),
	})
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path, err := s.PersistManifest(dir, view)
	if err != nil {
		t.Fatal(err)
	}
	// 二次写入同目标：若第一次未 Close，Windows 上 Rename/覆盖可能失败
	view.ID = view.ID + "-2"
	path2, err := s.PersistManifest(dir, view)
	if err != nil {
		t.Fatalf("second persist failed (likely temp handle leak): %v", err)
	}
	// 临时文件不应残留
	tmp := path + ".tmp"
	if _, err := os.Stat(tmp); err == nil {
		t.Fatalf("tmp still exists: %s", tmp)
	}
	if _, err := os.Stat(path2); err != nil {
		t.Fatal(err)
	}
	// 同目录再写第三份，并尝试删除第一份（句柄泄漏时 Windows 会失败）
	if err := os.Remove(path); err != nil {
		t.Fatalf("cannot remove persisted file (handle leak?): %v", err)
	}
	_ = filepath.Base(path2)
}
