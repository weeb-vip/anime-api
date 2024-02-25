ALTER TABLE `anime` ADD COLUMN formatted_timestamp VARCHAR(255);
UPDATE `anime` SET formatted_timestamp = DATE_FORMAT(start_date, '%Y-%m-%d %H:%i:%s') WHERE start_date IS NOT NULL;
ALTER TABLE `anime` DROP COLUMN start_date;
ALTER TABLE `anime` ADD COLUMN start_date VARCHAR(255);
UPDATE `anime` SET start_date = formatted_timestamp WHERE formatted_timestamp IS NOT NULL;
ALTER TABLE `anime` DROP COLUMN formatted_timestamp;

ALTER TABLE `anime` ADD COLUMN formatted_timestamp VARCHAR(255);
UPDATE `anime` SET formatted_timestamp = DATE_FORMAT(end_date, '%Y-%m-%d %H:%i:%s') WHERE end_date IS NOT NULL
ALTER TABLE `anime` DROP COLUMN end_date;
ALTER TABLE `anime` ADD COLUMN end_date VARCHAR(255);
UPDATE `anime` SET end_date = formatted_timestamp WHERE formatted_timestamp IS NOT NULL;
ALTER TABLE `anime` DROP COLUMN formatted_timestamp;

ALTER TABLE `anime` ADD COLUMN formatted_timestamp VARCHAR(255);
UPDATE `anime` SET formatted_timestamp = DATE_FORMAT(airing_start, '%Y-%m-%d %H:%i:%s') WHERE airing_start IS NOT NULL;
ALTER TABLE `anime` DROP COLUMN airing_start;
ALTER TABLE `anime` ADD COLUMN airing_start VARCHAR(255);
UPDATE `anime` SET airing_start = formatted_timestamp WHERE formatted_timestamp IS NOT NULL;
ALTER TABLE `anime` DROP COLUMN formatted_timestamp;
