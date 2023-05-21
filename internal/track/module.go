package track

import (
	"musicd/internal/cli"
	"musicd/internal/temporal"
	"musicd/internal/track/api"
	"musicd/internal/track/domain"
	"musicd/internal/track/download"
	"musicd/internal/track/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"track",
		temporal.Provide[*download.Workflow](download.NewWorkflow),
		fx.Provide(
			api.NewTrackLoaderFactory,
			api.NewTrackResolver,

			download.NewTrigger,
			download.NewYoutubeDownloader,
			func(downloader *download.YoutubeDownloader) download.TrackInfoFetcher {
				return downloader
			},

			storage.NewTracks,
			func(tracks *storage.Tracks) domain.Tracks {
				return tracks
			},
			func(pool *pgxpool.Pool) storage.DBTX {
				return pool
			},
			func(db storage.DBTX) *storage.Queries {
				return storage.New(db)
			},

			cli.ProvideCommand(
				func(trigger *download.Trigger) *cobra.Command {
					return trigger.Command()
				},
			),
		),
	)
}
