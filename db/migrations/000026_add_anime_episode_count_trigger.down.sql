DROP TRIGGER IF EXISTS update_anime_episode_count_after_insert ON episodes;
-- The function is dropped by 000026 rather than 000028 because it is created
-- here; 000027 and 000028 only attach triggers to it. CASCADE is deliberately
-- not used -- if a trigger still references it, that is worth an error.
DROP FUNCTION IF EXISTS update_anime_episode_count();
