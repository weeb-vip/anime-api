-- Adds `work` and `anime.source_work_id`: the manga, light novels and novels
-- anime are adapted from, and the pointer to them.
--
-- This is the read-store half of the scraper's `work` table. The scraper owns
-- the data; Debezium already captures the table with no configuration change,
-- because it runs with schema.include.list=public and no table filter, and
-- ANIMEDB already claims anime-db.> so the subject and its retention exist
-- too. What did not exist is somewhere for it to land, which is this.
--
-- Why it matters: `anime.source` records a category, never an identity. We
-- know 6,110 anime came from a manga and cannot say which one, so
-- re-adaptations of the same source look unrelated -- Fruits Basket 2001 and
-- 2019, Hunter x Hunter 1999 and 2011, Fullmetal Alchemist and Brotherhood.
-- They share no TheTVDB series id and frequently no cast, so both signals
-- behind relatedAnime miss them. The source work is the only thing connecting
-- them.

CREATE TABLE IF NOT EXISTS "work" (
    "id"              varchar(36) NOT NULL,
    "mal_id"          integer,
    "type"            varchar(32) NOT NULL DEFAULT 'MANGA',
    "title_en"        varchar(255),
    "title_jp"        varchar(255),
    "title_synonyms"  text,
    "synopsis"        text,
    "image_url"       varchar(512),
    "status"          varchar(64),
    "volumes"         integer,
    "chapters"        integer,
    "published_from"  timestamptz,
    "published_to"    timestamptz,
    "demographic"     varchar(64),
    "serialization"   varchar(255),
    "authors"         text,
    "score"           numeric(4,2),
    "ranking"         integer,
    "members"         integer,
    "favorites"       integer,
    "created_at"      timestamptz NOT NULL DEFAULT now(),
    "updated_at"      timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT "work_pkey" PRIMARY KEY ("id")
);

-- varchar(36), not char(36) or uuid. The rest of this schema is varchar(36)
-- and migration 000062 exists precisely because four char(36) columns forced a
-- cast that made their indexes unusable. A uuid column here would repeat that
-- against anime.source_work_id.
CREATE UNIQUE INDEX IF NOT EXISTS "idx_work_mal_id" ON "work" ("mal_id");

ALTER TABLE "anime" ADD COLUMN IF NOT EXISTS "source_work_id" varchar(36);

-- The index that earns its keep is this one, not the forward lookup. An anime's
-- own source is already on its row; what needs finding is every other anime
-- adapted from the same work, which is the relation itself.
CREATE INDEX IF NOT EXISTS "idx_anime_source_work_id" ON "anime" ("source_work_id");

-- No foreign key, matching the scraper. CDC makes no ordering promise across
-- tables, so an anime routinely arrives before the work it points at; anime-sync
-- already discards foreign key violations for exactly this reason, and a
-- constraint here would turn ordinary event ordering into dropped rows.
