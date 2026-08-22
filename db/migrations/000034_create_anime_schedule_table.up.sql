CREATE TABLE anime_schedule (
    id VARCHAR(36) PRIMARY KEY DEFAULT (gen_random_uuid()::text),
    anime_id VARCHAR(36) NOT NULL,
    animeschedule_route VARCHAR(255) NULL,
    jpn_time timestamptz NULL,
    sub_time timestamptz NULL,
    dub_time timestamptz NULL,
    premier timestamptz NULL,
    sub_premier timestamptz NULL,
    dub_premier timestamptz NULL,
    notes TEXT NULL,
    delayed_timetable VARCHAR(50) NULL,
    sub_delayed_timetable VARCHAR(50) NULL,
    dub_delayed_timetable VARCHAR(50) NULL,
    episode_override_date timestamptz NULL,
    episode_override_episode INT NULL,
    episode_override_episodes_aired INT NULL,
    sub_episode_override_date timestamptz NULL,
    sub_episode_override_episode INT NULL,
    sub_episode_override_episodes_aired INT NULL,
    dub_episode_override_date timestamptz NULL,
    dub_episode_override_episode INT NULL,
    dub_episode_override_episodes_aired INT NULL,
    last_synced_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX idx_anime_schedule_anime_id ON anime_schedule (anime_id);
CREATE INDEX idx_anime_schedule_route ON anime_schedule (animeschedule_route);

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

CREATE TRIGGER anime_schedule_set_updated_at BEFORE UPDATE ON anime_schedule
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
