CREATE TABLE sessions
(
    id         text,
    user_id     UUID,
    expiration timestamp,
    PRIMARY KEY (user_id, id)
);