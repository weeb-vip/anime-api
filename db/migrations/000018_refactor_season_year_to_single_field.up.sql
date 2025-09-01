-- Drop the old indexes
DROP INDEX IF EXISTS idx_anime_season;
DROP INDEX IF EXISTS idx_anime_year;  
DROP INDEX IF EXISTS idx_anime_season;

-- Drop the old separate columns
ALTER TABLE anime DROP COLUMN IF EXISTS year;
ALTER TABLE anime DROP COLUMN IF EXISTS season;

-- Add the new combined season column
ALTER TABLE anime ADD COLUMN season VARCHAR(50) NULL;

-- Create index for the new column
CREATE INDEX idx_anime_season ON anime(season);