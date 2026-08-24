-- Restores idx_anime_title_kanji. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_anime_title_kanji ON anime (title_kanji);
