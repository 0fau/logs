CREATE TABLE users
(
    id             UUID PRIMARY KEY,
    username       STRING,
    created_at     TIMESTAMP NOT NULL,
    updated_at     TIMESTAMP NOT NULL,
    access_token   STRING(64),

    discord_id     STRING    NOT NULL,
    discord_tag    STRING    NOT NULL,
    avatar         STRING    NOT NULL,

    friends        UUID ARRAY,
    settings       JSONB     NOT NULL,
    log_visibility JSONB,

    titles         STRING[],
    roles          STRING[]
);

CREATE TABLE encounters
(
    id           INT PRIMARY KEY,
    uploaded_by  UUID      NOT NULL REFERENCES users (id),
    uploaded_at  TIMESTAMP NOT NULL,
    settings     JSONB     NOT NULL,
    thumbnail    BOOLEAN   NOT NULL,
    tags         STRING[],

    header       JSONB     NOT NULL,
    data         JSON      NOT NULL,

    private      BOOLEAN   NOT NULL,

    unique_hash  STRING    NOT NULL,
    unique_group INT       NOT NULL,

    visibility   JSONB,

    boss         STRING    NOT NULL,
    difficulty   STRING    NOT NULL,
    date         TIMESTAMP NOT NULL,
    duration     INT       NOT NULL,
    version      INT       NOT NULL,
    local_player STRING    NOT NULL
);

CREATE TABLE players
(
    encounter  INT    NOT NULL REFERENCES encounters (id),
    boss       STRING NOT NULL,
    difficulty STRING NOT NULL,
    class      STRING NOT NULL,
    name       STRING NOT NULL,
    dead       BOOL   NOT NULL,
    dps        BIGINT NOT NULL,
    gear_score FLOAT  NOT NULL,
    place      INT    NOT NULL
);

CREATE TABLE grouped_encounters
(
    group_id  INT NOT NULL PRIMARY KEY REFERENCES encounters (id),
    uploaders UUID ARRAY
);

CREATE TABLE whitelist
(
    discord STRING NOT NULL,
    role    STRING NOT NULL
);

CREATE TABLE IF NOT EXISTS friends
(
    user1 UUID      NOT NULL REFERENCES users (id),
    user2 UUID      NOT NULL REFERENCES users (id),
    date  TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS friend_requests
(
    user1 string    NOT NULL REFERENCES users (discord_id),
    user2 string    NOT NULL REFERENCES users (discord_id),
    date  TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS roster
(
    user_id    UUID   NOT NULL REFERENCES users (id),
    character  STRING NOT NULL,
    class      STRING NOT NULL,
    gear_score FLOAT  NOT NULL
);

CREATE TABLE IF NOT EXISTS raids
(
    id         INT       NOT NULL PRIMARY KEY,
    boss       STRING    NOT NULL,
    difficulty STRING    NOT NULL,
    date       TIMESTAMP NOT NULL,
    duration   BIGINT    NOT NULL,
    uploaders  UUID[]    NOT NULL,
    players    STRING[]  NOT NULL
);