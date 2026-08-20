package sign

import "time"

// Record 签名上下文记录。
type Record struct {
	ID        string
	KeyID     string
	Digest    []byte
	Signature []byte
	Payload   []byte
	Algo      string
	Created   time.Time
	Name      string
}
