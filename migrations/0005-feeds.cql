DROP TABLE IF EXISTS twtr.feed_items;

CREATE TABLE twtr.feed_items (
  username text,
  tweet_id uuid,
  tweet_username text,
  tweet_content text,
  tweet_created_at timestamp,
  PRIMARY KEY (username, tweet_created_at)
);