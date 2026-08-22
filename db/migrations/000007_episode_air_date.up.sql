ALTER TABLE episodes RENAME COLUMN aired TO backup_aired;
ALTER TABLE episodes ADD COLUMN aired DATE;
UPDATE episodes SET aired = to_timestamp(backup_aired, 'YYYY-MM-DD HH24:MI:SS') WHERE backup_aired IS NOT NULL;
