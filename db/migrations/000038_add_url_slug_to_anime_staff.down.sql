DROP INDEX IF EXISTS idx_anime_staff_url_slug;
ALTER TABLE anime_staff DROP COLUMN IF EXISTS url_slug;
DROP FUNCTION IF EXISTS immutable_unaccent(text);
