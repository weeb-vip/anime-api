package resolvers

import (
	"context"
	"encoding/json"

	"github.com/weeb-vip/anime-api/graph/model"
	workrepo "github.com/weeb-vip/anime-api/internal/db/repositories/work"
	"github.com/weeb-vip/anime-api/internal/services/anime"
	"github.com/weeb-vip/anime-api/internal/services/work"
	"github.com/weeb-vip/anime-api/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// jsonList decodes the JSON-in-a-text-column the scraper writes for authors and
// title synonyms.
//
// Returns nil rather than an error on anything unparseable. These columns are
// scraped, and one malformed row should cost that row its author list, not the
// whole query -- the same reasoning the scraper's own transformer uses.
func jsonList(raw *string) []string {
	if raw == nil || *raw == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(*raw), &out); err != nil {
		return nil
	}

	return out
}

func transformWorkToGraphQL(w workrepo.Work) *model.Work {
	return &model.Work{
		ID:            w.ID,
		MalID:         w.MalID,
		Type:          w.Type,
		URLSlug:       w.UrlSlug,
		TitleEn:       w.TitleEn,
		TitleJp:       w.TitleJp,
		TitleSynonyms: jsonList(w.TitleSynonyms),
		Synopsis:      w.Synopsis,
		ImageURL:      w.ImageURL,
		Status:        w.Status,
		Volumes:       w.Volumes,
		Chapters:      w.Chapters,
		PublishedFrom: w.PublishedFrom,
		PublishedTo:   w.PublishedTo,
		Demographic:   w.Demographic,
		Serialization: w.Serialization,
		Authors:       jsonList(w.Authors),
		Score:         w.Score,
		Ranking:       w.Ranking,
		Members:       w.Members,
		Favorites:     w.Favorites,
		CreatedAt:     w.CreatedAt.String(),
		UpdatedAt:     w.UpdatedAt.String(),
	}
}

