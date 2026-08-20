package persist

import (
	"encoding/json"
	"os"

	"example.com/artifactsign/internal/revoke"
)

// SaveRevokeList 保存吊销清单。
func SaveRevokeList(path string, entries []revoke.Entry) error {
	return writeJSON(path, entries)
}

// LoadRevokeList 加载。
func LoadRevokeList(path string) ([]revoke.Entry, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []revoke.Entry
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FailStore 测试用：永远失败的吊销存储钩子（由测试通过目录权限模拟）。
type FailStore struct{}
