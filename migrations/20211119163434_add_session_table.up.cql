CREATE TABLE sessions
(
    sess_token text, // TODO: change to varchar?
    userId     UUID,
    expiration timestamp,
    PRIMARY KEY (sess_token)
);