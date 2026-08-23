-- Restores idx_anime_id. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_anime_id ON anime_character (anime_id);
