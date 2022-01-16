CREATE TABLE shweet_likes_count
(
    id UUID,
    count COUNTER,
    PRIMARY KEY ( id )
);

CREATE TABLE shweet_reshweets_count
(
    id UUID,
    count COUNTER,
    PRIMARY KEY ( id )
);

CREATE TABLE shweet_comments_count
(
    id UUID,
    count COUNTER,
    PRIMARY KEY ( id )
);

CREATE TABLE user_liked_shweets 
(
    user_id UUID,
    shweet_id UUID,
    PRIMARY KEY ( user_id, shweet_id )
);

CREATE TABLE shweet_liked_by_users
(
    shweet_id UUID,
    user_id UUID,
    PRIMARY KEY ( shweet_id, user_id )
);