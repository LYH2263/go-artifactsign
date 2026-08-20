package artifactsign_test

import (
	"bytes"
	"testing"

	"example.com/artifactsign"
)

func TestBug02_ExportTrustMaterialAlias(t *testing.T) {
	s := artifactsign.New()
	defer s.Close()
	mat := []byte("trust-root-secret-key")
	if err := s.RegisterTrust("root1", mat); err != nil {
		t.Fatal(err)
	}
	view, err := s.ExportTrust("root1")
	if err != nil {
		t.Fatal(err)
	}
	for i := range view.Material {
		view.Material[i] = 0xff
	}
	again, err := s.ExportTrust("root1")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(again.Material, []byte{0xff}) && bytes.Equal(again.Material, view.Material) {
		t.Fatalf("export shared underlying array with stock")
	}
	if !bytes.Equal(again.Material, []byte("trust-root-secret-key")) {
		t.Fatalf("trust material corrupted: %q", again.Material)
	}
}
