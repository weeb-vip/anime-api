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
