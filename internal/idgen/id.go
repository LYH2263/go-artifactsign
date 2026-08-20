package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

var seq uint64

// New 生成带前缀 ID。
func New(prefix string) string {
	n := atomic.AddUint64(&seq, 1)
	var b [4]byte
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("%s-%d-%d-%s", prefix, time.Now().UnixNano(), n, hex.EncodeToString(b[:]))
}
