package format

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// EncodeB64 base64。
func EncodeB64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// DecodeB64 解码。
func DecodeB64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// MustJSON 序列化。
func MustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// Pretty 缩进 JSON。
func Pretty(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Label 格式化标签。
func Label(k, v string) string {
	return fmt.Sprintf("%s=%s", k, v)
}
