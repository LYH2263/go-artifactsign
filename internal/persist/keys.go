package persist

import (
	"encoding/json"
	"os"
	"time"
)

// KeyMeta 持久化密钥元数据。
type KeyMeta struct {
	ID       string    `json:"id"`
	Algo     string    `json:"algo"`
	Created  time.Time `json:"created"`
	Active   bool      `json:"active"`
	Length   int       `json:"length"`
	Material []byte    `json:"material,omitempty"`
}

// SaveKeyMeta 保存。
func SaveKeyMeta(path string, meta []KeyMeta) error {
	return writeJSON(path, meta)
}

// LoadKeyMeta 加载。
func LoadKeyMeta(path string) ([]KeyMeta, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []KeyMeta
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			_ = f.Close()
		}
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
	return os.Rename(tmp, path)
}
