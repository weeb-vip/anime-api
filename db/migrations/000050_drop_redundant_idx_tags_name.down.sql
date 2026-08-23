-- Restores idx_tags_name. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tags_name ON tags (name);
