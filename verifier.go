package artifactsign

// Verifier 可插拔验签接口。
type Verifier interface {
	Verify(keyMaterial, payload, sig []byte) error
}
