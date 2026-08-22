-- Attaches the DELETE event to the function 000026 defines.
CREATE TRIGGER update_anime_episode_count_after_delete
  AFTER DELETE ON episodes
  FOR EACH ROW EXECUTE FUNCTION update_anime_episode_count();
