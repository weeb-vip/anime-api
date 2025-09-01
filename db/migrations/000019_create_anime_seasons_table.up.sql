-- Create enum for anime season status
CREATE TYPE anime_seasons_status_enum AS ENUM ('unknown', 'confirmed', 'announced', 'cancelled');

-- Create anime_seasons table
CREATE TABLE anime_seasons (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    season        VARCHAR NOT NULL,
    status        anime_seasons_status_enum DEFAULT 'unknown'::anime_seasons_status_enum NOT NULL,
    episode_count INTEGER,
    notes         TEXT,
    created_at    TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at    TIMESTAMP DEFAULT NOW() NOT NULL,
    anime_id      UUID,
    CONSTRAINT "FK_anime_seasons_anime_id" FOREIGN KEY (anime_id) REFERENCES anime(id) ON DELETE CASCADE,
    CONSTRAINT "UQ_anime_season" UNIQUE (anime_id, season)
);

-- Create indexes
CREATE INDEX "IDX_season" ON anime_seasons (season);
CREATE INDEX "IDX_status" ON anime_seasons (status);
CREATE INDEX "IDX_anime_id" ON anime_seasons (anime_id);

-- Remove old season column from anime table
ALTER TABLE anime DROP COLUMN IF EXISTS season;