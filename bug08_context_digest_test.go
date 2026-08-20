package artifactsign_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"example.com/artifactsign"
)

type slowReader struct {
	n   int
	max int
}

func (s *slowReader) Read(p []byte) (int, error) {
	if s.n >= s.max {
		return 0, io.EOF
	}
	time.Sleep(5 * time.Millisecond)
	p[0] = 'a'
	s.n++
	return 1, nil
}

func TestBug08_DigestReaderHonorsCancel(t *testing.T) {
	s := artifactsign.New(artifactsign.WithDigestChunk(1))
	defer s.Close()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	start := time.Now()
	_, err := s.DigestReaderContext(ctx, &slowReader{max: 200})
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected cancel")
	}
	if !errors.Is(err, artifactsign.ErrCanceled) && !errors.Is(err, context.Canceled) {
		t.Fatalf("want canceled, got %v", err)
	}
	if elapsed > 500*time.Millisecond {
		t.Fatalf("digest loop ignored cancel, elapsed=%s", elapsed)
	}
	_ = bytes.NewReader(nil) // keep import used if trimmed
}
