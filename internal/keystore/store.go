package keystore

import (
	"fmt"
	"sync"
	"time"

	"example.com/artifactsign/internal/bytesutil"
	"example.com/artifactsign/internal/persist"
)

// Entry 密钥条目。
type Entry struct {
	ID       string
	Material []byte
	Algo     string
	Created  time.Time
	Active   bool
}

// Store 密钥存储。
type Store struct {
	mu     sync.Mutex
	max    int
	byID   map[string]*Entry
	active string
}

// New 构造。
func New(max int) *Store {
	if max < 1 {
		max = 1
	}
	return &Store{max: max, byID: make(map[string]*Entry)}
}

// Put 写入。
func (s *Store) Put(e Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.byID) >= s.max {
		if _, ok := s.byID[e.ID]; !ok {
			return fmt.Errorf("keystore: full")
		}
	}
	cp := e
	cp.Material = bytesutil.Clone(e.Material)
	s.byID[e.ID] = &cp
	if e.Active || s.active == "" {
		s.active = e.ID
		for id, ent := range s.byID {
			ent.Active = id == s.active
		}
	}
	return nil
}

// Get 拷贝返回。
func (s *Store) Get(id string) *Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.byID[id]
	if !ok || e == nil {
		return nil
	}
	out := *e
	out.Material = bytesutil.Clone(e.Material)
	return &out
}

// Active 当前活动密钥拷贝。
func (s *Store) Active() *Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active == "" {
		return nil
	}
	e := s.byID[s.active]
	if e == nil {
		return nil
	}
	out := *e
	out.Material = bytesutil.Clone(e.Material)
	return &out
}

// ActiveID 活动 ID。
func (s *Store) ActiveID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.active
}

// SetActive 切换活动。
func (s *Store) SetActive(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byID[id]; !ok {
		return fmt.Errorf("keystore: missing %s", id)
	}
	s.active = id
	for k, e := range s.byID {
		e.Active = k == id
	}
	return nil
}

// Len 数量。
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.byID)
}

// Clear 清空。
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byID = make(map[string]*Entry)
	s.active = ""
}

// MetaSnapshot 元数据（不含明文材料滥用演示：只存长度与指纹）。
func (s *Store) MetaSnapshot() []persist.KeyMeta {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]persist.KeyMeta, 0, len(s.byID))
	for _, e := range s.byID {
		out = append(out, persist.KeyMeta{
			ID:        e.ID,
			Algo:      e.Algo,
			Created:   e.Created,
			Active:    e.Active,
			Length:    len(e.Material),
			Material:  bytesutil.Clone(e.Material), // 演示持久化；生产应加密
		})
	}
	return out
}

// LoadMeta 从元数据恢复。
func (s *Store) LoadMeta(meta []persist.KeyMeta) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byID = make(map[string]*Entry, len(meta))
	s.active = ""
	for _, m := range meta {
		e := &Entry{
			ID:       m.ID,
			Material: bytesutil.Clone(m.Material),
			Algo:     m.Algo,
			Created:  m.Created,
			Active:   m.Active,
		}
		s.byID[m.ID] = e
		if m.Active {
			s.active = m.ID
		}
	}
	if s.active == "" && len(meta) > 0 {
		s.active = meta[0].ID
		s.byID[s.active].Active = true
	}
}
