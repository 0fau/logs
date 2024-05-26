CREATE TABLE IF NOT EXISTS api_tokens
(
    user_id UUID   NOT NULL REFERENCES users (id),
    token   STRING NOT NULL PRIMARY KEY,
    grade   INT    NOT NULL
);