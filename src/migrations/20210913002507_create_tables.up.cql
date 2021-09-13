CREATE TABLE shweets
(
    id      UUID,
    user_id UUID,
    message text,
    PRIMARY KEY ( id )
);

CREATE TABLE users
(
    id    UUID,
    name  text,
    email text,
    bio   text,
    PRIMARY KEY (id)
);