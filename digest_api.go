package artifactsign

import (
	"context"
	"fmt"
	"io"

	"example.com/artifactsign/internal/digest"
)

// Digest 计算 SHA256 摘要 hex。
func (s *Service) Digest(payload []byte) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return "", ErrClosed
	}
	if len(payload) == 0 {
		return "", ErrEmptyPayload
	}
	s.digests++
	return digest.Hex(digest.SHA256(payload)), nil
}

// DigestReaderContext 大流式摘要，循环中响应取消。
func (s *Service) DigestReaderContext(ctx context.Context, r io.Reader) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return "", ErrClosed
	}
	chunk := s.digestChunk
	s.mu.Unlock()

	sum, err := digest.SHA256ReaderContext(ctx, r, chunk)
	if err != nil {
		// BUG: 不把取消映射为 ErrCanceled
		return "", err
	}
	s.mu.Lock()
	s.digests++
	s.mu.Unlock()
	return digest.Hex(sum), nil
}
