package storage

import (
	"context"
	"musicd/internal/pagination"
	"musicd/internal/track/domain"
	"time"

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

type trackListCursor struct {
	Id         uuid.UUID `json:"id"`
	UploadedAt time.Time `json:"uploadedAt"`
}

func (t *Tracks) PaginateTracks(ctx context.Context, first int, after null.String, before null.String) ([]domain.TrackEdge, error) {
	params := PaginateTracksParams{
		First: int32(first),
	}

	if after.Valid {
		cursor, err := pagination.DecodeCursor[trackListCursor](ctx, after.String)
		if err != nil {
			return nil, err
		}

		params.AfterID = uuid.NullUUID{UUID: cursor.Id, Valid: true}
		params.AfterUploadedAt = null.TimeFrom(cursor.UploadedAt)
	}

	tracks, err := t.queries.PaginateTracks(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]domain.TrackEdge, len(tracks))
	for i, track := range tracks {
		cursor, err := pagination.EncodeCursor(
			trackListCursor{
				Id:         track.ID,
				UploadedAt: track.UploadedAt,
			},
		)
		if err != nil {
			return nil, err
		}
		result[i] = domain.TrackEdge{
			Cursor: cursor,
			Node: &domain.Track{
				ID:         track.ID,
				Title:      track.Title,
				Duration:   int(track.Duration.Int64),
				Files:      track.Files,
				UploadedBy: track.UploadedBy,
				UploadedAt: track.UploadedAt,
			},
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
