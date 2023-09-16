CREATE TABLE IF NOT EXISTS buffs
(
    encounter INT           NOT NULL,
    player    STRING        NOT NULL,
    buff_id   INT           NOT NULL,
    damage    BIGINT        NOT NULL,
    percent   DECIMAL(2, 2) NOT NULL,
    UNIQUE (encounter, player, buff_id),
    FOREIGN KEY (encounter) REFERENCES encounters (id) ON DELETE CASCADE,
    INDEX (encounter)
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
    name         STRING NOT NULL,
    UNIQUE (encounter, player, skill_id),
    FOREIGN KEY (encounter) REFERENCES encounters (id) ON DELETE CASCADE,
    INDEX (encounter),
    INDEX (name),
    INDEX (dps)
);