-- Widens work's title columns to text, matching the scraper.
--
-- MyAnimeList titles light novels well past 255 characters -- manga/182702 is
-- 282, a full sentence of plot -- and the scraper started writing those once
-- its manga heading fix gave title_en a fallback. Rows that long arrive here
-- through CDC and would be rejected by a varchar(255) column, stalling the
-- consumer on a row it can never write.
--
-- text rather than a larger varchar, because this table's sibling already
-- settled it: anime has had title_en and title_romaji as text since 000001.
-- work's varchar(255) from 000063 was the inconsistent one, and in postgres
-- the two are the same storage with the same performance.
ALTER TABLE "work" ALTER COLUMN "title_en" TYPE text;
ALTER TABLE "work" ALTER COLUMN "title_jp" TYPE text;
