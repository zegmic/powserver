package db

import "sync"

// in real system there should be used a persistent external cache (e.g. Redis) as
// 1. if service is restarted the token could be used again (for a short period its timestamp is valid)
// 2. there should be a cleanup mechanism for unused or unchallenged tokens based on its validity period
type Mem struct {
	data map[string]interface{}
	mut  sync.Mutex
}

func New() *Mem {
	return &Mem{
		data: make(map[string]interface{}),
	}
}

func (m *Mem) Add(key string) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.data[key] = struct{}{}
}

func (m *Mem) Exists(key string) bool {
	_, ok := m.data[key]
	return ok
}

func (m *Mem) Delete(key string) {
	m.mut.Lock()
	defer m.mut.Unlock()
	delete(m.data, key)
}
