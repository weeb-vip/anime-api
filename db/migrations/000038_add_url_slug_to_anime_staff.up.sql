-- url_slug is a voice actor's public URL segment: /people/mary-elizabeth-mcglynn
-- instead of /people/<uuid>.
--
-- Unlike anime.url_slug, this one is NOT minted in postgres and carried over
-- CDC. It does not need to be. A staff slug is a pure function of given_name
-- and family_name, and those are already unique: the scraper deduplicates
-- anime_staff on exactly that pair, so one row is one person and the slug needs
-- no minting authority to keep it distinct. Measured over 21,545 staff rows,
-- the expression below yields 21,541 distinct slugs; all four collisions are
-- the same person recorded twice under spelling variants the scraper's
-- exact-match dedup missed (Daniel García / Daniel Garcia, Julien Haggége /
-- Haggège), not two different people.
--
-- Deriving it here rather than upstream is what keeps this change to one
-- repository. A plain column would have needed something to populate rows
-- arriving over CDC -- a scheduled backfill, and a window in which a newly
-- synced voice actor has no slug and 404s. A STORED generated column has no
-- such window: MySQL computes it on insert and recomputes it on update, so
-- rows the sync writes are correct the moment they land, and a name change
-- moves the slug with it. character-staff-sync needs no change at all; GORM
-- builds its column list from a struct that has no url_slug, so its upserts
-- never name this column, which is the one thing that would make them fail.
--
-- The REPLACE chain folds accents before the character class strips what is
-- left. Without it 8.5% of names (1,835 of 21,545) slug badly -- Danièle Hazan
-- became "dani-le-hazan" -- because a bare accented character is not in
-- [a-z0-9] and became a separator. LOWER runs first and handles accented
-- uppercase on its own, so only lowercase forms need mapping. The set covers
-- every non-ASCII character present in the data (é Á ó í ú ç ñ ß ş ÿ and a
-- combining tilde) plus the rest of the common Latin range, since new names
-- arrive continuously.
--
-- The index is deliberately NOT unique, for the same reason 000037 gave for
-- anime: this table is a CDC replica, and a constraint here could only reject a
-- write the source already accepted -- stalling the consumer. The four
-- duplicate-person rows above would do exactly that.

ALTER TABLE anime_staff
    ADD COLUMN url_slug VARCHAR(255)
        GENERATED ALWAYS AS (
            NULLIF(
                TRIM(BOTH '-' FROM
                    REGEXP_REPLACE(
                        REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(LOWER(CONCAT(given_name,' ',family_name)),'á','a'),'à','a'),'â','a'),'ä','a'),'ã','a'),'å','a'),'ā','a'),'ă','a'),'ą','a'),'é','e'),'è','e'),'ê','e'),'ë','e'),'ē','e'),'ĕ','e'),'ė','e'),'ę','e'),'ě','e'),'í','i'),'ì','i'),'î','i'),'ï','i'),'ī','i'),'į','i'),'ı','i'),'ó','o'),'ò','o'),'ô','o'),'ö','o'),'õ','o'),'ō','o'),'ø','o'),'ő','o'),'ú','u'),'ù','u'),'û','u'),'ü','u'),'ū','u'),'ů','u'),'ű','u'),'ų','u'),'ç','c'),'ć','c'),'č','c'),'ĉ','c'),'ñ','n'),'ń','n'),'ň','n'),'š','s'),'ś','s'),'ş','s'),'ș','s'),'ž','z'),'ź','z'),'ż','z'),'ý','y'),'ÿ','y'),'ğ','g'),'ĝ','g'),'ł','l'),'đ','d'),'ď','d'),'ť','t'),'ț','t'),'ř','r'),'ß','ss'),'æ','ae'),'œ','oe'),'̀',''),'́',''),'̂',''),'̃',''),'̈',''),'̊',''),'̧',''),
                        '[^a-z0-9]+', '-'
                    )
                ),
                ''
            )
        ) STORED;

CREATE INDEX idx_anime_staff_url_slug ON anime_staff (url_slug);
