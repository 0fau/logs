CREATE TABLE IF NOT EXISTS friends
(
    user1 UUID      NOT NULL REFERENCES users (id),
    user2 UUID      NOT NULL REFERENCES users (id),
    date  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user1, user2),
    INDEX (user1),
    INDEX (user2)
);

CREATE TABLE IF NOT EXISTS friend_requests
(
    user1 string    NOT NULL REFERENCES users (discord_id),
    user2 string    NOT NULL REFERENCES users (discord_id),
    date  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user1, user2),
    INDEX (user1),
    INDEX (user2)
);