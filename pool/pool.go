package pool

import (
	"fmt"
	"sync"
)

type Pool struct {
	jobs chan Runner
	exit chan struct{}
	wg   sync.WaitGroup

	mu   sync.Mutex
	size int
}

// NewPool returns a goroutine pool of the target size provided. Size
// must be positive.
func New(size int) (pool *Pool, err error) {
	p := &Pool{
		jobs: make(chan Runner),
		exit: make(chan struct{}),
	}
	err = p.Resize(size)
	return p, err
}

// Resize adjusts the pool size to the target provided. Adjusting up
// happens immediately, as newly started goroutines will be free to
// take from the queued jobs. Adjusting down will occur as running
// goroutines complete the current jobs. Target size must be positive.
func (p *Pool) Resize(size int) (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if size < 1 {
		return fmt.Errorf("bad pool size: %v", size)
	}

	for p.size < size {
		p.size++
		p.wg.Add(1)
		go p.runner()
	}
	for p.size > size {
		p.size--
		p.exit <- struct{}{}
	}

	return
}

// Run sends the given job to the pool's job queue.
func (p *Pool) Run(job Runner) {
	p.jobs <- job
}

// Size returns the current size of the job pool.
func (p *Pool) Size() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.size
}

// Wait closes the pool and blocks until all current jobs are completed.
func (p *Pool) Wait() {
	close(p.jobs)
	p.wg.Wait()
}

func (p *Pool) runner() {
	defer p.wg.Done()

	for {
		select {
		case job, ok := <-p.jobs:
			if !ok {
				return
			}
			job.Run()
		case <-p.exit:
			return
		}
	}
}
