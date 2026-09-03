package work

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/weeb-vip/anime-api/internal/db"
	"github.com/weeb-vip/anime-api/metrics"
)

type WorkRepositoryImpl interface {
	FindByID(ctx context.Context, id string) (*Work, error)
	FindBySlug(ctx context.Context, slug string) (*Work, error)
	FindByIDs(ctx context.Context, ids []string) ([]*Work, error)
	CurrentlyPublishing(ctx context.Context, limit int) ([]*Work, error)
	FindByTypes(ctx context.Context, types []string, excludeTypes []string, offset int, limit int, sortBy string) ([]*Work, error)
	CountByTypes(ctx context.Context, types []string, excludeTypes []string) (int64, error)
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
		//
		// title_en and title_jp only. The scraper's `work` also has a
		// title_romaji, and this ordering originally reached for it -- but
		// that column was never added to the read store, so the query failed
		// outright with `column "title_romaji" does not exist`. It is not
		// missed: since the manga heading fix, title_en carries the romanised
		// name whenever there is no English one, which is the case this
		// coalesce exists for.
		return "lower(coalesce(nullif(title_en, ''), title_jp)) asc nulls last, id"
	default:
		return "members desc nulls last, id"
	}
}

// scopeTypes applies the include and exclude lists.
//
// Shared by the page query and the count so the two can never disagree about
// what is on the shelf -- a total that counts rows the page cannot show is a
// pager that runs off the end.
func scopeTypes(q *gorm.DB, types []string, excludeTypes []string) *gorm.DB {
	if len(types) > 0 {
		q = q.Where("type IN ?", types)
	}
	if len(excludeTypes) > 0 {
		q = q.Where("type NOT IN ?", excludeTypes)
	}

	return q
}

// FindByTypes lists works of the given kinds, ordered and paged.
//
// Backs the /manga and /light-novels browse pages. The include case is served
// by idx_work_type_members from migration 000066 for the default ordering. The
// exclude case -- "everything that is not a novel" -- matches most of the
// table, so it reads the same index and filters rather than seeking within it;
// that is the honest cost of a shelf defined by what it is not, and it is
// bounded by the page size.
func (r *WorkRepository) FindByTypes(ctx context.Context, types []string, excludeTypes []string, offset int, limit int, sortBy string) ([]*Work, error) {
	startTime := time.Now()

	if limit <= 0 {
		return []*Work{}, nil
	}
	if offset < 0 {
		offset = 0
	}

	var found []*Work
	err := scopeTypes(r.db.DB.WithContext(ctx), types, excludeTypes).
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

// CountByTypes is the total behind the pager.
//
// A second query rather than a window function on the page above, because the
// count does not change between pages and the page query is the one that has
// to stay cheap.
func (r *WorkRepository) CountByTypes(ctx context.Context, types []string, excludeTypes []string) (int64, error) {
	startTime := time.Now()

	var total int64
	err := scopeTypes(r.db.DB.WithContext(ctx).Model(&Work{}), types, excludeTypes).
		Count(&total).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Error)
		return 0, err
	}

	metrics.GetAppMetrics().DatabaseSince(startTime, "work", metrics.MethodSelect, metrics.Success)
	return total, nil
}
