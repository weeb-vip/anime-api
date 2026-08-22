ALTER TABLE anime
    ADD COLUMN mal_id INT NULL;
CREATE INDEX IF NOT EXISTS idx_anime_mal_id ON anime (mal_id);
