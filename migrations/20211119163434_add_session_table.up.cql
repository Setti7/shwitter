CREATE TABLE sessions
(
    id         text,
    userId     UUID,
    expiration timestamp,
    PRIMARY KEY (id)
);