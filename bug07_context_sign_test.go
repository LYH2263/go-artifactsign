package artifactsign_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/artifactsign"
)

func TestBug07_SignContextHonorsCancel(t *testing.T) {
	s := artifactsign.New(artifactsign.WithSignStep(200 * time.Millisecond))
	defer s.Close()
	mustKey(t, s, "k1", []byte("secret-key-material!!"))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	_, err := s.SignContext(ctx, []byte("payload"), artifactsign.Meta{Name: "n"})
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected cancel error")
	}
	if !errors.Is(err, artifactsign.ErrCanceled) && !errors.Is(err, context.Canceled) {
		t.Fatalf("want canceled, got %v", err)
	}
	if elapsed > 100*time.Millisecond {
		t.Fatalf("SignContext ignored cancel, elapsed=%s", elapsed)
	}
}
