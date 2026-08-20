package trust

import (
	"fmt"
	"sync"
	"time"

	"example.com/artifactsign/internal/bytesutil"
)

// Root 信任根。
type Root struct {
	ID       string
	Material []byte
	Algo     string
	Created  time.Time
	Active   bool
}

// Roots 信任根集合。
type Roots struct {
	mu   sync.Mutex
	byID map[string]*Root
}

// New 构造。
func New() *Roots {
	return &Roots{byID: make(map[string]*Root)}
}

// Put 写入（深拷贝材料）。
func (r *Roots) Put(root Root) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if root.ID == "" || len(root.Material) == 0 {
		return fmt.Errorf("trust: invalid root")
	}
	cp := root
	cp.Material = bytesutil.Clone(root.Material)
	r.byID[root.ID] = &cp
	return nil
}

// Get 返回拷贝。
func (r *Roots) Get(id string) *Root {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.byID[id]
	if !ok || e == nil {
		return nil
	}
	out := *e
	// BUG: Get 不拷贝 Material
	out.Material = e.Material
	return &out
}

// List 全部拷贝。
func (r *Roots) List() []Root {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Root, 0, len(r.byID))
	for _, e := range r.byID {
		cp := *e
		cp.Material = bytesutil.Clone(e.Material)
		out = append(out, cp)
	}
	return out
}
