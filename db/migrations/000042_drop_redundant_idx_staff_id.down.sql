-- Restores idx_staff_id. Also alone and CONCURRENTLY, for the same reason.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_staff_id ON anime_character_staff_link (staff_id);
