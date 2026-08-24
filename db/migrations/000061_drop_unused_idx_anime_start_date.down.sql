-- Restores idx_anime_start_date. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_anime_start_date ON anime (start_date);
