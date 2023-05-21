package cli

import (
	"context"
	"sync"
)

type Runner struct {
	done chan struct{}
	wg   sync.WaitGroup
}

func NewRunner() *Runner {
	return &Runner{done: make(chan struct{})}
}

func (p *Runner) stop() {
	select {
	case p.done <- struct{}{}:
	default:
	}

	p.wg.Wait()
}

func (p *Runner) wait() {
	p.wg.Wait()
}

func (p *Runner) Run(ctx context.Context, fn func(ctx context.Context)) {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()
		<-p.done
	}()

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		fn(ctx)
	}()
}
