-- Restores idx_anime_mal_id. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_anime_mal_id ON anime (mal_id);
