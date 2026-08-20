package artifactsign

import "time"

// Meta 签名时的业务元数据。
type Meta struct {
	Name        string
	Labels      []string
	Description string
	Algo        string
}

// SignatureView 对外可见的签名结果。
type SignatureView struct {
	ID          string
	KeyID       string
	DigestHex   string
	Signature   []byte
	PayloadCopy []byte
	Algo        string
	CreatedAt   time.Time
	Meta        Meta
}

// TrustRootView 信任根视图。
type TrustRootView struct {
	KeyID     string
	Algo      string
	Material  []byte
	CreatedAt time.Time
	Active    bool
}

// RevokeEntryView 吊销条目视图。
type RevokeEntryView struct {
	Fingerprint string
	Reason      string
	At          time.Time
}

// ManifestEntry 清单中的单文件摘要。
type ManifestEntry struct {
	Path      string
	DigestHex string
	Size      int64
}

// ManifestView 多文件清单视图。
type ManifestView struct {
	ID        string
	Entries   []ManifestEntry
	RootHex   string
	CreatedAt time.Time
}

// Stats 运行计数。
type Stats struct {
	Signs    uint64
	Verifies uint64
	Digests  uint64
	Revokes  uint64
	Trusts   uint64
}
