-- Serves the homepage's "still publishing" row.
--
-- The query is `WHERE LOWER(status) = 'publishing' ORDER BY members DESC, id`,
-- and without an index it is a sequential scan over the whole table followed by
-- a top-N sort: 73ms measured on staging against 2,011 rows, which is far more
-- than a table this size should cost. The rows are wide -- synopsis is in
-- there -- so the scan reads a lot of heap to answer a question about three
-- columns.
--
-- Indexed on the expression rather than on `status` directly because the query
-- lowercases it. status is a scraped label rather than an enum, so the
-- comparison is deliberately case-insensitive and a plain btree on the raw
-- column could not be used.
--
-- All three ordering terms are included so the index answers the ORDER BY as
-- well as the filter, leaving nothing to sort.
CREATE INDEX IF NOT EXISTS "idx_work_status_members"
    ON "work" (LOWER("status"), "members" DESC, "id");
