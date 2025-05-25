package resolvers

import (
	"context"
	"github.com/weeb-vip/anime-api/graph/model"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_character_staff_link"
	anime_character_staff_link2 "github.com/weeb-vip/anime-api/internal/services/anime_character_staff_link"
)

func convertAnimeCharacterToGraphql(animeCharacterEntity *anime_character_staff_link.AnimeCharacterWithStaff) (*model.CharacterWithStaff, error) {
	if animeCharacterEntity == nil {
		return nil, nil
	}

	character := &model.AnimeCharacter{
		ID:            animeCharacterEntity.AnimeCharacter.ID,
		AnimeID:       animeCharacterEntity.AnimeCharacter.AnimeID,
		Name:          animeCharacterEntity.AnimeCharacter.Name,
		Role:          animeCharacterEntity.AnimeCharacter.Role,
		Birthday:      &animeCharacterEntity.AnimeCharacter.Birthday,
		Zodiac:        &animeCharacterEntity.AnimeCharacter.Zodiac,
		Gender:        &animeCharacterEntity.AnimeCharacter.Gender,
		Race:          &animeCharacterEntity.AnimeCharacter.Race,
		Height:        &animeCharacterEntity.AnimeCharacter.Height,
		Weight:        &animeCharacterEntity.AnimeCharacter.Weight,
		Title:         &animeCharacterEntity.AnimeCharacter.Title,
		MartialStatus: &animeCharacterEntity.AnimeCharacter.MartialStatus,
		Summary:       &animeCharacterEntity.AnimeCharacter.Summary,
		Image:         &animeCharacterEntity.AnimeCharacter.Image,
		CreatedAt:     &animeCharacterEntity.AnimeCharacter.CreatedAt,
		UpdatedAt:     &animeCharacterEntity.AnimeCharacter.UpdatedAt,
	}
	staffs := make([]*model.AnimeStaff, 0, len(animeCharacterEntity.VoiceActors))

	for _, staff := range animeCharacterEntity.VoiceActors {
		voiceActor := &model.AnimeStaff{
			ID:         staff.ID,
			GivenName:  staff.GivenName,
			FamilyName: staff.FamilyName,
			Image:      &staff.Image,
			Birthday:   &staff.Birthday,
			BirthPlace: &staff.BirthPlace,
			BloodType:  &staff.BloodType,
			Hobbies:    &staff.Hobbies,
			Summary:    &staff.Summary,
			CreatedAt:  &staff.CreatedAt,
			UpdatedAt:  &staff.UpdatedAt,
		}
		staffs = append(staffs, voiceActor)
	}

	return &model.CharacterWithStaff{
		Staff:     staffs,
		Character: character,
	}, nil
}

func CharactersAndStaffByAnimeID(ctx context.Context, animeCharacterStaffLinkService anime_character_staff_link2.AnimeCharacterStaffLinkImpl, animeId string) ([]*model.CharacterWithStaff, error) {
	animeCharacters, err := animeCharacterStaffLinkService.FindAnimeCharacterAndStaffByAnimeId(ctx, animeId)
	if err != nil {
		return nil, err
	}

	var characters []*model.CharacterWithStaff
	for _, animeCharacter := range animeCharacters {
		character, err := convertAnimeCharacterToGraphql(animeCharacter)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}
