package artifactsign_test

import (
	"bytes"
	"testing"

	"example.com/artifactsign"
)

func TestBug01_SignPayloadBufferAlias(t *testing.T) {
	s := artifactsign.New()
	defer s.Close()
	mustKey(t, s, "k1", []byte("secret-key-material!!"))
	payload := []byte("artifact-payload-v1")
	view, err := s.Sign(payload, artifactsign.Meta{Name: "a"})
	if err != nil {
		t.Fatal(err)
	}
	// 调用方改写缓冲
	for i := range payload {
		payload[i] = 'Z'
	}
	got, err := s.LookupSignature(view.ID)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(got.PayloadCopy, payload) {
		t.Fatalf("stored payload aliased with caller buffer")
	}
	if !bytes.Equal(got.PayloadCopy, []byte("artifact-payload-v1")) {
		t.Fatalf("stored payload corrupted: %q", got.PayloadCopy)
	}
}
