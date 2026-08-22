-- This file was empty, in the MySQL original too, so `migrate down -all` never
-- removed episodes and left it behind holding a foreign key to anime. MySQL
-- tolerated that until the very last migration; Postgres refuses to drop anime
-- while the reference exists, which is what surfaced it.
--
-- Filled in rather than left as-is: an empty down is only invisible until
-- someone actually runs it.
DROP TABLE IF EXISTS episodes;
