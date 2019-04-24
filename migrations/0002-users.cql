DROP TABLE IF EXISTS twtr.users;

CREATE TABLE twtr.users (
  username text PRIMARY KEY,
  email text,
  password text,
  refresh_token text
);