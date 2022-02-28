CREATE TABLE sessions
(
    id         UUID,
    user_id     UUID,
    expiration timestamp,
    PRIMARY KEY (user_id, id)
);