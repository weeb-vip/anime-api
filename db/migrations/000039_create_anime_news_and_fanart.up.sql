-- anime_news and anime_fanart, which existed in production without ever being
-- created by a migration.
--
-- news-ingest's store package says so itself: "writes into the anime-api-owned
-- MySQL tables (anime_news, anime_fanart). anime-api creates/migrates these
-- tables; we only upsert rows." anime-api did not, so both were made by hand at
-- some point and no migration chain has ever produced them.
--
-- That was survivable while every environment was long-lived. It stops being
-- survivable on Postgres, where the database is built from these files: without
-- this, anime-api starts against a database with no anime_fanart and the root
-- handler's fanart repository fails on first query.
--
-- Column definitions are taken from the production schema dumped out of
-- PlanetScale, so this reproduces what is actually there rather than what the
-- Go structs imply.

CREATE TABLE anime_news (
  id              varchar(40) NOT NULL,
  anime_id        varchar(36) NOT NULL,
  mal_id          integer,
  title           varchar(512) NOT NULL,
  summary         text,
  category        varchar(32) NOT NULL,
  source_url      text,
  source_name     varchar(255),
  language        varchar(8),
  -- MySQL json becomes jsonb. The Go struct carries it as a marshalled string
  -- either way, so this changes storage and indexing, not the code.
  reference_links jsonb,
  published_date  date,
  episode_number  integer,
  title_slug      varchar(255),
  researched_at   timestamptz,
  created_at      timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

-- The dedupe key is what makes the ingest idempotent: re-running over the same
-- feed upserts rather than duplicating.
CREATE UNIQUE INDEX idx_news_dedupe ON anime_news (anime_id, published_date, title_slug);
CREATE INDEX idx_news_anime ON anime_news (anime_id);
CREATE INDEX idx_news_latest ON anime_news (published_date DESC, id);

CREATE TABLE anime_fanart (
  id         varchar(40) NOT NULL,
  anime_id   varchar(36) NOT NULL,
  image_url  text NOT NULL,
  source_url text,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

CREATE INDEX idx_fanart_anime ON anime_fanart (anime_id);

-- anime_news.updated_at used ON UPDATE CURRENT_TIMESTAMP, which Postgres has no
-- equivalent for. set_updated_at() already exists by this point; CREATE OR
-- REPLACE keeps this migration standalone regardless.
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER anime_news_set_updated_at BEFORE UPDATE ON anime_news
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
