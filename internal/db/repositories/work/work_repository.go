package work

import (
	"context"
	"strings"
	"time"

	"github.com/weeb-vip/anime-api/internal/db"
	"github.com/weeb-vip/anime-api/metrics"
)

type WorkRepositoryImpl interface {
	FindByID(ctx context.Context, id string) (*Work, error)
	FindBySlug(ctx context.Context, slug string) (*Work, error)
	FindByIDs(ctx context.Context, ids []string) ([]*Work, error)
	CurrentlyPublishing(ctx context.Context, limit int) ([]*Work, error)
	FindByType(ctx context.Context, workType string, offset int, limit int, sortBy string) ([]*Work, error)
	CountByType(ctx context.Context, workType string) (int64, error)
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

// orderForSort maps the API's sort names onto SQL.
//
// Every one of them ends in `id` so a page boundary is stable: two works with
// the same member count in an unspecified order can swap between requests, and
// a reader paging through would see one row twice and never see the other.
//
// NULLS LAST on the columns that have them. Postgres sorts nulls first on DESC
// by default, so `score desc` would open the shelf with the ten thousand works
// nobody has rated rather than the best ones.
func orderForSort(sortBy string) string {
	switch strings.ToUpper(strings.TrimSpace(sortBy)) {
	case "SCORE":
		return "score desc nulls last, id"
	case "NEWEST":
		return "published_from desc nulls last, id"
	case "TITLE":
		// COALESCE because title_en is null on the works the scraper has not
		// caught up with; sorting on it alone would group them all together
		// under whatever nulls-last does rather than by the name shown.
		return "lower(coalesce(nullif(title_en, ''), nullif(title_romaji, ''), title_jp)) asc nulls last, id"
	default:
		return "members desc nulls last, id"
	}
}

// FindByType lists one kind of work, ordered and paged.
//
// Backs the /manga and /light-novels browse pages. Served by
// idx_work_type_members from migration 000066 for the default ordering; the
// other three sorts fall back to a scan of the type's rows, which is the
// trade this makes deliberately -- indexing all four would be four indexes on
// a table that is written by a scraper running for days at a time.
func (r *WorkRepository) FindByType(ctx context.Context, workType string, offset int, limit int, sortBy string) ([]*Work, error) {
	startTime := time.Now()

	if workType == "" || limit <= 0 {
		return []*Work{}, nil
	}
	if offset < 0 {
		offset = 0
	}

	var found []*Work
	err := r.db.DB.WithContext(ctx).
		Where("type = ?", workType).
		Order(orderForSort(sortBy)).
		Offset(offset).
		Limit(limit).
		Find(&found).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Error)
		return nil, err
	}

	metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Success)
	return found, nil
}

// CountByType is the total behind the pager.
//
// A second query rather than a window function on the page above, because the
// count does not change between pages and the page query is the one that has
// to stay cheap.
func (r *WorkRepository) CountByType(ctx context.Context, workType string) (int64, error) {
	startTime := time.Now()

	if workType == "" {
		return 0, nil
	}

	var total int64
	err := r.db.DB.WithContext(ctx).
		Model(&Work{}).
		Where("type = ?", workType).
		Count(&total).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Error)
		return 0, err
	}

	metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Success)
	return total, nil
}
