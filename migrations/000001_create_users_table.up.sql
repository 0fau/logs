CREATE TABLE IF NOT EXISTS users
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discord_id     STRING NOT NULL,
    discord_tag    STRING NOT NULL,
    avatar         STRING NOT NULL,
    username       STRING,
    access_token   STRING(64),

    friends        UUID ARRAY,
    settings       JSONB  NOT NULL,
    log_visibility JSONB,

    titles         STRING ARRAY,
    roles          STRING ARRAY,

    created_at     TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP        DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE (discord_id),
    UNIQUE (UPPER(username))
);