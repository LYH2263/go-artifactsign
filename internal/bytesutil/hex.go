package bytesutil

import "encoding/hex"

// ToHex 编码。
func ToHex(b []byte) string {
	return hex.EncodeToString(b)
}

// FromHex 解码。
func FromHex(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
