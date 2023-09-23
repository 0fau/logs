CREATE TABLE IF NOT EXISTS skills
(
    encounter INT    NOT NULL,
    player    STRING NOT NULL,
    skill_id  INT    NOT NULL,
    name      STRING NOT NULL,
    dps       BIGINT NOT NULL,
    damage    BIGINT NOT NULL,
    tripods   JSONB  NOT NULL,
    fields    JSON   NOT NULL,
    UNIQUE (encounter, player, skill_id),
    FOREIGN KEY (encounter) REFERENCES encounters (id) ON DELETE CASCADE,
    INDEX (encounter),
    INDEX (name),
    INDEX (dps)
);