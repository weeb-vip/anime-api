-- Add database indexes for improved query performance.
--
-- MySQL had no CREATE INDEX IF NOT EXISTS, so this migration defined a stored
-- procedure to emulate it and called it once per index. Postgres has the real
-- thing, so the procedure and its 30 CALLs collapse to 30 plain statements.
--
-- Prefix lengths are gone too: MySQL required (status(50)) to index a TEXT
-- column at all. Postgres indexes text directly, so the truncation those
-- lengths imposed -- and the partial matches it caused -- simply do not apply.

CREATE INDEX IF NOT EXISTS idx_anime_status ON anime (status);
CREATE INDEX IF NOT EXISTS idx_anime_rating ON anime (rating);
CREATE INDEX IF NOT EXISTS idx_anime_ranking ON anime (ranking);
CREATE INDEX IF NOT EXISTS idx_anime_created_at ON anime (created_at);
CREATE INDEX IF NOT EXISTS idx_anime_type ON anime (type);
CREATE INDEX IF NOT EXISTS idx_anime_source ON anime (source);
CREATE INDEX IF NOT EXISTS idx_anime_start_date ON anime (start_date);
CREATE INDEX IF NOT EXISTS idx_anime_end_date ON anime (end_date);
CREATE INDEX IF NOT EXISTS idx_anime_anidbid ON anime (anidbid);
CREATE INDEX IF NOT EXISTS idx_anime_thetvdbid ON anime (thetvdbid);

CREATE INDEX IF NOT EXISTS idx_anime_title_en ON anime (title_en);
CREATE INDEX IF NOT EXISTS idx_anime_title_jp ON anime (title_jp);
CREATE INDEX IF NOT EXISTS idx_anime_title_romaji ON anime (title_romaji);
CREATE INDEX IF NOT EXISTS idx_anime_title_kanji ON anime (title_kanji);

CREATE INDEX IF NOT EXISTS idx_episodes_anime_id ON episodes (anime_id);
CREATE INDEX IF NOT EXISTS idx_episodes_aired ON episodes (aired);
CREATE INDEX IF NOT EXISTS idx_episodes_episode_number ON episodes (episode);
CREATE INDEX IF NOT EXISTS idx_episodes_created_at ON episodes (created_at);
CREATE INDEX IF NOT EXISTS idx_episodes_anime_aired ON episodes (anime_id, aired);

CREATE INDEX IF NOT EXISTS idx_anime_character_anime_id ON anime_character (anime_id);
CREATE INDEX IF NOT EXISTS idx_anime_character_name ON anime_character (name);
CREATE INDEX IF NOT EXISTS idx_anime_character_role ON anime_character (role);

CREATE INDEX IF NOT EXISTS idx_anime_staff_given_name ON anime_staff (given_name);
CREATE INDEX IF NOT EXISTS idx_anime_staff_family_name ON anime_staff (family_name);
CREATE INDEX IF NOT EXISTS idx_anime_staff_language ON anime_staff (language);

CREATE INDEX IF NOT EXISTS idx_character_staff_character_id ON anime_character_staff_link (character_id);
CREATE INDEX IF NOT EXISTS idx_character_staff_staff_id ON anime_character_staff_link (staff_id);

CREATE INDEX IF NOT EXISTS idx_anime_seasons_episode_count ON anime_seasons (episode_count);
CREATE INDEX IF NOT EXISTS idx_anime_seasons_created_at ON anime_seasons (created_at);

-- The original targeted a table called "relations", which does not exist -- the
-- table 000009 creates is anime_relations. Its own comment hedges with "(if
-- exists)", so the intent was conditional; that intent is expressed here rather
-- than the typo reproduced. Guarded because a bare CREATE INDEX on a missing
-- table is a hard error in Postgres, where MySQL's dynamic SQL only failed at
-- CALL time.
DO $$
BEGIN
  IF to_regclass('public.anime_relations') IS NOT NULL THEN
    CREATE INDEX IF NOT EXISTS idx_relations_anime_id ON anime_relations (anime_id);
    CREATE INDEX IF NOT EXISTS idx_relations_related_anime_id ON anime_relations (related_anime_id);
  END IF;
END $$;
