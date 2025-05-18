CREATE TABLE anime_relations (
                                 id UUID DEFAULT (UUID()) PRIMARY KEY,
                                 anime_id VARCHAR(36) NOT NULL,
                                 related_anime_id VARCHAR(36) NOT NULL,
                                 relation_type VARCHAR(30), -- e.g., 'sequel', 'prequel', 'special', 'same_series'
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
