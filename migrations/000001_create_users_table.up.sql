CREATE TABLE IF NOT EXISTS users
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discord_id   STRING,
    discord_name STRING(32) NOT NULL,
    access_token STRING(64),
    roles        STRING ARRAY,
    created_at   TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP        DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE (discord_id)
);