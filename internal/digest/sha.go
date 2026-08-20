package digest

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256 计算摘要。
func SHA256(payload []byte) []byte {
	sum := sha256.Sum256(payload)
	out := make([]byte, len(sum))
	copy(out, sum[:])
	return out
}

// Hex 转十六进制。
func Hex(sum []byte) string {
	return hex.EncodeToString(sum)
}

// Fingerprint 制品指纹（摘要 hex）。
func Fingerprint(payload []byte) string {
	return Hex(SHA256(payload))
}

// Aggregate 多摘要再哈希。
func Aggregate(parts [][]byte) []byte {
	h := sha256.New()
	for _, p := range parts {
		_, _ = h.Write(p)
	}
	return h.Sum(nil)
}
