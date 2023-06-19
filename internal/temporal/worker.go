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

const (
	FlagQueue = "queue"
)

type Worker struct {
	runner      *cli.Runner
	logger      logger.Logger
	temporal    client.Client
	registerers []Registerer
	queue       string
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
	cmd := &cobra.Command{
		Use:  "worker",
		Args: cobra.NoArgs,
		RunE: w.Run,
	}

	cmd.Flags().StringVarP(&w.queue, FlagQueue, "q", "", "queue name")
	cmd.MarkFlagRequired(FlagQueue)

	return cmd
}

func (w *Worker) Run(cmd *cobra.Command, args []string) error {
	return w.runner.Run(
		cmd.Context(),
		func(ctx context.Context) error {
			tw := worker.New(w.temporal, w.queue, worker.Options{})

			for _, r := range w.registerers {
				r.Register(tw)
			}

			return tw.Run(w.interruptCh(ctx))
		},
	)
}

func (w *Worker) interruptCh(ctx context.Context) <-chan interface{} {
	interruptCh := make(chan interface{}, 1)
	go func() {
		<-ctx.Done()

		interruptCh <- struct{}{}
	}()

	return interruptCh
}
