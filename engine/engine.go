package engine

import (
	"sync"

	"github.com/v0id-user/vkv/protocol"
	"github.com/v0id-user/vkv/storage"
)

// Engine wires protocol-level commands to the storage backend (memtable, SSTables, WAL).
type Engine struct {
	memtable *storage.Memtable
	wal      *storage.WAL

	mu       sync.RWMutex
	sstables []*storage.SSTable
}

// New creates a new Engine instance.
func New(mem *storage.Memtable, wal *storage.WAL) *Engine {
	return &Engine{
		memtable: mem,
		wal:      wal,
		sstables: make([]*storage.SSTable, 0),
	}
}

// AddSSTable registers a new SSTable (e.g. after a flush).
// Newer tables should be added last so reads check them first.
func (e *Engine) AddSSTable(sst *storage.SSTable) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.sstables = append(e.sstables, sst)
}

// Execute is the main entrypoint: take a Command, produce a Response.
func (e *Engine) Execute(cmd protocol.Command) protocol.Response {
	return Route(e, cmd)
}

// readFromStorage tries memtable first, then SSTables (newest → oldest).
func (e *Engine) readFromStorage(key string) (string, bool) {
	// First: in-memory memtable
	if v, ok := e.memtable.Get(key); ok {
		return v, true
	}

	// Then: on-disk SSTables (newest first)
	e.mu.RLock()
	defer e.mu.RUnlock()

	for i := len(e.sstables) - 1; i >= 0; i-- {
		sst := e.sstables[i]
		if v, ok := sst.Get(key); ok {
			return v, true
		}
	}

	return "", false
}
