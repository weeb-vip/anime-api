-- Restores idx_anime_rating_desc. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_anime_rating_desc ON anime (rating DESC);
