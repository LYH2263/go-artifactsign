package persist

import (
	"os"
	"path/filepath"
)

// KeyMetaPath 密钥元数据路径。
func KeyMetaPath(dir string) string {
	return filepath.Join(dir, "keys.json")
}

// RevokePath 吊销清单路径。
func RevokePath(dir string) string {
	return filepath.Join(dir, "revoke.json")
}

// IsNotExist 是否不存在。
func IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
