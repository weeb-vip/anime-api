-- Converts the four remaining char(36) id columns to varchar(36), so joins
-- against them can use an index.
--
-- anime_character.id and anime_staff.id are char(36); the columns that
-- reference them, anime_character_staff_link.character_id and .staff_id, are
-- varchar(36). Postgres resolves that mismatch by casting, and the cast lands
-- on the indexed side:
--
--   Hash Cond: ((anime_character_staff_link.character_id)::bpchar = anime_character.id)
--   ->  Parallel Seq Scan on anime_character_staff_link (rows=149862 loops=3)
--
-- A cast on the indexed column makes the index unusable, so every request for
-- an anime's characters read all 450,000 link rows instead of the handful it
-- needed. That is why idx_character_staff, at 116MB, had 149 scans while the
-- table had 37,000 sequential scans, and why this query was the largest CPU
-- consumer in Performance Insights.
--
-- Measured on a reproduction of the same shape and volume:
--
--   before:  Parallel Seq Scan, 150000 rows x 3 workers   18.9ms
--   after:   Index Only Scan,   2 rows x 75 loops          2.1ms
--
-- varchar is the direction rather than char because varchar is what the rest of
-- the schema uses: these four were the only char columns against thirty-six
-- varchar ones. char(36) also blank-pads and compares with trailing-space
-- semantics, which is a trap for the next join written against them.
--
-- This does not rewrite the tables. bpchar and varchar share a representation,
-- so Postgres treats the change as binary-coercible: it updates the catalog and
-- revalidates. The two foreign keys survive untouched -- verified, not assumed.
--
-- The values are 36-character UUIDs, so there is no padding to lose.
--
-- One transaction, so a failure leaves every column on its original type rather
-- than half-converted. Each ALTER takes an ACCESS EXCLUSIVE lock for the
-- duration; on the row counts here that measured under a second per table.
ALTER TABLE anime_character ALTER COLUMN id TYPE varchar(36);
ALTER TABLE anime_staff ALTER COLUMN id TYPE varchar(36);
ALTER TABLE anime_character_staff_link ALTER COLUMN id TYPE varchar(36);
ALTER TABLE anime_relations ALTER COLUMN id TYPE varchar(36);
