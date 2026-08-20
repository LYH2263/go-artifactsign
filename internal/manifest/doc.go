package manifest

// FileEntry 清单文件项。
type FileEntry struct {
	Path      string `json:"path"`
	DigestHex string `json:"digest"`
	Size      int64  `json:"size"`
}

// Document 清单文档。
type Document struct {
	ID      string      `json:"id"`
	RootHex string      `json:"root"`
	Entries []FileEntry `json:"entries"`
}
