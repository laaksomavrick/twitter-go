DROP TABLE IF EXISTS twtr.tweets;

CREATE TABLE twtr.tweets (
  username text,
  content text,
  created_at timestamp,
  PRIMARY KEY ((username), created_at)
);