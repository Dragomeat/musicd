package cli

import (
	"context"
	"sync"

	"go.uber.org/fx"
)

type Runner struct {
	shutdowner fx.Shutdowner
	errors     chan<- error `name:"errors.channel"`
	done       chan struct{}
	stopOnce   sync.Once
	waitOnce   sync.Once
	wg         sync.WaitGroup
}

type NewRunnerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Errors     chan<- error `name:"errors.channel"`
}

func NewRunner(params NewRunnerParams) *Runner {
	return &Runner{shutdowner: params.Shutdowner, errors: params.Errors, done: make(chan struct{})}
}

func (p *Runner) start(fn func() error) error {
	go func() {
		defer func() {
			// Shutdown app when all goroutines are done.
			p.waitOnce.Do(func() {
				p.wg.Wait()
				err := p.shutdowner.Shutdown()
				if err != nil {
					p.errors <- err
				}
			})
		}()

		err := fn()
		if err != nil {
			p.errors <- err
		}
	}()

	return nil
}

func (p *Runner) stop() error {
	p.stopOnce.Do(
		func() {
			close(p.done)
			p.wg.Wait()
		},
	)

	return nil
}

func (p *Runner) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()
		select {
		case <-ctx.Done():
		case <-p.done:
		}
	}()

	p.wg.Add(1)
	defer p.wg.Done()

	return fn(ctx)
}
