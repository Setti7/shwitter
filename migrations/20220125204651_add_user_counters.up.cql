CREATE TABLE user_counters
(
    id UUID,
    followers COUNTER,
    friends COUNTER,
    shweets COUNTER,
    PRIMARY KEY ( id )
);
