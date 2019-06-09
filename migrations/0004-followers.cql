DROP TABLE IF EXISTS twtr.user_followers;
DROP TABLE IF EXISTS twtr.user_followings;

CREATE TABLE twtr.user_followers (
  username text,
  follower_username text,
  PRIMARY KEY (username, follower_username)
);

CREATE TABLE twtr.user_followings (
  username text,
  following_username text,
  PRIMARY KEY (username, following_username)
);