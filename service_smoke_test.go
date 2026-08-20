package artifactsign_test

import (
	"testing"

	"example.com/artifactsign"
)

func TestSmoke_SignVerify(t *testing.T) {
	s := artifactsign.New()
	defer s.Close()
	mustKey(t, s, "k1", []byte("secret-key-material!!"))
	payload := []byte("hello-artifact")
	view, err := s.Sign(payload, artifactsign.Meta{Name: "demo"})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Verify(view.KeyID, payload, view.Signature); err != nil {
		t.Fatal(err)
	}
}
