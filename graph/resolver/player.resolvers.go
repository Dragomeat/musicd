package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"context"
	"fmt"
	"musicd/graph/generated"
	"musicd/graph/model"

	"github.com/gofrs/uuid"
)

// Track is the resolver for the track field.
func (r *currentTrackResolver) Track(ctx context.Context, obj *model.CurrentTrack) (*model.Track, error) {
	return r.trackResolver.GetTrack(ctx, obj.Track.ID)
}

// CreatePlayer is the resolver for the createPlayer field.
func (r *mutationResolver) CreatePlayer(ctx context.Context) (*model.Player, error) {
	return r.playerResolver.CreatePlayer(ctx)
}

// StartPlayer is the resolver for the startPlayer field.
func (r *mutationResolver) StartPlayer(ctx context.Context, playerID uuid.UUID) (*model.Player, error) {
	return r.playerResolver.Start(ctx, playerID)
}

// StopPlayer is the resolver for the stopPlayer field.
func (r *mutationResolver) StopPlayer(ctx context.Context, playerID uuid.UUID) (*model.Player, error) {
	return r.playerResolver.Stop(ctx, playerID)
}

// SeekTo is the resolver for the seekTo field.
func (r *mutationResolver) SeekTo(ctx context.Context, playerID uuid.UUID, positionInSeconds int) (*model.Player, error) {
	return r.playerResolver.SeekTo(ctx, playerID, positionInSeconds)
}

// QueueTrack is the resolver for the queueTrack field.
func (r *mutationResolver) QueueTrack(ctx context.Context, playerID uuid.UUID, trackID uuid.UUID) (*model.Player, error) {
	return r.playerResolver.QueueTrack(ctx, playerID, trackID)
}

// RemoveTrackFromQueue is the resolver for the removeTrackFromQueue field.
func (r *mutationResolver) RemoveTrackFromQueue(ctx context.Context, playerID uuid.UUID, trackID uuid.UUID) (*model.Player, error) {
	return r.playerResolver.RemoveTrackFromQueue(ctx, playerID, trackID)
}

// MoveTrackInQueue is the resolver for the moveTrackInQueue field.
func (r *mutationResolver) MoveTrackInQueue(ctx context.Context, playerID uuid.UUID, trackID uuid.UUID, position int) (*model.Player, error) {
	panic(fmt.Errorf("not implemented: MoveTrackInQueue - moveTrackInQueue"))
}

// PreviousTrack is the resolver for the previousTrack field.
func (r *mutationResolver) PreviousTrack(ctx context.Context, playerID uuid.UUID) (*model.Player, error) {
	return r.playerResolver.PreviousTrack(ctx, playerID)
}

// NextTrack is the resolver for the nextTrack field.
func (r *mutationResolver) NextTrack(ctx context.Context, playerID uuid.UUID) (*model.Player, error) {
	return r.playerResolver.NextTrack(ctx, playerID)
}

// Player is the resolver for the player field.
func (r *queryResolver) Player(ctx context.Context, playerID uuid.UUID) (*model.Player, error) {
	return r.playerResolver.GetPlayer(ctx, playerID)
}

// Queue is the resolver for the queue field.
func (r *queryResolver) Queue(ctx context.Context, playerID uuid.UUID, first int, after *string, before *string) (*model.QueuedTrackList, error) {
	panic(fmt.Errorf("not implemented: Queue - queue"))
}

// Track is the resolver for the track field.
func (r *queuedTrackResolver) Track(ctx context.Context, obj *model.QueuedTrack) (*model.Track, error) {
	return r.trackResolver.GetTrack(ctx, obj.Track.ID)
}

// CurrentTrack returns generated.CurrentTrackResolver implementation.
func (r *Resolver) CurrentTrack() generated.CurrentTrackResolver { return &currentTrackResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// QueuedTrack returns generated.QueuedTrackResolver implementation.
func (r *Resolver) QueuedTrack() generated.QueuedTrackResolver { return &queuedTrackResolver{r} }

type currentTrackResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type queuedTrackResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) SeekTrack(ctx context.Context, playerID uuid.UUID, trackID uuid.UUID, positionInSeconds int) (*model.Player, error) {
	panic(fmt.Errorf("not implemented: SeekTrack - seekTrack"))
}
func (r *mutationResolver) UpdateCursor(ctx context.Context, playerID uuid.UUID, cursor int) (*model.Player, error) {
	panic(fmt.Errorf("not implemented: UpdateCursor - updateCursor"))
}
func (r *mutationResolver) SkipTrack(ctx context.Context, playerID uuid.UUID, trackID uuid.UUID) (*model.Player, error) {
	panic(fmt.Errorf("not implemented: SkipTrack - skipTrack"))
}
