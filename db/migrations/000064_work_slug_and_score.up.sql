-- Brings the read store's `work` into line with the scraper's, which gained a
-- url_slug and changed score's type after 000063 landed.
--
-- url_slug is the public URL segment behind /manga/<slug>. It is assigned by a
-- trigger in the scraper and never recomputed, so a title correction upstream
-- cannot silently move a live URL. No trigger is needed here: this side only
-- ever receives the value the scraper already decided, and generating a second
-- one would produce a different slug for the same work.
ALTER TABLE "work" ADD COLUMN IF NOT EXISTS "url_slug" varchar(255);

-- Unique because it is a URL. Nullable rows do not collide under a unique index
-- in postgres, so rows that arrive before the scraper has assigned one sit
-- happily alongside each other.
CREATE UNIQUE INDEX IF NOT EXISTS "idx_work_url_slug" ON "work" ("url_slug");

-- double precision, matching the scraper.
--
-- Debezium's decimal.handling.mode defaults to `precise`, which encodes numeric
-- columns as base64 bytes and a scale rather than as a number. The consumer
-- reads this into a float, so the column has to be a float on both sides --
-- see anime_scraper migration 1787360400000 for the full reasoning.
--
-- Empty table, so the type change rewrites nothing.
ALTER TABLE "work" ALTER COLUMN "score" TYPE double precision;
