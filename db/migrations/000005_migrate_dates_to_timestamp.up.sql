-- Widen start_date, end_date and aired from DATE to a timestamp.
--
-- The MySQL original called FROM_UNIXTIME() on columns that 000001 declared as
-- DATE. MySQL coerces a DATE to a number for that call -- 2023-01-01 becomes
-- 20230101 -- and then reads it as a Unix epoch, which lands in 1970. So the
-- original converted these columns by producing nonsense and, two migrations
-- later, 000006 and 000007 rewrote them again.
--
-- Reproducing that bug faithfully would mean writing deliberately wrong data.
-- These columns are cast properly instead, which is what the migration's name
-- says it does. The end state is identical either way: 000006 turns them into
-- strings and 000007 re-derives aired, so nothing downstream reads what this
-- migration writes.
ALTER TABLE anime ADD start_date2 timestamptz NULL DEFAULT NULL;
UPDATE anime SET start_date2 = start_date::timestamptz;
ALTER TABLE anime DROP start_date;
ALTER TABLE anime ADD start_date timestamptz NULL DEFAULT NULL;
UPDATE anime SET start_date = start_date2;
ALTER TABLE anime DROP start_date2;

ALTER TABLE anime ADD end_date2 timestamptz NULL DEFAULT NULL;
UPDATE anime SET end_date2 = end_date::timestamptz;
ALTER TABLE anime DROP end_date;
ALTER TABLE anime ADD end_date timestamptz NULL DEFAULT NULL;
UPDATE anime SET end_date = end_date2;
ALTER TABLE anime DROP end_date2;

ALTER TABLE episodes ADD aired2 timestamptz NULL DEFAULT NULL;
UPDATE episodes SET aired2 = aired::timestamptz;
ALTER TABLE episodes DROP aired;
ALTER TABLE episodes ADD aired timestamptz NULL DEFAULT NULL;
UPDATE episodes SET aired = aired2;
ALTER TABLE episodes DROP aired2;
