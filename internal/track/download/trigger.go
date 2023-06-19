package download

import (
	context "context"
	"crypto/sha1"
	"fmt"
	"musicd/internal/cli"

	"github.com/spf13/cobra"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

type Trigger struct {
	runner   *cli.Runner
	temporal client.Client
}

func NewTrigger(runner *cli.Runner, temporal client.Client) *Trigger {
	return &Trigger{
		runner:   runner,
		temporal: temporal,
	}
}

func (s *Trigger) Command() *cobra.Command {
	return &cobra.Command{
		Use:  "download [track id]",
		Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: s.Run,
	}
}

func (s *Trigger) Run(cmd *cobra.Command, args []string) error {
	return s.runner.Run(
		cmd.Context(),
		func(ctx context.Context) error {
			trackId := args[0]
			workflowId := fmt.Sprintf("track.download.%x", sha1.Sum([]byte(trackId)))
			workflowOptions := client.StartWorkflowOptions{
				ID:                    workflowId,
				WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY,
				TaskQueue:             "track",
			}

			_, err := s.temporal.ExecuteWorkflow(ctx, workflowOptions, WorkflowName, WorkflowInput{Url: trackId})
			return err
		},
	)
}
