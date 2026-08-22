CREATE TABLE anime_staff (
                             id CHAR(36) NOT NULL PRIMARY KEY,
                             given_name VARCHAR(255) NOT NULL,
                             family_name VARCHAR(255) NOT NULL,
                             image TEXT,
                             birthday VARCHAR(255),
                             birth_place VARCHAR(255),
                             blood_type VARCHAR(255),
                             hobbies VARCHAR(255),
                             summary TEXT,
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

CREATE TRIGGER anime_staff_set_updated_at BEFORE UPDATE ON anime_staff
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
