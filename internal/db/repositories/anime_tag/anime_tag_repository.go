package anime_tag

import (
	"github.com/weeb-vip/anime-api/internal/db"
)

type AnimeTagRepositoryImpl interface {
	SetTagsForAnime(animeID string, tagIDs []int64) error
	GetTagIDsForAnime(animeID string) ([]int64, error)
	GetTagNamesForAnime(animeID string) ([]string, error)
	GetTagNamesForAnimeIDs(animeIDs []string) (map[string][]string, error)
	AddTagToAnime(animeID string, tagID int64) error
	RemoveTagFromAnime(animeID string, tagID int64) error
	DeleteAllTagsForAnime(animeID string) error
}

type AnimeTagRepository struct {
	db *db.DB
}

func NewAnimeTagRepository(db *db.DB) AnimeTagRepositoryImpl {
	return &AnimeTagRepository{db: db}
}

// SetTagsForAnime replaces all tags for an anime with the given tag IDs
func (r *AnimeTagRepository) SetTagsForAnime(animeID string, tagIDs []int64) error {
	// Delete existing tag associations
	err := r.db.DB.Where("anime_id = ?", animeID).Delete(&AnimeTag{}).Error
	if err != nil {
		return err
	}

	// Insert new tag associations
	if len(tagIDs) > 0 {
		animeTags := make([]AnimeTag, len(tagIDs))
		for i, tagID := range tagIDs {
			animeTags[i] = AnimeTag{
				AnimeID: animeID,
				TagID:   tagID,
			}
		}
		err = r.db.DB.Create(&animeTags).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// GetTagIDsForAnime returns all tag IDs associated with an anime
func (r *AnimeTagRepository) GetTagIDsForAnime(animeID string) ([]int64, error) {
	var animeTags []AnimeTag
	err := r.db.DB.Where("anime_id = ?", animeID).Find(&animeTags).Error
	if err != nil {
		return nil, err
	}

	tagIDs := make([]int64, len(animeTags))
	for i, at := range animeTags {
		tagIDs[i] = at.TagID
	}
	return tagIDs, nil
}

// AddTagToAnime adds a single tag to an anime
func (r *AnimeTagRepository) AddTagToAnime(animeID string, tagID int64) error {
	animeTag := AnimeTag{
		AnimeID: animeID,
		TagID:   tagID,
	}
	return r.db.DB.Create(&animeTag).Error
}

// RemoveTagFromAnime removes a single tag from an anime
func (r *AnimeTagRepository) RemoveTagFromAnime(animeID string, tagID int64) error {
	return r.db.DB.Where("anime_id = ? AND tag_id = ?", animeID, tagID).Delete(&AnimeTag{}).Error
}

// DeleteAllTagsForAnime removes all tags for an anime
func (r *AnimeTagRepository) DeleteAllTagsForAnime(animeID string) error {
	return r.db.DB.Where("anime_id = ?", animeID).Delete(&AnimeTag{}).Error
}

// GetTagNamesForAnime returns all tag names associated with an anime
func (r *AnimeTagRepository) GetTagNamesForAnime(animeID string) ([]string, error) {
	var tagNames []string
	err := r.db.DB.Table("anime_tags").
		Select("tags.name").
		Joins("JOIN tags ON tags.id = anime_tags.tag_id").
		Where("anime_tags.anime_id = ?", animeID).
		Pluck("name", &tagNames).Error
	if err != nil {
		return nil, err
	}
	return tagNames, nil
}

// AnimeTagName represents a row with anime_id and tag name
type AnimeTagName struct {
	AnimeID string `gorm:"column:anime_id"`
	Name    string `gorm:"column:name"`
}

// GetTagNamesForAnimeIDs returns a map of anime ID to tag names for multiple anime
func (r *AnimeTagRepository) GetTagNamesForAnimeIDs(animeIDs []string) (map[string][]string, error) {
	if len(animeIDs) == 0 {
		return make(map[string][]string), nil
	}

	var results []AnimeTagName
	err := r.db.DB.Table("anime_tags").
		Select("anime_tags.anime_id, tags.name").
		Joins("JOIN tags ON tags.id = anime_tags.tag_id").
		Where("anime_tags.anime_id IN ?", animeIDs).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Build map from results
	tagMap := make(map[string][]string)
	for _, r := range results {
		tagMap[r.AnimeID] = append(tagMap[r.AnimeID], r.Name)
	}

	return tagMap, nil
}