// WorkBySlug backs /manga/<slug>.
//
// Null for an unknown slug rather than an error: that is a 404 for the page
// above, not a failure.
func WorkBySlug(ctx context.Context, workService work.WorkServiceImpl, slug string) (*model.Work, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "WorkBySlug",
		trace.WithAttributes(
			attribute.String("work.slug", slug),
			attribute.String("resolver.name", "WorkBySlug"),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	found, err := workService.FindBySlug(ctx, slug)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	if found == nil {
		span.SetAttributes(attribute.Bool("work.found", false))
		return nil, nil
	}

	return transformWorkToGraphQL(*found), nil
}

// WorkByID resolves a work from a federation reference.
//
// The router calls this when another subgraph returns a Work it does not own --
// list-service extending it with the viewer's reading progress, for instance.
// It is the same lookup workBySlug performs, keyed on id because that is what a
// federation key carries.
//
// A missing work is nil rather than an error, matching WorkBySlug: the id comes
// from another service's row, which can point at something this store has not
// received yet.
func WorkByID(ctx context.Context, workService work.WorkServiceImpl, id string) (*model.Work, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "WorkByID",
		trace.WithAttributes(
			attribute.String("work.id", id),
			attribute.String("resolver.name", "WorkByID"),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	found, err := workService.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	if found == nil {
		span.SetAttributes(attribute.Bool("work.found", false))
		return nil, nil
	}

	return transformWorkToGraphQL(*found), nil
}

// SourceWorkForAnime resolves Anime.sourceWork.
//
// Null is the ordinary answer. Two thirds of the catalogue is either original
// or adapted from something MyAnimeList's manga database does not cover, and a
// work that has not been scraped yet resolves to nothing until it is.
func SourceWorkForAnime(ctx context.Context, workService work.WorkServiceImpl, obj *model.Anime) (*model.Work, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "SourceWorkForAnime",
		trace.WithAttributes(
			attribute.String("anime.id", obj.ID),
			attribute.String("resolver.name", "SourceWorkForAnime"),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	if obj.SourceWorkID == nil || *obj.SourceWorkID == "" {
		span.SetAttributes(attribute.Bool("anime.has_source_work", false))
		return nil, nil
	}

	found, err := workService.FindByID(ctx, *obj.SourceWorkID)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	if found == nil {
		// The anime points at a work the read store has not received yet. Not
		// an error: CDC makes no ordering promise across tables.
		span.SetAttributes(attribute.Bool("work.found", false))
		return nil, nil
	}

	return transformWorkToGraphQL(*found), nil
}

// AdaptationsForWork resolves Work.adaptations: every anime adapted from it.
//
// Usually empty, and that is the expected answer rather than a gap. MyAnimeList
// holds far more manga than there are anime adapted from one, so most works
// have never been adapted.
func AdaptationsForWork(ctx context.Context, animeService anime.AnimeServiceImpl, obj *model.Work, limit *int) ([]*model.Anime, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "AdaptationsForWork",
		trace.WithAttributes(
			attribute.String("work.id", obj.ID),
			attribute.String("resolver.name", "AdaptationsForWork"),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	resultLimit := 25
	if limit != nil && *limit > 0 {
		resultLimit = *limit
	}

	// The same query the SHARED_SOURCE relation uses, with no anime to exclude
	// because here the work is the subject rather than one of its adaptations.
	found, err := animeService.RelatedAnimeBySourceWorkID(ctx, obj.ID, "", resultLimit)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	adaptations := make([]*model.Anime, 0, len(found))
	for _, entity := range found {
		if entity == nil {
			continue
		}
		transformed, err := transformAnimeToGraphQL(*entity)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		adaptations = append(adaptations, transformed)
	}

	span.SetAttributes(attribute.Int("work.adaptation_count", len(adaptations)))
	return adaptations, nil
}

// CurrentlyPublishingWorks backs the homepage's reading row.
//
// The counterpart of the currently-airing row above it: what someone can pick
// up and keep following, rather than a shelf of finished classics that never
// changes. An empty list is a valid answer -- it means the scraper has not
// reached any ongoing series yet, not that anything failed.
func CurrentlyPublishingWorks(ctx context.Context, workService work.WorkServiceImpl, limit *int) ([]*model.Work, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "CurrentlyPublishingWorks",
		trace.WithAttributes(
			attribute.String("resolver.name", "CurrentlyPublishingWorks"),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	// Capped as well as defaulted: this feeds a horizontal row, and an
	// unbounded limit would let one caller pull the whole table through the
	// router for a strip nobody scrolls that far along.
	resultLimit := 12
	if limit != nil && *limit > 0 {
		resultLimit = *limit
	}
	if resultLimit > 50 {
		resultLimit = 50
	}

	found, err := workService.CurrentlyPublishing(ctx, resultLimit)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	works := make([]*model.Work, 0, len(found))
	for _, entity := range found {
		if entity == nil {
			continue
		}
		works = append(works, transformWorkToGraphQL(*entity))
	}

	span.SetAttributes(attribute.Int("work.count", len(works)))
	return works, nil
}

// Works backs the /manga and /light-novels browse pages.
//
// Types are passed through as given rather than validated against a list.
// Work.type is a scraped label -- the schema says so -- and a kind nobody has
// heard of should come back as an empty page, which is true, rather than as an
// error the page has to special-case. Both lists are optional; omitting each
// means "every kind" and "exclude nothing" respectively.
func Works(ctx context.Context, workService work.WorkServiceImpl, input model.WorksInput) (*model.WorkPage, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "Works",
		trace.WithAttributes(
			attribute.String("resolver.name", "Works"),
			attribute.StringSlice("work.types", input.Types),
			attribute.StringSlice("work.exclude_types", input.ExcludeTypes),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	// Capped for the same reason CurrentlyPublishingWorks is: without it one
	// caller can pull all 53,000 manga through the router in a single request.
	perPage := 24
	if input.PerPage != nil && *input.PerPage > 0 {
		perPage = *input.PerPage
	}
	if perPage > 100 {
		perPage = 100
	}

	page := 0
	if input.Page != nil && *input.Page > 0 {
		page = *input.Page
	}

	sortBy := ""
	if input.SortBy != nil {
		sortBy = *input.SortBy
	}

	total, err := workService.CountByTypes(ctx, input.Types, input.ExcludeTypes)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	found, err := workService.FindByTypes(ctx, input.Types, input.ExcludeTypes, page*perPage, perPage, sortBy)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	works := make([]*model.Work, 0, len(found))
	for _, entity := range found {
		if entity == nil {
			continue
		}
		works = append(works, transformWorkToGraphQL(*entity))
	}

	span.SetAttributes(
		attribute.Int("work.count", len(works)),
		attribute.Int("work.total", int(total)),
	)

	return &model.WorkPage{
		Works:   works,
		Total:   int(total),
		Page:    page,
		PerPage: perPage,
	}, nil
}
