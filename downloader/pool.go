package downloader

import (
	"context"
	"errors"
	"log"

	"golang.org/x/sync/errgroup"
)

type Job func() error

type Pool struct {
	size    int
	bufsize int
	ch      chan Job
	eg      *errgroup.Group
	cancel  context.CancelFunc
}

func NewPool() *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		size:    1,
		bufsize: 100,
		eg:      &errgroup.Group{},
		cancel:  cancel,
	}
	p.ch = make(chan Job, p.bufsize)
	for i := 0; i < p.size; i++ {
		p.eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case job, ok := <-p.ch:
					if !ok {
						return errors.New("channel closed")
					}
					if err := job(); err != nil {
						log.Printf("job failed: %+v", err)
					}
				}
			}
		})
	}
	return p
}

// This can potentially block if buffer has all been used up
func (p *Pool) Do(job Job) {
	p.ch <- job
}

func (p *Pool) Close() error {
	log.Printf("closing pool")
	p.cancel()
	return p.eg.Wait()
}
