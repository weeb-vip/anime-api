CREATE TABLE episode_air_time (
    id VARCHAR(36) PRIMARY KEY DEFAULT (gen_random_uuid()::text),
    anime_id VARCHAR(36) NOT NULL,
    episode_number INT NOT NULL,
    air_type text NOT NULL CHECK (air_type IN ('raw', 'sub', 'dub')),
    air_datetime timestamptz NOT NULL,
    streams_json JSON NULL,
    last_synced_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_episode_air_time_air_datetime ON episode_air_time (air_datetime);
CREATE UNIQUE INDEX idx_episode_air_time_unique ON episode_air_time (anime_id, episode_number, air_type);

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

CREATE TRIGGER episode_air_time_set_updated_at BEFORE UPDATE ON episode_air_time
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
