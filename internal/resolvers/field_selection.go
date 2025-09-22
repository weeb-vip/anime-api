package resolvers

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	anime_repo "github.com/weeb-vip/anime-api/internal/db/repositories/anime"
)

// ExtractFieldSelection extracts requested fields from GraphQL context for the Anime type
func ExtractFieldSelection(ctx context.Context) *anime_repo.FieldSelection {
	return extractFieldsFromContext(ctx)
}

// ExtractAnimeFieldSelection specifically extracts field selection for nested Anime fields
func ExtractAnimeFieldSelection(ctx context.Context) *anime_repo.FieldSelection {
	return extractFieldsFromContext(ctx)
}

// extractFieldsFromContext parses the GraphQL context to determine which fields are requested
func extractFieldsFromContext(ctx context.Context) *anime_repo.FieldSelection {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return &anime_repo.FieldSelection{
			Fields: make(map[string]bool),
		}
	}

	fieldMap := make(map[string]bool)

	// Check if episodes field is requested
	for _, field := range graphql.CollectFieldsCtx(ctx, nil) {
		fieldMap[field.Name] = true
	}

	return &anime_repo.FieldSelection{
		Fields: fieldMap,
	}
}

// isEpisodesRequested checks if episodes field is requested in the selection
func isEpisodesRequested(selection *anime_repo.FieldSelection) bool {
	if selection == nil || selection.Fields == nil {
		return false
	}
	return selection.Fields["episodes"]
}