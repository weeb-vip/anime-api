package work

import (
	"context"
	"errors"

	"gorm.io/gorm"

	workrepo "github.com/weeb-vip/anime-api/internal/db/repositories/work"
	"github.com/weeb-vip/anime-api/tracing"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type WorkServiceImpl interface {
	FindByID(ctx context.Context, id string) (*workrepo.Work, error)
	FindBySlug(ctx context.Context, slug string) (*workrepo.Work, error)
	FindByIDs(ctx context.Context, ids []string) ([]*workrepo.Work, error)
	CurrentlyPublishing(ctx context.Context, limit int) ([]*workrepo.Work, error)
}

type WorkService struct {
	Repository workrepo.WorkRepositoryImpl
}

func NewWorkService(repository workrepo.WorkRepositoryImpl) WorkServiceImpl {
	return &WorkService{Repository: repository}
}

// notFoundToNil turns "no such row" into a null result rather than an error.
//
// A slug that does not resolve is an ordinary 404 for the page above, not a
// failure worth an error in the response and a span marked bad. Every other
// database error still propagates.
func notFoundToNil(w *workrepo.Work, err error) (*workrepo.Work, error) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return w, nil
}

func (s *WorkService) FindByID(ctx context.Context, id string) (*workrepo.Work, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "FindByID")
	span.SetTag("service", "work")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return notFoundToNil(s.Repository.FindByID(spanCtx, id))
}

func (s *WorkService) FindBySlug(ctx context.Context, slug string) (*workrepo.Work, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "FindBySlug")
	span.SetTag("service", "work")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return notFoundToNil(s.Repository.FindBySlug(spanCtx, slug))
}

func (s *WorkService) FindByIDs(ctx context.Context, ids []string) ([]*workrepo.Work, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "FindByIDs")
	span.SetTag("service", "work")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return s.Repository.FindByIDs(spanCtx, ids)
}

func (s *WorkService) CurrentlyPublishing(ctx context.Context, limit int) ([]*workrepo.Work, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "CurrentlyPublishing")
	span.SetTag("service", "work")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return s.Repository.CurrentlyPublishing(spanCtx, limit)
}
