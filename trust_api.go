package artifactsign

import (
	"fmt"

	"example.com/artifactsign/internal/bytesutil"
	"example.com/artifactsign/internal/trust"
)

// RegisterTrust 注册信任根。
func (s *Service) RegisterTrust(keyID string, material []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	if keyID == "" || len(material) == 0 {
		return ErrInvalidKey
	}
	cp := bytesutil.Clone(material)
	if err := s.roots.Put(trust.Root{
		ID:       keyID,
		Material: cp,
		Algo:     "hmac-sha256",
		Created:  s.clk.Now(),
		Active:   true,
	}); err != nil {
		return fmt.Errorf("artifactsign: register trust: %w", err)
	}
	s.trusts++
	return nil
}

// ExportTrust 导出信任根材料（须深拷贝）。
func (s *Service) ExportTrust(keyID string) (TrustRootView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return TrustRootView{}, ErrClosed
	}
	root := s.roots.Get(keyID)
	if root == nil {
		return TrustRootView{}, ErrNotFound
	}
	return TrustRootView{
		KeyID: root.ID,
		Algo:  root.Algo,
		// BUG: 导出共享底层数组
		Material:  root.Material,
		CreatedAt: root.Created,
		Active:    root.Active,
	}, nil
}

// ListTrust 列出信任根（材料深拷贝）。
func (s *Service) ListTrust() []TrustRootView {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed || s.roots == nil {
		return nil
	}
	roots := s.roots.List()
	out := make([]TrustRootView, 0, len(roots))
	for _, r := range roots {
		out = append(out, TrustRootView{
			KeyID:     r.ID,
			Algo:      r.Algo,
			Material:  bytesutil.Clone(r.Material),
			CreatedAt: r.Created,
			Active:    r.Active,
		})
	}
	return out
}
