package artifactsign

// Close 先刷盘再释放密钥存储；之后 Sign/Verify 返回 ErrClosed。
func (s *Service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	err := s.flushLocked()
	s.closed = true
	// 释放密钥存储引用；Sign/Verify 等已在 closed 态提前返回，不再解引用 ks
	s.ks = nil
	return err
}
