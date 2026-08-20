package artifactsign

import (
	"context"
	"fmt"
	"time"

	"example.com/artifactsign/internal/bytesutil"
	"example.com/artifactsign/internal/digest"
	"example.com/artifactsign/internal/idgen"
	"example.com/artifactsign/internal/sign"
)

// Sign 对 payload 签名（同步）。
func (s *Service) Sign(payload []byte, meta Meta) (SignatureView, error) {
	return s.SignContext(context.Background(), payload, meta)
}

// SignContext 支持取消的签名。
func (s *Service) SignContext(ctx context.Context, payload []byte, meta Meta) (SignatureView, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	// 入口先检查已取消 ctx，避免无谓加锁与等待。
	if err := ctx.Err(); err != nil {
		return SignatureView{}, ErrCanceled
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return SignatureView{}, ErrClosed
	}
	if s.ks == nil {
		return SignatureView{}, ErrClosed
	}
	if len(payload) == 0 {
		return SignatureView{}, ErrEmptyPayload
	}

	// 可取消的步进等待（模拟 I/O）
	step := s.signStep
	s.mu.Unlock()
	if err := waitStep(ctx, step); err != nil {
		s.mu.Lock()
		return SignatureView{}, err
	}
	s.mu.Lock()
	if s.closed || s.ks == nil {
		return SignatureView{}, ErrClosed
	}

	key := s.ks.Active()
	if key == nil {
		return SignatureView{}, ErrInvalidKey
	}

	payloadCopy := bytesutil.Clone(payload)
	dg := digest.SHA256(payloadCopy)
	sig, err := sign.HMACSHA256(key.Material, payloadCopy)
	if err != nil {
		return SignatureView{}, fmt.Errorf("artifactsign: sign: %w", err)
	}
	id := idgen.New("sig")
	rec := &sign.Record{
		ID:        id,
		KeyID:     key.ID,
		Digest:    dg,
		Signature: bytesutil.Clone(sig),
		Payload:   payloadCopy,
		Algo:      "hmac-sha256",
		Created:   s.clk.Now(),
		Name:      meta.Name,
	}
	s.lastSigs[id] = rec
	s.signs++

	return SignatureView{
		ID:          id,
		KeyID:       key.ID,
		DigestHex:   digest.Hex(dg),
		Signature:   bytesutil.Clone(sig),
		PayloadCopy: bytesutil.Clone(payloadCopy),
		Algo:        "hmac-sha256",
		CreatedAt:   rec.Created,
		Meta:        meta,
	}, nil
}

// LookupSignature 按 ID 取回签名上下文（含 payload 副本）。
func (s *Service) LookupSignature(id string) (SignatureView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return SignatureView{}, ErrClosed
	}
	rec, ok := s.lastSigs[id]
	if !ok || rec == nil {
		return SignatureView{}, ErrNotFound
	}
	return SignatureView{
		ID:          rec.ID,
		KeyID:       rec.KeyID,
		DigestHex:   digest.Hex(rec.Digest),
		Signature:   bytesutil.Clone(rec.Signature),
		PayloadCopy: bytesutil.Clone(rec.Payload),
		Algo:        rec.Algo,
		CreatedAt:   rec.Created,
		Meta:        Meta{Name: rec.Name},
	}, nil
}

func waitStep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	// 听 ctx：取消立即失败返回，不再盲目 Sleep。
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ErrCanceled
	case <-t.C:
		return nil
	}
}
