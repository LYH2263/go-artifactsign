package artifactsign

import (
	"fmt"
	"os"
	"path/filepath"

	"example.com/artifactsign/internal/digest"
	"example.com/artifactsign/internal/idgen"
	"example.com/artifactsign/internal/manifest"
)

// BuildManifest 聚合多文件摘要。
func (s *Service) BuildManifest(files map[string][]byte) (ManifestView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ManifestView{}, ErrClosed
	}
	if len(files) == 0 {
		return ManifestView{}, ErrBadRequest
	}
	entries := make([]ManifestEntry, 0, len(files))
	var parts [][]byte
	for path, body := range files {
		dg := digest.SHA256(body)
		entries = append(entries, ManifestEntry{
			Path:      path,
			DigestHex: digest.Hex(dg),
			Size:      int64(len(body)),
		})
		parts = append(parts, dg)
	}
	root := digest.Aggregate(parts)
	id := idgen.New("mf")
	return ManifestView{
		ID:        id,
		Entries:   entries,
		RootHex:   digest.Hex(root),
		CreatedAt: s.clk.Now(),
	}, nil
}

// PersistManifest 将清单写入目录（原子临时文件）。
func (s *Service) PersistManifest(dir string, view ManifestView) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return "", ErrClosed
	}
	if dir == "" {
		return "", ErrBadRequest
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	final := filepath.Join(dir, view.ID+".json")
	doc := manifest.Document{
		ID:      view.ID,
		RootHex: view.RootHex,
		Entries: make([]manifest.FileEntry, 0, len(view.Entries)),
	}
	for _, e := range view.Entries {
		doc.Entries = append(doc.Entries, manifest.FileEntry{
			Path:      e.Path,
			DigestHex: e.DigestHex,
			Size:      e.Size,
		})
	}
	if err := manifest.WriteAtomic(final, doc); err != nil {
		return "", fmt.Errorf("artifactsign: manifest persist: %w", err)
	}
	return final, nil
}
