package work

import (
	"context"
	"time"

	"github.com/weeb-vip/anime-api/internal/db"
	"github.com/weeb-vip/anime-api/metrics"
)

type WorkRepositoryImpl interface {
	FindByID(ctx context.Context, id string) (*Work, error)
	FindBySlug(ctx context.Context, slug string) (*Work, error)
	FindByIDs(ctx context.Context, ids []string) ([]*Work, error)
	CurrentlyPublishing(ctx context.Context, limit int) ([]*Work, error)
}

type WorkRepository struct {
	db *db.DB
}

func NewWorkRepository(db *db.DB) WorkRepositoryImpl {
	return &WorkRepository{db: db}
}

func (r *WorkRepository) FindByID(ctx context.Context, id string) (*Work, error) {
	startTime := time.Now()

	if id == "" {
		return nil, nil
	}

	var found Work
	err := r.db.DB.WithContext(ctx).Where("id = ?", id).First(&found).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Error)
		return nil, err
	}

	metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Success)
	return &found, nil
}

// FindBySlug backs /manga/<slug>. Served by idx_work_url_slug from migration
// 000064.
//
// The empty-string guard matters for the same reason it does on the anime
// relations: url_slug is nullable while the scraper catches up, and ” = ” is
// true, so a blank slug would match every other blank one and return an
// arbitrary work.
func (r *WorkRepository) FindBySlug(ctx context.Context, slug string) (*Work, error) {
	startTime := time.Now()

	if slug == "" {
		return nil, nil
	}

	var found Work
	err := r.db.DB.WithContext(ctx).Where("url_slug = ?", slug).First(&found).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Error)
		return nil, err
	}

	metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Success)
	return &found, nil
}

// FindByIDs exists so a list of anime can resolve its source works in one
// query rather than one per row.
func (r *WorkRepository) FindByIDs(ctx context.Context, ids []string) ([]*Work, error) {
	startTime := time.Now()

	if len(ids) == 0 {
		return []*Work{}, nil
	}

	var found []*Work
	err := r.db.DB.WithContext(ctx).Where("id IN ?", ids).Find(&found).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Error)
		return nil, err
	}

	metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Success)
	return found, nil
}

// CurrentlyPublishing returns ongoing works, most widely read first.
//
// The status match is case-insensitive because it is a scraped label rather
// than an enum: MyAnimeList writes "Publishing" today, and a row this query
// silently stops matching would empty the homepage section with nothing to
// point at. "On Hiatus" is excluded on purpose -- a paused series is not
// something a reader can follow right now, which is the whole question this
// answers.
//
// Served by idx_work_status_members from migration 000065, which indexes the
// same LOWER(status) expression and both ordering terms, so the whole query is
// one index scan with nothing left to sort. Measured on staging, the
// unindexed version cost 73ms against 2,011 rows -- the rows are wide enough
// (synopsis is in there) that scanning them all to answer a question about
// three columns is not free, even at this size.
func (r *WorkRepository) CurrentlyPublishing(ctx context.Context, limit int) ([]*Work, error) {
	startTime := time.Now()

	if limit <= 0 {
		return []*Work{}, nil
	}

	var found []*Work
	err := r.db.DB.WithContext(ctx).
		Where("LOWER(status) = ?", "publishing").
		// Tie-break on id so a page reload does not reshuffle works that share
		// a member count.
		Order("members desc, id").
		Limit(limit).
		Find(&found).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Error)
		return nil, err
	}

	metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Success)
	return found, nil
}
