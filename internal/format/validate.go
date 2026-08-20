package format

import (
	"fmt"
	"strings"
	"unicode"
)

// CheckKeyID 校验密钥 ID。
func CheckKeyID(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("format: empty key id")
	}
	if len(id) > 128 {
		return fmt.Errorf("format: key id too long")
	}
	for _, r := range id {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			continue
		}
		return fmt.Errorf("format: invalid key id rune")
	}
	return nil
}

// CheckAlgo 校验算法名。
func CheckAlgo(algo string) error {
	switch strings.ToLower(strings.TrimSpace(algo)) {
	case "hmac-sha256", "sha256", "":
		return nil
	default:
		return fmt.Errorf("format: unsupported algo %q", algo)
	}
}
