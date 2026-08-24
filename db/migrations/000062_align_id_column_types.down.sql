-- Restores char(36).
--
-- Note this reintroduces the problem: joins from anime_character_staff_link
-- back to these tables will cast the indexed column again and fall back to a
-- sequential scan of every link row.
ALTER TABLE anime_relations ALTER COLUMN id TYPE char(36);
ALTER TABLE anime_character_staff_link ALTER COLUMN id TYPE char(36);
ALTER TABLE anime_staff ALTER COLUMN id TYPE char(36);
ALTER TABLE anime_character ALTER COLUMN id TYPE char(36);
