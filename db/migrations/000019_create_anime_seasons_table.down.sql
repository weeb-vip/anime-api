-- Add back old season column to anime table
ALTER TABLE anime ADD COLUMN season VARCHAR(50) NULL;

-- Drop indexes
DROP INDEX IF EXISTS "IDX_anime_id";
DROP INDEX IF EXISTS "IDX_status";
DROP INDEX IF EXISTS "IDX_season";

-- Drop anime_seasons table
DROP TABLE IF EXISTS anime_seasons;

-- Drop enum type
DROP TYPE IF EXISTS anime_seasons_status_enum;