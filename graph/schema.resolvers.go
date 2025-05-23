package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"

	"github.com/weeb-vip/anime-api/graph/generated"
	"github.com/weeb-vip/anime-api/graph/model"
	"github.com/weeb-vip/anime-api/internal/resolvers"
)

// Episodes is the resolver for the episodes field.
func (r *animeResolver) Episodes(ctx context.Context, obj *model.Anime) ([]*model.Episode, error) {
	animeID := obj.ID
	return resolvers.EpisodesByAnimeID(ctx, r.AnimeEpisodeService, animeID)
}

// NextEpisode is the resolver for the nextEpisode field.
func (r *animeResolver) NextEpisode(ctx context.Context, obj *model.Anime) (*model.Episode, error) {
	if obj.NextEpisode != nil {
		return obj.NextEpisode, nil
	}
	animeID := obj.ID
	return resolvers.NextEpisode(ctx, r.AnimeEpisodeService, animeID)
}

// AnimeAPI is the resolver for the animeApi field.
func (r *apiInfoResolver) AnimeAPI(ctx context.Context, obj *model.APIInfo) (*model.AnimeAPI, error) {
	return resolvers.AnimeAPI(r.Config)
}

// DbSearch is the resolver for the dbSearch field.
func (r *queryResolver) DbSearch(ctx context.Context, searchQuery model.AnimeSearchInput) ([]*model.Anime, error) {
	return resolvers.DBSearchAnime(ctx, r.AnimeService, searchQuery.Query, searchQuery.Page, searchQuery.PerPage)
}

// APIInfo is the resolver for the apiInfo field.
func (r *queryResolver) APIInfo(ctx context.Context) (*model.APIInfo, error) {
	return resolvers.APIInfo(r.Config)
}

// Anime is the resolver for the anime field.
func (r *queryResolver) Anime(ctx context.Context, id string) (*model.Anime, error) {
	return resolvers.AnimeByID(ctx, r.AnimeService, id)
}

// NewestAnime is the resolver for the newestAnime field.
func (r *queryResolver) NewestAnime(ctx context.Context, limit *int) ([]*model.Anime, error) {
	return resolvers.NewestAnime(ctx, r.AnimeService, limit)
}

// TopRatedAnime is the resolver for the topRatedAnime field.
func (r *queryResolver) TopRatedAnime(ctx context.Context, limit *int) ([]*model.Anime, error) {
	return resolvers.TopRatedAnime(ctx, r.AnimeService, limit)
}

// MostPopularAnime is the resolver for the mostPopularAnime field.
func (r *queryResolver) MostPopularAnime(ctx context.Context, limit *int) ([]*model.Anime, error) {
	return resolvers.MostPopularAnime(ctx, r.AnimeService, limit)
}

// Episode is the resolver for the episode field.
func (r *queryResolver) Episode(ctx context.Context, id string) (*model.Episode, error) {
	panic(fmt.Errorf("not implemented: Episode - episode"))
}

// EpisodesByAnimeID is the resolver for the episodesByAnimeId field.
func (r *queryResolver) EpisodesByAnimeID(ctx context.Context, animeID string) ([]*model.Episode, error) {
	return resolvers.EpisodesByAnimeID(ctx, r.AnimeEpisodeService, animeID)
}

// CurrentlyAiring is the resolver for the currentlyAiring field.
func (r *queryResolver) CurrentlyAiring(ctx context.Context, input *model.CurrentlyAiringInput) ([]*model.Anime, error) {
	return resolvers.CurrentlyAiring(ctx, r.AnimeService, input)
}

// Anime is the resolver for the anime field.
func (r *userAnimeResolver) Anime(ctx context.Context, obj *model.UserAnime) (*model.Anime, error) {
	animeID := obj.AnimeID
	return resolvers.AnimeByID(ctx, r.AnimeService, animeID)
}

// Anime returns generated.AnimeResolver implementation.
func (r *Resolver) Anime() generated.AnimeResolver { return &animeResolver{r} }

// ApiInfo returns generated.ApiInfoResolver implementation.
func (r *Resolver) ApiInfo() generated.ApiInfoResolver { return &apiInfoResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// UserAnime returns generated.UserAnimeResolver implementation.
func (r *Resolver) UserAnime() generated.UserAnimeResolver { return &userAnimeResolver{r} }

type animeResolver struct{ *Resolver }
type apiInfoResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userAnimeResolver struct{ *Resolver }
