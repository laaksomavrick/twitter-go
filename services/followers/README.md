## Followers service

The followers service provides functionality for managing user-following and user-follower relationships. Each user follows 0 or more users, subscribing to their tweets for their tweet feed. Correspondingly, each user also has 0 or more followers

#### What is a follower/following in the data model?

```
user_followers (
  username text,
  follower_username text
  PRIMARY KEY (username, follower_username)
);

user_followings (
  username text,
  following_username text
  PRIMARY KEY (username, follower_username)
);
```

#### Access patterns

- Retrieve all the users a user is following from the `user_followers` table
    - E.g, to support a "my followers" list

- Retrieve a count of the number of users a user is following from the `user_followers`
    - E.g, to support a count of the # of followers a user has

- Retrieve all the users following a user from the `user_followings` table
    - E.g, to support a "i'm following" list
    - E.g, for finding all the users subscribed to a particular user's tweets

- Retreive a count of the number of users following a user from the `user_followings` table
    - E.g, to support a count of the # of users a particular user is following