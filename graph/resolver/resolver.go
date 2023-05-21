package resolver

import (
	"musicd/graph/generated"
	playerApi "musicd/internal/player/api"
	trackApi "musicd/internal/track/api"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Resolver struct {
	playerResolver *playerApi.PlayerResolver
	trackResolver  *trackApi.TrackResolver
}

func NewResolver(
	playerResolver *playerApi.PlayerResolver,
	trackResolver *trackApi.TrackResolver,
	trackLoaderFactory *trackApi.TrackLoaderFactory,
	s3 *s3.Client,
) *Resolver {
	return &Resolver{
		playerResolver: playerResolver,
		trackResolver:  trackResolver,
	}
}

func (r *Resolver) GetDirectives() generated.DirectiveRoot {
	return generated.DirectiveRoot{}
}
