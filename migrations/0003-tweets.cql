DROP TABLE IF EXISTS twtr.tweets;
DROP TABLE IF EXISTS twtr.tweets_by_user;

CREATE TABLE twtr.tweets (
  id uuid PRIMARY KEY,
  username text,
  content text,
  created_at timestamp,
);


CREATE TABLE twtr.tweets_by_user (
  id uuid,
  username text,
  content text,
  created_at timestamp,
  PRIMARY KEY (username, created_at)
);