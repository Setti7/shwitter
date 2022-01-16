CREATE TABLE user_followers_count
(
    id UUID,
    count COUNTER,
    PRIMARY KEY ( id )
);

CREATE TABLE user_friends_count
(
    id UUID,
    count COUNTER,
    PRIMARY KEY ( id )
);

CREATE TABLE user_shweets_count
(
    id UUID,
    count COUNTER,
    PRIMARY KEY ( id )
);