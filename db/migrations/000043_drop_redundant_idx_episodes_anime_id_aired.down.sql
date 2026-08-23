-- Restores idx_episodes_anime_id_aired. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_episodes_anime_id_aired ON episodes (anime_id, aired);
