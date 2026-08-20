package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

// HMACSHA256 计算 HMAC。
func HMACSHA256(key, payload []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("sign: empty key")
	}
	m := hmac.New(sha256.New, key)
	_, _ = m.Write(payload)
	return m.Sum(nil), nil
}

// HMACVerifier 内置验签。
type HMACVerifier struct{}

func (HMACVerifier) Verify(keyMaterial, payload, sig []byte) error {
	want, err := HMACSHA256(keyMaterial, payload)
	if err != nil {
		return err
	}
	if !hmac.Equal(want, sig) {
		return fmt.Errorf("sign: hmac mismatch")
	}
	return nil
}
