package artifactsign_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"example.com/artifactsign"
	"example.com/artifactsign/internal/digest"
)

func TestBug06_RevokePersistFailureNoEffect(t *testing.T) {
	dir := t.TempDir()
	// 让 revoke.json 的父路径不可写：放入只读目录中的文件作为路径冲突
	// 更稳：persistDir 指向一个「同名文件」，Mkdir/Create 失败
	block := filepath.Join(dir, "block")
	if err := os.WriteFile(block, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := artifactsign.New(artifactsign.WithPersistDir(block)) // 文件当目录 → 持久化失败
	defer s.Close()
	fp := digest.Fingerprint([]byte("artifact"))
	err := s.RevokeFingerprint(fp, "leak")
	if err == nil {
		t.Fatal("expected persist error")
	}
	if !errors.Is(err, artifactsign.ErrPersist) {
		t.Fatalf("want ErrPersist, got %v", err)
	}
	if s.IsRevoked(fp) {
		t.Fatal("revoke must not take effect when persist fails")
	}
}
