package engine

import (
	"fmt"
	"time"

	"github.com/v0id-user/vkv/protocol"
	"github.com/v0id-user/vkv/storage"
)

// Route dispatches a parsed Command to the correct handler on the Engine.
func Route(e *Engine, cmd protocol.Command) protocol.Response {
	switch c := cmd.(type) {

	case protocol.Set:
		// 1) Log to WAL for durability (if present)
		if e.wal != nil {
			if err := e.wal.AppendSet(c.Key, c.Value); err != nil {
				// Internal failure -> surface as protocol error
				return protocol.ResponseErr("internal write error")
			}
		}

		// 2) Apply to in-memory storage
		e.memtable.Set(c.Key, c.Value)
		return protocol.ResponseOK()

	case protocol.Get:
		// Read from memtable, then SSTables
		if v, ok := e.readFromStorage(c.Key); ok {
			return protocol.ResponseValue(v)
		}
		return protocol.ResponseNil()

	case protocol.Del:
		// 1) Log deletion
		if e.wal != nil {
			if err := e.wal.AppendDel(c.Key); err != nil {
				return protocol.ResponseErr("internal delete error")
			}
		}

		// 2) Apply deletion in memory
		e.memtable.Del(c.Key)
		return protocol.ResponseOK()
	case protocol.Flush:
		snapshot := e.memtable.Snapshot()
		// Import "time" package at the top, since time is used here.
		path := fmt.Sprintf("data/%d.sst", time.Now().UnixNano())
		storage.BuildSSTable(path, snapshot)
		sst, err := storage.OpenSSTable(path)
		if err != nil {
			return protocol.ResponseErr("failed to open SSTable")
		}
		e.AddSSTable(sst)
		if e.wal != nil {
			e.wal.Reset()
		}
		e.memtable = storage.NewMemtable()
		return protocol.ResponseOK()
	default:
		// Should not happen, but guard.
		return protocol.ResponseErr("unknown command")
	}
}
