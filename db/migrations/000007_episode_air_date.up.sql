-- Rename the original 'aired' column
ALTER TABLE episodes RENAME COLUMN aired TO backup_aired;

-- Add the new 'aired' column as DATE
ALTER TABLE episodes ADD COLUMN aired DATE;

-- Set 'aired' only where 'backup_aired' is not null
UPDATE episodes
SET aired = STR_TO_DATE(backup_aired, '%Y-%m-%d %H:%i:%s')
WHERE backup_aired IS NOT NULL;