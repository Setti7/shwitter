CREATE TABLE user_followers_count
(
    user_id UUID,
    count COUNTER,
    PRIMARY KEY ( user_id )
);

CREATE TABLE user_friends_count
(
    user_id UUID,
    count COUNTER,
    PRIMARY KEY ( user_id )
);

CREATE TABLE user_shweets_count
(
    user_id UUID,
    count COUNTER,
    PRIMARY KEY ( user_id )
);