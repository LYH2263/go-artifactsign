package artifactsign

import "errors"

var (
	ErrClosed         = errors.New("artifactsign: service closed")
	ErrNotFound       = errors.New("artifactsign: not found")
	ErrInvalidKey     = errors.New("artifactsign: invalid key")
	ErrNoVerifier     = errors.New("artifactsign: verifier not configured")
	ErrVerifyFailed   = errors.New("artifactsign: verify failed")
	ErrRevoked        = errors.New("artifactsign: fingerprint revoked")
	ErrPersist        = errors.New("artifactsign: persist failed")
	ErrCanceled       = errors.New("artifactsign: canceled")
	ErrBadRequest     = errors.New("artifactsign: bad request")
	ErrDuplicate      = errors.New("artifactsign: duplicate")
	ErrDigestMismatch = errors.New("artifactsign: digest mismatch")
	ErrEmptyPayload   = errors.New("artifactsign: empty payload")
)
