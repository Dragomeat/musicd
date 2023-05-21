package api

import (
	"context"
	"musicd/graph/model"
	"musicd/internal/graphql"
	"musicd/internal/track/domain"
	"musicd/internal/validator"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofrs/uuid"
	"gopkg.in/guregu/null.v4"
)

type TrackResolver struct {
	tracks             domain.Tracks
	trackLoaderFactory *TrackLoaderFactory
	s3                 *s3.Client
}

func NewTrackResolver(
	tracks domain.Tracks,
	trackLoaderFactory *TrackLoaderFactory,
	s3 *s3.Client,
) *TrackResolver {
	return &TrackResolver{
		tracks:             tracks,
		trackLoaderFactory: trackLoaderFactory,
		s3:                 s3,
	}
}

func (r *TrackResolver) GetTrack(ctx context.Context, id uuid.UUID) (*model.Track, error) {
	trackLoader := graphql.GetLoader[uuid.UUID, *domain.Track](ctx, r.trackLoaderFactory)

	track, err := trackLoader.Load(ctx, id)()
	if err != nil {
		return nil, err
	}

	presignClient := s3.NewPresignClient(r.s3)
	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("musicd"),
		Key:    aws.String(track.Files[0].Sha256),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 1 * time.Hour
	})
	if err != nil {
		return nil, err
	}

	return &model.Track{
		ID:                track.ID,
		Title:             track.Title,
		DurationInSeconds: track.Duration,
		URL:               req.URL,
	}, nil
}

func (r *TrackResolver) Tracks(ctx context.Context, first int, after *string, before *string) (*model.TrackList, error) {
	output := &model.TrackList{PageInfo: &model.PageInfo{}}
	input := struct {
		First  int
		After  *string
		Before *string
	}{First: first, After: after, Before: before}
	err := validation.ValidateStruct(
		&input,
		validation.Field(&input.First, validation.Required, validation.Min(1), validation.Max(100)),
		validation.Field(&input.After, validation.NilOrNotEmpty, is.Base64),
		validation.Field(&input.Before, validation.When(input.After == nil, validation.NilOrNotEmpty, is.Base64).Else(validation.Nil)),
	)
	if err != nil {
		return output, validator.NewErrInvalidInputFromOzzo(ctx, err)
	}

	tracks, err := r.tracks.PaginateTracks(
		ctx,
		first+1,
		null.StringFromPtr(after),
		null.StringFromPtr(before),
	)
	if err != nil {
		return output, err
	}

	edges := make([]*model.TrackEdge, 0, len(tracks))
	for i, track := range tracks {
		if i == input.First {
			break
		}

		edges = append(
			edges,
			&model.TrackEdge{
				Cursor: track.Cursor,
				Node: &model.Track{
					ID:                track.Node.ID,
					Title:             track.Node.Title,
					DurationInSeconds: track.Node.Duration,
				},
			},
		)
	}
	output.Edges = edges
	if len(edges) > 0 {
		output.PageInfo.StartCursor = edges[0].Cursor
		output.PageInfo.EndCursor = edges[len(edges)-1].Cursor
	}
	output.PageInfo.HasNextPage = len(tracks) > len(edges)
	output.PageInfo.HasPreviousPage = after != nil

	return output, nil
}
