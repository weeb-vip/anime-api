-- Add back old season column to anime table
ALTER TABLE anime ADD COLUMN season VARCHAR(50) NULL;

-- Drop indexes
DROP INDEX IF EXISTS IDX_anime_id ON anime_seasons;
DROP INDEX IF EXISTS IDX_status ON anime_seasons;
DROP INDEX IF EXISTS IDX_season ON anime_seasons;

-- Drop anime_seasons table
DROP TABLE IF EXISTS anime_seasons;