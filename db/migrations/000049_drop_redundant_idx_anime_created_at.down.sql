-- Restores idx_anime_created_at. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_anime_created_at ON anime (created_at);
