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
	// 刷盘完成后再清空密钥，避免写出空快照
	if s.ks != nil {
		s.ks.Clear()
	}
	return err
}
