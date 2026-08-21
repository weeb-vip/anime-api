package anime_staff

import (
	"context"

	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_staff"
)

type AnimeStaffServiceImpl interface {
	StaffByID(ctx context.Context, id string) (*anime_staff.AnimeStaff, error)
	StaffBySlug(ctx context.Context, slug string) (*anime_staff.AnimeStaff, error)
}

type AnimeStaffService struct {
	Repository anime_staff.AnimeStaffRepositoryImpl
}

func NewAnimeStaffService(repository anime_staff.AnimeStaffRepositoryImpl) AnimeStaffServiceImpl {
	return &AnimeStaffService{
		Repository: repository,
	}
}

func (a *AnimeStaffService) StaffByID(ctx context.Context, id string) (*anime_staff.AnimeStaff, error) {
	return a.Repository.FindStaffByID(ctx, id)
}

func (a *AnimeStaffService) StaffBySlug(ctx context.Context, slug string) (*anime_staff.AnimeStaff, error) {
	return a.Repository.FindStaffBySlug(ctx, slug)
}
