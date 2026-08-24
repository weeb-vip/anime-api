-- Restores idx_anime_anidbid. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_anime_anidbid ON anime (anidbid);
