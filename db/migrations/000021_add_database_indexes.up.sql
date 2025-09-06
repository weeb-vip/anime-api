-- Add database indexes for improved query performance

-- Helper procedure to safely create indexes
DELIMITER //
CREATE PROCEDURE CreateIndexIfNotExists(
    IN indexName VARCHAR(255),
    IN tableName VARCHAR(255), 
    IN columnDef VARCHAR(255)
)
BEGIN
    DECLARE index_exists INT DEFAULT 0;
    
    SELECT COUNT(*) INTO index_exists
    FROM information_schema.statistics 
    WHERE table_schema = DATABASE() 
    AND table_name = tableName 
    AND index_name = indexName;
    
    IF index_exists = 0 THEN
        SET @sql = CONCAT('CREATE INDEX ', indexName, ' ON ', tableName, ' ', columnDef);
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END IF;
END//
DELIMITER ;

-- Anime table indexes
CALL CreateIndexIfNotExists('idx_anime_status', 'anime', '(status)');
CALL CreateIndexIfNotExists('idx_anime_rating', 'anime', '(rating)');
CALL CreateIndexIfNotExists('idx_anime_ranking', 'anime', '(ranking)');
CALL CreateIndexIfNotExists('idx_anime_created_at', 'anime', '(created_at)');
CALL CreateIndexIfNotExists('idx_anime_type', 'anime', '(type)');
CALL CreateIndexIfNotExists('idx_anime_source', 'anime', '(source)');
CALL CreateIndexIfNotExists('idx_anime_start_date', 'anime', '(start_date)');
CALL CreateIndexIfNotExists('idx_anime_end_date', 'anime', '(end_date)');
CALL CreateIndexIfNotExists('idx_anime_anidbid', 'anime', '(anidbid)');
CALL CreateIndexIfNotExists('idx_anime_thetvdbid', 'anime', '(thetvdbid)');

-- Full-text search indexes for title fields
CALL CreateIndexIfNotExists('idx_anime_title_en', 'anime', '(title_en(255))');
CALL CreateIndexIfNotExists('idx_anime_title_jp', 'anime', '(title_jp(255))');
CALL CreateIndexIfNotExists('idx_anime_title_romaji', 'anime', '(title_romaji(255))');
CALL CreateIndexIfNotExists('idx_anime_title_kanji', 'anime', '(title_kanji(255))');

-- Episodes table indexes
CALL CreateIndexIfNotExists('idx_episodes_anime_id', 'episodes', '(anime_id)');
CALL CreateIndexIfNotExists('idx_episodes_aired', 'episodes', '(aired)');
CALL CreateIndexIfNotExists('idx_episodes_episode_number', 'episodes', '(episode)');
CALL CreateIndexIfNotExists('idx_episodes_created_at', 'episodes', '(created_at)');

-- Composite index for episodes by anime and air date
CALL CreateIndexIfNotExists('idx_episodes_anime_aired', 'episodes', '(anime_id, aired)');

-- Anime characters table indexes
CALL CreateIndexIfNotExists('idx_anime_character_anime_id', 'anime_character', '(anime_id)');
CALL CreateIndexIfNotExists('idx_anime_character_name', 'anime_character', '(name)');
CALL CreateIndexIfNotExists('idx_anime_character_role', 'anime_character', '(role)');

-- Anime staff table indexes
CALL CreateIndexIfNotExists('idx_anime_staff_given_name', 'anime_staff', '(given_name)');
CALL CreateIndexIfNotExists('idx_anime_staff_family_name', 'anime_staff', '(family_name)');
CALL CreateIndexIfNotExists('idx_anime_staff_language', 'anime_staff', '(language)');

-- Anime character staff link table indexes (if exists)
CALL CreateIndexIfNotExists('idx_character_staff_character_id', 'anime_character_staff_link', '(character_id)');
CALL CreateIndexIfNotExists('idx_character_staff_staff_id', 'anime_character_staff_link', '(staff_id)');

-- Anime seasons table indexes (already has some indexes but adding missing ones)
CALL CreateIndexIfNotExists('idx_anime_seasons_episode_count', 'anime_seasons', '(episode_count)');
CALL CreateIndexIfNotExists('idx_anime_seasons_created_at', 'anime_seasons', '(created_at)');

-- Relations table indexes (if exists)
CALL CreateIndexIfNotExists('idx_relations_anime_id', 'relations', '(anime_id)');
CALL CreateIndexIfNotExists('idx_relations_related_anime_id', 'relations', '(related_anime_id)');

-- Drop the helper procedure
DROP PROCEDURE CreateIndexIfNotExists;