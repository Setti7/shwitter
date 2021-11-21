CREATE TABLE followers
(
    user_id      UUID,
    follower_id UUID,
    since       TIMESTAMP,
    PRIMARY KEY ( user_id, follower_id )
);

CREATE TABLE friends
(
    user_id    UUID,
    friend_id UUID,
    since     TIMESTAMP,
    PRIMARY KEY ( user_id, friend_id )
);