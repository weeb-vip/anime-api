-- Restores idx_anime_ranking. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_anime_ranking ON anime (ranking);
