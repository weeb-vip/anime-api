package resolvers

import (
	"context"

	anime_repo "github.com/weeb-vip/anime-api/internal/db/repositories/anime"
)

// ExtractFieldSelection extracts requested fields from GraphQL context for the Anime type
func ExtractFieldSelection(ctx context.Context) *anime_repo.FieldSelection {
	// For now, return nil to use standard queries
	// This is where you would implement GraphQL field parsing
	return nil
}

// ExtractAnimeFieldSelection specifically extracts field selection for nested Anime fields
func ExtractAnimeFieldSelection(ctx context.Context) *anime_repo.FieldSelection {
	// For now, return nil to fall back to standard optimization
	// In a production implementation, you would parse the GraphQL AST here
	// to extract the specific fields requested for anime objects
	return nil
}