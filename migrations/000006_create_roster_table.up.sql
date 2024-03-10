CREATE TABLE IF NOT EXISTS roster
(
    user_id    UUID   NOT NULL REFERENCES users (id),
    character  STRING NOT NULL,
    class      STRING NOT NULL,
    gear_score FLOAT  NOT NULL,
    PRIMARY KEY (user_id, character)
);