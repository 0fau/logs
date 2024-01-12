CREATE TABLE users
(
    id           UUID PRIMARY KEY,
    username     STRING,
    created_at   TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP NOT NULL,
    access_token STRING(64),

    discord_id   STRING    NOT NULL,
    discord_tag  STRING    NOT NULL,
    avatar       STRING,

    friends      UUID ARRAY,
    settings     JSONB     NOT NULL,

    titles       STRING[],
    roles        STRING[]
);

CREATE TABLE encounters
(
    id           INT PRIMARY KEY,
    uploaded_by  UUID      NOT NULL REFERENCES users (id),
    uploaded_at  TIMESTAMP NOT NULL,
    settings     JSONB     NOT NULL,
    tags         STRING[],

    header       JSONB     NOT NULL,
    data         JSON      NOT NULL,

    unique_hash  STRING    NOT NULL,
    unique_group INT       NOT NULL,

    boss         STRING    NOT NULL,
    difficulty   STRING    NOT NULL,
    date         TIMESTAMP NOT NULL,
    duration     INT       NOT NULL,
    local_player STRING    NOT NULL
);

CREATE TABLE players
(
    encounter INT    NOT NULL REFERENCES encounters (id),
    class     STRING NOT NULL,
    name      STRING NOT NULL,
    dead      BOOL   NOT NULL,
    data      JSONB  NOT NULL,
    dps       BIGINT NOT NULL,
    place     INT    NOT NULL
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