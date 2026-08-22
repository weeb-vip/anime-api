-- Keep anime.episodes in step with the rows in episodes.
--
-- MySQL spelled the body inline in each of three triggers. Postgres separates
-- the function from the triggers that call it, so the shared body is defined
-- once here and 000027 and 000028 attach their own events to it.
--
-- COALESCE(NEW.anime_id, OLD.anime_id) lets one body serve INSERT, UPDATE and
-- DELETE: NEW is null on delete, OLD is null on insert.
CREATE OR REPLACE FUNCTION update_anime_episode_count() RETURNS trigger AS $$
DECLARE
  target text;
BEGIN
  target := COALESCE(NEW.anime_id, OLD.anime_id);
  UPDATE anime SET episodes = (
    SELECT COUNT(*) FROM episodes WHERE anime_id = target
  ) WHERE id = target;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_anime_episode_count_after_insert
  AFTER INSERT ON episodes
  FOR EACH ROW EXECUTE FUNCTION update_anime_episode_count();
