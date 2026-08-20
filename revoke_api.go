package artifactsign

import (
	"fmt"

	"example.com/artifactsign/internal/persist"
	"example.com/artifactsign/internal/revoke"
)

// RevokeFingerprint 吊销指纹；持久化失败则不生效。
func (s *Service) RevokeFingerprint(fp, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	if fp == "" {
		return ErrBadRequest
	}
	entry := revoke.Entry{
		Fingerprint: fp,
		Reason:      reason,
		At:          s.clk.Now(),
	}
	// 先尝试持久化，成功后再写入内存
	if s.persistDir != "" {
		path := persist.RevokePath(s.persistDir)
		snapshot := append(s.crl.All(), entry)
		if err := persist.SaveRevokeList(path, snapshot); err != nil {
			return fmt.Errorf("%w: %v", ErrPersist, err)
		}
	}
	s.crl.Add(entry)
	s.revokes++
	return nil
}

// IsRevoked 查询吊销。
func (s *Service) IsRevoked(fp string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.crl == nil {
		return false
	}
	return s.crl.Has(fp)
}

// ListRevoked 列出吊销条目。
func (s *Service) ListRevoked() []RevokeEntryView {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.crl == nil {
		return nil
	}
	all := s.crl.All()
	out := make([]RevokeEntryView, 0, len(all))
	for _, e := range all {
		out = append(out, RevokeEntryView{
			Fingerprint: e.Fingerprint,
			Reason:      e.Reason,
			At:          e.At,
		})
	}
	return out
}
