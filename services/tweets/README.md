## Tweets service

The tweets service provides functionality for creating and retrieving tweets.

#### What is a tweet in the data model?

```
tweets (
  id uuid PRIMARY KEY,
  username text,
  content text,
  created_at timestamp,
);

tweets_by_user (
  id uuid,
  username text,
  content text,
  created_at timestamp,
  PRIMARY KEY (username, created_at)
);
```

#### Access patterns

- Get a tweet by ID from the `tweets` table
    - E.g, to see a tweet and its comments should this be supported
- Get all tweets by username, sorted by created_at from the `tweets_by_user` table
    - E.g, to support a "my tweets" page for a given user

