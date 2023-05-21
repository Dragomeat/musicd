package temporal

import (
	"context"
	"musicd/internal/cli"
	"musicd/internal/logger"

	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
)

type Worker struct {
	runner      *cli.Runner
	logger      logger.Logger
	temporal    client.Client
	registerers []Registerer
}

type WorkersParams struct {
	fx.In

	Runner      *cli.Runner
	Logger      logger.Logger
	Temporal    client.Client
	Registerers []Registerer `group:"temporal.registerers"`
}

func NewWorker(params WorkersParams) *Worker {
	return &Worker{
		runner:      params.Runner,
		logger:      params.Logger,
		temporal:    params.Temporal,
		registerers: params.Registerers,
	}
}

func (w *Worker) Command() *cobra.Command {
	return &cobra.Command{
		Use:  "worker",
		Args: cobra.NoArgs,
		Run:  w.Run,
	}
}

func (w *Worker) Run(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	w.runner.Run(ctx, func(ctx context.Context) {
		tw := worker.New(w.temporal, "track", worker.Options{})

		for _, r := range w.registerers {
			r.Register(tw)
		}

		err := tw.Run(w.interruptCh(ctx))
		if err != nil {
			w.logger.Error(context.Background(), "temporal error", logger.Field("err", err))
		}
	})
}

func (w *Worker) interruptCh(ctx context.Context) <-chan interface{} {
	interruptCh := make(chan interface{}, 1)
	go func() {
		<-ctx.Done()

		interruptCh <- struct{}{}
	}()

	return interruptCh
}
