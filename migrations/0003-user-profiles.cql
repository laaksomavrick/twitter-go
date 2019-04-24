DROP TABLE IF EXISTS twtr.user_profiles;

CREATE TABLE twtr.user_profiles (
  username text PRIMARY KEY,
  email text,
  display_name text
);