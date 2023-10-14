CREATE SEQUENCE IF NOT EXISTS encounter_id_seq START 1 INCREMENT 1;

CREATE TABLE IF NOT EXISTS encounters
(
    id           INT                DEFAULT nextval('encounter_id_seq'),
    uploaded_by  UUID      NOT NULL REFERENCES users (id),
    uploaded_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    settings     JSONB     NOT NULL,
    tags         STRING ARRAY,

    header       JSONB     NOT NULL,
    data         JSON      NOT NULL,

    boss         STRING    NOT NULL,
    date         TIMESTAMP NOT NULL,
    duration     INT       NOT NULL,
    local_player STRING    NOT NULL,

    PRIMARY KEY (id) USING HASH,
    INDEX (uploaded_by),
    INDEX (date),
    INDEX (boss),
    INDEX (local_player)
);

CREATE TABLE IF NOT EXISTS players
(
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter INT    NOT NULL REFERENCES encounters (id) ON DELETE CASCADE,
    class     STRING NOT NULL,
    name      STRING NOT NULL,
    damage    BIGINT NOT NULL,
    dps       BIGINT NOT NULL,
    dead      BOOL   NOT NULL,
    tags      STRING ARRAY,
    data      JSON   NOT NULL,
    UNIQUE (encounter, name),
    INDEX (name),
    INDEX (class),
    INDEX (dps)
);