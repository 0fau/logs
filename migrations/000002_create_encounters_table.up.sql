CREATE SEQUENCE IF NOT EXISTS encounter_id_seq START 1 INCREMENT 1;

CREATE TABLE IF NOT EXISTS encounters
(
    id             INT                DEFAULT nextval('encounter_id_seq'),
    uploaded_by    UUID      NOT NULL REFERENCES users (id),
    uploaded_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    settings       JSONB     NOT NULL,
    tags           STRING ARRAY,

    header         JSONB     NOT NULL,
    data           JSON      NOT NULL,

    boss           STRING    NOT NULL,
    difficulty     STRING    NOT NULL,
    date           TIMESTAMP NOT NULL,
    duration       INT       NOT NULL,
    local_player   STRING    NOT NULL,

    unique_hash    STRING    NOT NULL,
    unique_group   INT       NOT NULL,

    PRIMARY KEY (id) USING HASH,
    UNIQUE (date, local_player, boss),
    INDEX (uploaded_by),
    INDEX (date),
    INDEX (boss),
    INDEX (unique_hash)
);

CREATE TABLE IF NOT EXISTS players
(
    encounter INT    NOT NULL REFERENCES encounters (id) ON DELETE CASCADE,
    class     STRING NOT NULL,
    name      STRING NOT NULL,
    dead      BOOL   NOT NULL,
    data      JSONB  NOT NULL,
    place     INT    NOT NULL,
    PRIMARY KEY (encounter, name),
    INDEX (name),
    INDEX (class),
    INDEX (((data ->> 'dps')::BIGINT), place)
);