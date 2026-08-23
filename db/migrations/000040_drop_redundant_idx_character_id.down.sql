-- Restores idx_character_id. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_character_id ON anime_character_staff_link (character_id);
