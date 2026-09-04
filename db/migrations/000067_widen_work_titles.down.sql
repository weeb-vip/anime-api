-- Truncating rather than failing: by the time anyone runs this down there are
-- likely titles that no longer fit, and a migration that cannot be reversed
-- without hand-editing data is worse than a lossy one.
UPDATE "work" SET "title_en" = left("title_en", 255) WHERE length("title_en") > 255;
UPDATE "work" SET "title_jp" = left("title_jp", 255) WHERE length("title_jp") > 255;
ALTER TABLE "work" ALTER COLUMN "title_en" TYPE character varying(255);
ALTER TABLE "work" ALTER COLUMN "title_jp" TYPE character varying(255);
