package artifactsign_test

import (
	"errors"
	"testing"

	"example.com/artifactsign"
)

func TestBug03_SignAfterCloseNoPanic(t *testing.T) {
	s := artifactsign.New()
	mustKey(t, s, "k1", []byte("secret-key-material!!"))
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	var panicked any
	func() {
		defer func() { panicked = recover() }()
		_, err := s.Sign([]byte("x"), artifactsign.Meta{Name: "n"})
		if !errors.Is(err, artifactsign.ErrClosed) {
			t.Fatalf("want ErrClosed, got %v", err)
		}
	}()
	if panicked != nil {
		t.Fatalf("panic after Close: %v", panicked)
	}
}
