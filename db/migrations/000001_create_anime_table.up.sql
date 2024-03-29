CREATE TABLE anime
(
    id             varchar(36) PRIMARY KEY,
    type           varchar(30) DEFAULT 'Anime',
    title_en       text,
    title_jp       text,
    title_romaji   text,
    title_kanji    text,
    title_synonyms text,
    image_url      text,
    synopsis       text,
    episodes       int,
    status         text,
    start_date     date,
    end_date       date,
    genres         text,
    duration       text,
    broadcast      text,
    source         text,
    licensors      text,
    studios        text,
    rating         text,
    created_at     timestamp,
    updated_at     timestamp
);
