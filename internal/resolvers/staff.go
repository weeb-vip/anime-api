package resolvers

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/weeb-vip/anime-api/graph/model"
	anime_character2 "github.com/weeb-vip/anime-api/internal/db/repositories/anime_character"
	anime_staff2 "github.com/weeb-vip/anime-api/internal/db/repositories/anime_staff"
	"github.com/weeb-vip/anime-api/internal/services/anime"
	anime_character_staff_link2 "github.com/weeb-vip/anime-api/internal/services/anime_character_staff_link"
	anime_staff3 "github.com/weeb-vip/anime-api/internal/services/anime_staff"
	"github.com/weeb-vip/anime-api/tracing"
)

func transformStaffToGraphQL(staff anime_staff2.AnimeStaff) *model.AnimeStaff {
	return &model.AnimeStaff{
		ID:         staff.ID,
		GivenName:  staff.GivenName,
		FamilyName: staff.FamilyName,
		Language:   &staff.Language,
		Image:      &staff.Image,
		Birthday:   &staff.Birthday,
		BirthPlace: &staff.BirthPlace,
		BloodType:  &staff.BloodType,
		Hobbies:    &staff.Hobbies,
		Summary:    &staff.Summary,
		CreatedAt:  &staff.CreatedAt,
		UpdatedAt:  &staff.UpdatedAt,
	}
}

func transformCharacterToGraphQL(character anime_character2.AnimeCharacter) *model.AnimeCharacter {
	return &model.AnimeCharacter{
		ID:            character.ID,
		AnimeID:       character.AnimeID,
		Name:          character.Name,
		Role:          character.Role,
		Birthday:      &character.Birthday,
		Zodiac:        &character.Zodiac,
		Gender:        &character.Gender,
		Race:          &character.Race,
		Height:        &character.Height,
		Weight:        &character.Weight,
		Title:         &character.Title,
		MartialStatus: &character.MartialStatus,
		Summary:       &character.Summary,
		Image:         &character.Image,
		CreatedAt:     &character.CreatedAt,
		UpdatedAt:     &character.UpdatedAt,
	}
}

// StaffByID resolves the staff query. A missing staff member is a null result,
// not an error -- the field is nullable so a stale link from an old page does
// not fail the whole query.
func StaffByID(ctx context.Context, staffService anime_staff3.AnimeStaffServiceImpl, id string) (*model.AnimeStaff, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "StaffByID",
		trace.WithAttributes(
			attribute.String("staff.id", id),
			attribute.String("resolver.name", "StaffByID"),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	staff, err := staffService.StaffByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetAttributes(attribute.Bool("staff.found", false))
			return nil, nil
		}
		span.RecordError(err)
		return nil, err
	}

	return transformStaffToGraphQL(*staff), nil
}

// RolesByStaffID resolves AnimeStaff.roles: every character this person played,
// each paired with the anime it was in.
//
// Two queries, not one per role: the characters come back in one join, then
// their anime are fetched in a single batch by id. Prolific voice actors have
// hundreds of credits, so a per-role anime lookup would be a few hundred round
// trips for one page.
func RolesByStaffID(
	ctx context.Context,
	linkService anime_character_staff_link2.AnimeCharacterStaffLinkImpl,
	animeService anime.AnimeServiceImpl,
	staffID string,
) ([]*model.StaffRole, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "RolesByStaffID",
		trace.WithAttributes(
			attribute.String("staff.id", staffID),
			attribute.String("resolver.name", "RolesByStaffID"),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	characters, err := linkService.FindCharactersByStaffId(ctx, staffID)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(attribute.Int("staff.role_count", len(characters)))
	if len(characters) == 0 {
		return []*model.StaffRole{}, nil
	}

	animeIDs := make([]string, 0, len(characters))
	seen := make(map[string]struct{}, len(characters))
	for _, character := range characters {
		if character.AnimeID == "" {
			continue
		}
		if _, ok := seen[character.AnimeID]; ok {
			continue
		}
		seen[character.AnimeID] = struct{}{}
		animeIDs = append(animeIDs, character.AnimeID)
	}

	animeByID := make(map[string]*model.Anime, len(animeIDs))
	if len(animeIDs) > 0 {
		foundAnime, err := animeService.AnimeByIDs(ctx, animeIDs)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		for _, entity := range foundAnime {
			if entity == nil {
				continue
			}
			transformed, err := transformAnimeToGraphQL(*entity)
			if err != nil {
				span.RecordError(err)
				return nil, err
			}
			animeByID[entity.ID] = transformed
		}
	}

	roles := make([]*model.StaffRole, 0, len(characters))
	for _, character := range characters {
		if character == nil {
			continue
		}
		roles = append(roles, &model.StaffRole{
			Character: transformCharacterToGraphQL(*character),
			Anime:     animeByID[character.AnimeID],
		})
	}

	return roles, nil
}
