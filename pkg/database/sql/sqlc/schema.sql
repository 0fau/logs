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

    boss         STRING    NOT NULL,
    date         TIMESTAMP NOT NULL,
    duration     INT       NOT NULL,
    local_player STRING    NOT NULL
);

CREATE TABLE players
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