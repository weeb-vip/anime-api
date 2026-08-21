package anime_character_staff_link

import (
	"context"
	"time"
	"github.com/weeb-vip/anime-api/internal/db"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_character"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_staff"
	"github.com/weeb-vip/anime-api/metrics"
	metrics_lib "github.com/weeb-vip/go-metrics-lib"
)

type AnimeCharacterWithStaff struct {
	anime_character.AnimeCharacter
	VoiceActors []anime_staff.AnimeStaff `json:"voice_actors"`
}

func (AnimeCharacterWithStaff) TableName() string {
	table := AnimeCharacterStaffLink{}
	return table.TableName()
}

type AnimeCharacterStaffLinkRepositoryImpl interface {
	FindAnimeCharacterAndStaffByAnimeId(ctx context.Context, animeId string) ([]*AnimeCharacterWithStaff, error)
	FindCharactersByStaffId(ctx context.Context, staffId string) ([]*anime_character.AnimeCharacter, error)
}

type AnimeCharacterStaffLinkRepository struct {
	db *db.DB
}

func NewAnimeCharacterStaffLinkRepository(db *db.DB) AnimeCharacterStaffLinkRepositoryImpl {
	return &AnimeCharacterStaffLinkRepository{db: db}
}

func (a *AnimeCharacterStaffLinkRepository) FindAnimeCharacterAndStaffByAnimeId(ctx context.Context, animeId string) ([]*AnimeCharacterWithStaff, error) {
	startTime := time.Now()

	type joinResult struct {
		CharacterID            string
		CharacterName          string
		CharacterRole          string
		CharacterImage         string
		StaffID                string
		StaffURLSlug           *string
		GivenName              string
		FamilyName             string
		StaffImage             string
		Language               string
		Birthday               string
		BirthPlace             string
		BloodType              string
		Hobbies                string
		Summary                string
		CharacterWeight        string
		CharacterHeight        string
		CharacterBirthday      string
		CharacterZodiac        string
		CharacterGender        string
		CharacterMartialStatus string
		CharacterTitle         string
		CharacterSummary       string
	}

	var rows []joinResult

	err := a.db.DB.WithContext(ctx).
		Table("anime_character_staff_link").
		Select(`
			anime_character.id as character_id,
			anime_character.name as character_name,
			anime_character.role as character_role,
			anime_character.image as character_image,
			anime_character.weight as character_weight,
			anime_character.height as character_height,
			anime_character.birthday as character_birthday,
			anime_character.zodiac as character_zodiac,
			anime_character.gender as character_gender,
			anime_character.martial_status as character_martial_status,
			anime_character.title as character_title,
			anime_character.summary as character_summary,
			anime_staff.id as staff_id,
			anime_staff.url_slug as staff_url_slug,
			anime_staff.given_name,
			anime_staff.family_name,
			anime_staff.image as staff_image,
			anime_staff.language as language,
			anime_staff.birthday as birthday,
			anime_staff.birth_place as birth_place,
			anime_staff.blood_type as blood_type,
			anime_staff.hobbies as hobbies,
			anime_staff.summary as summary
		`).
		Joins("JOIN anime_character ON anime_character.id = anime_character_staff_link.character_id").
		Joins("JOIN anime_staff ON anime_staff.id = anime_character_staff_link.staff_id").
		Where("anime_character.anime_id = ?", animeId).
		Scan(&rows).Error

	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_character_staff_link",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	// Group by character ID
	characterMap := make(map[string]*AnimeCharacterWithStaff)
	for _, row := range rows {
		if _, exists := characterMap[row.CharacterID]; !exists {
			characterMap[row.CharacterID] = &AnimeCharacterWithStaff{
				AnimeCharacter: anime_character.AnimeCharacter{
					ID:            row.CharacterID,
					Name:          row.CharacterName,
					Role:          row.CharacterRole,
					Image:         row.CharacterImage,
					Weight:        row.CharacterWeight,
					Height:        row.CharacterHeight,
					Birthday:      row.CharacterBirthday,
					Zodiac:        row.CharacterZodiac,
					Gender:        row.CharacterGender,
					MartialStatus: row.CharacterMartialStatus,
					Title:         row.CharacterTitle,
					Summary:       row.CharacterSummary,
				},
				VoiceActors: []anime_staff.AnimeStaff{},
			}
		}

		characterMap[row.CharacterID].VoiceActors = append(characterMap[row.CharacterID].VoiceActors, anime_staff.AnimeStaff{
			ID:         row.StaffID,
			URLSlug:    row.StaffURLSlug,
			GivenName:  row.GivenName,
			FamilyName: row.FamilyName,
			Image:      row.StaffImage,
			Language:   row.Language,
			Birthday:   row.Birthday,
			BirthPlace: row.BirthPlace,
			BloodType:  row.BloodType,
			Hobbies:    row.Hobbies,
			Summary:    row.Summary,
		})
	}

	var result []*AnimeCharacterWithStaff
	for _, char := range characterMap {
		result = append(result, char)
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_character_staff_link",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return result, nil
}


// FindCharactersByStaffId returns every character this staff member is linked
// to, across every anime.
//
// anime_staff is deduplicated on (given_name, family_name) by the scraper and
// carries no anime_id, so one row is one real person and this join is the whole
// of their credited work -- no cross-anime name matching needed.
//
// Ordered by the anime's start date, newest first. MySQL sorts NULL first on
// ASC and last on DESC, so undated anime land at the end of the DESC ordering
// on their own; the explicit IS NULL key keeps that true if the direction ever
// changes.
func (a *AnimeCharacterStaffLinkRepository) FindCharactersByStaffId(ctx context.Context, staffId string) ([]*anime_character.AnimeCharacter, error) {
	startTime := time.Now()

	var characters []*anime_character.AnimeCharacter

	err := a.db.DB.WithContext(ctx).
		Table("anime_character_staff_link").
		Select("anime_character.*").
		Joins("JOIN anime_character ON anime_character.id = anime_character_staff_link.character_id").
		Joins("LEFT JOIN anime ON anime.id = anime_character.anime_id").
		Where("anime_character_staff_link.staff_id = ?", staffId).
		Order("anime.start_date IS NULL, anime.start_date DESC").
		Find(&characters).Error

	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_character_staff_link",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_character_staff_link",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return characters, nil
}
