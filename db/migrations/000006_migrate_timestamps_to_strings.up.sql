ALTER TABLE anime ADD COLUMN formatted_timestamp VARCHAR(255);
UPDATE anime SET formatted_timestamp = to_char(start_date, 'YYYY-MM-DD HH24:MI:SS') WHERE start_date IS NOT NULL;
ALTER TABLE anime DROP COLUMN start_date;
ALTER TABLE anime ADD COLUMN start_date VARCHAR(255);
UPDATE anime SET start_date = formatted_timestamp WHERE formatted_timestamp IS NOT NULL;
ALTER TABLE anime DROP COLUMN formatted_timestamp;

ALTER TABLE anime ADD COLUMN formatted_timestamp VARCHAR(255);
UPDATE anime SET formatted_timestamp = to_char(end_date, 'YYYY-MM-DD HH24:MI:SS') WHERE end_date IS NOT NULL;
ALTER TABLE anime DROP COLUMN end_date;
ALTER TABLE anime ADD COLUMN end_date VARCHAR(255);
UPDATE anime SET end_date = formatted_timestamp WHERE formatted_timestamp IS NOT NULL;
ALTER TABLE anime DROP COLUMN formatted_timestamp;

ALTER TABLE episodes ADD COLUMN formatted_timestamp VARCHAR(255);
UPDATE episodes SET formatted_timestamp = to_char(aired, 'YYYY-MM-DD HH24:MI:SS') WHERE aired IS NOT NULL;
ALTER TABLE episodes DROP COLUMN aired;
ALTER TABLE episodes ADD COLUMN aired VARCHAR(255);
UPDATE episodes SET aired = formatted_timestamp WHERE formatted_timestamp IS NOT NULL;
ALTER TABLE episodes DROP COLUMN formatted_timestamp;
