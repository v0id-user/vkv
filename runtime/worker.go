package runtime

import (
	"sync"
)

type Job func()

type WorkerPool struct {
	jobs    chan Job
	wg      sync.WaitGroup
	workers int
	stopCh  chan struct{}
}

func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		jobs:    make(chan Job, 128), // small, safe queue
		workers: workers,
		stopCh:  make(chan struct{}),
	}
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case job := <-p.jobs:
					if job != nil {
						job()
					}
				case <-p.stopCh:
					return
				}
			}
		}()
	}
}

func (p *WorkerPool) Submit(job Job) {
	p.jobs <- job
}

func (p *WorkerPool) Stop() {
	close(p.stopCh)
	p.wg.Wait()
}
