-- Restores idx_anime_title_romaji. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_anime_title_romaji ON anime (title_romaji);
