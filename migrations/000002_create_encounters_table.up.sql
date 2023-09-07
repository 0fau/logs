CREATE SEQUENCE IF NOT EXISTS encounter_id_seq START 1 INCREMENT 1;

CREATE TABLE IF NOT EXISTS encounters
(
    id                 INT                DEFAULT nextval('encounter_id_seq'),
    uploaded_by        UUID      NOT NULL REFERENCES users (id),
    title              STRING,
    description        STRING,
    visibility         STRING    NOT NULL,
    raid               STRING    NOT NULL,
    date               TIMESTAMP NOT NULL,
    duration           INT       NOT NULL,
    total_damage_dealt BIGINT    NOT NULL,
    cleared            BOOLEAN   NOT NULL,
    uploaded_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id) USING HASH,
    INDEX (uploaded_by),
    INDEX (date),
    INDEX (raid)
);

CREATE TABLE IF NOT EXISTS entities
(
    encounter INTEGER NOT NULL REFERENCES encounters (id),
    class     STRING  NOT NULL,
    name      STRING  NOT NULL,
    enttype   STRING  NOT NULL,
    damage    BIGINT  NOT NULL,
    dps       INT     NOT NULL,
    PRIMARY KEY (encounter, enttype, name) USING HASH,
    INDEX (name),
    INDEX (class),
    INDEX (dps)
);