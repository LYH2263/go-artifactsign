package digest

import (
	"context"
	"crypto/sha256"
	"io"
)

// SHA256ReaderContext 分块读并响应取消。
func SHA256ReaderContext(ctx context.Context, r io.Reader, chunk int) ([]byte, error) {
	if chunk < 16 {
		chunk = 16
	}
	h := sha256.New()
	buf := make([]byte, chunk)
	for {
		// 每读完一块就听 ctx，用户取消时立刻退出而不是把整段读完。
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		n, err := r.Read(buf)
		if n > 0 {
			_, _ = h.Write(buf[:n])
		}
		if err == io.EOF {
			return h.Sum(nil), nil
		}
		if err != nil {
			return nil, err
		}
	}
}
