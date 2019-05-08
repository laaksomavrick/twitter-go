# twitter-go

Twitter-go is an application api intended to back a minimal feature set of twitter. Its development serves as a fun learning exercise to explore an event driven microservice architecture, using Go. I've been curious about Go for a variety of reasons (performance, native type system, binaries, pragmatic ecosystem) and want to explore it and its ecosystem. Moreover, microservice backends are becoming more ubiquitous due to the organizational benefits they offer (independent deployability, independent scalability, fault tolerance), and an event driven architecture is a common way of developing loosely coupled services.

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

- API gateway (Entry point into the backend; maps http to n rpc calls)
- User service (CRUD for users; user authorization)
- Post service (Adding to user post list; "my posts")
- Follower service (Managing user - user follows/followers relationships; "my followers/follower count")
- User feed service (Aggregating user activity feed; "my feed")

#### Service ethos:

- Services are responsible solely for their domain (biz logic, tables)
- Services will publish events about their domain for other services to subscribe to as required
- Services that require data not belonging to their domain will embrace denormalization and eventual consistency
- Services should be written as dumb as possible and avoid pre-emptive abstractions; YAGNI

#### Technical choices:

- Go for all service level code
- RabbitMQ for the message bus
- Cassandra for a nosql datastore
- Docker and Kubernetes for deployment
- Some bash scripts for god knows what
- More things as they come up
