CREATE TABLE followers
(
    userid      UUID,
    follower_id UUID,
    since       TIMESTAMP,
    PRIMARY KEY ( userid, follower_id )
);

CREATE TABLE friends
(
    userid    UUID,
    friend_id UUID,
    since     TIMESTAMP,
    PRIMARY KEY ( userid, friend_id )
);