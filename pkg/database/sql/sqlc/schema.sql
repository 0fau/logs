CREATE TABLE users
(
    id           UUID,
    discord_id   TEXT,
    discord_name STRING(32),
    access_token STRING(64),
    roles        STRING[],
    created_at   TIMESTAMP,
    updated_at   TIMESTAMP
);