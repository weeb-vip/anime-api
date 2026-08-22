-- Create anime_seasons table
CREATE TABLE anime_seasons (
    id            VARCHAR(36) PRIMARY KEY DEFAULT (gen_random_uuid()::text),
    season        VARCHAR(255) NOT NULL,
    status text DEFAULT 'unknown' NOT NULL CHECK (status IN ('unknown', 'confirmed', 'announced', 'cancelled')),
    episode_count INTEGER,
    notes         TEXT,
    created_at    timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at    timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    anime_id      VARCHAR(36),
    CONSTRAINT UQ_anime_season UNIQUE (anime_id, season)
);

-- Create indexes
CREATE INDEX IDX_season ON anime_seasons (season);
CREATE INDEX IDX_status ON anime_seasons (status);
CREATE INDEX idx_anime_seasons_anime_id ON anime_seasons (anime_id);

-- Remove old season column from anime table
ALTER TABLE anime DROP COLUMN season;

-- MySQL's ON UPDATE CURRENT_TIMESTAMP has no Postgres equivalent, so the columns
-- that relied on it need a trigger to keep behaving the same way. Without it
-- updated_at silently stops advancing on UPDATE: nothing errors, the value just
-- goes stale.
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER anime_seasons_set_updated_at BEFORE UPDATE ON anime_seasons
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
