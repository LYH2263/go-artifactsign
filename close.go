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
	// BUG: 置 nil，Sign 若未检查会 panic
	s.ks = nil
	return err
}
