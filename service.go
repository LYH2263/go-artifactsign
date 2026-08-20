package artifactsign

import (
	"sync"
	"time"

	"example.com/artifactsign/internal/clock"
	"example.com/artifactsign/internal/keystore"
	"example.com/artifactsign/internal/revoke"
	"example.com/artifactsign/internal/sign"
	"example.com/artifactsign/internal/trust"
)

const (
	defaultMaxKeys     = 256
	defaultDigestChunk = 64 * 1024
	defaultSignStep    = time.Millisecond
)

// Service 制品签名服务门面。零值不可用，须 New。
type Service struct {
	mu sync.Mutex

	closed           bool
	clk              clock.Clock
	ks               *keystore.Store
	roots            *trust.Roots
	crl              *revoke.List
	verifier         Verifier
	allowNilVerifier bool

	persistDir  string
	maxKeys     int
	digestChunk int
	signStep    time.Duration

	signs    uint64
	verifies uint64
	digests  uint64
	revokes  uint64
	trusts   uint64

	// lastSigs 保存签名上下文（含 payload 副本），供调试/重放。
	lastSigs map[string]*sign.Record
}

// New 构造 Service。
func New(opts ...Option) *Service {
	s := &Service{
		clk:         clock.Real{},
		maxKeys:     defaultMaxKeys,
		digestChunk: defaultDigestChunk,
		signStep:    defaultSignStep,
		lastSigs:    make(map[string]*sign.Record),
	}
	for _, o := range opts {
		if o != nil {
			o(s)
		}
	}
	if s.clk == nil {
		s.clk = clock.Real{}
	}
	if s.maxKeys < 1 {
		s.maxKeys = 1
	}
	if s.digestChunk < 16 {
		s.digestChunk = 16
	}
	s.ks = keystore.New(s.maxKeys)
	s.roots = trust.New()
	s.crl = revoke.New()
	if !s.allowNilVerifier && s.verifier == nil {
		s.verifier = sign.HMACVerifier{}
	}
	return s
}

// activeVerifier 返回已配置的验签器；未配置时返回 ErrNoVerifier。
func (s *Service) activeVerifier() (Verifier, error) {
	if s.verifier == nil {
		return nil, ErrNoVerifier
	}
	return s.verifier, nil
}

// Stats 返回计数快照。
func (s *Service) Stats() Stats {
	s.mu.Lock()
	defer s.mu.Unlock()
	return Stats{
		Signs:    s.signs,
		Verifies: s.verifies,
		Digests:  s.digests,
		Revokes:  s.revokes,
		Trusts:   s.trusts,
	}
}

// KeyCount 当前密钥数。
func (s *Service) KeyCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ks == nil {
		return 0
	}
	return s.ks.Len()
}

// Closed 是否已关闭。
func (s *Service) Closed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}
