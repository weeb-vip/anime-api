-- Rollback composite indexes

-- Drop the composite indexes
DROP INDEX IF EXISTS idx_anime_created_at_id;
DROP INDEX IF EXISTS idx_anime_ranking_id;