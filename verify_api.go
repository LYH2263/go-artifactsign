package artifactsign

import (
	"fmt"

	"example.com/artifactsign/internal/digest"
	"example.com/artifactsign/internal/sign"
)

// Verify 验签；检查吊销指纹。
func (s *Service) Verify(keyID string, payload, sig []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	verifier, err := s.activeVerifier()
	if err != nil {
		return err
	}
	if s.ks == nil {
		return ErrClosed
	}
	ent := s.ks.Get(keyID)
	if ent == nil {
		// 也可从信任根取
		root := s.roots.Get(keyID)
		if root == nil {
			return ErrNotFound
		}
		fp := digest.Fingerprint(payload)
		if s.crl.Has(fp) {
			return ErrRevoked
		}
		if err := verifier.Verify(root.Material, payload, sig); err != nil {
			return fmt.Errorf("%w: %w", ErrVerifyFailed, err)
		}
		s.verifies++
		return nil
	}
	fp := digest.Fingerprint(payload)
	if s.crl.Has(fp) {
		return ErrRevoked
	}
	if err := verifier.Verify(ent.Material, payload, sig); err != nil {
		return fmt.Errorf("%w: %w", ErrVerifyFailed, err)
	}
	s.verifies++
	return nil
}

// VerifyWithDigest 先核对摘要再验签。
func (s *Service) VerifyWithDigest(keyID string, payload, sig []byte, wantHex string) error {
	got := digest.Hex(digest.SHA256(payload))
	if got != wantHex {
		return fmt.Errorf("%w: got %s want %s", ErrDigestMismatch, got, wantHex)
	}
	return s.Verify(keyID, payload, sig)
}

// DefaultHMACVerifier 返回内置验签器。
func DefaultHMACVerifier() Verifier {
	return sign.HMACVerifier{}
}
