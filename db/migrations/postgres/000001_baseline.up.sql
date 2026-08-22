-- Postgres baseline for anime-api.
--
-- Replaces 37 MySQL migrations with one schema rather than porting the history.
-- Two reasons that is not laziness: several of them do not translate one for
-- one, and the chain cannot build a database from scratch anyway -- the trigger
-- migrations (000026-000028) are bare CREATE TRIGGER ... BEGIN ... END with no
-- DELIMITER, so any client that splits on ';' dies partway through.
--
-- Generated from the schema those migrations actually produce, captured out of
-- information_schema after applying them to a throwaway MySQL. Translating the
-- migration text by hand would have carried their bugs across; see the
-- idx_relations_* note below for one it caught.

-- Index changes made deliberately during conversion:
--   anime_character (anime_id): kept idx_anime_character_anime_id, dropped idx_anime_id
--   anime_character_staff_link (character_id): kept idx_character_staff_character_id, dropped idx_character_id
--   anime_character_staff_link (staff_id): kept idx_character_staff_staff_id, dropped idx_staff_id
--   episodes (anime_id, aired): kept idx_episodes_anime_id_aired, dropped idx_episodes_anime_aired
--   tags (name): kept name UNIQUE, dropped idx_tags_name

CREATE EXTENSION IF NOT EXISTS unaccent;

-- unaccent() is STABLE, not IMMUTABLE: the one-argument form looks up the
-- default dictionary in the catalog, so Postgres refuses it in a generated
-- column. Naming the dictionary explicitly removes that lookup, which is what
-- makes this wrapper honest enough to mark IMMUTABLE -- the same pattern the
-- Postgres docs use for unaccent-based indexes.
--
-- The caveat is real: if the unaccent dictionary is ever redefined, stored
-- values here will not recompute. That is a deployment-time concern, not a
-- runtime one, and the alternative is not having the column at all.
CREATE OR REPLACE FUNCTION immutable_unaccent(text) RETURNS text
  LANGUAGE sql IMMUTABLE STRICT PARALLEL SAFE AS
$$ SELECT unaccent('unaccent', $1) $$;

-- ----------------------------------------------------------------------
CREATE TABLE anime (
  id varchar(36) NOT NULL,
  type varchar(30) DEFAULT 'Anime',
  title_en text,
  title_jp text,
  title_romaji text,
  title_kanji text,
  title_synonyms text,
  image_url text,
  synopsis text,
  episodes integer,
  status text,
  ranking integer,
  genres text,
  duration text,
  broadcast text,
  source text,
  licensors text,
  studios text,
  created_at timestamptz,
  updated_at timestamptz,
  anidbid varchar(10),
  mal_id integer,
  thetvdbid varchar(255),
  start_date varchar(255),
  end_date varchar(255),
  rating numeric(3,1),
  url_slug varchar(255),
  PRIMARY KEY (id)
);

CREATE INDEX idx_anime_anidbid ON anime (anidbid);
CREATE INDEX idx_anime_created_at ON anime (created_at);
CREATE INDEX idx_anime_created_at_id ON anime (created_at, id);
CREATE INDEX idx_anime_end_date ON anime (end_date);
CREATE INDEX idx_anime_end_date_id ON anime (end_date, id);
CREATE INDEX idx_anime_mal_id ON anime (mal_id);
CREATE INDEX idx_anime_ranking ON anime (ranking);
CREATE INDEX idx_anime_ranking_id ON anime (ranking, id);
CREATE INDEX idx_anime_rating_desc ON anime (rating);
CREATE INDEX idx_anime_rating_id ON anime (rating, id);
CREATE INDEX idx_anime_source ON anime (source);
CREATE INDEX idx_anime_start_date ON anime (start_date);
CREATE INDEX idx_anime_status ON anime (status);
CREATE INDEX idx_anime_thetvdbid ON anime (thetvdbid);
CREATE INDEX idx_anime_title_en ON anime (title_en);
CREATE INDEX idx_anime_title_jp ON anime (title_jp);
CREATE INDEX idx_anime_title_kanji ON anime (title_kanji);
CREATE INDEX idx_anime_title_romaji ON anime (title_romaji);
CREATE INDEX idx_anime_type ON anime (type);
CREATE INDEX idx_anime_url_slug ON anime (url_slug);

-- ----------------------------------------------------------------------
CREATE TABLE anime_character (
  id char(36) NOT NULL,
  anime_id varchar(36) NOT NULL,
  name varchar(255) NOT NULL,
  role varchar(255) NOT NULL,
  birthday varchar(255),
  zodiac varchar(255),
  gender varchar(255),
  race varchar(255),
  height varchar(255),
  weight varchar(255),
  title varchar(255),
  martial_status varchar(255),
  summary text,
  image text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (id)
);

