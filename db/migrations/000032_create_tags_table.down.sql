-- Drop anime_tags junction table
DROP TABLE IF EXISTS anime_tags;

-- Drop tags table. Its trigger goes with it.
DROP TABLE IF EXISTS tags;

-- set_updated_at() is deliberately NOT dropped here. It is created with CREATE
-- OR REPLACE by every migration whose table needs it -- anime_character,
-- anime_staff, anime_schedule and others, all of which predate this one -- so
-- dropping it here would either fail on the dependency or, with CASCADE, take
-- their triggers with it and leave those tables silently not updating
-- updated_at.
