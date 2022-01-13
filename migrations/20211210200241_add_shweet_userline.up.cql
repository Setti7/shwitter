CREATE TABLE userline
(
    user_id   UUID,
    shweet_id UUID,
    shweet_message text,
    posted_by UUID,
    created_at timestamp,
    PRIMARY KEY ( user_id, created_at )
) WITH CLUSTERING ORDER BY (created_at DESC);