CREATE INDEX idx_anime_character_anime_id ON anime_character (anime_id);
CREATE INDEX idx_anime_character_name ON anime_character (name);
CREATE INDEX idx_anime_character_role ON anime_character (role);

-- ----------------------------------------------------------------------
CREATE TABLE anime_character_staff_link (
  id char(36) NOT NULL,
  character_id varchar(36) NOT NULL,
  staff_id varchar(36) NOT NULL,
  character_name varchar(255) NOT NULL,
  staff_given_name varchar(255) NOT NULL,
  staff_family_name varchar(255) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (id)
);

CREATE INDEX idx_character_staff_character_id ON anime_character_staff_link (character_id);
CREATE INDEX idx_character_staff ON anime_character_staff_link (character_id, staff_id);
CREATE INDEX idx_character_staff_staff_id ON anime_character_staff_link (staff_id);

-- ----------------------------------------------------------------------
CREATE TABLE anime_relations (
  id char(36) NOT NULL,
  anime_id varchar(36) NOT NULL,
  related_anime_id varchar(36) NOT NULL,
  relation_type varchar(30),
  created_at timestamptz DEFAULT now(),
  PRIMARY KEY (id)
);


-- ----------------------------------------------------------------------
CREATE TABLE anime_schedule (
  id varchar(36) NOT NULL DEFAULT 'uuid()',
  anime_id varchar(36) NOT NULL,
  animeschedule_route varchar(255),
  jpn_time timestamptz,
  sub_time timestamptz,
  dub_time timestamptz,
  premier timestamptz,
  sub_premier timestamptz,
  dub_premier timestamptz,
  notes text,
  delayed_timetable varchar(50),
  sub_delayed_timetable varchar(50),
  dub_delayed_timetable varchar(50),
  episode_override_date timestamptz,
  episode_override_episode integer,
  episode_override_episodes_aired integer,
  sub_episode_override_date timestamptz,
  sub_episode_override_episode integer,
  sub_episode_override_episodes_aired integer,
  dub_episode_override_date timestamptz,
  dub_episode_override_episode integer,
  dub_episode_override_episodes_aired integer,
  last_synced_at timestamptz NOT NULL DEFAULT now(),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_anime_schedule_anime_id ON anime_schedule (anime_id);
CREATE INDEX idx_anime_schedule_route ON anime_schedule (animeschedule_route);

-- ----------------------------------------------------------------------
CREATE TABLE anime_seasons (
  id varchar(36) NOT NULL DEFAULT 'uuid()',
  season varchar(255) NOT NULL,
  status text NOT NULL DEFAULT 'unknown',
  episode_count integer,
  notes text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  anime_id varchar(36) NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT anime_seasons_status_check CHECK (status IN ('unknown', 'confirmed', 'announced', 'cancelled'))
);

CREATE INDEX idx_anime_id ON anime_seasons (anime_id);
CREATE INDEX idx_anime_seasons_created_at ON anime_seasons (created_at);
CREATE INDEX idx_anime_seasons_episode_count ON anime_seasons (episode_count);
CREATE INDEX idx_anime_seasons_season_anime_id ON anime_seasons (season, anime_id);
CREATE INDEX idx_season ON anime_seasons (season);
CREATE INDEX idx_status ON anime_seasons (status);
CREATE UNIQUE INDEX uq_anime_season ON anime_seasons (anime_id, season);

-- ----------------------------------------------------------------------
CREATE TABLE anime_staff (
  id char(36) NOT NULL,
  given_name varchar(255) NOT NULL,
  family_name varchar(255) NOT NULL,
  image text,
  birthday varchar(255),
  birth_place varchar(255),
  blood_type varchar(255),
  hobbies varchar(255),
  summary text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  language varchar(30),
  url_slug varchar(255) GENERATED ALWAYS AS (nullif(trim(both '-' from regexp_replace(lower(immutable_unaccent(coalesce(given_name,'') || ' ' || coalesce(family_name,''))), '[^a-z0-9]+', '-', 'g')), '')) STORED,
  PRIMARY KEY (id)
);

CREATE INDEX idx_anime_staff_family_name ON anime_staff (family_name);
CREATE INDEX idx_anime_staff_given_name ON anime_staff (given_name);
CREATE INDEX idx_anime_staff_language ON anime_staff (language);
CREATE INDEX idx_anime_staff_url_slug ON anime_staff (url_slug);

-- ----------------------------------------------------------------------
CREATE TABLE anime_streaming_platform (
  id varchar(36) NOT NULL DEFAULT 'uuid()',
  anime_id varchar(36) NOT NULL,
  platform varchar(100) NOT NULL,
  name varchar(255),
  url text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_streaming_anime_platform ON anime_streaming_platform (anime_id, platform);

-- ----------------------------------------------------------------------
CREATE TABLE anime_tags (
  anime_id varchar(36) NOT NULL,
  tag_id bigint NOT NULL,
  created_at timestamptz DEFAULT now(),
  PRIMARY KEY (anime_id, tag_id)
);

CREATE INDEX idx_anime_tags_anime_id ON anime_tags (anime_id);
CREATE INDEX idx_anime_tags_tag_id ON anime_tags (tag_id);

-- ----------------------------------------------------------------------
CREATE TABLE episode_air_time (
  id varchar(36) NOT NULL DEFAULT 'uuid()',
  anime_id varchar(36) NOT NULL,
  episode_number integer NOT NULL,
  air_type text NOT NULL,
  air_datetime timestamptz NOT NULL,
  streams_json jsonb,
  last_synced_at timestamptz NOT NULL DEFAULT now(),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  CONSTRAINT episode_air_time_air_type_check CHECK (air_type IN ('raw', 'sub', 'dub'))
);

CREATE INDEX idx_episode_air_time_air_datetime ON episode_air_time (air_datetime);
CREATE UNIQUE INDEX idx_episode_air_time_unique ON episode_air_time (anime_id, episode_number, air_type);

-- ----------------------------------------------------------------------
CREATE TABLE episodes (
  id varchar(36) NOT NULL,
  anime_id varchar(36),
  episode integer,
  title_en text,
  title_jp text,
  synopsis text,
  created_at timestamptz,
  updated_at timestamptz,
  backup_aired varchar(255),
  aired date,
  PRIMARY KEY (id)
);

CREATE INDEX idx_episodes_aired ON episodes (aired);
CREATE INDEX idx_episodes_aired_anime_id ON episodes (aired, anime_id);
CREATE INDEX idx_episodes_anime_id_aired ON episodes (anime_id, aired);
CREATE INDEX idx_episodes_anime_id ON episodes (anime_id);
CREATE INDEX idx_episodes_anime_id_episode ON episodes (anime_id, episode);
CREATE INDEX idx_episodes_created_at ON episodes (created_at);
CREATE INDEX idx_episodes_episode_number ON episodes (episode);

-- ----------------------------------------------------------------------
CREATE TABLE tags (
  id bigint GENERATED BY DEFAULT AS IDENTITY,
  name varchar(100) NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX name ON tags (name);

-- ----------------------------------------------------------------------
-- MySQL's ON UPDATE CURRENT_TIMESTAMP has no Postgres equivalent, so the
-- columns that relied on it need a trigger to keep behaving the same way.
-- Without this they would silently stop advancing on UPDATE -- a difference
-- that does not fail anything, it just quietly makes updated_at wrong.
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER anime_character_set_updated_at BEFORE UPDATE ON anime_character
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER anime_character_staff_link_set_updated_at BEFORE UPDATE ON anime_character_staff_link
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER anime_schedule_set_updated_at BEFORE UPDATE ON anime_schedule
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER anime_seasons_set_updated_at BEFORE UPDATE ON anime_seasons
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER anime_staff_set_updated_at BEFORE UPDATE ON anime_staff
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER anime_streaming_platform_set_updated_at BEFORE UPDATE ON anime_streaming_platform
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER episode_air_time_set_updated_at BEFORE UPDATE ON episode_air_time
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER tags_set_updated_at BEFORE UPDATE ON tags
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ----------------------------------------------------------------------
-- anime.episodes is a denormalised count maintained by triggers on episodes.
-- One plpgsql function covers all three events; the MySQL original was three
-- near-identical trigger bodies.
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
CREATE TRIGGER update_anime_episode_count_after_update
  AFTER UPDATE ON episodes
  FOR EACH ROW EXECUTE FUNCTION update_anime_episode_count();
CREATE TRIGGER update_anime_episode_count_after_delete
  AFTER DELETE ON episodes
  FOR EACH ROW EXECUTE FUNCTION update_anime_episode_count();

-- ----------------------------------------------------------------------
-- 000021 tried to index a table named 'relations'. The table is
-- anime_relations (000009), so these two indexes were never created -- not
-- here and not in production, where that migration must have failed at the
-- same line. Created correctly here rather than faithfully reproducing the
-- typo.
CREATE INDEX idx_relations_anime_id ON anime_relations (anime_id);
CREATE INDEX idx_relations_related_anime_id ON anime_relations (related_anime_id);
