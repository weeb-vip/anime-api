-- Remove database indexes

-- Anime table indexes
DROP INDEX IF EXISTS idx_anime_id ON anime;
DROP INDEX IF EXISTS idx_anime_status ON anime;
DROP INDEX IF EXISTS idx_anime_rating ON anime;
DROP INDEX IF EXISTS idx_anime_ranking ON anime;
DROP INDEX IF EXISTS idx_anime_created_at ON anime;
DROP INDEX IF EXISTS idx_anime_type ON anime;
DROP INDEX IF EXISTS idx_anime_source ON anime;
DROP INDEX IF EXISTS idx_anime_start_date ON anime;
DROP INDEX IF EXISTS idx_anime_end_date ON anime;
DROP INDEX IF EXISTS idx_anime_anidbid ON anime;
DROP INDEX IF EXISTS idx_anime_thetvdbid ON anime;

-- Full-text search indexes for title fields
DROP INDEX IF EXISTS idx_anime_title_en ON anime;
DROP INDEX IF EXISTS idx_anime_title_jp ON anime;
DROP INDEX IF EXISTS idx_anime_title_romaji ON anime;
DROP INDEX IF EXISTS idx_anime_title_kanji ON anime;

-- Episodes table indexes
DROP INDEX IF EXISTS idx_episodes_id ON episodes;
DROP INDEX IF EXISTS idx_episodes_anime_id ON episodes;
DROP INDEX IF EXISTS idx_episodes_aired ON episodes;
DROP INDEX IF EXISTS idx_episodes_episode_number ON episodes;
DROP INDEX IF EXISTS idx_episodes_created_at ON episodes;

-- Composite index for episodes by anime and air date
DROP INDEX IF EXISTS idx_episodes_anime_aired ON episodes;

-- Anime characters table indexes
DROP INDEX IF EXISTS idx_anime_character_id ON anime_character;
DROP INDEX IF EXISTS idx_anime_character_anime_id ON anime_character;
DROP INDEX IF EXISTS idx_anime_character_name ON anime_character;
DROP INDEX IF EXISTS idx_anime_character_role ON anime_character;

-- Anime staff table indexes
DROP INDEX IF EXISTS idx_anime_staff_id ON anime_staff;
DROP INDEX IF EXISTS idx_anime_staff_given_name ON anime_staff;
DROP INDEX IF EXISTS idx_anime_staff_family_name ON anime_staff;
DROP INDEX IF EXISTS idx_anime_staff_language ON anime_staff;

-- Anime character staff link table indexes
DROP INDEX IF EXISTS idx_character_staff_character_id ON anime_character_staff_link;
DROP INDEX IF EXISTS idx_character_staff_staff_id ON anime_character_staff_link;

-- Anime seasons table indexes
DROP INDEX IF EXISTS idx_anime_seasons_id ON anime_seasons;
DROP INDEX IF EXISTS idx_anime_seasons_episode_count ON anime_seasons;
DROP INDEX IF EXISTS idx_anime_seasons_created_at ON anime_seasons;

-- Relations table indexes
DROP INDEX IF EXISTS idx_relations_anime_id ON relations;
DROP INDEX IF EXISTS idx_relations_related_anime_id ON relations;