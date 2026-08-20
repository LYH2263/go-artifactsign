package revoke

import (
	"sync"
	"time"
)

// Entry 吊销条目。
type Entry struct {
	Fingerprint string
	Reason      string
	At          time.Time
}

// List 吊销清单。
type List struct {
	mu   sync.Mutex
	byFP map[string]Entry
}

// New 构造。
func New() *List {
	return &List{byFP: make(map[string]Entry)}
}

// Add 添加。
func (l *List) Add(e Entry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	// 无回滚钩子；调用方必须在 persist 成功后再 Add
	l.byFP[e.Fingerprint] = e
}

// Has 查询。
func (l *List) Has(fp string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, ok := l.byFP[fp]
	return ok
}

// All 快照。
func (l *List) All() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Entry, 0, len(l.byFP))
	for _, e := range l.byFP {
		out = append(out, e)
	}
	return out
}

// Replace 整体替换。
func (l *List) Replace(entries []Entry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.byFP = make(map[string]Entry, len(entries))
	for _, e := range entries {
		l.byFP[e.Fingerprint] = e
	}
}

// Len 数量。
func (l *List) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.byFP)
}
