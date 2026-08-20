package artifactsign

// Close 先刷盘再释放密钥存储；之后 Sign/Verify 返回 ErrClosed。
func (s *Service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	// 先刷盘（密钥仍在内存），再清空；否则写出空快照，重启后 LoadPersist 为空。
	err := s.flushLocked()
	if s.ks != nil {
		s.ks.Clear()
	}
	return err
}
