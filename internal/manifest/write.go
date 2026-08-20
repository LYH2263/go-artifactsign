package manifest

import (
	"encoding/json"
	"os"
)

// WriteAtomic 写临时文件、Sync、Close 后再 Rename，错误路径清理 .tmp。
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
	// 错误路径必须 Close 并清理 tmp，否则句柄泄漏/残留 .tmp
	if _, err := f.Write(b); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	// Windows 上未 Close 的句柄会占用 tmp，导致 Rename 失败并残留 .tmp
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, final)
}
