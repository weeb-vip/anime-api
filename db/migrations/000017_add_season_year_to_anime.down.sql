DROP INDEX IF EXISTS idx_anime_season;
DROP INDEX IF EXISTS idx_anime_year;
DROP INDEX IF EXISTS idx_anime_season;

ALTER TABLE anime DROP COLUMN year;
ALTER TABLE anime DROP COLUMN season;