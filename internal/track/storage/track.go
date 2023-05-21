package storage

import (
	"context"
	"musicd/internal/track/domain"

	"github.com/gofrs/uuid"
	"gopkg.in/guregu/null.v4"
)

type Tracks struct {
	queries *Queries
}

func NewTracks(q *Queries) *Tracks {
	return &Tracks{queries: q}
}

func (t *Tracks) FindTracks(ctx context.Context, trackIds []uuid.UUID) (map[uuid.UUID]*domain.Track, error) {
	tracks, err := t.queries.FindTracks(ctx, trackIds)
	if err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID]*domain.Track, len(tracks))
	for _, track := range tracks {
		result[track.ID] = &domain.Track{
			ID:         track.ID,
			Title:      track.Title,
			Duration:   int(track.Duration.Int64),
			Files:      track.Files,
			UploadedBy: track.UploadedBy,
			UploadedAt: track.UploadedAt,
		}
	}

	return result, nil
}

func (t *Tracks) Create(ctx context.Context, track domain.Track) error {
	return t.queries.CreateTrack(ctx, CreateTrackParams{
		ID:         track.ID,
		Title:      track.Title,
		Duration:   null.IntFrom(int64(track.Duration)),
		Files:      track.Files,
		UploadedBy: track.UploadedBy,
		UploadedAt: track.UploadedAt,
	})
}
