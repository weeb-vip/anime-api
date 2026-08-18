-- url_slug is the anime's public URL segment, replacing /show/<uuid> with
-- /anime/<slug>. It is generated in postgres, the source of truth, and arrives
-- here over CDC like every other column -- this migration only makes somewhere
-- for it to land.
--
-- Nullable on purpose: rows exist here before the postgres backfill reaches
-- them, and a NOT NULL column would reject those CDC writes outright.
--
-- The unique index is what stops two anime ever claiming the same URL. MySQL
-- permits any number of NULLs under a unique index, so unbackfilled rows do not
-- collide with each other while the backfill is in flight.
ALTER TABLE anime ADD COLUMN url_slug VARCHAR(255) NULL;
CREATE UNIQUE INDEX idx_anime_url_slug ON anime (url_slug);
