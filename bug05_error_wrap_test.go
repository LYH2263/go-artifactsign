package artifactsign_test

import (
	"errors"
	"testing"

	"example.com/artifactsign"
)

func TestBug05_VerifyErrorIsSentinel(t *testing.T) {
	s := artifactsign.New()
	defer s.Close()
	mustKey(t, s, "k1", []byte("secret-key-material!!"))
	payload := []byte("good")
	view, err := s.Sign(payload, artifactsign.Meta{Name: "n"})
	if err != nil {
		t.Fatal(err)
	}
	bad := append([]byte{}, view.Signature...)
	bad[0] ^= 0xff
	err = s.Verify(view.KeyID, payload, bad)
	if err == nil {
		t.Fatal("expected verify error")
	}
	if !errors.Is(err, artifactsign.ErrVerifyFailed) {
		t.Fatalf("want errors.Is ErrVerifyFailed, got %v", err)
	}
}
