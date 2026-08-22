-- url_slug is a voice actor's public URL segment: /people/mary-elizabeth-mcglynn
-- instead of /people/<uuid>.
--
-- Unlike anime.url_slug, this one is NOT minted in postgres and carried over
-- CDC. It does not need to be. A staff slug is a pure function of given_name
-- and family_name, and those are already unique: the scraper deduplicates
-- anime_staff on exactly that pair, so one row is one person and the slug needs
-- no minting authority to keep it distinct. Measured over 21,545 staff rows,
-- the expression below yields 21,541 distinct slugs; all four collisions are
-- the same person recorded twice under spelling variants the scraper's
-- exact-match dedup missed (Daniel García / Daniel Garcia, Julien Haggége /
-- Haggège), not two different people.
--
-- Deriving it here rather than upstream is what keeps this change to one
-- repository. A plain column would have needed something to populate rows
-- arriving over CDC -- a scheduled backfill, and a window in which a newly
-- synced voice actor has no slug and 404s. A STORED generated column has no
-- such window: MySQL computes it on insert and recomputes it on update, so
-- rows the sync writes are correct the moment they land, and a name change
-- moves the slug with it. character-staff-sync needs no change at all; GORM

-- In Postgres the whole REPLACE chain collapses to unaccent(). One caveat: the
-- one-argument unaccent() is STABLE, not IMMUTABLE, because it looks the default
-- dictionary up in the catalog -- and a generated column requires an immutable
-- expression. Naming the dictionary explicitly removes that lookup, which is
-- what makes the wrapper below honest enough to mark IMMUTABLE. It is the same
-- pattern the Postgres documentation uses for unaccent-based indexes.
CREATE EXTENSION IF NOT EXISTS unaccent;

CREATE OR REPLACE FUNCTION immutable_unaccent(text) RETURNS text
  LANGUAGE sql IMMUTABLE STRICT PARALLEL SAFE AS
$$ SELECT unaccent('unaccent', $1) $$;

ALTER TABLE anime_staff
    ADD COLUMN url_slug VARCHAR(255)
        GENERATED ALWAYS AS (
            NULLIF(
                TRIM(BOTH '-' FROM
                    REGEXP_REPLACE(
                        LOWER(immutable_unaccent(
                            COALESCE(given_name,'') || ' ' || COALESCE(family_name,'')
                        )),
                        '[^a-z0-9]+', '-', 'g'
                    )
                ),
            '')
        ) STORED;

-- Not unique, for the reason given above: this table is a CDC replica, and a
-- constraint here could only reject a write the source already accepted.
CREATE INDEX idx_anime_staff_url_slug ON anime_staff (url_slug);
