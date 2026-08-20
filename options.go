package artifactsign

import (
	"time"

	"example.com/artifactsign/internal/clock"
)

// Option 配置 Service。
type Option func(*Service)

// WithClock 注入时钟（测试用）。
func WithClock(c clock.Clock) Option {
	return func(s *Service) {
		if c != nil {
			s.clk = c
		}
	}
}

// WithPersistDir 设置持久化目录（密钥元数据 / 吊销清单）。
func WithPersistDir(dir string) Option {
	return func(s *Service) {
		s.persistDir = dir
	}
}

// WithVerifier 注入验签器；缺省为内置 HMAC 验签。
func WithVerifier(v Verifier) Option {
	return func(s *Service) {
		s.verifier = v
	}
}

// WithNilVerifier 显式不配置验签器（用于复现缺省 nil 路径）。
func WithNilVerifier() Option {
	return func(s *Service) {
		s.verifier = nil
		s.allowNilVerifier = true
	}
}

// WithSignStep 签名路径每步等待（便于 ctx 取消测试）。
func WithSignStep(d time.Duration) Option {
	return func(s *Service) {
		if d > 0 {
			s.signStep = d
		}
	}
}

// WithDigestChunk 大文件摘要块大小。
func WithDigestChunk(n int) Option {
	return func(s *Service) {
		if n >= 16 {
			s.digestChunk = n
		}
	}
}

// WithMaxKeys 密钥上限。
func WithMaxKeys(n int) Option {
	return func(s *Service) {
		if n >= 1 {
			s.maxKeys = n
		}
	}
}
