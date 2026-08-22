CREATE TABLE anime_character_staff_link (
                                            id CHAR(36) NOT NULL PRIMARY KEY,
                                            character_id VARCHAR(36) NOT NULL,
                                            staff_id VARCHAR(36) NOT NULL,
                                            created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                            updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- MySQL's ON UPDATE CURRENT_TIMESTAMP has no Postgres equivalent, so the columns
-- that relied on it need a trigger to keep behaving the same way. Without it
-- updated_at silently stops advancing on UPDATE: nothing errors, the value just
-- goes stale.
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER anime_character_staff_link_set_updated_at BEFORE UPDATE ON anime_character_staff_link
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
