-- Drop the old indexes
DROP INDEX IF EXISTS idx_anime_season;
DROP INDEX IF EXISTS idx_anime_year;

-- Drop the old separate columns
ALTER TABLE anime DROP COLUMN year;
ALTER TABLE anime DROP COLUMN season;

-- Add the new combined season column
ALTER TABLE anime ADD COLUMN season VARCHAR(50) NULL;

-- Create index for the new column
CREATE INDEX idx_anime_season ON anime(season);