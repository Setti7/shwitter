CREATE TABLE timeline
(
    user_id   UUID,
    shweet_id UUID,
    shweet_message text,
    posted_by UUID,
    created_at timestamp,
    PRIMARY KEY ( user_id, shweet_id )
);