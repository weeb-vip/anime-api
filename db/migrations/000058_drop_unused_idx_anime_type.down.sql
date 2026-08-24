-- Restores idx_anime_type. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_anime_type ON anime (type);
