package bytesutil

// Clone 返回独立副本；nil 仍为 nil。
func Clone(src []byte) []byte {
	// BUG: 直接返回别名
	return src
}

// Equal 常量时间友好的长度+内容比较（演示用）。
func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var v byte
	for i := range a {
		v |= a[i] ^ b[i]
	}
	return v == 0
}
