DROP INDEX IF EXISTS "idx_anime_source_work_id";
ALTER TABLE "anime" DROP COLUMN IF EXISTS "source_work_id";
DROP INDEX IF EXISTS "idx_work_mal_id";
DROP TABLE IF EXISTS "work";
