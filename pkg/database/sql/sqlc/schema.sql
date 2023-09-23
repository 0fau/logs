CREATE TABLE users
(
    id           UUID PRIMARY KEY,
    discord_id   TEXT      NOT NULL,
    discord_name STRING(32) NOT NULL,
    access_token STRING(64),
    roles        STRING[],
    created_at   TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP NOT NULL
);

CREATE TABLE encounters
(
    id           INT PRIMARY KEY,
    uploaded_by  UUID      NOT NULL REFERENCES users (id),
    visibility   STRING    NOT NULL,
    title        STRING,
    description  STRING,
    raid         STRING    NOT NULL,
    date         TIMESTAMP NOT NULL,
    duration     INT       NOT NULL,
    damage       BIGINT    NOT NULL,
    fields       JSON      NOT NULL,
    cleared      BOOLEAN   NOT NULL,
    uploaded_at  TIMESTAMP NOT NULL,
    tags         STRING[],
    local_player STRING    NOT NULL
);

CREATE TABLE entities
(
    encounter INTEGER NOT NULL REFERENCES encounters (id),
    enttype   STRING  NOT NULL,
    name      STRING  NOT NULL,
    class     STRING  NOT NULL,
    damage    BIGINT  NOT NULL,
    dps       BIGINT  NOT NULL,
    dead      BOOL    NOT NULL,
    fields    JSON    NOT NULL
);

CREATE TABLE IF NOT EXISTS skills
(
    encounter INT    NOT NULL,
    player    STRING NOT NULL,
    skill_id  INT    NOT NULL,
    name      STRING NOT NULL,
    dps       BIGINT NOT NULL,
    damage    BIGINT NOT NULL,
    tripods   JSONB  NOT NULL,
    fields    JSON   NOT NULL
);