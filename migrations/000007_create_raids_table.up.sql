CREATE TABLE IF NOT EXISTS raids
(
    id         INT       NOT NULL PRIMARY KEY,
    boss       STRING    NOT NULL,
    difficulty STRING    NOT NULL,
    date       TIMESTAMP NOT NULL,
    duration   BIGINT    NOT NULL,
    uploaders  UUID[]    NOT NULL,
    players    STRING[]  NOT NULL
);