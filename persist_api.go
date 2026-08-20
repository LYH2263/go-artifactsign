package artifactsign

import (
	"fmt"

	"example.com/artifactsign/internal/persist"
)

// Flush 将密钥元数据与吊销清单刷盘。
func (s *Service) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.flushLocked()
}

func (s *Service) flushLocked() error {
	if s.persistDir == "" {
		return nil
	}
	if s.ks == nil {
		return ErrClosed
	}
	meta := s.ks.MetaSnapshot()
	if err := persist.SaveKeyMeta(persist.KeyMetaPath(s.persistDir), meta); err != nil {
		return fmt.Errorf("%w: key meta: %v", ErrPersist, err)
	}
	if s.crl != nil {
		if err := persist.SaveRevokeList(persist.RevokePath(s.persistDir), s.crl.All()); err != nil {
			return fmt.Errorf("%w: revoke: %v", ErrPersist, err)
		}
	}
	return nil
}

// LoadPersist 从目录加载。
func (s *Service) LoadPersist() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	if s.persistDir == "" {
		return nil
	}
	meta, err := persist.LoadKeyMeta(persist.KeyMetaPath(s.persistDir))
	if err != nil && !persist.IsNotExist(err) {
		return err
	}
	if meta != nil && s.ks != nil {
		s.ks.LoadMeta(meta)
	}
	list, err := persist.LoadRevokeList(persist.RevokePath(s.persistDir))
	if err != nil && !persist.IsNotExist(err) {
		return err
	}
	if list != nil && s.crl != nil {
		s.crl.Replace(list)
	}
	return nil
}
