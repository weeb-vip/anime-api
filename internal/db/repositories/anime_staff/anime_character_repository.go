package anime_staff

import (
	"github.com/weeb-vip/anime-api/internal/db"
)

type AnimeStaffRepositoryImpl interface {
}

type AnimeStaffRepository struct {
	db *db.DB
}

func NewAnimeStaffRepository(db *db.DB) AnimeStaffRepositoryImpl {
	return &AnimeStaffRepository{db: db}
}
