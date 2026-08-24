-- Restores idx_anime_title_en. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_anime_title_en ON anime (title_en);
