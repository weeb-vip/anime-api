-- Reverse of the baseline.
--
-- Dropping the tables is enough for everything except the two functions, which
-- live at schema level rather than being owned by any table. immutable_unaccent
-- in particular would survive a table drop and then collide with a re-run of the
-- up migration.
--
-- The unaccent extension is deliberately NOT dropped: other schemas may use it,
-- and creating it needs privileges that dropping it does not, so removing it
-- here can leave a database that cannot migrate up again.

DROP TABLE IF EXISTS anime_tags;
DROP TABLE IF EXISTS anime_character_staff_link;
DROP TABLE IF EXISTS anime_character;
DROP TABLE IF EXISTS anime_staff;
DROP TABLE IF EXISTS anime_streaming_platform;
DROP TABLE IF EXISTS anime_seasons;
DROP TABLE IF EXISTS anime_relations;
DROP TABLE IF EXISTS anime_schedule;
DROP TABLE IF EXISTS episode_air_time;
DROP TABLE IF EXISTS episodes;
DROP TABLE IF EXISTS anime;
DROP TABLE IF EXISTS tags;

DROP FUNCTION IF EXISTS update_anime_episode_count();
DROP FUNCTION IF EXISTS set_updated_at();
DROP FUNCTION IF EXISTS immutable_unaccent(text);
