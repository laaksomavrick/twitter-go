# twitter-go

Twitter-go is an application api intended to back a minimal feature set of twitter. Its development serves as a fun learning exercise to explore serverless technologies.
A serverless approach boasts an enticing value proposition (no operations work required; pay for what you use; high scalability), hence my curiosity. Moreover, as a primarily js/ts dev, I've been curious about Go for a variety of reasons (performance, native type system, binaries, pragmatic ecosystem) and want to explore it and its ecosystem as well.

#### What is the app?

- A user makes post
- A user has a viewable list of their own posts
- A user can subscribe to other users
- A user has an activity feed of those they follow's posts (chronological order)

#### This app needs to provide the ability:

- To create a user
- To login a user
- To get a list of a user's posts
- To follow other users
- To retrieve an aggregated activity feed of posts from those they follow

#### Ergo, service breakdown:

- User service (CRUD for users; user authorization)
- User profile service (CRUD for user profiles)
- Post service (Adding to user post list; "my posts")
- Follower service (Managing user - user follows/followers relationships)
- User feed service (Aggregating user activity feed; "my feed")

#### Service ethos:

- A logical service is a set of lambda functions and dynamo tables belonging to a particular domain entity
- Each logical service will publish events for other logical services to subscribe to
- Each logical service should have all data it requires to serve requests, ie denormalized data; no rpc or tight coupling
  - Thus, each logical service will be eventually consistent

#### Technical choices:

- AWS Lambda for function handlers with Go (fast, minimal, modern, simple...get shit done with no fuss)
- AWS SQS for an event queue
- AWS DynamoDB for a NoSQL datastore
- AWS API gateway for an api gateway
- Serverless framework for managing the development of all this
- Some bash scripts for god knows what
