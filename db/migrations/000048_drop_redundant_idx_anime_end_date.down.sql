-- Restores idx_anime_end_date. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_anime_end_date ON anime (end_date);
