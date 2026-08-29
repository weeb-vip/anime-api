ALTER TABLE "work" ALTER COLUMN "score" TYPE numeric(4,2);
DROP INDEX IF EXISTS "idx_work_url_slug";
ALTER TABLE "work" DROP COLUMN IF EXISTS "url_slug";
