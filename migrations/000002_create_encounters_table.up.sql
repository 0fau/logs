CREATE SEQUENCE IF NOT EXISTS encounter_id_seq START 1 INCREMENT 1;

CREATE TABLE IF NOT EXISTS encounters
(
    id           INT                DEFAULT nextval('encounter_id_seq'),
    uploaded_by  UUID      NOT NULL REFERENCES users (id),
    uploaded_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    settings     JSONB     NOT NULL,
    thumbnail    BOOLEAN   NOT NULL,
    tags         STRING ARRAY,

    header       JSONB     NOT NULL,
    data         JSON      NOT NULL,

    private      BOOLEAN   NOT NULL,

    boss         STRING    NOT NULL,
    difficulty   STRING    NOT NULL,
    date         TIMESTAMP NOT NULL,
    duration     INT       NOT NULL,
    version      INT       NOT NULL DEFAULT 0,
    local_player STRING    NOT NULL,

    visibility   JSONB,

    unique_hash  STRING    NOT NULL,
    unique_group INT       NOT NULL,

    PRIMARY KEY (id) USING HASH,
    UNIQUE (date, local_player, boss),
    INDEX (uploaded_by),
    INDEX (date),
    INDEX (boss),
    INDEX (version),
    INDEX (unique_hash)
);

CREATE TABLE IF NOT EXISTS players
(
    encounter  INT    NOT NULL REFERENCES encounters (id) ON DELETE CASCADE,
    class      STRING NOT NULL,
    name       STRING NOT NULL,
    dead       BOOL   NOT NULL,
    dps        BIGINT NOT NULL,
    place      INT    NOT NULL,
    gear_score FLOAT  NOT NULL,

    boss       STRING NOT NULL,
    difficulty STRING NOT NULL,

    PRIMARY KEY (encounter, name),
    INDEX (name),
    INDEX (class),
    INDEX (gear_score),
    INDEX (dps)
);