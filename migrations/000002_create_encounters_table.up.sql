CREATE SEQUENCE IF NOT EXISTS encounter_id_seq START 1 INCREMENT 1;

CREATE TABLE IF NOT EXISTS encounters
(
    id           INT                DEFAULT nextval('encounter_id_seq'),
    uploaded_by  UUID      NOT NULL REFERENCES users (id),
    title        STRING,
    description  STRING,
    visibility   STRING    NOT NULL,
    raid         STRING    NOT NULL,
    date         TIMESTAMP NOT NULL,
    duration     INT       NOT NULL,
    damage       BIGINT    NOT NULL,
    fields       JSON      NOT NULL,
    cleared      BOOLEAN   NOT NULL,
    uploaded_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    tags         STRING ARRAY,
    local_player STRING    NOT NULL,
    PRIMARY KEY (id) USING HASH,
    INDEX (tags),
    INDEX (uploaded_by),
    INDEX (date),
    INDEX (raid),
    INDEX (local_player)
);

CREATE TABLE IF NOT EXISTS entities
(
    encounter INT    NOT NULL REFERENCES encounters (id) ON DELETE CASCADE,
    class     STRING NOT NULL,
    name      STRING NOT NULL,
    enttype   STRING NOT NULL,
    damage    BIGINT NOT NULL,
    dps       BIGINT NOT NULL,
    dead      BOOL   NOT NULL,
    tag       STRING ARRAY,
    fields    JSON   NOT NULL,
    UNIQUE (encounter, enttype, name),
    INDEX (encounter, name),
    INDEX (name),
    INDEX (class),
    INDEX (dps)
);