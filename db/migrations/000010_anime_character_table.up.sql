CREATE TABLE anime_character (
                                 id CHAR(36) NOT NULL PRIMARY KEY,
                                 anime_id VARCHAR(36) NOT NULL,
                                 name VARCHAR(255) NOT NULL,
                                 role VARCHAR(255) NOT NULL,
                                 birthday VARCHAR(255),
                                 zodiac VARCHAR(255),
                                 gender VARCHAR(255),
                                 race VARCHAR(255),
                                 height VARCHAR(255),
                                 weight VARCHAR(255),
                                 title VARCHAR(255),
                                 martial_status VARCHAR(255),
                                 summary TEXT,
                                 image TEXT,
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

CREATE TRIGGER anime_character_set_updated_at BEFORE UPDATE ON anime_character
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
