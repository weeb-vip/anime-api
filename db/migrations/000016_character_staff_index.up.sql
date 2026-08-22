CREATE INDEX IF NOT EXISTS idx_character_id ON anime_character_staff_link (character_id);
CREATE INDEX IF NOT EXISTS idx_staff_id ON anime_character_staff_link (staff_id);
CREATE INDEX IF NOT EXISTS idx_character_staff ON anime_character_staff_link (character_id, staff_id);CREATE INDEX IF NOT EXISTS idx_anime_id ON anime_character (anime_id);
