package g

import (
	"context"
	"golang.org/x/sync/errgroup"
	"sync"
)

type Job func() error

/*
	goroutine pool
    控制协程数量 顺序执行 job
*/
type Pool interface {
	Wait() error
	JJ(Job)
}

type flag int

const (
	running flag = 1
	closed  flag = -1
)

type pool struct {
	f    flag
	m    sync.Mutex
	once sync.Once
	err  error

	wc  int
	cap int
	jc  chan Job
	g   *errgroup.Group
}

func NewPool(workerCount, cap int) Pool {
	if workerCount < 1 || cap < 1 {
		panic("the workerCount or cap should gt 1")
	}
	if cap < workerCount {
		cap = workerCount
	}
	g, _ := errgroup.WithContext(context.Background())
	jc := make(chan Job, cap)
	return &pool{
		g:  g,
		wc: workerCount,
		jc: jc,
	}
}

// join job
func (p *pool) JJ(job Job) {
	p.run()
	p.jj(job)
}

func (p *pool) jj(job Job) {
	p.jc <- job
}

func (p *pool) run() {
	if p.f == running {
		return
	}
	if p.f == closed {
		panic(" can not run, the flag it`s closed ")
	}
	p.m.Lock()
	if p.f == running {
		p.m.Unlock()
		return
	}
	p.f = running
	p.runWorkers()
	p.m.Unlock()
}
func (p *pool) runWorker() {
	for job := range p.jc {
		err := job()
		if err != nil {
			p.once.Do(func() {
				p.err = err
			})
		}
	}
}

func (p *pool) runWorkers() {
	for i := 0; i < p.wc; i++ {
		p.g.Go(func() error {
			p.runWorker()
			return nil
		})
	}
}

func (p *pool) close() error {

	close(p.jc)
	_ = p.g.Wait()
	p.m.Lock()
	p.f = closed
	p.m.Unlock()
	return p.err
}

func (p *pool) Close() error {
	return p.close()
}

func (p *pool) Wait() error {
	return p.Close()
}
