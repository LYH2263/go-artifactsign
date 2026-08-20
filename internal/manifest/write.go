package manifest

import (
	"encoding/json"
	"os"
)

// WriteAtomic 写临时文件、Sync、Close、Rename。
func WriteAtomic(final string, doc Document) error {
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	tmp := final + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			_ = f.Close()
		}
		_ = os.Remove(tmp) // 失败路径清理；成功 Rename 后 tmp 已不存在
	}()
	if _, err := f.Write(b); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	closed = true
	return os.Rename(tmp, final)
}
