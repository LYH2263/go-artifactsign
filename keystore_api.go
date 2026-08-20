package artifactsign

import (
	"fmt"

	"example.com/artifactsign/internal/bytesutil"
	"example.com/artifactsign/internal/keystore"
)

// RegisterKey 注册 HMAC 演示密钥，返回 keyID。
func (s *Service) RegisterKey(keyID string, material []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	if s.ks == nil {
		return ErrClosed
	}
	if keyID == "" || len(material) == 0 {
		return ErrInvalidKey
	}
	cp := bytesutil.Clone(material)
	if err := s.ks.Put(keystore.Entry{
		ID:       keyID,
		Material: cp,
		Algo:     "hmac-sha256",
		Created:  s.clk.Now(),
		Active:   true,
	}); err != nil {
		return fmt.Errorf("artifactsign: register key: %w", err)
	}
	return nil
}

// ActivateKey 激活指定密钥（轮换）。
func (s *Service) ActivateKey(keyID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed || s.ks == nil {
		return ErrClosed
	}
	if err := s.ks.SetActive(keyID); err != nil {
		return fmt.Errorf("%w: %v", ErrNotFound, err)
	}
	return nil
}

// ActiveKeyID 当前活动密钥。
func (s *Service) ActiveKeyID() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed || s.ks == nil {
		return "", ErrClosed
	}
	id := s.ks.ActiveID()
	if id == "" {
		return "", ErrNotFound
	}
	return id, nil
}
