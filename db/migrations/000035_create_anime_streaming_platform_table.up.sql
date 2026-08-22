CREATE TABLE anime_streaming_platform (
    id VARCHAR(36) PRIMARY KEY DEFAULT (gen_random_uuid()::text),
    anime_id VARCHAR(36) NOT NULL,
    platform VARCHAR(100) NOT NULL,
    name VARCHAR(255) NULL,
    url TEXT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX idx_streaming_anime_platform ON anime_streaming_platform (anime_id, platform);

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

CREATE TRIGGER anime_streaming_platform_set_updated_at BEFORE UPDATE ON anime_streaming_platform
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
