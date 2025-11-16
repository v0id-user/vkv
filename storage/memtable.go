package storage


import (
	"sync"
)

// bucket is one shard of the memtable.
// It has its own map and its own lock.
type bucket struct {
	mu sync.RWMutex
	m  map[string]string
}

// Memtable is an in memory key-value store
// sharded into 256 buckets for better concurrency.
type Memtable struct {
	buckets [256]*bucket
}

// NewMemtable initializes the sharded memtable.
func NewMemtable() *Memtable {
	mt := &Memtable{}
	for i := range mt.buckets {
		mt.buckets[i] = &bucket{
			m: make(map[string]string),
		}
	}
	return mt
}

// getBucket picks a bucket based on the first byte of the key.
func (m *Memtable) getBucket(key string) *bucket {
	if key == "" {
		// Empty keys should not really happen.
		// but just in case, send them to bucket 0.
		return m.buckets[0]
	}
	idx := uint8(key[0]) // first byte of key
	return m.buckets[idx]
}


// Set stores the value for the given key.
func (m *Memtable) Set(key, value string) {
	b := m.getBucket(key)

	b.mu.Lock()
	b.m[key] = value
	b.mu.Unlock()
}

// Get returns the value for a key and whether it existed.
func (m *Memtable) Get(key string) (string, bool) {
	b := m.getBucket(key)

	b.mu.RLock()
	val, ok := b.m[key]
	b.mu.RUnlock()

	return val, ok
}

// Del removes a key if it exists.
func (m *Memtable) Del(key string) {
	b := m.getBucket(key)

	b.mu.Lock()
	delete(b.m, key)
	b.mu.Unlock()
}


func (m *Memtable) Snapshot() map[string]string {
	out := make(map[string]string)
	for i := range m.buckets {
		b := m.buckets[i]
		b.mu.RLock()
		for k, v := range b.m {
			out[k] = v
		}
		b.mu.RUnlock()
	}
	return out
}
