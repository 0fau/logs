CREATE TABLE users
(
    id           UUID PRIMARY KEY,
    discord_id   TEXT       NOT NULL,
    discord_name STRING(32) NOT NULL,
    access_token STRING(64),
    roles        STRING[],
    created_at   TIMESTAMP  NOT NULL,
    updated_at   TIMESTAMP  NOT NULL
);

CREATE TABLE encounters
(
    id                 INT PRIMARY KEY,
    uploaded_by        UUID      NOT NULL REFERENCES users (id),
    visibility         STRING    NOT NULL,
    title              STRING,
    description        STRING,
    raid               STRING    NOT NULL,
    date               TIMESTAMP NOT NULL,
    duration           INT       NOT NULL,
    total_damage_dealt BIGINT    NOT NULL,
    cleared            BOOLEAN   NOT NULL,
    uploaded_at        TIMESTAMP NOT NULL,
    tags               STRING[],
    local_player       STRING    NOT NULL
);

CREATE TABLE entities
(
    encounter INTEGER NOT NULL REFERENCES encounters (id),
    class     STRING  NOT NULL,
    enttype   STRING  NOT NULL,
    name      STRING  NOT NULL,
    damage    BIGINT  NOT NULL,
    dps       BIGINT  NOT NULL
);

CREATE TABLE IF NOT EXISTS buffs
(
    encounter INT           NOT NULL,
    player    STRING        NOT NULL,
    buff_id   INT           NOT NULL,
    percent   DECIMAL(2, 2) NOT NULL,
    damage    BIGINT        NOT NULL
);

CREATE TABLE IF NOT EXISTS skills
(
    encounter    INT    NOT NULL,
    player       STRING NOT NULL,
    skill_id     INT    NOT NULL,
    casts        INT    NOT NULL,
    crits        INT    NOT NULL,
    dps          BIGINT NOT NULL,
    hits         INT    NOT NULL,
    max_damage   BIGINT NOT NULL,
    total_damage BIGINT NOT NULL,
    name         STRING NOT NULL
);