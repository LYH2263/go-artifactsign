package artifactsign_test

import (
	"errors"
	"testing"

	"example.com/artifactsign"
)

func TestBug04_VerifyWithoutVerifier(t *testing.T) {
	s := artifactsign.New(artifactsign.WithNilVerifier())
	defer s.Close()
	mustKey(t, s, "k1", []byte("secret-key-material!!"))
	var panicked any
	var err error
	func() {
		defer func() { panicked = recover() }()
		err = s.Verify("k1", []byte("p"), []byte("sig"))
	}()
	if panicked != nil {
		t.Fatalf("panic: %v", panicked)
	}
	if !errors.Is(err, artifactsign.ErrNoVerifier) {
		t.Fatalf("want ErrNoVerifier, got %v", err)
	}
}
