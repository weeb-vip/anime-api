CREATE TABLE episodes
(
    id         varchar(36) PRIMARY KEY,
    anime_id   varchar(36) REFERENCES anime (id),
    episode    int,
    title_en   text,
    title_jp   text,
    aired      timestamptz,
    synopsis   text,
    created_at timestamptz,
    updated_at timestamptz
);
