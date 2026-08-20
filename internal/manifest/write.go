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
	// BUG: Write/Sync 后未 Close 就 Rename
	if _, err := f.Write(b); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	return os.Rename(tmp, final)
}
