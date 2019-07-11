## Feeds service

The feeds service provides functionality to retrieve a user's feed alongside asynchronously aggregating subscribed user tweets to a user's feed. Plainly put: it lets you get your feed, and it handles creating and updating all user feeds. This is known as a _fan-out on write_ strategy, in contrast to a fan-out on read strategy.

For the expected application workload (0, given this is a toy application for learning purposes), either choice would be fine. Fan out on write is typically best suited for users who have a limited number of subscribers (e.g, 100 writes vs 1,000,000 writes on tweet), and makes querying a user's feed blazing fast (grab whatever is in the database at the moment). For a real-world twitter app, they likely use a combination of read and write strategies given the performance heuristics: users with many subscribers have their tweets added to a feed at read with another query, and users with few subscribers have their tweets added to a feed on write.

The application is consciously accepting that there may be a delay between a tweet being made, and that tweet appearing in a user's feed. Moreover, a user's feed will still be viewable even if other services are down (e.g, creating a tweet).

#### What is a feed in the data model?

```
twtr.feed_items (
  username text,
  tweet_id uuid,
  tweet_username text,
  tweet_content text,
  tweet_created_at timestamp,
  PRIMARY KEY (username, tweet_created_at)
);
```


#### Access patterns

- Retrieve a user's feed
- Write a tweet to all relevant feeds