ALTER TABLE anime_character_staff_link
    ADD COLUMN character_name VARCHAR(255) NOT NULL,
    ADD COLUMN staff_given_name VARCHAR(255) NOT NULL,
    ADD COLUMN staff_family_name VARCHAR(255) NOT NULL;
