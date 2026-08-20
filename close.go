package artifactsign

// Close 先刷盘再释放密钥存储；之后 Sign/Verify 返回 ErrClosed。
func (s *Service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	// BUG: 先清空密钥再 Flush，写出空快照
	if s.ks != nil {
		s.ks.Clear()
	}
	err := s.flushLocked()
	return err
}
