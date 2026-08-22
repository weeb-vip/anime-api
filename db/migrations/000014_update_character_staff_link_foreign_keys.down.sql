ALTER TABLE anime_character_staff_link
    DROP CONSTRAINT IF EXISTS fk_character_id,
    DROP CONSTRAINT IF EXISTS fk_staff_id;
