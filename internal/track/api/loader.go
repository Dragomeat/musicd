package api

import (
	"context"
	"musicd/internal/track/domain"

	"github.com/gofrs/uuid"
	"github.com/graph-gophers/dataloader/v7"
)

type TrackLoaderFactory struct {
	tracks domain.Tracks
}

func NewTrackLoaderFactory(tracks domain.Tracks) *TrackLoaderFactory {
	return &TrackLoaderFactory{tracks: tracks}
}

func (f *TrackLoaderFactory) Create() *dataloader.Loader[uuid.UUID, *domain.Track] {
	return dataloader.NewBatchedLoader(
		func(ctx context.Context, trackIds []uuid.UUID) []*dataloader.Result[*domain.Track] {
			tracks, err := f.tracks.FindTracks(ctx, trackIds)

			result := make([]*dataloader.Result[*domain.Track], len(trackIds))
			for i, trackId := range trackIds {
				if err != nil {
					result[i] = &dataloader.Result[*domain.Track]{
						Error: err,
					}
					continue
				}

				result[i] = &dataloader.Result[*domain.Track]{
					Data: tracks[trackId],
				}
			}
			return result
		},
	)
}
