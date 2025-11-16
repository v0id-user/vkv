package runtime

import (
	"net"
)

type Reactor struct {
	workerPool *WorkerPool
	connCh     chan net.Conn
	stopCh     chan struct{}
}

func NewReactor(pool *WorkerPool) *Reactor {
	return &Reactor{
		workerPool: pool,
		connCh:     make(chan net.Conn, 128),
		stopCh:     make(chan struct{}),
	}
}

// Register receives a new connection to be handled by a worker.
func (r *Reactor) Register(conn net.Conn, handler func(net.Conn)) {
	select {
	case r.connCh <- conn:
	default:
		// if channel is full, reject (rare for simple KV store)
		conn.Close()
	}

	// Wrap the handler into a Job for the worker pool
	r.workerPool.Submit(func() {
		handler(conn)
	})
}

func (r *Reactor) Stop() {
	close(r.stopCh)
}
