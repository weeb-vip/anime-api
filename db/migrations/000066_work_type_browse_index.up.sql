-- Serves the /manga and /light-novels browse pages.
--
-- The query is `WHERE type = ? ORDER BY members DESC NULLS LAST, id`, paged.
-- Without an index that is a sequential scan of all 81,000 works followed by a
-- sort, on every page view -- and the rows are wide, synopsis included, so the
-- scan reads a great deal of heap to answer a question about three columns.
-- The same reasoning as idx_work_status_members in 000065.
--
-- Both ordering terms are in the index so it answers the ORDER BY as well as
-- the filter. NULLS LAST is spelled out because it is what the query asks for:
-- postgres puts nulls first on DESC by default, and an index built the other
-- way cannot serve it.
--
-- Only the default sort is indexed. SCORE, NEWEST and TITLE fall back to
-- scanning the type's own rows, which is a fraction of the table and is the
-- deliberate trade -- four indexes on a table a scraper writes to for days at
-- a stretch costs more on the write side than those three sorts are worth.
CREATE INDEX IF NOT EXISTS "idx_work_type_members"
    ON "work" ("type", "members" DESC NULLS LAST, "id");
