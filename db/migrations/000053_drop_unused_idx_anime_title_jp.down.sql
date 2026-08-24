-- Restores idx_anime_title_jp. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_anime_title_jp ON anime (title_jp);
