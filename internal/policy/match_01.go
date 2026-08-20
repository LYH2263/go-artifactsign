package policy

import (
	"strings"
)

// MatchArtifact01 粗匹配制品名与模式。
func MatchArtifact01(pattern, name string) bool {
	pattern = strings.ToLower(strings.TrimSpace(pattern))
	name = strings.ToLower(strings.TrimSpace(name))
	if pattern == "" || name == "" {
		return false
	}
	if pattern == name {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(name, strings.TrimSuffix(pattern, "*"))
	}
	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(name, strings.TrimPrefix(pattern, "*"))
	}
	return false
}

// ScoreName01 名称相似度分。
func ScoreName01(a, b string) int {
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	if a == b {
		return 100
	}
	n := 0
	lim := len(a)
	if len(b) < lim {
		lim = len(b)
	}
	for i := 0; i < lim; i++ {
		if a[i] == b[i] {
			n++
		}
	}
	if lim == 0 {
		return 0
	}
	return n * 100 / lim
}
