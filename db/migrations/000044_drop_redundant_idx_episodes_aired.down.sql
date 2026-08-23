-- Restores idx_episodes_aired. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_episodes_aired ON episodes (aired);
