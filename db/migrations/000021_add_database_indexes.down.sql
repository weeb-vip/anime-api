-- Remove database indexes.
--
-- The MySQL version needed a stored procedure because DROP INDEX had no IF
-- EXISTS. Postgres has it, and drops by index name alone -- index names are
-- schema-global here, not per-table.

DROP INDEX IF EXISTS idx_anime_status;
DROP INDEX IF EXISTS idx_anime_rating;
DROP INDEX IF EXISTS idx_anime_ranking;
DROP INDEX IF EXISTS idx_anime_created_at;
DROP INDEX IF EXISTS idx_anime_type;
DROP INDEX IF EXISTS idx_anime_source;
DROP INDEX IF EXISTS idx_anime_start_date;
DROP INDEX IF EXISTS idx_anime_end_date;
DROP INDEX IF EXISTS idx_anime_anidbid;
DROP INDEX IF EXISTS idx_anime_thetvdbid;
DROP INDEX IF EXISTS idx_anime_title_en;
DROP INDEX IF EXISTS idx_anime_title_jp;
DROP INDEX IF EXISTS idx_anime_title_romaji;
DROP INDEX IF EXISTS idx_anime_title_kanji;
DROP INDEX IF EXISTS idx_episodes_anime_id;
DROP INDEX IF EXISTS idx_episodes_aired;
DROP INDEX IF EXISTS idx_episodes_episode_number;
DROP INDEX IF EXISTS idx_episodes_created_at;
DROP INDEX IF EXISTS idx_episodes_anime_aired;
DROP INDEX IF EXISTS idx_anime_character_anime_id;
DROP INDEX IF EXISTS idx_anime_character_name;
DROP INDEX IF EXISTS idx_anime_character_role;
DROP INDEX IF EXISTS idx_anime_staff_given_name;
DROP INDEX IF EXISTS idx_anime_staff_family_name;
DROP INDEX IF EXISTS idx_anime_staff_language;
DROP INDEX IF EXISTS idx_character_staff_character_id;
DROP INDEX IF EXISTS idx_character_staff_staff_id;
DROP INDEX IF EXISTS idx_anime_seasons_episode_count;
DROP INDEX IF EXISTS idx_anime_seasons_created_at;
DROP INDEX IF EXISTS idx_relations_anime_id;
DROP INDEX IF EXISTS idx_relations_related_anime_id;
