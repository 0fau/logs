CREATE TABLE IF NOT EXISTS whitelist
(
    discord STRING NOT NULL,
    role    STRING NOT NULL,
    UNIQUE (discord)
);