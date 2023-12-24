CREATE TABLE IF NOT EXISTS grouped_encounters
(
    group_id  INT NOT NULL PRIMARY KEY REFERENCES encounters (id),
    uploaders UUID ARRAY
);