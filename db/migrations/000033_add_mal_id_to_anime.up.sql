ALTER TABLE anime
    ADD COLUMN mal_id INT NULL AFTER anidbid,
    ADD INDEX idx_anime_mal_id (mal_id);
